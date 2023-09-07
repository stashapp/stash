FROM --platform=$BUILDPLATFORM alpine:latest AS binary
ARG TARGETPLATFORM
WORKDIR /
COPY stash-*  /
RUN if [ "$TARGETPLATFORM" = "linux/arm/v6" ];   then BIN=stash-linux-arm32v6; \
    elif [ "$TARGETPLATFORM" = "linux/arm/v7" ]; then BIN=stash-linux-arm32v7; \
    elif [ "$TARGETPLATFORM" = "linux/arm64" ];  then BIN=stash-linux-arm64v8; \
    elif [ "$TARGETPLATFORM" = "linux/amd64" ];  then BIN=stash-linux; \
    fi; \
    mv $BIN /stash

FROM --platform=$TARGETPLATFORM alpine:latest AS app
COPY --from=binary /stash /usr/bin/
RUN apk add --no-cache --virtual .build-deps gcc python3-dev musl-dev \
    && apk add --no-cache ca-certificates python3 py3-requests py3-requests-toolbelt py3-lxml py3-pip ffmpeg vips-tools ruby tzdata \
    && pip install mechanicalsoup cloudscraper bencoder.pyx \
    && gem install faraday \
    && apk del .build-deps
ENV STASH_CONFIG_FILE=/root/.stash/config.yml

EXPOSE 9999
CMD ["stash"]

FROM --platform=$TARGETPLATFORM nvidia/cuda:12.2.0-base-ubuntu22.04 AS cuda_app
COPY --from=binary /stash /usr/bin/
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get install -y \
    gcc \
    python3-dev \
    python3 \
    python3-pip \
    python3-requests \
    python3-requests-toolbelt \
    python3-lxml \
    ffmpeg \
    libvips-tools \
    ruby \
    gem \
    tzdata \
    build-essential \
    wget \
    && pip3 install mechanicalsoup cloudscraper bencoder.pyx \
    && gem install faraday
ENV LANG C.UTF-8
ENV NVIDIA_VISIBLE_DEVICES all
ENV NVIDIA_DRIVER_CAPABILITIES=video,utility
ENV STASH_CONFIG_FILE=/root/.stash/config.yml

# NVENC Patch
RUN mkdir -p /usr/local/bin /patched-lib \
    && wget https://raw.githubusercontent.com/keylase/nvidia-patch/master/patch.sh -O /usr/local/bin/patch.sh \
    && wget https://raw.githubusercontent.com/keylase/nvidia-patch/master/docker-entrypoint.sh -O /usr/local/bin/docker-entrypoint.sh \
    && chmod +x /usr/local/bin/patch.sh /usr/local/bin/docker-entrypoint.sh /usr/bin/stash

# Clean
RUN apt-get remove -y build-essential wget && apt-get autoremove -y && apt-get clean && rm -rf /var/lib/apt/lists/*

EXPOSE 9999
ENTRYPOINT ["docker-entrypoint.sh", "stash"]
