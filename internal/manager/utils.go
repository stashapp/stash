package manager

import "github.com/stashapp/stash/internal/utils"

func IsStreamable(videoPath string, checksum string) (bool, error) {
	fileType, err := utils.FileType(videoPath)
	if err != nil {
		return false, err
	}

	if fileType.MIME.Value == "video/quicktime" || fileType.MIME.Value == "video/mp4" || fileType.MIME.Value == "video/webm" || fileType.MIME.Value == "video/x-m4v" {
		return true, nil
	} else {
		transcodePath := instance.Paths.Scene.GetTranscodePath(checksum)
		return utils.FileExists(transcodePath)
	}
}