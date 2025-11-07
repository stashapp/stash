# Building from Source

## Pre-requisites

* [Go](https://golang.org/dl/)
* [GolangCI](https://golangci-lint.run/) - A meta-linter which runs several linters in parallel
  * To install, follow the [local installation instructions](https://golangci-lint.run/welcome/install/#local-installation)
* [nodejs](https://nodejs.org/en/download) - nodejs runtime
  * corepack/[pnpm](https://pnpm.io/installation) - nodejs package manager (included with nodejs)

## Environment

### Windows

1. Download and install [Go for Windows](https://golang.org/dl/)
2. Download and extract [MinGW64](https://sourceforge.net/projects/mingw-w64/files/) (scroll down and select x86_64-posix-seh, don't use the autoinstaller, it doesn't work)
3. Search for "Advanced System Settings" and open the System Properties dialog.
    1. Click the `Environment Variables` button
    2. Under System Variables find `Path`. Edit and add `C:\MinGW\bin` (replace with the correct path to where you extracted MingW64).

NOTE: The `make` command in Windows will be `mingw32-make` with MinGW. For example, `make pre-ui` will be `mingw32-make pre-ui`.

### macOS

1. If you don't have it already, install the [Homebrew package manager](https://brew.sh).
2. Install dependencies: `brew install go git gcc make node ffmpeg`

### Linux

#### Arch Linux

1. Install dependencies: `sudo pacman -S go git gcc make nodejs ffmpeg --needed`

#### Ubuntu

1. Install dependencies: `sudo apt-get install golang git gcc nodejs ffmpeg -y`

### OpenBSD

1. Install dependencies `doas pkg_add gmake go git node cmake ffmpeg`
2. Follow the instructions below to build a release, but replace the final step `make build-release` with `gmake flags-release stash`, to [avoid the PIE buildmode](https://github.com/golang/go/issues/59866).

NOTE: The `make` command in OpenBSD will be `gmake`. For example, `make pre-ui` will be `gmake pre-ui`.

## Commands

* `make pre-ui` - Installs the UI dependencies. This only needs to be run once after cloning the repository, or if the dependencies are updated.
* `make generate` - Generates Go and UI GraphQL files. Requires `make pre-ui` to have been run.
* `make generate-stash-box-client` - Generate Go files for the Stash-box client code.
* `make ui` - Builds the UI. Requires `make pre-ui` to have been run.
* `make stash` - Builds the `stash` binary (make sure to build the UI as well... see below)
* `make stash-macapp` - Builds the `Stash.app` macOS app (only works when on macOS, for cross-compilation see below)
* `make phasher` - Builds the `phasher` binary
* `make build` - Builds both the `stash` and `phasher` binaries, alias for `make stash phasher`
* `make build-release` - Builds release versions (debug information removed) of both the `stash` and `phasher` binaries, alias for `make flags-release flags-pie build`
* `make docker-build` - Locally builds and tags a complete 'stash/build' docker image
* `make docker-cuda-build` - Locally builds and tags a complete 'stash/cuda-build' docker image
* `make validate` - Runs all of the tests and checks required to submit a PR
* `make lint` - Runs `golangci-lint` on the backend
* `make it` - Runs all unit and integration tests
* `make fmt` - Formats the Go source code
* `make fmt-ui` - Formats the UI source code
* `make validate-ui` - Runs tests and checks for the UI only
* `make fmt-ui-quick` - (experimental) Formats only changed UI source code
* `make validate-ui-quick` - (experimental) Runs tests and checks of changed UI code
* `make server-start` - Runs a development stash server in the `.local` directory
* `make server-clean` - Removes the `.local` directory and all of its contents
* `make ui-start` - Runs the UI in development mode. Requires a running Stash server to connect to - the server URL can be changed from the default of `http://localhost:9999` using the environment variable `VITE_APP_PLATFORM_URL`, but keep in mind that authentication cannot be used since the session authorization cookie cannot be sent cross-origin. The UI runs on port `3000` or the next available port.

When building, you can optionally prepend `flags-*` targets to the target list in your `make` command to use different build flags:

* `flags-release` (e.g. `make flags-release stash`) - Remove debug information from the binary.
* `flags-pie` (e.g. `make flags-pie build`) - Build a PIE (Position Independent Executable) binary. This provides increased security, but it is unsupported on some systems (notably 32-bit ARM and OpenBSD).
* `flags-static` (e.g. `make flags-static phasher`) - Build a statically linked binary (the default is a dynamically linked binary).
* `flags-static-pie` (e.g. `make flags-static-pie stash`) - Build a statically linked PIE binary (using `flags-static` and `flags-pie` separately will not work).
* `flags-static-windows` (e.g. `make flags-static-windows build`) - Identical to `flags-static-pie`, but does not enable the `netgo` build tag, which is not needed for static builds on Windows.

## Local development quickstart

1. Run `make pre-ui` to install UI dependencies
2. Run `make generate` to create generated files
3. In one terminal, run `make server-start` to run the server code
4. In a separate terminal, run `make ui-start` to run the UI in development mode
5. Open the UI in a browser: `http://localhost:3000/`

Changes to the UI code can be seen by reloading the browser page.

Changes to the backend code require a server restart (`CTRL-C` in the server terminal, followed by `make server-start` again) to be seen.

On first launch:

1. On the "Stash Setup Wizard" screen, choose a directory with some files to test with
2. Press "Next" to use the default locations for the database and generated content
3. Press the "Confirm" and "Finish" buttons to get into the UI
4. On the side menu, navigate to "Tasks -> Library -> Scan" and press the "Scan" button
5. You're all set! Set any other configurations you'd like and test your code changes.

To start fresh with new configuration:

1. Stop the server (`CTRL-C` in the server terminal)
2. Run `make server-clean` to clear all config, database, and generated files (under `.local`)
3. Run `make server-start` to restart the server
4. Follow the "On first launch" steps above

## Building a release

Simply run `make` or `make release`, or equivalently:

1. Run `make pre-ui` to install UI dependencies
2. Run `make generate` to create generated files
3. Run `make ui` to build the frontend
4. Run `make build-release` to build a release executable for your current platform

## Cross-compiling

This project uses a modification of the [CI-GoReleaser](https://github.com/bep/dockerfiles/tree/master/ci-goreleaser) Docker container for cross-compilation, defined in `docker/compiler/Dockerfile`.

To cross-compile the app yourself:

1. Run `make pre-ui`, `make generate` and `make ui` outside the container, to generate files and build the UI.
2. Pull the latest compiler image from Docker Hub: `docker pull stashapp/compiler`
3. Run `docker run --rm --mount type=bind,source="$(pwd)",target=/stash -w /stash -it stashapp/compiler /bin/bash` to open a shell inside the container.
4. From inside the container, run `make build-cc-all` to build for all platforms, or run `make build-cc-{platform}` to build for a specific platform (have a look at the `Makefile` for the list of targets).
5. You will find the compiled binaries in `dist/`.

NOTE: Since the container is run as UID 0 (root), the resulting binaries (and the `dist/` folder itself, if it had to be created) will be owned by root.

## Profiling

Stash can be profiled using the `--cpuprofile <output profile filename>` command line flag.

The resulting file can then be used with pprof as follows:

`go tool pprof <path to binary> <path to profile filename>`

With `graphviz` installed and in the path, a call graph can be generated with:

`go tool pprof -svg <path to binary> <path to profile filename> > <output svg file>`
