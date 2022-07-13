package image

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestIsCover(t *testing.T) {
	type test struct {
		fn      string
		isCover bool
	}

	tests := []test{
		{"cover.jpg", true},
		{"covernot.jpg", false},
		{"Cover.jpg", false},
		{fmt.Sprintf("subDir%scover.jpg", string(filepath.Separator)), true},
		{"endsWithcover.jpg", true},
		{"cover.png", false},
	}

	assert := assert.New(t)
	for _, tc := range tests {
		img := &models.Image{
			Files: []*file.ImageFile{
				{
					BaseFile: &file.BaseFile{
						Path: tc.fn,
					},
				},
			},
		}
		assert.Equal(tc.isCover, IsCover(img), "expected: %t for %s", tc.isCover, tc.fn)
	}
}
