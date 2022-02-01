#!/bin/bash
# Update the Stash icon throughout the project from a master stash-logo.png

# Imagemagick, and go packages icns and rsrc are required.
# Copy a high-resolution stash-logo.png to this stash/scripts folder
# and run this script from said folder, commit the result.

if [ ! -f "stash-logo.png" ]; then
    echo "stash-logo.png not found."
    exit
fi

if [ -z "$GOPATH" ]; then
    echo "GOPATH environment variable not set"
    exit
fi

if [ ! -e "$GOPATH/bin/rsrc" ]; then
    echo "Missing Dependency:"
    echo "Please run the following /outside/ of the stash folder:"
    echo "go install github.com/akavel/rsrc@latest" 
    exit
fi

if [ ! -e "$GOPATH/bin/icnsify" ]; then
    echo "Missing Dependency:"
    echo "Please run the following /outside/ of the stash folder:"
    echo "go install github.com/jackmordaunt/icns/v2/cmd/icnsify@latest" 
    exit
fi

# Favicon, used for web favicon, windows systray icon, windows executable icon
convert stash-logo.png -define icon:auto-resize=256,64,48,32,16 favicon.ico
cp favicon.ico ../ui/v2.5/public/

# Build .syso for Windows icon, consumed by linker while building stash-win.exe
"$GOPATH"/bin/rsrc -ico favicon.ico -o icon_windows.syso
mv icon_windows.syso ../pkg/desktop/

# *nixes systray icon
convert stash-logo.png -resize x256 favicon.png
cp favicon.png ../ui/v2.5/public/

# MacOS, used for bundle icon
# https://developer.apple.com/library/archive/documentation/CoreFoundation/Conceptual/CFBundles/BundleTypes/BundleTypes.html
"$GOPATH"/bin/icnsify -i stash-logo.png -o icon.icns
mv icon.icns macos-bundle/Contents/Resources/icon.icns

# cleanup
rm favicon.png favicon.ico