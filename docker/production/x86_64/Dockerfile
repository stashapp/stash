FROM ubuntu:18.04 as prep
LABEL MAINTAINER="leopere [at] nixc [dot] us"

RUN apt-get update && \
    apt-get -y install curl xz-utils && \
    apt-get autoclean -y && \
    rm -rf /var/lib/apt/lists/*
WORKDIR /
SHELL ["/bin/bash", "-o", "pipefail", "-c"]
RUN curl -L -o /stash $(curl -s https://api.github.com/repos/stashapp/stash/releases/latest | awk '/browser_download_url/ && /stash-linux/' | sed -e 's/.*: "\(.*\)"/\1/') && \
    chmod +x /stash && \
    curl --http1.1 -o /ffmpeg.tar.xz https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz && \
    tar xf /ffmpeg.tar.xz && \
    rm ffmpeg.tar.xz && \
    mv /ffmpeg*/ /ffmpeg/

FROM ubuntu:18.04 as app
RUN apt-get update && apt-get -y install ca-certificates
COPY --from=prep /stash /ffmpeg/ffmpeg /ffmpeg/ffprobe /usr/bin/
EXPOSE 9999
CMD ["stash"]
