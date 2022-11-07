#!/bin/bash

COMPILER_CONTAINER="stashapp/compiler:6"

BUILD_DATE=`go run -mod=vendor scripts/getDate.go`
GITHASH=`git rev-parse --short HEAD`
STASH_VERSION=`git describe --tags --exclude latest_develop`

SETENV="BUILD_DATE=\"$BUILD_DATE\" GITHASH=$GITHASH STASH_VERSION=\"$STASH_VERSION\""
SETUP="export CGO_ENABLED=1;"
WINDOWS="echo '=== Building Windows binary ==='; $SETENV make cross-compile-windows;"
DARWIN="echo '=== Building OSX binary ==='; $SETENV make cross-compile-macos-intel;"
DARWIN_ARM64="echo '=== Building OSX (arm64) binary ==='; $SETENV make cross-compile-macos-applesilicon;"
LINUX_AMD64="echo '=== Building Linux (amd64) binary ==='; $SETENV make cross-compile-linux;"
LINUX_ARM64v8="echo '=== Building Linux (armv8/arm64) binary ==='; $SETENV make cross-compile-linux-arm64v8;"
LINUX_ARM32v7="echo '=== Building Linux (armv7/armhf) binary ==='; $SETENV make cross-compile-linux-arm32v7;"
LINUX_ARM32v6="echo '=== Building Linux (armv6 | Raspberry Pi 1) binary ==='; $SETENV make cross-compile-pi;"
BUILD_COMPLETE="echo '=== Build complete ==='"

BUILD=`echo "$1" | cut -d - -f 1`
if [ "$BUILD" == "windows" ]
then
  echo "Building Windows"
  COMMAND="$SETUP $WINDOWS $BUILD_COMPLETE"
elif [ "$BUILD" == "darwin" ]
then
  echo "Building Darwin(MacOSX)"
  COMMAND="$SETUP $DARWIN $BUILD_COMPLETE"
elif [ "$BUILD" == "amd64" ]
then
  echo "Building Linux AMD64"
  COMMAND="$SETUP $LINUX_AMD64 $BUILD_COMPLETE"
elif [ "$BUILD" == "arm64v8" ]
then
  echo "Building Linux ARM64v8"
  COMMAND="$SETUP $LINUX_ARM64v8 $BUILD_COMPLETE"
elif [ "$BUILD" == "arm32v6" ]
then
  echo "Building Linux ARM32v6"
  COMMAND="$SETUP $LINUX_ARM32v6 $BUILD_COMPLETE"
elif [ "$BUILD" == "arm32v7" ]
then
  echo "Building Linux ARM32v7"
  COMMAND="$SETUP $LINUX_ARM32v7 $BUILD_COMPLETE"
else
  echo "Building All"
  COMMAND="$SETUP $WINDOWS $DARWIN $DARWIN_ARM64 $LINUX_AMD64 $LINUX_ARM64v8 $LINUX_ARM32v7 $LINUX_ARM32v6 $BUILD_COMPLETE"
fi

# Pull Latest Image
docker pull $COMPILER_CONTAINER

# Changed consistency to delegated since this is being used as a build tool. The binded volume shouldn't be changing during its run.
docker run --rm --mount type=bind,source="$(pwd)",target=/stash,consistency=delegated -w /stash $COMPILER_CONTAINER /bin/bash -c "$COMMAND"

