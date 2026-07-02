#!/bin/sh

PUID=${PUID:-1000}
PGID=${PGID:-1000}

# Create group if not exists
if ! getent group abc >/dev/null; then
    addgroup -g "${PGID}" abc
fi

# Create user if not exists
if ! getent passwd abc >/dev/null; then
    adduser -u "${PUID}" -G abc -h /app -D -s /bin/sh abc
fi

chown -R abc:abc /app

exec su-exec abc "$@"
