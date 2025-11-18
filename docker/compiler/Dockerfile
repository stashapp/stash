FROM golang:1.24.3

LABEL maintainer="https://discord.gg/2TsNFKt"

RUN apt-get update && apt-get install -y apt-transport-https ca-certificates gnupg

RUN mkdir -p /etc/apt/keyrings

ADD https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key nodesource.gpg.key
RUN cat nodesource.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg && rm nodesource.gpg.key
RUN echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_24.x nodistro main" | tee /etc/apt/sources.list.d/nodesource.list

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        git make tar bash nodejs zip \
        clang llvm-dev cmake patch libxml2-dev uuid-dev libssl-dev xz-utils \
        bzip2 gzip sed cpio libbz2-dev zlib1g-dev \
        gcc-mingw-w64 \
        gcc-arm-linux-gnueabi libc-dev-armel-cross linux-libc-dev-armel-cross \
        gcc-aarch64-linux-gnu libc-dev-arm64-cross && \
    rm -rf /var/lib/apt/lists/*;

# pnpm install with npm
RUN npm install -g pnpm

# FreeBSD cross-compilation setup
# https://github.com/smartmontools/docker-build/blob/6b8c92560d17d325310ba02d9f5a4b250cb0764a/Dockerfile#L66
ENV FREEBSD_VERSION       13.4
ENV FREEBSD_DOWNLOAD_URL  http://ftp.plusline.de/FreeBSD/releases/amd64/${FREEBSD_VERSION}-RELEASE/base.txz
ENV FREEBSD_SHA           8e13b0a93daba349b8d28ad246d7beb327659b2ef4fe44d89f447392daec5a7c

RUN cd /tmp && \
    curl -o base.txz $FREEBSD_DOWNLOAD_URL && \
    echo "$FREEBSD_SHA base.txz" | sha256sum -c - && \
    mkdir -p /opt/cross-freebsd && \
    cd /opt/cross-freebsd && \
    tar -xf /tmp/base.txz ./lib/ ./usr/lib/ ./usr/include/ && \
    rm -f /tmp/base.txz && \
    cd /opt/cross-freebsd/usr/lib && \
    find . -xtype l | xargs ls -l | grep ' /lib/' | awk '{print "ln -sf /opt/cross-freebsd"$11 " " $9}' | /bin/sh && \
    ln -s libc++.a libstdc++.a && \
    ln -s libc++.so libstdc++.so

# macOS cross-compilation setup
ENV OSX_SDK_VERSION 	   11.3
ENV OSX_SDK_DOWNLOAD_FILE  MacOSX${OSX_SDK_VERSION}.sdk.tar.xz
ENV OSX_SDK_DOWNLOAD_URL   https://github.com/phracker/MacOSX-SDKs/releases/download/${OSX_SDK_VERSION}/${OSX_SDK_DOWNLOAD_FILE}
ENV OSX_SDK_SHA            cd4f08a75577145b8f05245a2975f7c81401d75e9535dcffbb879ee1deefcbf4
ENV OSXCROSS_REVISION      5e1b71fcceb23952f3229995edca1b6231525b5b
ENV OSXCROSS_DOWNLOAD_URL  https://codeload.github.com/tpoechtrager/osxcross/tar.gz/${OSXCROSS_REVISION}
ENV OSXCROSS_SHA           d3f771bbc20612fea577b18a71be3af2eb5ad2dd44624196cf55de866d008647

RUN cd /tmp && \
    curl -o osxcross.tar.gz $OSXCROSS_DOWNLOAD_URL && \
    echo "$OSXCROSS_SHA osxcross.tar.gz" | sha256sum -c - && \
    mkdir osxcross && \
    tar --strip=1 -C osxcross -xf osxcross.tar.gz && \
    rm -f osxcross.tar.gz && \
    curl -Lo $OSX_SDK_DOWNLOAD_FILE $OSX_SDK_DOWNLOAD_URL && \
    echo "$OSX_SDK_SHA $OSX_SDK_DOWNLOAD_FILE" | sha256sum -c - && \
    mv $OSX_SDK_DOWNLOAD_FILE osxcross/tarballs/ && \
    UNATTENDED=yes SDK_VERSION=$OSX_SDK_VERSION OSX_VERSION_MIN=10.10 osxcross/build.sh && \
    cp osxcross/target/lib/* /usr/lib/ && \
    mv osxcross/target /opt/osx-ndk-x86 && \
    rm -rf /tmp/osxcross

ENV PATH /opt/osx-ndk-x86/bin:$PATH

RUN mkdir -p /root/.ssh && \
    chmod 0700 /root/.ssh && \
    ssh-keyscan github.com > /root/.ssh/known_hosts

# ignore "dubious ownership" errors
RUN git config --global safe.directory '*'

# To test locally:
# make generate
# make ui
# cd docker/compiler
# make build
# docker run --rm -v /PATH_TO_STASH:/stash -w /stash -i -t stashapp/compiler:latest make build-cc-all
# # binaries will show up in /dist
