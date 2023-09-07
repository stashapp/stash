#!/bin/bash

echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin

# Build the main image from the dist directory
docker buildx build \
    --platform linux/amd64,linux/arm64,linux/arm/v7,linux/arm/v6 --push \
    $(for TAG in "$@"; do echo -n "-t stashapp/stash:$TAG "; done) \
    -f docker/ci/x86_64/Dockerfile dist/ --target app

# Build the CUDA image from the dist directory
docker buildx build \
    --platform linux/amd64,linux/arm64 --push \
    $(for TAG in "$@"; do echo -n "-t stashapp/stash:$TAG-cuda "; done) \
    -f docker/ci/x86_64/Dockerfile dist/ --target cuda_app
