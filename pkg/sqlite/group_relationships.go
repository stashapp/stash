package sqlite

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/models"
	"gopkg.in/guregu/null.v4"
	"gopkg.in/guregu/null.v4/zero"
)

type groupRelationshipRow struct {
	ContainingID int         `db:"containing_id"`
	SubID        int         `db:"sub_id"`
	OrderIndex   int         `db:"order_index"`
	Description  zero.String `db:"description"`
}

func (r groupRelationshipRow) resolve(useContainingID bool) models.GroupIDDescription {
	id := r.ContainingID
	if !useContainingID {
		id = r.SubID
	}

	return models.GroupIDDescription{
		GroupID:     id,
		Description: r.Description.String,
	}
}

type groupRelationshipStore struct {
	table *table
}

func (s *groupRelationshipStore) GetContainingGroupDescriptions(ctx context.Context, id int) ([]models.GroupIDDescription, error) {
	const idIsContaining = false
	return s.getGroupRelationships(ctx, id, idIsContaining)
}

func (s *groupRelationshipStore) GetSubGroupDescriptions(ctx context.Context, id int) ([]models.GroupIDDescription, error) {
	const idIsContaining = true
	return s.getGroupRelationships(ctx, id, idIsContaining)
}

func (s *groupRelationshipStore) getGroupRelationships(ctx context.Context, id int, idIsContaining bool) ([]models.GroupIDDescription, error) {
	col := "containing_id"
	if !idIsContaining {
		col = "sub_id"
	}

	table := s.table.table
	q := dialect.Select(table.All()).
		From(table).
		Where(table.Col(col).Eq(id)).
		Order(table.Col("order_index").Asc())

	const single = false
	var ret []models.GroupIDDescription
	if err := queryFunc(ctx, q, single, func(rows *sqlx.Rows) error {
		var row groupRelationshipRow
		if err := rows.StructScan(&row); err != nil {
			return err
		}

		ret = append(ret, row.resolve(!idIsContaining))

		return nil
	}); err != nil {
		return nil, fmt.Errorf("getting group relationships from %s: %w", table.GetTable(), err)
	}

	return ret, nil
}

// getMaxOrderIndex gets the maximum order index for the containing group with the given id
func (s *groupRelationshipStore) getMaxOrderIndex(ctx context.Context, containingID int) (int, error) {
	idColumn := s.table.table.Col("containing_id")

	q := dialect.Select(goqu.MAX("order_index")).
		From(s.table.table).
		Where(idColumn.Eq(containingID))

	var maxOrderIndex zero.Int
	if err := querySimple(ctx, q, &maxOrderIndex); err != nil {
		return 0, fmt.Errorf("getting max order index: %w", err)
	}

	return int(maxOrderIndex.Int64), nil
}

// createRelationships creates relationships between a group and other groups.
// If idIsContaining is true, the provided id is the containing group.
func (s *groupRelationshipStore) createRelationships(ctx context.Context, id int, d models.RelatedGroupDescriptions, idIsContaining bool) error {
	if d.Loaded() {
		for i, v := range d.List() {
			orderIndex := i + 1

			r := groupRelationshipRow{
				ContainingID: id,
				SubID:        v.GroupID,
				OrderIndex:   orderIndex,
				Description:  zero.StringFrom(v.Description),
			}

			if !idIsContaining {
				// get the max order index of the containing groups sub groups
				containingID := v.GroupID
				maxOrderIndex, err := s.getMaxOrderIndex(ctx, containingID)
				if err != nil {
					return err
				}

				r.ContainingID = v.GroupID
				r.SubID = id
				r.OrderIndex = maxOrderIndex + 1
			}

			_, err := s.table.insert(ctx, r)
			if err != nil {
				return fmt.Errorf("inserting into %s: %w", s.table.table.GetTable(), err)
			}
		}

		return nil
	}

	return nil
}

// createRelationships creates relationships between a group and other groups.
// If idIsContaining is true, the provided id is the containing group.
func (s *groupRelationshipStore) createContainingRelationships(ctx context.Context, id int, d models.RelatedGroupDescriptions) error {
	const idIsContaining = false
	return s.createRelationships(ctx, id, d, idIsContaining)
}

// createRelationships creates relationships between a group and other groups.
// If idIsContaining is true, the provided id is the containing group.
func (s *groupRelationshipStore) createSubRelationships(ctx context.Context, id int, d models.RelatedGroupDescriptions) error {
	const idIsContaining = true
	return s.createRelationships(ctx, id, d, idIsContaining)
}

func (s *groupRelationshipStore) replaceRelationships(ctx context.Context, id int, d models.RelatedGroupDescriptions, idIsContaining bool) error {
	// always destroy the existing relationships even if the new list is empty
	if err := s.destroyAllJoins(ctx, id, idIsContaining); err != nil {
		return err
	}

	return s.createRelationships(ctx, id, d, idIsContaining)
}

func (s *groupRelationshipStore) replaceContainingRelationships(ctx context.Context, id int, d models.RelatedGroupDescriptions) error {
	const idIsContaining = false
	return s.replaceRelationships(ctx, id, d, idIsContaining)
}

func (s *groupRelationshipStore) replaceSubRelationships(ctx context.Context, id int, d models.RelatedGroupDescriptions) error {
	const idIsContaining = true
	return s.replaceRelationships(ctx, id, d, idIsContaining)
}

func (s *groupRelationshipStore) modifyRelationships(ctx context.Context, id int, v *models.UpdateGroupDescriptions, idIsContaining bool) error {
	if v == nil {
		return nil
	}

	switch v.Mode {
	case models.RelationshipUpdateModeSet:
		return s.replaceJoins(ctx, id, *v, idIsContaining)
	case models.RelationshipUpdateModeAdd:
		return s.addJoins(ctx, id, v.Groups, idIsContaining)
	case models.RelationshipUpdateModeRemove:
		toRemove := make([]int, len(v.Groups))
		for i, vv := range v.Groups {
			toRemove[i] = vv.GroupID
		}
		return s.destroyJoins(ctx, id, toRemove, idIsContaining)
	}

	return nil
}

func (s *groupRelationshipStore) modifyContainingRelationships(ctx context.Context, id int, v *models.UpdateGroupDescriptions) error {
	const idIsContaining = false
	return s.modifyRelationships(ctx, id, v, idIsContaining)
}

func (s *groupRelationshipStore) modifySubRelationships(ctx context.Context, id int, v *models.UpdateGroupDescriptions) error {
	const idIsContaining = true
	return s.modifyRelationships(ctx, id, v, idIsContaining)
}

func (s *groupRelationshipStore) addJoins(ctx context.Context, id int, groups []models.GroupIDDescription, idIsContaining bool) error {
	// if we're adding to a containing group, get the max order index first
	var maxOrderIndex int
	if idIsContaining {
		var err error
		maxOrderIndex, err = s.getMaxOrderIndex(ctx, id)
		if err != nil {
			return err
		}
	}

	for i, vv := range groups {
		r := groupRelationshipRow{
			Description: zero.StringFrom(vv.Description),
		}

		if idIsContaining {
			r.ContainingID = id
			r.SubID = vv.GroupID
			r.OrderIndex = maxOrderIndex + (i + 1)
		} else {
			// get the max order index of the containing groups sub groups
			containingMaxOrderIndex, err := s.getMaxOrderIndex(ctx, vv.GroupID)
			if err != nil {
				return err
			}

			r.ContainingID = vv.GroupID
			r.SubID = id
			r.OrderIndex = containingMaxOrderIndex + 1
		}

		_, err := s.table.insert(ctx, r)
		if err != nil {
			return fmt.Errorf("inserting into %s: %w", s.table.table.GetTable(), err)
		}
	}

	return nil
}

func (s *groupRelationshipStore) destroyAllJoins(ctx context.Context, id int, idIsContaining bool) error {
	table := s.table.table
	idColumn := table.Col("containing_id")
	if !idIsContaining {
		idColumn = table.Col("sub_id")
	}

	q := dialect.Delete(table).Where(idColumn.Eq(id))

	if _, err := exec(ctx, q); err != nil {
		return fmt.Errorf("destroying %s: %w", table.GetTable(), err)
	}

	return nil
}

func (s *groupRelationshipStore) replaceJoins(ctx context.Context, id int, v models.UpdateGroupDescriptions, idIsContaining bool) error {
	if err := s.destroyAllJoins(ctx, id, idIsContaining); err != nil {
		return err
	}

	// convert to RelatedGroupDescriptions
	rgd := models.NewRelatedGroupDescriptions(v.Groups)
	return s.createRelationships(ctx, id, rgd, idIsContaining)
}

func (s *groupRelationshipStore) destroyJoins(ctx context.Context, id int, toRemove []int, idIsContaining bool) error {
	table := s.table.table
	idColumn := table.Col("containing_id")
	fkColumn := table.Col("sub_id")
	if !idIsContaining {
		idColumn = table.Col("sub_id")
		fkColumn = table.Col("containing_id")
	}

	q := dialect.Delete(table).Where(idColumn.Eq(id), fkColumn.In(toRemove))

	if _, err := exec(ctx, q); err != nil {
		return fmt.Errorf("destroying %s: %w", table.GetTable(), err)
	}

	return nil
}

func (s *groupRelationshipStore) getOrderIndexOfSubGroup(ctx context.Context, containingGroupID int, subGroupID int) (int, error) {
	table := s.table.table
	q := dialect.Select("order_index").
		From(table).
		Where(
			table.Col("containing_id").Eq(containingGroupID),
			table.Col("sub_id").Eq(subGroupID),
		)

	var orderIndex null.Int
	if err := querySimple(ctx, q, &orderIndex); err != nil {
		return 0, fmt.Errorf("getting order index: %w", err)
	}

	if !orderIndex.Valid {
		return 0, fmt.Errorf("sub-group %d not found in containing group %d", subGroupID, containingGroupID)
	}

	return int(orderIndex.Int64), nil
}

func (s *groupRelationshipStore) getGroupIDAtOrderIndex(ctx context.Context, containingGroupID int, orderIndex int) (*int, error) {
	table := s.table.table
	q := dialect.Select(table.Col("sub_id")).From(table).Where(
		table.Col("containing_id").Eq(containingGroupID),
		table.Col("order_index").Eq(orderIndex),
	)

	var ret null.Int
	if err := querySimple(ctx, q, &ret); err != nil {
		return nil, fmt.Errorf("getting sub id for order index: %w", err)
	}

	if !ret.Valid {
		return nil, nil
	}

	intRet := int(ret.Int64)
	return &intRet, nil
}

func (s *groupRelationshipStore) getOrderIndexAfterOrderIndex(ctx context.Context, containingGroupID int, orderIndex int) (int, error) {
	table := s.table.table
	q := dialect.Select(goqu.MIN("order_index")).From(table).Where(
		table.Col("containing_id").Eq(containingGroupID),
		table.Col("order_index").Gt(orderIndex),
	)

	var ret null.Int
	if err := querySimple(ctx, q, &ret); err != nil {
		return 0, fmt.Errorf("getting order index: %w", err)
	}

	if !ret.Valid {
		return orderIndex + 1, nil
	}

	return int(ret.Int64), nil
}

// incrementOrderIndexes increments the order_index value of all sub-groups in the containing group at or after the given index
func (s *groupRelationshipStore) incrementOrderIndexes(ctx context.Context, groupID int, indexBefore int) error {
	table := s.table.table

	// WORKAROUND - sqlite won't allow incrementing the value directly since it causes a
	// unique constraint violation.
	// Instead, we first set the order index to a negative value temporarily
	// see https://stackoverflow.com/a/7703239/695786
	q := dialect.Update(table).Set(exp.Record{
		"order_index": goqu.L("-order_index"),
	}).Where(
		table.Col("containing_id").Eq(groupID),
		table.Col("order_index").Gte(indexBefore),
	)

	if _, err := exec(ctx, q); err != nil {
		return fmt.Errorf("updating %s: %w", table.GetTable(), err)
	}

	q = dialect.Update(table).Set(exp.Record{
		"order_index": goqu.L("1-order_index"),
	}).Where(
		table.Col("containing_id").Eq(groupID),
		table.Col("order_index").Lt(0),
	)

	if _, err := exec(ctx, q); err != nil {
		return fmt.Errorf("updating %s: %w", table.GetTable(), err)
	}

	return nil
}

func (s *groupRelationshipStore) reorderSubGroup(ctx context.Context, groupID int, subGroupID int, insertPointID int, insertAfter bool) error {
	insertPointIndex, err := s.getOrderIndexOfSubGroup(ctx, groupID, insertPointID)
	if err != nil {
		return err
	}

	// if we're setting before
	if insertAfter {
		insertPointIndex, err = s.getOrderIndexAfterOrderIndex(ctx, groupID, insertPointIndex)
		if err != nil {
			return err
		}
	}

	// increment the order index of all sub-groups after and including the insertion point
	if err := s.incrementOrderIndexes(ctx, groupID, int(insertPointIndex)); err != nil {
		return err
	}

	// set the order index of the sub-group to the insertion point
	table := s.table.table
	q := dialect.Update(table).Set(exp.Record{
		"order_index": insertPointIndex,
	}).Where(
		table.Col("containing_id").Eq(groupID),
		table.Col("sub_id").Eq(subGroupID),
	)

	if _, err := exec(ctx, q); err != nil {
		return fmt.Errorf("updating %s: %w", table.GetTable(), err)
	}

	return nil
}

func (s *groupRelationshipStore) AddSubGroups(ctx context.Context, groupID int, subGroups []models.GroupIDDescription, insertIndex *int) error {
	const idIsContaining = true

	if err := s.addJoins(ctx, groupID, subGroups, idIsContaining); err != nil {
		return err
	}

	ids := make([]int, len(subGroups))
	for i, v := range subGroups {
		ids[i] = v.GroupID
	}

	if insertIndex != nil {
		// get the id of the sub-group at the insert index
		insertPointID, err := s.getGroupIDAtOrderIndex(ctx, groupID, *insertIndex)
		if err != nil {
			return err
		}

		if insertPointID == nil {
			// if the insert index is out of bounds, just assume adding to the end
			return nil
		}

		// reorder the sub-groups
		const insertAfter = false
		if err := s.ReorderSubGroups(ctx, groupID, ids, *insertPointID, insertAfter); err != nil {
			return err
		}
	}

	return nil
}

func (s *groupRelationshipStore) RemoveSubGroups(ctx context.Context, groupID int, subGroupIDs []int) error {
	const idIsContaining = true
	return s.destroyJoins(ctx, groupID, subGroupIDs, idIsContaining)
}

func (s *groupRelationshipStore) ReorderSubGroups(ctx context.Context, groupID int, subGroupIDs []int, insertPointID int, insertAfter bool) error {
	for _, id := range subGroupIDs {
		if err := s.reorderSubGroup(ctx, groupID, id, insertPointID, insertAfter); err != nil {
			return err
		}
	}

	return nil
}
