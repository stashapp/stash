package manager

import (
	"fmt"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

func IsStreamable(scene *models.Scene) (bool, error) {
	if scene == nil {
		return false, fmt.Errorf("nil scene")
	}
	fileType, err := utils.FileType(scene.Path)
	if err != nil {
		return false, err
	}

	switch fileType.MIME.Value {
	case "video/quicktime", "video/mp4", "video/webm", "video/x-m4v":
		return true, nil
	default:
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
