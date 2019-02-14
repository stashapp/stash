package manager

import (
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/jsonschema"
	"github.com/stashapp/stash/pkg/manager/paths"
	"github.com/stashapp/stash/pkg/utils"
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
		configFile := jsonschema.LoadConfigFile(paths.StaticPaths.ConfigFile)
		instance = &singleton{
			Status:      Idle,
			Paths:       paths.NewPaths(configFile),
			StaticPaths: &paths.StaticPaths,
			JSON:        &jsonUtils{},
		}

		instance.refreshConfig(configFile)

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

func HasValidConfig() bool {
	configFileExists, _ := utils.FileExists(instance.StaticPaths.ConfigFile) // TODO: Verify JSON is correct
	if configFileExists && instance.Paths.Config != nil {
		return true
	}
	return false
}

func (s *singleton) SaveConfig(config *jsonschema.Config) error {
	if err := jsonschema.SaveConfigFile(s.StaticPaths.ConfigFile, config); err != nil {
		return err
	}

	// Reload the config
	s.refreshConfig(config)

	return nil
}

func (s *singleton) refreshConfig(config *jsonschema.Config) {
	if config == nil {
		config = jsonschema.LoadConfigFile(s.StaticPaths.ConfigFile)
	}
	s.Paths = paths.NewPaths(config)

	if HasValidConfig() {
		_ = utils.EnsureDir(s.Paths.Generated.Screenshots)
		_ = utils.EnsureDir(s.Paths.Generated.Vtt)
		_ = utils.EnsureDir(s.Paths.Generated.Markers)
		_ = utils.EnsureDir(s.Paths.Generated.Transcodes)

		_ = utils.EnsureDir(s.Paths.JSON.Performers)
		_ = utils.EnsureDir(s.Paths.JSON.Scenes)
		_ = utils.EnsureDir(s.Paths.JSON.Galleries)
		_ = utils.EnsureDir(s.Paths.JSON.Studios)
	}
}
