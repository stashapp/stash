# facerec.py — last-resort performer tagging by face recognition

When neither a stash-box fingerprint (`external_identify.py`) nor a filename token
(`external_filename_parse.py`) can identify a scene, the only signal left is **who is
on screen**. This tool recognizes performers' faces with **InsightFace** (the
`buffalo_l` model: RetinaFace detector + 512-d ArcFace embedding) on the GPU and
links the matched, **already-existing** performers to those scenes.

It is the lowest-fidelity identifier in the chain — run it **last**, only over the
dark tail (the ~446 scenes still missing performers), and always review a dry-run
before applying.

## Two-step workflow

```
build-embeddings   ──>  performers.npz  ──>   tag-scenes
(once / after new        reference             (frames -> faces ->
 performers+images)      library               cosine match -> link)
```

### Step 1 — build the reference library

Fetches every performer (`findPerformers` -> `id, name, image_path`), downloads each
performer image, detects the face, and stores the L2-normalized 512-d embedding.
Performers with no/blank image are skipped; if an image has several faces the
**largest** (the portrait subject) is used.

```bash
python facerec.py build-embeddings \
    --stash-url http://overwatch-stash:9999 \
    --embeddings-file performers.npz
```

Output is a compressed `.npz` with three arrays: `ids` (object), `names` (object),
`embeddings` (`float32` shape `(N, 512)`, each row L2-normalized). Re-run it whenever
you add performers or change their images.

### Step 2 — tag scenes

For target scenes (default `is_missing:"performers"`), pulls `files{path duration}`,
extracts frames via ffmpeg at fixed fractions of the duration (default
**10/30/50/70/90%**), detects + embeds every face, and matches each against the
library by **cosine similarity**. Confident matches (>= `--min-similarity`) become
proposed performer links.

```bash
# dry-run (default): prints proposed scene -> performers, writes nothing
python facerec.py tag-scenes \
    --stash-url http://overwatch-stash:9999 \
    --embeddings-file performers.npz \
    --media-prefix "/data=\\overwatch-stash\share" \
    --ffmpeg ffmpeg \
    --min-similarity 0.5 --limit 20

# apply once the dry-run looks right (MERGES onto existing performers)
python facerec.py tag-scenes ... --apply
```

Like `external_filename_parse.py`, it **links existing performers only — never
creates** — and writes with **MERGE** semantics: proposed performer ids are added on
top of the scene's existing performers via
`sceneUpdate(input:{id, performer_ids:[...]})`, never replacing them.

## Flags

| flag | default | effect |
|---|---|---|
| `--stash-url` | `http://overwatch-stash:9999` | stash GraphQL base URL |
| `--stash-api-key` | none | sent as `ApiKey` header + `?apikey=` on image GETs |
| `--embeddings-file` | `performers.npz` | reference library path |
| `--ctx-id` | `0` | InsightFace device: `0` = first GPU, `-1` = CPU |
| `--det-size` | `640` | detector input size (square); raise to catch smaller faces |
| `--media-prefix` | none | `STASH_PREFIX=WORKER_PREFIX` path rewrite (tag-scenes) |
| `--ffmpeg` | `ffmpeg` | path to the ffmpeg binary (tag-scenes) |
| `--frames` | `10,30,50,70,90` | % of duration to sample (tag-scenes) |
| `--min-similarity` | `0.5` | cosine threshold for a confident match (tag-scenes) |
| `--max-performers` | `4` | cap performers proposed per scene, highest-confidence first |
| `--scene-filter` | none | raw JSON `SceneFilterType` overriding the default |
| `--limit` | `0` | cap scenes processed (great for a first test) |
| `--apply` | off | write changes (default is dry-run) |
| `--min-image-bytes` | `2048` | (build-embeddings) treat tiny payloads as "no image" |

### `--media-prefix` (path translation)

stash reports media paths from inside its container (e.g. `/data/completed/x.mp4`);
the RTX box reaches the same file over SMB. The rewrite mirrors the Go worker's
`PrefixRewriter`:

```
--media-prefix "/data=\\overwatch-stash\share"
   /data/completed/x.mp4  ->  \\overwatch-stash\share\completed\x.mp4
```

Separators are flipped to backslashes when the worker side is a Windows/UNC path.
(Both bash and PowerShell collapse `\\` to `\`, so a half-escaped `\server\share`
is normalized back to `\\server\share`.)

## GPU setup (RTX box)

1. **NVIDIA driver + CUDA runtime.** `onnxruntime-gpu` needs a matching CUDA +
   cuDNN. For `onnxruntime-gpu` 1.17–1.18 that's the **CUDA 12.x** runtime + cuDNN 9.
   Confirm `nvidia-smi` works. (If you hit a cuDNN/CUDA version mismatch, either
   install the matching cuDNN or pin `onnxruntime-gpu` to a build that matches your
   CUDA — see the onnxruntime CUDA requirements matrix.)
2. **Virtualenv + deps:**
   ```powershell
   python -m venv .venv ; .\.venv\Scripts\activate
   pip install -r requirements.txt
   ```
   Install **only** `onnxruntime-gpu`, never plain `onnxruntime` alongside it.
3. **Model download (automatic, first run).** InsightFace auto-downloads `buffalo_l`
   to `~/.insightface/models/buffalo_l/` on the first `FaceAnalysis(...).prepare()`.
   To pre-stage it offline, drop the `buffalo_l` model files there yourself. It is a
   ~300 MB pack (detector + recognition + a few aux models).
4. **ffmpeg** on `PATH` (or pass `--ffmpeg C:\path\to\ffmpeg.exe`). No GPU build is
   needed — single-frame CPU decode is fine (same reasoning as the worker's cover
   extraction: NVDEC's CUDA-context overhead would outweigh per-frame savings).
5. **Verify the GPU is actually used.** On startup InsightFace prints the active
   providers; you want `CUDAExecutionProvider` listed first. If you only see
   `CPUExecutionProvider`, onnxruntime-gpu didn't load (CUDA/cuDNN mismatch or plain
   `onnxruntime` shadowing it) — it will still run, just slowly.

This repo deliberately ships **no models and no installed deps** — provisioning
happens on the RTX box.

## Accuracy & threshold caveats

- **Face rec is the weakest signal here.** Adult scenes have hard lighting, motion
  blur, oblique/occluded faces, makeup, and lookalikes — expect both misses and
  false positives. Treat every proposal as a suggestion, not ground truth.
- **Threshold.** ArcFace cosine similarity runs roughly: ~0.28 is InsightFace's own
  same/different boundary on clean web photos, but video frames are noisier, so the
  default here is a deliberately **conservative 0.5**. Calibrate on *your* data: run
  `tag-scenes --limit 30` against scenes whose performers you already know, eyeball
  the printed scores, and set `--min-similarity` just above where wrong names start
  appearing. Raising it trades recall for precision.
- **Reference-image quality dominates.** A performer whose stash image is a logo,
  a body shot with no clear face, or a group photo will match poorly or wrongly.
  build-embeddings skips no-face images, but a *bad* single face still gets stored.
- **More frames = more recall, more cost.** The 5 default samples are a balance;
  add frames for stubborn scenes (`--frames 5,15,25,...`), accept the extra ffmpeg +
  inference time.
- **`--max-performers` guards against scene-wide false matches.** A frame full of
  faces can otherwise propose a crowd; the cap keeps the top-N most confident.
- Do a `--dry-run` pass and skim it before every `--apply`. Don't run this against
  scenes that already have correct performers (the default `is_missing:"performers"`
  filter avoids that).

## GraphQL used

- `findPerformers(filter){ performers { id name image_path } }` — build the library.
- Performer image: the `image_path` URL stash returns, falling back to
  `GET <stash>/performer/<id>/image` (API key added as header + `?apikey=`).
- `findScenes(filter, scene_filter){ scenes { id title performers{id} files{path duration} } }`
  — target scenes (default `scene_filter: {is_missing:"performers"}`).
- `sceneUpdate(input:{id, performer_ids:[...]})` — apply (merged) performer links.

## Requirements

See `requirements.txt`: `insightface`, `onnxruntime-gpu`, `numpy<2`, `opencv-python`,
`Pillow`. Python 3.9+.
