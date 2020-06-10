#!/bin/bash

# This script creates platform specific Dockerfiles from Dockerfile.template.

cd $(dirname $0)

arch_list="arm32v7 arm64v8 amd64"
arm32v7="stash-linux-arm32v7 arm32v7/ubuntu armel"
arm64v8="stash-linux-arm64v8 arm64v8/ubuntu armhf"
amd64="stash-linux ubuntu amd64"

render() {
    local template_file=${1}
    local stash_binary=${2}
    local base_image=${3}
    local ffmpeg_arch=${4}

    sed -E "
        s!%%STASH_BINARY%%!$stash_binary!g;
        s!%%BASE_IMAGE%%!$base_image!g;
        s!%%FFMPEG_ARCH%%!$ffmpeg_arch!g;
        " $template_file
}

for arch in $arch_list; do
    mkdir -p $arch
    render ./Dockerfile.template ${!arch} > "$arch/Dockerfile"
done
