package models

type Repository interface {
	Gallery() GalleryReaderWriter
	Image() ImageReaderWriter
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
	Movie() MovieReader
	Performer() PerformerReader
	SceneMarker() SceneMarkerReader
	Scene() SceneReader
	Studio() StudioReader
	Tag() TagReader
}
