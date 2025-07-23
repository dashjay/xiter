#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail
set -o errtrace


if [[ $(git diff --exit-code) ]]; then
  echo "workspace is not clean"
  exit 1
fi
