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

type HookTriggerEnum string

// Scan-related hooks are current disabled until post-hook execution is
// integrated.

const (
	SceneMarkerCreatePost  HookTriggerEnum = "SceneMarker.Create.Post"
	SceneMarkerUpdatePost  HookTriggerEnum = "SceneMarker.Update.Post"
	SceneMarkerDestroyPost HookTriggerEnum = "SceneMarker.Destroy.Post"

	SceneCreatePost  HookTriggerEnum = "Scene.Create.Post"
	SceneUpdatePost  HookTriggerEnum = "Scene.Update.Post"
	SceneDestroyPost HookTriggerEnum = "Scene.Destroy.Post"

	ImageCreatePost  HookTriggerEnum = "Image.Create.Post"
	ImageUpdatePost  HookTriggerEnum = "Image.Update.Post"
	ImageDestroyPost HookTriggerEnum = "Image.Destroy.Post"

	GalleryCreatePost  HookTriggerEnum = "Gallery.Create.Post"
	GalleryUpdatePost  HookTriggerEnum = "Gallery.Update.Post"
	GalleryDestroyPost HookTriggerEnum = "Gallery.Destroy.Post"

	MovieCreatePost  HookTriggerEnum = "Movie.Create.Post"
	MovieUpdatePost  HookTriggerEnum = "Movie.Update.Post"
	MovieDestroyPost HookTriggerEnum = "Movie.Destroy.Post"

	PerformerCreatePost  HookTriggerEnum = "Performer.Create.Post"
	PerformerUpdatePost  HookTriggerEnum = "Performer.Update.Post"
	PerformerDestroyPost HookTriggerEnum = "Performer.Destroy.Post"

	StudioCreatePost  HookTriggerEnum = "Studio.Create.Post"
	StudioUpdatePost  HookTriggerEnum = "Studio.Update.Post"
	StudioDestroyPost HookTriggerEnum = "Studio.Destroy.Post"

	TagCreatePost  HookTriggerEnum = "Tag.Create.Post"
	TagUpdatePost  HookTriggerEnum = "Tag.Update.Post"
	TagMergePost   HookTriggerEnum = "Tag.Merge.Post"
	TagDestroyPost HookTriggerEnum = "Tag.Destroy.Post"
)

var AllHookTriggerEnum = []HookTriggerEnum{
	SceneMarkerCreatePost,
	SceneMarkerUpdatePost,
	SceneMarkerDestroyPost,

	SceneCreatePost,
	SceneUpdatePost,
	SceneDestroyPost,

	ImageCreatePost,
	ImageUpdatePost,
	ImageDestroyPost,

	GalleryCreatePost,
	GalleryUpdatePost,
	GalleryDestroyPost,

	MovieCreatePost,
	MovieUpdatePost,
	MovieDestroyPost,

	PerformerCreatePost,
	PerformerUpdatePost,
	PerformerDestroyPost,

	StudioCreatePost,
	StudioUpdatePost,
	StudioDestroyPost,

	TagCreatePost,
	TagUpdatePost,
	TagMergePost,
	TagDestroyPost,
}

func (e HookTriggerEnum) IsValid() bool {

	switch e {
	case SceneMarkerCreatePost,
		SceneMarkerUpdatePost,
		SceneMarkerDestroyPost,

		SceneCreatePost,
		SceneUpdatePost,
		SceneDestroyPost,

		ImageCreatePost,
		ImageUpdatePost,
		ImageDestroyPost,

		GalleryCreatePost,
		GalleryUpdatePost,
		GalleryDestroyPost,

		MovieCreatePost,
		MovieUpdatePost,
		MovieDestroyPost,

		PerformerCreatePost,
		PerformerUpdatePost,
		PerformerDestroyPost,

		StudioCreatePost,
		StudioUpdatePost,
		StudioDestroyPost,

		TagCreatePost,
		TagUpdatePost,
		TagDestroyPost:
		return true
	}
	return false
}

func (e HookTriggerEnum) String() string {
	return string(e)
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
