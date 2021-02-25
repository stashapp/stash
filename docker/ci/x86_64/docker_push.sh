#!/bin/bash

DOCKER_TAG=$1
echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin

# must build the image from dist directory
docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 --push --output type=image,name=stashapp/stash:$DOCKER_TAG,push=true -f docker/ci/x86_64/Dockerfile dist/

