#!/bin/sh
set -eu

if [ "$(id -u)" -eq 0 ]; then
  exec /usr/local/bin/nvidia-patch-entrypoint.sh "$@"
fi

exec "$@"
