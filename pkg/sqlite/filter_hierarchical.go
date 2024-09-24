package sqlite

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

// hierarchicalRelationshipHandler provides handlers for parent, children, parent count, and child count criteria.
type hierarchicalRelationshipHandler struct {
	primaryTable  string
	relationTable string
	aliasPrefix   string
	parentIDCol   string
	childIDCol    string
}

func (h hierarchicalRelationshipHandler) validateModifier(m models.CriterionModifier) error {
	switch m {
	case models.CriterionModifierIncludesAll, models.CriterionModifierIncludes, models.CriterionModifierExcludes, models.CriterionModifierIsNull, models.CriterionModifierNotNull:
		// valid
		return nil
	default:
		return fmt.Errorf("invalid modifier %s", m)
	}
}

func (h hierarchicalRelationshipHandler) handleNullNotNull(f *filterBuilder, m models.CriterionModifier, isParents bool) {
	var notClause string
	if m == models.CriterionModifierNotNull {
		notClause = "NOT"
	}

	as := h.aliasPrefix + "_parents"
	col := h.childIDCol
	if !isParents {
		as = h.aliasPrefix + "_children"
		col = h.parentIDCol
	}

	// Based on:
	// f.addLeftJoin("tags_relations", "parent_relations", "tags.id = parent_relations.child_id")
	// f.addWhere(fmt.Sprintf("parent_relations.parent_id IS %s NULL", notClause))

	f.addLeftJoin(h.relationTable, as, fmt.Sprintf("%s.id = %s.%s", h.primaryTable, as, col))
	f.addWhere(fmt.Sprintf("%s.%s IS %s NULL", as, col, notClause))
}

func (h hierarchicalRelationshipHandler) parentsAlias() string {
	return h.aliasPrefix + "_parents"
}

func (h hierarchicalRelationshipHandler) childrenAlias() string {
	return h.aliasPrefix + "_children"
}

func (h hierarchicalRelationshipHandler) valueQuery(value []string, depth int, alias string, isParents bool) string {
	var depthCondition string
	if depth != -1 {
		depthCondition = fmt.Sprintf("WHERE depth < %d", depth)
	}

	queryTempl := `{alias} AS (
SELECT {root_id_col} AS root_id, {item_id_col} AS item_id, 0 AS depth FROM {relation_table} WHERE {root_id_col} IN` + getInBinding(len(value)) + `
UNION
SELECT root_id, {item_id_col}, depth + 1 FROM {relation_table} INNER JOIN {alias} ON item_id = {root_id_col} ` + depthCondition + `
)`

	var queryMap utils.StrFormatMap
	if isParents {
		queryMap = utils.StrFormatMap{
			"root_id_col": h.parentIDCol,
			"item_id_col": h.childIDCol,
		}
	} else {
		queryMap = utils.StrFormatMap{
			"root_id_col": h.childIDCol,
			"item_id_col": h.parentIDCol,
		}
	}

	queryMap["alias"] = alias
	queryMap["relation_table"] = h.relationTable

	return utils.StrFormat(queryTempl, queryMap)
}

func (h hierarchicalRelationshipHandler) handleValues(f *filterBuilder, c models.HierarchicalMultiCriterionInput, isParents bool, aliasSuffix string) {
	if len(c.Value) == 0 {
		return
	}

	var args []interface{}
	for _, val := range c.Value {
		args = append(args, val)
	}

	depthVal := 0
	if c.Depth != nil {
		depthVal = *c.Depth
	}

	tableAlias := h.parentsAlias()
	if !isParents {
		tableAlias = h.childrenAlias()
	}
	tableAlias += aliasSuffix

	query := h.valueQuery(c.Value, depthVal, tableAlias, isParents)
	f.addRecursiveWith(query, args...)

	f.addLeftJoin(tableAlias, "", fmt.Sprintf("%s.item_id = %s.id", tableAlias, h.primaryTable))
	addHierarchicalConditionClauses(f, c, tableAlias, "root_id")
}

func (h hierarchicalRelationshipHandler) handleValuesSimple(f *filterBuilder, value string, isParents bool) {
	joinCol := h.childIDCol
	valueCol := h.parentIDCol
	if !isParents {
		joinCol = h.parentIDCol
		valueCol = h.childIDCol
	}

	tableAlias := h.parentsAlias()
	if !isParents {
		tableAlias = h.childrenAlias()
	}

	f.addInnerJoin(h.relationTable, tableAlias, fmt.Sprintf("%s.%s = %s.id", tableAlias, joinCol, h.primaryTable))
	f.addWhere(fmt.Sprintf("%s.%s = ?", tableAlias, valueCol), value)
}

func (h hierarchicalRelationshipHandler) hierarchicalCriterionHandler(criterion *models.HierarchicalMultiCriterionInput, isParents bool) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if criterion != nil {
			c := criterion.CombineExcludes()

			// validate the modifier
			if err := h.validateModifier(c.Modifier); err != nil {
				f.setError(err)
				return
			}

			if c.Modifier == models.CriterionModifierIsNull || c.Modifier == models.CriterionModifierNotNull {
				h.handleNullNotNull(f, c.Modifier, isParents)
				return
			}

			if len(c.Value) == 0 && len(c.Excludes) == 0 {
				return
			}

			depth := 0
			if c.Depth != nil {
				depth = *c.Depth
			}

			// if we have a single include, no excludes, and no depth, we can use a simple join and where clause
			if (c.Modifier == models.CriterionModifierIncludes || c.Modifier == models.CriterionModifierIncludesAll) && len(c.Value) == 1 && len(c.Excludes) == 0 && depth == 0 {
				h.handleValuesSimple(f, c.Value[0], isParents)
				return
			}

			aliasSuffix := ""
			h.handleValues(f, c, isParents, aliasSuffix)

			if len(c.Excludes) > 0 {
				exCriterion := models.HierarchicalMultiCriterionInput{
					Value:    c.Excludes,
					Depth:    c.Depth,
					Modifier: models.CriterionModifierExcludes,
				}

				aliasSuffix := "2"
				h.handleValues(f, exCriterion, isParents, aliasSuffix)
			}
		}
	}
}

func (h hierarchicalRelationshipHandler) ParentsCriterionHandler(criterion *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	const isParents = true
	return h.hierarchicalCriterionHandler(criterion, isParents)
}

func (h hierarchicalRelationshipHandler) ChildrenCriterionHandler(criterion *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	const isParents = false
	return h.hierarchicalCriterionHandler(criterion, isParents)
}

func (h hierarchicalRelationshipHandler) countCriterionHandler(c *models.IntCriterionInput, isParents bool) criterionHandlerFunc {
	tableAlias := h.parentsAlias()
	col := h.childIDCol
	otherCol := h.parentIDCol
	if !isParents {
		tableAlias = h.childrenAlias()
		col = h.parentIDCol
		otherCol = h.childIDCol
	}
	tableAlias += "_count"

	return func(ctx context.Context, f *filterBuilder) {
		if c != nil {
			f.addLeftJoin(h.relationTable, tableAlias, fmt.Sprintf("%s.%s = %s.id", tableAlias, col, h.primaryTable))
			clause, args := getIntCriterionWhereClause(fmt.Sprintf("count(distinct %s.%s)", tableAlias, otherCol), *c)

			f.addHaving(clause, args...)
		}
	}
}

func (h hierarchicalRelationshipHandler) ParentCountCriterionHandler(parentCount *models.IntCriterionInput) criterionHandlerFunc {
	const isParents = true
	return h.countCriterionHandler(parentCount, isParents)
}

func (h hierarchicalRelationshipHandler) ChildCountCriterionHandler(childCount *models.IntCriterionInput) criterionHandlerFunc {
	const isParents = false
	return h.countCriterionHandler(childCount, isParents)
}
