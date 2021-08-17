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
ENV DEBIAN_FRONTEND=noninteractive
RUN apt update && apt install -y python3 python-is-python3 python3-requests python3-requests-toolbelt python3-lxml python3-pip && pip3 install cloudscraper
FROM ubuntu:20.04 as app
RUN apt update && apt install -y python3 python-is-python3 python3-requests python3-requests-toolbelt python3-lxml ffmpeg && rm -rf /var/lib/apt/lists/*
COPY --from=prep /stash /usr/bin/
COPY --from=prep /usr/local/lib/python3.8/dist-packages /usr/local/lib/python3.8/dist-packages

ENV STASH_CONFIG_FILE=/root/.stash/config.yml

EXPOSE 9999
CMD ["stash"]

