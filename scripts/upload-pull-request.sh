#!/bin/sh

# assumes cross-compile.sh has already been run successfully
function uploadFile(file) {
    BASENAME="$(basename "${FILE}")"
    curl --upload-file $FILE "https://transfer.sh/$BASENAME"
}

uploadFile("dist/stash-osx")
uploadFile("dist/stash-win.exe")
uploadFile("dist/stash-linux")
