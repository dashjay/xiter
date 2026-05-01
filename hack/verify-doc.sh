#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail
set -o errtrace

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)
ROOT_DIR="$SCRIPT_DIR"/..

cd "$ROOT_DIR"

# Generate docs into a temp directory and diff against committed docs
tmpdir=$(mktemp -d)
trap 'rm -rf "$tmpdir"' EXIT

while IFS= read -r -d '' dir; do
    i=$(basename "$dir")
    [[ "$i" == "hack" || "$i" == "doc" || "$i" == ".git" ]] && continue
    mkdir -p "$tmpdir/$i"
    $(go env GOPATH)/bin/gomarkdoc \
        --repository.url=https://github.com/dashjay/xiter \
        --repository.default-branch=main \
        ./$i/... \
        --output \
        "$tmpdir/$i/README.md"
    if [ -f "$tmpdir/$i/README.md" ]; then
        if ! diff -q "./$i/README.md" "$tmpdir/$i/README.md" &>/dev/null; then
            echo "ERROR: $i/README.md is out of date. Run hack/update-doc.sh to regenerate."
            diff "./$i/README.md" "$tmpdir/$i/README.md"
            exit 1
        fi
    fi
done < <(find . -maxdepth 1 -type d -not -name '.' -print0)

echo "All docs are up to date."
