package manager

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/hash/imagephash"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type GenerateImagePhashTask struct {
	repository models.Repository
	File       *models.ImageFile
	Overwrite  bool
}

func (t *GenerateImagePhashTask) GetDescription() string {
	return fmt.Sprintf("Generating phash for %s", t.File.Path)
}

func (t *GenerateImagePhashTask) Start(ctx context.Context) {
	if !t.required() {
		return
	}

	var hash int64
	set := false

	// #4393 - if there is a file with the same md5, we can use the same phash
	// only use this if we're not overwriting
	if !t.Overwrite {
		existing, err := t.findExistingPhash(ctx)
		if err != nil {
			logger.Warnf("Error finding existing phash: %v", err)
		} else if existing != nil {
			logger.Infof("Using existing phash for %s", t.File.Path)
			hash = existing.(int64)
			set = true
		}
	}

	if !set {
		generated, err := imagephash.Generate(t.File)
		if err != nil {
			logger.Errorf("Error generating phash for %q: %v", t.File.Path, err)
			logErrorOutput(err)
			return
		}

		hash = int64(*generated)
	}

	r := t.repository
	if err := r.WithTxn(ctx, func(ctx context.Context) error {
		t.File.Fingerprints = t.File.Fingerprints.AppendUnique(models.Fingerprint{
			Type:        models.FingerprintTypePhash,
			Fingerprint: hash,
		})

		return r.File.Update(ctx, t.File)
	}); err != nil && ctx.Err() == nil {
		logger.Errorf("Error setting phash: %v", err)
	}
}

func (t *GenerateImagePhashTask) findExistingPhash(ctx context.Context) (interface{}, error) {
	r := t.repository
	var ret interface{}
	if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
		md5 := t.File.Fingerprints.Get(models.FingerprintTypeMD5)

		// find other files with the same md5
		files, err := r.File.FindByFingerprint(ctx, models.Fingerprint{
			Type:        models.FingerprintTypeMD5,
			Fingerprint: md5,
		})
		if err != nil {
			return fmt.Errorf("finding files by md5: %w", err)
		}

		// find the first file with a phash
		for _, file := range files {
			if phash := file.Base().Fingerprints.Get(models.FingerprintTypePhash); phash != nil {
				ret = phash
				return nil
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (t *GenerateImagePhashTask) required() bool {
	if t.Overwrite {
		return true
	}

	return t.File.Fingerprints.Get(models.FingerprintTypePhash) == nil
}
