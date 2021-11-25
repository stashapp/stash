#!/bin/bash
# Update the Stash icon throughout the project from a master stash-logo.png

# Imagemagick, and go packages 2goarray and rsrc are required.
# Copy a high-resolution stash-logo.png to this stash/scripts folder
# and run this script from said folder and can commit the result.

if [ ! -f "stash-logo.png" ]; then
    echo "stash-logo.png not found."
    exit
fi

if [ -z "$GOPATH" ]; then
    echo "GOPATH environment variable not set"
    exit
fi

if [ ! -e "$GOPATH/bin/2goarray" ]; then
    echo "Missing Dependency:"
    echo "Please run the following /outside/ of the stash folder:"
    echo "go get github.com/cratonica/2goarray" 
    exit
fi

if [ ! -e "$GOPATH/bin/rsrc" ]; then
    echo "Missing Dependency:"
    echo "Please run the following /outside/ of the stash folder:"
    echo "go get github.com/akavel/rsrc" 
    exit
fi

# Favicon, used for web and for windows executable icon
convert stash-logo.png -define icon:auto-resize=256,64,48,32,16 favicon.ico
cp favicon.ico ../ui/v2.5/public/

# Linux / Macos 
convert stash-logo.png -resize x256 favicon.png

# Add icons for systray / notifications
echo "//go:build linux || darwin" > ../pkg/desktop/favicon_unix.go
echo "// +build linux darwin" >> ../pkg/desktop/favicon_unix.go
echo >> ../pkg/desktop/favicon_unix.go
"$GOPATH"/bin/2goarray favicon desktop < favicon.png >> ../pkg/desktop/favicon_unix.go

echo "//go:build windows" > ../pkg/desktop/favicon_windows.go
echo "// +build windows" >> ../pkg/desktop/favicon_windows.go
echo >> ../pkg/desktop/favicon_windows.go
"$GOPATH"/bin/2goarray favicon desktop < favicon.ico >> ../pkg/desktop/favicon_windows.go

# Build .syso for Windows icon, consumed by linker while building stash-win.exe
"$GOPATH"/bin/rsrc -ico favicon.ico -o icon_windows.syso
mv icon_windows.syso ../pkg/desktop/

# cleanup
rm favicon.png favicon.ico