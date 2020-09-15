SHELL=bash

all: clean fmt test fixture install docker integration

clean:
	rm -rf mocks

fmt:
	go fmt ./...

test:
	go test ./...

fixture:
	mockery --print --dir pkg/fixtures --name RequesterVariadic > pkg/fixtures/mocks/requester_variadic.go

install:
	go install ./...

docker:
	docker build -t vektra/mockery .

integration: docker install
	./hack/run-e2e.sh
