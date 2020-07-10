# Stash

[![Build Status](https://travis-ci.org/stashapp/stash.svg?branch=master)](https://travis-ci.org/stashapp/stash)
[![Go Report Card](https://goreportcard.com/badge/github.com/stashapp/stash)](https://goreportcard.com/report/github.com/stashapp/stash)
[![Discord](https://img.shields.io/discord/559159668438728723.svg?logo=discord)](https://discord.gg/2TsNFKt)

https://stashapp.cc

**Stash is a Go app which organizes and serves your porn.**

See a demo [here](https://vimeo.com/275537038) (password is stashapp).

An in-app manual is available, and the manual pages can be viewed [here](https://github.com/stashapp/stash/tree/develop/ui/v2.5/src/docs/en).

# Docker install

Follow [this README.md in the docker directory.](docker/production/README.md)

# Bare-metal Install

Stash supports macOS, Windows, and Linux.  Download the [latest release here](https://github.com/stashapp/stash/releases).

Run the executable (double click the exe on windows or run `./stash-osx` / `./stash-linux` from the terminal on macOS / Linux) and navigate to either https://localhost:9999 or http://localhost:9999 to get started.

*Note for Windows users:* Running the app might present a security prompt since the binary isn't signed yet.  Just click more info and then the "run anyway" button.

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

This command would need customizing for your environment.  [This link](https://stackoverflow.com/questions/10175812/how-to-create-a-self-signed-certificate-with-openssl) might be useful.

Once you have a certificate and key file name them `stash.crt` and `stash.key` and place them in the `~/.stash` directory.  Stash detects these and starts up using HTTPS rather than HTTP.

# FAQ

> I'm unable to run the app on OSX or Linux

Try running `chmod u+x stash-osx` or `chmod u+x stash-linux` to make the file executable.

> I have a question not answered here.

Join the [Discord server](https://discord.gg/2TsNFKt).

# Development

## Install

* [Go](https://golang.org/dl/)
* [Revive](https://github.com/mgechev/revive) - Configurable linter
    * Go Install: `go get github.com/mgechev/revive`
* [Packr2](https://github.com/gobuffalo/packr/tree/v2.0.2/v2) - Static asset bundler
    * Go Install: `go get github.com/gobuffalo/packr/v2/packr2@v2.0.2`
    * [Binary Download](https://github.com/gobuffalo/packr/releases)
* [Yarn](https://yarnpkg.com/en/docs/install) - Yarn package manager
    * Run `yarn install --frozen-lockfile` in the `stash/ui/v2.5` folder (before running make generate for first time).

NOTE: You may need to run the `go get` commands outside the project directory to avoid modifying the projects module file.

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

NOTE: The `make` command in Windows will be `mingw32-make` with MingW.

## Commands

* `make generate` - Generate Go and UI GraphQL files
* `make build` - Builds the binary (make sure to build the UI as well... see below)
* `make pre-ui` - Installs the UI dependencies. Only needs to be run once before building the UI for the first time, or if the dependencies are updated
* `make fmt-ui` - Formats the UI source code.
* `make ui` - Builds the frontend and the packr2 files
* `make packr` - Generate packr2 files (sub-target of `ui`. Use to regenerate packr2 files without rebuilding UI)
* `make vet` - Run `go vet`
* `make lint` - Run the linter
* `make fmt` - Run `go fmt`
* `make fmt-check` - Ensure changed files are formatted correctly
* `make it` - Run the unit and integration tests
* `make validate` - Run all of the tests and checks required to submit a PR

## Building a release

1. Run `make generate` to create generated files 
2. Run `make ui` to compile the frontend
3. Run `make build` to build the executable for your current platform

## Cross compiling

This project uses a modification of [this](https://github.com/bep/dockerfiles/tree/master/ci-goreleaser) docker container to create an environment
where the app can be cross-compiled.  This process is kicked off by CI via the `scripts/cross-compile.sh` script.  Run the following
command to open a bash shell to the container to poke around:

`docker run --rm --mount type=bind,source="$(pwd)",target=/stash -w /stash -i -t stashappdev/compiler:latest /bin/bash`

## Customization

You can make Stash interface fit your desired style with  [Custom CSS snippets](https://github.com/stashapp/stash/wiki/Custom-CSS-snippets) and [CSS Tweaks](https://github.com/stashapp/stash/wiki/CSS-Tweaks).

[Stash Plex Theme](https://github.com/stashapp/stash/wiki/Stash-Plex-Theme) is a community created theme inspired by popular Plex Interface.

