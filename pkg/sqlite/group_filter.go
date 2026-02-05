package sqlite

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
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

var groupHierarchyHandler = hierarchicalRelationshipHandler{
	primaryTable:  groupTable,
	relationTable: groupRelationsTable,
	aliasPrefix:   groupTable,
	parentIDCol:   "containing_id",
	childIDCol:    "sub_id",
}

func (qb *groupFilterHandler) criterionHandler() criterionHandler {
	groupFilter := qb.groupFilter
	return compoundHandler{
		stringCriterionHandler(groupFilter.Name, "groups.name"),
		stringCriterionHandler(groupFilter.Director, "groups.director"),
		stringCriterionHandler(groupFilter.Synopsis, "groups.description"),
		intCriterionHandler(groupFilter.Rating100, "groups.rating", nil),
		floatIntCriterionHandler(groupFilter.Duration, "groups.duration", nil),
		qb.missingCriterionHandler(groupFilter.IsMissing),
		qb.urlsCriterionHandler(groupFilter.URL),
		studioCriterionHandler(groupTable, groupFilter.Studios),
		qb.performersCriterionHandler(groupFilter.Performers),
		qb.tagsCriterionHandler(groupFilter.Tags),
		qb.tagCountCriterionHandler(groupFilter.TagCount),
		qb.groupOCounterCriterionHandler(groupFilter.OCounter),
		&dateCriterionHandler{groupFilter.Date, "groups.date", nil},
		groupHierarchyHandler.ParentsCriterionHandler(groupFilter.ContainingGroups),
		groupHierarchyHandler.ChildrenCriterionHandler(groupFilter.SubGroups),
		groupHierarchyHandler.ParentCountCriterionHandler(groupFilter.ContainingGroupCount),
		groupHierarchyHandler.ChildCountCriterionHandler(groupFilter.SubGroupCount),
		&timestampCriterionHandler{groupFilter.CreatedAt, "groups.created_at", nil},
		&timestampCriterionHandler{groupFilter.UpdatedAt, "groups.updated_at", nil},

		&relatedFilterHandler{
			relatedIDCol:   "groups_scenes.scene_id",
			relatedRepo:    sceneRepository.repository,
			relatedHandler: &sceneFilterHandler{groupFilter.ScenesFilter},
			joinFn: func(f *filterBuilder) {
				groupRepository.scenes.innerJoin(f, "", "groups.id")
			},
		},

		&relatedFilterHandler{
			relatedIDCol:   "groups.studio_id",
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
				f.addWhere("groups.front_image_blob IS NULL")
			case "back_image":
				f.addWhere("groups.back_image_blob IS NULL")
			case "scenes":
				f.addLeftJoin("groups_scenes", "", "groups_scenes.group_id = groups.id")
				f.addWhere("groups_scenes.scene_id IS NULL")
			default:
				f.addWhere("(groups." + *isMissing + " IS NULL OR TRIM(groups." + *isMissing + ") = '')")
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
			groupsURLsTableMgr.join(f, "", "groups.id")
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

				f.addLeftJoin("groups_scenes", "", "groups.id = groups_scenes.group_id")
				f.addLeftJoin("performers_scenes", "", "groups_scenes.scene_id = performers_scenes.scene_id")

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
			f.addWith(`groups_performers AS (
				SELECT groups_scenes.group_id, performers_scenes.performer_id
				FROM groups_scenes
				INNER JOIN performers_scenes ON groups_scenes.scene_id = performers_scenes.scene_id
				WHERE performers_scenes.performer_id IN`+getInBinding(len(performers.Value))+`
			)`, args...)
			f.addLeftJoin("groups_performers", "", "groups.id = groups_performers.group_id")

			switch performers.Modifier {
			case models.CriterionModifierIncludes:
				f.addWhere("groups_performers.performer_id IS NOT NULL")
			case models.CriterionModifierIncludesAll:
				f.addWhere("groups_performers.performer_id IS NOT NULL")
				f.addHaving("COUNT(DISTINCT groups_performers.performer_id) = ?", len(performers.Value))
			case models.CriterionModifierExcludes:
				f.addWhere("groups_performers.performer_id IS NULL")
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
		joinAs:         "group_tag",
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

// used for sorting and filtering on group o-count
var selectGroupOCountSQL = utils.StrFormat(
	"SELECT SUM(o_counter) "+
		"FROM ("+
		"SELECT COUNT({scenes_o_dates}.{o_date}) as o_counter from {groups_scenes} s "+
		"LEFT JOIN {scenes} ON {scenes}.id = s.{scene_id} "+
		"LEFT JOIN {scenes_o_dates} ON {scenes_o_dates}.{scene_id} = {scenes}.id "+
		"WHERE s.{group_id} = {group}.id "+
		")",
	map[string]interface{}{
		"group":          groupTable,
		"group_id":       groupIDColumn,
		"groups_scenes":  groupsScenesTable,
		"scenes":         sceneTable,
		"scene_id":       sceneIDColumn,
		"scenes_o_dates": scenesODatesTable,
		"o_date":         sceneODateColumn,
	},
)

func (qb *groupFilterHandler) groupOCounterCriterionHandler(count *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if count == nil {
			return
		}

		lhs := "(" + selectGroupOCountSQL + ")"
		clause, args := getIntCriterionWhereClause(lhs, *count)

		f.addWhere(clause, args...)
	}

}
