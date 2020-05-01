// +build integration
package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stashapp/stash/pkg/models"
)

func TestGalleryFind(t *testing.T) {
	gqb := models.NewGalleryQueryBuilder()

	const galleryIdx = 0
	gallery, err := gqb.Find(galleryIDs[galleryIdx])

	if err != nil {
		t.Fatalf("Error finding gallery: %s", err.Error())
	}

	assert.Equal(t, getGalleryStringValue(galleryIdx, "Path"), gallery.Path)

	gallery, err = gqb.Find(0)

	if err != nil {
		t.Fatalf("Error finding gallery: %s", err.Error())
	}

	assert.Nil(t, gallery)
}

func TestGalleryFindByChecksum(t *testing.T) {
	gqb := models.NewGalleryQueryBuilder()

	const galleryIdx = 0
	galleryChecksum := getGalleryStringValue(galleryIdx, "Checksum")
	gallery, err := gqb.FindByChecksum(galleryChecksum, nil)

	if err != nil {
		t.Fatalf("Error finding gallery: %s", err.Error())
	}

	assert.Equal(t, getGalleryStringValue(galleryIdx, "Path"), gallery.Path)

	galleryChecksum = "not exist"
	gallery, err = gqb.FindByChecksum(galleryChecksum, nil)

	if err != nil {
		t.Fatalf("Error finding gallery: %s", err.Error())
	}

	assert.Nil(t, gallery)
}

func TestGalleryFindByPath(t *testing.T) {
	gqb := models.NewGalleryQueryBuilder()

	const galleryIdx = 0
	galleryPath := getGalleryStringValue(galleryIdx, "Path")
	gallery, err := gqb.FindByPath(galleryPath)

	if err != nil {
		t.Fatalf("Error finding gallery: %s", err.Error())
	}

	assert.Equal(t, galleryPath, gallery.Path)

	galleryPath = "not exist"
	gallery, err = gqb.FindByPath(galleryPath)

	if err != nil {
		t.Fatalf("Error finding gallery: %s", err.Error())
	}

	assert.Nil(t, gallery)
}

func TestGalleryFindBySceneID(t *testing.T) {
	gqb := models.NewGalleryQueryBuilder()

	sceneID := sceneIDs[sceneIdxWithGallery]
	gallery, err := gqb.FindBySceneID(sceneID, nil)

	if err != nil {
		t.Fatalf("Error finding gallery: %s", err.Error())
	}

	assert.Equal(t, getGalleryStringValue(galleryIdxWithScene, "Path"), gallery.Path)

	gallery, err = gqb.FindBySceneID(0, nil)

	if err != nil {
		t.Fatalf("Error finding gallery: %s", err.Error())
	}

	assert.Nil(t, gallery)
}

// TODO ValidGalleriesForScenePath
// TODO Count
// TODO All
// TODO Query
// TODO Update
// TODO Destroy
// TODO ClearGalleryId
