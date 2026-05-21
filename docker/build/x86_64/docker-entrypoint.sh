#!/bin/sh
set -eu

# If the container is started as a non-root user, honor that directly.
if [ "$(id -u)" -ne 0 ]; then
  exec "$@"
fi

: "${PUID:=1000}"
: "${PGID:=1000}"
: "${UMASK:=002}"
: "${STASH_USER:=stash}"
: "${STASH_GROUP:=stash}"
: "${STASH_CONFIG_FILE:=/config/config.yml}"

export HOME=/config
export STASH_CONFIG_FILE

umask "${UMASK}" 2>/dev/null || true

mkdir -p /config /cache /metadata /generated /blobs

existing_group="$(awk -F: -v gid="${PGID}" '$3 == gid { print $1; exit }' /etc/group)"
if [ -n "${existing_group}" ]; then
  STASH_GROUP="${existing_group}"
elif ! grep -q "^${STASH_GROUP}:" /etc/group; then
  addgroup -g "${PGID}" -S "${STASH_GROUP}"
fi

existing_user="$(awk -F: -v uid="${PUID}" '$3 == uid { print $1; exit }' /etc/passwd)"
if [ -n "${existing_user}" ]; then
  STASH_USER="${existing_user}"
elif ! id -u "${STASH_USER}" >/dev/null 2>&1; then
  adduser -u "${PUID}" -S -D -h /config -G "${STASH_GROUP}" "${STASH_USER}"
fi

addgroup "${STASH_USER}" "${STASH_GROUP}" >/dev/null 2>&1 || true

chown -R "${PUID}:${PGID}" /config /cache /metadata /generated /blobs

exec su-exec "${PUID}:${PGID}" "$@"
