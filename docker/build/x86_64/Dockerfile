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
FROM golang:1.22-bullseye as backend
RUN apt update && apt install -y build-essential golang
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
FROM ubuntu:22.04
RUN apt update && apt upgrade -y && apt install -y --no-install-recommends curl gnupg ca-certificates libvips-tools wget intel-media-va-driver-non-free vainfo

# Optional: Install google-chrome
# RUN wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb && apt install -y ./google-chrome-stable_current_amd64.deb && apt --fix-broken install

# Add Jellyfin GPG key
RUN curl -fsSL https://repo.jellyfin.org/debian/jellyfin_team.gpg.key | gpg --dearmor --batch --yes -o /etc/apt/trusted.gpg.d/debian-jellyfin.gpg

# Add Jellyfin repository
RUN os_id=$(awk -F'=' '/^ID=/{ print $NF }' /etc/os-release) && \
    os_codename=$(awk -F'=' '/^VERSION_CODENAME=/{ print $NF }' /etc/os-release) && \
    echo "deb [arch=$(dpkg --print-architecture)] https://repo.jellyfin.org/$os_id $os_codename main" > /etc/apt/sources.list.d/jellyfin.list

# Update package lists and install Jellyfin FFmpeg
RUN apt-get update && \
    apt-get install --no-install-recommends --no-install-suggests -y jellyfin-ffmpeg6 && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# Link ffmpeg to PATH
RUN ln -s /usr/lib/jellyfin-ffmpeg/ffmpeg /usr/local/bin/ffmpeg && \
    ln -s /usr/lib/jellyfin-ffmpeg/ffprobe /usr/local/bin/ffprobe

COPY --from=backend /stash/stash /usr/bin/
ENV STASH_CONFIG_FILE=/root/.stash/config.yml
EXPOSE 9999
ENTRYPOINT ["stash"]

# vim: ft=dockerfile
