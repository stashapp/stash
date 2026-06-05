package urlbuilders

import (
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

type MangaURLBuilder struct {
	BaseURL   string
	MangaID   string
	UpdatedAt string
}

func NewMangaURLBuilder(baseURL string, manga *models.Manga) MangaURLBuilder {
	return MangaURLBuilder{
		BaseURL:   baseURL,
		MangaID:   strconv.Itoa(manga.ID),
		UpdatedAt: strconv.FormatInt(manga.UpdatedAt.Unix(), 10),
	}
}

func (b MangaURLBuilder) GetCoverURL() string {
	return b.BaseURL + "/manga/" + b.MangaID + "/cover?t=" + b.UpdatedAt
}
