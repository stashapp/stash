package edit

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/models"
)

var ErrUnsupportedFavorite = errors.New("scene does not support the favorite field")

type Update struct {
	Title     *string
	Date      *string
	Rating    *int
	Organized *bool
	Watched   *bool
	Favorite  *bool
}

type Service struct {
	repo models.Repository
}

func New(repo models.Repository) *Service {
	return &Service{repo: repo}
}

func ParseArgs(args []string) (Update, error) {
	var update Update
	for _, arg := range args {
		key, value, ok := strings.Cut(arg, "=")
		if !ok {
			return Update{}, fmt.Errorf("argument %q must be key=value", arg)
		}

		switch strings.ToLower(key) {
		case "title":
			update.Title = &value
		case "date":
			update.Date = &value
		case "rating":
			rating, err := strconv.Atoi(value)
			if err != nil {
				return Update{}, fmt.Errorf("rating must be a number: %w", err)
			}
			update.Rating = &rating
		case "organized":
			organized, err := strconv.ParseBool(value)
			if err != nil {
				return Update{}, fmt.Errorf("organized must be true/false: %w", err)
			}
			update.Organized = &organized
		case "watched":
			watched, err := strconv.ParseBool(value)
			if err != nil {
				return Update{}, fmt.Errorf("watched must be true/false: %w", err)
			}
			update.Watched = &watched
		case "favorite":
			favorite, err := strconv.ParseBool(value)
			if err != nil {
				return Update{}, fmt.Errorf("favorite must be true/false: %w", err)
			}
			update.Favorite = &favorite
		default:
			return Update{}, fmt.Errorf("unknown edit field %q", key)
		}
	}

	return update, nil
}

func (s *Service) Apply(ctx context.Context, sceneID int, update Update) error {
	partial, hasSceneFields, err := BuildPartial(update)
	if err != nil {
		return err
	}

	return s.repo.WithTxn(ctx, func(ctx context.Context) error {
		if hasSceneFields {
			if _, err := s.repo.Scene.UpdatePartial(ctx, sceneID, partial); err != nil {
				return fmt.Errorf("update scene: %w", err)
			}
		}

		if update.Watched != nil {
			if *update.Watched {
				if _, err := s.repo.Scene.AddViews(ctx, sceneID, []time.Time{time.Now()}); err != nil {
					return fmt.Errorf("mark watched: %w", err)
				}
			} else {
				if _, err := s.repo.Scene.DeleteAllViews(ctx, sceneID); err != nil {
					return fmt.Errorf("clear watched: %w", err)
				}
			}
		}

		return nil
	})
}

func BuildPartial(update Update) (models.ScenePartial, bool, error) {
	if update.Favorite != nil {
		return models.ScenePartial{}, false, ErrUnsupportedFavorite
	}

	partial := models.NewScenePartial()
	hasSceneFields := false

	if update.Title != nil {
		partial.Title = models.NewOptionalString(*update.Title)
		hasSceneFields = true
	}
	if update.Date != nil {
		date, err := models.ParseDate(*update.Date)
		if err != nil {
			return models.ScenePartial{}, false, err
		}
		partial.Date = models.NewOptionalDate(date)
		hasSceneFields = true
	}
	if update.Rating != nil {
		if *update.Rating < 0 || *update.Rating > 100 {
			return models.ScenePartial{}, false, fmt.Errorf("rating must be between 0 and 100")
		}
		partial.Rating = models.NewOptionalInt(*update.Rating)
		hasSceneFields = true
	}
	if update.Organized != nil {
		partial.Organized = models.NewOptionalBool(*update.Organized)
		hasSceneFields = true
	}

	return partial, hasSceneFields, nil
}
