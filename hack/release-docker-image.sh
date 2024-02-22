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
## For testing or production environment, pass the second variable
version=`git describe --exact-match --tags $(git log -n1 --pretty='%h')`
if [ ! $? -eq 0 ]; then
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

registry=uhub.service.ucloud.cn
if [ "x" != $2 ]; then
  registry=$2
fi

service_name=$1

echo "Release docker image for $PLATFORM -- $version"

user=`whoami`
if [ "$user" == "root" ]; then
    docker push $registry/entropypool/$service_name:$version
else
    sudo docker push $registry/entropypool/$service_name:$version
fi
