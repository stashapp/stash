package browse

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/stashapp/stash/pkg/models"
)

type Query struct {
	Text      string
	Tag       string
	Performer string
	Page      int
	PerPage   int
	Sort      string
	Direction models.SortDirectionEnum
}

type SceneItem struct {
	ID        int
	Title     string
	Path      string
	Duration  float64
	Date      string
	Rating    *int
	Organized bool
}

type Result struct {
	Items []SceneItem
	Total int
}

type Service struct {
	repo models.Repository
}

func New(repo models.Repository) *Service {
	return &Service{repo: repo}
}

func ParseQuery(raw string) Query {
	var query Query
	for _, part := range strings.Fields(raw) {
		key, value, ok := strings.Cut(part, ":")
		if !ok {
			query.Text = strings.TrimSpace(strings.Join([]string{query.Text, part}, " "))
			continue
		}

		switch strings.ToLower(key) {
		case "tag":
			query.Tag = value
		case "performer", "actress":
			query.Performer = value
		default:
			query.Text = strings.TrimSpace(strings.Join([]string{query.Text, part}, " "))
		}
	}

	query.Text = strings.TrimSpace(query.Text)
	return query
}

func (s *Service) Search(ctx context.Context, query Query) (Result, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.PerPage == 0 {
		query.PerPage = 40
	}
	if query.Sort == "" {
		query.Sort = "title"
	}
	if query.Direction == "" {
		query.Direction = models.SortDirectionEnumAsc
	}

	var ret Result
	err := s.repo.WithReadTxn(ctx, func(ctx context.Context) error {
		sceneFilter, err := s.buildSceneFilter(ctx, query)
		if err != nil {
			return err
		}

		findFilter := &models.FindFilterType{
			Page:      &query.Page,
			PerPage:   &query.PerPage,
			Sort:      &query.Sort,
			Direction: &query.Direction,
		}
		if query.Text != "" {
			findFilter.Q = &query.Text
		}

		result, err := s.repo.Scene.Query(ctx, models.SceneQueryOptions{
			QueryOptions: models.QueryOptions{
				FindFilter: findFilter,
				Count:      true,
			},
			SceneFilter: sceneFilter,
		})
		if err != nil {
			return err
		}

		scenes, err := result.Resolve(ctx)
		if err != nil {
			return err
		}

		ret.Total = result.Count
		ret.Items = make([]SceneItem, 0, len(scenes))
		for _, scene := range scenes {
			item, err := s.sceneItem(ctx, scene)
			if err != nil {
				return err
			}
			ret.Items = append(ret.Items, item)
		}

		return nil
	})

	return ret, err
}

func (s *Service) buildSceneFilter(ctx context.Context, query Query) (*models.SceneFilterType, error) {
	filter := &models.SceneFilterType{}

	if query.Tag != "" {
		tag, err := s.repo.Tag.FindByName(ctx, query.Tag, true)
		if err != nil {
			return nil, fmt.Errorf("find tag %q: %w", query.Tag, err)
		}
		if tag == nil {
			filter.Tags = &models.HierarchicalMultiCriterionInput{
				Value:    []string{"-1"},
				Modifier: models.CriterionModifierIncludes,
			}
		} else {
			filter.Tags = &models.HierarchicalMultiCriterionInput{
				Value:    []string{strconv.Itoa(tag.ID)},
				Modifier: models.CriterionModifierIncludes,
			}
		}
	}

	if query.Performer != "" {
		performers, err := s.repo.Performer.FindByNames(ctx, []string{query.Performer}, true)
		if err != nil {
			return nil, fmt.Errorf("find performer %q: %w", query.Performer, err)
		}

		values := make([]string, 0, len(performers))
		for _, performer := range performers {
			values = append(values, strconv.Itoa(performer.ID))
		}
		if len(values) == 0 {
			values = []string{"-1"}
		}

		filter.Performers = &models.MultiCriterionInput{
			Value:    values,
			Modifier: models.CriterionModifierIncludes,
		}
	}

	if filter.Tags == nil && filter.Performers == nil {
		return nil, nil
	}

	return filter, nil
}

func (s *Service) sceneItem(ctx context.Context, scene *models.Scene) (SceneItem, error) {
	item := SceneItem{
		ID:        scene.ID,
		Title:     scene.GetTitle(),
		Path:      scene.Path,
		Rating:    scene.Rating,
		Organized: scene.Organized,
	}
	if scene.Date != nil {
		item.Date = scene.Date.String()
	}

	if err := scene.LoadPrimaryFile(ctx, s.repo.File); err != nil {
		return SceneItem{}, err
	}
	if primary := scene.Files.Primary(); primary != nil {
		item.Path = primary.Path
		item.Duration = primary.Duration
	}

	return item, nil
}
