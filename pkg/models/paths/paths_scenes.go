package paths

import (
	"fmt"
	"path/filepath"

	"github.com/stashapp/stash/pkg/fsutil"
)

type scenePaths struct {
	generated generatedPaths
	data      dataPaths
}

func newScenePaths(p Paths) *scenePaths {
	sp := scenePaths{}
	sp.generated = *p.Generated
	sp.data = *p.Data
	return &sp
}

func (sp *scenePaths) GetCoverPath(id int) string {
	fname := fmt.Sprintf("%d_cover.jpg", id)
	return filepath.Join(sp.data.Scenes, fsutil.GetIntraDirID(id, dataDirLength), fname)
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
