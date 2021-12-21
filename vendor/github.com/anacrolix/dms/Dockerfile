FROM golang

RUN \
  apt-get update && \
  DEBIAN_FRONTEND=noninteractive \
    apt-get install -y --no-install-recommends \
      ffmpeg \
      ffmpegthumbnailer \
  && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/* && \
  touch /root/.dms-ffprobe-cache

COPY . /go/src/github.com/anacrolix/dms/
WORKDIR /go/src/github.com/anacrolix/dms/
RUN \
  go build -v .

ENTRYPOINT [ "./dms" ]
