package paths

import "path/filepath"

type dataPaths struct {
	Scenes string
}

func newDataPaths(path string) *dataPaths {
	gp := dataPaths{}
	gp.Scenes = filepath.Join(path, "scenes")
	return &gp
}
