#!/bin/sh
set -e

REPO="tomiok/rend"
BINARY="rend"
INSTALL_DIR="/usr/local/bin"

#  OS & arch
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
  x86_64)  ARCH="amd64" ;;
  arm64)   ARCH="arm64" ;;
  aarch64) ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

case $OS in
  linux|darwin) ;;
  *)
    echo "Unsupported OS: $OS"
    exit 1
    ;;
esac

# obtener la ultima version
VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" \
  | grep '"tag_name"' \
  | sed 's/.*"tag_name": *"\([^"]*\)".*/\1/')

if [ -z "$VERSION" ]; then
  echo "Could not determine latest version"
  exit 1
fi

FILENAME="rend-${OS}-${ARCH}"
URL="https://github.com/$REPO/releases/download/$VERSION/$FILENAME"

echo "Installing rend $VERSION for $OS/$ARCH..."

# download
TMP=$(mktemp)
curl -sL "$URL" -o "$TMP"
chmod +x "$TMP"

# install
if [ -w "$INSTALL_DIR" ]; then
  mv "$TMP" "$INSTALL_DIR/$BINARY"
else
  sudo mv "$TMP" "$INSTALL_DIR/$BINARY"
fi

echo ""
echo "rend installed to $INSTALL_DIR/$BINARY"
echo "Run: rend add clock"