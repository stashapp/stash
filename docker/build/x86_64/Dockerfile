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
RUN BUILD_DATE=$(date +"%Y-%m-%d %H:%M:%S") make ui-only

# Build Backend
FROM golang:1.17-alpine as backend
RUN apk add --no-cache xz make alpine-sdk
## install packr, ffmpeg
ENV PACKR2_VERSION=2.8.1
ENV PACKR2_SHA=1cb2a0113550bc7962a8fda31a29877fcbbabd56b46c25e1fffbc225334162e7
ENV PACKR2_DOWNLOAD_FILE=packr_${PACKR2_VERSION}_linux_amd64.tar.gz
ENV PACKR2_DOWNLOAD_URL=https://github.com/gobuffalo/packr/releases/download/v${PACKR2_VERSION}/${PACKR2_DOWNLOAD_FILE}
WORKDIR /
RUN wget ${PACKR2_DOWNLOAD_URL}; \
	echo "$PACKR2_SHA  $PACKR2_DOWNLOAD_FILE" | sha256sum -c - || exit 1; \
	tar -xzf $PACKR2_DOWNLOAD_FILE -C /usr/bin/ packr2; \
	rm $PACKR2_DOWNLOAD_FILE;
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
RUN make packr
RUN make build

# Final Runnable Image
FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=backend /stash/stash /ffmpeg/ffmpeg /ffmpeg/ffprobe /usr/bin/
ENV STASH_CONFIG_FILE=/root/.stash/config.yml
EXPOSE 9999
ENTRYPOINT ["stash"]