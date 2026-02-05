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

RUN apk add --no-cache ca-certificates python3 py3-requests py3-requests-toolbelt py3-lxml py3-pip ffmpeg tzdata vips vips-tools \
    && pip install --break-system-packages mechanicalsoup cloudscraper stashapp-tools
ENV STASH_CONFIG_FILE=/root/.stash/config.yml

# Basic build-time metadata as defined at https://github.com/opencontainers/image-spec/blob/main/annotations.md#pre-defined-annotation-keys
LABEL org.opencontainers.image.title="Stash" \
    org.opencontainers.image.description="An organizer for your porn, written in Go." \
    org.opencontainers.image.url="https://stashapp.cc" \
    org.opencontainers.image.documentation="https://docs.stashapp.cc" \
    org.opencontainers.image.source="https://github.com/stashapp/stash" \
    org.opencontainers.image.licenses="AGPL-3.0"

EXPOSE 9999
CMD ["stash"]
