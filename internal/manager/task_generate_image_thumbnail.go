package manager

import (
	"context"
	"errors"
	"fmt"
	"os/exec"

	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type GenerateImageThumbnailTask struct {
	Image     models.Image
	Overwrite bool
}

func (t *GenerateImageThumbnailTask) GetDescription() string {
	return fmt.Sprintf("Generating Thumbnail for image %s", t.Image.Path)
}

func (t *GenerateImageThumbnailTask) logStderr(err error) {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		logger.Debugf("[generator] error output: %s", exitErr.Stderr)
	}
}

func (t *GenerateImageThumbnailTask) Start(ctx context.Context) {
	if !t.required() {
		return
	}

	thumbPath := GetInstance().Paths.Generated.GetThumbnailPath(t.Image.Checksum, models.DefaultGthumbWidth)
	f := t.Image.Files.Primary()
	path := f.Base().Path

	logger.Debugf("Generating thumbnail for %s", path)

	mgr := GetInstance()
	c := mgr.Config

	clipPreviewOptions := image.ClipPreviewOptions{
		InputArgs:  c.GetTranscodeInputArgs(),
		OutputArgs: c.GetTranscodeOutputArgs(),
		Preset:     c.GetPreviewPreset().String(),
	}

	encoder := image.NewThumbnailEncoder(mgr.FFMpeg, mgr.FFProbe, clipPreviewOptions)
	data, err := encoder.GetThumbnail(f, models.DefaultGthumbWidth)

	if err != nil {
		// don't log for animated images
		if !errors.Is(err, image.ErrNotSupportedForThumbnail) {
			logger.Errorf("[generator] getting thumbnail for image %s: %s", path, err.Error())
			t.logStderr(err)
		}
		return
	}

	err = fsutil.WriteFile(thumbPath, data)
	if err != nil {
		logger.Errorf("[generator] writing thumbnail for image %s: %s", path, err.Error())
		return
	}
}

func (t *GenerateImageThumbnailTask) required() bool {
	vf, ok := t.Image.Files.Primary().(models.VisualFile)
	if !ok {
		return false
	}

	if vf.GetHeight() <= models.DefaultGthumbWidth && vf.GetWidth() <= models.DefaultGthumbWidth {
		return false
	}

	if t.Overwrite {
		return true
	}

	thumbPath := GetInstance().Paths.Generated.GetThumbnailPath(t.Image.Checksum, models.DefaultGthumbWidth)
	exists, _ := fsutil.FileExists(thumbPath)

	return !exists
}
