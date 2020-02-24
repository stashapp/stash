#!/bin/sh

# assumes cross-compile.sh has already been run successfully
uploadFile() 
{
    FILE=$1
    BASENAME="$(basename "${FILE}")"
    # abort if it takes more than two minutes to upload
    uploadedTo=`curl -m 120 --upload-file $FILE "https://transfer.sh/$BASENAME"`
    echo "$BASENAME uploaded to url: $uploadedTo"
}

uploadFile "dist/stash-osx"
uploadFile "dist/stash-win.exe"
uploadFile "dist/stash-linux"
