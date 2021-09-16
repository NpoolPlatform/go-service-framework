#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

PLATFORMS=(
    linux/amd64
    windows/amd64
    darwin/amd64
)
OUTPUT=./output

pkg=github.com/NpoolPlatform/version

for PLATFORM in "${PLATFORMS[@]}"; do
    OS="${PLATFORM%/*}"
    ARCH=$(basename "$PLATFORM")

    if git_status=$(git status --porcelain --untracked=no 2>/dev/null) && [[ -z "${git_status}" ]]; then
        git_tree_state=clean
    fi

    go build -v -ldflags "-s -w \
        -X $pkg.buildDate=$(date -u +'%Y-%m-%dT%H:%M:%SZ') \
        -X $pkg.gitCommit=$(git rev-parse HEAD 2>/dev/null || echo unknown) \
        -X $pkg.gitVersion=$(git describe --tags --abbrev=0 || echo unknown)" \
        -o "${OUTPUT}/${OS}/${ARCH}/" "$(pwd)/cmd/..." \
        || return 1

    echo "Building project for $PLATFORM"
    GOARCH="$ARCH" GOOS="$OS" go build -o output/ ./...
done
