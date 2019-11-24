ifeq ($(OS),Windows_NT)
  SEPARATOR := &&
  SET := set
endif

build:
	$(eval DATE := $(shell go run scripts/getDate.go))
	$(eval GITHASH := $(shell git rev-parse --short HEAD))
	$(SET) CGO_ENABLED=1 $(SEPARATOR) go build -mod=vendor -v -ldflags "-X 'github.com/stashapp/stash/pkg/api.buildstamp=$(DATE)' -X 'github.com/stashapp/stash/pkg/api.githash=$(GITHASH)'"

install:
	packr2 install

clean:
	packr2 clean

# Regenerates GraphQL files
.PHONY: generate
generate:
	go generate -mod=vendor
	cd ui/v2 && yarn run gqlgen

# Runs gofmt -w on the project's source code, modifying any files that do not match its style.
.PHONY: fmt
fmt:
	go fmt ./...

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

.PHONY: ui
ui:
	cd ui/v2 && yarn build
	packr2
