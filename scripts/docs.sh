#!/usr/bin/env bash
set -euo pipefail

PORT=6060
DOC_URL="http://localhost:$PORT/pkg/github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/"

# Install godoc if not present
if ! command -v godoc &>/dev/null; then
  echo "godoc not found, installing..."
  go install golang.org/x/tools/cmd/godoc@latest
fi

# Open browser after a short delay to let the server start
(sleep 1 && xdg-open "$DOC_URL") &

echo "Starting godoc server at $DOC_URL"
godoc -http=":$PORT" &
