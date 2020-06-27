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
	$(SET) CGO_ENABLED=1 $(SEPARATOR) go build $(OUTPUT) -mod=vendor -v -ldflags "$(LDFLAGS) $(EXTRA_LDFLAGS)"

# strips debug symbols from the release build
# consider -trimpath in go build if we move to go 1.13+
build-release: EXTRA_LDFLAGS := -s -w
build-release: build

install:
	packr2 install

clean:
	packr2 clean

# Regenerates GraphQL files
.PHONY: generate
generate:
	go generate -mod=vendor
	cd ui/v2.5 && yarn run gqlgen

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
