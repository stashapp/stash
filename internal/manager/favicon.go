package manager

import (
	"embed"
	"runtime"
)

const faviconDir = "v2.5/build/"

type FaviconProvider struct {
	UIBox embed.FS
}

func (p *FaviconProvider) GetFavicon() []byte {
	if runtime.GOOS == "windows" {
		faviconPath := faviconDir + "favicon.ico"
		ret, _ := p.UIBox.ReadFile(faviconPath)
		return ret
	}

	return p.GetFaviconPng()
}

func (p *FaviconProvider) GetFaviconPng() []byte {
	faviconPath := faviconDir + "favicon.png"
	ret, _ := p.UIBox.ReadFile(faviconPath)
	return ret
}
