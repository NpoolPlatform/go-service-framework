#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

PLATFORM=linux/amd64
OUTPUT=./output

pkg=github.com/NpoolPlatform/go-service-framework/pkg/version

OS="${PLATFORM%/*}"
ARCH=$(basename "$PLATFORM")

if git_status=$(git status --porcelain --untracked=no 2>/dev/null) && [[ -z "${git_status}" ]]; then
    git_tree_state=clean
fi

set +e
version=`git describe --tags --abbrev=0`
if [ ! $? -eq 0 ]; then
    version=latest
fi
set -e

service_name=$1

echo "Release docker image for $PLATFORM -- $version"
kubectl apply -k ./cmd/$service_name/k8s
