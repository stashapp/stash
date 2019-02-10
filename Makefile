gqlgen:
	go run scripts/gqlgen.go

build:
	CGO_ENABLED=1 packr2 build -mod=vendor -v