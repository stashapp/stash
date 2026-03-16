### OSXCROSS
FROM debian:bookworm AS osxcross
# add osxcross
WORKDIR /tmp/osxcross
ARG OSXCROSS_REVISION=5e1b71fcceb23952f3229995edca1b6231525b5b
ADD --checksum=sha256:d3f771bbc20612fea577b18a71be3af2eb5ad2dd44624196cf55de866d008647 https://codeload.github.com/tpoechtrager/osxcross/tar.gz/${OSXCROSS_REVISION} /tmp/osxcross.tar.gz

ARG OSX_SDK_VERSION=11.3
ARG OSX_SDK_DOWNLOAD_FILE=MacOSX${OSX_SDK_VERSION}.sdk.tar.xz
ARG OSX_SDK_DOWNLOAD_URL=https://github.com/phracker/MacOSX-SDKs/releases/download/${OSX_SDK_VERSION}/${OSX_SDK_DOWNLOAD_FILE}
ADD --checksum=sha256:cd4f08a75577145b8f05245a2975f7c81401d75e9535dcffbb879ee1deefcbf4 ${OSX_SDK_DOWNLOAD_URL} /tmp/osxcross/tarballs/${OSX_SDK_DOWNLOAD_FILE}

ENV UNATTENDED=yes \
    SDK_VERSION=${OSX_SDK_VERSION} \
    OSX_VERSION_MIN=10.10
RUN apt update && \
  apt install -y --no-install-recommends \
  bash ca-certificates clang cmake git patch libssl-dev bzip2 cpio libbz2-dev libxml2-dev make python3 xz-utils zlib1g-dev
# lzma-dev libxml2-dev xz
RUN tar --strip=1 -C /tmp/osxcross -xf /tmp/osxcross.tar.gz
RUN ./build.sh

### FREEBSD cross-compilation stage
# use alpine for cacheable image since apt is notorous for not caching
FROM alpine:3 AS freebsd
# match golang latest
# https://go.dev/wiki/FreeBSD
ARG FREEBSD_VERSION=12.4
ADD --checksum=sha256:581c7edacfd2fca2bdf5791f667402d22fccd8a5e184635e0cac075564d57aa8 \
  http://ftp-archive.freebsd.org/mirror/FreeBSD-Archive/old-releases/amd64/${FREEBSD_VERSION}-RELEASE/base.txz \
  /tmp/base.txz

WORKDIR /opt/cross-freebsd
RUN apk add --no-cache tar xz
RUN tar -xf /tmp/base.txz --strip-components=1 ./usr/lib ./usr/include ./lib
RUN cd /opt/cross-freebsd/usr/lib && \
    find . -type l -exec sh -c ' \
      for link; do \
        target=$(readlink "$link"); \
        case "$target" in \
          /lib/*) ln -sf "/opt/cross-freebsd$target" "$link";; \
        esac; \
      done \
    ' sh {} + && \
    ln -s libc++.a libstdc++.a && \
    ln -s libc++.so libstdc++.so

### BUILDER
FROM golang:1.24.3 AS builder
ENV PATH=/opt/osx-ndk-x86/bin:$PATH

# copy in nodejs instead of using nodesource :thumbsup:
COPY --from=docker.io/library/node:24-bookworm /usr/local /usr/local
# copy in osxcross
COPY --from=osxcross /tmp/osxcross/target/lib /usr/lib
COPY --from=osxcross /tmp/osxcross/target /opt/osx-ndk-x86
# copy in cross-freebsd 
COPY --from=freebsd /opt/cross-freebsd /opt/cross-freebsd

# pnpm install with npm
RUN npm install -g pnpm

# git for getting hash
# make and bash for building

# clang for macos
# zip for stashapp.zip
# gcc-extensions for cross-arch build
# we still target arm soft float?
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        git make bash \
        clang zip \
        gcc-mingw-w64 \
        gcc-arm-linux-gnueabi \
        libc-dev-armel-cross linux-libc-dev-armel-cross \
        gcc-aarch64-linux-gnu libc-dev-arm64-cross && \
    rm -rf /var/lib/apt/lists/*;
RUN git config --global safe.directory '*'
# To test locally:
# make generate
# make ui
# cd docker/compiler
# docker build . -t ghcr.io/stashapp/compiler:latest
# docker run --rm -v /PATH_TO_STASH:/stash -w /stash -i -t ghcr.io/stashapp/compiler:latest make build-cc-all
# # binaries will show up in /dist