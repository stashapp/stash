// TODO(audio): update this file

// Package audio provides the application logic for audio functionality.
// Most functionality is provided by [Service].
package audio

import (
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/paths"
	"github.com/stashapp/stash/pkg/plugin"
)

type Config interface {
	GetVideoFileNamingAlgorithm() models.HashAlgorithm
}

type Service struct {
	File             models.FileReaderWriter
	Repository       models.AudioReaderWriter
	MarkerRepository models.AudioMarkerReaderWriter
	PluginCache      *plugin.Cache

	Paths  *paths.Paths
	Config Config
}
