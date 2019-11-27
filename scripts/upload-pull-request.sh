#!/bin/sh

# assumes cross-compile.sh has already been run successfully
uploadFile() 
{
    FILE=$1
    BASENAME="$(basename "${FILE}")"
    uploadedTo=`curl --upload-file $FILE "https://transfer.sh/$BASENAME"`
    echo "$BASENAME uploaded to url: $uploadedTo"
}

uploadFile "dist/stash-osx"
uploadFile "dist/stash-win.exe"
uploadFile "dist/stash-linux"