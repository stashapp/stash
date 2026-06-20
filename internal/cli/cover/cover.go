package cover

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/stashapp/stash/pkg/models"
)

type Source string

const (
	SourceMissing  Source = "missing"
	SourceDatabase Source = "database"
)

type Request struct {
	SceneID  int
	Path     string
	Duration float64
}

type Cover struct {
	SceneID int
	Data    []byte
	Path    string
	Source  Source
}

type Service struct {
	repo     models.Repository
	hasRepo  bool
	cacheDir string
}

func New(repo models.Repository, cacheDir string) *Service {
	return &Service{repo: repo, hasRepo: true, cacheDir: cacheDir}
}

func (s *Service) Load(ctx context.Context, req Request) (Cover, error) {
	if s.hasRepo {
		var data []byte
		if err := s.repo.WithReadTxn(ctx, func(ctx context.Context) error {
			var err error
			data, err = s.repo.Scene.GetCover(ctx, req.SceneID)
			return err
		}); err != nil {
			return Cover{}, fmt.Errorf("read scene cover: %w", err)
		}

		if len(data) > 0 {
			return Cover{SceneID: req.SceneID, Data: data, Source: SourceDatabase}, nil
		}
	}

	return Cover{SceneID: req.SceneID, Source: SourceMissing}, nil
}

func (s *Service) WriteCache(sceneID int, data []byte) (string, error) {
	if len(data) == 0 {
		return "", fmt.Errorf("cover data is empty")
	}

	path := s.CachePath(sceneID)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return "", err
	}

	return path, nil
}

func (s *Service) CachePath(sceneID int) string {
	sum := sha256.Sum256([]byte(fmt.Sprintf("scene:%d", sceneID)))
	return filepath.Join(s.cacheDir, "covers", hex.EncodeToString(sum[:])+".jpg")
}
