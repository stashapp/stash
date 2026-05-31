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
  'mkdir -p /volume1/docker/stash/{config,metadata,cache,blobs,generated} \
            /volume1/docker/stash/cliproxy/auth'
```
Copy the compose + gateway/bridge configs into place (via tar, since no scp):
```bash
cd /b/Downloads2/Projects/stash-llm/docker/llm
tar cf - docker-compose.nas.yml litellm/config.yaml cliproxy/config.example.yaml \
  | ssh -p 3239 shadowshark@overwatch-stash \
    'mkdir -p /volume1/docker/stash/compose && tar xf - -C /volume1/docker/stash/compose'
# on the NAS: cp cliproxy/config.example.yaml cliproxy/config.yaml and set a random api-key
```
Then create the env files in `/volume1/docker/stash/compose/` from their `.env.example` templates
(**do not commit them**):
- `stash.env` — `STASH_ASSISTANT_API_KEY` = the LiteLLM master key (must match `litellm.env`), and
  `STASH_ASSISTANT_MODEL` (`minimax` or `grok`). Base URL is preset in the compose file.
- `litellm.env` — `LITELLM_MASTER_KEY` (random), `MINIMAX_API_KEY`, and for the `grok` model:
  `GROK_PROXY_URL=http://cli-proxy-api:8317/v1` + `GROK_PROXY_KEY` = the api-key you put in
  `cliproxy/config.yaml`.
- `cliproxy/config.yaml` — set one random string under `api-keys` (this is `GROK_PROXY_KEY`).

### 3a. One-time Grok OAuth login (only for the `grok` model)
The `cli-proxy-api` bridge needs your **grok.com OAuth token**. The login is an interactive
browser flow, so do it on a machine with a browser, then copy the token to the NAS:
```bash
# on your laptop (has a browser) — run the bridge's Grok login (see help.router-for.me for the
# exact flag; it writes the token into ~/.cli-proxy-api):
docker run --rm -it -p 8317:8317 -v "$HOME/.cli-proxy-api:/root/.cli-proxy-api" \
  eceasy/cli-proxy-api:latest <grok-login-command>
# then ship the resulting token dir to the NAS auth-dir (no scp → tar over ssh):
tar cf - -C "$HOME/.cli-proxy-api" . | ssh -p 3239 shadowshark@overwatch-stash \
  'tar xf - -C /volume1/docker/stash/cliproxy/auth'
```
> **MiniMax needs none of this** — `STASH_ASSISTANT_MODEL=minimax` works with just `MINIMAX_API_KEY`.
> Use `grok` only if you want the subscription path (gray-area vs xAI ToS; OAuth tokens refresh and
> can expire — re-run the login if Grok calls start failing). If you skip Grok, comment out the
> `grok` model in `litellm/config.yaml` and drop the `cli-proxy-api` service from the compose file.

## 4. Bring it up
Synology Container Manager supports compose projects; via CLI (verify the compose binary name on the
box — `docker compose` plugin vs `docker-compose`):
```bash
ssh -p 3239 shadowshark@overwatch-stash \
  'cd /volume1/docker/stash/compose && /usr/local/bin/docker compose -f docker-compose.nas.yml up -d'
```
**Compose is the recommended path** — it brings up `stash`, `litellm`, and `cli-proxy-api` on a private
network so stash reaches the gateway at `http://litellm:4000/v1` and litellm reaches the bridge at
`http://cli-proxy-api:8317/v1`. (If you're only using `minimax`, the `cli-proxy-api` service is
unused and can be removed.) The plain `docker run` fallback below starts **only stash**; if you use it
you must run litellm (and, for grok, cli-proxy-api) yourself on a shared network and set
`STASH_ASSISTANT_BASE_URL` accordingly.

Fallback plain `docker run` (stash only — see note above):
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
point the library at `/data` and let it scan. The assistant is configured via the gateway env (or rely on
the env var) and open the Assistant panel.

## Gotchas
- **Config volume is mandatory.** Without `…/config:/root/.stash`, stash dies with
  `could not write to provided config path` — the dir must exist (step 3 handles it).
- **Media stays read-only.** The assistant's write tools mutate *metadata in stash's DB*, never the
  media files. The `:ro` mounts enforce that the files Jellyfin serves are never touched.
- **Rebuilds drop nothing** — all state is on the bind mounts; `docker compose up -d` with a new image
  is a clean upgrade. Keep the previous image tag around for rollback.
