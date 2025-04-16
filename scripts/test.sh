#!/bin/bash

cd $(dirname "$0")/..

go test -v -count=1 ./...
