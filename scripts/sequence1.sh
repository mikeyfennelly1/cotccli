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
source .env.local

function title() {
  echo ""
  echo ""
  echo "########################"
  echo $1
  echo "########################"
  echo ""
}

title "checking health"
cotc health

title "view stream structure"
cotc tree

title "create a new stream"

cotc rmstream --name=$STREAM_NAME

STREAM_NAME="test"
title "starting stream: $STREAM_NAME"
cotc mkstream --name=$STREAM_NAME

title "listing all producers"
cotc lsproducers

title "registering producer"
cotc mkproducer -g test -n clitest

title "starting producer"
cotc produce -n sysinfo-reader--2 -t sysinfo --log-level=debug

