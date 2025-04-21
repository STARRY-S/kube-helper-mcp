#!/bin/bash

set -euo pipefail

cd $(dirname $0)/..

rm -rf bin &> /dev/null
rm -rf dist &> /dev/null
