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

	hash, err := videophash.Generate(instance.FFMpeg, t.File)
	if err != nil {
		logger.Errorf("error generating phash: %s", err.Error())
		logErrorOutput(err)
		return
	}

	r := t.repository
	if err := r.WithTxn(ctx, func(ctx context.Context) error {
		hashValue := int64(*hash)
		t.File.Fingerprints = t.File.Fingerprints.AppendUnique(models.Fingerprint{
			Type:        models.FingerprintTypePhash,
			Fingerprint: hashValue,
		})

		return r.File.Update(ctx, t.File)
	}); err != nil && ctx.Err() == nil {
		logger.Errorf("Error setting phash: %v", err)
	}
}

func (t *GeneratePhashTask) required() bool {
	if t.Overwrite {
		return true
	}

	return t.File.Fingerprints.Get(models.FingerprintTypePhash) == nil
}
