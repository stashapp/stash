package api

import (
	"github.com/stashapp/stash/pkg/models"
)

type BaseFile interface{}

type GalleryFile struct {
	*models.BaseFile
}
