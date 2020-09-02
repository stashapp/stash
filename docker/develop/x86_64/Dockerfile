FROM ubuntu:18.04 as prep
LABEL MAINTAINER="https://discord.gg/Uz29ny"

RUN apt-get update && \
    apt-get -y install curl xz-utils && \
    apt-get autoclean -y && \
    rm -rf /var/lib/apt/lists/*
WORKDIR /
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# added " to end of stash-linux clause so that it doesn't pick up the arm builds
RUN curl -L -o /stash $(curl -s https://api.github.com/repos/stashapp/stash/releases/tags/latest_develop | awk '/browser_download_url/ && /stash-linux"/' | sed -e 's/.*: "\(.*\)"/\1/') && \
    chmod +x /stash 

RUN curl --http1.1 -o /ffmpeg.tar.xz https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz && \
    tar xf /ffmpeg.tar.xz && \
    rm ffmpeg.tar.xz && \
    mv /ffmpeg*/ /ffmpeg/

FROM ubuntu:18.04 as app
RUN apt-get update && apt-get -y install ca-certificates
COPY --from=prep /stash /ffmpeg/ffmpeg /ffmpeg/ffprobe /usr/bin/
EXPOSE 9999
CMD ["stash"]
