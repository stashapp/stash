# Stash

[![Build Status](https://travis-ci.org/stashapp/stash.svg?branch=master)](https://travis-ci.org/stashapp/stash)
[![Go Report Card](https://goreportcard.com/badge/github.com/stashapp/stash)](https://goreportcard.com/report/github.com/stashapp/stash)
[![Discord](https://img.shields.io/discord/559159668438728723.svg?logo=discord)](https://discord.gg/2TsNFKt)

**Stash is a Go app which organizes and serves your porn.**

See a demo [here](https://vimeo.com/275537038) (password is stashapp).

# Install

Stash supports macOS, Windows, and Linux.  Download the [latest release here](https://github.com/stashapp/stash/releases).

Simply run the executable (double click the exe on windows or run `./stash-osx` / `./stash-linux` from the terminal on macOS / Linux) and navigate to either https://localhost:9999 or http://localhost:9998 to get started.

*Note for Windows users:* Running the app might present a security prompt since the binary isn't signed yet.  Just click more info and then the run anyway button.

#### FFMPEG

If stash is unable to find or download FFMPEG then download it yourself from the link for your platform:

* [macOS](https://ffmpeg.zeranoe.com/builds/macos64/static/ffmpeg-4.0-macos64-static.zip)
* [Windows](https://ffmpeg.zeranoe.com/builds/win64/static/ffmpeg-4.0-win64-static.zip)
* [Linux](https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz)

The `ffmpeg(.exe)` and `ffprobe(.exe)` files should be placed in `~/.stash` on macOS / Linux or `C:\Users\YourUsername\.stash` on Windows.

# Usage

## CLI

Stash provides some command line options.  See what is currently available by running `stash --help`.

For example, to run stash locally on port 80 run it like this (OSX / Linux) `stash --host 127.0.0.1 --port 80`

## SSL (HTTPS)

Stash supports HTTPS with some additional work.  First you must generate a SSL certificate and key combo.  Here is an example using openssl:

`openssl req -x509 -newkey rsa:4096 -sha256 -days 7300 -nodes -keyout stash.key -out stash.crt -extensions san -config <(echo "[req]"; echo distinguished_name=req; echo "[san]"; echo subjectAltName=DNS:stash.server,IP:127.0.0.1) -subj /CN=stash.server`

This command would need to be customized for your environment.  [This link](https://stackoverflow.com/questions/10175812/how-to-create-a-self-signed-certificate-with-openssl) might be useful.

Once you have a certificate and key file name them `stash.crt` and `stash.key` and place them in the `~/.stash` directory.  Stash will detect these and start up using HTTPS rather than HTTP.

# FAQ

> I'm unable to run the app on OSX or Linux

Try running `chmod u+x stash-osx` or `chmod u+x stash-linux` to make the file executable.

> I have a question not answered here.

Join the [Discord server](https://discord.gg/2TsNFKt).

# Development

## Install

* [Revive](https://github.com/mgechev/revive) - Configurable linter `go get github.com/mgechev/revive`
* [Yarn](https://yarnpkg.com/en/docs/install) - Yarn package manager

## Environment

### macOS

TODO

### Windows

1. Download and install [Go for Windows](https://golang.org/dl/)
2. Download and install [MingW](https://sourceforge.net/projects/mingw-w64/)
3. Search for "advanced system settings" and open the system properties dialog.
	1. Click the `Environment Variables` button
	2. Add `GO111MODULE=on`
	3. Under system variables find the `Path`.  Edit and add `C:\Program Files\mingw-w64\*\mingw64\bin` (replace * with the correct path).

## Commands

* `make build` - Builds the binary
* `make gqlgen` - Regenerate Go GraphQL files
* `make vet` - Run `go vet`
* `make lint` - Run the linter

## Building a release

1. cd into the `ui/v2` directory and run `yarn build` to compile the frontend
2. cd back to the root directory and run `make build` to build the executable for your current platform

## Cross compiling

This project uses a modification of [this](https://github.com/bep/dockerfiles/tree/master/ci-goreleaser) docker container to create an environment
where the app can be cross compiled.  This process is kicked off by CI via the `scripts/cross-compile.sh` script.  Run the following
command to open a bash shell to the container to poke around:

`docker run --rm --mount type=bind,source="$(pwd)",target=/stash -w /stash -i -t stashappdev/compiler:latest /bin/bash`
