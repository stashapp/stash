package paths

import (
	"fmt"
	"github.com/stashapp/stash/utils"
	"path/filepath"
	"runtime"
	"strings"
)

type fixedPaths struct {
	ExecutionDirectory string
	ConfigDirectory    string
	ConfigFile         string
	DatabaseFile       string

	FFMPEG  string
	FFProbe string
}

func newFixedPaths() *fixedPaths {
	fp := fixedPaths{}
	fp.ExecutionDirectory = getExecutionDirectory()
	fp.ConfigDirectory = filepath.Join(getHomeDirectory(), ".stash")
	fp.ConfigFile = filepath.Join(fp.ConfigDirectory, "config.json")
	fp.DatabaseFile = filepath.Join(fp.ConfigDirectory, "stash-go.sqlite")

	ffmpegDirectories := []string{fp.ExecutionDirectory, fp.ConfigDirectory}
	ffmpegFileName := func() string {
		if runtime.GOOS == "windows" {
			return "ffmpeg.exe"
		} else {
			return "ffmpeg"
		}
	}()
	ffprobeFileName := func() string {
		if runtime.GOOS == "windows" {
			return "ffprobe.exe"
		} else {
			return "ffprobe"
		}
	}()
	for _, directory := range ffmpegDirectories {
		ffmpegPath := filepath.Join(directory, ffmpegFileName)
		ffprobePath := filepath.Join(directory, ffprobeFileName)
		if exists, _ := utils.FileExists(ffmpegPath); exists {
			fp.FFMPEG = ffmpegPath
		}
		if exists, _ := utils.FileExists(ffprobePath); exists {
			fp.FFProbe = ffprobePath
		}
	}

	errorText := fmt.Sprintf(
		"FFMPEG or FFProbe not found.  Place it in one of the following folders:\n\n%s",
		strings.Join(ffmpegDirectories, ","),
	)
	if exists, _ := utils.FileExists(fp.FFMPEG); !exists {
		panic(errorText)
	}
	if exists, _ := utils.FileExists(fp.FFProbe); !exists {
		panic(errorText)
	}

	return &fp
}