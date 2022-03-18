package paths

import (
	"path/filepath"
	"strconv"
)

type sceneMarkerPaths struct {
	generated generatedPaths
}

func newSceneMarkerPaths(p Paths) *sceneMarkerPaths {
	sp := sceneMarkerPaths{}
	sp.generated = *p.Generated
	return &sp
}

func (sp *sceneMarkerPaths) GetStreamPath(checksum string, seconds int) string {
	return filepath.Join(sp.generated.Markers, checksum, strconv.Itoa(seconds)+".mp4")
}

func (sp *sceneMarkerPaths) GetStreamPreviewImagePath(checksum string, seconds int) string {
	return filepath.Join(sp.generated.Markers, checksum, strconv.Itoa(seconds)+".webp")
}

func (sp *sceneMarkerPaths) GetStreamScreenshotPath(checksum string, seconds int) string {
	return filepath.Join(sp.generated.Markers, checksum, strconv.Itoa(seconds)+".jpg")
}
