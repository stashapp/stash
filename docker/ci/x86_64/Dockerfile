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
RUN apk add --no-cache shadow && \
    groupadd -g 1000 stash && \
    useradd -u 1000 -g stash -m -s /bin/sh stash && \
    apk del shadow
COPY --from=binary /stash /usr/bin/
RUN mkdir -p /home/stash/.stash && chown -R stash:stash /home/stash/.stash

RUN apk add --no-cache ca-certificates python3 py3-requests py3-requests-toolbelt py3-lxml py3-pip ffmpeg tzdata vips vips-tools vips-heif \
    && pip install --break-system-packages mechanicalsoup cloudscraper stashapp-tools
ENV STASH_CONFIG_FILE=/home/stash/.stash/config.yml
USER stash

# Basic build-time metadata as defined at https://github.com/opencontainers/image-spec/blob/main/annotations.md#pre-defined-annotation-keys
LABEL org.opencontainers.image.title="Stash" \
    org.opencontainers.image.description="An organizer for your porn, written in Go." \
    org.opencontainers.image.url="https://stashapp.cc" \
    org.opencontainers.image.documentation="https://docs.stashapp.cc" \
    org.opencontainers.image.source="https://github.com/stashapp/stash" \
    org.opencontainers.image.licenses="AGPL-3.0"

EXPOSE 9999
CMD ["stash"]
