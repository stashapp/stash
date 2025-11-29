package postgres

import "github.com/stashapp/stash/pkg/database"

func (db *Database) Blobs() database.BlobStore {
	return db.storeRepository.Blobs
}
func (db *Database) File() database.FileStore {
	return db.storeRepository.File
}
func (db *Database) Folder() database.FolderStore {
	return db.storeRepository.Folder
}
func (db *Database) Image() database.ImageStore {
	return db.storeRepository.Image
}
func (db *Database) Gallery() database.GalleryStore {
	return db.storeRepository.Gallery
}
func (db *Database) GalleryChapter() database.GalleryChapterStore {
	return db.storeRepository.GalleryChapter
}
func (db *Database) Scene() database.SceneStore {
	return db.storeRepository.Scene
}
func (db *Database) SceneMarker() database.SceneMarkerStore {
	return db.storeRepository.SceneMarker
}
func (db *Database) Performer() database.PerformerStore {
	return db.storeRepository.Performer
}
func (db *Database) SavedFilter() database.SavedFilterStore {
	return db.storeRepository.SavedFilter
}
func (db *Database) Studio() database.StudioStore {
	return db.storeRepository.Studio
}
func (db *Database) Tag() database.TagStore {
	return db.storeRepository.Tag
}
func (db *Database) Group() database.GroupStore {
	return db.storeRepository.Group
}
func (db *Database) NewMigrator() (database.MigrateStore, error) {
	return NewMigrator(db)
}
