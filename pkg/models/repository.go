package models

type Repository interface {
	Gallery() GalleryReaderWriter
	Image() ImageReaderWriter
	Join() JoinReaderWriter
	Movie() MovieReaderWriter
	Performer() PerformerReaderWriter
	SceneMarker() SceneMarkerReaderWriter
	Scene() SceneReaderWriter
	Studio() StudioReaderWriter
	Tag() TagReaderWriter
}
