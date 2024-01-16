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

# vips version 8.15.0-r0 breaks thumbnail generation on arm32v6
# need to use 8.14.3-r0 from alpine 3.18 instead

RUN apk add --no-cache --virtual .build-deps gcc python3-dev musl-dev \
    && apk add --no-cache ca-certificates python3 py3-requests py3-requests-toolbelt py3-lxml py3-pip ffmpeg ruby tzdata \
    && apk add --no-cache --repository http://dl-cdn.alpinelinux.org/alpine/v3.18/community vips=8.14.3-r0 vips-tools=8.14.3-r0 \
    && pip install --user --break-system-packages mechanicalsoup cloudscraper bencoder.pyx \
    && gem install faraday \
    && apk del .build-deps
ENV STASH_CONFIG_FILE=/root/.stash/config.yml

EXPOSE 9999
CMD ["stash"]
