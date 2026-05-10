#!/bin/sh
set -eu

if [ "$#" -eq 0 ]; then
  set -- stash
elif [ "${1#-}" != "$1" ]; then
  set -- stash "$@"
fi

if [ "$(id -u)" = "0" ]; then
  : "${PUID:=1000}"
  : "${PGID:=1000}"
  : "${HOME:=/root}"
  : "${STASH_CONFIG_FILE:=/root/.stash/config.yml}"
  : "${STASH_CHOWN_PATHS:=$HOME /metadata /cache /blobs /generated}"
  export HOME STASH_CONFIG_FILE

  if ! getent group "$PGID" >/dev/null 2>&1; then
    if command -v addgroup >/dev/null 2>&1; then
      addgroup -g "$PGID" -S stash >/dev/null 2>&1 || addgroup --gid "$PGID" stash >/dev/null 2>&1 || true
    elif command -v groupadd >/dev/null 2>&1; then
      groupadd -g "$PGID" stash >/dev/null 2>&1 || true
    fi
  fi

  group_name="$(getent group "$PGID" 2>/dev/null | cut -d: -f1 || true)"
  group_name="${group_name:-stash}"

  if ! getent passwd "$PUID" >/dev/null 2>&1; then
    nologin=/sbin/nologin
    [ -x "$nologin" ] || nologin=/bin/false

    if command -v adduser >/dev/null 2>&1; then
      adduser -S -D -H -h "$HOME" -s "$nologin" -u "$PUID" -G "$group_name" stash >/dev/null 2>&1 || \
        adduser --system --disabled-password --no-create-home --home "$HOME" --uid "$PUID" --gid "$PGID" --shell "$nologin" stash >/dev/null 2>&1 || true
    elif command -v useradd >/dev/null 2>&1; then
      useradd --system --no-create-home --home-dir "$HOME" --uid "$PUID" --gid "$PGID" --shell "$nologin" stash >/dev/null 2>&1 || true
    fi
  fi

  config_dir="$(dirname "$STASH_CONFIG_FILE")"
  mkdir -p "$HOME" "$config_dir"

  for path in "$HOME" $STASH_CHOWN_PATHS "$config_dir"; do
    [ -n "$path" ] || continue
    mkdir -p "$path" 2>/dev/null || true
    if [ -e "$path" ]; then
      chown -R "$PUID:$PGID" "$path"
    fi
  done

  if command -v su-exec >/dev/null 2>&1; then
    exec su-exec "$PUID:$PGID" "$@"
  fi

  if command -v gosu >/dev/null 2>&1; then
    exec gosu "$PUID:$PGID" "$@"
  fi
fi

exec "$@"
