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
FROM golang:1.24.3-alpine AS backend
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
COPY --from=backend /stash/stash /usr/bin/
ENV STASH_CONFIG_FILE=/root/.stash/config.yml
EXPOSE 9999
ENTRYPOINT ["stash"]
