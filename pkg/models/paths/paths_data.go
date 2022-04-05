package paths

import (
	"fmt"
	"path/filepath"

	"github.com/stashapp/stash/pkg/fsutil"
)

const dataDirLength = 3

type dataPaths struct {
	Scenes string
}

func newDataPaths(path string) *dataPaths {
	gp := dataPaths{}
	gp.Scenes = filepath.Join(path, "scenes")
	return &gp
}

func (gp *dataPaths) GetSceneCoverPath(id int) string {
	fname := fmt.Sprintf("%d_cover.jpg", id)
	return filepath.Join(gp.Scenes, fsutil.GetIntraDirID(id, dataDirLength), fname)
}
