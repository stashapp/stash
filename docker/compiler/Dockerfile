FROM golang:1.11.5

LABEL maintainer="stashappdev@gmail.com"

ENV GORELEASER_VERSION=0.95.0
ENV GORELEASER_SHA=4f3b9fc978a3677806ebd959096a1f976a7c7bb5fbdf7a9a1d01554c8c5c31c5

ENV GORELEASER_DOWNLOAD_FILE=goreleaser_Linux_x86_64.tar.gz
ENV GORELEASER_DOWNLOAD_URL=https://github.com/goreleaser/goreleaser/releases/download/v${GORELEASER_VERSION}/${GORELEASER_DOWNLOAD_FILE}

ENV PACKR2_VERSION=2.0.2
ENV PACKR2_SHA=f95ff4c96d7a28813220df030ad91700b8464fe292ab3e1dc9582305c2a338d2
ENV PACKR2_DOWNLOAD_FILE=packr_${PACKR2_VERSION}_linux_amd64.tar.gz
ENV PACKR2_DOWNLOAD_URL=https://github.com/gobuffalo/packr/releases/download/v${PACKR2_VERSION}/${PACKR2_DOWNLOAD_FILE}

# Install tools
RUN apt-get update && \
    apt-get install -y automake autogen \
    libtool libxml2-dev uuid-dev libssl-dev bash \
    patch make tar xz-utils bzip2 gzip sed cpio \
	gcc-multilib g++-multilib gcc-mingw-w64 g++-mingw-w64 clang llvm-dev --no-install-recommends || exit 1; \
	rm -rf /var/lib/apt/lists/*;

# Cross compile setup
ENV OSX_SDK_VERSION 	10.11
ENV OSX_SDK     		MacOSX$OSX_SDK_VERSION.sdk
ENV OSX_NDK_X86 		/usr/local/osx-ndk-x86
ENV OSX_SDK_PATH 		/$OSX_SDK.tar.gz

COPY $OSX_SDK.tar.gz /go

RUN git clone https://github.com/tpoechtrager/osxcross.git && \
    git -C osxcross checkout c47ff0aeed1a7d0e1f884812fc170e415f05be5a || exit 1; \
    mv $OSX_SDK.tar.gz osxcross/tarballs/ && \
    UNATTENDED=yes SDK_VERSION=${OSX_SDK_VERSION} OSX_VERSION_MIN=10.9 osxcross/build.sh || exit 1; \
    mv osxcross/target $OSX_NDK_X86; \
    rm -rf osxcross;

ENV PATH $OSX_NDK_X86/bin:$PATH

RUN mkdir -p /root/.ssh; \
    chmod 0700 /root/.ssh; \
    ssh-keyscan github.com > /root/.ssh/known_hosts;

RUN  wget ${GORELEASER_DOWNLOAD_URL}; \
			echo "$GORELEASER_SHA $GORELEASER_DOWNLOAD_FILE" | sha256sum -c - || exit 1; \
			tar -xzf $GORELEASER_DOWNLOAD_FILE -C /usr/bin/ goreleaser; \
			rm $GORELEASER_DOWNLOAD_FILE;

RUN  wget ${PACKR2_DOWNLOAD_URL}; \
			echo "$PACKR2_SHA $PACKR2_DOWNLOAD_FILE" | sha256sum -c - || exit 1; \
			tar -xzf $PACKR2_DOWNLOAD_FILE -C /usr/bin/ packr2; \
			rm $PACKR2_DOWNLOAD_FILE;

CMD ["goreleaser", "-v"]
CMD ["packr2", "version"]


# Notes for self:
# Windows:
# GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++  go build -ldflags "-extldflags '-static'" -tags extended


# Darwin
# CC=o64-clang CXX=o64-clang++ GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -tags extended
# env GO111MODULE=on goreleaser --config=goreleaser-extended.yml --skip-publish --skip-validate --rm-dist --release-notes=temp/0.48-relnotes-ready.md