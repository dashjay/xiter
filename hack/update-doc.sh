#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail
set -o errtrace

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)
ROOT_DIR="$SCRIPT_DIR"/..

cd "$ROOT_DIR"

while IFS= read -r -d '' dir; do
    i=$(basename "$dir")
    [[ "$i" == "hack" || "$i" == "doc" || "$i" == ".git" ]] && continue
    $(go env GOPATH)/bin/gomarkdoc \
        --repository.url=https://github.com/dashjay/xiter \
        --repository.default-branch=main \
        ./$i/... \
        --output \
        ./$i/README.md
done < <(find . -maxdepth 1 -type d -not -name '.' -print0)
