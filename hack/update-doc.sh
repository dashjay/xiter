#!/usr/bin/env bash

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)
ROOT_DIR="$SCRIPT_DIR"/..

cd "$ROOT_DIR"


for i in $(ls pkg/)
do gomarkdoc ./pkg/$i/... --output ./pkg/$i/README.md
done
