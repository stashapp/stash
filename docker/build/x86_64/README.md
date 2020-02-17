# Introduction

This dockerfile is used to build a stash docker container using the current source code.

# Building the docker container

From the top-level directory (should contain `main.go` file):

```
docker build -t stash/build -f ./docker/build/x86_64/Dockerfile .

```

# Running the docker container

## Using docker-compose

See the `README.md` file in `docker/production` for instructions on how to get docker-compose if needed.

The `stash/build` container can be run with the `docker-compose.yml` file in `docker/production` by changing the `image` value to be `stash/build`. See the instructions in `docker/production` for how to run docker-compose.

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
 stash/build:latest 
```
