package sqlite

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

type sceneMarkerFilterHandler struct {
	sceneMarkerFilter *models.SceneMarkerFilterType
}

func (qb *sceneMarkerFilterHandler) validate() error {
	return nil
}

func (qb *sceneMarkerFilterHandler) handle(ctx context.Context, f *filterBuilder) {
	sceneMarkerFilter := qb.sceneMarkerFilter
	if sceneMarkerFilter == nil {
		return
	}

	if err := qb.validate(); err != nil {
		f.setError(err)
		return
	}

	f.handleCriterion(ctx, qb.criterionHandler())
}

func (qb *sceneMarkerFilterHandler) joinScenes(f *filterBuilder) {
	sceneMarkerRepository.scenes.innerJoin(f, "", "scene_markers.scene_id")
}

func (qb *sceneMarkerFilterHandler) criterionHandler() criterionHandler {
	sceneMarkerFilter := qb.sceneMarkerFilter
	return compoundHandler{
		qb.tagIDCriterionHandler(sceneMarkerFilter.TagID),
		qb.tagsCriterionHandler(sceneMarkerFilter.Tags),
		qb.sceneTagsCriterionHandler(sceneMarkerFilter.SceneTags),
		qb.performersCriterionHandler(sceneMarkerFilter.Performers),
		qb.scenesCriterionHandler(sceneMarkerFilter.Scenes),
		floatCriterionHandler(sceneMarkerFilter.Duration, "COALESCE(scene_markers.end_seconds - scene_markers.seconds, NULL)", nil),
		&timestampCriterionHandler{sceneMarkerFilter.CreatedAt, "scene_markers.created_at", nil},
		&timestampCriterionHandler{sceneMarkerFilter.UpdatedAt, "scene_markers.updated_at", nil},
		&dateCriterionHandler{sceneMarkerFilter.SceneDate, "scenes.date", qb.joinScenes},
		&timestampCriterionHandler{sceneMarkerFilter.SceneCreatedAt, "scenes.created_at", qb.joinScenes},
		&timestampCriterionHandler{sceneMarkerFilter.SceneUpdatedAt, "scenes.updated_at", qb.joinScenes},

		&relatedFilterHandler{
			relatedIDCol:   "scenes.id",
			relatedRepo:    sceneRepository.repository,
			relatedHandler: &sceneFilterHandler{sceneMarkerFilter.SceneFilter},
			joinFn: func(f *filterBuilder) {
				qb.joinScenes(f)
			},
		},
	}
}

func (qb *sceneMarkerFilterHandler) tagIDCriterionHandler(tagID *string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if tagID != nil {
			f.addLeftJoin("scene_markers_tags", "", "scene_markers_tags.scene_marker_id = scene_markers.id")

			f.addWhere("(scene_markers.primary_tag_id = ? OR scene_markers_tags.tag_id = ?)", *tagID, *tagID)
		}
	}
}

func (qb *sceneMarkerFilterHandler) tagsCriterionHandler(criterion *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if criterion != nil {
			tags := criterion.CombineExcludes()

			if tags.Modifier == models.CriterionModifierIsNull || tags.Modifier == models.CriterionModifierNotNull {
				var notClause string
				if tags.Modifier == models.CriterionModifierNotNull {
					notClause = "NOT"
				}

				f.addLeftJoin("scene_markers_tags", "", "scene_markers.id = scene_markers_tags.scene_marker_id")

				f.addWhere(fmt.Sprintf("%s scene_markers_tags.tag_id IS NULL", notClause))
				return
			}

			if tags.Modifier == models.CriterionModifierEquals && tags.Depth != nil && *tags.Depth != 0 {
				f.setError(fmt.Errorf("depth is not supported for equals modifier for marker tag filtering"))
				return
			}

			if len(tags.Value) == 0 && len(tags.Excludes) == 0 {
				return
			}

			if len(tags.Value) > 0 {
				valuesClause, err := getHierarchicalValues(ctx, tags.Value, tagTable, "tags_relations", "parent_id", "child_id", tags.Depth)
				if err != nil {
					f.setError(err)
					return
				}

				f.addWith(`marker_tags AS (
	SELECT mt.scene_marker_id, t.column1 AS root_tag_id FROM scene_markers_tags mt
	INNER JOIN (` + valuesClause + `) t ON t.column2 = mt.tag_id
	UNION
	SELECT m.id, t.column1 FROM scene_markers m
	INNER JOIN (` + valuesClause + `) t ON t.column2 = m.primary_tag_id
	)`)

				f.addLeftJoin("marker_tags", "", "marker_tags.scene_marker_id = scene_markers.id")

				switch tags.Modifier {
				case models.CriterionModifierEquals:
					// includes only the provided ids
					f.addWhere("marker_tags.root_tag_id IS NOT NULL")
					tagsLen := len(tags.Value)
					f.addHaving(fmt.Sprintf("count(distinct marker_tags.root_tag_id) IS %d", tagsLen))
					// decrement by one to account for primary tag id
					f.addWhere("(SELECT COUNT(*) FROM scene_markers_tags s WHERE s.scene_marker_id = scene_markers.id) = ?", tagsLen-1)
				case models.CriterionModifierNotEquals:
					f.setError(fmt.Errorf("not equals modifier is not supported for scene marker tags"))
				default:
					addHierarchicalConditionClauses(f, tags, "marker_tags", "root_tag_id")
				}
			}

			if len(criterion.Excludes) > 0 {
				valuesClause, err := getHierarchicalValues(ctx, tags.Excludes, tagTable, "tags_relations", "parent_id", "child_id", tags.Depth)
				if err != nil {
					f.setError(err)
					return
				}

				clause := "scene_markers.id NOT IN (SELECT scene_markers_tags.scene_marker_id FROM scene_markers_tags WHERE scene_markers_tags.tag_id IN (SELECT column2 FROM (%s)))"
				f.addWhere(fmt.Sprintf(clause, valuesClause))

				f.addWhere(fmt.Sprintf("scene_markers.primary_tag_id NOT IN (SELECT column2 FROM (%s))", valuesClause))
			}
		}
	}
}

func (qb *sceneMarkerFilterHandler) sceneTagsCriterionHandler(tags *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if tags != nil {
			f.addLeftJoin("scenes_tags", "", "scene_markers.scene_id = scenes_tags.scene_id")

			h := joinedHierarchicalMultiCriterionHandlerBuilder{
				primaryTable: "scene_markers",
				primaryKey:   sceneIDColumn,
				foreignTable: tagTable,
				foreignFK:    tagIDColumn,

				relationsTable: "tags_relations",
				joinTable:      "scenes_tags",
				joinAs:         "marker_scenes_tags",
				primaryFK:      sceneIDColumn,
			}

			h.handler(tags).handle(ctx, f)
		}
	}
}

func (qb *sceneMarkerFilterHandler) performersCriterionHandler(performers *models.MultiCriterionInput) criterionHandlerFunc {
	h := joinedMultiCriterionHandlerBuilder{
		primaryTable: sceneTable,
		joinTable:    performersScenesTable,
		joinAs:       "performers_join",
		primaryFK:    sceneIDColumn,
		foreignFK:    performerIDColumn,

		addJoinTable: func(f *filterBuilder) {
			f.addLeftJoin(performersScenesTable, "performers_join", "performers_join.scene_id = scene_markers.scene_id")
		},
	}

	handler := h.handler(performers)
	return func(ctx context.Context, f *filterBuilder) {
		if performers == nil {
			return
		}

		// Make sure scenes is included, otherwise excludes filter fails
		qb.joinScenes(f)
		handler(ctx, f)
	}
}

func (qb *sceneMarkerFilterHandler) scenesCriterionHandler(scenes *models.MultiCriterionInput) criterionHandlerFunc {
	addJoinsFunc := func(f *filterBuilder) {
		f.addLeftJoin(sceneTable, "markers_scenes", "markers_scenes.id = scene_markers.scene_id")
	}
	h := multiCriterionHandlerBuilder{
		primaryTable: sceneMarkerTable,
		foreignTable: "markers_scenes",
		joinTable:    "",
		primaryFK:    sceneIDColumn,
		foreignFK:    sceneIDColumn,
		addJoinsFunc: addJoinsFunc,
	}
	return h.handler(scenes)
}
