#!/bin/bash

set -euo pipefail

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
ROOT_DIR="$SCRIPT_DIR"/..

set -x

GOLANG_CI_LINT=$(go env GOPATH)/bin/golangci-lint

cd "$ROOT_DIR"

if [[ ! -f "$GOLANG_CI_LINT" ]];then
  echo 'golangci-lint not exists, install by:'
  echo 'curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.1.6'
  exit 1
fi

"$GOLANG_CI_LINT" run --config "$SCRIPT_DIR"/.golangci-lint.yml "$@"
