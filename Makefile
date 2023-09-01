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
else
  RM := rm -f
  RMDIR := rm -rf
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
GO_BUILD_TAGS += sqlite_stat4

# set STASH_NOLEGACY environment variable or uncomment to disable legacy browser support
# STASH_NOLEGACY := true

# set STASH_SOURCEMAPS environment variable or uncomment to enable UI sourcemaps
# STASH_SOURCEMAPS := true

export CGO_ENABLED := 1

.PHONY: release
release: pre-ui generate ui build-release

# targets to set various build flags

.PHONY: flags-release
flags-release: 
	$(eval LDFLAGS += -s -w)
	$(eval GO_BUILD_FLAGS += -trimpath)

.PHONY: flags-pie
flags-pie:
	$(eval GO_BUILD_FLAGS += -buildmode=pie)

.PHONY: flags-static
flags-static: 
	$(eval LDFLAGS += -extldflags=-static)
	$(eval GO_BUILD_TAGS += sqlite_omit_load_extension osusergo netgo)

.PHONY: flags-static-pie
flags-static-pie: 
	$(eval LDFLAGS += -extldflags=-static-pie)
	$(eval GO_BUILD_FLAGS += -buildmode=pie)
	$(eval GO_BUILD_TAGS += sqlite_omit_load_extension osusergo netgo)

.PHONY: flags-static-windows
flags-static-windows:
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

.PHONY: stash-release
stash-release: flags-release
stash-release: flags-pie
stash-release: stash

.PHONY: stash-release-static
stash-release-static: flags-release
stash-release-static: flags-static-pie
stash-release-static: stash

.PHONY: stash-release-static-windows
stash-release-static-windows: flags-release
stash-release-static-windows: flags-static-windows
stash-release-static-windows: stash

.PHONY: phasher
phasher: build-flags
	go build $(PHASHER_OUTPUT) $(BUILD_FLAGS) ./cmd/phasher

.PHONY: phasher-release
phasher-release: flags-release
phasher-release: flags-pie
phasher-release: phasher

.PHONY: phasher-release-static
phasher-release-static: flags-release
phasher-release-static: flags-static-pie
phasher-release-static: phasher

.PHONY: phasher-release-static-windows
phasher-release-static-windows: flags-release
phasher-release-static-windows: flags-static-windows
phasher-release-static-windows: phasher

# builds dynamically-linked debug binaries
.PHONY: build
build: stash phasher

# builds dynamically-linked release binaries
.PHONY: build-release
build-release: stash-release phasher-release

# builds statically-linked release binaries
.PHONY: build-release-static
build-release-static: stash-release-static phasher-release-static

# build-release-static, but excluding netgo, which is not needed on windows
.PHONY: build-release-static-windows
build-release-static-windows: stash-release-static-windows phasher-release-static-windows

# cross-compile- targets should be run within the compiler docker container
.PHONY: cross-compile-windows
cross-compile-windows: export GOOS := windows
cross-compile-windows: export GOARCH := amd64
cross-compile-windows: export CC := x86_64-w64-mingw32-gcc
cross-compile-windows: export CXX := x86_64-w64-mingw32-g++
cross-compile-windows: STASH_OUTPUT := -o dist/stash-win.exe
cross-compile-windows: PHASHER_OUTPUT := -o dist/phasher-win.exe
cross-compile-windows: flags-release
cross-compile-windows: flags-static-windows
cross-compile-windows: build

.PHONY: cross-compile-macos-intel
cross-compile-macos-intel: export GOOS := darwin
cross-compile-macos-intel: export GOARCH := amd64
cross-compile-macos-intel: export CC := o64-clang
cross-compile-macos-intel: export CXX := o64-clang++
cross-compile-macos-intel: STASH_OUTPUT := -o dist/stash-macos-intel
cross-compile-macos-intel: PHASHER_OUTPUT := -o dist/phasher-macos-intel
cross-compile-macos-intel: flags-release
# can't use static build for OSX
cross-compile-macos-intel: flags-pie
cross-compile-macos-intel: build

.PHONY: cross-compile-macos-applesilicon
cross-compile-macos-applesilicon: export GOOS := darwin
cross-compile-macos-applesilicon: export GOARCH := arm64
cross-compile-macos-applesilicon: export CC := oa64e-clang
cross-compile-macos-applesilicon: export CXX := oa64e-clang++
cross-compile-macos-applesilicon: STASH_OUTPUT := -o dist/stash-macos-applesilicon
cross-compile-macos-applesilicon: PHASHER_OUTPUT := -o dist/phasher-macos-applesilicon
cross-compile-macos-applesilicon: flags-release
# can't use static build for OSX
cross-compile-macos-applesilicon: flags-pie
cross-compile-macos-applesilicon: build

.PHONY: cross-compile-macos
cross-compile-macos:
	rm -rf dist/Stash.app dist/Stash-macos.zip
	make cross-compile-macos-applesilicon
	make cross-compile-macos-intel
	# Combine into one universal binary
	lipo -create -output dist/stash-macos-universal dist/stash-macos-intel dist/stash-macos-applesilicon
	rm dist/stash-macos-intel dist/stash-macos-applesilicon
	# Place into bundle and zip up
	cp -R scripts/macos-bundle dist/Stash.app
	mkdir dist/Stash.app/Contents/MacOS
	mv dist/stash-macos-universal dist/Stash.app/Contents/MacOS/stash
	cd dist && zip -r Stash-macos.zip Stash.app && cd ..
	rm -rf dist/Stash.app

.PHONY: cross-compile-freebsd
cross-compile-freebsd: export GOOS := freebsd
cross-compile-freebsd: export GOARCH := amd64
cross-compile-freebsd: STASH_OUTPUT := -o dist/stash-freebsd
cross-compile-freebsd: PHASHER_OUTPUT := -o dist/phasher-freebsd
cross-compile-freebsd: flags-release
cross-compile-freebsd: flags-static-pie
cross-compile-freebsd: build

.PHONY: cross-compile-linux
cross-compile-linux: export GOOS := linux
cross-compile-linux: export GOARCH := amd64
cross-compile-linux: STASH_OUTPUT := -o dist/stash-linux
cross-compile-linux: PHASHER_OUTPUT := -o dist/phasher-linux
cross-compile-linux: flags-release
cross-compile-linux: flags-static-pie
cross-compile-linux: build

.PHONY: cross-compile-linux-arm64v8
cross-compile-linux-arm64v8: export GOOS := linux
cross-compile-linux-arm64v8: export GOARCH := arm64
cross-compile-linux-arm64v8: export CC := aarch64-linux-gnu-gcc
cross-compile-linux-arm64v8: STASH_OUTPUT := -o dist/stash-linux-arm64v8
cross-compile-linux-arm64v8: PHASHER_OUTPUT := -o dist/phasher-linux-arm64v8
cross-compile-linux-arm64v8: flags-release
cross-compile-linux-arm64v8: flags-static-pie
cross-compile-linux-arm64v8: build

.PHONY: cross-compile-linux-arm32v7
cross-compile-linux-arm32v7: export GOOS := linux
cross-compile-linux-arm32v7: export GOARCH := arm
cross-compile-linux-arm32v7: export GOARM := 7
cross-compile-linux-arm32v7: export CC := arm-linux-gnueabihf-gcc
cross-compile-linux-arm32v7: STASH_OUTPUT := -o dist/stash-linux-arm32v7
cross-compile-linux-arm32v7: PHASHER_OUTPUT := -o dist/phasher-linux-arm32v7
cross-compile-linux-arm32v7: flags-release
cross-compile-linux-arm32v7: flags-static
cross-compile-linux-arm32v7: build

.PHONY: cross-compile-linux-arm32v6
cross-compile-linux-arm32v6: export GOOS := linux
cross-compile-linux-arm32v6: export GOARCH := arm
cross-compile-linux-arm32v6: export GOARM := 6
cross-compile-linux-arm32v6: export CC := arm-linux-gnueabi-gcc
cross-compile-linux-arm32v6: STASH_OUTPUT := -o dist/stash-linux-arm32v6
cross-compile-linux-arm32v6: PHASHER_OUTPUT := -o dist/phasher-linux-arm32v6
cross-compile-linux-arm32v6: flags-release
cross-compile-linux-arm32v6: flags-static
cross-compile-linux-arm32v6: build

.PHONY: cross-compile-all
cross-compile-all:
	make cross-compile-windows
	make cross-compile-macos-intel
	make cross-compile-macos-applesilicon
	make cross-compile-linux
	make cross-compile-linux-arm64v8
	make cross-compile-linux-arm32v7
	make cross-compile-linux-arm32v6

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
	go test -tags=integration ./...

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

.PHONY: ui-nolegacy
ui-nolegacy: STASH_NOLEGACY := true
ui-nolegacy: ui

.PHONY: ui-sourcemaps
ui-sourcemaps: STASH_SOURCEMAPS := true
ui-sourcemaps: ui

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
