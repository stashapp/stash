FROM golang:1.16-alpine as builder

COPY mockery /usr/local/bin

# Explicitly set a writable cache path when running --user=$(id -u):$(id -g)
# see: https://github.com/golang/go/issues/26280#issuecomment-445294378
ENV GOCACHE /tmp/.cache

ENTRYPOINT ["/usr/local/bin/mockery"]
