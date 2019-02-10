# Stash

**Stash is a rails app which organizes and serves your porn.**

See a demo [here](https://vimeo.com/275537038) (password is stashapp).

TODO: This is not match the features of the Rails project quite yet.  Consider using that until this project is complete.

## Setup

TODO: This is not final.  There is more work to be done to ease this process.

### OSX / Linux

1. `mkdir ~/.stash` && `cd ~/.stash`
2. Download FFMPEG ([macOS](https://ffmpeg.zeranoe.com/builds/macos64/static/ffmpeg-4.0-macos64-static.zip), [Linux](https://www.johnvansickle.com/ffmpeg/old-releases/ffmpeg-4.0.3-64bit-static.tar.xz)) and extract so that just `ffmpeg` and `ffprobe` are in `~/.stash`
3. Create a `config.json` file (see below).
4. Run stash with `./stash` and visit `http://localhost:9998` or `https://localhost:9999`

### Windows

1. Create a new folder at `C:\Users\YourUsername\.stash`
2. Download [FFMPEG](https://ffmpeg.zeranoe.com/builds/win64/static/ffmpeg-4.0-win64-static.zip) and extract so that just `ffmpeg.exe` and `ffprobe.exe` are in `C:\Users\YourUsername\.stash`
3. Create a `config.json` file (see below)
4. Run stash with `./stash` and visit `http://localhost:9998` or `https://localhost:9999`

### Config.json

Example:

*OSX / Linux*
```
{
  "stash": "/Volumes/Drobo/videos",
  "metadata": "/Volumes/Drobo/stash/metadata",
  "cache": "/Volumes/Drobo/stash/cache",
  "downloads": "/Volumes/Drobo/stash/downloads"
}
```

*Windows*
```
{
  "stash": "C:\\Videos",
  "metadata": "C:\\stash\\metadata",
  "cache": "C:\\stash\\cache",
  "downloads": "C:\\stash\\downloads"
}
```

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