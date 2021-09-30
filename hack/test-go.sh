#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)

# Default timeout is 1800s
TEST_TIMEOUT=${TIMEOUT:-1800}

# Write go test artifacts here
ARTIFACTS=${ARTIFACTS:-"${REPO_ROOT}/tmp"}

for arg in "$@"
do
    case $arg in
        -t=*|--timeout=*)
        TEST_TIMEOUT="${arg#*=}"
        shift
        ;;
        -t|--timeout)
        TEST_TIMEOUT="$2"
        shift
        shift
    esac
done

cd "${REPO_ROOT}"

mkdir -p "${ARTIFACTS}"

go_test_flags=(
    -v
    -count=1
    -timeout="${TEST_TIMEOUT}s"
    -cover -coverprofile "${ARTIFACTS}/coverage.out"
)

go test ./... -coverprofile ${ARTIFACTS}/coverage.out
go tool cover -html "${ARTIFACTS}/coverage.out" -o "${ARTIFACTS}/coverage.html"
