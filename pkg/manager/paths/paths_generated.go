package paths

import (
	"io/ioutil"
	"path/filepath"

	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/utils"
)

type generatedPaths struct {
	Screenshots string
	Vtt         string
	Markers     string
	Transcodes  string
	Downloads   string
	Tmp         string
}

func newGeneratedPaths() *generatedPaths {
	gp := generatedPaths{}
	gp.Screenshots = filepath.Join(config.GetGeneratedPath(), "screenshots")
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
