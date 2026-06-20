package coverfetch

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/stashbox"
	"github.com/stashapp/stash/pkg/utils"
)

var (
	ErrNotConfigured  = errors.New("stash-box is not configured")
	ErrNoFingerprints = errors.New("scene has no fingerprints")
	ErrNoMatch        = errors.New("no stash-box match found")
	ErrNoImage        = errors.New("stash-box match has no image")
)

type Finder interface {
	FindSceneByFingerprints(context.Context, models.Fingerprints) ([]*models.ScrapedScene, error)
	QueryScene(context.Context, string) ([]*models.ScrapedScene, error)
}

type Result struct {
	SceneID      int
	RemoteSiteID string
	Bytes        int
}

type Service struct {
	finder Finder

	getFingerprints func(context.Context, int) (models.Fingerprints, error)
	getSceneInfo    func(context.Context, int) (sceneInfo, error)
	updateCover     func(context.Context, int, []byte) error
	processImage    func(context.Context, string) ([]byte, error)
}

type sceneInfo struct {
	Title    string
	Basename string
	Duration float64
}

func New(repo models.Repository, box models.StashBox) *Service {
	var finder Finder
	if box.Endpoint != "" {
		finder = stashbox.NewClient(box)
	}

	fingerprintService := scene.Service{Repository: repo.Scene}
	return &Service{
		finder: finder,
		getFingerprints: func(ctx context.Context, sceneID int) (models.Fingerprints, error) {
			var fingerprints []models.Fingerprints
			if err := repo.WithReadTxn(ctx, func(ctx context.Context) error {
				var err error
				fingerprints, err = fingerprintService.GetScenesFingerprints(ctx, []int{sceneID})
				return err
			}); err != nil {
				return nil, err
			}
			if len(fingerprints) == 0 {
				return nil, nil
			}
			return fingerprints[0], nil
		},
		getSceneInfo: func(ctx context.Context, sceneID int) (sceneInfo, error) {
			var ret sceneInfo
			if err := repo.WithReadTxn(ctx, func(ctx context.Context) error {
				scene, err := repo.Scene.Find(ctx, sceneID)
				if err != nil {
					return err
				}
				if scene == nil {
					return fmt.Errorf("scene with id %d not found", sceneID)
				}
				ret.Title = scene.Title
				if err := scene.LoadFiles(ctx, repo.Scene); err != nil {
					return err
				}
				if files := scene.Files.List(); len(files) > 0 {
					file := files[0]
					ret.Basename = file.Base().Basename
					ret.Duration = file.Duration
				}
				return nil
			}); err != nil {
				return sceneInfo{}, err
			}
			return ret, nil
		},
		updateCover: func(ctx context.Context, sceneID int, cover []byte) error {
			return repo.WithTxn(ctx, func(ctx context.Context) error {
				return repo.Scene.UpdateCover(ctx, sceneID, cover)
			})
		},
		processImage: defaultProcessImage,
	}
}

func (s *Service) Fetch(ctx context.Context, sceneID int) (Result, error) {
	if s.finder == nil {
		return Result{}, ErrNotConfigured
	}
	if s.getFingerprints == nil || s.updateCover == nil {
		return Result{}, ErrNotConfigured
	}
	if s.processImage == nil {
		s.processImage = defaultProcessImage
	}

	fingerprints, err := s.getFingerprints(ctx, sceneID)
	if err != nil {
		return Result{}, fmt.Errorf("load scene fingerprints: %w", err)
	}
	if len(fingerprints) == 0 {
		return Result{}, ErrNoFingerprints
	}

	scenes, err := s.finder.FindSceneByFingerprints(ctx, fingerprints)
	if err != nil {
		return Result{}, fmt.Errorf("query stash-box: %w", err)
	}

	match, err := firstSceneWithImage(scenes)
	if err != nil {
		match, err = s.findSceneByCode(ctx, sceneID)
		if err != nil {
			return Result{}, err
		}
	}

	cover, err := s.processImage(ctx, *match.Image)
	if err != nil {
		return Result{}, fmt.Errorf("process stash-box image: %w", err)
	}
	if len(cover) == 0 {
		return Result{}, ErrNoImage
	}

	if err := s.updateCover(ctx, sceneID, cover); err != nil {
		return Result{}, fmt.Errorf("update scene cover: %w", err)
	}

	result := Result{SceneID: sceneID, Bytes: len(cover)}
	if match.RemoteSiteID != nil {
		result.RemoteSiteID = *match.RemoteSiteID
	}
	return result, nil
}

func (s *Service) findSceneByCode(ctx context.Context, sceneID int) (*models.ScrapedScene, error) {
	if s.getSceneInfo == nil {
		return nil, ErrNoMatch
	}

	info, err := s.getSceneInfo(ctx, sceneID)
	if err != nil {
		return nil, fmt.Errorf("load scene info: %w", err)
	}

	for _, term := range sceneCodeSearchTerms(info) {
		scenes, err := s.finder.QueryScene(ctx, term)
		if err != nil {
			return nil, fmt.Errorf("query stash-box by code: %w", err)
		}
		if match := bestExactCodeMatch(scenes, term, info.Duration); match != nil {
			return match, nil
		}
	}

	return nil, ErrNoMatch
}

func firstSceneWithImage(scenes []*models.ScrapedScene) (*models.ScrapedScene, error) {
	if len(scenes) == 0 {
		return nil, ErrNoMatch
	}
	for _, scene := range scenes {
		if scene != nil && scene.Image != nil && *scene.Image != "" {
			return scene, nil
		}
	}
	return nil, ErrNoImage
}

var sceneCodePattern = regexp.MustCompile(`(?i)([a-z]{2,12})[-_ ]*(0*[0-9]{2,6})`)

func sceneCodeSearchTerms(info sceneInfo) []string {
	var ret []string
	add := func(term string) {
		term = strings.ToUpper(strings.TrimSpace(term))
		term = strings.TrimSuffix(term, strings.ToUpper(filepath.Ext(term)))
		if term != "" && !slices.Contains(ret, term) {
			ret = append(ret, term)
		}
	}

	for _, source := range []string{info.Title, info.Basename} {
		stem := strings.TrimSuffix(source, filepath.Ext(source))
		matches := sceneCodePattern.FindStringSubmatch(stem)
		if len(matches) == 0 {
			continue
		}
		prefix := strings.ToUpper(matches[1])
		digits := matches[2]
		add(prefix + "-" + digits)
		trimmed := strings.TrimLeft(digits, "0")
		if trimmed != "" {
			add(prefix + "-" + trimmed)
		}
	}

	return ret
}

func bestExactCodeMatch(scenes []*models.ScrapedScene, term string, duration float64) *models.ScrapedScene {
	var best *models.ScrapedScene
	bestDelta := 0.0

	for _, scene := range scenes {
		if scene == nil || scene.Image == nil || *scene.Image == "" {
			continue
		}
		if !sceneHasExactCode(scene, term) || !sceneDurationMatches(scene, duration) {
			continue
		}
		delta := sceneDurationDelta(scene, duration)
		if best == nil || delta < bestDelta {
			best = scene
			bestDelta = delta
		}
	}

	return best
}

func sceneHasExactCode(scene *models.ScrapedScene, term string) bool {
	term = strings.ToUpper(strings.TrimSpace(term))
	for _, value := range []*string{scene.Code, scene.Title} {
		if value != nil && strings.EqualFold(strings.TrimSpace(*value), term) {
			return true
		}
	}
	return false
}

func sceneDurationMatches(scene *models.ScrapedScene, duration float64) bool {
	if duration <= 0 || scene.Duration == nil || *scene.Duration <= 0 {
		return true
	}
	return sceneDurationDelta(scene, duration) <= max(180, duration*0.05)
}

func sceneDurationDelta(scene *models.ScrapedScene, duration float64) float64 {
	if scene.Duration == nil {
		return 0
	}
	delta := float64(*scene.Duration) - duration
	if delta < 0 {
		return -delta
	}
	return delta
}

func defaultProcessImage(ctx context.Context, image string) ([]byte, error) {
	return utils.ProcessImageInput(ctx, image)
}
