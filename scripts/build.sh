#!/bin/bash

# Builds the Go programs and places the binaries into the `bin` dir.
cd $(dirname $0)/..

set -euo pipefail

cmds=(
    "calculator"
    "ipaddr"
    "kubecheck"
    "weather"
)

mkdir -p bin && cd bin

COMMIT=${COMMIT:-"UNKNOW"} TAG=${TAG:-"HEAD"}

for c in ${cmds[@]}; do
    echo "Building $c..."
    go build \
        -buildmode=pie \
        -ldflags="-extldflags='-static' -s -w -X github.com/STARRY-S/learn-mcp/pkg/utils.Version=${TAG} -X github.com/STARRY-S/learn-mcp/pkg/utils.Commit=${COMMIT}" \
        -o ./$c \
        ../cmd/$c/main.go \

done
