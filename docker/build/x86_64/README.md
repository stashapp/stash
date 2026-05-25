# Introduction

This directory contains Dockerfiles for building Stash container images using the current source code. There are three variants:

| Dockerfile | Built-in GPU Support | Use Case |
|---|---|---|
| `Dockerfile` | None | Default — alpine-based frontend/backend, alpine final image, no GPU drivers |
| `Dockerfile-debian` | Intel/AMD | Debian-based final image with VA-API support (mesa-va-drivers + intel-media-va-driver-non-free) |
| `Dockerfile-CUDA` | NVIDIA | NVIDIA GPU support with CUDA, NVENC patch, and VA-API for Intel GPUs |

# Building the docker containers

From the top-level directory (should contain `tools.go` file):

```
# Default variant (alpine-based, no GPU)
make docker-build

# Debian-based variant with high Python binary compatibility and built-in Intel/AMD drivers
make docker-build-debian

# NVIDIA CUDA variant
make docker-cuda-build
```

# Running the docker containers

## Using docker-compose

See the `README.md` file in `docker/production` for instructions on how to get docker-compose if needed.

The `stash/build` container can be run with the `docker-compose.yml` file in `docker/production` by changing the `image` value to the appropriate variant:

- `stash/build` — default variant (alpine-based, no GPU)
- `stash/build-debian` — Debian-based variant with VA-API GPU support
- `stash/cuda-build` — NVIDIA GPU variant

See the instructions in `docker/production` for how to run docker-compose.

## Using `docker run`

After building the container:

```
docker run \
 -e STASH_STASH=/data/ \
 -e STASH_METADATA=/metadata/ \
 -e STASH_CACHE=/cache/ \
 -e STASH_GENERATED=/generated/ \
 -v <path to config dir>:/root/.stash \
 -v <path to media>:/data \
 -v <path to metadata>:/metadata \
 -v <path to cache>:/cache \
 -v <path to generated>:/generated \
 -p 9999:9999 \
 stash/build:latest 
```

Change the `<xxx>` to the appropriate paths. Note that the `<path to media>` directory should be separate from the cache, generated and metadata directories. It is recommended to have the cache, generated and metadata directories in the same parent directory, for example:

```
/stash
  /config
  /metadata
  /generated
  /cache
/media
```

Using this example directory structure, the above command would be:

```
docker run \
 -e STASH_STASH=/data/ \
 -e STASH_METADATA=/metadata/ \
 -e STASH_CACHE=/cache/ \
 -e STASH_GENERATED=/generated/ \
 -v /stash/config:/root/.stash \
 -v /media:/data \
 -v /stash/metadata:/metadata \
 -v /stash/cache:/cache \
 -v /stash/generated:/generated \
 -p 9999:9999 \
 stash/build:latest 
```
