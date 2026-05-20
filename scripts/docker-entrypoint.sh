#!/bin/sh
set -eu

if [ "$#" -eq 0 ]; then
  set -- stash
elif [ "${1#-}" != "$1" ]; then
  set -- stash "$@"
fi

if [ "$(id -u)" != "0" ]; then
  exec "$@"
fi

: "${PUID:=1000}"
: "${PGID:=1000}"
: "${HOME:=/root}"
: "${STASH_CONFIG_FILE:=/root/.stash/config.yml}"
: "${STASH_CHOWN:=true}"

case "$PUID" in
  '' | *[!0-9]*)
    echo "PUID must be a numeric value." >&2
    exit 1
    ;;
esac

case "$PGID" in
  '' | *[!0-9]*)
    echo "PGID must be a numeric value." >&2
    exit 1
    ;;
esac

lookup_group() {
  command -v getent >/dev/null 2>&1 || return 1
  getent group "$PGID" 2>/dev/null | cut -d: -f1
}

create_group() {
  group_name="$(lookup_group || true)"
  if [ -n "$group_name" ]; then
    echo "$group_name"
    return
  fi

  if command -v addgroup >/dev/null 2>&1; then
    addgroup -S -g "$PGID" stash >/dev/null 2>&1 || \
      addgroup --system --gid "$PGID" stash >/dev/null 2>&1 || true
  elif command -v groupadd >/dev/null 2>&1; then
    groupadd -g "$PGID" stash >/dev/null 2>&1 || true
  fi

  group_name="$(lookup_group || true)"
  if [ -n "$group_name" ]; then
    echo "$group_name"
  else
    echo stash
  fi
}

create_user() {
  group_name="$1"

  if command -v getent >/dev/null 2>&1 && getent passwd "$PUID" >/dev/null 2>&1; then
    return
  fi

  nologin=/sbin/nologin
  [ -x "$nologin" ] || nologin=/bin/false

  if command -v adduser >/dev/null 2>&1; then
    adduser -S -D -H -h "$HOME" -s "$nologin" -u "$PUID" -G "$group_name" stash >/dev/null 2>&1 || \
      adduser --system --disabled-password --no-create-home --home "$HOME" --uid "$PUID" --gid "$PGID" --shell "$nologin" stash >/dev/null 2>&1 || true
  elif command -v useradd >/dev/null 2>&1; then
    useradd --system --no-create-home --home-dir "$HOME" --uid "$PUID" --gid "$PGID" --shell "$nologin" stash >/dev/null 2>&1 || true
  fi
}

drop_privileges() {
  if command -v su-exec >/dev/null 2>&1; then
    exec su-exec "$PUID:$PGID" "$@"
  fi

  if command -v gosu >/dev/null 2>&1; then
    exec gosu "$PUID:$PGID" "$@"
  fi

  echo "No privilege-drop helper found; starting as root." >&2
  exec "$@"
}

group_name="$(create_group)"
create_user "$group_name"

config_dir="$(dirname "$STASH_CONFIG_FILE")"
mkdir -p "$HOME" "$config_dir"

: "${STASH_CHOWN_PATHS:=$HOME $config_dir /metadata /cache /blobs /generated}"
if [ "$STASH_CHOWN" != "false" ]; then
  for path in $STASH_CHOWN_PATHS; do
    [ -n "$path" ] || continue
    mkdir -p "$path" 2>/dev/null || true
    if [ -e "$path" ]; then
      chown -R "$PUID:$PGID" "$path"
    fi
  done
fi

export HOME STASH_CONFIG_FILE
drop_privileges "$@"
