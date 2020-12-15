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

type ReaderRepository interface {
	Gallery() GalleryReader
	Image() ImageReader
	Join() JoinReader
	Movie() MovieReader
	Performer() PerformerReader
	SceneMarker() SceneMarkerReader
	Scene() SceneReader
	Studio() StudioReader
	Tag() TagReader
}
