package image

import (
	"github.com/stashapp/stash/pkg/models"
)

type Service struct {
	File       models.FileReaderWriter
	Repository models.ImageReaderWriter
}
