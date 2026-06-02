# stash-llm — operations runbook (every tool, every command)

One page for all the external/off-queue tooling. Per-tool detail lives in each tool's
README (linked per section); this is the copy-paste command reference + the operational
rules (run order, what's safe in parallel) that otherwise only live in people's heads.

Everything here is **additive to upstream stash** and talks to it over its GraphQL API or
the SMB shares — nothing here is a stash core change.

---

## 0. Setup & shared facts

| Thing | Value |
|---|---|
| stash API (over Tailscale) | `http://overwatch-stash:9999` |
| stash API (when running ON the NAS) | `http://localhost:9999` |
| auth | none currently — `--stash-api-key` is optional (needed only if auth gets enabled) |
| media path rewrite | `--media-prefix '/data=\\overwatch-stash\data\torrents'` |
| generated path rewrite | `--generated-prefix '/generated=\\overwatch-stash\docker\stash\generated'` |
| GPU worker box | Windows + RTX 5080, tailnet `ryokushen-throne`; repo at `B:\Downloads2\Projects\stash-llm` |
| ffmpeg (worker) | `./ffmpeg-master-latest-win64-gpl/bin/ffmpeg.exe` (BtbN win64-gpl, NVENC+NVDEC), from `worker/` |
| NAS shell | `ssh -p 3239 shadowshark@overwatch-stash` (custom port; `docker` = `/usr/local/bin/docker`; not a sudoer; no scp/SFTP — use `tar`/`rsync`) |
| NAS tool copies | `/volume1/docker/stash/tools/` and `~/stash-tools/`; NAS python is `python3` (3.8, stdlib only) |

**Windows/bash gotcha:** prefix worker and facerec commands with `MSYS_NO_PATHCONV=1` so the
`/data=...` and `\\server\share` args aren't path-mangled. The Python tools all send a browser
`User-Agent` (stash returns an empty body to the default `Python-urllib` UA).

**Everything defaults to safe:** the worker skips already-generated files; every metadata-writing
Python tool defaults to **dry-run** (needs `--apply`); `db_backup` needs `--apply`.

---

## 1. GPU worker — media generation (run on the RTX box, from `worker/`)

`stash-worker.exe`. Doc: [`worker/README.md`](../../worker/README.md) · design: [`EXTERNAL-WORKERS.md`](EXTERNAL-WORKERS.md).
Tasks (in `--tasks`, comma-sep, in order): `previews, covers, sprites, phash, transcode, image-phash`.

```bash
cd worker
# previews + covers + sprites + scene phash (the common bulk set)
MSYS_NO_PATHCONV=1 ./stash-worker.exe \
  --stash-url http://overwatch-stash:9999 \
  --media-prefix '/data=\\overwatch-stash\data\torrents' \
  --generated-prefix '/generated=\\overwatch-stash\docker\stash\generated' \
  --ffmpeg ./ffmpeg-master-latest-win64-gpl/bin/ffmpeg.exe \
  --tasks previews,covers,sprites,phash --max-failures 30
```

| Task | What | `--generated-prefix`? |
|---|---|---|
| `previews` | 12-seg NVENC preview mp4 → `screenshots/<hash>.mp4` | required |
| `covers` | cover JPEG via `sceneUpdate(cover_image)` | not needed (API) |
| `sprites` | scrubber sprite + VTT → `vtt/` | required |
| `phash` | per-file scene phash via `fileSetFingerprints` (bit-exact replica) | not needed (API) |
| `transcode` | browser-friendly h264/aac → `transcodes/<hash>.mp4` (HEVC/mkv/ac3/…); h264 sources remuxed | required |
| `image-phash` | per-image phash (enumerates all images, skips phashed) → image dedup | not needed (API) |

Other flags: `--limit N` (cap encoded items), `--max-failures N`, `--per-page N`, `--watch`, `--dry-run`.

**Phash bit-exactness gate (always run once before a bulk phash if ffmpeg changed):**
```bash
MSYS_NO_PATHCONV=1 ./stash-worker.exe --stash-url http://overwatch-stash:9999 \
  --media-prefix '/data=\\overwatch-stash\data\torrents' \
  --ffmpeg ./ffmpeg-master-latest-win64-gpl/bin/ffmpeg.exe --verify-phash 25
```

**Run a 2nd instance for a parallel light task** (e.g. covers) alongside a heavy one — same .exe, different `--tasks`.

---

## 2. Off-queue identify / metadata (Python, run anywhere with API access)

Doc: [`tools/identify/README.md`](../../tools/identify/README.md). All default to dry-run; add `--apply` to write.

```bash
# (a) fingerprint identify — studio/performers/date/tags + stash_id (StashDB + TPDB)
python3 tools/identify/external_identify.py --stash-url http://overwatch-stash:9999            # dry-run
python3 tools/identify/external_identify.py --stash-url http://overwatch-stash:9999 --apply
#   flags: --all-scenes  --phashed-only  --set-organized  --no-tags  --allow-multiple  --batch N  --limit N

# (b) filename parser — link to EXISTING studios/performers from the path (never creates)
python3 tools/identify/external_filename_parse.py --stash-url http://overwatch-stash:9999      # dry-run
python3 tools/identify/external_filename_parse.py --stash-url http://overwatch-stash:9999 --apply
#   flags: --all-scenes  --set-title  --min-name-len N  --limit N

# (c) performer enrichment — fill performer metadata from stash-box (refresh/search/both)
python3 tools/identify/external_identify_performers.py --stash-url http://overwatch-stash:9999 --mode both        # dry-run
python3 tools/identify/external_identify_performers.py --stash-url http://overwatch-stash:9999 --mode both --apply
#   flags: --mode {refresh,search,both}  --allow-multiple  --per-page N  --limit N

# (d) orchestrator — runs (a) then (b) in the safe order, with gap snapshots
python3 tools/maintenance/ingest_pipeline.py --stash-url http://overwatch-stash:9999            # dry-run both
python3 tools/maintenance/ingest_pipeline.py --stash-url http://overwatch-stash:9999 --apply
#   flags: --steps identify,filename  --stop-on-error  --no-snapshot
```

---

## 3. Library maintenance (Python). Doc: [`tools/maintenance/README.md`](../../tools/maintenance/README.md)

```bash
# tag consolidation — auto-merge SAFE dup tags; report subjective families
python3 tools/maintenance/tag_consolidate.py --stash-url http://overwatch-stash:9999            # dry-run report
python3 tools/maintenance/tag_consolidate.py --stash-url http://overwatch-stash:9999 --apply     # merge SAFE
python3 tools/maintenance/tag_consolidate.py --stash-url http://overwatch-stash:9999 --emit-plan # JSON plan of subjective merges

# scene dedup — phash Hamming clustering of near-duplicate scenes
python3 tools/maintenance/scene_dedup.py --stash-url http://overwatch-stash:9999                 # report
python3 tools/maintenance/scene_dedup.py --stash-url http://overwatch-stash:9999 --apply          # tag non-keepers _dupe-candidate
#   flags: --max-distance N (default 8)  --limit N

# generated/ orphan + integrity sweeper — RUN ON THE NAS (local walk)
ssh -p 3239 shadowshark@overwatch-stash 'python3 /volume1/docker/stash/tools/generated_sweeper.py \
  --stash-url http://localhost:9999 --generated-dir /volume1/docker/stash/generated'             # report
#   add --prune to move orphans into generated/_quarantine/ (reversible)

# metadata DB + config backup — RUN ON THE NAS (dry-run by default; --apply to write)
ssh -p 3239 shadowshark@overwatch-stash 'python3 /volume1/docker/stash/tools/db_backup.py \
  --db /volume1/docker/stash/config/stash-go.sqlite --config /volume1/docker/stash/config/config.yml \
  --out /volume1/docker/stash/backups --apply'
#   flags: --keep N (default 14)  --remote user@host:/path (off-NAS rsync)
#   SCHEDULE: DSM Control Panel → Task Scheduler → user-defined script (no crontab as shadowshark), nightly --apply.
```

---

## 4. Face recognition — dark-tail performer tagging (RTX box, `tools/facerec/`)

Doc: [`tools/facerec/README.md`](../../tools/facerec/README.md). Uses InsightFace; runs in the `.venv`.
**GPU note:** the CUDA EP loads but a cuDNN-frontend failure on the Blackwell (sm_120) card forces
CPU fallback — so pass **`--ctx-id -1`** (explicit CPU) for now. CPU is fine here (build is network-bound,
tagging is a low-frequency batch). ffmpeg must be an **absolute** path (facerec runs from `tools/facerec/`).

```bash
cd tools/facerec
# (1) build the per-performer embedding library (one-time; ~985 performers w/ images)
.venv/Scripts/python.exe facerec.py build-embeddings \
  --stash-url http://overwatch-stash:9999 --ctx-id -1 --embeddings-file performer_embeddings.npz

# (2) tag the dark tail (is_missing:performers). DRY-RUN FIRST, then --apply.
MSYS_NO_PATHCONV=1 .venv/Scripts/python.exe -u facerec.py tag-scenes \
  --stash-url http://overwatch-stash:9999 \
  --media-prefix '/data=\\overwatch-stash\data\torrents' \
  --ffmpeg 'B:/Downloads2/Projects/stash-llm/worker/ffmpeg-master-latest-win64-gpl/bin/ffmpeg.exe' \
  --ctx-id -1 --embeddings-file performer_embeddings.npz \
  --min-similarity 0.45 --max-performers 4 --limit 10            # dry-run; add --apply to write
#   --min-similarity 0.45 is calibrated (≥0.44 = true positives; false positives sat at 0.21–0.29)
#   links EXISTING performers only (never creates). --scene-filter '<raw JSON>' to retarget.
```

---

## 5. In-app assistant (no CLI — driven from the Stash web UI's Assistant panel)

Grok/MiniMax via the LiteLLM gateway. It has the hand-coded library/scraper/identify tools **plus**
generic GraphQL (`graphql_schema`/`graphql_query`/`graphql_mutate`), persisted self-defined tools
(`define_tool`/`list_dynamic_tools`/`delete_dynamic_tool`), and a confirm-gated shell (`run_command`,
gated behind `assistant_dev_loop_enabled`). All mutations are confirm-gated (`write_policy=ask`). Detail:
[`DESIGN.md`](DESIGN.md) §4/§4a/§6. `identify_scenes_fast` = the bundled `external_identify.py` run from inside the container.

---

## 6. Run order & what's safe in parallel

**Pipeline order for a fresh/new batch of scenes:**
1. stash **Scan** (native) → 2. worker **covers/previews/sprites/phash** → 3. worker **transcode** (after phash) →
4. **identify** (fingerprint) → 5. **filename-parse** (link-only, the leftovers) → 6. **facerec** (the dark tail) →
7. **tag_consolidate** / **scene_dedup** cleanup. `ingest_pipeline.py` does steps 4–5 in order.

**Hard rules:**
- **Never** run `external_identify` and `external_filename_parse` at the same time (both write studio/performers/date — they'd race). Identify first.
- **Never** run `external_identify --apply` alongside a *native* stash Identify (last-write-wins).

**Concurrency (the real cap is NAS disk I/O):**
- **Heavy (read media)** — `previews, sprites, transcode, phash, image-phash`, `facerec tag-scenes`: only **~2 at once** before the NAS SHR disk saturates.
- **Light (API/DB only, no media reads)** — `external_identify`, `external_filename_parse`, `external_identify_performers`, `tag_consolidate`, `scene_dedup`, `generated_sweeper`, `db_backup`, `ingest_pipeline`: stack freely alongside the heavy ones (just minor SQLite write-lock contention).
- A `docker compose up -d stash` redeploy is a ~10-15s API blip: the **worker tasks** and `generated_sweeper` ride through it (retry/resumable); the standalone identify/parse/tag scripts are not hardened — just rerun (all idempotent/resumable).
