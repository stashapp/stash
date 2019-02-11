package paths

import (
	"os"
	"os/user"
	"path/filepath"
)

type StaticPathsType struct {
	ExecutionDirectory string
	ConfigDirectory    string
	ConfigFile         string
	DatabaseFile       string

	FFMPEG  string
	FFProbe string
}

var StaticPaths = StaticPathsType{
	ExecutionDirectory: getExecutionDirectory(),
	ConfigDirectory: getConfigDirectory(),
	ConfigFile: filepath.Join(getConfigDirectory(), "config.json"),
	DatabaseFile: filepath.Join(getConfigDirectory(), "stash-go.sqlite"),
}

func getExecutionDirectory() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}

func getHomeDirectory() string {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	return currentUser.HomeDir
}

func getConfigDirectory() string {
	return filepath.Join(getHomeDirectory(), ".stash")
}