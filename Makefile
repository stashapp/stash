build:
	CGO_ENABLED=1 packr2 build -mod=vendor -v

install:
	packr2 install

gqlgen:
	go run scripts/gqlgen.go