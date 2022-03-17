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
RUN apk add --no-cache ca-certificates python3 py3-requests py3-requests-toolbelt py3-lxml py3-pip ffmpeg vips-tools && pip install --no-cache-dir mechanicalsoup cloudscraper
RUN ln -s /usr/bin/python3 /usr/bin/python
ENV STASH_CONFIG_FILE=/root/.stash/config.yml

EXPOSE 9999
CMD ["stash"]
