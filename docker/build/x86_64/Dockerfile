# This dockerfile must be built from the top-level stash directory
# ie from top-level stash:
# docker build -t stash/build -f docker/build/x86_64/Dockerfile .

# Build Frontend
FROM node:alpine as frontend
RUN apk add --no-cache make git
## cache node_modules separately
COPY ./ui/v2.5/package.json ./ui/v2.5/yarn.lock /stash/ui/v2.5/
WORKDIR /stash
RUN yarn --cwd ui/v2.5 install --frozen-lockfile.
COPY Makefile /stash/
COPY ./.git /stash/.git
COPY ./graphql /stash/graphql/
COPY ./ui /stash/ui/
RUN make generate-frontend
RUN BUILD_DATE=$(date +"%Y-%m-%d %H:%M:%S") make ui

# Build Backend
FROM golang:1.17-alpine as backend
RUN apk add --no-cache xz make alpine-sdk
## install ffmpeg
WORKDIR /
RUN wget -O /ffmpeg.tar.xz https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz && \
    tar xf /ffmpeg.tar.xz && \
    rm ffmpeg.tar.xz && \
    mv /ffmpeg*/ /ffmpeg/
WORKDIR /stash
COPY ./go* ./*.go Makefile gqlgen.yml .gqlgenc.yml /stash/
COPY ./static /stash/static/
COPY ./scripts /stash/scripts/
COPY ./vendor /stash/vendor/
COPY ./pkg /stash/pkg/
COPY --from=frontend /stash /stash/
RUN make generate-backend
RUN make build

# Final Runnable Image
FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=backend /stash/stash /ffmpeg/ffmpeg /ffmpeg/ffprobe /usr/bin/
ENV STASH_CONFIG_FILE=/root/.stash/config.yml
EXPOSE 9999
ENTRYPOINT ["stash"]