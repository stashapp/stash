#!/bin/sh
set -e

# Support PUID/PGID environment variables to run as a specific UID/GID.
# This allows the container to run as a non-root user while ensuring
# proper ownership of mounted volumes when PUID/PGID match the host user.
#
# Usage: PUID=1000 PGID=1000 /entrypoint.sh stash [args...]
# For custom UID/GID, rebuild the image with:
#   docker build --build-arg PUID=<uid> --build-arg PGID=<gid> .

STASH_UID=${PUID:-9999}
STASH_GID=${PGID:-9999}

if [ "$(id -u)" = "0" ]; then
    # Running as root — ensure the stash user exists with the right UID/GID
    # and switch to it using gosu.
    if [ "$STASH_UID" != "0" ]; then
        # Create or update the stash group and user with the configured UID/GID
        addgroup -g ${STASH_GID} stash 2>/dev/null || true
        adduser -u ${STASH_UID} -G stash -h /home/stash -s /bin/sh -D stash 2>/dev/null || true
        # Ensure config directory exists and is owned correctly
        mkdir -p /root/.stash
        chown stash:stash /root/.stash
        exec gosu stash "$@"
    else
        exec "$@"
    fi
elif [ "$(id -u)" = "$STASH_UID" ] && [ "$(id -g)" = "$STASH_GID" ]; then
    # Already running as the correct user — just execute
    exec "$@"
else
    # Running as some other non-root user — execute as-is
    exec "$@"
fi
