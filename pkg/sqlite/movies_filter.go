package sqlite

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

func (qb *MovieStore) makeFilter(ctx context.Context, movieFilter *models.MovieFilterType) *filterBuilder {
	if movieFilter == nil {
		return nil
	}

	query := &filterBuilder{}

	query.handleCriterion(ctx, stringCriterionHandler(movieFilter.Name, "movies.name"))
	query.handleCriterion(ctx, stringCriterionHandler(movieFilter.Director, "movies.director"))
	query.handleCriterion(ctx, stringCriterionHandler(movieFilter.Synopsis, "movies.synopsis"))
	query.handleCriterion(ctx, intCriterionHandler(movieFilter.Rating100, "movies.rating", nil))
	query.handleCriterion(ctx, floatIntCriterionHandler(movieFilter.Duration, "movies.duration", nil))
	query.handleCriterion(ctx, qb.missingCriterionHandler(movieFilter.IsMissing))
	query.handleCriterion(ctx, stringCriterionHandler(movieFilter.URL, "movies.url"))
	query.handleCriterion(ctx, studioCriterionHandler(movieTable, movieFilter.Studios))
	query.handleCriterion(ctx, qb.performersCriterionHandler(movieFilter.Performers))
	query.handleCriterion(ctx, dateCriterionHandler(movieFilter.Date, "movies.date"))
	query.handleCriterion(ctx, timestampCriterionHandler(movieFilter.CreatedAt, "movies.created_at"))
	query.handleCriterion(ctx, timestampCriterionHandler(movieFilter.UpdatedAt, "movies.updated_at"))

	return query
}

func (qb *MovieStore) missingCriterionHandler(isMissing *string) criterionHandlerFunc {
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

func (qb *MovieStore) performersCriterionHandler(performers *models.MultiCriterionInput) criterionHandlerFunc {
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
