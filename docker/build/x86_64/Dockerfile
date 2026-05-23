# This dockerfile should be built with `make docker-build` from the stash root.

# Build Frontend
FROM node:24-alpine AS frontend
RUN apk add --no-cache make git
## cache node_modules separately
COPY ./ui/v2.5/package.json ./ui/v2.5/pnpm-lock.yaml /stash/ui/v2.5/
WORKDIR /stash
COPY Makefile /stash/
COPY ./graphql /stash/graphql/
COPY ./ui /stash/ui/
# pnpm install with npm
RUN npm install -g pnpm
RUN make pre-ui
RUN make generate-ui
ARG GITHASH
ARG STASH_VERSION
RUN BUILD_DATE=$(date +"%Y-%m-%d %H:%M:%S") make ui-only

# Build Backend
FROM golang:1.25.9-alpine AS backend
RUN apk add --no-cache make alpine-sdk
WORKDIR /stash
COPY ./go* ./*.go Makefile gqlgen.yml .gqlgenc.yml /stash/
COPY ./graphql /stash/graphql/
COPY ./scripts /stash/scripts/
COPY ./pkg /stash/pkg/
COPY ./cmd /stash/cmd/
COPY ./internal /stash/internal/
# needed for generate-login-locale
COPY ./ui /stash/ui/
RUN make generate-backend generate-login-locale
COPY --from=frontend /stash /stash/
ARG GITHASH
ARG STASH_VERSION
RUN make flags-release flags-pie stash

# Final Runnable Image
FROM alpine:latest
RUN apk add --no-cache ca-certificates vips-tools ffmpeg

# Create non-privileged user for security best practices
# PUID/PGID can be overridden at build time if needed
ARG PUID=1000
ARG PGID=1000
RUN addgroup -g ${PGID} stash && \
    adduser -D -G stash -u ${PUID} stash

COPY --from=backend /stash/stash /usr/bin/

# Create default config directory with proper ownership
RUN mkdir -p /config && chown stash:stash /config

ENV STASH_CONFIG_FILE=/config/config.yml
ENV PUID=1000
ENV PGID=1000

EXPOSE 9999
USER stash
ENTRYPOINT ["stash"]
