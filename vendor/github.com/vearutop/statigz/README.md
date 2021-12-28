# statigz

[![Build Status](https://github.com/vearutop/statigz/workflows/test-unit/badge.svg)](https://github.com/vearutop/statigz/actions?query=branch%3Amaster+workflow%3Atest-unit)
[![Coverage Status](https://codecov.io/gh/vearutop/statigz/branch/master/graph/badge.svg)](https://codecov.io/gh/vearutop/statigz)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/github.com/vearutop/statigz)
[![Time Tracker](https://wakatime.com/badge/github/vearutop/statigz.svg)](https://wakatime.com/badge/github/vearutop/statigz)
![Code lines](https://sloc.xyz/github/vearutop/statigz/?category=code)
![Comments](https://sloc.xyz/github/vearutop/statigz/?category=comments)

`statigz` serves pre-compressed embedded files with http in Go 1.16 and later.

## Why?

Since version 1.16 Go provides [standard way](https://tip.golang.org/pkg/embed/) to embed static assets. This API has
advantages over previous solutions:

* assets are processed during build, so there is no need for manual generation step,
* embedded data does not need to be kept in residential memory (as opposed to previous solutions that kept data in
  regular byte slices).

A common case for embedding is to serve static assets of a web application. In order to save bandwidth and improve
latency, those assets are often served compressed. Compression concerns are out of `embed` responsibilities, yet they
are quite important. Previous solutions (for example [`vfsgen`](https://github.com/shurcooL/vfsgen)
with [`httpgzip`](https://github.com/shurcooL/httpgzip)) can optimize performance by storing compressed assets and
serving them directly to capable user agents. This library implements such functionality for embedded file systems.

Read more in a [blog post](https://dev.to/vearutop/serving-compressed-static-assets-with-http-in-go-1-16-55bb).

> **_NOTE:_** Guarding new api (`embed`) with build tags is not a viable option, since it imposes
> [issue](https://github.com/golang/go/issues/40067) in older versions of Go.

## Example

```go
package main

import (
	"embed"
	"log"
	"net/http"

	"github.com/vearutop/statigz"
	"github.com/vearutop/statigz/brotli"
)

// Declare your embedded assets.

//go:embed static/*
var st embed.FS

func main() {
	// Plug static assets handler to your server or router.
	err := http.ListenAndServe(":80", statigz.FileServer(st, brotli.AddEncoding))
	if err != nil {
		log.Fatal(err)
	}
}
```

## Usage

Behavior is based on [nginx gzip static module](http://nginx.org/en/docs/http/ngx_http_gzip_static_module.html) and
[`github.com/lpar/gzipped`](https://github.com/lpar/gzipped).

Static assets have to be manually compressed with additional file extension, e.g. `bundle.js` would
become `bundle.js.gz` (compressed with gzip) or `index.html` would become `index.html.br` (compressed with brotli).

> **_NOTE:_** [`zopfli`](https://github.com/google/zopfli) provides better compression than `gzip` while being
> backwards compatible with it.

Upon request server checks if there is a compressed file matching `Accept-Encoding` and serves it directly.

If user agent does not support available compressed data, server uses an uncompressed file if it is available (
e.g. `bundle.js`). If uncompressed file is not available, then server would decompress a compressed file into response.

Responses have `ETag` headers (64-bit FNV-1 hash of file contents) to enable caching. Responses that are not dynamically
decompressed are served with [`http.ServeContent`](https://golang.org/pkg/net/http/#ServeContent) for ranges support.

### Brotli support

Support for `brotli` is optional. Using `brotli` adds about 260 KB to binary size, that's why it is moved to a separate
package.

> **_NOTE:_** Although [`brotli`](https://github.com/google/brotli) has better compression than `gzip` and already
> has wide support in browsers, it has limitations for non-https servers,
> see [this](https://bugs.chromium.org/p/chromium/issues/detail?id=452335)
> and [this](https://bugzilla.mozilla.org/show_bug.cgi?id=1218924).

### Runtime encoding

Recommended way of embedding assets is to compress assets before the build, so that binary includes `*.gz` or `*.br`
files. This can be inconvenient in some cases, there is `EncodeOnInit` option to compress assets in runtime when
creating file server. Once compressed, assets will be served directly without additional dynamic compression.

Files with extensions ".gz", ".br", ".gif", ".jpg", ".png", ".webp" are excluded from runtime encoding by default.

> **_NOTE:_** Compressing assets in runtime can degrade startup performance and increase memory usage to prepare and store compressed data.

### Mounting a subdirectory

It may be convenient to strip leading directory from an embedded file system, you can do that with `fs.Sub` and a type
assertion.

```go
package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/vearutop/statigz"
	"github.com/vearutop/statigz/brotli"
)

// Declare your embedded assets.

//go:embed static/*
var st embed.FS

func main() {
	// Retrieve sub directory.
	sub, err := fs.Sub(st, "static")
	if err != nil {
		log.Fatal(err)
	}

	// Plug static assets handler to your server or router.
	err = http.ListenAndServe(":80", statigz.FileServer(sub.(fs.ReadDirFS), brotli.AddEncoding))
	if err != nil {
		log.Fatal(err)
	}
}
```