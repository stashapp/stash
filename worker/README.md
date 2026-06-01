# stash-worker — external GPU media generator

A Windows + NVIDIA worker that offloads stash's preview / cover / sprite generation from the slow NAS,
writing outputs directly into stash's `generated/` share so stash auto-detects them. Runs outside
stash's job queue (so it works in parallel with whatever stash is doing).

**Design spec:** [`../docs/llm/EXTERNAL-WORKERS.md`](../docs/llm/EXTERNAL-WORKERS.md). This README
is for builders/operators; the design doc explains why.

## Current state — live-validated 2026-06-01 on RTX 5080

Five task types are shipped and live-tested against the production NAS:

| Task | Output | Per-scene time | Notes |
|---|---|---|---|
| **previews** | `<generated>/screenshots/<hash>.mp4` | ~7s | 12 NVENC-encoded segments + concat demuxer (matches stash's `transcoder.Splice`). Two-pass slow-seek retry + VFR `-vsync 2` detection both implemented. |
| **covers** | base64 JPEG via `sceneUpdate(cover_image: …)` mutation | ~3s | Single-frame extract at 20% of duration (mirrors `screenshotDurationProportion` in stash). CPU JPEG (no NVENC overhead). |
| **sprites** | `<generated>/vtt/<hash>_sprite.jpg` + `_thumbs.vtt` | ~30s | Per-frame `-ss` seek with NVDEC, tiled in Go via `image/draw`. WebVTT cues match stash's `gridSize = ⌈√N⌉` math. |
| **phash** | `phash` fingerprint via `fileSetFingerprints` mutation | ~15-25s | **Bit-for-bit** replica of stash's `pkg/hash/videophash`: 25 BMP frames (5×5) at 5%-offset timestamps, `imaging` montage, `goimagehash.PerceptionHash`. **CPU decode (no NVDEC)** — hardware decode would change pixels and break the hash. Gated by `--verify-phash` (see below). |
| **transcode** | `<generated>/transcodes/<hash>.mp4` | remux ~fast / NVENC varies | Pre-generates a browser-friendly h264/aac MP4 for scenes stash can't stream directly (HEVC, mpeg4/wmv, ac3/dts audio, mkv/avi/wmv containers); stash then serves it directly (`GetStreamPath`) instead of live-transcoding on the weak NAS. If the video is already h264 it **stream-copies** (lossless, just fixes audio/container); otherwise full **NVENC h264** re-encode at full resolution (CPU decode → universal codec support). HEVC is treated as needs-transcode even though stash considers it streamable. Heavy on NAS I/O — run after previews/phash, in `--limit` batches. |

**What's NOT shipped:**
- Concurrent multi-scene encoding (sequential per worker; run two `.exe` processes in parallel as a poor-man's concurrency — proven to work, see "Parallel workers" below).

## Build

The worker has no CGO and no native deps. You can either build natively (needs Go installed locally)
or cross-build via a Docker container (no Go install needed):

```bash
# Native (if you have Go)
cd worker
go build -ldflags="-s -w" -o stash-worker.exe ./cmd/stash-worker

# Cross-build via Docker container (works on any host with Docker Desktop)
MSYS_NO_PATHCONV=1 docker run --rm \
  -v "$(pwd)":/worker -w /worker \
  golang:1.25-alpine \
  sh -c 'GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o stash-worker.exe ./cmd/stash-worker'
```

Either path produces a single statically-linked Windows executable (~6.3MB).

## Run

```powershell
.\stash-worker.exe `
  --stash-url http://overwatch-stash:9999 `
  --media-prefix "/data=\\overwatch-stash\data\torrents" `
  --generated-prefix "/generated=\\overwatch-stash\docker\stash\generated" `
  --ffmpeg .\ffmpeg-master-latest-win64-gpl\bin\ffmpeg.exe `
  --tasks previews,covers `
  --limit 5
```

### Flags

| Flag | What it does | Default |
|---|---|---|
| `--stash-url` | base URL of your stash instance (over Tailscale) | `http://localhost:9999` |
| `--stash-api-key` | optional; sent as `ApiKey` header (also via `STASH_API_KEY` env) | empty |
| `--media-prefix` `STASH=WORKER` | translate stash's media paths to the worker's view (e.g. SMB share). UNC paths with one or two leading backslashes are both accepted — see "UNC path quoting" below. | empty |
| `--generated-prefix` `STASH=WORKER` | same, for the `generated/` dir. The worker writes here. **Required for `previews`/`sprites`** (they write files); optional for `covers`/`phash` (those write via the API). | empty |
| `--ffmpeg` | path to `ffmpeg.exe`. Must be a build with NVENC + NVDEC. | `ffmpeg` (PATH) |
| `--tasks` | comma-separated, in order: `previews`, `covers`, `sprites`, `phash`, `transcode` | `previews` |
| `--verify-phash N` | **gate, not a task.** Recompute `N` files that already have a native stash phash and compare. Exits non-zero on any mismatch. Run this once before trusting `--tasks phash`. | `0` (off) |
| `--limit` | cap items ENCODED per run (skips don't count). `--limit 1` encodes exactly one missing item then exits. | `0` (unbounded) |
| `--max-failures` | abort after N **consecutive** failures (catches systemic problems vs thrashing the whole library) | `5` |
| `--per-page` | GraphQL pagination size for scene enumeration | `200` |
| `--watch` | keep polling instead of exit-when-done | off |
| `--dry-run` | print what would happen without writing | off |

### Path translation

Stash sees mounted media at paths like `/data/completed/foo.mp4`. The worker accesses the same file
over SMB at `\\overwatch-stash\data\torrents\completed\foo.mp4`. The `--media-prefix` flag tells the
worker how to rewrite:

- `/data=\\overwatch-stash\data\torrents` — stash sees `/data`; worker sees `\\overwatch-stash\data\torrents`
- `/generated=\\overwatch-stash\docker\stash\generated` — same idea for outputs

#### UNC path quoting

Both `bash` and PowerShell collapse `\\` to `\` before args reach a binary. The worker has a
`normalizeUNC` step that detects a mangled `\server\share` form and restores the missing leading
backslash — so passing `\\` in your shell quoting layer **works regardless** of which shell escapes
it. Drive-letter paths (`C:\`) and already-correct UNC (`\\server\share`) are left alone.

### Idle / exit

The worker exits cleanly when two consecutive enumeration passes return zero candidates 60s apart.
`--watch` keeps it running indefinitely. `Ctrl-C` once = graceful cancel (finish current scene),
`Ctrl-C` twice = force exit.

## ffmpeg

Pin a dated release from [BtbN/FFmpeg-Builds](https://github.com/BtbN/FFmpeg-Builds). Use the
**`win64-gpl`** static build — one `ffmpeg.exe`, no DLLs to hunt down. NVENC + NVDEC + libwebp are
compiled in. Verified working: `ffmpeg-master-latest-win64-gpl` (download from the Latest release
page).

## What each task does, end-to-end

### previews
1. Read stash's `configuration.general` (video_file_naming_algorithm + preview knobs).
2. Page through `findScenes`, projecting `id`, `files { id, path, fingerprints { type value } }`.
   **No `is_missing` filter** — stash doesn't track preview existence (`pkg/sqlite/scene_filter.go:396-436`).
3. For each scene: compute `<generated>/screenshots/<hash>.mp4`; skip if it exists.
4. Otherwise probe with ffprobe (duration + framerate), encode each segment via NVENC fast-seek
   (`-ss` before `-i`, `-xerror`). On per-segment failure, retry once with slow-seek. If framerate
   ≤ 0.01 fps, add `-vsync 2`. Concat via demuxer.
5. Write to `<generated>/tmp/<hash>.mp4.partial`, then atomically rename to the final destination.

### covers
1. Page through `findScenes` with `scene_filter: {is_missing: "cover"}` — server-side filter.
2. For each scene: ffmpeg single-frame extract at `0.2 × duration` (matches stash's
   `screenshotDurationProportion`).
3. Base64-encode the JPEG, send via `sceneUpdate(input:{id, cover_image: "data:image/jpeg;base64,…"})`.
   Stash unpacks and writes to its blob storage.

### phash

Replicates `pkg/hash/videophash/phash.go` exactly — the output **must** be a bit-identical
`uint64` or it won't match StashDB/TPDB fingerprints and the work is wasted.

1. Page through `findScenes` with `scene_filter: {is_missing: "phash"}` — server-side filter that
   **shrinks** as phashes are applied (same re-fetch-page-1 pagination as covers).
2. For each **file** lacking a phash (phash is per-file, not per-scene): use the **stash-stored
   duration** (from the API — NOT a fresh ffprobe; a few-ms drift would shift every sampled
   timestamp and change the hash).
3. Extract 25 frames (5×5 grid) at `offset + i·step` where `offset = 0.05·dur`, `step = 0.9·dur/25`.
   Each frame: `ffmpeg -v error -y -ss <t> -i <src> -frames:v 1 -vf scale=160:-2 -c:v bmp -f rawvideo -`
   (lossless BMP, **CPU decode** — exactly stash's `transcoder.ScreenshotTime`).
4. Montage the 25 frames with `disintegration/imaging` (same lib/version as stash), perceptual-hash
   with `corona10/goimagehash` (same), take `GetHash()`.
5. Send `fileSetFingerprints(input:{id, fingerprints:[{type:"phash", value:<hex>}]})`. Value is
   `strconv.FormatUint(hash, 16)` — stash parses it back with `ParseUint(_, 16, 64)`, so the
   round-trip matches stash's own.

**Always run the gate first:** `--verify-phash 25` recomputes 25 files that already have a native
phash and compares. Bit-exactness depends on the ffmpeg build's decoder/scaler producing identical
pixels; if your ffmpeg differs from the NAS's, the gate catches it before you write 700 wrong hashes.

### sprites
1. Page through `findScenes` (no filter — stash doesn't track sprite existence).
2. For each scene: skip if both `<hash>_sprite.jpg` and `<hash>_thumbs.vtt` already exist.
3. Compute `N = clamp(duration/30, 10, 500)` frame count and `gridSize = ⌈√N⌉`.
4. Run **4 parallel ffmpeg seeks** (bounded), each extracting one scaled JPEG via NVDEC fast-seek.
5. Tile the JPEGs in Go via `image/draw` into one composite JPEG.
6. Write a matching WebVTT cue file (cell dimensions derived from the decoded frames).
7. Atomic rename for both files.

## Parallel workers

Running **two `.exe` instances** concurrently is supported and effective:

```powershell
# Window 1: previews + cover sweep
.\stash-worker.exe --tasks previews,covers ...

# Window 2: cover-only catch-up (faster covers visible in stash UI immediately)
.\stash-worker-dev.exe --tasks covers ...
```

Each task type touches different output paths and stash doesn't lock at the resource level, so
concurrent workers do not conflict on writes. **However**, the bottleneck is usually the **NAS
disk** (Synology SHR + 4-core V1500B), not your GPU. Watch `iostat` on the NAS:

```bash
ssh shadowshark@overwatch-stash 'iostat -x 1 2 | tail -8'
```

If `%util` is >85% you're at the disk's limit; more workers will just slow each other down.

## Known limitations

- **Sprite gen is the slowest task** (~30s/scene even after optimization). 1417 scenes = ~12h
  uncontested. The bottleneck is NAS read I/O across 4 parallel ffmpeg seeks per scene.
- **phash is bit-exactness-sensitive.** It depends on your ffmpeg build decoding + scaling to the
  same pixels the NAS's ffmpeg does. The `--verify-phash` gate exists precisely because this can't
  be assumed — always run it after changing ffmpeg builds. CPU decode is mandatory (no NVDEC).
- **Corrupt H.264 sources** (Invalid NAL unit, etc.) fail both fast-seek AND slow-seek passes. Seen
  in ~1% of scenes in testing. They show up as `[stash-worker] scene N: encode failed: ...` in the
  log and don't block the rest of the run (unless you hit `--max-failures` consecutively).
- **Windows SMB rename** isn't POSIX-atomic; brief 1-2 ms window where the destination doesn't
  exist. Stash's lazy detection treats absent files as "needs generating" so this is harmless in
  practice.
