package sqlite

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

func (qb *StudioStore) validateFilter(filter *models.StudioFilterType) error {
	const and = "AND"
	const or = "OR"
	const not = "NOT"

	if filter.And != nil {
		if filter.Or != nil {
			return illegalFilterCombination(and, or)
		}
		if filter.Not != nil {
			return illegalFilterCombination(and, not)
		}

		return qb.validateFilter(filter.And)
	}

	if filter.Or != nil {
		if filter.Not != nil {
			return illegalFilterCombination(or, not)
		}

		return qb.validateFilter(filter.Or)
	}

	if filter.Not != nil {
		return qb.validateFilter(filter.Not)
	}

	return nil
}

func (qb *StudioStore) makeFilter(ctx context.Context, studioFilter *models.StudioFilterType) *filterBuilder {
	query := &filterBuilder{}

	if studioFilter.And != nil {
		query.and(qb.makeFilter(ctx, studioFilter.And))
	}
	if studioFilter.Or != nil {
		query.or(qb.makeFilter(ctx, studioFilter.Or))
	}
	if studioFilter.Not != nil {
		query.not(qb.makeFilter(ctx, studioFilter.Not))
	}

	query.handleCriterion(ctx, qb.criterionHandler(studioFilter))

	return query
}

func (qb *StudioStore) criterionHandler(studioFilter *models.StudioFilterType) criterionHandler {
	return compoundHandler{
		stringCriterionHandler(studioFilter.Name, studioTable+".name"),
		stringCriterionHandler(studioFilter.Details, studioTable+".details"),
		stringCriterionHandler(studioFilter.URL, studioTable+".url"),
		intCriterionHandler(studioFilter.Rating100, studioTable+".rating", nil),
		boolCriterionHandler(studioFilter.Favorite, studioTable+".favorite", nil),
		boolCriterionHandler(studioFilter.IgnoreAutoTag, studioTable+".ignore_auto_tag", nil),

		criterionHandlerFunc(func(ctx context.Context, f *filterBuilder) {
			if studioFilter.StashID != nil {
				qb.stashIDRepository().join(f, "studio_stash_ids", "studios.id")
				stringCriterionHandler(studioFilter.StashID, "studio_stash_ids.stash_id")(ctx, f)
			}
		}),
		&stashIDCriterionHandler{
			c:                 studioFilter.StashIDEndpoint,
			stashIDRepository: qb.stashIDRepository(),
			stashIDTableAs:    "studio_stash_ids",
			parentIDCol:       "studios.id",
		},

		qb.isMissingCriterionHandler(studioFilter.IsMissing),
		qb.sceneCountCriterionHandler(studioFilter.SceneCount),
		qb.imageCountCriterionHandler(studioFilter.ImageCount),
		qb.galleryCountCriterionHandler(studioFilter.GalleryCount),
		qb.parentCriterionHandler(studioFilter.Parents),
		qb.aliasCriterionHandler(studioFilter.Aliases),
		qb.childCountCriterionHandler(studioFilter.ChildCount),
		timestampCriterionHandler(studioFilter.CreatedAt, studioTable+".created_at"),
		timestampCriterionHandler(studioFilter.UpdatedAt, studioTable+".updated_at"),
	}
}

func (qb *StudioStore) isMissingCriterionHandler(isMissing *string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if isMissing != nil && *isMissing != "" {
			switch *isMissing {
			case "image":
				f.addWhere("studios.image_blob IS NULL")
			case "stash_id":
				qb.stashIDRepository().join(f, "studio_stash_ids", "studios.id")
				f.addWhere("studio_stash_ids.studio_id IS NULL")
			default:
				f.addWhere("(studios." + *isMissing + " IS NULL OR TRIM(studios." + *isMissing + ") = '')")
			}
		}
	}
}

func (qb *StudioStore) sceneCountCriterionHandler(sceneCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if sceneCount != nil {
			f.addLeftJoin("scenes", "", "scenes.studio_id = studios.id")
			clause, args := getIntCriterionWhereClause("count(distinct scenes.id)", *sceneCount)

			f.addHaving(clause, args...)
		}
	}
}

func (qb *StudioStore) imageCountCriterionHandler(imageCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if imageCount != nil {
			f.addLeftJoin("images", "", "images.studio_id = studios.id")
			clause, args := getIntCriterionWhereClause("count(distinct images.id)", *imageCount)

			f.addHaving(clause, args...)
		}
	}
}

func (qb *StudioStore) galleryCountCriterionHandler(galleryCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if galleryCount != nil {
			f.addLeftJoin("galleries", "", "galleries.studio_id = studios.id")
			clause, args := getIntCriterionWhereClause("count(distinct galleries.id)", *galleryCount)

			f.addHaving(clause, args...)
		}
	}
}

func (qb *StudioStore) parentCriterionHandler(parents *models.MultiCriterionInput) criterionHandlerFunc {
	addJoinsFunc := func(f *filterBuilder) {
		f.addLeftJoin("studios", "parent_studio", "parent_studio.id = studios.parent_id")
	}
	h := multiCriterionHandlerBuilder{
		primaryTable: studioTable,
		foreignTable: "parent_studio",
		joinTable:    "",
		primaryFK:    studioIDColumn,
		foreignFK:    "parent_id",
		addJoinsFunc: addJoinsFunc,
	}
	return h.handler(parents)
}

func (qb *StudioStore) aliasCriterionHandler(alias *models.StringCriterionInput) criterionHandlerFunc {
	h := stringListCriterionHandlerBuilder{
		joinTable:    studioAliasesTable,
		stringColumn: studioAliasColumn,
		addJoinTable: func(f *filterBuilder) {
			studiosAliasesTableMgr.join(f, "", "studios.id")
		},
	}

	return h.handler(alias)
}

func (qb *StudioStore) childCountCriterionHandler(childCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if childCount != nil {
			f.addLeftJoin("studios", "children_count", "children_count.parent_id = studios.id")
			clause, args := getIntCriterionWhereClause("count(distinct children_count.id)", *childCount)

			f.addHaving(clause, args...)
		}
	}
}
