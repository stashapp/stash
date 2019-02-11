package manager

import (
	"github.com/stashapp/stash/ffmpeg"
	"github.com/stashapp/stash/logger"
	"github.com/stashapp/stash/manager/paths"
	"sync"
)

type singleton struct {
	Status      JobStatus
	Paths       *paths.Paths
	StaticPaths *paths.StaticPathsType
	JSON        *jsonUtils
}

var instance *singleton
var once sync.Once

func GetInstance() *singleton {
	Initialize()
	return instance
}

func Initialize() *singleton {
	once.Do(func() {
		instance = &singleton{
			Status:      Idle,
			Paths:       paths.RefreshPaths(),
			StaticPaths: &paths.StaticPaths,
			JSON:        &jsonUtils{},
		}

		initFFMPEG()
	})

	return instance
}

func initFFMPEG() {
	ffmpegPath, ffprobePath := ffmpeg.GetPaths(instance.StaticPaths.ConfigDirectory)
	if ffmpegPath == "" || ffprobePath == "" {
		logger.Infof("couldn't find FFMPEG, attempting to download it")
		if err := ffmpeg.Download(instance.StaticPaths.ConfigDirectory); err != nil {
			msg := `Unable to locate / automatically download FFMPEG

Check the readme for download links.
The FFMPEG and FFProbe binaries should be placed in %s

The error was: %s
`
			logger.Fatalf(msg, instance.StaticPaths.ConfigDirectory, err)
		}
	}

	instance.StaticPaths.FFMPEG = ffmpegPath
	instance.StaticPaths.FFProbe = ffprobePath
}
