#!/bin/sh
set -e

export HOME=/home/nonroot

if [ "$(id -u)" = "0" ]; then
    if [ -d "/home/nonroot/.swag2mcp" ]; then
        chown -R 65532:65532 /home/nonroot/.swag2mcp 2>/dev/null || true
    fi
    exec su-exec 65532:65532 swag2mcp "$@"
else
    exec swag2mcp "$@"
fi
