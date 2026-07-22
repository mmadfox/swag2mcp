#!/usr/bin/env bash
set -euo pipefail

REPO="mmadfox/swag2mcp"
VERSION=""

while getopts "v:" opt; do
  case "$opt" in
    v) VERSION="$OPTARG" ;;
    *) echo "Usage: $0 [-v version]" >&2; exit 1 ;;
  esac
done

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH" >&2; exit 1 ;;
esac

if [ -z "$VERSION" ]; then
  VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | sed 's/.*"tag_name": "\(.*\)",/\1/')
  if [ -z "$VERSION" ]; then
    echo "Failed to fetch latest version" >&2
    exit 1
  fi
fi

FILENAME="swag2mcp_${VERSION}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/download/$VERSION/$FILENAME"

echo "Downloading swag2mcp $VERSION ($OS/$ARCH)..."
curl -fsSL "$URL" -o "/tmp/$FILENAME"

echo "Extracting..."
tar xzf "/tmp/$FILENAME" -C /tmp

INSTALL_DIR="/usr/local/bin"
if [ ! -w "$INSTALL_DIR" ]; then
  INSTALL_DIR="$HOME/.local/bin"
  mkdir -p "$INSTALL_DIR"
fi

mv /tmp/swag2mcp "$INSTALL_DIR/swag2mcp"
chmod +x "$INSTALL_DIR/swag2mcp"
rm -f "/tmp/$FILENAME"

echo "Installed to $INSTALL_DIR/swag2mcp"
echo "Run 'swag2mcp --version' to verify."
