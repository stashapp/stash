// Package image provides the application logic for images.
// The functionality is exposed via the [Service] type.
package image

import (
	"github.com/stashapp/stash/pkg/models"
)

type Service struct {
	File       models.FileReaderWriter
	Repository models.ImageReaderWriter
}
