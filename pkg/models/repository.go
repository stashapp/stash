package models

type Repository interface {
	Gallery() GalleryReaderWriter
	Image() ImageReaderWriter
	Join() JoinReaderWriter
	Movie() MovieReaderWriter
	SceneMarker() SceneMarkerReaderWriter
	Scene() SceneReaderWriter
	Studio() StudioReaderWriter
	Tag() TagReaderWriter
}
