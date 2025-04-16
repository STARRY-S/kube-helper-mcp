#!/bin/bash

# Builds the Go programs and places the binaries into the `bin` dir.
cd $(dirname $0)

set -euo pipefail

cmds=(
    "calculator"
    "ipaddr"
    "kubecheck"
    "weather"
)

mkdir -p bin && cd bin

for c in ${cmds[@]}; do
    echo "Building $c..."
    go build -o ../cmd/$c/main.go -o ./$c
done
