SHELL=bash

.PHONY: all
all: clean fmt mocks test install docker integration

.PHONY: clean
clean:
	rm -rf mocks

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: test
test: mocks
	go test ./...

mocks: $(shell find . -type f -name '*.go' -not -name '*_test.go')
	go run . --dir pkg/fixtures --output mocks/pkg/fixtures
	go run . --all=false --print --dir pkg/fixtures --name RequesterVariadic --structname RequesterVariadicOneArgument --unroll-variadic=False > mocks/pkg/fixtures/RequesterVariadicOneArgument.go
	@touch mocks

.PHONY: install
install:
	go install .

.PHONY: docker
docker:
	docker build -t vektra/mockery .

.PHONY: integration
integration: docker install
	./hack/run-e2e.sh
