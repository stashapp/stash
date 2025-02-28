package hook

type TriggerEnum string

// Scan-related hooks are current disabled until post-hook execution is
// integrated.

const (
	SceneMarkerCreatePost  TriggerEnum = "SceneMarker.Create.Post"
	SceneMarkerUpdatePost  TriggerEnum = "SceneMarker.Update.Post"
	SceneMarkerDestroyPost TriggerEnum = "SceneMarker.Destroy.Post"

	SceneCreatePost  TriggerEnum = "Scene.Create.Post"
	SceneUpdatePost  TriggerEnum = "Scene.Update.Post"
	SceneDestroyPost TriggerEnum = "Scene.Destroy.Post"

	ImageCreatePost  TriggerEnum = "Image.Create.Post"
	ImageUpdatePost  TriggerEnum = "Image.Update.Post"
	ImageDestroyPost TriggerEnum = "Image.Destroy.Post"

	GalleryCreatePost  TriggerEnum = "Gallery.Create.Post"
	GalleryUpdatePost  TriggerEnum = "Gallery.Update.Post"
	GalleryDestroyPost TriggerEnum = "Gallery.Destroy.Post"

	GalleryChapterCreatePost  TriggerEnum = "GalleryChapter.Create.Post"
	GalleryChapterUpdatePost  TriggerEnum = "GalleryChapter.Update.Post"
	GalleryChapterDestroyPost TriggerEnum = "GalleryChapter.Destroy.Post"

	// deprecated - use Group hooks instead
	// for now, both movie and group hooks will be executed
	MovieCreatePost  TriggerEnum = "Movie.Create.Post"
	MovieUpdatePost  TriggerEnum = "Movie.Update.Post"
	MovieDestroyPost TriggerEnum = "Movie.Destroy.Post"

	GroupCreatePost  TriggerEnum = "Group.Create.Post"
	GroupUpdatePost  TriggerEnum = "Group.Update.Post"
	GroupDestroyPost TriggerEnum = "Group.Destroy.Post"

	PerformerCreatePost  TriggerEnum = "Performer.Create.Post"
	PerformerUpdatePost  TriggerEnum = "Performer.Update.Post"
	PerformerDestroyPost TriggerEnum = "Performer.Destroy.Post"

	StudioCreatePost  TriggerEnum = "Studio.Create.Post"
	StudioUpdatePost  TriggerEnum = "Studio.Update.Post"
	StudioDestroyPost TriggerEnum = "Studio.Destroy.Post"

	TagCreatePost  TriggerEnum = "Tag.Create.Post"
	TagUpdatePost  TriggerEnum = "Tag.Update.Post"
	TagMergePost   TriggerEnum = "Tag.Merge.Post"
	TagDestroyPost TriggerEnum = "Tag.Destroy.Post"
)

var AllHookTriggerEnum = []TriggerEnum{
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

	GalleryChapterCreatePost,
	GalleryChapterUpdatePost,
	GalleryChapterDestroyPost,

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

func (e TriggerEnum) IsValid() bool {

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

		GalleryChapterCreatePost,
		GalleryChapterUpdatePost,
		GalleryChapterDestroyPost,

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

func (e TriggerEnum) String() string {
	return string(e)
}
