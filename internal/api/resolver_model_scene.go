package api

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/stashapp/stash/internal/api/urlbuilders"
	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

func (r *sceneResolver) FileModTime(ctx context.Context, obj *models.Scene) (*time.Time, error) {
	if obj.PrimaryFile() != nil {
		return &obj.PrimaryFile().ModTime, nil
	}
	return nil, nil
}

func (r *sceneResolver) Date(ctx context.Context, obj *models.Scene) (*string, error) {
	if obj.Date != nil {
		result := obj.Date.String()
		return &result, nil
	}
	return nil, nil
}

// File is deprecated
func (r *sceneResolver) File(ctx context.Context, obj *models.Scene) (*models.SceneFileType, error) {
	f := obj.PrimaryFile()
	if f == nil {
		return nil, nil
	}

	bitrate := int(f.BitRate)
	size := strconv.FormatInt(f.Size, 10)

	return &models.SceneFileType{
		Size:       &size,
		Duration:   handleFloat64(f.Duration),
		VideoCodec: &f.VideoCodec,
		AudioCodec: &f.AudioCodec,
		Width:      &f.Width,
		Height:     &f.Height,
		Framerate:  handleFloat64(f.FrameRate),
		Bitrate:    &bitrate,
	}, nil
}

func (r *sceneResolver) Files(ctx context.Context, obj *models.Scene) ([]*VideoFile, error) {
	ret := make([]*VideoFile, len(obj.Files))

	for i, f := range obj.Files {
		ret[i] = &VideoFile{
			ID:             strconv.Itoa(int(f.ID)),
			Path:           f.Path,
			Basename:       f.Basename,
			ParentFolderID: strconv.Itoa(int(f.ParentFolderID)),
			ModTime:        f.ModTime,
			Format:         f.Format,
			Size:           f.Size,
			Duration:       handleFloat64Value(f.Duration),
			VideoCodec:     f.VideoCodec,
			AudioCodec:     f.AudioCodec,
			Width:          f.Width,
			Height:         f.Height,
			FrameRate:      handleFloat64Value(f.FrameRate),
			BitRate:        int(f.BitRate),
			CreatedAt:      f.CreatedAt,
			UpdatedAt:      f.UpdatedAt,
			Fingerprints:   resolveFingerprints(f.Base()),
		}

		if f.ZipFileID != nil {
			zipFileID := strconv.Itoa(int(*f.ZipFileID))
			ret[i].ZipFileID = &zipFileID
		}
	}

	return ret, nil
}

func resolveFingerprints(f *file.BaseFile) []*Fingerprint {
	ret := make([]*Fingerprint, len(f.Fingerprints))

	for i, fp := range f.Fingerprints {
		ret[i] = &Fingerprint{
			Type:  fp.Type,
			Value: formatFingerprint(fp.Fingerprint),
		}
	}

	return ret
}

func formatFingerprint(fp interface{}) string {
	switch v := fp.(type) {
	case int64:
		return strconv.FormatUint(uint64(v), 16)
	default:
		return fmt.Sprintf("%v", fp)
	}
}

func (r *sceneResolver) Paths(ctx context.Context, obj *models.Scene) (*ScenePathsType, error) {
	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	config := manager.GetInstance().Config
	builder := urlbuilders.NewSceneURLBuilder(baseURL, obj.ID)
	builder.APIKey = config.GetAPIKey()
	screenshotPath := builder.GetScreenshotURL(obj.UpdatedAt)
	previewPath := builder.GetStreamPreviewURL()
	streamPath := builder.GetStreamURL()
	webpPath := builder.GetStreamPreviewImageURL()
	vttPath := builder.GetSpriteVTTURL()
	spritePath := builder.GetSpriteURL()
	chaptersVttPath := builder.GetChaptersVTTURL()
	funscriptPath := builder.GetFunscriptURL()
	captionBasePath := builder.GetCaptionURL()
	interactiveHeatmap := builder.GetInteractiveHeatmapURL()

	return &ScenePathsType{
		Screenshot:         &screenshotPath,
		Preview:            &previewPath,
		Stream:             &streamPath,
		Webp:               &webpPath,
		Vtt:                &vttPath,
		ChaptersVtt:        &chaptersVttPath,
		Sprite:             &spritePath,
		Funscript:          &funscriptPath,
		InteractiveHeatmap: &interactiveHeatmap,
		Caption:            &captionBasePath,
	}, nil
}

func (r *sceneResolver) SceneMarkers(ctx context.Context, obj *models.Scene) (ret []*models.SceneMarker, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.SceneMarker.FindBySceneID(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *sceneResolver) Captions(ctx context.Context, obj *models.Scene) (ret []*models.VideoCaption, err error) {
	primaryFile := obj.PrimaryFile()
	if primaryFile == nil {
		return nil, nil
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.File.GetCaptions(ctx, primaryFile.Base().ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, err
}

func (r *sceneResolver) Galleries(ctx context.Context, obj *models.Scene) (ret []*models.Gallery, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Gallery.FindBySceneID(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *sceneResolver) Studio(ctx context.Context, obj *models.Scene) (ret *models.Studio, err error) {
	if obj.StudioID == nil {
		return nil, nil
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Studio.Find(ctx, *obj.StudioID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *sceneResolver) Movies(ctx context.Context, obj *models.Scene) (ret []*SceneMovie, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		mqb := r.repository.Movie

		for _, sm := range obj.Movies {
			movie, err := mqb.Find(ctx, sm.MovieID)
			if err != nil {
				return err
			}

			sceneIdx := sm.SceneIndex
			sceneMovie := &SceneMovie{
				Movie:      movie,
				SceneIndex: sceneIdx,
			}

			ret = append(ret, sceneMovie)
		}

		return nil
	}); err != nil {
		return nil, err
	}
	return ret, nil
}

func (r *sceneResolver) Tags(ctx context.Context, obj *models.Scene) (ret []*models.Tag, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Tag.FindBySceneID(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *sceneResolver) Performers(ctx context.Context, obj *models.Scene) (ret []*models.Performer, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Performer.FindBySceneID(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *sceneResolver) Phash(ctx context.Context, obj *models.Scene) (*string, error) {
	phash := obj.Phash()
	if phash != 0 {
		hexval := utils.PhashToString(phash)
		return &hexval, nil
	}
	return nil, nil
}

func (r *sceneResolver) SceneStreams(ctx context.Context, obj *models.Scene) ([]*manager.SceneStreamEndpoint, error) {
	config := manager.GetInstance().Config

	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	builder := urlbuilders.NewSceneURLBuilder(baseURL, obj.ID)

	return manager.GetSceneStreamPaths(obj, builder.GetStreamURL(), config.GetMaxStreamingTranscodeSize())
}

func (r *sceneResolver) Interactive(ctx context.Context, obj *models.Scene) (bool, error) {
	primaryFile := obj.PrimaryFile()
	if primaryFile == nil {
		return false, nil
	}

	return primaryFile.Interactive, nil
}

func (r *sceneResolver) InteractiveSpeed(ctx context.Context, obj *models.Scene) (*int, error) {
	primaryFile := obj.PrimaryFile()
	if primaryFile == nil {
		return nil, nil
	}

	return primaryFile.InteractiveSpeed, nil
}
