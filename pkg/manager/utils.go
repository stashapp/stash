package manager

import (
	"fmt"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

// IsStreamable returns true if the scene video can be streamed directly
// without live transcoding. A video can be streamed directly if it has
// a transcode already, or if video and audio codecs and container are
// compatible with browser viewing.
func IsStreamable(scene *models.Scene, useMD5 bool) (bool, error) {
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
	audioCodec := ffmpeg.MissingUnsupported
	if scene.AudioCodec.Valid {
		audioCodec = ffmpeg.AudioCodec(scene.AudioCodec.String)
	}

	if ffmpeg.IsValidCodec(videoCodec) && ffmpeg.IsValidCombo(videoCodec, container) && ffmpeg.IsValidAudioForContainer(audioCodec, container) {
		logger.Debugf("File is streamable %s, %s, %s\n", videoCodec, audioCodec, container)
		return true, nil
	} else {
		hasTranscode := HasTranscode(scene, useMD5)
		logger.Debugf("File is not streamable , transcode is needed  %s, %s, %s\n", videoCodec, audioCodec, container)
		return hasTranscode, nil
	}
}

// HasTranscode returns true if a transcoded video exists for the provided
// scene. It will check using the OSHash of the scene first, then fall back
// to the checksum.
func HasTranscode(scene *models.Scene, useMD5 bool) bool {
	if scene == nil {
		return false
	}

	sceneHash := scene.GetHash(useMD5)
	if sceneHash == "" {
		return false
	}

	transcodePath := instance.Paths.Scene.GetTranscodePath(sceneHash)
	ret, _ := utils.FileExists(transcodePath)
	return ret
}
