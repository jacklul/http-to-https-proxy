#!/bin/bash
#shellcheck disable=SC2001,SC2155,SC2206

ARCHITECTURES_TO_BUILD=( "linux/amd64" "linux/386" "linux/arm64" "linux/arm/7" "linux/arm/6" "linux/arm/5" "windows/amd64" "windows/386" "darwin/amd64" "darwin/arm64" )

###################################

[ -z "$GITHUB_WORKSPACE" ] && GITHUB_WORKSPACE="$(readlink -f "$(dirname "$0")")"
[ -n "$1" ] && ARCHITECTURES_TO_BUILD=( $* )

set -e
umask 0022

export CGO_ENABLED=0

for ARCHITECTURE in "${ARCHITECTURES_TO_BUILD[@]}"; do
    export GOOS="$(echo "$ARCHITECTURE" | cut -d/ -f 1)"
    export GOARCH="$(echo "$ARCHITECTURE" | cut -d/ -f 2)"
    export GOARM="$(echo "$ARCHITECTURE" | cut -d/ -f 3)"

    VARIANT="$GOOS-$GOARCH"
    [ -n "$GOARM" ] && VARIANT="$GOOS-$GOARCH-v$GOARM"
    [ "$GOOS" = "windows" ] && VARIANT="win-$GOARCH.exe"

    echo "Building for $ARCHITECTURE..."
    go build -o "$GITHUB_WORKSPACE/http-to-https-proxy-$VARIANT" .
done
