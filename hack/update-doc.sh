#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail
set -o errtrace

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)
ROOT_DIR="$SCRIPT_DIR"/..

cd "$ROOT_DIR"

for i in $(ls pkg/)
do $(go env GOPATH)/bin/gomarkdoc \
    --repository.url=https://github.com/dashjay/xiter \
    --repository.default-branch=main \
    ./pkg/$i/... \
    --output \
    ./pkg/$i/README.md
done
