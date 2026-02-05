package sqlite

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

type tagFilterHandler struct {
	tagFilter *models.TagFilterType
}

func (qb *tagFilterHandler) validate() error {
	tagFilter := qb.tagFilter
	if tagFilter == nil {
		return nil
	}

	if err := validateFilterCombination(tagFilter.OperatorFilter); err != nil {
		return err
	}

	if subFilter := tagFilter.SubFilter(); subFilter != nil {
		sqb := &tagFilterHandler{tagFilter: subFilter}
		if err := sqb.validate(); err != nil {
			return err
		}
	}

	return nil
}

func (qb *tagFilterHandler) handle(ctx context.Context, f *filterBuilder) {
	tagFilter := qb.tagFilter
	if tagFilter == nil {
		return
	}

	if err := qb.validate(); err != nil {
		f.setError(err)
		return
	}

	sf := tagFilter.SubFilter()
	if sf != nil {
		sub := &tagFilterHandler{sf}
		handleSubFilter(ctx, sub, f, tagFilter.OperatorFilter)
	}

	f.handleCriterion(ctx, qb.criterionHandler())
}

var tagHierarchyHandler = hierarchicalRelationshipHandler{
	primaryTable:  tagTable,
	relationTable: tagRelationsTable,
	aliasPrefix:   tagTable,
	parentIDCol:   "parent_id",
	childIDCol:    "child_id",
}

func (qb *tagFilterHandler) criterionHandler() criterionHandler {
	tagFilter := qb.tagFilter
	return compoundHandler{
		stringCriterionHandler(tagFilter.Name, tagTable+".name"),
		stringCriterionHandler(tagFilter.SortName, tagTable+".sort_name"),
		qb.aliasCriterionHandler(tagFilter.Aliases),

		boolCriterionHandler(tagFilter.Favorite, tagTable+".favorite", nil),
		stringCriterionHandler(tagFilter.Description, tagTable+".description"),
		boolCriterionHandler(tagFilter.IgnoreAutoTag, tagTable+".ignore_auto_tag", nil),

		qb.isMissingCriterionHandler(tagFilter.IsMissing),
		qb.sceneCountCriterionHandler(tagFilter.SceneCount),
		qb.imageCountCriterionHandler(tagFilter.ImageCount),
		qb.galleryCountCriterionHandler(tagFilter.GalleryCount),
		qb.performerCountCriterionHandler(tagFilter.PerformerCount),
		qb.studioCountCriterionHandler(tagFilter.StudioCount),

		qb.groupCountCriterionHandler(tagFilter.GroupCount),
		qb.groupCountCriterionHandler(tagFilter.MovieCount),

		qb.markerCountCriterionHandler(tagFilter.MarkerCount),
		tagHierarchyHandler.ParentsCriterionHandler(tagFilter.Parents),
		tagHierarchyHandler.ChildrenCriterionHandler(tagFilter.Children),
		tagHierarchyHandler.ParentCountCriterionHandler(tagFilter.ParentCount),
		tagHierarchyHandler.ChildCountCriterionHandler(tagFilter.ChildCount),

		&stashIDCriterionHandler{
			c:                 tagFilter.StashIDEndpoint,
			stashIDRepository: &tagRepository.stashIDs,
			stashIDTableAs:    "tag_stash_ids",
			parentIDCol:       "tags.id",
		},
		&stashIDsCriterionHandler{
			c:                 tagFilter.StashIDsEndpoint,
			stashIDRepository: &tagRepository.stashIDs,
			stashIDTableAs:    "tag_stash_ids",
			parentIDCol:       "tags.id",
		},

		&timestampCriterionHandler{tagFilter.CreatedAt, "tags.created_at", nil},
		&timestampCriterionHandler{tagFilter.UpdatedAt, "tags.updated_at", nil},

		&relatedFilterHandler{
			relatedIDCol:   "scenes_tags.scene_id",
			relatedRepo:    sceneRepository.repository,
			relatedHandler: &sceneFilterHandler{tagFilter.ScenesFilter},
			joinFn: func(f *filterBuilder) {
				tagRepository.scenes.innerJoin(f, "", "tags.id")
			},
		},

		&relatedFilterHandler{
			relatedIDCol:   "images_tags.image_id",
			relatedRepo:    imageRepository.repository,
			relatedHandler: &imageFilterHandler{tagFilter.ImagesFilter},
			joinFn: func(f *filterBuilder) {
				tagRepository.images.innerJoin(f, "", "tags.id")
			},
		},

		&relatedFilterHandler{
			relatedIDCol:   "galleries_tags.gallery_id",
			relatedRepo:    galleryRepository.repository,
			relatedHandler: &galleryFilterHandler{tagFilter.GalleriesFilter},
			joinFn: func(f *filterBuilder) {
				tagRepository.galleries.innerJoin(f, "", "tags.id")
			},
		},

		&relatedFilterHandler{
			relatedIDCol:   "performers_tags.performer_id",
			relatedRepo:    performerRepository.repository,
			relatedHandler: &performerFilterHandler{tagFilter.PerformersFilter},
			joinFn: func(f *filterBuilder) {
				tagRepository.performers.innerJoin(f, "", "tags.id")
			},
		},
	}
}

func (qb *tagFilterHandler) aliasCriterionHandler(alias *models.StringCriterionInput) criterionHandlerFunc {
	h := stringListCriterionHandlerBuilder{
		primaryTable: tagTable,
		primaryFK:    tagIDColumn,
		joinTable:    tagAliasesTable,
		stringColumn: tagAliasColumn,
		addJoinTable: func(f *filterBuilder) {
			tagRepository.aliases.join(f, "", "tags.id")
		},
	}

	return h.handler(alias)
}

func (qb *tagFilterHandler) isMissingCriterionHandler(isMissing *string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if isMissing != nil && *isMissing != "" {
			switch *isMissing {
			case "image":
				f.addWhere("tags.image_blob IS NULL")
			default:
				f.addWhere("(tags." + *isMissing + " IS NULL OR TRIM(tags." + *isMissing + ") = '')")
			}
		}
	}
}

func (qb *tagFilterHandler) sceneCountCriterionHandler(sceneCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if sceneCount != nil {
			f.addLeftJoin("scenes_tags", "", "scenes_tags.tag_id = tags.id")
			clause, args := getIntCriterionWhereClause("count(distinct scenes_tags.scene_id)", *sceneCount)

			f.addHaving(clause, args...)
		}
	}
}

func (qb *tagFilterHandler) imageCountCriterionHandler(imageCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if imageCount != nil {
			f.addLeftJoin("images_tags", "", "images_tags.tag_id = tags.id")
			clause, args := getIntCriterionWhereClause("count(distinct images_tags.image_id)", *imageCount)

			f.addHaving(clause, args...)
		}
	}
}

func (qb *tagFilterHandler) galleryCountCriterionHandler(galleryCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if galleryCount != nil {
			f.addLeftJoin("galleries_tags", "", "galleries_tags.tag_id = tags.id")
			clause, args := getIntCriterionWhereClause("count(distinct galleries_tags.gallery_id)", *galleryCount)

			f.addHaving(clause, args...)
		}
	}
}

func (qb *tagFilterHandler) performerCountCriterionHandler(performerCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if performerCount != nil {
			f.addLeftJoin("performers_tags", "", "performers_tags.tag_id = tags.id")
			clause, args := getIntCriterionWhereClause("count(distinct performers_tags.performer_id)", *performerCount)

			f.addHaving(clause, args...)
		}
	}
}

func (qb *tagFilterHandler) studioCountCriterionHandler(studioCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if studioCount != nil {
			f.addLeftJoin("studios_tags", "", "studios_tags.tag_id = tags.id")
			clause, args := getIntCriterionWhereClause("count(distinct studios_tags.studio_id)", *studioCount)

			f.addHaving(clause, args...)
		}
	}
}

func (qb *tagFilterHandler) groupCountCriterionHandler(groupCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if groupCount != nil {
			f.addLeftJoin("groups_tags", "", "groups_tags.tag_id = tags.id")
			clause, args := getIntCriterionWhereClause("count(distinct groups_tags.group_id)", *groupCount)

			f.addHaving(clause, args...)
		}
	}
}

func (qb *tagFilterHandler) markerCountCriterionHandler(markerCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if markerCount != nil {
			f.addLeftJoin("scene_markers_tags", "", "scene_markers_tags.tag_id = tags.id")
			f.addLeftJoin("scene_markers", "", "scene_markers_tags.scene_marker_id = scene_markers.id OR scene_markers.primary_tag_id = tags.id")
			clause, args := getIntCriterionWhereClause("count(distinct scene_markers.id)", *markerCount)

			f.addHaving(clause, args...)
		}
	}
}
