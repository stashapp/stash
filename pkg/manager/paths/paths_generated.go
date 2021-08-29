package paths

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/stashapp/stash/pkg/utils"
)

const thumbDirDepth int = 2
const thumbDirLength int = 2 // thumbDirDepth * thumbDirLength must be smaller than the length of checksum

type generatedPaths struct {
	Screenshots string
	Thumbnails  string
	Heatmaps    string
	Vtt         string
	Markers     string
	Transcodes  string
	Downloads   string
	Tmp         string
}

func newGeneratedPaths(path string) *generatedPaths {
	gp := generatedPaths{}
	gp.Screenshots = filepath.Join(path, "screenshots")
	gp.Thumbnails = filepath.Join(path, "thumbnails")
	gp.Heatmaps = filepath.Join(path, "heatmaps")
	gp.Vtt = filepath.Join(path, "vtt")
	gp.Markers = filepath.Join(path, "markers")
	gp.Transcodes = filepath.Join(path, "transcodes")
	gp.Downloads = filepath.Join(path, "download_stage")
	gp.Tmp = filepath.Join(path, "tmp")
	return &gp
}

func (gp *generatedPaths) GetTmpPath(fileName string) string {
	return filepath.Join(gp.Tmp, fileName)
}

func (gp *generatedPaths) EnsureTmpDir() {
	utils.EnsureDir(gp.Tmp)
}

func (gp *generatedPaths) EmptyTmpDir() {
	utils.EmptyDir(gp.Tmp)
}

func (gp *generatedPaths) RemoveTmpDir() {
	utils.RemoveDir(gp.Tmp)
}

func (gp *generatedPaths) TempDir(pattern string) (string, error) {
	gp.EnsureTmpDir()
	ret, err := ioutil.TempDir(gp.Tmp, pattern)
	if err != nil {
		return "", err
	}

	utils.EmptyDir(ret)

	return ret, nil
}

func (gp *generatedPaths) GetThumbnailPath(checksum string, width int) string {
	fname := fmt.Sprintf("%s_%d.jpg", checksum, width)
	return filepath.Join(gp.Thumbnails, utils.GetIntraDir(checksum, thumbDirDepth, thumbDirLength), fname)
}
