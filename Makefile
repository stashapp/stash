ifeq ($(OS),Windows_NT)
  SEPARATOR := &&
  SET := set
endif

release: generate ui build

build:
	$(eval DATE := $(shell go run scripts/getDate.go))
	$(eval GITHASH := $(shell git rev-parse --short HEAD))
	$(eval STASH_VERSION := $(shell git describe --tags --exclude latest_develop))
	$(SET) CGO_ENABLED=1 $(SEPARATOR) go build -mod=vendor -v -ldflags "-X 'github.com/stashapp/stash/pkg/api.version=$(STASH_VERSION)' -X 'github.com/stashapp/stash/pkg/api.buildstamp=$(DATE)' -X 'github.com/stashapp/stash/pkg/api.githash=$(GITHASH)'"

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

.PHONY: ui
ui:
	cd ui/v2.5 && yarn build
	packr2

# just repacks the packr files - use when updating migrations and packed files without 
# rebuilding the UI
.PHONY: packr
packr:
	packr2
