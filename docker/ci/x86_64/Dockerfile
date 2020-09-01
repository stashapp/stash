
# must be built from /dist directory

FROM ubuntu:20.04 as prep
LABEL MAINTAINER="https://discord.gg/Uz29ny"

RUN apt-get update && \
    apt-get -y install curl xz-utils && \
    apt-get autoclean -y && \
    rm -rf /var/lib/apt/lists/*
WORKDIR /
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

RUN curl --http1.1 -o /ffmpeg.tar.xz https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz && \
    tar xf /ffmpeg.tar.xz && \
    rm ffmpeg.tar.xz && \
    mv /ffmpeg*/ /ffmpeg/

FROM ubuntu:20.04 as app
RUN apt-get update && apt-get -y install ca-certificates
COPY --from=prep /ffmpeg/ffmpeg /ffmpeg/ffprobe /usr/bin/
COPY /stash-linux /usr/bin/stash

EXPOSE 9999
CMD ["stash"]
