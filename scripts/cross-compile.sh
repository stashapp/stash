#!/bin/sh

BUILD_DATE=`go run -mod=vendor scripts/getDate.go`
GITHASH=`git rev-parse --short HEAD`
STASH_VERSION=`git describe --tags --exclude latest_develop`
SETENV="BUILD_DATE=\"$BUILD_DATE\" GITHASH=$GITHASH STASH_VERSION=\"$STASH_VERSION\""
SETUP="export GO111MODULE=on; export CGO_ENABLED=1; make packr;"
WINDOWS="echo '=== Building Windows binary ==='; $SETENV GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ LDFLAGS=\"-extldflags '-static' \" OUTPUT=\"dist/stash-win.exe\" make build-release;"
DARWIN="echo '=== Building OSX binary ==='; $SETENV GOOS=darwin GOARCH=amd64 CC=o64-clang CXX=o64-clang++ OUTPUT=\"dist/stash-osx\" make build-release;"
LINUX="echo '=== Building Linux binary ==='; $SETENV OUTPUT=\"dist/stash-linux\" make build-release;"
RASPPI="echo '=== Building Raspberry Pi binary ==='; $SETENV GOOS=linux GOARCH=arm GOARM=5 CC=arm-linux-gnueabi-gcc OUTPUT=\"dist/stash-pi\" make build-release;"

COMMAND="$SETUP $WINDOWS $DARWIN $LINUX $RASPPI"

docker run --rm --mount type=bind,source="$(pwd)",target=/stash -w /stash stashapp/compiler:develop /bin/bash -c "$COMMAND"
