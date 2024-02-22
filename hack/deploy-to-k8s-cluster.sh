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
## Get tag we're on
# version=`git describe --tags --abbrev=0`
version=`git describe --exact-match --tags $(git log -n1 --pretty='%h')`
if [ ! $? -eq 0 ]; then
  ## branch=`git branch --show-current` // Only for git 2.22^
  branch=`git rev-parse --abbrev-ref HEAD | grep -v ^HEAD$ || git rev-parse HEAD`
  if [ "x$branch" == "xmaster" ]; then
    version=latest
  else
    version=`echo $branch | awk -F '/' '{print $2}'`
  fi
  ## Do we need commit ?
  # commit=`git rev-parse HEAD`
  # version=$version-$commit
fi
set -e

service_name=$1

echo "Deploy docker image for $PLATFORM -- $version"
kubectl apply -k ./cmd/$service_name/k8s
