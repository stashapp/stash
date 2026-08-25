// Package audio provides the application logic for audio functionality.
// Most functionality is provided by [Service].
package audio

import (
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/paths"
	"github.com/stashapp/stash/pkg/plugin"
)

type Service struct {
	File        models.FileReaderWriter
	Repository  models.AudioReaderWriter
	PluginCache *plugin.Cache

	Paths *paths.Paths
}
