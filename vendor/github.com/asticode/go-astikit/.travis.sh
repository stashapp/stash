#!/bin/sh

if [ "$(go list -m all)" != "github.com/asticode/go-astikit" ]; then
    echo "This repo doesn't allow any external dependencies"
    exit 1
else
    echo "cheers!"
fi