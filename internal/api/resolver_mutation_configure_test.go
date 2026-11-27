package api

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/models"
)

func TestSetVideoCodec(t *testing.T) {
	hevc := models.GeneratedVideoCodecHevc
	av1 := models.GeneratedVideoCodecAv1
	onlyH264AndHevc := func(c models.GeneratedVideoCodec) bool { return c != models.GeneratedVideoCodecAv1 }

	c := config.InitializeEmpty()

	assert.NoError(t, setVideoCodec(c, config.PreviewCodec, nil, c.GetPreviewCodec(), onlyH264AndHevc), "nil is a no-op")
	assert.Equal(t, models.GeneratedVideoCodecH264, c.GetPreviewCodec())

	assert.NoError(t, setVideoCodec(c, config.PreviewCodec, &hevc, c.GetPreviewCodec(), onlyH264AndHevc))
	assert.Equal(t, hevc, c.GetPreviewCodec())

	assert.Error(t, setVideoCodec(c, config.PreviewCodec, &av1, c.GetPreviewCodec(), onlyH264AndHevc), "unsupported codec rejected")
	assert.Equal(t, hevc, c.GetPreviewCodec(), "rejected codec not stored")

	// ffmpeg lost HEVC after it was saved: re-sending the unchanged value must not fail the save
	noneSupported := func(models.GeneratedVideoCodec) bool { return false }
	assert.NoError(t, setVideoCodec(c, config.PreviewCodec, &hevc, c.GetPreviewCodec(), noneSupported))
}
