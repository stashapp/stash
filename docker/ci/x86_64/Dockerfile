#FROM --platform=$BUILDPLATFORM ubuntu:20.04  AS prep
FROM --platform=$BUILDPLATFORM alpine  AS prep
ARG TARGETPLATFORM
ARG BUILDPLATFORM
RUN apk add --no-cache curl xz
WORKDIR /
COPY stash-*  /
RUN if [ "$TARGETPLATFORM" = "linux/arm/v6" ]; then FFMPEG_URL=https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-armel-static.tar.xz;BIN=stash-pi; elif [ "$TARGETPLATFORM" = "linux/arm/v7" ]; then FFMPEG_URL=https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-armhf-static.tar.xz;BIN=stash-linux-arm32v7; elif [ "$TARGETPLATFORM" = "linux/arm64" ]; then FFMPEG_URL=https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-arm64-static.tar.xz;BIN=stash-linux-arm64v8;elif [ "$TARGETPLATFORM" = "linux/amd64" ]; then FFMPEG_URL=https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz;BIN=stash-linux;  fi; \
    curl --http1.1 -o /ffmpeg.tar.xz $FFMPEG_URL && \
    tar xf /ffmpeg.tar.xz && \
    rm ffmpeg.tar.xz && \
    mv /ffmpeg*/ /ffmpeg/ && \
    mv $BIN /stash
#FROM ubuntu:20.04 as app
#RUN apt-get update && apt-get -y install ca-certificates
FROM alpine as app
run apk add --no-cache python3 py3-requests sqlite-libs && ln -s /usr/bin/python3 /usr/bin/python
COPY --from=prep /ffmpeg/ffmpeg /ffmpeg/ffprobe /stash /usr/bin/

EXPOSE 9999
CMD ["stash"]

