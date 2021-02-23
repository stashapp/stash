package image

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestIsCover(t *testing.T) {
	const cover = "cover.jpg"
	const notCover = "notcover.jpg"
	const capitalCover = "Cover.jpg"
	subDirCover := fmt.Sprintf("subDir%scover.jpg", string(filepath.Separator))

	img := &models.Image{
		Path: cover,
	}

	assert := assert.New(t)
	assert.True(IsCover(img))

	img.Path = notCover
	assert.False(IsCover(img))

	img.Path = capitalCover
	assert.False(IsCover(img))

	img.Path = subDirCover
	assert.True(IsCover(img))
}
