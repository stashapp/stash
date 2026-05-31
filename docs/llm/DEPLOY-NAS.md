# stash-llm — NAS deployment

Target: **`overwatch-stash`** (Synology DSM, x86_64), as a Docker container alongside Jellyfin, reusing
Jellyfin's media folders **read-only**. Files: [`docker/llm/docker-compose.nas.yml`](../../docker/llm/docker-compose.nas.yml),
[`.env.example`](../../docker/llm/.env.example).

## NAS facts (don't relearn these the hard way)
- SSH: `ssh -p 3239 shadowshark@overwatch-stash` — **custom port 3239**.
- `docker` is at **`/usr/local/bin/docker`** (not on the non-interactive SSH PATH — use the full path).
- `shadowshark` is **not a sudoer** but is in the `docker` group → no `sudo`; `docker exec` runs as root.
- **No scp/SFTP** — move bytes with `tar | ssh` or `docker save | ssh … docker load`.
- Media (read-only, already used by Jellyfin): `/volume2/data/torrents/{completed,incoming}`.
- App state goes under `/volume1/docker/stash/` (Jellyfin uses `/volume1/docker/jellyfin/`).
- Port `9999` (Jellyfin is `8096`, Mission Control `5119`, agents `9119/9120` — no clash).

## 1. Build the image (on the Mac/dev box)
```bash
cd /b/Downloads2/Projects/stash-llm
./docker/llm/build-and-push.sh            # local image stash-llm:dev
```

## 2. Get the image onto the NAS — pick one

**A. Registry (ghcr)** — convenient, but the ghcr package must be pullable by the NAS:
```bash
./docker/llm/build-and-push.sh --push      # → ghcr.io/ryokushen/stash:dev
# on the NAS (only if the package is private):
#   echo <GH_PAT> | /usr/local/bin/docker login ghcr.io -u ryokushen --password-stdin
/usr/local/bin/docker pull ghcr.io/ryokushen/stash:dev
```

**B. Save + load over SSH** — no registry, no scp, works with the NAS's constraints (recommended first time):
```bash
./docker/llm/build-and-push.sh --save                       # → stash-llm-dev.tar
docker save stash-llm:dev | ssh -p 3239 shadowshark@overwatch-stash '/usr/local/bin/docker load'
```
If you use path **B**, set `image: stash-llm:dev` in the compose file instead of the ghcr ref.

## 3. Prepare state dirs + secrets on the NAS
```bash
ssh -p 3239 shadowshark@overwatch-stash \
  'mkdir -p /volume1/docker/stash/{config,metadata,cache,blobs,generated}'
```
Copy the compose + env into place (via tar, since no scp):
```bash
cd /b/Downloads2/Projects/stash-llm/docker/llm
tar cf - docker-compose.nas.yml | ssh -p 3239 shadowshark@overwatch-stash \
  'mkdir -p /volume1/docker/stash/compose && tar xf - -C /volume1/docker/stash/compose'
```
Then create `/volume1/docker/stash/compose/stash.env` on the NAS from `.env.example`, filling
`STASH_ANTHROPIC_API_KEY` from your secret store. **Do not commit it.**

## 4. Bring it up
Synology Container Manager supports compose projects; via CLI (verify the compose binary name on the
box — `docker compose` plugin vs `docker-compose`):
```bash
ssh -p 3239 shadowshark@overwatch-stash \
  'cd /volume1/docker/stash/compose && /usr/local/bin/docker compose -f docker-compose.nas.yml up -d'
```
Fallback plain `docker run` (if compose isn't available):
```bash
ssh -p 3239 shadowshark@overwatch-stash '/usr/local/bin/docker run -d --name stash-llm \
  --restart unless-stopped -p 9999:9999 \
  --env-file /volume1/docker/stash/compose/stash.env \
  -e STASH_STASH=/data/ -e STASH_GENERATED=/generated/ -e STASH_METADATA=/metadata/ -e STASH_CACHE=/cache/ \
  -v /volume1/docker/stash/config:/root/.stash \
  -v /volume1/docker/stash/metadata:/metadata -v /volume1/docker/stash/cache:/cache \
  -v /volume1/docker/stash/blobs:/blobs -v /volume1/docker/stash/generated:/generated \
  -v /volume2/data/torrents/completed:/data/completed:ro \
  -v /volume2/data/torrents/incoming:/data/incoming:ro \
  stash-llm:dev'
```

## 5. Verify
```bash
ssh -p 3239 shadowshark@overwatch-stash '/usr/local/bin/docker logs --tail 20 stash-llm'
# expect: "stash is listening on 0.0.0.0:9999" and the correct version line
```
Then browse to `http://overwatch-stash:9999/` (over Tailscale). First run shows Stash's setup wizard;
point the library at `/data` and let it scan. Configure the Anthropic key under Settings (or rely on
the env var) and open the Assistant panel.

## Gotchas
- **Config volume is mandatory.** Without `…/config:/root/.stash`, stash dies with
  `could not write to provided config path` — the dir must exist (step 3 handles it).
- **Media stays read-only.** The assistant's write tools mutate *metadata in stash's DB*, never the
  media files. The `:ro` mounts enforce that the files Jellyfin serves are never touched.
- **Rebuilds drop nothing** — all state is on the bind mounts; `docker compose up -d` with a new image
  is a clean upgrade. Keep the previous image tag around for rollback.
