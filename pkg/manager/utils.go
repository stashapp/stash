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
	var container ffmpeg.Container
	if scene.Format.Valid {
		container = ffmpeg.Container(scene.Format.String)
	} else { // container isn't in the DB
		// shouldn't happen, fallback to ffprobe reading from file
		tmpVideoFile, err := ffmpeg.NewVideoFile(instance.FFProbePath, scene.Path)
		if err != nil {
			return false, fmt.Errorf("error reading video file: %s", err.Error())
		}
		container = ffmpeg.MatchContainer(tmpVideoFile.Container, scene.Path)
	}

	videoCodec := scene.VideoCodec.String
	if ffmpeg.IsValidCodec(videoCodec) && ffmpeg.IsValidCombo(videoCodec, container) {
		fmt.Printf("File is streamable %s,%s\n", videoCodec, container)
		return true, nil
	} else {
		hasTranscode, _ := HasTranscode(scene)
		fmt.Printf("File is not streamable , transcode is needed  %s,%s\n", videoCodec, container)
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
