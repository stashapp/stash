# This dockerfile should be built with `make docker-build` from the stash root.

# Build Frontend
FROM node:alpine as frontend
RUN apk add --no-cache make git
## cache node_modules separately
COPY ./ui/v2.5/package.json ./ui/v2.5/yarn.lock /stash/ui/v2.5/
WORKDIR /stash
COPY Makefile /stash/
COPY ./graphql /stash/graphql/
COPY ./ui /stash/ui/
RUN make pre-ui
RUN make generate-ui
ARG GITHASH
ARG STASH_VERSION
RUN BUILD_DATE=$(date +"%Y-%m-%d %H:%M:%S") make ui

# Build Backend
FROM golang:1.22-alpine as backend
RUN apk add --no-cache make alpine-sdk
WORKDIR /stash
COPY ./go* ./*.go Makefile gqlgen.yml .gqlgenc.yml /stash/
COPY ./scripts /stash/scripts/
COPY ./pkg /stash/pkg/
COPY ./cmd /stash/cmd
COPY ./internal /stash/internal
COPY --from=frontend /stash /stash/
RUN make generate-backend
ARG GITHASH
ARG STASH_VERSION
RUN make flags-release flags-pie stash

# Final Runnable Image
FROM alpine:latest
RUN apk add --no-cache --virtual .build-deps gcc python3-dev musl-dev \
    && apk add --no-cache ca-certificates python3 py3-requests py3-requests-toolbelt py3-lxml py3-pip ffmpeg ruby tzdata \
    && apk add --no-cache --repository http://dl-cdn.alpinelinux.org/alpine/v3.18/community vips=8.14.3-r0 vips-tools=8.14.3-r0 \
    && pip install --user --break-system-packages mechanicalsoup cloudscraper bencoder.pyx \
    && gem install faraday \
    && apk del .build-deps
COPY --from=backend /stash/stash /usr/bin/
ENV STASH_CONFIG_FILE=/root/.stash/config.yml
EXPOSE 9999
ENTRYPOINT ["stash"]
