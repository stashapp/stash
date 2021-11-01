# Introduction

This dockerfile is used to build a stash docker container using the current source code. This is ideal for testing your current branch in docker. Note that it does not include python, so python-based scrapers will not work in this image. The production docker images distributed by the project contain python and the necessary packages.

# Building the docker container

From the top-level directory (should contain `main.go` file):

```
make docker-build

```

# Running
The following command should be tweaked to update paths as necessary. You can omit all of the volumes if you'd like to run a fresh instance that doesn't persist your changes.

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