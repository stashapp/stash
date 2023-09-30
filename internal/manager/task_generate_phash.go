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
	hashValue := int64(*hash)
	t.File.Fingerprints = t.File.Fingerprints.AppendUnique(models.Fingerprint{
		Type:        models.FingerprintTypePhash,
		Fingerprint: hashValue,
	})

	// For a temporary time, if it doesn't already exist, also generate old style
	// phashes for videos shorter than 2.5 min to aid in tagging with Stashbox
	if t.File.Fingerprints.Get(models.FingerprintTypePhashOld) == nil {
		hashOld, err := videophash.Generate(instance.FFMPEG, t.File, true)
		if err != nil {
			logger.Errorf("error generating phash-old: %s", err.Error())
			logErrorOutput(err)
			return
		}
		hashOldValue := int64(*hashOld)
		t.File.Fingerprints = t.File.Fingerprints.AppendUnique(models.Fingerprint{
			Type:        models.FingerprintTypePhashOld,
			Fingerprint: hashOldValue,
		})
	}

	if err := txn.WithTxn(ctx, t.txnManager, func(ctx context.Context) error {
		qb := t.fileUpdater
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
