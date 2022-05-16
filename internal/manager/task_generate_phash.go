package manager

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/hash/videophash"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/txn"
)

type GeneratePhashTask struct {
	File                *file.VideoFile
	Overwrite           bool
	fileNamingAlgorithm models.HashAlgorithm
	txnManager          txn.Manager
	fileUpdater         file.Updater
}

func (t *GeneratePhashTask) GetDescription() string {
	return fmt.Sprintf("Generating phash for %s", t.File.Path)
}

func (t *GeneratePhashTask) Start(ctx context.Context) {
	if !t.shouldGenerate() {
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
		t.File.Fingerprints = t.File.Fingerprints.AppendUnique(file.Fingerprint{
			Type:        file.FingerprintTypePhash,
			Fingerprint: hashValue,
		})

		return qb.Update(ctx, t.File)
	}); err != nil {
		logger.Error(err.Error())
	}
}

func (t *GeneratePhashTask) shouldGenerate() bool {
	return t.Overwrite || t.File.Fingerprints.Get(file.FingerprintTypePhash) == nil
}
