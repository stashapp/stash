package sqlite

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

type movieFilterHandler struct {
	movieFilter *models.MovieFilterType
}

func (qb *movieFilterHandler) validate() error {
	movieFilter := qb.movieFilter
	if movieFilter == nil {
		return nil
	}

	if err := validateFilterCombination(movieFilter.OperatorFilter); err != nil {
		return err
	}

	if subFilter := movieFilter.SubFilter(); subFilter != nil {
		sqb := &movieFilterHandler{movieFilter: subFilter}
		if err := sqb.validate(); err != nil {
			return err
		}
	}

	return nil
}

func (qb *movieFilterHandler) handle(ctx context.Context, f *filterBuilder) {
	movieFilter := qb.movieFilter
	if movieFilter == nil {
		return
	}

	if err := qb.validate(); err != nil {
		f.setError(err)
		return
	}

	sf := movieFilter.SubFilter()
	if sf != nil {
		sub := &movieFilterHandler{sf}
		handleSubFilter(ctx, sub, f, movieFilter.OperatorFilter)
	}

	f.handleCriterion(ctx, qb.criterionHandler())
}

func (qb *movieFilterHandler) criterionHandler() criterionHandler {
	movieFilter := qb.movieFilter
	return compoundHandler{
		stringCriterionHandler(movieFilter.Name, "movies.name"),
		stringCriterionHandler(movieFilter.Director, "movies.director"),
		stringCriterionHandler(movieFilter.Synopsis, "movies.synopsis"),
		intCriterionHandler(movieFilter.Rating100, "movies.rating", nil),
		floatIntCriterionHandler(movieFilter.Duration, "movies.duration", nil),
		qb.missingCriterionHandler(movieFilter.IsMissing),
		qb.urlsCriterionHandler(movieFilter.URL),
		studioCriterionHandler(movieTable, movieFilter.Studios),
		qb.performersCriterionHandler(movieFilter.Performers),
		&dateCriterionHandler{movieFilter.Date, "movies.date", nil},
		&timestampCriterionHandler{movieFilter.CreatedAt, "movies.created_at", nil},
		&timestampCriterionHandler{movieFilter.UpdatedAt, "movies.updated_at", nil},

		&relatedFilterHandler{
			relatedIDCol:   "movies_scenes.scene_id",
			relatedRepo:    sceneRepository.repository,
			relatedHandler: &sceneFilterHandler{movieFilter.ScenesFilter},
			joinFn: func(f *filterBuilder) {
				movieRepository.scenes.innerJoin(f, "", "movies.id")
			},
		},

		&relatedFilterHandler{
			relatedIDCol:   "movies.studio_id",
			relatedRepo:    studioRepository.repository,
			relatedHandler: &studioFilterHandler{movieFilter.StudiosFilter},
		},
	}
}

func (qb *movieFilterHandler) missingCriterionHandler(isMissing *string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if isMissing != nil && *isMissing != "" {
			switch *isMissing {
			case "front_image":
				f.addWhere("movies.front_image_blob IS NULL")
			case "back_image":
				f.addWhere("movies.back_image_blob IS NULL")
			case "scenes":
				f.addLeftJoin("movies_scenes", "", "movies_scenes.movie_id = movies.id")
				f.addWhere("movies_scenes.scene_id IS NULL")
			default:
				f.addWhere("(movies." + *isMissing + " IS NULL OR TRIM(movies." + *isMissing + ") = '')")
			}
		}
	}
}

func (qb *movieFilterHandler) urlsCriterionHandler(url *models.StringCriterionInput) criterionHandlerFunc {
	h := stringListCriterionHandlerBuilder{
		primaryTable: movieTable,
		primaryFK:    movieIDColumn,
		joinTable:    movieURLsTable,
		stringColumn: movieURLColumn,
		addJoinTable: func(f *filterBuilder) {
			moviesURLsTableMgr.join(f, "", "movies.id")
		},
	}

	return h.handler(url)
}

func (qb *movieFilterHandler) performersCriterionHandler(performers *models.MultiCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if performers != nil {
			if performers.Modifier == models.CriterionModifierIsNull || performers.Modifier == models.CriterionModifierNotNull {
				var notClause string
				if performers.Modifier == models.CriterionModifierNotNull {
					notClause = "NOT"
				}

				f.addLeftJoin("movies_scenes", "", "movies.id = movies_scenes.movie_id")
				f.addLeftJoin("performers_scenes", "", "movies_scenes.scene_id = performers_scenes.scene_id")

				f.addWhere(fmt.Sprintf("performers_scenes.performer_id IS %s NULL", notClause))
				return
			}

			if len(performers.Value) == 0 {
				return
			}

			var args []interface{}
			for _, arg := range performers.Value {
				args = append(args, arg)
			}

			// Hack, can't apply args to join, nor inner join on a left join, so use CTE instead
			f.addWith(`movies_performers AS (
				SELECT movies_scenes.movie_id, performers_scenes.performer_id
				FROM movies_scenes
				INNER JOIN performers_scenes ON movies_scenes.scene_id = performers_scenes.scene_id
				WHERE performers_scenes.performer_id IN`+getInBinding(len(performers.Value))+`
			)`, args...)
			f.addLeftJoin("movies_performers", "", "movies.id = movies_performers.movie_id")

			switch performers.Modifier {
			case models.CriterionModifierIncludes:
				f.addWhere("movies_performers.performer_id IS NOT NULL")
			case models.CriterionModifierIncludesAll:
				f.addWhere("movies_performers.performer_id IS NOT NULL")
				f.addHaving("COUNT(DISTINCT movies_performers.performer_id) = ?", len(performers.Value))
			case models.CriterionModifierExcludes:
				f.addWhere("movies_performers.performer_id IS NULL")
			}
		}
	}
}
