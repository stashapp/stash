FROM ubuntu:18.04 as prep
LABEL MAINTAINER="leopere [at] nixc [dot] us"

RUN apt-get update && \
    apt-get -y install curl xz-utils ca-certificates -y && \
    update-ca-certificates && \
    apt-get autoclean -y && \
    rm -rf /var/lib/apt/lists/*
WORKDIR /
SHELL ["/bin/bash", "-o", "pipefail", "-c"]
RUN curl -L -o /stash $(curl -s https://api.github.com/repos/stashapp/stash/releases | grep -F 'stash-linux' | grep download | head -n 1 | cut -d'"' -f4) && \
    chmod +x /stash && \
    curl -o /ffmpeg.tar.xz https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz && \
    tar xf /ffmpeg.tar.xz && \
    rm ffmpeg.tar.xz && \
    mv /ffmpeg*/ /ffmpeg/

FROM ubuntu:18.04 as app
RUN adduser stash --gecos GECOS --shell /bin/bash --disabled-password --home /home/stash
COPY --from=prep /stash /ffmpeg/ffmpeg /ffmpeg/ffprobe /usr/bin/
EXPOSE 9999
CMD ["stash"]
