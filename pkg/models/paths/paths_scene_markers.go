package paths

import (
	"path/filepath"
	"strconv"
)

type sceneMarkerPaths struct {
	generatedPaths
}

func newSceneMarkerPaths(p Paths) *sceneMarkerPaths {
	sp := sceneMarkerPaths{
		generatedPaths: *p.Generated,
	}
	return &sp
}

func (sp *sceneMarkerPaths) GetVideoPreviewPath(checksum string, seconds int) string {
	return filepath.Join(sp.Markers, checksum, strconv.Itoa(seconds)+".mp4")
}

func (sp *sceneMarkerPaths) GetWebpPreviewPath(checksum string, seconds int) string {
	return filepath.Join(sp.Markers, checksum, strconv.Itoa(seconds)+".webp")
}

func (sp *sceneMarkerPaths) GetScreenshotPath(checksum string, seconds int) string {
	return filepath.Join(sp.Markers, checksum, strconv.Itoa(seconds)+".jpg")
}
