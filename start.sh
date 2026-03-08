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
cd "${BASEDIR}"

source "${BASEDIR}"/../.env
source "${BASEDIR}"/scripts/helpers.sh

var_must_exist DESKTOP_SYSINFO_LISTEN_PORT DESKTOP_SYSINFO_LOG_LEVEL

go run "${BASEDIR}"/main.go producer -n 1 -s test -t sysinfo --log-level=debug -p "${COLLECTOR_LISTEN_PORT}"