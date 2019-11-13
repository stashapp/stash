FROM golang:1.11.5

LABEL maintainer="stashappdev@gmail.com"

ENV PACKR2_VERSION=2.0.2
ENV PACKR2_SHA=f95ff4c96d7a28813220df030ad91700b8464fe292ab3e1dc9582305c2a338d2
ENV PACKR2_DOWNLOAD_FILE=packr_${PACKR2_VERSION}_linux_amd64.tar.gz
ENV PACKR2_DOWNLOAD_URL=https://github.com/gobuffalo/packr/releases/download/v${PACKR2_VERSION}/${PACKR2_DOWNLOAD_FILE}

# Install tools
RUN apt-get update && apt-get install -y apt-transport-https
RUN curl -sL https://deb.nodesource.com/setup_10.x | bash -
RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add - && \
    echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list


RUN apt-get update && \
    apt-get install -y automake autogen \
    libtool libxml2-dev uuid-dev libssl-dev bash \
    patch make tar xz-utils bzip2 gzip sed cpio \
	gcc-6-multilib g++-6-multilib gcc-mingw-w64 g++-mingw-w64 clang llvm-dev \
	gcc-arm-linux-gnueabi libc-dev-armel-cross linux-libc-dev-armel-cross \
	nodejs yarn --no-install-recommends || exit 1; \
	rm -rf /var/lib/apt/lists/*;

# Cross compile setup
ENV OSX_SDK_VERSION 	10.11
ENV OSX_SDK_DOWNLOAD_FILE=MacOSX${OSX_SDK_VERSION}.sdk.tar.xz
ENV OSX_SDK_DOWNLOAD_URL=https://github.com/ndeloof/golang-cross/raw/113fix/${OSX_SDK_DOWNLOAD_FILE}
ENV OSX_SDK_SHA=98cdd56e0f6c1f9e1af25e11dd93d2e7d306a4aa50430a2bc6bc083ac67efbb8
ENV OSX_SDK     		MacOSX$OSX_SDK_VERSION.sdk
ENV OSX_NDK_X86 		/usr/local/osx-ndk-x86

RUN  wget ${OSX_SDK_DOWNLOAD_URL}
RUN  echo "$OSX_SDK_SHA $OSX_SDK_DOWNLOAD_FILE" | sha256sum -c - || exit 1; \
     git clone https://github.com/tpoechtrager/osxcross.git && \
     git -C osxcross checkout a9317c18a3a457ca0a657f08cc4d0d43c6cf8953 || exit 1; \
     mv $OSX_SDK_DOWNLOAD_FILE osxcross/tarballs/ && \
     UNATTENDED=yes SDK_VERSION=${OSX_SDK_VERSION} OSX_VERSION_MIN=10.9 osxcross/build.sh || exit 1; \
     mv osxcross/target $OSX_NDK_X86; \
     rm -rf osxcross;

ENV PATH $OSX_NDK_X86/bin:$PATH

RUN mkdir -p /root/.ssh; \
    chmod 0700 /root/.ssh; \
    ssh-keyscan github.com > /root/.ssh/known_hosts;

RUN  wget ${PACKR2_DOWNLOAD_URL}; \
			echo "$PACKR2_SHA $PACKR2_DOWNLOAD_FILE" | sha256sum -c - || exit 1; \
			tar -xzf $PACKR2_DOWNLOAD_FILE -C /usr/bin/ packr2; \
			rm $PACKR2_DOWNLOAD_FILE;

CMD ["packr2", "version"]


# Notes for self:
# Windows:
# GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++  go build -ldflags "-extldflags '-static'" -tags extended


# Darwin
# CC=o64-clang CXX=o64-clang++ GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -tags extended
# env GO111MODULE=on goreleaser --config=goreleaser-extended.yml --skip-publish --skip-validate --rm-dist --release-notes=temp/0.48-relnotes-ready.md
