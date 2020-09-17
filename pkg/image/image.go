package image

import (
	"database/sql"
	"image"
	"os"

	"github.com/stashapp/stash/pkg/models"
)

func GetSourceImage(i *models.Image) (image.Image, error) {
	f, err := os.Open(i.Path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	srcImage, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return srcImage, nil
}

func SetFileDetails(i *models.Image) error {
	f, err := os.Stat(i.Path)
	if err != nil {
		return err
	}

	src, err := GetSourceImage(i)
	if err != nil {
		return err
	}

	i.Width = sql.NullInt64{
		Int64: int64(src.Bounds().Max.X),
		Valid: true,
	}
	i.Height = sql.NullInt64{
		Int64: int64(src.Bounds().Max.Y),
		Valid: true,
	}
	i.Size = sql.NullInt64{
		Int64: int64(f.Size()),
	}

	return nil
}
