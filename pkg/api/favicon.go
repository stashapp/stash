package api

import (
	"embed"
	"runtime"
)

const faviconDir = "ui/v2.5/build/"

type FaviconProvider struct {
	uiBox embed.FS
}

func (p *FaviconProvider) GetFavicon() []byte {
	if runtime.GOOS == "windows" {
		faviconPath := faviconDir + "favicon.ico"
		ret, _ := p.uiBox.ReadFile(faviconPath)
		return ret
	}

	return p.GetFaviconPng()
}

func (p *FaviconProvider) GetFaviconPng() []byte {
	faviconPath := faviconDir + "favicon.png"
	ret, _ := p.uiBox.ReadFile(faviconPath)
	return ret
}
