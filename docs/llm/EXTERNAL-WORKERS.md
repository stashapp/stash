# External GPU workers for stash media generation

A design for offloading the heavy media-generation work (previews, scrubber sprites, eventually
phashes) from the NAS to a Windows + NVIDIA GPU worker, with outputs landing back in the NAS's stash
`generated/` directory so stash auto-detects them. No code yet — this doc is the spec.

Target operator: 1 user, 1 worker box (Windows, NVIDIA GPU). Designed to scale to multiple workers
later without a rewrite.

---

## 1. Goals & non-goals

### Goals
- Push the slowest stash tasks off the V1500B NAS onto a worker with **NVENC + NVDEC** (10–30× per
  scene for transcode-heavy work).
- Produce **bit-compatible outputs** that stash treats as if its own Generate task had produced them
  (no UI gap, no DB mismatch).
- Run **alongside** stash's own job queue — same parallel model as `identify_scenes_fast`. Stash can
  keep doing its work; the worker is independent.
- Be **resumable** and **idempotent** — interrupting/restarting the worker is safe; rerunning skips
  finished scenes.

### Non-goals
- Replacing covers (single-frame extract, <1s per scene, blob storage — not worth the engineering).
- Multi-worker distributed orchestration (design accommodates it later; the v1 spec doesn't build it).
- Replacing stash's Generate UI/task — the worker is a parallel path, not a replacement.

---

## 2. Architecture

```
┌────────────── NAS (overwatch-stash) ──────────────┐
│  /volume2/data/torrents/{completed,incoming}        │  ← media files (SMB share, RO from worker)
│  /volume1/docker/stash/generated/                   │  ← stash's generated/ dir (SMB share, RW)
│  /volume1/docker/stash/config/config.yml            │  ← knobs the worker needs to read
│  http://overwatch-stash:9999/graphql                │  ← scene list, "refresh" trigger
└─────────────────────────────┬───────────────────────┘
                              │  SMB (read media + write generated)
                              │  HTTPS (GraphQL)
                              ▼
┌──────────────── Worker (Windows + NVIDIA) ───────────────┐
│  worker.exe (or python worker.py):                       │
│   • polls stash for "scenes missing X"                   │
│   • runs ffmpeg with NVENC/NVDEC                          │
│   • writes output to the NAS share at the exact path    │
│   • tells stash to refresh (lightweight GraphQL call)   │
└──────────────────────────────────────────────────────────┘
```

**Key design choice — mount-based, not sync-based.** The worker writes directly to the NAS's
`generated/` share. No local scratch, no rsync afterwards. Simpler, instant visibility in stash. The
tradeoff is SMB write latency for many small files; in practice OK for the file sizes here
(sprite ~500KB, preview ~3–10MB, VTT a few KB).

---

## 3. Output file layout (verbatim from `pkg/models/paths/paths_scenes.go`)

Stash's `generated/` directory has fixed subdirectories. The worker MUST write to these exact paths:

| Artifact | Path (relative to `generated/`) | Format |
|---|---|---|
| Video preview | `screenshots/<checksum>.mp4` | H.264 MP4 |
| Animated WebP preview | `screenshots/<checksum>.webp` | WebP |
| Scrubber sprite image | `vtt/<checksum>_sprite.jpg` | JPEG montage |
| Scrubber sprite VTT | `vtt/<checksum>_thumbs.vtt` | WebVTT cue file |
| (Phash) | *no file output — DB column* | uint64 |

`<checksum>` is determined by stash's `video_file_naming_algorithm` config setting. The enum values
are **uppercase** (`pkg/models/model_file.go:14-24`):
- `OSHASH` (default) — stash's own oshash hash of the file
- `MD5` — MD5 sum

**The worker must read that config setting at startup and use the matching hash for every scene.**
Mixing causes orphaned files stash won't find.

For multi-file scenes (real in stash), the **primary file's** checksum is used for preview/sprite
paths — see `pkg/models/model_scene.go:257-266`. Phase A/B handle this naturally because they key on
the scene's `Checksum`/`OSHash` field which already resolves to the primary file. Phase C (phash)
must iterate **all** files of a scene because phashes are stored per file_id, not per scene.

---

## 4. Stash config knobs the worker must mirror

Stash's Generate task is highly configurable. The worker has to **read these from stash** (via
GraphQL `configuration { general { ... } }`) and feed them into ffmpeg or it produces non-matching
output. Defaults from `internal/manager/config/config.go`:

### For previews:
| Key | Default | Effect |
|---|---|---|
| `previewSegments` | 12 | how many short clips to extract per video |
| `previewSegmentDuration` | 0.75s | length of each clip |
| `previewExcludeStart` | "0" | seconds (or %) to skip at start |
| `previewExcludeEnd` | "0" | seconds (or %) to skip at end |
| `previewAudio` | true | include audio |
| `previewPreset` | "slow" | x264 preset (ignored when using NVENC; pick comparable NVENC preset) |

### For sprites:
| Key | Default | Effect |
|---|---|---|
| `use_custom_sprite_interval` | false | use custom interval vs auto |
| `sprite_interval` | 30s | seconds between sprite frames (when custom) |
| `minimum_sprites` | 10 | minimum frame count |
| `maximum_sprites` | 500 | maximum frame count |
| `sprite_screenshot_width` | 160 | width of each thumbnail in the sprite |

Auto-mode (default): stash computes the count such that `count ≈ duration / 30`, clamped to
`[min, max]`. The worker mirrors this formula.

### For phash (Phase C):
Fixed in `pkg/hash/videophash/phash.go`:
- **25 screenshots** in a 5 × 5 grid (`columns = rows = 5`).
- Each frame extracted at **width = 160 px, height preserves source aspect ratio** (the screenshot
  options set only `Width`; ffmpeg scales the other dimension proportionally). A 16:9 source →
  160×90 frames → 800×450 sprite. A 4:3 source → 160×120 → 800×600. **NOT 160×160.**
- Sampling timestamps: stash trims **5% off each end of the video** before sampling. From the source
  (`phash.go:84-85`):
  ```
  offset   = 0.05 * duration
  stepSize = 0.9  * duration / 25     // = duration * 0.036
  t_i      = offset + i * stepSize    for i in 0..24
  ```
- `goimagehash.PerceptionHash(sprite)` → uint64.

**This is the spec for the phash sprite — different from the scrubber sprite. The worker must use
this exact frame layout and sample timing for hash compatibility. Even a 1-pixel shift in framing
changes the perceptual hash.**

---

## 5. ffmpeg pipelines (NVIDIA hardware-accelerated)

### Tooling
- **Windows ffmpeg with full NVENC + NVDEC support.** Use a **pinned dated release** from
  [BtbN/FFmpeg-Builds](https://github.com/BtbN/FFmpeg-Builds) (NOT `master-latest` — that's a rolling
  tag and reproducibility breaks). Prefer the **`win64-gpl`** static build over `gpl-shared` — one
  `ffmpeg.exe`, no DLL hunt, simpler for end-users.
- Pin the build + SHA256 in the worker repo; a tiny download script grabs it on first build.
- NVENC encoder names: `h264_nvenc`, `hevc_nvenc`. NVDEC decoder: `-hwaccel cuda -hwaccel_output_format cuda`.
- The NVIDIA driver provides the runtime; no CUDA SDK install needed on the worker.

### Preview (the heavy one)
Reference: `pkg/scene/generate/preview.go`. Stash's actual params (from lines 182-190, 215):

| Param | Value | Note |
|---|---|---|
| codec | `libx264` | replaced with `h264_nvenc` |
| **crf** | **`21`** | NVENC analog: `-rc vbr -cq 21` (NOT 23 — the original doc draft had this wrong) |
| pix_fmt | `yuv420p` | required for browser compatibility |
| profile | `high` | required |
| level | `4.2` | required |
| preset | `slow` | NVENC `-preset p6` is the slow-quality analog (`p4` is medium = faster but lower quality) |
| scale | `-vf scale=640:-1` | preview width is fixed to **640** (`scenePreviewWidth` const) |
| audio | `-c:a aac -b:a 128k -ac 2` | applied when `previewAudio: true` |
| movflags | `+faststart` | for web playback |

**Exclude-window parsing** (`preview.go:38-48`): `previewExcludeStart`/`previewExcludeEnd` accept
plain seconds (e.g. `"30"`) **or a percent suffix** (e.g. `"5%"` = 5% of duration). The worker must
parse both. Effective segment timestamps:
```
total    = duration - excludeStartSec - excludeEndSec
stepSize = total / Segments
t_i      = excludeStartSec + i * stepSize    for i in 0..(Segments-1)
```

**Use stash's two-step approach (NVENC encode per-segment → concat demuxer), not single-pass
filter_complex.** Stash's preview.go calls `transcoder.Splice` which uses ffmpeg's `concat` *demuxer*
(byte-level stream copy of pre-encoded segments) — this gives near-identical output to stash. A
single-pass `filter_complex concat` re-encodes through filters and produces visually different
output, and breaks NVDEC pipelining because trim filters force `hwdownload` per segment unless
explicitly placed.

Per-segment encode:
```
ffmpeg -hwaccel cuda -hwaccel_output_format cuda \
  -ss <t_i> -i <input> -t 0.75 \
  -vf "hwdownload,format=nv12,scale=640:-1" \
  -c:v h264_nvenc -preset p6 -profile:v high -level 4.2 -rc vbr -cq 21 -pix_fmt yuv420p \
  -c:a aac -b:a 128k -ac 2 \
  seg_<i>.mp4
```

Concat (no re-encode):
```
ffmpeg -f concat -safe 0 -i seglist.txt -c copy -movflags +faststart <checksum>.mp4
```

**Two-pass retry, mirroring stash** (`task_generate_preview.go:67-72`): the first attempt uses
`-xerror` (fail on warnings) with fast-seek (`-ss` before `-i`). If it fails — common with corrupt
or VFR sources — retry once with **slow-seek** (`-ss` after `-i`) and no `-xerror`. Additionally, if
ffprobe reports framerate ≤ 0.01 fps (VFR/probe failure), add `-vsync 2` to the encode. The single-
shot pipeline in the original draft of this doc skipped all of this and would crash on the first
weird file.

**Expected speedup:** 30s+ on the NAS CPU → ~5s on NVENC for a typical scene (the per-segment +
concat approach is slightly slower than the single-pass theoretical max, but still ~6× faster than
the NAS and produces drop-in-compatible output).

### WebP preview
Same segment-extraction logic, encode with `-c:v libwebp -lossless 0 -compression_level 4` (NVENC
doesn't encode WebP — this stays CPU, but it's much cheaper than video transcode).

### Scrubber sprite
1. Compute frame timestamps: `count = clamp(duration / 30, 10, 500)`; `interval = duration / count`.
2. Extract each frame at width 160, source aspect preserved:
   ```
   ffmpeg -hwaccel cuda -ss <t> -i <input> -vframes 1 -vf scale=160:-1 frame_<i>.jpg
   ```
3. Tile into one big JPEG. **Grid dimensions follow `gridSize = ceil(sqrt(N))`** — from
   `pkg/scene/generate/sprite.go:113-115`. For 25 frames, gridSize = 5 → 5×5. For 12 frames,
   gridSize = 4 → 4×4 with the last 4 cells blank (VTT only emits N cues).
4. Generate the **VTT file**. Cell dimensions are derived from the sprite image, not hard-coded:
   ```
   cell_w = sprite_image_width  / gridSize         // both ints; sprite_image_width = 160 * gridSize
   cell_h = sprite_image_height / gridSize         // height per cell = source-aspect-preserved
   ```
   So for a 16:9 source with 25 frames: cell_w=160, cell_h=90 (sprite 800×450). For 4:3 it changes.
   VTT format:
   ```vtt
   WEBVTT

   00:00:00.000 --> 00:00:30.000
   <checksum>_sprite.jpg#xywh=0,0,<cell_w>,<cell_h>

   00:00:30.000 --> 00:01:00.000
   <checksum>_sprite.jpg#xywh=<cell_w>,0,<cell_w>,<cell_h>
   ...
   ```
   The exact grid math + VTT generation lives in `pkg/scene/generate/sprite.go` — read that for
   any edge cases (very short videos, etc.).

### Phash sprite (Phase C)
5×5 grid, 25 frames, frames at width=160 with source aspect preserved (NOT 160×160). Sampling
formula in §4 above (`offset=0.05·duration`, `step=0.9·duration/25`). The sprite image is built but
**not written to disk** — the worker hashes it in-memory with `goimagehash.PerceptionHash` and
applies the result via the existing `fileSetFingerprints` mutation (no fork addition needed — see
Phase C below).

---

## 6. Phase plan (the build order)

### Phase A — preview worker
Smallest useful version: one task type, validate the whole loop end-to-end.

1. Worker reads stash's `video_file_naming_algorithm` + preview config (via GraphQL `configuration { general }`).
2. Worker SMB-mounts NAS shares (media RO, generated/ RW). Confirms the mounts resolve before starting.
3. **Enumerate candidate scenes by filesystem walk, NOT by GraphQL "missing" filter.** Stash's
   `is_missing` enum does *not* include `videoPreview` (see `pkg/sqlite/scene_filter.go:396-436`) —
   stash determines preview existence purely by `fsutil.FileExists` at task runtime
   (`task_generate_preview.go:101`). So:
   - Page through `findScenes` with no missing-filter, batches of e.g. 500.
   - For each scene, read its primary file's hash and check `<generated>/screenshots/<hash>.mp4`.
   - Skip if the file exists and is non-empty; otherwise add to the worker's queue.
4. For each queued scene:
   - Run the per-segment NVENC pipeline from §5, with the two-pass retry + VFR detection.
   - Write each segment + the concat output to **`<generated>/tmp/<hash>.mp4.partial`** (stash
     provides this scratch dir — `paths_generated.go:34`). After fsync, `os.Rename` to the final
     destination at `<generated>/screenshots/<hash>.mp4`. Renaming from a sibling dir means stash
     never sees a half-written file via path-existence checks.
   - On rename failure (destination exists from a concurrent run), prefer skip-and-leave-existing.
5. **No `metadataScan` call needed** — preview detection is lazy. Stash serves the file as soon as
   it appears on disk. (The original doc draft suggested calling `metadataScan` here; that's for
   discovering new media files, it does nothing for generated-artifact detection.)

Verify: open the scene in stash; preview plays. No new code on the stash side.

### Phase B — sprite worker
Same worker process, add a second task. Most of the work is the **VTT format**: it must match
stash's expected grid layout exactly. The sprite image format is forgiving (any JPEG with the right
dimensions); the VTT must reference its cells with the right `#xywh=` regions.

### Phase C — phash worker (no fork addition needed)
This unblocks `identify_scenes_fast` for the long tail.

**Use the existing `fileSetFingerprints` mutation.** Phash in stash is stored **per file**
(`files_fingerprints` table, keyed by `file_id + type='phash'`), NOT per scene. The mutation
already exists at `internal/api/resolver_mutation_file.go:283` — schema:
```graphql
fileSetFingerprints(input: FileSetFingerprintsInput!): Boolean!
# input: { id: FileID, fingerprints: [{ type: "phash", value: "<16-hex>" }] }
```
Stash's resolver parses the hex string into a uint64 (`strconv.ParseUint(value, 16, 64)`) and
writes the row. **No new mutation needed.** The original doc draft proposed adding `sceneSetPhash`;
that was the wrong shape (scenes can have multiple files; phash is per-file).

#### Worker flow
1. Pull scenes missing phash. The filter `scene_filter: {phash: {modifier: IS_NULL, value: ""}}` is
   accepted — verified in `pkg/sqlite/scene_filter.go:422-425`. Include `files { id fingerprints { type value } }`
   in the projection so the worker can decide which files need a phash.
2. For each scene, iterate **all files** (multi-file scenes are real). For each file missing phash:
   a. Generate the 5×5 phash sprite per the §4 spec — width=160 aspect-preserved, sampling at
      `offset+i·step` per the formula. Decode with NVDEC, scale on the GPU where possible.
   b. Tile in-memory (`image/draw` in Go; no file written).
   c. Hash with `goimagehash.PerceptionHash`. This guarantees bit-compatibility with stash's
      native phash — same library, same algorithm.
3. Call `fileSetFingerprints` with the file's ID and the 16-char hex hash.

#### Verification (do this before bulk-running)
Generate phash for ~20 files that **already have a native phash from stash**. Compare hex output.
If even one differs, the sprite generation isn't matching stash's params — fix before bulk-running
or you'll write wrong hashes to the DB. Likely culprits: frame timestamps off (5%/95% window),
scale dimensions wrong, color space mismatch, JPEG vs BMP intermediate (stash uses BMP via
`transcoder.ScreenshotOutputTypeBMP`).

---

## 7. Worker tech stack

Recommendation: **Go + ffmpeg binary + Python NOT required**.

Why Go over Python:
- Phash compatibility wants goimagehash → already Go.
- Single statically-linked exe on Windows; simpler operationally than Python + venv.
- We're already a Go-fork dev environment.
- The stash-box GraphQL types are already in the repo if useful.

Worker layout:
```
worker/
  cmd/stash-worker/main.go      ← entrypoint + CLI flags
  internal/worker/
    config.go                    ← reads stash config + worker config
    scenes.go                    ← GraphQL: list scenes missing X
    preview.go                   ← NVENC ffmpeg invocation for previews
    sprite.go                    ← sprite + VTT generation
    phash.go                     ← phash sprite + goimagehash
    refresh.go                   ← notify stash
  third-party/ffmpeg/            ← BtbN Windows ffmpeg+NVENC binary, pinned
  README.md
```

Build target: `GOOS=windows GOARCH=amd64 go build` → single `.exe`. Ships with the ffmpeg dir.
Run as: `stash-worker.exe --stash-url http://overwatch-stash:9999 --media-root \\overwatch-stash\media --generated-root \\overwatch-stash\generated --tasks previews,sprites --parallel 2`.

---

## 8. Concurrency on the worker

The old "consumer GPUs cap NVENC at 2-3 sessions" is **no longer true** — NVIDIA removed the session
cap entirely in driver R555 (mid-2024). On RTX 30/40/50, parallelism is limited by VRAM and encoder
throughput, not session count.

- Default `--parallel 2` for previews — a *responsiveness* knob, not a hardware limit. Keep the host
  PC usable; bump higher on the actual hardware if benching shows headroom.
- Sprite/phash generation is decode-heavy not encode-heavy; can run higher concurrency safely.
- GPU memory: 4–8GB is plenty for 4K transcodes at this concurrency.

---

## 9. Operations

### Run model
A **long-running Windows service** is overkill for v1. Just a CLI you start when you want it to work,
with sane interrupt handling (Ctrl-C stops cleanly, partial files are cleaned up). **Exit
condition:** worker exits when two consecutive enumeration passes (60s apart) return zero candidate
scenes. A `--watch` flag keeps it running indefinitely; default is exit-when-done.

### Logging
Stdout for the operator; a rotating file log in `%LOCALAPPDATA%\stash-worker\`. Per-scene success/
failure log so you can post-mortem.

### Resume / atomicity
Idempotency comes from file-existence check. Crashes mid-scene leave `.partial` files in
`<generated>/tmp/` — stash never sees them (different directory than `screenshots/` or `vtt/`).
Worker startup sweeps `tmp/` to clean orphans. Windows SMB → Synology Samba: `os.Rename` over SMB
is delete-then-rename, not POSIX-atomic; the millisecond window is acceptable because both source
and destination live in stash-owned directories and stash treats destination-absence the same as
"needs generating" (lazy detection). SMB 3.0+ on the Synology side handles persistent handles
correctly.

### Error policy
- ffmpeg failure on one scene → log + skip; don't kill the whole run.
- SMB write failure → exponential backoff retry on that scene (transient network blips).
- GraphQL failure → retry the call, then give up after N attempts; print summary at the end.

### Auth
Stash localhost on the NAS doesn't currently require auth. If you add an API key later, the worker
sends it via the **`ApiKey` header** (stash's convention — NOT `Authorization: Bearer`), driven by a
`STASH_API_KEY` env var (same pattern as `external_identify.py`).

---

## 10. Open questions / risks

| Question | Risk | Mitigation |
|---|---|---|
| ~~Does stash detect new preview files lazily?~~ | — | **Confirmed lazy** — `task_generate_preview.go:101` does `fsutil.FileExists` at task time. File drop is sufficient; no refresh call needed. Closed. |
| Visual difference: NVENC + CQ vs libx264 + CRF | Side-by-side preview look different even at matched quality target | Document as expected. Use `-cq 21` to approximate stash's `-crf 21`. Not a correctness issue. |
| `measurements` / tattoos string formats (only relevant to performer enrichment, not this doc) | — | See `EXTERNAL-PERFORMER-IDENTIFY.md`. Out of scope here. |
| SMB write performance for many small files (sprite + VTT) | Slow on poor SMB tuning | Bulk-write strategy: hold sprite + VTT until both done, then write both within ~100ms. Single MP4 writes are big enough that SMB handles them fine. |
| Stash's `video_file_naming_algorithm` could change mid-run | Existing files become orphaned with the wrong hash; new files diverge | Worker re-reads on startup; warns if it changes between runs. |
| Multiple workers writing the same file (future) | Race conditions | v1 is single-worker; v2 adds a claim/lease via a stash tag or a separate sqlite-on-NAS for work queue. |
| Per-segment + concat demuxer requires segment files on disk | Disk IO + cleanup | Write segments to `<generated>/tmp/<hash>/seg_*.mp4`; delete after concat. Tmp dir is stash's own scratch — won't interfere. |
| Phash sprite parameters mismatch stash's exact output | Hashes don't match → bad writes to file fingerprints | Verification gate before bulk-running (§ Phase C). 20-scene cross-check is mandatory. |

---

## 11. What to build first

Smallest valuable shippable: **Phase A only**, preview worker, ~3–5 days of work:

1. ffmpeg-NVENC per-segment + concat pipeline working locally on Windows against one scene file
   (validate output plays in stash side-by-side with a stash-generated preview).
2. GraphQL client in Go for `findScenes` (paginated) + config queries.
3. End-to-end loop on one scene: NAS share → per-segment encode → concat → atomic rename → stash
   serves it.
4. Scale to all scenes; concurrency; idempotency (existence check); two-pass retry; VFR detection.
5. Polish: logs, summary, signal handling.

**Phase B** (sprites) is +1-2 days — same architecture, adds the VTT grid math + sqrt-grid tiling.
**Phase C** (phash) is **~3-4 days** (down from the original "week" — no fork mutation needed; use
the existing `fileSetFingerprints`). Most of Phase C's time is the **verification gate**: prove the
worker's phash output matches stash's bit-for-bit against ~20 known files. If they don't match,
fixing the sprite-gen params can eat days.

---

## 12. Approval gate

Before any code: review this doc, flag anything wrong, and pick the scope. Once approved, I'll
implement Phase A end-to-end on a branch, you test against one scene, then we extend.
