#!/bin/sh

DATE=`go run scripts/getDate.go`
GITHASH=`git rev-parse --short HEAD`
VERSION_FLAGS="-X 'github.com/stashapp/stash/pkg/api.buildstamp=$DATE' -X 'github.com/stashapp/stash/pkg/api.githash=$GITHASH'"

SETUP="export GO111MODULE=on; export CGO_ENABLED=1;"
WINDOWS="GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ packr2 build -o dist/stash-win.exe -ldflags \"-extldflags '-static' $VERSION_FLAGS\" -tags extended -v -mod=vendor;"
DARWIN="GOOS=darwin GOARCH=amd64 CC=o64-clang CXX=o64-clang++ packr2 build -o dist/stash-osx -ldflags \"$VERSION_FLAGS\" -tags extended -v -mod=vendor;"
LINUX="packr2 build -o dist/stash-linux -ldflags \"$VERSION_FLAGS\" -v -mod=vendor;"

COMMAND="$SETUP $WINDOWS $DARWIN $LINUX"

docker run --rm --mount type=bind,source="$(pwd)",target=/stash -w /stash stashappdev/compiler:1 /bin/bash -c "$COMMAND"