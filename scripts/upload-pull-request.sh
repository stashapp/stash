#!/bin/sh

# assumes cross-compile.sh has already been run successfully
uploadFile() 
{
    FILE=$1
    BASENAME="$(basename "${FILE}")"
    
    # get available server from gofile api
    serverApi=$(curl -m 15 https://apiv2.gofile.io/getServer)
    resp=$(echo "$serverApi" | cut -d "\"" -f 4)
    
    # if no server is available abort
    if [ $resp != "ok" ] ; then
	echo "Upload of $BASENAME failed! Server not available."
	echo
	return
    fi
    server=$(echo "$serverApi" | cut -d "," -f 2  | cut -d "\"" -f 6)

    # abort if it takes more than two minutes to upload
    uploadedTo=$(curl -m 120 -F "email=stash@stashapp.cc" -F "file=@$FILE" "https://$server.gofile.io/uploadFile")
    resp=$(echo "$uploadedTo" | cut -d "\"" -f 4)
    if [ $resp = "ok" ] ; then
	URL=$(echo "$uploadedTo"|cut -d "," -f 2 | cut -d "\"" -f 6)
	echo "$BASENAME uploaded to url: \"https://gofile.io/d/$URL\""
    fi
    # print an extra newline
    echo

}

uploadFile "dist/stash-osx"
uploadFile "dist/stash-win.exe"
uploadFile "dist/stash-linux"

echo "SHA1 Checksums"
cat CHECKSUMS_SHA1 | grep -v '\-pi\|\-arm'
