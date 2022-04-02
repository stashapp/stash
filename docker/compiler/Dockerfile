FROM golang:1.17

LABEL maintainer="https://discord.gg/2TsNFKt"

# Install tools
RUN apt-get update && apt-get install -y apt-transport-https
RUN curl -sL https://deb.nodesource.com/setup_lts.x | bash -

# prevent caching of the key
ADD https://dl.yarnpkg.com/debian/pubkey.gpg yarn.gpg
RUN cat yarn.gpg | apt-key add - && \
    echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list && \
    rm yarn.gpg

RUN apt-get update && \
    apt-get install -y automake autogen cmake \
    libtool libxml2-dev uuid-dev libssl-dev bash \
    patch make tar xz-utils bzip2 gzip zlib1g-dev sed cpio \
	gcc-10-multilib gcc-mingw-w64 g++-mingw-w64 clang llvm-dev \
	gcc-arm-linux-gnueabi libc-dev-armel-cross linux-libc-dev-armel-cross \
    gcc-arm-linux-gnueabihf libc-dev-armhf-cross \
    gcc-aarch64-linux-gnu libc-dev-arm64-cross \
	nodejs yarn zip --no-install-recommends || exit 1; \
	rm -rf /var/lib/apt/lists/*;

# Cross compile setup
ENV OSX_SDK_VERSION 	11.3
ENV OSX_SDK_DOWNLOAD_FILE=MacOSX${OSX_SDK_VERSION}.sdk.tar.xz
ENV OSX_SDK_DOWNLOAD_URL=https://github.com/phracker/MacOSX-SDKs/releases/download/${OSX_SDK_VERSION}/${OSX_SDK_DOWNLOAD_FILE}
ENV OSX_SDK_SHA=cd4f08a75577145b8f05245a2975f7c81401d75e9535dcffbb879ee1deefcbf4
ENV OSX_SDK     		MacOSX$OSX_SDK_VERSION.sdk
ENV OSX_NDK_X86 		/usr/local/osx-ndk-x86

RUN  wget ${OSX_SDK_DOWNLOAD_URL}
RUN  echo "$OSX_SDK_SHA $OSX_SDK_DOWNLOAD_FILE" | sha256sum -c - || exit 1; \
     git clone https://github.com/tpoechtrager/osxcross.git; \
     mv $OSX_SDK_DOWNLOAD_FILE osxcross/tarballs/

RUN     UNATTENDED=yes SDK_VERSION=${OSX_SDK_VERSION} OSX_VERSION_MIN=10.10 osxcross/build.sh || exit 1;
RUN     cp osxcross/target/lib/* /usr/lib/ ; \
        mv osxcross/target $OSX_NDK_X86; \
        rm -rf osxcross;

ENV PATH $OSX_NDK_X86/bin:$PATH

RUN mkdir -p /root/.ssh; \
    chmod 0700 /root/.ssh; \
    ssh-keyscan github.com > /root/.ssh/known_hosts;

# Notes for self:

# To test locally:
# make generate
# make ui
# cd docker/compiler
# make build
# docker run -it -v /PATH_TO_STASH:/go/stash stashapp/compiler:latest /bin/bash
# cd stash
# make cross-compile-all
# # binaries will show up in /dist

# Windows:
# GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++  go build -ldflags "-extldflags '-static'" -tags extended

# Darwin
# CC=o64-clang CXX=o64-clang++ GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -tags extended
# env goreleaser --config=goreleaser-extended.yml --skip-publish --skip-validate --rm-dist --release-notes=temp/0.48-relnotes-ready.md
