# syntax=docker/dockerfile:1

ARG GO_VERSION=1.24.3

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-bookworm AS builder

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y \
    ca-certificates \
    curl \
    build-essential \
    git \
    gnupg \
    python3 \
    python3-pip \
    tzdata && \
  rm -rf /var/lib/apt/lists/*

RUN <<'EOF'
set -eux
arch="$(dpkg --print-architecture)"
case "$arch" in
  arm64) goarch=arm64 ;;
  amd64) goarch=amd64 ;;
  armhf) goarch=armv6l ;;
  *) echo "unsupported architecture: $arch"; exit 1 ;;
esac
EOF

RUN <<'EOF'
set -eux
mkdir -p /usr/share/keyrings
curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /usr/share/keyrings/nodesource.gpg
echo "deb [signed-by=/usr/share/keyrings/nodesource.gpg] https://deb.nodesource.com/node_20.x nodistro main" > /etc/apt/sources.list.d/nodesource.list
apt-get update
apt-get install -y nodejs
corepack enable
corepack prepare yarn@1.22.19 --activate
rm -rf /var/lib/apt/lists/*
EOF

ENV PATH="/usr/local/go/bin:${PATH}" \
    CGO_ENABLED=1

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 安装依赖
RUN cd ui/v2.5 && yarn install --frozen-lockfile

# 先生成 GraphQL 类型/ hooks 等
RUN cd ui/v2.5 && yarn run gqlgen

# 再构建
RUN cd ui/v2.5 && yarn build

# 为 Go 后端生成 GraphQL 代码
RUN go run github.com/99designs/gqlgen generate

RUN <<'EOF'
set -eux
GOOS="${TARGETOS:-linux}"
GOARCH="${TARGETARCH}"
GOARM=""
if [ "${GOARCH}" = "arm" ]; then
  case "${TARGETVARIANT}" in
    v7) GOARM=7 ;;
    v6) GOARM=6 ;;
  esac
fi
CGO_ENABLED=1 GOOS="${GOOS}" GOARCH="${GOARCH}" GOARM="${GOARM}" \
  go build -trimpath -ldflags="-s -w" -o /out/stash ./cmd/stash
EOF

FROM ubuntu:24.04

ARG DEBIAN_FRONTEND="noninteractive"

ENV \
  HOME="/root" \
  TZ="Etc/UTC" \
  STASH_CONFIG_FILE="/root/.stash/config.yml" \
  PY_VENV="/pip-install/venv" \
  PATH="/pip-install/venv/bin:$PATH" \
  LIBVA_DRIVER_NAME="rkmpp" \
  LD_LIBRARY_PATH="/usr/lib/jellyfin-ffmpeg:/usr/local/lib"

RUN touch /var/mail/ubuntu && chown ubuntu /var/mail/ubuntu && userdel -r ubuntu

RUN apt-get update && apt-get install -y \
    apt-utils \
    locales && \
  rm -rf /var/lib/apt/lists/* && \
  sed -i '/en_US.UTF-8/s/^# //g' /etc/locale.gen && \
  locale-gen

ENV LANG=en_US.UTF-8
ENV LANGUAGE=en_US:en
ENV LC_ALL=en_US.UTF-8

RUN apt-get update && apt-get install -y --no-install-recommends --no-install-suggests \
    gnupg \
    ca-certificates \
    libvips-tools \
    python3 \
    python3-pip \
    python3-venv \
    tzdata \
    wget \
    curl \
    yq && \
  rm -rf /var/lib/apt/lists/*

# Jellyfin APT 源（用 keyring + signed-by）
RUN set -eux; \
  install -m 0755 -d /usr/share/keyrings; \
  curl -fsSL https://repo.jellyfin.org/jellyfin_team.gpg.key \
    | gpg --dearmor -o /usr/share/keyrings/jellyfin.gpg; \
  . /etc/os-release; \
  echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/jellyfin.gpg] https://repo.jellyfin.org/${ID} ${VERSION_CODENAME} main" \
    > /etc/apt/sources.list.d/jellyfin.list

# Rockchip 多媒体 PPA（提供 librga2 / librga-dev / mpp / libv4l-rkmpp）
RUN set -eux; \
  curl -fsSL "https://keyserver.ubuntu.com/pks/lookup?op=get&search=0x0B2F0747E3BD546820A639B68065BE1FC67AABDE" \
    | gpg --dearmor -o /usr/share/keyrings/rkmm.gpg; \
  . /etc/os-release; \
  echo "deb [signed-by=/usr/share/keyrings/rkmm.gpg] http://ppa.launchpadcontent.net/liujianfeng1994/rockchip-multimedia/ubuntu ${VERSION_CODENAME} main" \
    > /etc/apt/sources.list.d/rockchip-multimedia.list

RUN set -eux; \
  apt-get update; \
  apt-get install -y --no-install-recommends --no-install-suggests \
    jellyfin-ffmpeg7 \
    libdrm2 \
    libdrm-dev \
    libopencl1 \
    ocl-icd-opencl-dev \
    librga2 \
    librockchip-mpp1 \
    libv4l-rkmpp; \
  rm -rf /var/lib/apt/lists/*

COPY --from=builder /out/stash /usr/bin/stash

RUN useradd -u 1000 -U -d /config -s /bin/false stash && \
  usermod -G users stash && \
  usermod -G video stash

RUN apt-get purge -qq wget gnupg curl apt-utils && \
  apt-get autoremove -qq && \
  apt-get clean -qq && \
  rm -rf /tmp/* /var/lib/apt/lists/* /var/tmp/* /var/log/*

ENV PATH="${PATH}:/usr/lib/jellyfin-ffmpeg"

COPY --chmod=755 entrypoint.sh /usr/local/bin/entrypoint.sh

EXPOSE 9999
ENTRYPOINT ["entrypoint.sh"]