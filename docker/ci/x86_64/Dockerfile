FROM --platform=$BUILDPLATFORM ubuntu:20.04  AS prep
ARG TARGETPLATFORM
WORKDIR /
COPY stash-*  /
RUN if [ "$TARGETPLATFORM" = "linux/arm/v6" ];   then BIN=stash-pi; \
    elif [ "$TARGETPLATFORM" = "linux/arm/v7" ]; then BIN=stash-linux-arm32v7; \
    elif [ "$TARGETPLATFORM" = "linux/arm64" ];  then BIN=stash-linux-arm64v8; \
    elif [ "$TARGETPLATFORM" = "linux/amd64" ];  then BIN=stash-linux; \
    fi; \
    mv $BIN /stash
FROM ubuntu:20.04 as app
run apt update && apt install -y python3 python3 python-is-python3  python3-requests && rm -rf /var/lib/apt/lists/*
COPY --from=prep /stash /usr/bin/

EXPOSE 9999
CMD ["stash"]

