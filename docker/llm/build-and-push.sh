#!/usr/bin/env bash
# Build the stash-llm image (UI + Go backend, all from source inside Docker) and optionally
# publish it. Run from anywhere — it cd's to the repo root itself.
#
#   ./docker/llm/build-and-push.sh            # build local image  stash-llm:dev
#   ./docker/llm/build-and-push.sh --push     # build + push       ghcr.io/ryokushen/stash:dev
#   ./docker/llm/build-and-push.sh --save     # build + save tar   stash-llm-dev.tar (for NAS load)
#
# The --save path suits this NAS (no scp/SFTP): pipe the tar over ssh and `docker load` — see
# docs/llm/DEPLOY-NAS.md.
set -euo pipefail

cd "$(dirname "$0")/../.."          # repo root

IMAGE_LOCAL="stash-llm:dev"
IMAGE_REMOTE="ghcr.io/ryokushen/stash:dev"
DOCKERFILE="docker/build/x86_64/Dockerfile"

GITHASH="$(git rev-parse --short HEAD)"
STASH_VERSION="$(git describe --tags --exclude latest_develop 2>/dev/null || echo dev)"

echo ">> building ${IMAGE_LOCAL}  (githash=${GITHASH} version=${STASH_VERSION})"
docker build \
  --build-arg GITHASH="${GITHASH}" \
  --build-arg STASH_VERSION="${STASH_VERSION}" \
  -t "${IMAGE_LOCAL}" \
  -f "${DOCKERFILE}" .

case "${1:-}" in
  --push)
    echo ">> tagging + pushing ${IMAGE_REMOTE}"
    docker tag "${IMAGE_LOCAL}" "${IMAGE_REMOTE}"
    docker push "${IMAGE_REMOTE}"
    ;;
  --save)
    OUT="${2:-stash-llm-dev.tar}"
    echo ">> saving image to ${OUT}"
    docker save "${IMAGE_LOCAL}" -o "${OUT}"
    echo "   load on the NAS with:  docker load < ${OUT}"
    ;;
  "")
    echo ">> done (local image only). Use --push or --save to publish."
    ;;
  *)
    echo "unknown option: ${1}  (use --push or --save)"; exit 2
    ;;
esac
