IS_WIN_SHELL =
ifeq (${SHELL}, sh.exe)
  IS_WIN_SHELL = true
endif
ifeq (${SHELL}, cmd)
  IS_WIN_SHELL = true
endif

ifdef IS_WIN_SHELL
  RM := del /s /q
  RMDIR := rmdir /s /q
  NOOP := @@
else
  RM := rm -f
  RMDIR := rm -rf
  NOOP := @:
endif

# set LDFLAGS environment variable to any extra ldflags required
LDFLAGS := $(LDFLAGS)

# set OUTPUT environment variable to generate a specific binary name
# this will apply to both `stash` and `phasher`, so build them separately
# alternatively use STASH_OUTPUT or PHASHER_OUTPUT to set the value individually
ifdef OUTPUT
  STASH_OUTPUT := $(OUTPUT)
  PHASHER_OUTPUT := $(OUTPUT)
endif
ifdef STASH_OUTPUT
  STASH_OUTPUT := -o $(STASH_OUTPUT)
endif
ifdef PHASHER_OUTPUT
  PHASHER_OUTPUT := -o $(PHASHER_OUTPUT)
endif

# set GO_BUILD_FLAGS environment variable to any extra build flags required
GO_BUILD_FLAGS := $(GO_BUILD_FLAGS)

# set GO_BUILD_TAGS environment variable to any extra build tags required
GO_BUILD_TAGS := $(GO_BUILD_TAGS)
GO_BUILD_TAGS += sqlite_stat4 sqlite_math_functions

# set STASH_NOLEGACY environment variable or uncomment to disable legacy browser support
# STASH_NOLEGACY := true

# set STASH_SOURCEMAPS environment variable or uncomment to enable UI sourcemaps
# STASH_SOURCEMAPS := true

export CGO_ENABLED := 1

# define COMPILER_IMAGE for cross-compilation docker container
ifndef COMPILER_IMAGE
  COMPILER_IMAGE := stashapp/compiler:latest
endif

.PHONY: release
release: pre-ui generate ui build-release

# targets to set various build flags
# use combinations on the make command-line to configure a build, e.g.:
# for a static-pie release build: `make flags-static-pie flags-release stash`
# for a static windows debug build: `make flags-static-windows stash`

# $(NOOP) prevents "nothing to be done" warnings

.PHONY: flags-release
flags-release:
	$(NOOP)
	$(eval LDFLAGS += -s -w)
	$(eval GO_BUILD_FLAGS += -trimpath)

.PHONY: flags-pie
flags-pie:
	$(NOOP)
	$(eval GO_BUILD_FLAGS += -buildmode=pie)

.PHONY: flags-static
flags-static:
	$(NOOP)
	$(eval LDFLAGS += -extldflags=-static)
	$(eval GO_BUILD_TAGS += sqlite_omit_load_extension osusergo netgo)

.PHONY: flags-static-pie
flags-static-pie:
	$(NOOP)
	$(eval LDFLAGS += -extldflags=-static-pie)
	$(eval GO_BUILD_FLAGS += -buildmode=pie)
	$(eval GO_BUILD_TAGS += sqlite_omit_load_extension osusergo netgo)

# identical to flags-static-pie, but excluding netgo, which is not needed on windows
.PHONY: flags-static-windows
flags-static-windows:
	$(NOOP)
	$(eval LDFLAGS += -extldflags=-static-pie)
	$(eval GO_BUILD_FLAGS += -buildmode=pie)
	$(eval GO_BUILD_TAGS += sqlite_omit_load_extension osusergo)

.PHONY: build-info
build-info:
ifndef BUILD_DATE
	$(eval BUILD_DATE := $(shell go run scripts/getDate.go))
endif
ifndef GITHASH
	$(eval GITHASH := $(shell git rev-parse --short HEAD))
endif
ifndef STASH_VERSION
	$(eval STASH_VERSION := $(shell git describe --tags --exclude latest_develop))
endif
ifndef OFFICIAL_BUILD
	$(eval OFFICIAL_BUILD := false)
endif

.PHONY: build-flags
build-flags: build-info
	$(eval BUILD_LDFLAGS := $(LDFLAGS))
	$(eval BUILD_LDFLAGS += -X 'github.com/stashapp/stash/internal/build.buildstamp=$(BUILD_DATE)')
	$(eval BUILD_LDFLAGS += -X 'github.com/stashapp/stash/internal/build.githash=$(GITHASH)')
	$(eval BUILD_LDFLAGS += -X 'github.com/stashapp/stash/internal/build.version=$(STASH_VERSION)')
	$(eval BUILD_LDFLAGS += -X 'github.com/stashapp/stash/internal/build.officialBuild=$(OFFICIAL_BUILD)')
	$(eval BUILD_FLAGS := -v -tags "$(GO_BUILD_TAGS)" $(GO_BUILD_FLAGS) -ldflags "$(BUILD_LDFLAGS)")

.PHONY: stash
stash: build-flags
	go build $(STASH_OUTPUT) $(BUILD_FLAGS) ./cmd/stash

.PHONY: phasher
phasher: build-flags
	go build $(PHASHER_OUTPUT) $(BUILD_FLAGS) ./cmd/phasher

# builds dynamically-linked debug binaries
.PHONY: build
build: stash phasher

# builds dynamically-linked PIE release binaries
.PHONY: build-release
build-release: flags-release flags-pie build

# compile and bundle into Stash.app
# for when on macOS itself
.PHONY: stash-macapp
stash-macapp: STASH_OUTPUT := -o stash
stash-macapp: flags-release flags-pie stash
	rm -rf Stash.app
	cp -R scripts/macos-bundle Stash.app
	mkdir Stash.app/Contents/MacOS
	cp stash Stash.app/Contents/MacOS/stash

# build-cc- targets should be run within the compiler docker container

.PHONY: build-cc-windows
build-cc-windows: export GOOS := windows
build-cc-windows: export GOARCH := amd64
build-cc-windows: export CC := x86_64-w64-mingw32-gcc
build-cc-windows: STASH_OUTPUT := -o dist/stash-win.exe
build-cc-windows: PHASHER_OUTPUT :=-o dist/phasher-win.exe
build-cc-windows: flags-release
build-cc-windows: flags-static-windows
build-cc-windows: build

.PHONY: build-cc-macos-intel
build-cc-macos-intel: export GOOS := darwin
build-cc-macos-intel: export GOARCH := amd64
build-cc-macos-intel: export CC := o64-clang
build-cc-macos-intel: STASH_OUTPUT := -o dist/stash-macos-intel
build-cc-macos-intel: PHASHER_OUTPUT := -o dist/phasher-macos-intel
build-cc-macos-intel: flags-release
# can't use static build for macOS
build-cc-macos-intel: flags-pie
build-cc-macos-intel: build

.PHONY: build-cc-macos-arm
build-cc-macos-arm: export GOOS := darwin
build-cc-macos-arm: export GOARCH := arm64
build-cc-macos-arm: export CC := oa64e-clang
build-cc-macos-arm: STASH_OUTPUT := -o dist/stash-macos-arm
build-cc-macos-arm: PHASHER_OUTPUT := -o dist/phasher-macos-arm
build-cc-macos-arm: flags-release
# can't use static build for macOS
build-cc-macos-arm: flags-pie
build-cc-macos-arm: build

.PHONY: build-cc-macos
build-cc-macos:
	make build-cc-macos-arm
	make build-cc-macos-intel

	# Combine into universal binaries
	lipo -create -output dist/stash-macos dist/stash-macos-intel dist/stash-macos-arm
	rm dist/stash-macos-intel dist/stash-macos-arm
	lipo -create -output dist/phasher-macos dist/phasher-macos-intel dist/phasher-macos-arm
	rm dist/phasher-macos-intel dist/phasher-macos-arm

	# Place into bundle and zip up
	rm -rf dist/Stash.app
	cp -R scripts/macos-bundle dist/Stash.app
	mkdir dist/Stash.app/Contents/MacOS
	cp dist/stash-macos dist/Stash.app/Contents/MacOS/stash
	cd dist && rm -f Stash.app.zip && zip -r Stash.app.zip Stash.app
	rm -rf dist/Stash.app

.PHONY: build-cc-freebsd
build-cc-freebsd: export GOOS := freebsd
build-cc-freebsd: export GOARCH := amd64
build-cc-freebsd: export CC := clang -target x86_64-unknown-freebsd12.0 --sysroot=/opt/cross-freebsd
build-cc-freebsd: STASH_OUTPUT := -o dist/stash-freebsd
build-cc-freebsd: PHASHER_OUTPUT := -o dist/phasher-freebsd
build-cc-freebsd: flags-release
build-cc-freebsd: flags-static-pie
build-cc-freebsd: build

.PHONY: build-cc-linux
build-cc-linux: export GOOS := linux
build-cc-linux: export GOARCH := amd64
build-cc-linux: STASH_OUTPUT := -o dist/stash-linux
build-cc-linux: PHASHER_OUTPUT := -o dist/phasher-linux
build-cc-linux: flags-release
build-cc-linux: flags-static-pie
build-cc-linux: build

.PHONY: build-cc-linux-arm64v8
build-cc-linux-arm64v8: export GOOS := linux
build-cc-linux-arm64v8: export GOARCH := arm64
build-cc-linux-arm64v8: export CC := aarch64-linux-gnu-gcc
build-cc-linux-arm64v8: STASH_OUTPUT := -o dist/stash-linux-arm64v8
build-cc-linux-arm64v8: PHASHER_OUTPUT := -o dist/phasher-linux-arm64v8
build-cc-linux-arm64v8: flags-release
build-cc-linux-arm64v8: flags-static-pie
build-cc-linux-arm64v8: build

.PHONY: build-cc-linux-arm32v7
build-cc-linux-arm32v7: export GOOS := linux
build-cc-linux-arm32v7: export GOARCH := arm
build-cc-linux-arm32v7: export GOARM := 7
build-cc-linux-arm32v7: export CC := arm-linux-gnueabi-gcc -march=armv7-a
build-cc-linux-arm32v7: STASH_OUTPUT := -o dist/stash-linux-arm32v7
build-cc-linux-arm32v7: PHASHER_OUTPUT := -o dist/phasher-linux-arm32v7
build-cc-linux-arm32v7: flags-release
build-cc-linux-arm32v7: flags-static
build-cc-linux-arm32v7: build

.PHONY: build-cc-linux-arm32v6
build-cc-linux-arm32v6: export GOOS := linux
build-cc-linux-arm32v6: export GOARCH := arm
build-cc-linux-arm32v6: export GOARM := 6
build-cc-linux-arm32v6: export CC := arm-linux-gnueabi-gcc
build-cc-linux-arm32v6: STASH_OUTPUT := -o dist/stash-linux-arm32v6
build-cc-linux-arm32v6: PHASHER_OUTPUT := -o dist/phasher-linux-arm32v6
build-cc-linux-arm32v6: flags-release
build-cc-linux-arm32v6: flags-static
build-cc-linux-arm32v6: build

.PHONY: build-cc-all
build-cc-all:
	make build-cc-windows
	make build-cc-macos
	make build-cc-linux
	make build-cc-linux-arm64v8
	make build-cc-linux-arm32v7
	make build-cc-linux-arm32v6
	make build-cc-freebsd

.PHONY: touch-ui
touch-ui:
ifdef IS_WIN_SHELL
	@if not exist "ui\\v2.5\\build" mkdir ui\\v2.5\\build
	@type nul >> ui/v2.5/build/index.html
else
	@mkdir -p ui/v2.5/build
	@touch ui/v2.5/build/index.html
endif

# Regenerates GraphQL files
.PHONY: generate
generate: generate-backend generate-ui

.PHONY: generate-ui
generate-ui:
	cd ui/v2.5 && yarn run gqlgen

.PHONY: generate-backend
generate-backend: touch-ui
	go generate ./cmd/stash

.PHONY: generate-dataloaders
generate-dataloaders:
	go generate ./internal/api/loaders

# Regenerates stash-box client files
.PHONY: generate-stash-box-client
generate-stash-box-client:
	go run github.com/Yamashou/gqlgenc

# Runs gofmt -w on the project's source code, modifying any files that do not match its style.
.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint:
	golangci-lint run

# runs unit tests - excluding integration tests
.PHONY: test
test:
	go test ./...

# runs all tests - including integration tests
.PHONY: it
it:
	$(eval GO_BUILD_TAGS += integration)
	go test -tags "$(GO_BUILD_TAGS)" ./...

# generates test mocks
.PHONY: generate-test-mocks
generate-test-mocks:
	go run github.com/vektra/mockery/v2

# runs server
# sets the config file to use the local dev config
.PHONY: server-start
server-start: export STASH_CONFIG_FILE := config.yml
server-start: build-flags
ifdef IS_WIN_SHELL
	@if not exist ".local" mkdir .local
else
	@mkdir -p .local
endif
	cd .local && go run $(BUILD_FLAGS) ../cmd/stash

# removes local dev config files
.PHONY: server-clean
server-clean:
	$(RMDIR) .local

# installs UI dependencies. Run when first cloning repository, or if UI
# dependencies have changed
.PHONY: pre-ui
pre-ui:
	cd ui/v2.5 && yarn install --frozen-lockfile

.PHONY: ui-env
ui-env: build-info
	$(eval export VITE_APP_DATE := $(BUILD_DATE))
	$(eval export VITE_APP_GITHASH := $(GITHASH))
	$(eval export VITE_APP_STASH_VERSION := $(STASH_VERSION))
ifdef STASH_NOLEGACY
	$(eval export VITE_APP_NOLEGACY := true)
endif
ifdef STASH_SOURCEMAPS
	$(eval export VITE_APP_SOURCEMAPS := true)
endif

.PHONY: ui
ui: ui-env
	cd ui/v2.5 && yarn build

.PHONY: zip-ui
zip-ui:
	rm -f dist/stash-ui.zip
	cd ui/v2.5/build && zip -r ../../../dist/stash-ui.zip .

.PHONY: ui-start
ui-start: ui-env
	cd ui/v2.5 && yarn start --host

.PHONY: fmt-ui
fmt-ui:
	cd ui/v2.5 && yarn format

# runs all of the frontend PR-acceptance steps
.PHONY: validate-ui
validate-ui:
	cd ui/v2.5 && yarn run validate

# these targets run the same steps as fmt-ui and validate-ui, but only on files that have changed
fmt-ui-quick:
	cd ui/v2.5 && yarn run prettier --write $$(git diff --name-only --relative --diff-filter d . ../../graphql)

# does not run tsc checks, as they are slow
validate-ui-quick:
	cd ui/v2.5 && \
	yarn run eslint $$(git diff --name-only --relative --diff-filter d src | grep -e "\.tsx\?\$$") && \
	yarn run stylelint $$(git diff --name-only --relative --diff-filter d src | grep "\.scss") && \
	yarn run prettier --check $$(git diff --name-only --relative --diff-filter d . ../../graphql)

# runs all of the backend PR-acceptance steps
.PHONY: validate-backend
validate-backend: lint it

# runs all of the tests and checks required for a PR to be accepted
.PHONY: validate
validate: validate-ui validate-backend

# locally builds and tags a 'stash/build' docker image
.PHONY: docker-build
docker-build: build-info
	docker build --build-arg GITHASH=$(GITHASH) --build-arg STASH_VERSION=$(STASH_VERSION) -t stash/build -f docker/build/x86_64/Dockerfile .

# locally builds and tags a 'stash/cuda-build' docker image
.PHONY: docker-cuda-build
docker-cuda-build: build-info
	docker build --build-arg GITHASH=$(GITHASH) --build-arg STASH_VERSION=$(STASH_VERSION) -t stash/cuda-build -f docker/build/x86_64/Dockerfile-CUDA .

# start the build container - for cross compilation
# this is adapted from the github actions build.yml file
.PHONY: start-compiler-container
start-compiler-container:
	docker run -d --name build --mount type=bind,source="$(PWD)",target=/stash,consistency=delegated $(EXTRA_CONTAINER_ARGS) -w /stash $(COMPILER_IMAGE) tail -f /dev/null

# run the cross-compilation using
# docker exec -t build /bin/bash -c "make build-cc-<platform>"

.PHONY: remove-compiler-container
remove-compiler-container:
	docker rm -f -v build