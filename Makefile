gqlgen:
	go run scripts/gqlgen.go

build:
	packr2 build

build-win:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 packr2 build -o stash.exe -v