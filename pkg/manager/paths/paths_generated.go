package paths

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/utils"
)

const thumbDirDepth int = 2
const thumbDirLength int = 2 // thumbDirDepth * thumbDirLength must be smaller than the length of checksum

type generatedPaths struct {
	Screenshots string
	Thumbnails  string
	Vtt         string
	Markers     string
	Transcodes  string
	Downloads   string
	Tmp         string
}

func newGeneratedPaths() *generatedPaths {
	gp := generatedPaths{}
	gp.Screenshots = filepath.Join(config.GetGeneratedPath(), "screenshots")
	gp.Thumbnails = filepath.Join(config.GetGeneratedPath(), "thumbnails")
	gp.Vtt = filepath.Join(config.GetGeneratedPath(), "vtt")
	gp.Markers = filepath.Join(config.GetGeneratedPath(), "markers")
	gp.Transcodes = filepath.Join(config.GetGeneratedPath(), "transcodes")
	gp.Downloads = filepath.Join(config.GetGeneratedPath(), "downloads")
	gp.Tmp = filepath.Join(config.GetGeneratedPath(), "tmp")
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
