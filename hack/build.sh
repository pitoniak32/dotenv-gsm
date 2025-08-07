#!/bin/bash -ue

if [ "$#" -ne 2 ]; then
  echo "Illegal number of parameters!"
  echo "USAGE: ./build.sh OS ARCH"
fi

GIT_COMMIT=$(git rev-parse HEAD)
GIT_COMMIT_SHORT=$(git rev-parse --short HEAD)
GIT_COMMIT_SHORT=${GIT_COMMIT_SHORT:-"dirty"}
VERSION_RELEASE_TAG=${VERSION_RELEASE_TAG:-"v0.0.0"}

RELEASE_TAG="$VERSION_RELEASE_TAG-$GIT_COMMIT_SHORT"

BUILD_DATE=$(date ${SOURCE_DATE_EPOCH:+"--date=@${SOURCE_DATE_EPOCH}"} -u +'%Y-%m-%dT%H:%M:%SZ')
BUILD_DATE=${BUILD_DATE:-""}

export GOOS=$1
export GOARCH=$2
export CGO_ENABLED=0

echo "building os=$GOOS arch=$GOARCH"

go build \
  -ldflags="-X github.com/pitoniak32/dotenv_gsm/internal/version.commit=$GIT_COMMIT -X github.com/pitoniak32/dotenv_gsm/internal/version.version=$RELEASE_TAG -X github.com/pitoniak32/dotenv_gsm/internal/version.buildDate=$BUILD_DATE" \
  -o ./bin/dotenv_gsm-$GOOS-$GOARCH .