#!/bin/sh
#
set -e

tmpFile=$(mktemp)

( cd $(dirname "$0") &&
	go build -buildvcs="false" -o "$tmpFile" ./cmd/leGit )

exec "$tmpFile" "$@"
