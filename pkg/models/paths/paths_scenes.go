package paths

import (
	"path/filepath"

	"github.com/stashapp/stash/pkg/fsutil"
)

type scenePaths struct {
	generated generatedPaths
}

func newScenePaths(p Paths) *scenePaths {
	sp := scenePaths{}
	sp.generated = *p.Generated
	return &sp
}

func (sp *scenePaths) GetScreenshotPath(checksum string) string {
	return filepath.Join(sp.generated.Screenshots, checksum+".jpg")
}

func (sp *scenePaths) GetThumbnailScreenshotPath(checksum string) string {
	return filepath.Join(sp.generated.Screenshots, checksum+".thumb.jpg")
}

func (sp *scenePaths) GetTranscodePath(checksum string) string {
	return filepath.Join(sp.generated.Transcodes, checksum+".mp4")
}

func (sp *scenePaths) GetStreamPath(scenePath string, checksum string) string {
	transcodePath := sp.GetTranscodePath(checksum)
	transcodeExists, _ := fsutil.FileExists(transcodePath)
	if transcodeExists {
		return transcodePath
	}
	return scenePath
}

func (sp *scenePaths) GetStreamPreviewPath(checksum string) string {
	return filepath.Join(sp.generated.Screenshots, checksum+".mp4")
}

func (sp *scenePaths) GetStreamPreviewImagePath(checksum string) string {
	return filepath.Join(sp.generated.Screenshots, checksum+".webp")
}

func (sp *scenePaths) GetSpriteImageFilePath(checksum string) string {
	return filepath.Join(sp.generated.Vtt, checksum+"_sprite.jpg")
}

func (sp *scenePaths) GetSpriteVttFilePath(checksum string) string {
	return filepath.Join(sp.generated.Vtt, checksum+"_thumbs.vtt")
}

func (sp *scenePaths) GetInteractiveHeatmapPath(checksum string) string {
	return filepath.Join(sp.generated.InteractiveHeatmap, checksum+".png")
}
