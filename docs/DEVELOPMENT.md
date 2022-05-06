# Building from Source

## Pre-requisites

* [Go](https://golang.org/dl/)
* [GolangCI](https://golangci-lint.run/) - A meta-linter which runs several linters in parallel
  * To install, follow the [local installation instructions](https://golangci-lint.run/usage/install/#local-installation)
* [Yarn](https://yarnpkg.com/en/docs/install) - Yarn package manager
  * Run `yarn install --frozen-lockfile` in the `stash/ui/v2.5` folder (before running make generate for first time).

NOTE: You may need to run the `go get` commands outside the project directory to avoid modifying the projects module file.

## Environment

### Windows

1. Download and install [Go for Windows](https://golang.org/dl/)
2. Download and install [MingW](https://sourceforge.net/projects/mingw/) and select packages `mingw32-base`
3. Search for "advanced system settings" and open the system properties dialog.
    1. Click the `Environment Variables` button
    2. Under system variables find the `Path`.  Edit and add `C:\MinGW\bin` (replace * with the correct path).

NOTE: The `make` command in Windows will be `mingw32-make` with MingW. For example `make pre-ui` will be `mingw32-make pre-ui`

### macOS

1. If you don't have it already, install the [Homebrew package manager](https://brew.sh).
2. Install dependencies: `brew install go git yarn gcc make`

## Commands

* `make pre-ui` - Installs the UI dependencies. Only needs to be run once before building the UI for the first time, or if the dependencies are updated
* `make generate` - Generate Go and UI GraphQL files
* `make fmt-ui` - Formats the UI source code
* `make ui` - Builds the frontend
* `make build` - Builds the binary (make sure to build the UI as well... see below)
* `make docker-build` - Locally builds and tags a complete 'stash/build' docker image
* `make lint` - Run the linter on the backend
* `make fmt` - Run `go fmt`
* `make it` - Run the unit and integration tests
* `make validate` - Run all of the tests and checks required to submit a PR
* `make ui-start` - Runs the UI in development mode. Requires a running stash server to connect to. Stash server port can be changed from the default of `9999` using environment variable `VITE_APP_PLATFORM_PORT`. UI runs on port `3000` or the next available port.

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
