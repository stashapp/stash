package paths

import (
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/utils"
	"path/filepath"
)

type generatedPaths struct {
	Screenshots string
	Vtt         string
	Markers     string
	Transcodes  string
	Tmp         string
}

func newGeneratedPaths() *generatedPaths {
	gp := generatedPaths{}
	gp.Screenshots = filepath.Join(config.GetGeneratedPath(), "screenshots")
	gp.Vtt = filepath.Join(config.GetGeneratedPath(), "vtt")
	gp.Markers = filepath.Join(config.GetGeneratedPath(), "markers")
	gp.Transcodes = filepath.Join(config.GetGeneratedPath(), "transcodes")
	gp.Tmp = filepath.Join(config.GetGeneratedPath(), "tmp")
	return &gp
}

func (gp *generatedPaths) GetTmpPath(fileName string) string {
	return filepath.Join(gp.Tmp, fileName)
}

func (gp *generatedPaths) EnsureTmpDir() {
	_ = utils.EnsureDir(gp.Tmp)
}

func (gp *generatedPaths) EmptyTmpDir() {
	_ = utils.EmptyDir(gp.Tmp)
}

func (gp *generatedPaths) RemoveTmpDir() {
	_ = utils.RemoveDir(gp.Tmp)
}
