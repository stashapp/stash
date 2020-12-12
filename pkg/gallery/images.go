package gallery

import (
	"github.com/stashapp/stash/pkg/api/urlbuilders"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sqlite"
)

func GetFiles(g *models.Gallery, baseURL string) []*models.GalleryFilesType {
	var galleryFiles []*models.GalleryFilesType

	qb := sqlite.NewImageQueryBuilder()
	images, err := qb.FindByGalleryID(g.ID)
	if err != nil {
		return nil
	}

	for i, img := range images {
		builder := urlbuilders.NewImageURLBuilder(baseURL, img.ID)
		imageURL := builder.GetImageURL()

		galleryFile := models.GalleryFilesType{
			Index: i,
			Name:  &img.Title.String,
			Path:  &imageURL,
		}
		galleryFiles = append(galleryFiles, &galleryFile)
	}

	return galleryFiles
}
