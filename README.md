# Stash

**Stash is a rails app which organizes and serves your porn.**

See a demo [here](https://vimeo.com/275537038) (password is stashapp).

TODO

## Setup

TODO

# Development

## Environment

### macOS

TODO

### Windows

1. Download and install [Go for Windows](https://golang.org/dl/)
2. Download and install [MingW](https://sourceforge.net/projects/mingw-w64/)
3. Search for "advanced system settings" and open the system properties dialog.
	1. Click the `Environment Variables` button
	2. Add `GO111MODULE=on`
	3. Under system variables find the `Path`.  Edit and add `C:\Program Files\mingw-w64\*\mingw64\bin` (replace * with the correct path).

## Commands

* `make build` - Builds the binary
* `make gqlgen` - Regenerate Go GraphQL files

## Building a release

1. cd into the UI directory and run `ng build --prod`
2. cd back to the root directory and run `make build` to build the executable

#### Notes for the dev

https://blog.filippo.io/easy-windows-and-linux-cross-compilers-for-macos/