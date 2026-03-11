#!/usr/bin/env bash
set -euo pipefail

BINARY_NAME="cotc"
INSTALL_DIR="/usr/local/bin"

echo "Building $BINARY_NAME..."
go build -o "$BINARY_NAME" .

echo "Installing to $INSTALL_DIR/$BINARY_NAME..."
sudo mv "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"

echo "Installed: $(which $BINARY_NAME)"
