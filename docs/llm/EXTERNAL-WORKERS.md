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

`<checksum>` is determined by stash's `video_file_naming_algorithm` config setting:
- `oshash` (default) — stash's own oshash hash of the file
- `md5` — MD5 sum

**The worker must read that config setting at startup and use the matching hash for every scene.**
Mixing causes orphaned files stash won't find.

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
| `useCustomSpriteInterval` | false | use custom interval vs auto |
| `spriteInterval` | 30s | seconds between sprite frames (when custom) |
| `minimumSprites` | 10 | minimum frame count |
| `maximumSprites` | 500 | maximum frame count |
| `spriteScreenshotSize` | 160 | width of each thumbnail in the sprite |

Auto-mode (default): stash computes the count such that `count ≈ duration / 30`, clamped to
`[min, max]`. The worker mirrors this formula.

### For phash (Phase C):
Fixed in `pkg/hash/videophash/phash.go`:
- 5 × 5 = **25 screenshots** evenly distributed across the video
- Each **160 × 160 px**, tiled into an 800 × 800 sprite
- `goimagehash.PerceptionHash(sprite)` → uint64

**This is the spec for the phash sprite — different from the scrubber sprite. The worker must use
this exact layout for hash compatibility.**

---

## 5. ffmpeg pipelines (NVIDIA hardware-accelerated)

### Tooling
- **Windows ffmpeg with full NVENC + NVDEC support.** Use the [BtbN/FFmpeg-Builds](https://github.com/BtbN/FFmpeg-Builds)
  Windows builds — `ffmpeg-master-latest-win64-gpl-shared` includes NVENC, NVDEC, libwebp, all the
  needed filters. Ship it alongside the worker (single zipped folder) so versioning is pinned.
- NVENC encoder names: `h264_nvenc`, `hevc_nvenc`. NVDEC decoder: pass `-hwaccel cuda -hwaccel_output_format cuda`.
- The NVIDIA driver provides the runtime; no CUDA SDK install needed on the worker.

### Preview (the heavy one)
Reference: `pkg/scene/generate/preview.go`. Stash's CPU pipeline is approximately:

```
ffmpeg -ss <start> -i <input> -t <segDuration> -map 0:v:0 -an -c:v libx264 -preset slow ... segN.mp4
ffmpeg -i seg0.mp4 -i seg1.mp4 ... -filter_complex concat=n=12:v=1:a=0 output.mp4
```

NVIDIA-accelerated version:

```
# fast single-pass: select N segments via the trim+concat filter, encode once with NVENC
ffmpeg -hwaccel cuda -hwaccel_output_format cuda -i <input> \
  -filter_complex "[0:v]trim=start=t0:end=t0+0.75,setpts=PTS-STARTPTS[v0]; \
                   [0:v]trim=start=t1:end=t1+0.75,setpts=PTS-STARTPTS[v1]; \
                   ... \
                   [v0][v1]...[v11]concat=n=12:v=1[outv]" \
  -map "[outv]" -c:v h264_nvenc -preset p4 -tune hq -rc constqp -qp 23 \
  -movflags +faststart screenshots/<checksum>.mp4
```

Audio is added with a parallel filter chain when `previewAudio: true`. Expected: 30s+ on CPU → ~2-4s
on NVENC for a typical scene.

### WebP preview
Same segment-extraction logic, encode with `-c:v libwebp -lossless 0 -compression_level 4` (NVENC
doesn't encode WebP — this stays CPU, but it's much cheaper than video transcode).

### Scrubber sprite
1. Compute frame timestamps: `count = clamp(duration / 30, 10, 500)`; `interval = duration / count`.
2. Extract each frame at its timestamp (NVDEC for decode → CPU for scaling/save, since the
   intermediate is a JPEG):
   ```
   ffmpeg -hwaccel cuda -ss <t> -i <input> -vframes 1 -vf scale=160:-1 frame.jpg
   ```
3. Tile into a single image with **ImageMagick `montage`** or `ffmpeg -filter_complex tile=...`.
4. Generate the **VTT file** by emitting one cue per frame:
   ```vtt
   WEBVTT

   00:00:00.000 --> 00:00:30.000
   <checksum>_sprite.jpg#xywh=0,0,160,90

   00:00:30.000 --> 00:01:00.000
   <checksum>_sprite.jpg#xywh=160,0,160,90
   ...
   ```
   The exact grid (columns × rows) and per-frame dimensions are in the VTT — they have to match the
   sprite image. Stash's layout convention is in `pkg/scene/generate/sprite.go`.

### Phash sprite (Phase C)
Fixed 5×5 grid, 25 frames evenly spaced, each scaled to 160×160. Same pattern as the scrubber
sprite but with hard-coded dimensions and **no VTT**. The image gets hashed (not stored as a file).

---

## 6. Phase plan (the build order)

### Phase A — preview worker
Smallest useful version: one task type, validate the whole loop end-to-end.

1. Worker reads stash's `video_file_naming_algorithm` + preview config.
2. Worker SMB-mounts NAS shares (media RO, generated/ RW).
3. Polls stash for scenes whose preview is missing (GraphQL: `findScenes` with filter
   `is_missing: "videoPreview"` — confirm exact field name from schema).
4. For each scene:
   - Compute target path: `<generated>/screenshots/<checksum>.mp4`.
   - Skip if file already exists and is non-empty (idempotency).
   - Run the NVENC ffmpeg pipeline.
   - Write atomically (write to `.partial`, then rename) to avoid stash seeing half-written files.
5. After every N scenes, call stash's `metadataScan` (with generate flags) to make sure the player
   picks up the new files — or just rely on stash's lazy detection.

Verify: open the scene in stash; preview plays. No new code on the stash side.

### Phase B — sprite worker
Same worker process, add a second task. Most of the work is the **VTT format**: it must match
stash's expected grid layout exactly. The sprite image format is forgiving (any JPEG with the right
dimensions); the VTT must reference its cells with the right `#xywh=` regions.

### Phase C — phash worker + new mutation
This unblocks `identify_scenes_fast` for the long tail.

#### The fork addition (small)
Add a GraphQL mutation to the fork: `sceneSetPhash(scene_id: ID!, phash: String!): Boolean`.
- `phash` is the 16-character hex string of a uint64 (how stash stores it).
- Implementation: lookup scene → set the `phash` column via `pkg/sqlite` → return.
- ~30 lines in `internal/api/resolver_mutation_scene.go` + schema declaration. The same low-touch
  pattern as the LLM additions.

#### The worker
1. Pull scenes missing phash (`scene_filter: {phash: {modifier: IS_NULL, value: ""}}`).
2. For each scene, generate the **5×5 phash sprite** with NVDEC + ffmpeg + Pillow (or ffmpeg `tile=`).
3. **Compute the hash with goimagehash** (the Go library, to guarantee bit-compatible output with
   stash's native phash). Options:
   - Compile a tiny Go helper binary (~20 lines) that the worker invokes per sprite, OR
   - Make the worker itself Go from the start (probably simpler in the end).
4. Call the new `sceneSetPhash` mutation.

Phash sensitivity to image differences is real — even a 1-pixel shift in framing changes the hash.
Need a verification step: generate phash for ~20 known scenes, compare to stash's native phashes for
the same scenes, ensure they match. If they don't, debug the sprite-generation params until they do.

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

NVENC has hardware concurrency limits — consumer GPUs (RTX 30/40 etc.) can run **2–3 encoding
sessions** simultaneously. The worker:
- Default `--parallel 2` for previews (NVENC-bound).
- Sprite generation is decode-heavy not encode-heavy; can run higher concurrency safely.
- Phash sprite gen similar to scrubber sprite.

GPU memory: 4–8GB is plenty for 4K transcodes.

---

## 9. Operations

### Run model
A **long-running Windows service** is overkill for v1. Just a CLI you start when you want it to work,
with sane interrupt handling (Ctrl-C stops cleanly, partial files are cleaned up).

### Logging
Stdout for the operator; a rotating file log in `%LOCALAPPDATA%\stash-worker\`. Per-scene success/
failure log so you can post-mortem.

### Resume
Idempotency comes from file-existence check. If the worker crashes mid-scene, the partial file is
named `.partial` and either gets cleaned on next start or finished. Stash never sees a half-written
preview.

### Error policy
- ffmpeg failure on one scene → log + skip; don't kill the whole run.
- SMB write failure → exponential backoff retry on that scene (transient network blips).
- GraphQL failure → retry the call, then give up after N attempts; print summary at the end.

### Auth
Stash localhost on the NAS doesn't currently require auth. If you add an API key later, the worker
honors a `STASH_API_KEY` env var (same pattern as `external_identify.py`).

---

## 10. Open questions / risks

| Question | Risk | Mitigation |
|---|---|---|
| Does stash detect new preview files lazily, or does it cache "missing" state? | Files might not appear until next Scan/Generate. | If lazy, no problem. If cached: call `metadataGenerate` with all flags false (a cheap no-op that refreshes detection), or just include a periodic minimal `metadataScan` of metadata only. |
| `goimagehash.PerceptionHash` vs Python `imagehash.phash` compatibility | Hashes don't match → stash-box doesn't recognize our phashes. | Pick Go worker for phash specifically (uses the same library). Cross-verify against 20 known scenes before bulk-running. |
| ffmpeg NVENC output vs stash's libx264 output for previews | Bitrate/quality differences. They're playable either way, but visually different. | Document the visual difference; acceptable. Match `-rc constqp -qp 23` approximately to stash's `-crf 23`. Not a correctness issue. |
| SMB write performance for many small files (sprite + VTT) | Slow on poor SMB tuning. | Bulk-write strategy: hold sprite + VTT until both done, then write both within ~100ms. Single MP4 writes are big enough that SMB handles them fine. |
| Stash's `videoFileNamingAlgorithm` could change mid-run | Existing files become orphaned with the wrong hash; new files diverge. | Worker re-reads on startup; warns if it changes between runs. |
| Multiple workers writing the same file (future) | Race conditions. | v1 is single-worker; v2 adds a claim/lease via a stash tag or a separate sqlite-on-NAS for work queue. |

---

## 11. What to build first

Smallest valuable shippable: **Phase A only**, preview worker, ~3–5 days of work spread:

1. ffmpeg-NVENC pipeline working locally on Windows against one scene file (validate output plays in stash).
2. GraphQL client in Go for `findScenes` + config queries.
3. End-to-end loop on one scene end-to-end (NAS share, write back, stash recognizes).
4. Scale to all scenes; concurrency; idempotency; error handling.
5. Polish: logs, summary, signal handling.

Phase B builds on this in a few more days. Phase C is the biggest — it includes the fork mutation +
phash verification work — likely a week.

---

## 12. Approval gate

Before any code: review this doc, flag anything wrong, and pick the scope. Once approved, I'll
implement Phase A end-to-end on a branch, you test against one scene, then we extend.
