package sqlite

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

type groupFilterHandler struct {
	groupFilter *models.GroupFilterType
}

func (qb *groupFilterHandler) validate() error {
	groupFilter := qb.groupFilter
	if groupFilter == nil {
		return nil
	}

	if err := validateFilterCombination(groupFilter.OperatorFilter); err != nil {
		return err
	}

	if subFilter := groupFilter.SubFilter(); subFilter != nil {
		sqb := &groupFilterHandler{groupFilter: subFilter}
		if err := sqb.validate(); err != nil {
			return err
		}
	}

	return nil
}

func (qb *groupFilterHandler) handle(ctx context.Context, f *filterBuilder) {
	groupFilter := qb.groupFilter
	if groupFilter == nil {
		return
	}

	if err := qb.validate(); err != nil {
		f.setError(err)
		return
	}

	sf := groupFilter.SubFilter()
	if sf != nil {
		sub := &groupFilterHandler{sf}
		handleSubFilter(ctx, sub, f, groupFilter.OperatorFilter)
	}

	f.handleCriterion(ctx, qb.criterionHandler())
}

func (qb *groupFilterHandler) criterionHandler() criterionHandler {
	groupFilter := qb.groupFilter
	return compoundHandler{
		stringCriterionHandler(groupFilter.Name, "movies.name"),
		stringCriterionHandler(groupFilter.Director, "movies.director"),
		stringCriterionHandler(groupFilter.Synopsis, "movies.synopsis"),
		intCriterionHandler(groupFilter.Rating100, "movies.rating", nil),
		floatIntCriterionHandler(groupFilter.Duration, "movies.duration", nil),
		qb.missingCriterionHandler(groupFilter.IsMissing),
		qb.urlsCriterionHandler(groupFilter.URL),
		studioCriterionHandler(groupTable, groupFilter.Studios),
		qb.performersCriterionHandler(groupFilter.Performers),
		qb.tagsCriterionHandler(groupFilter.Tags),
		qb.tagCountCriterionHandler(groupFilter.TagCount),
		&dateCriterionHandler{groupFilter.Date, "movies.date", nil},
		&timestampCriterionHandler{groupFilter.CreatedAt, "movies.created_at", nil},
		&timestampCriterionHandler{groupFilter.UpdatedAt, "movies.updated_at", nil},

		&relatedFilterHandler{
			relatedIDCol:   "movies_scenes.scene_id",
			relatedRepo:    sceneRepository.repository,
			relatedHandler: &sceneFilterHandler{groupFilter.ScenesFilter},
			joinFn: func(f *filterBuilder) {
				groupRepository.scenes.innerJoin(f, "", "movies.id")
			},
		},

		&relatedFilterHandler{
			relatedIDCol:   "movies.studio_id",
			relatedRepo:    studioRepository.repository,
			relatedHandler: &studioFilterHandler{groupFilter.StudiosFilter},
		},
	}
}

func (qb *groupFilterHandler) missingCriterionHandler(isMissing *string) criterionHandlerFunc {
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

func (qb *groupFilterHandler) urlsCriterionHandler(url *models.StringCriterionInput) criterionHandlerFunc {
	h := stringListCriterionHandlerBuilder{
		primaryTable: groupTable,
		primaryFK:    groupIDColumn,
		joinTable:    groupURLsTable,
		stringColumn: groupURLColumn,
		addJoinTable: func(f *filterBuilder) {
			groupsURLsTableMgr.join(f, "", "movies.id")
		},
	}

	return h.handler(url)
}

func (qb *groupFilterHandler) performersCriterionHandler(performers *models.MultiCriterionInput) criterionHandlerFunc {
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

func (qb *groupFilterHandler) tagsCriterionHandler(tags *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	h := joinedHierarchicalMultiCriterionHandlerBuilder{
		primaryTable: groupTable,
		foreignTable: tagTable,
		foreignFK:    "tag_id",

		relationsTable: "tags_relations",
		joinAs:         "movie_tag",
		joinTable:      groupsTagsTable,
		primaryFK:      groupIDColumn,
	}

	return h.handler(tags)
}

func (qb *groupFilterHandler) tagCountCriterionHandler(count *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: groupTable,
		joinTable:    groupsTagsTable,
		primaryFK:    groupIDColumn,
	}

	return h.handler(count)
}
