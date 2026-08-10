package sqlite

import (
	"context"
	"fmt"

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

	f.handleCriterion(ctx, qb.criterionHandler())

	sf := tagFilter.SubFilter()
	if sf != nil {
		sub := &tagFilterHandler{sf}
		handleSubFilter(ctx, sub, f, tagFilter.OperatorFilter)
	}
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
		qb.audioCountCriterionHandler(tagFilter.AudioCount),
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

		&customFieldsFilterHandler{
			table: tagsCustomFieldsTable.GetTable(),
			fkCol: tagIDColumn,
			c:     tagFilter.CustomFields,
			idCol: "tags.id",
		},

		&relatedFilterHandler{
			relatedIDCol:   "scenes_tags.scene_id",
			relatedRepo:    sceneRepository.repository,
			relatedHandler: &sceneFilterHandler{tagFilter.ScenesFilter},
			joinFn: func(f *filterBuilder) {
				tagRepository.scenes.innerJoin(f, "", "tags.id")
			},
		},

		&relatedFilterHandler{
			relatedIDCol:   "audios_tags.audio_id",
			relatedRepo:    audioRepository.repository,
			relatedHandler: &audioFilterHandler{tagFilter.AudiosFilter},
			joinFn: func(f *filterBuilder) {
				tagRepository.audios.innerJoin(f, "", "tags.id")
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
			relatedIDCol:   "groups_tags.group_id",
			relatedRepo:    groupRepository.repository,
			relatedHandler: &groupFilterHandler{tagFilter.GroupsFilter},
			joinFn: func(f *filterBuilder) {
				tagRepository.groups.innerJoin(f, "", "tags.id")
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

		&relatedFilterHandler{
			relatedIDCol:   "studios_tags.studio_id",
			relatedRepo:    studioRepository.repository,
			relatedHandler: &studioFilterHandler{tagFilter.StudiosFilter},
			joinFn: func(f *filterBuilder) {
				tagRepository.studios.innerJoin(f, "", "tags.id")
			},
		},

		&relatedFilterHandler{
			relatedIDCol:   "markers_tags.marker_id",
			relatedRepo:    sceneMarkerRepository.repository,
			relatedHandler: &sceneMarkerFilterHandler{tagFilter.MarkersFilter},
			joinFn: func(f *filterBuilder) {
				f.addWith(`markers_tags AS (
				SELECT mt.scene_marker_id AS marker_id, mt.tag_id AS tag_id FROM scene_markers_tags mt
				UNION
				SELECT m.id, m.primary_tag_id FROM scene_markers m
				)`)
				f.addInnerJoin("markers_tags", "", "markers_tags.tag_id = tags.id")
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
		addJoinTable: func(f *filterBuilder, joinType joinType) {
			tagRepository.aliases.join(f, joinType, "", "tags.id")
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
			case "aliases":
				tagRepository.aliases.leftJoin(f, "", "tags.id")
				f.addWhere("tag_aliases.alias IS NULL")
			case "stash_id":
				tagRepository.stashIDs.leftJoin(f, "tag_stash_ids", "tags.id")
				f.addWhere("tag_stash_ids.tag_id IS NULL")
			default:
				if err := validateIsMissing(*isMissing, []string{
					"description",
				}); err != nil {
					f.setError(err)
					return
				}
				f.addWhere("(tags." + *isMissing + " IS NULL OR TRIM(tags." + *isMissing + ") = '')")
			}
		}
	}
}

// addHierarchicalCountCTE adds a recursive CTE to walk the tag hierarchy.
// depth < 0 includes all descendant levels (unlimited recursion).
// depth >= 0 limits the recursion to that many levels from the root tag.
func (qb *tagFilterHandler) addHierarchicalCountCTE(f *filterBuilder, cteAlias string, depth int) {
	if depth < 0 {
		// unlimited recursion — no depth tracking needed
		f.addRecursiveWith(fmt.Sprintf(
			`%[1]s(root_id, descendant_id) AS (
				SELECT id, id FROM tags
				UNION ALL
				SELECT td.root_id, tr.child_id
				FROM tags_relations tr
				INNER JOIN %[1]s td ON td.descendant_id = tr.parent_id
			)`, cteAlias))
	} else {
		// depth-limited: track recursion level as a CTE column
		f.addRecursiveWith(fmt.Sprintf(
			`%[1]s(root_id, descendant_id, depth) AS (
				SELECT id, id, 0 FROM tags
				UNION ALL
				SELECT td.root_id, tr.child_id, td.depth + 1
				FROM tags_relations tr
				INNER JOIN %[1]s td ON td.descendant_id = tr.parent_id
				WHERE td.depth < %[2]d
			)`, cteAlias, depth))
	}
}

func (qb *tagFilterHandler) hierarchicalCountHandler(input *models.HierarchicalCountInput, joinTable, entityIDCol string) criterionHandlerFunc {
	return qb.hierarchicalCountHandlerOnCol(input, joinTable, "tag_id", entityIDCol)
}

func (qb *tagFilterHandler) hierarchicalCountHandlerOnCol(input *models.HierarchicalCountInput, joinTable, tagCol, entityIDCol string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if input == nil {
			return
		}

		intInput := models.IntCriterionInput{
			Value:    input.Value,
			Value2:   input.Value2,
			Modifier: input.Modifier,
		}

		if input.Depth != nil {
			cteAlias := joinTable + "_desc"
			qb.addHierarchicalCountCTE(f, cteAlias, *input.Depth)
			f.addLeftJoin(cteAlias, "", fmt.Sprintf("%s.root_id = tags.id", cteAlias))
			f.addLeftJoin(joinTable, "", fmt.Sprintf("%[1]s.%[2]s = %[3]s.descendant_id", joinTable, tagCol, cteAlias))
		} else {
			f.addLeftJoin(joinTable, "", fmt.Sprintf("%s.%s = tags.id", joinTable, tagCol))
		}

		clause, args := getIntCriterionWhereClause(fmt.Sprintf("count(distinct %s.%s)", joinTable, entityIDCol), intInput)
		f.addHaving(clause, args...)
	}
}

func (qb *tagFilterHandler) sceneCountCriterionHandler(sceneCount *models.HierarchicalCountInput) criterionHandlerFunc {
	return qb.hierarchicalCountHandler(sceneCount, "scenes_tags", "scene_id")
}

func (qb *tagFilterHandler) audioCountCriterionHandler(audioCount *models.HierarchicalCountInput) criterionHandlerFunc {
	return qb.hierarchicalCountHandler(audioCount, "audios_tags", "audio_id")
}

func (qb *tagFilterHandler) imageCountCriterionHandler(imageCount *models.HierarchicalCountInput) criterionHandlerFunc {
	return qb.hierarchicalCountHandler(imageCount, "images_tags", "image_id")
}

func (qb *tagFilterHandler) galleryCountCriterionHandler(galleryCount *models.HierarchicalCountInput) criterionHandlerFunc {
	return qb.hierarchicalCountHandler(galleryCount, "galleries_tags", "gallery_id")
}

func (qb *tagFilterHandler) performerCountCriterionHandler(performerCount *models.HierarchicalCountInput) criterionHandlerFunc {
	return qb.hierarchicalCountHandler(performerCount, "performers_tags", "performer_id")
}

func (qb *tagFilterHandler) studioCountCriterionHandler(studioCount *models.HierarchicalCountInput) criterionHandlerFunc {
	return qb.hierarchicalCountHandler(studioCount, "studios_tags", "studio_id")
}

func (qb *tagFilterHandler) groupCountCriterionHandler(groupCount *models.HierarchicalCountInput) criterionHandlerFunc {
	return qb.hierarchicalCountHandler(groupCount, "groups_tags", "group_id")
}

func (qb *tagFilterHandler) markerCountCriterionHandler(markerCount *models.HierarchicalCountInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if markerCount == nil {
			return
		}

		intInput := models.IntCriterionInput{
			Value:    markerCount.Value,
			Value2:   markerCount.Value2,
			Modifier: markerCount.Modifier,
		}

		if markerCount.Depth != nil {
			cteAlias := "scene_markers_desc"
			qb.addHierarchicalCountCTE(f, cteAlias, *markerCount.Depth)
			f.addLeftJoin(cteAlias, "", fmt.Sprintf("%s.root_id = tags.id", cteAlias))
			f.addLeftJoin("scene_markers_tags", "", fmt.Sprintf("scene_markers_tags.tag_id = %s.descendant_id", cteAlias))
			f.addLeftJoin("scene_markers", "", "scene_markers_tags.scene_marker_id = scene_markers.id OR scene_markers.primary_tag_id = "+cteAlias+".descendant_id")
		} else {
			f.addLeftJoin("scene_markers_tags", "", "scene_markers_tags.tag_id = tags.id")
			f.addLeftJoin("scene_markers", "", "scene_markers_tags.scene_marker_id = scene_markers.id OR scene_markers.primary_tag_id = tags.id")
		}

		clause, args := getIntCriterionWhereClause("count(distinct scene_markers.id)", intInput)
		f.addHaving(clause, args...)
	}
}
