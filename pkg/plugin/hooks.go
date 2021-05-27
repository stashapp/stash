package plugin

import (
	"github.com/stashapp/stash/pkg/plugin/common"
)

type HookTypeEnum string

const (
	SceneMarkerCreatePost  HookTypeEnum = "SceneMarker.Create.Post"
	SceneMarkerUpdatePost  HookTypeEnum = "SceneMarker.Update.Post"
	SceneMarkerDestroyPost HookTypeEnum = "SceneMarker.Destroy.Post"

	SceneCreatePost  HookTypeEnum = "Scene.Create.Post"
	SceneUpdatePost  HookTypeEnum = "Scene.Update.Post"
	SceneDestroyPost HookTypeEnum = "Scene.Destroy.Post"

	ImageCreatePost  HookTypeEnum = "Image.Create.Post"
	ImageUpdatePost  HookTypeEnum = "Image.Update.Post"
	ImageDestroyPost HookTypeEnum = "Image.Destroy.Post"

	GalleryCreatePost  HookTypeEnum = "Gallery.Create.Post"
	GalleryUpdatePost  HookTypeEnum = "Gallery.Update.Post"
	GalleryDestroyPost HookTypeEnum = "Gallery.Destroy.Post"

	MovieCreatePost  HookTypeEnum = "Movie.Create.Post"
	MovieUpdatePost  HookTypeEnum = "Movie.Update.Post"
	MovieDestroyPost HookTypeEnum = "Movie.Destroy.Post"

	PerformerCreatePost  HookTypeEnum = "Performer.Create.Post"
	PerformerUpdatePost  HookTypeEnum = "Performer.Update.Post"
	PerformerDestroyPost HookTypeEnum = "Performer.Destroy.Post"

	StudioCreatePost  HookTypeEnum = "Studio.Create.Post"
	StudioUpdatePost  HookTypeEnum = "Studio.Update.Post"
	StudioDestroyPost HookTypeEnum = "Studio.Destroy.Post"

	TagCreatePost  HookTypeEnum = "Tag.Create.Post"
	TagUpdatePost  HookTypeEnum = "Tag.Update.Post"
	TagDestroyPost HookTypeEnum = "Tag.Destroy.Post"
)

var AllHookTypeEnum = []HookTypeEnum{
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
	TagDestroyPost,
}

func (e HookTypeEnum) IsValid() bool {

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

func (e HookTypeEnum) String() string {
	return string(e)
}

func addHookContext(argsMap common.ArgsMap, hookType HookTypeEnum, input interface{}) {
	argsMap[common.HookContextKey] = common.HookContext{
		Type:  string(hookType),
		Input: input,
	}
}
