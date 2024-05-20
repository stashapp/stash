package manager

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/hash/videophash"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type GeneratePhashTask struct {
	repository          models.Repository
	File                *models.VideoFile
	Overwrite           bool
	fileNamingAlgorithm models.HashAlgorithm
}

func (t *GeneratePhashTask) GetDescription() string {
	return fmt.Sprintf("Generating phash for %s", t.File.Path)
}

func (t *GeneratePhashTask) Start(ctx context.Context) {
	if !t.required() {
		return
	}

	var hash int64
	set := false

	// #4393 - if there is a file with the same oshash, we can use the same phash
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
		generated, err := videophash.Generate(instance.FFMpeg, t.File)
		if err != nil {
			logger.Errorf("Error generating phash: %v", err)
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

func (t *GeneratePhashTask) findExistingPhash(ctx context.Context) (interface{}, error) {
	r := t.repository
	var ret interface{}
	if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
		oshash := t.File.Fingerprints.Get(models.FingerprintTypeOshash)

		// find other files with the same oshash
		files, err := r.File.FindByFingerprint(ctx, models.Fingerprint{
			Type:        models.FingerprintTypeOshash,
			Fingerprint: oshash,
		})
		if err != nil {
			return fmt.Errorf("finding files by oshash: %w", err)
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

func (t *GeneratePhashTask) required() bool {
	if t.Overwrite {
		return true
	}

	return t.File.Fingerprints.Get(models.FingerprintTypePhash) == nil
}
