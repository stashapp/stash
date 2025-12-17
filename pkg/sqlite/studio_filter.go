package sqlite

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

type studioFilterHandler struct {
	studioFilter *models.StudioFilterType
}

func (qb *studioFilterHandler) validate() error {
	studioFilter := qb.studioFilter
	if studioFilter == nil {
		return nil
	}

	if err := validateFilterCombination(studioFilter.OperatorFilter); err != nil {
		return err
	}

	if subFilter := studioFilter.SubFilter(); subFilter != nil {
		sqb := &studioFilterHandler{studioFilter: subFilter}
		if err := sqb.validate(); err != nil {
			return err
		}
	}

	return nil
}

func (qb *studioFilterHandler) handle(ctx context.Context, f *filterBuilder) {
	studioFilter := qb.studioFilter
	if studioFilter == nil {
		return
	}

	if err := qb.validate(); err != nil {
		f.setError(err)
		return
	}

	sf := studioFilter.SubFilter()
	if sf != nil {
		sub := &studioFilterHandler{sf}
		handleSubFilter(ctx, sub, f, studioFilter.OperatorFilter)
	}

	f.handleCriterion(ctx, qb.criterionHandler())
}

func (qb *studioFilterHandler) criterionHandler() criterionHandler {
	studioFilter := qb.studioFilter
	return compoundHandler{
		stringCriterionHandler(studioFilter.Name, studioTable+".name"),
		stringCriterionHandler(studioFilter.Details, studioTable+".details"),
		qb.urlsCriterionHandler(studioFilter.URL),
		intCriterionHandler(studioFilter.Rating100, studioTable+".rating", nil),
		boolCriterionHandler(studioFilter.Favorite, studioTable+".favorite", nil),
		boolCriterionHandler(studioFilter.IgnoreAutoTag, studioTable+".ignore_auto_tag", nil),

		criterionHandlerFunc(func(ctx context.Context, f *filterBuilder) {
			if studioFilter.StashID != nil {
				studioRepository.stashIDs.join(f, "studio_stash_ids", "studios.id")
				stringCriterionHandler(studioFilter.StashID, "studio_stash_ids.stash_id")(ctx, f)
			}
		}),
		&stashIDCriterionHandler{
			c:                 studioFilter.StashIDEndpoint,
			stashIDRepository: &studioRepository.stashIDs,
			stashIDTableAs:    "studio_stash_ids",
			parentIDCol:       "studios.id",
		},

		qb.isMissingCriterionHandler(studioFilter.IsMissing),
		qb.tagCountCriterionHandler(studioFilter.TagCount),
		qb.sceneCountCriterionHandler(studioFilter.SceneCount),
		qb.imageCountCriterionHandler(studioFilter.ImageCount),
		qb.galleryCountCriterionHandler(studioFilter.GalleryCount),
		qb.parentCriterionHandler(studioFilter.Parents),
		qb.aliasCriterionHandler(studioFilter.Aliases),
		qb.tagsCriterionHandler(studioFilter.Tags),
		qb.childCountCriterionHandler(studioFilter.ChildCount),
		&timestampCriterionHandler{studioFilter.CreatedAt, studioTable + ".created_at", nil},
		&timestampCriterionHandler{studioFilter.UpdatedAt, studioTable + ".updated_at", nil},

		&relatedFilterHandler{
			relatedIDCol:   "scenes.id",
			relatedRepo:    sceneRepository.repository,
			relatedHandler: &sceneFilterHandler{studioFilter.ScenesFilter},
			joinFn: func(f *filterBuilder) {
				studioRepository.scenes.innerJoin(f, "", "studios.id")
			},
		},

		&relatedFilterHandler{
			relatedIDCol:   "images.id",
			relatedRepo:    imageRepository.repository,
			relatedHandler: &imageFilterHandler{studioFilter.ImagesFilter},
			joinFn: func(f *filterBuilder) {
				studioRepository.images.innerJoin(f, "", "studios.id")
			},
		},

		&relatedFilterHandler{
			relatedIDCol:   "galleries.id",
			relatedRepo:    galleryRepository.repository,
			relatedHandler: &galleryFilterHandler{studioFilter.GalleriesFilter},
			joinFn: func(f *filterBuilder) {
				studioRepository.galleries.innerJoin(f, "", "studios.id")
			},
		},
	}
}

func (qb *studioFilterHandler) isMissingCriterionHandler(isMissing *string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if isMissing != nil && *isMissing != "" {
			switch *isMissing {
			case "url":
				studiosURLsTableMgr.join(f, "", "studios.id")
				f.addWhere("studio_urls.url IS NULL")
			case "image":
				f.addWhere("studios.image_blob IS NULL")
			case "stash_id":
				studioRepository.stashIDs.join(f, "studio_stash_ids", "studios.id")
				f.addWhere("studio_stash_ids.studio_id IS NULL")
			default:
				f.addWhere("(studios." + *isMissing + " IS NULL OR TRIM(studios." + *isMissing + ") = '')")
			}
		}
	}
}

func (qb *studioFilterHandler) sceneCountCriterionHandler(sceneCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if sceneCount != nil {
			f.addLeftJoin("scenes", "", "scenes.studio_id = studios.id")
			clause, args := getIntCriterionWhereClause("count(distinct scenes.id)", *sceneCount)

			f.addHaving(clause, args...)
		}
	}
}

func (qb *studioFilterHandler) imageCountCriterionHandler(imageCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if imageCount != nil {
			f.addLeftJoin("images", "", "images.studio_id = studios.id")
			clause, args := getIntCriterionWhereClause("count(distinct images.id)", *imageCount)

			f.addHaving(clause, args...)
		}
	}
}

func (qb *studioFilterHandler) galleryCountCriterionHandler(galleryCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if galleryCount != nil {
			f.addLeftJoin("galleries", "", "galleries.studio_id = studios.id")
			clause, args := getIntCriterionWhereClause("count(distinct galleries.id)", *galleryCount)

			f.addHaving(clause, args...)
		}
	}
}

func (qb *studioFilterHandler) tagCountCriterionHandler(tagCount *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: studioTable,
		joinTable:    studiosTagsTable,
		primaryFK:    studioIDColumn,
	}

	return h.handler(tagCount)
}

func (qb *studioFilterHandler) parentCriterionHandler(parents *models.MultiCriterionInput) criterionHandlerFunc {
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

func (qb *studioFilterHandler) aliasCriterionHandler(alias *models.StringCriterionInput) criterionHandlerFunc {
	h := stringListCriterionHandlerBuilder{
		primaryTable: studioTable,
		primaryFK:    studioIDColumn,
		joinTable:    studioAliasesTable,
		stringColumn: studioAliasColumn,
		addJoinTable: func(f *filterBuilder) {
			studiosAliasesTableMgr.join(f, "", "studios.id")
		},
	}

	return h.handler(alias)
}

func (qb *studioFilterHandler) urlsCriterionHandler(url *models.StringCriterionInput) criterionHandlerFunc {
	h := stringListCriterionHandlerBuilder{
		primaryTable: studioTable,
		primaryFK:    studioIDColumn,
		joinTable:    studioURLsTable,
		stringColumn: studioURLColumn,
		addJoinTable: func(f *filterBuilder) {
			studiosURLsTableMgr.join(f, "", "studios.id")
		},
	}

	return h.handler(url)
}

func (qb *studioFilterHandler) childCountCriterionHandler(childCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if childCount != nil {
			f.addLeftJoin("studios", "children_count", "children_count.parent_id = studios.id")
			clause, args := getIntCriterionWhereClause("count(distinct children_count.id)", *childCount)

			f.addHaving(clause, args...)
		}
	}
}

func (qb *studioFilterHandler) tagsCriterionHandler(tags *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	h := joinedHierarchicalMultiCriterionHandlerBuilder{
		primaryTable: studioTable,
		foreignTable: tagTable,
		foreignFK:    "tag_id",

		relationsTable: "tags_relations",
		joinTable:      studiosTagsTable,
		joinAs:         "studio_tag",
		primaryFK:      studioIDColumn,
	}

	return h.handler(tags)
}
