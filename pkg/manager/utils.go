package manager

import (
	"fmt"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

func IsStreamable(scene *models.Scene) (bool, error) {
	if scene == nil {
		return false, fmt.Errorf("nil scene")
	}

	videoCodec := scene.VideoCodec.String
	if ffmpeg.IsValidCodec(videoCodec) {
		return true, nil
	} else {
		hasTranscode, _ := HasTranscode(scene)
		return hasTranscode, nil
	}
}

func HasTranscode(scene *models.Scene) (bool, error) {
	if scene == nil {
		return false, fmt.Errorf("nil scene")
	}
	transcodePath := instance.Paths.Scene.GetTranscodePath(scene.Checksum)
	return utils.FileExists(transcodePath)
}
