package plugin

import (
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin/common"
)

type PluginHook struct {
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	Hooks       []string `json:"hooks"`
	Plugin      *Plugin  `json:"plugin"`
}

func addHookContext(argsMap common.ArgsMap, hookContext common.HookContext) {
	argsMap[common.HookContextKey] = hookContext
}

// types for destroy hooks, to provide a little more information
type SceneDestroyInput struct {
	models.SceneDestroyInput
	Checksum string `json:"checksum"`
	OSHash   string `json:"oshash"`
	Path     string `json:"path"`
}

type ScenesDestroyInput struct {
	models.ScenesDestroyInput
	Checksum string `json:"checksum"`
	OSHash   string `json:"oshash"`
	Path     string `json:"path"`
}

type GalleryDestroyInput struct {
	models.GalleryDestroyInput
	Checksum string `json:"checksum"`
	Path     string `json:"path"`
}

type ImageDestroyInput struct {
	models.ImageDestroyInput
	Checksum string `json:"checksum"`
	Path     string `json:"path"`
}

type ImagesDestroyInput struct {
	models.ImagesDestroyInput
	Checksum string `json:"checksum"`
	Path     string `json:"path"`
}
