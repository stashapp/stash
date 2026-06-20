package player

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/stashapp/stash/internal/cli/browse"
	"github.com/stashapp/stash/pkg/models"
)

type Runner interface {
	Run(context.Context, string, []string) error
}

type ViewRecorder interface {
	AddView(context.Context, int, time.Time) error
}

type Service struct {
	path     string
	args     []string
	runner   Runner
	recorder ViewRecorder
}

func New(repo models.Repository, path string, args []string) *Service {
	return NewWithDeps(path, args, commandRunner{}, repoRecorder{repo: repo})
}

func NewWithDeps(path string, args []string, runner Runner, recorder ViewRecorder) *Service {
	return &Service{
		path:     strings.TrimSpace(path),
		args:     append([]string(nil), args...),
		runner:   runner,
		recorder: recorder,
	}
}

func (s *Service) Play(ctx context.Context, item browse.SceneItem) error {
	if item.Path == "" {
		return errors.New("selected scene has no playable file path")
	}
	if s.path == "" {
		return errors.New("ffplay is unavailable: configure ffplay_path")
	}
	if s.runner == nil {
		return errors.New("ffplay runner is not configured")
	}

	args := append([]string(nil), s.args...)
	args = append(args, item.Path)
	if err := s.runner.Run(ctx, s.path, args); err != nil {
		return fmt.Errorf("play scene: %w", err)
	}
	if s.recorder != nil {
		if err := s.recorder.AddView(ctx, item.ID, time.Now()); err != nil {
			return fmt.Errorf("record scene view: %w", err)
		}
	}

	return nil
}

type commandRunner struct{}

func (commandRunner) Run(ctx context.Context, path string, args []string) error {
	cmd := exec.CommandContext(ctx, path, args...)
	return cmd.Run()
}

type repoRecorder struct {
	repo models.Repository
}

func (r repoRecorder) AddView(ctx context.Context, sceneID int, at time.Time) error {
	return r.repo.WithTxn(ctx, func(ctx context.Context) error {
		_, err := r.repo.Scene.AddViews(ctx, sceneID, []time.Time{at})
		return err
	})
}
