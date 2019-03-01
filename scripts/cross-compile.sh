#!/bin/sh

SETUP="go install github.com/gobuffalo/packr/v2/packr2; export GO111MODULE=on; export CGO_ENABLED=1;"
WINDOWS="GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ packr2 build -o dist/stash-win.exe -ldflags \"-extldflags '-static'\" -tags extended -v -mod=vendor;"
DARWIN="GOOS=darwin GOARCH=amd64 CC=o64-clang CXX=o64-clang++ packr2 build -o dist/stash-osx -tags extended -v -mod=vendor;"
LINUX="packr2 build -o dist/stash-linux -v -mod=vendor;"

COMMAND="$SETUP $WINDOWS $DARWIN $LINUX"

docker run --rm --mount type=bind,source="$(pwd)",target=/stash -w /stash bepsays/ci-goreleaser:1.11-2 /bin/bash -c "$COMMAND"