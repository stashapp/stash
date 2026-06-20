package covergen

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/ffmpeg/transcoder"
	"github.com/stashapp/stash/pkg/models"
)

type Request struct {
	SceneID  int
	Path     string
	Duration float64
}

type Result struct {
	SceneID int
	Bytes   int
}

type Service struct {
	repo    models.Repository
	encoder *ffmpeg.FFMpeg
}

func New(repo models.Repository, ffmpegPath string) *Service {
	if ffmpegPath == "" {
		return nil
	}
	return &Service{repo: repo, encoder: ffmpeg.NewEncoder(ffmpegPath)}
}

func (s *Service) Generate(ctx context.Context, req Request) (Result, error) {
	if s == nil || s.encoder == nil {
		return Result{}, fmt.Errorf("ffmpeg cover generator is not configured")
	}
	if req.Path == "" {
		return Result{}, fmt.Errorf("scene path is empty")
	}

	tmp, err := os.CreateTemp("", "stash-cli-cover-*.jpg")
	if err != nil {
		return Result{}, err
	}
	tmpPath := tmp.Name()
	_ = tmp.Close()
	defer os.Remove(tmpPath)

	at := req.Duration * 0.2
	if at <= 0 {
		at = 1
	}

	args := transcoder.ScreenshotTime(req.Path, at, transcoder.ScreenshotOptions{
		OutputPath: tmpPath,
		OutputType: transcoder.ScreenshotOutputTypeImage2,
		Width:      1280,
		Quality:    3,
	})
	if err := s.encoder.Generate(ctx, args); err != nil {
		return Result{}, err
	}

	data, err := os.ReadFile(filepath.Clean(tmpPath))
	if err != nil {
		return Result{}, err
	}
	if len(data) == 0 {
		return Result{}, fmt.Errorf("generated cover is empty")
	}

	if err := s.repo.WithTxn(ctx, func(ctx context.Context) error {
		return s.repo.Scene.UpdateCover(ctx, req.SceneID, data)
	}); err != nil {
		return Result{}, err
	}

	return Result{SceneID: req.SceneID, Bytes: len(data)}, nil
}
