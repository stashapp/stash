//go:generate go run -tags=dev ../scripts/generateLoginLocales.go
package ui

import (
	"embed"
	"io/fs"
	"runtime"
)

//go:embed v2.5/build
var uiBox embed.FS
var UIBox fs.FS

//go:embed login
var loginUIBox embed.FS
var LoginUIBox fs.FS

func init() {
	var err error
	UIBox, err = fs.Sub(uiBox, "v2.5/build")
	if err != nil {
		panic(err)
	}

	LoginUIBox, err = fs.Sub(loginUIBox, "login")
	if err != nil {
		panic(err)
	}
}

type faviconProvider struct{}

var FaviconProvider = faviconProvider{}

func (p *faviconProvider) GetFavicon() []byte {
	if runtime.GOOS == "windows" {
		ret, _ := fs.ReadFile(UIBox, "favicon.ico")
		return ret
	}

	return p.GetFaviconPng()
}

func (p *faviconProvider) GetFaviconPng() []byte {
	ret, _ := fs.ReadFile(UIBox, "favicon.png")
	return ret
}
