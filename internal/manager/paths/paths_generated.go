package paths

import (
	"github.com/stashapp/stash/internal/utils"
	"path/filepath"
)

type generatedPaths struct {
	Screenshots string
	Vtt         string
	Markers     string
	Transcodes  string
	Tmp         string
}

func newGeneratedPaths(p Paths) *generatedPaths {
	gp := generatedPaths{}
	gp.Screenshots = filepath.Join(p.Config.Metadata, "screenshots")
	gp.Vtt = filepath.Join(p.Config.Metadata, "vtt")
	gp.Markers = filepath.Join(p.Config.Metadata, "markers")
	gp.Transcodes = filepath.Join(p.Config.Metadata, "transcodes")
	gp.Tmp = filepath.Join(p.Config.Metadata, "tmp")

	_ = utils.EnsureDir(gp.Screenshots)
	_ = utils.EnsureDir(gp.Vtt)
	_ = utils.EnsureDir(gp.Markers)
	_ = utils.EnsureDir(gp.Transcodes)
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
