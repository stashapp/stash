IS_WIN_SHELL =
ifeq (${SHELL}, sh.exe)
  IS_WIN_SHELL = true
endif
ifeq (${SHELL}, cmd)
  IS_WIN_SHELL = true
endif

ifdef IS_WIN_SHELL
  SEPARATOR := &&
  SET := set
else
  SEPARATOR := ;
  SET := export
endif

IS_WIN_OS =
ifeq ($(OS),Windows_NT)
	IS_WIN_OS = true
endif

# set LDFLAGS environment variable to any extra ldflags required
# set OUTPUT to generate a specific binary name

LDFLAGS := $(LDFLAGS)
ifdef OUTPUT
  OUTPUT := -o $(OUTPUT)
endif

export CGO_ENABLED = 1

.PHONY: release pre-build

release: generate ui build-release

pre-build:
ifndef BUILD_DATE
	$(eval BUILD_DATE := $(shell go run -mod=vendor scripts/getDate.go))
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

build: pre-build
build:
	$(eval LDFLAGS := $(LDFLAGS) -X 'github.com/stashapp/stash/internal/api.version=$(STASH_VERSION)' -X 'github.com/stashapp/stash/internal/api.buildstamp=$(BUILD_DATE)' -X 'github.com/stashapp/stash/internal/api.githash=$(GITHASH)')
	$(eval LDFLAGS := $(LDFLAGS) -X 'github.com/stashapp/stash/internal/manager/config.officialBuild=$(OFFICIAL_BUILD)')
	go build $(OUTPUT) -mod=vendor -v -tags "sqlite_omit_load_extension osusergo netgo" $(GO_BUILD_FLAGS) -ldflags "$(LDFLAGS) $(EXTRA_LDFLAGS) $(PLATFORM_SPECIFIC_LDFLAGS)" ./cmd/stash

# strips debug symbols from the release build
build-release: EXTRA_LDFLAGS := -s -w
build-release: GO_BUILD_FLAGS := -trimpath
build-release: build

build-release-static: EXTRA_LDFLAGS := -extldflags=-static -s -w
build-release-static: GO_BUILD_FLAGS := -trimpath
build-release-static: build

# cross-compile- targets should be run within the compiler docker container
cross-compile-windows: export GOOS := windows
cross-compile-windows: export GOARCH := amd64
cross-compile-windows: export CC := x86_64-w64-mingw32-gcc
cross-compile-windows: export CXX := x86_64-w64-mingw32-g++
cross-compile-windows: OUTPUT := -o dist/stash-win.exe
cross-compile-windows: build-release-static

cross-compile-macos-intel: export GOOS := darwin
cross-compile-macos-intel: export GOARCH := amd64
cross-compile-macos-intel: export CC := o64-clang
cross-compile-macos-intel: export CXX := o64-clang++
cross-compile-macos-intel: OUTPUT := -o dist/stash-macos-intel
# can't use static build for OSX
cross-compile-macos-intel: build-release

cross-compile-macos-applesilicon: export GOOS := darwin
cross-compile-macos-applesilicon: export GOARCH := arm64
cross-compile-macos-applesilicon: export CC := oa64e-clang
cross-compile-macos-applesilicon: export CXX := oa64e-clang++
cross-compile-macos-applesilicon: OUTPUT := -o dist/stash-macos-applesilicon
# can't use static build for OSX
cross-compile-macos-applesilicon: build-release

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

cross-compile-freebsd: export GOOS := freebsd
cross-compile-freebsd: export GOARCH := amd64
cross-compile-freebsd: OUTPUT := -o dist/stash-freebsd
cross-compile-freebsd: build-release-static

cross-compile-linux: export GOOS := linux
cross-compile-linux: export GOARCH := amd64
cross-compile-linux: OUTPUT := -o dist/stash-linux
cross-compile-linux: build-release-static

cross-compile-linux-arm64v8: export GOOS := linux
cross-compile-linux-arm64v8: export GOARCH := arm64
cross-compile-linux-arm64v8: export CC := aarch64-linux-gnu-gcc
cross-compile-linux-arm64v8: OUTPUT := -o dist/stash-linux-arm64v8
cross-compile-linux-arm64v8: build-release-static

cross-compile-linux-arm32v7: export GOOS := linux
cross-compile-linux-arm32v7: export GOARCH := arm
cross-compile-linux-arm32v7: export GOARM := 7
cross-compile-linux-arm32v7: export CC := arm-linux-gnueabihf-gcc
cross-compile-linux-arm32v7: OUTPUT := -o dist/stash-linux-arm32v7
cross-compile-linux-arm32v7: build-release-static

cross-compile-linux-arm32v6: export GOOS := linux
cross-compile-linux-arm32v6: export GOARCH := arm
cross-compile-linux-arm32v6: export GOARM := 6
cross-compile-linux-arm32v6: export CC := arm-linux-gnueabi-gcc
cross-compile-linux-arm32v6: OUTPUT := -o dist/stash-linux-arm32v6
cross-compile-linux-arm32v6: build-release-static

cross-compile-all:
	make cross-compile-windows
	make cross-compile-macos
	make cross-compile-linux
	make cross-compile-linux-arm64v8
	make cross-compile-linux-arm32v7
	make cross-compile-linux-arm32v6

.PHONY: touch-ui
touch-ui:
ifndef IS_WIN_SHELL
	@mkdir -p ui/v2.5/build
	@touch ui/v2.5/build/index.html
else
	@if not exist "ui\\v2.5\\build" mkdir ui\\v2.5\\build
	@type nul >> ui/v2.5/build/index.html
endif

# Regenerates GraphQL files
generate: generate-backend generate-frontend

.PHONY: generate-frontend
generate-frontend:
	cd ui/v2.5 && yarn run gqlgen

.PHONY: generate-backend
generate-backend: touch-ui 
	go generate -mod=vendor ./cmd/stash

# Regenerates stash-box client files
.PHONY: generate-stash-box-client
generate-stash-box-client:
	go run -mod=vendor github.com/Yamashou/gqlgenc

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
	go test -mod=vendor ./...

# runs all tests - including integration tests
.PHONY: it
it:
	go test -mod=vendor -tags=integration ./...

# generates test mocks
.PHONY: generate-test-mocks
generate-test-mocks:
	go run -mod=vendor github.com/vektra/mockery/v2 --dir ./pkg/models --name '.*ReaderWriter' --outpkg mocks --output ./pkg/models/mocks

# installs UI dependencies. Run when first cloning repository, or if UI
# dependencies have changed
.PHONY: pre-ui
pre-ui:
	cd ui/v2.5 && yarn install --frozen-lockfile

.PHONY: ui
ui: pre-build
	$(SET) VITE_APP_DATE="$(BUILD_DATE)" $(SEPARATOR) \
	$(SET) VITE_APP_GITHASH=$(GITHASH) $(SEPARATOR) \
	$(SET) VITE_APP_STASH_VERSION=$(STASH_VERSION) $(SEPARATOR) \
	cd ui/v2.5 && yarn build

.PHONY: ui-start
ui-start: pre-build
	$(SET) VITE_APP_DATE="$(BUILD_DATE)" $(SEPARATOR) \
	$(SET) VITE_APP_GITHASH=$(GITHASH) $(SEPARATOR) \
	$(SET) VITE_APP_STASH_VERSION=$(STASH_VERSION) $(SEPARATOR) \
	cd ui/v2.5 && yarn start --host

.PHONY: fmt-ui
fmt-ui:
	cd ui/v2.5 && yarn format

# runs tests and checks on the UI and builds it
.PHONY: ui-validate
ui-validate:
	cd ui/v2.5 && yarn run validate

# runs all of the tests and checks required for a PR to be accepted
.PHONY: validate
validate: validate-frontend validate-backend

# runs all of the frontend PR-acceptance steps
.PHONY: validate-frontend
validate-frontend: ui-validate

# runs all of the backend PR-acceptance steps
.PHONY: validate-backend
validate-backend: lint it

# locally builds and tags a 'stash/build' docker image
.PHONY: docker-build
docker-build: pre-build
	docker build --build-arg GITHASH=$(GITHASH) --build-arg STASH_VERSION=$(STASH_VERSION) -t stash/build -f docker/build/x86_64/Dockerfile .
