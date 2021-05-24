IS_WIN =
ifeq (${SHELL}, sh.exe)
  IS_WIN = true
endif
ifeq (${SHELL}, cmd)
  IS_WIN = true
endif

ifdef IS_WIN
  SEPARATOR := &&
  SET := set
else 
  SEPARATOR := ;
  SET := export
endif

# set LDFLAGS environment variable to any extra ldflags required
# set OUTPUT to generate a specific binary name

LDFLAGS := $(LDFLAGS)
ifdef OUTPUT
  OUTPUT := -o $(OUTPUT)
endif

export CGO_ENABLED = 1
export GO111MODULE = on

.PHONY: release pre-build install clean 

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

build: pre-build
	$(eval LDFLAGS := $(LDFLAGS) -X 'github.com/stashapp/stash/pkg/api.version=$(STASH_VERSION)' -X 'github.com/stashapp/stash/pkg/api.buildstamp=$(BUILD_DATE)' -X 'github.com/stashapp/stash/pkg/api.githash=$(GITHASH)')
	go build $(OUTPUT) -mod=vendor -v -tags "sqlite_omit_load_extension osusergo netgo" -ldflags "$(LDFLAGS) $(EXTRA_LDFLAGS)"

# strips debug symbols from the release build
# consider -trimpath in go build if we move to go 1.13+
build-release: EXTRA_LDFLAGS := -s -w
build-release: build

build-release-static: EXTRA_LDFLAGS := -extldflags=-static -s -w
build-release-static: build

# cross-compile- targets should be run within the compiler docker container
cross-compile-windows: export GOOS := windows
cross-compile-windows: export GOARCH := amd64
cross-compile-windows: export CC := x86_64-w64-mingw32-gcc
cross-compile-windows: export CXX := x86_64-w64-mingw32-g++
cross-compile-windows: OUTPUT := -o dist/stash-win.exe
cross-compile-windows: build-release-static

cross-compile-osx: export GOOS := darwin
cross-compile-osx: export GOARCH := amd64
cross-compile-osx: export CC := o64-clang
cross-compile-osx: export CXX := o64-clang++
cross-compile-osx: OUTPUT := -o dist/stash-osx
# can't use static build for OSX
cross-compile-osx: build-release

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

cross-compile-pi: export GOOS := linux
cross-compile-pi: export GOARCH := arm
cross-compile-pi: export GOARM := 6
cross-compile-pi: export CC := arm-linux-gnueabi-gcc
cross-compile-pi: OUTPUT := -o dist/stash-pi
cross-compile-pi: build-release-static

cross-compile-all: cross-compile-windows cross-compile-osx cross-compile-linux cross-compile-linux-arm64v8 cross-compile-linux-arm32v7 cross-compile-pi

install:
	packr2 install

clean:
	packr2 clean

# Regenerates GraphQL files
.PHONY: generate
generate:
	go generate -mod=vendor
	cd ui/v2.5 && yarn run gqlgen

# Regenerates stash-box client files
.PHONY: generate-stash-box-client
generate-stash-box-client:
	go run -mod=vendor github.com/Yamashou/gqlgenc

# Runs gofmt -w on the project's source code, modifying any files that do not match its style.
.PHONY: fmt
fmt:
	go fmt ./...

# Ensures that changed files have had gofmt run on them
.PHONY: fmt-check
fmt-check:
	sh ./scripts/check-gofmt.sh

# Runs go vet on the project's source code.
.PHONY: vet
vet:
	go vet -mod=vendor ./...

.PHONY: lint
lint:
	revive -config revive.toml -exclude ./vendor/...  ./...

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

.PHONY: ui-only
ui-only: pre-build
	$(SET) REACT_APP_DATE="$(BUILD_DATE)" $(SEPARATOR) \
	$(SET) REACT_APP_GITHASH=$(GITHASH) $(SEPARATOR) \
	$(SET) REACT_APP_STASH_VERSION=$(STASH_VERSION) $(SEPARATOR) \
	cd ui/v2.5 && yarn build

.PHONY: ui
ui: ui-only
	packr2

.PHONY: ui-start
ui-start: pre-build
	$(SET) REACT_APP_DATE="$(BUILD_DATE)" $(SEPARATOR) \
	$(SET) REACT_APP_GITHASH=$(GITHASH) $(SEPARATOR) \
	$(SET) REACT_APP_STASH_VERSION=$(STASH_VERSION) $(SEPARATOR) \
	cd ui/v2.5 && yarn start

.PHONY: fmt-ui
fmt-ui:
	cd ui/v2.5 && yarn format

# runs tests and checks on the UI and builds it
.PHONY: ui-validate
ui-validate:
	cd ui/v2.5 && yarn run validate

# just repacks the packr files - use when updating migrations and packed files without 
# rebuilding the UI
.PHONY: packr
packr:
	packr2

# runs all of the tests and checks required for a PR to be accepted
.PHONY: validate
validate: ui-validate fmt-check vet lint it
