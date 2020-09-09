#!/bin/sh

BUILD_DATE=`go run -mod=vendor scripts/getDate.go`
GITHASH=`git rev-parse --short HEAD`
STASH_VERSION=`git describe --tags --exclude latest_develop`
SETENV="BUILD_DATE=\"$BUILD_DATE\" GITHASH=$GITHASH STASH_VERSION=\"$STASH_VERSION\""
SETUP="export GO111MODULE=on; export CGO_ENABLED=1; set -e; echo '=== Running packr ==='; make packr;"
WINDOWS="echo '=== Building Windows binary ==='; $SETENV GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ LDFLAGS=\"-extldflags '-static' \" OUTPUT=\"dist/stash-win.exe\" make build-release;"
DARWIN="echo '=== Building OSX binary ==='; $SETENV GOOS=darwin GOARCH=amd64 CC=o64-clang CXX=o64-clang++ OUTPUT=\"dist/stash-osx\" make build-release;"
LINUX_AMD64="echo '=== Building Linux (amd64) binary ==='; $SETENV GOOS=linux GOARCH=amd64 OUTPUT=\"dist/stash-linux\" make build-release;"
LINUX_ARM64v8="echo '=== Building Linux (armv8/arm64) binary ==='; $SETENV GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc OUTPUT=\"dist/stash-linux-arm64v8\" make build-release;"
LINUX_ARM32v7="echo '=== Building Linux (armv7/armhf) binary ==='; $SETENV GOOS=linux GOARCH=arm GOARM=7 CC=arm-linux-gnueabihf-gcc OUTPUT=\"dist/stash-linux-arm32v7\" make build-release;"
LINUX_ARM32v6="echo '=== Building Linux (armv6 | Raspberry Pi 1) binary ==='; $SETENV GOOS=linux GOARCH=arm GOARM=6 CC=arm-linux-gnueabi-gcc OUTPUT=\"dist/stash-pi\" make build-release;"

COMMAND="$SETUP $WINDOWS $DARWIN $LINUX_AMD64 $LINUX_ARM64v8 $LINUX_ARM32v7 $LINUX_ARM32v6 echo '=== Build complete ==='"

docker run --rm --mount type=bind,source="$(pwd)",target=/stash -w /stash stashapp/compiler:3 /bin/bash -c "$COMMAND"
