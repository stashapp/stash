#!/bin/bash

DOCKER_TAG=$1

# must build the image from dist directory
echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v6,linux/arm/v7 --push --output type=registry,name=tweeticoats/stash:$DOCKER_TAG,push=true -f docker/ci/x86_64/Dockerfile dist/

