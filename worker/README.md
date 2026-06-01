# stash-worker — external GPU media generator

A Windows + NVIDIA worker that offloads stash's preview/sprite/phash generation from the slow NAS,
writing outputs directly into stash's `generated/` share so stash auto-detects them. Runs outside
stash's job queue (so it works in parallel with whatever stash is doing).

**Design spec:** [`../docs/llm/EXTERNAL-WORKERS.md`](../docs/llm/EXTERNAL-WORKERS.md). This README
is for builders/operators; the design doc explains why.

## Current state — Phase A skeleton

Phase A = previews only. This is a scaffold; the ffmpeg pipeline runs end-to-end but has only been
syntax/type-checked, not validated against real NVENC hardware. Phase B (sprites) and Phase C
(phash + `fileSetFingerprints`) are not yet implemented.

## Build

```powershell
cd worker
go build -o stash-worker.exe ./cmd/stash-worker
```

That produces a single statically-linked Windows executable. No CGO, no DLLs.

## Run

```powershell
.\stash-worker.exe `
  --stash-url http://overwatch-stash:9999 `
  --media-prefix "/data=\\overwatch-stash\torrents" `
  --generated-prefix "/generated=\\overwatch-stash\generated" `
  --ffmpeg .\ffmpeg\bin\ffmpeg.exe `
  --limit 5
```

Flags:

| Flag | What it does |
|---|---|
| `--stash-url` | base URL of your stash instance (over Tailscale) |
| `--stash-api-key` | optional; sent as `ApiKey` header if set (also via `STASH_API_KEY` env) |
| `--media-prefix` `STASH=WORKER` | translate stash's media paths to the worker's view (e.g. SMB share). Multiple flags accepted. |
| `--generated-prefix` `STASH=WORKER` | same, for the `generated/` dir. The worker writes here. |
| `--ffmpeg` | path to `ffmpeg.exe`. Defaults to `ffmpeg` on `$PATH`. Must be a build with NVENC + NVDEC. |
| `--limit` | cap scenes processed per run (0 = unbounded). Useful for first tests. |
| `--watch` | keep running and re-poll instead of exit-when-done. |
| `--dry-run` | print the ffmpeg commands but don't execute or write any files. |

## ffmpeg

Use a pinned dated release from [BtbN/FFmpeg-Builds](https://github.com/BtbN/FFmpeg-Builds) — the
`win64-gpl` **static** build (one `.exe`, no DLL hunt). NVENC + NVDEC are compiled in. Don't use the
rolling `master-latest` tag — pin a date for reproducibility. Download script + checksums will live
under `./ffmpeg/` once Phase A goes past the scaffold.

## What it does, end-to-end

1. Read stash's `configuration.general` (video_file_naming_algorithm + preview knobs).
2. Page through `findScenes`, projecting `id`, `files { id, path, fingerprints { type value } }`,
   and the primary file's hash. **No `is_missing` filter** — that enum doesn't have a "preview"
   value (`pkg/sqlite/scene_filter.go:396-436`). Worker decides what's missing by checking the file
   on disk.
3. For each scene: compute `<generated>/screenshots/<hash>.mp4`; skip if it exists.
4. Otherwise run the NVENC pipeline (per-segment encode → concat demuxer) and atomically rename
   from `<generated>/tmp/` to the final destination.

See [`internal/preview.go`](internal/preview.go) for the exact ffmpeg invocation.

## What it doesn't do (yet)

- Sprites (Phase B)
- Phash via `fileSetFingerprints` (Phase C)
- Concurrent encoding (single-threaded for now; correctness first)
- Two-pass slow-seek retry on first-pass failure (TODO)
- VFR `-vsync 2` detection via ffprobe (TODO)
