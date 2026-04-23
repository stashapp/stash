// TODO(audio): update this file
package paths

import (
	"path/filepath"

	"github.com/stashapp/stash/pkg/fsutil"
)

type audioPaths struct {
	generatedPaths
}

func newAudioPaths(p Paths) *audioPaths {
	sp := audioPaths{
		generatedPaths: *p.Generated,
	}
	return &sp
}

func (sp *audioPaths) GetLegacyScreenshotPath(checksum string) string {
	return filepath.Join(sp.Screenshots, checksum+".jpg")
}

func (sp *audioPaths) GetTranscodePath(checksum string) string {
	return filepath.Join(sp.Transcodes, checksum+".mp4")
}

func (sp *audioPaths) GetStreamPath(audioPath string, checksum string) string {
	transcodePath := sp.GetTranscodePath(checksum)
	transcodeExists, _ := fsutil.FileExists(transcodePath)
	if transcodeExists {
		return transcodePath
	}
	return audioPath
}

func (sp *audioPaths) GetVideoPreviewPath(checksum string) string {
	return filepath.Join(sp.Screenshots, checksum+".mp4")
}

func (sp *audioPaths) GetWebpPreviewPath(checksum string) string {
	return filepath.Join(sp.Screenshots, checksum+".webp")
}

func (sp *audioPaths) GetSpriteImageFilePath(checksum string) string {
	return filepath.Join(sp.Vtt, checksum+"_sprite.jpg")
}

func (sp *audioPaths) GetSpriteVttFilePath(checksum string) string {
	return filepath.Join(sp.Vtt, checksum+"_thumbs.vtt")
}

func (sp *audioPaths) GetInteractiveHeatmapPath(checksum string) string {
	return filepath.Join(sp.InteractiveHeatmap, checksum+".png")
}
