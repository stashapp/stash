#!/bin/sh

# assumes cross-compile.sh has already been run successfully
uploadFile() 
{
    FILE=$1
    BASENAME="$(basename "${FILE}")"
    # abort if it takes more than two minutes to upload
    uploadedTo=`curl -m 120 --upload-file $FILE "https://oshi.at/$BASENAME/20160"`
    CDN=`echo "$uploadedTo"|grep CDN`
    echo "$BASENAME uploaded to url: $CDN"
}

uploadFile "dist/stash-osx"
uploadFile "dist/stash-win.exe"
uploadFile "dist/stash-linux"
