package manager

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/hash/videophash"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/txn"
)

type GeneratePhashTask struct {
	File                *models.VideoFile
	Overwrite           bool
	fileNamingAlgorithm models.HashAlgorithm
	txnManager          txn.Manager
	fileUpdater         models.FileUpdater
}

func (t *GeneratePhashTask) GetDescription() string {
	return fmt.Sprintf("Generating phash for %s", t.File.Path)
}

func (t *GeneratePhashTask) Start(ctx context.Context) {
	if !t.required() {
		return
	}

	hash, err := videophash.Generate(instance.FFMPEG, t.File)
	if err != nil {
		logger.Errorf("error generating phash: %s", err.Error())
		logErrorOutput(err)
		return
	}

	if err := txn.WithTxn(ctx, t.txnManager, func(ctx context.Context) error {
		qb := t.fileUpdater
		hashValue := int64(*hash)
		t.File.Fingerprints = t.File.Fingerprints.AppendUnique(models.Fingerprint{
			Type:        models.FingerprintTypePhash,
			Fingerprint: hashValue,
		})

		return qb.Update(ctx, t.File)
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
