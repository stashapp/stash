package api

import (
	"context"
	"fmt"
	"time"

	"github.com/stashapp/stash/internal/api/loaders"
	"github.com/stashapp/stash/internal/api/urlbuilders"
	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/models"
)

func convertAudioFile(f models.File) (*models.AudioFile, error) {
	vf, ok := f.(*models.AudioFile)
	if !ok {
		return nil, fmt.Errorf("file %T is not a audio file", f)
	}
	return vf, nil
}

func (r *audioResolver) getPrimaryFile(ctx context.Context, obj *models.Audio) (*models.AudioFile, error) {
	if obj.PrimaryFileID != nil {
		f, err := loaders.From(ctx).FileByID.Load(*obj.PrimaryFileID)
		if err != nil {
			return nil, err
		}

		ret, err := convertAudioFile(f)
		if err != nil {
			return nil, err
		}

		obj.Files.SetPrimary(ret)

		return ret, nil
	} else {
		_ = obj.LoadPrimaryFile(ctx, r.repository.File)
	}

	return nil, nil
}

func (r *audioResolver) getFiles(ctx context.Context, obj *models.Audio) ([]*models.AudioFile, error) {
	fileIDs, err := loaders.From(ctx).AudioFiles.Load(obj.ID)
	if err != nil {
		return nil, err
	}

	files, errs := loaders.From(ctx).FileByID.LoadAll(fileIDs)
	err = firstError(errs)
	if err != nil {
		return nil, err
	}

	ret := make([]*models.AudioFile, len(files))
	for i, f := range files {
		ret[i], err = convertAudioFile(f)
		if err != nil {
			return nil, err
		}
	}

	obj.Files.Set(ret)

	return ret, nil
}

func (r *audioResolver) Date(ctx context.Context, obj *models.Audio) (*string, error) {
	if obj.Date != nil {
		result := obj.Date.String()
		return &result, nil
	}
	return nil, nil
}

func (r *audioResolver) Files(ctx context.Context, obj *models.Audio) ([]*AudioFile, error) {
	files, err := r.getFiles(ctx, obj)
	if err != nil {
		return nil, err
	}

	ret := make([]*AudioFile, len(files))

	for i, f := range files {
		ret[i] = &AudioFile{
			AudioFile: f,
		}
	}

	return ret, nil
}

func (r *audioResolver) Rating(ctx context.Context, obj *models.Audio) (*int, error) {
	if obj.Rating != nil {
		rating := models.Rating100To5(*obj.Rating)
		return &rating, nil
	}
	return nil, nil
}

func (r *audioResolver) Rating100(ctx context.Context, obj *models.Audio) (*int, error) {
	return obj.Rating, nil
}

func (r *audioResolver) Paths(ctx context.Context, obj *models.Audio) (*AudioPathsType, error) {
	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	config := manager.GetInstance().Config
	builder := urlbuilders.NewAudioURLBuilder(baseURL, obj)
	streamPath := builder.GetStreamURL(config.GetAPIKey()).String()
	captionBasePath := builder.GetCaptionURL()

	return &AudioPathsType{
		Stream:  &streamPath,
		Caption: &captionBasePath,
	}, nil
}

// TODO(audio|AudioCaption): need to update IF AudioCaption required
func (r *audioResolver) Captions(ctx context.Context, obj *models.Audio) (ret []*models.VideoCaption, err error) {
	primaryFile, err := r.getPrimaryFile(ctx, obj)
	if err != nil {
		return nil, err
	}
	if primaryFile == nil {
		return nil, nil
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.File.GetCaptions(ctx, primaryFile.Base().ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, err
}

func (r *audioResolver) Studio(ctx context.Context, obj *models.Audio) (ret *models.Studio, err error) {
	if obj.StudioID == nil {
		return nil, nil
	}

	return loaders.From(ctx).StudioByID.Load(*obj.StudioID)
}

func (r *audioResolver) Groups(ctx context.Context, obj *models.Audio) (ret []*AudioGroup, err error) {
	if !obj.Groups.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			qb := r.repository.Audio

			return obj.LoadGroups(ctx, qb)
		}); err != nil {
			return nil, err
		}
	}

	loader := loaders.From(ctx).GroupByID

	for _, sm := range obj.Groups.List() {
		group, err := loader.Load(sm.GroupID)
		if err != nil {
			return nil, err
		}

		audioIdx := sm.AudioIndex
		audioGroup := &AudioGroup{
			Group:      group,
			AudioIndex: audioIdx,
		}

		ret = append(ret, audioGroup)
	}

	return ret, nil
}

func (r *audioResolver) Tags(ctx context.Context, obj *models.Audio) (ret []*models.Tag, err error) {
	if !obj.TagIDs.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadTagIDs(ctx, r.repository.Audio)
		}); err != nil {
			return nil, err
		}
	}

	var errs []error
	ret, errs = loaders.From(ctx).TagByID.LoadAll(obj.TagIDs.List())
	return ret, firstError(errs)
}

func (r *audioResolver) Performers(ctx context.Context, obj *models.Audio) (ret []*models.Performer, err error) {
	if !obj.PerformerIDs.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadPerformerIDs(ctx, r.repository.Audio)
		}); err != nil {
			return nil, err
		}
	}

	var errs []error
	ret, errs = loaders.From(ctx).PerformerByID.LoadAll(obj.PerformerIDs.List())
	return ret, firstError(errs)
}

func (r *audioResolver) AudioStreams(ctx context.Context, obj *models.Audio) ([]*manager.AudioStreamEndpoint, error) {
	// load the primary file into the audio
	_, err := r.getPrimaryFile(ctx, obj)
	if err != nil {
		return nil, err
	}

	config := manager.GetInstance().Config

	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	builder := urlbuilders.NewAudioURLBuilder(baseURL, obj)
	apiKey := config.GetAPIKey()

	return manager.GetAudioStreamPaths(obj, builder.GetStreamURL(apiKey), config.GetMaxStreamingTranscodeSize())
}

func (r *audioResolver) Urls(ctx context.Context, obj *models.Audio) ([]string, error) {
	if !obj.URLs.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadURLs(ctx, r.repository.Audio)
		}); err != nil {
			return nil, err
		}
	}

	return obj.URLs.List(), nil
}

func (r *audioResolver) OCounter(ctx context.Context, obj *models.Audio) (*int, error) {
	ret, err := loaders.From(ctx).AudioOCount.Load(obj.ID)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func (r *audioResolver) LastPlayedAt(ctx context.Context, obj *models.Audio) (*time.Time, error) {
	ret, err := loaders.From(ctx).AudioLastPlayed.Load(obj.ID)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *audioResolver) PlayCount(ctx context.Context, obj *models.Audio) (*int, error) {
	ret, err := loaders.From(ctx).AudioPlayCount.Load(obj.ID)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func (r *audioResolver) PlayHistory(ctx context.Context, obj *models.Audio) ([]*time.Time, error) {
	ret, err := loaders.From(ctx).AudioPlayHistory.Load(obj.ID)
	if err != nil {
		return nil, err
	}

	// convert to pointer slice
	ptrRet := make([]*time.Time, len(ret))
	for i, t := range ret {
		tt := t
		ptrRet[i] = &tt
	}

	return ptrRet, nil
}

func (r *audioResolver) OHistory(ctx context.Context, obj *models.Audio) ([]*time.Time, error) {
	ret, err := loaders.From(ctx).AudioOHistory.Load(obj.ID)
	if err != nil {
		return nil, err
	}

	// convert to pointer slice
	ptrRet := make([]*time.Time, len(ret))
	for i, t := range ret {
		tt := t
		ptrRet[i] = &tt
	}

	return ptrRet, nil
}

func (r *audioResolver) CustomFields(ctx context.Context, obj *models.Audio) (map[string]interface{}, error) {
	m, err := loaders.From(ctx).AudioCustomFields.Load(obj.ID)
	if err != nil {
		return nil, err
	}

	if m == nil {
		return make(map[string]interface{}), nil
	}

	return m, nil
}
