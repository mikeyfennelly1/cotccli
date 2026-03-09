#!/bin/bash

################################################
# SETUP
################################################
OS=$(uname)
if [[ "$OS" == "Darwin" ]]; then
        # OSX uses BSD readlink
        BASEDIR="$(dirname "$0")"
else
        BASEDIR=$(readlink -e "$(dirname "$0")/")
fi
pushd "${BASEDIR}/.."

function title() {
  echo ""
  echo ""
  echo "########################"
  echo $1
  echo "########################"
  echo ""
}

title "checking health"
go run main.go health

title "view stream structure"
go run main.go tree

title "create a new stream"
STREAM_NAME="test.sequence1child"
echo "starting stream: $STREAM_NAME"
go run main.go mkstream --name=$STREAM_NAME

echo "listing all producers"
go run main.go lsproducers


