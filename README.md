# Stash

[![Build Status](https://travis-ci.org/stashapp/stash.svg?branch=master)](https://travis-ci.org/stashapp/stash)
[![Go Report Card](https://goreportcard.com/badge/github.com/stashapp/stash)](https://goreportcard.com/report/github.com/stashapp/stash)
[![Discord](https://img.shields.io/discord/559159668438728723.svg?logo=discord)](https://discord.gg/2TsNFKt)

https://stashapp.cc

**Stash is a locally hosted web-based app written in Go which organizes and serves your porn.**

* It can gather information about videos in your collection from the internet, and is extensible through the use of community-built plugins for a large number of content producers.
* It supports a wide variety of both video and image formats.
* You can tag videos and find them later.
* It provides statistics about performers, tags, studios and other things.

You can [watch a SFW demo video](https://vimeo.com/545323354) to see it in action.

For further information you can [read the in-app manual](ui/v2.5/src/docs/en).

# Installing stash

## via Docker

Follow [this README.md in the docker directory.](docker/production/README.md)

## Pre-Compiled Binaries

The Stash server runs on macOS, Windows, and Linux.  Download the [latest release here](https://github.com/stashapp/stash/releases).

Run the executable (double click the exe on windows or run `./stash-osx` / `./stash-linux` from the terminal on macOS / Linux) and navigate to either https://localhost:9999 or http://localhost:9999 to get started.

*Note for Windows users:* Running the app might present a security prompt since the binary isn't yet signed.  Bypass this by clicking "more info" and then the "run anyway" button.

#### FFMPEG

If stash is unable to find or download FFMPEG then download it yourself from the link for your platform:

* [macOS ffmpeg](https://evermeet.cx/ffmpeg/ffmpeg-4.3.1.zip), [macOS ffprobe](https://evermeet.cx/ffmpeg/ffprobe-4.3.1.zip)
* [Windows](https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-essentials.zip)
* [Linux](https://www.johnvansickle.com/ffmpeg/)

The `ffmpeg(.exe)` and `ffprobe(.exe)` files should be placed in `~/.stash` on macOS / Linux or `C:\Users\YourUsername\.stash` on Windows.

# Usage

## Quickstart Guide
1) Download and install Stash and its dependencies
2) Run Stash. It will prompt you for some configuration options and a directory to index (you can also do this step afterward)
3) After configuration, launch your web browser and navigate to the URL shown within the Stash app.

**Note that Stash does not currently retrieve and organize information about your entire library automatically.**  You will need to help it along through the use of [scrapers](blob/develop/ui/v2.5/src/docs/en/Scraping.md).  The Stash community has developed scrapers for many popular data sources which can be downloaded and installed from [this repository](https://github.com/stashapp/CommunityScrapers).

The simplest way to tag a large number of files is by using the [Tagger](https://github.com/stashapp/stash/blob/develop/ui/v2.5/src/docs/en/Tagger.md) which uses filename keywords to help identify the file and pull in scene and performer information from our stash-box database. Note that this data source is not comprehensive and you may need to use the scrapers to identify some of your media.

## CLI

Stash runs as a command-line app and local web server.  There are some command-line options available, which you can see by running `stash --help`.

For example, to run stash locally on port 80 run it like this (OSX / Linux) `stash --host 127.0.0.1 --port 80`

## SSL (HTTPS)

Stash can run over HTTPS with some additional work.  First you must generate a SSL certificate and key combo.  Here is an example using openssl:

`openssl req -x509 -newkey rsa:4096 -sha256 -days 7300 -nodes -keyout stash.key -out stash.crt -extensions san -config <(echo "[req]"; echo distinguished_name=req; echo "[san]"; echo subjectAltName=DNS:stash.server,IP:127.0.0.1) -subj /CN=stash.server`

This command would need customizing for your environment.  [This link](https://stackoverflow.com/questions/10175812/how-to-create-a-self-signed-certificate-with-openssl) might be useful.

Once you have a certificate and key file name them `stash.crt` and `stash.key` and place them in the same directory as the `config.yml` file, or the `~/.stash` directory.  Stash detects these and starts up using HTTPS rather than HTTP.

# Customization

## Themes and CSS Customization
There is a [directory of community-created themes](https://github.com/stashapp/stash/wiki/Themes) on our Wiki, along with instructions on how to install them.

You can also make Stash interface fit your desired style with [Custom CSS snippets](https://github.com/stashapp/stash/wiki/Custom-CSS-snippets) and [CSS Tweaks](https://github.com/stashapp/stash/wiki/CSS-Tweaks).

# Support (FAQ)

Answers to other Frequently Asked Questions can be found [on our Wiki](https://github.com/stashapp/stash/wiki/FAQ)

For issues not addressed there, there are a few options.

* Read the [Wiki](https://github.com/stashapp/stash/wiki)
* Check the in-app documentation (also available [here](https://github.com/stashapp/stash/tree/develop/ui/v2.5/src/docs/en)
* Join the [Discord server](https://discord.gg/2TsNFKt), where the community can offer support.

# Compiling From Source Code

## Pre-requisites

* [Go](https://golang.org/dl/)
* [Revive](https://github.com/mgechev/revive) - Configurable linter
    * Go Install: `go get github.com/mgechev/revive`
* [Packr2](https://github.com/gobuffalo/packr/) - Static asset bundler
    * Go Install: `go get github.com/gobuffalo/packr/v2/packr2`
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
    2. Under system variables find the `Path`.  Edit and add `C:\Program Files\mingw-w64\*\mingw64\bin` (replace * with the correct path).

NOTE: The `make` command in Windows will be `mingw32-make` with MingW.

## Commands

* `make generate` - Generate Go and UI GraphQL files
* `make build` - Builds the binary (make sure to build the UI as well... see below)
* `make docker-build` - Locally builds and tags a complete 'stash/build' docker image
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
* `make ui-start` - Runs the UI in development mode. Requires a running stash server to connect to. Stash port can be changed from the default of `9999` with environment variable `REACT_APP_PLATFORM_PORT`.

## Building a release

1. Run `make generate` to create generated files 
2. Run `make ui` to compile the frontend
3. Run `make build` to build the executable for your current platform

## Cross compiling

This project uses a modification of the [CI-GoReleaser](https://github.com/bep/dockerfiles/tree/master/ci-goreleaser) docker container to create an environment
where the app can be cross-compiled.  This process is kicked off by CI via the `scripts/cross-compile.sh` script.  Run the following
command to open a bash shell to the container to poke around:

`docker run --rm --mount type=bind,source="$(pwd)",target=/stash -w /stash -i -t stashappdev/compiler:latest /bin/bash`

## Profiling

Stash can be profiled using the `--cpuprofile <output profile filename>` command line flag.

The resulting file can then be used with pprof as follows: 

`go tool pprof <path to binary> <path to profile filename>`

With `graphviz` installed and in the path, a call graph can be generated with:

`go tool pprof -svg <path to binary> <path to profile filename> > <output svg file>`
