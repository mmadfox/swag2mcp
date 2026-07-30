#!/usr/bin/env bash
set -euo pipefail

REPO="mmadfox/swag2mcp"
VERSION=""
MODE="sudo"  # sudo or local

while [ $# -gt 0 ]; do
  case "$1" in
    --sudo) MODE="sudo"; shift ;;
    --local) MODE="local"; shift ;;
    -v) VERSION="$2"; shift 2 ;;
    -v*) VERSION="${1#-v}"; shift ;;
    *) echo "Usage: $0 [--sudo|--local] [-v version]" >&2; exit 1 ;;
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

FILENAME="swag2mcp-mock_${VERSION#v}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/download/$VERSION/$FILENAME"

if [ "$MODE" = "local" ]; then
  echo "Installing to ~/.local/bin (no sudo)..."
  echo "After install, add to your shell config (~/.zshrc, ~/.bashrc):"
  echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
  echo ""
fi

echo "Downloading swag2mcp-mock $VERSION ($OS/$ARCH)..."
curl -fsSL "$URL" -o "/tmp/$FILENAME"

echo "Extracting..."
tar xzf "/tmp/$FILENAME" -C /tmp

if [ "$MODE" = "local" ]; then
  INSTALL_DIR="$HOME/.local/bin"
  mkdir -p "$INSTALL_DIR"
  mv /tmp/swag2mcp-mock "$INSTALL_DIR/swag2mcp-mock"
elif [ -w "/usr/local/bin" ]; then
  INSTALL_DIR="/usr/local/bin"
  mv /tmp/swag2mcp-mock "$INSTALL_DIR/swag2mcp-mock"
elif command -v sudo &>/dev/null; then
  INSTALL_DIR="/usr/local/bin"
  if ! sudo mv /tmp/swag2mcp-mock "$INSTALL_DIR/swag2mcp-mock" 2>/dev/null; then
    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p "$INSTALL_DIR"
    mv /tmp/swag2mcp-mock "$INSTALL_DIR/swag2mcp-mock"
  fi
else
  INSTALL_DIR="$HOME/.local/bin"
  mkdir -p "$INSTALL_DIR"
  mv /tmp/swag2mcp-mock "$INSTALL_DIR/swag2mcp-mock"
fi

chmod +x "$INSTALL_DIR/swag2mcp-mock"
rm -f "/tmp/$FILENAME"

echo "Installed to $INSTALL_DIR/swag2mcp-mock"

if [ "$INSTALL_DIR" = "$HOME/.local/bin" ]; then
  echo ""
  echo "  Add to your shell config (~/.zshrc, ~/.bashrc):"
  echo "    export PATH=\"\$HOME/.local/bin:\$PATH\""
  echo ""
  echo "  Or move to /usr/local/bin:"
  echo "    sudo mv $INSTALL_DIR/swag2mcp-mock /usr/local/bin/"
fi

echo "Run 'hash -r' or open a new terminal, then 'swag2mcp-mock --version' to verify."
