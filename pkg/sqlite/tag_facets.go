package sqlite

import (
	"context"
	stdsql "database/sql"
	"fmt"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

// GetFacets returns aggregated facet counts for tags.
// Note: TagStore doesn't have complex filtering like other stores,
// so we use a simpler approach that queries all tags.
func (qb *TagStore) GetFacets(ctx context.Context, tagFilter *models.TagFilterType, limit int) (*models.TagFacets, error) {
	result := &models.TagFacets{
		Parents:  []models.FacetCount{},
		Children: []models.FacetCount{},
		Favorite: []models.BooleanFacetCount{},
	}

	// Single query to get all facets
	sql := `
		SELECT * FROM (
			SELECT 'parent' as facet_type, parent.id, parent.name as label, NULL as enum_value, COUNT(DISTINCT tr.child_id) as count
			FROM tags parent
			INNER JOIN tags_relations tr ON parent.id = tr.parent_id
			GROUP BY parent.id
			ORDER BY count DESC
			LIMIT ?
		)
		
		UNION ALL
		
		SELECT * FROM (
			SELECT 'child' as facet_type, child.id, child.name as label, NULL as enum_value, COUNT(DISTINCT tr.parent_id) as count
			FROM tags child
			INNER JOIN tags_relations tr ON child.id = tr.child_id
			GROUP BY child.id
			ORDER BY count DESC
			LIMIT ?
		)
		
		UNION ALL
		
		SELECT * FROM (
			SELECT 'favorite' as facet_type, 0 as id, '' as label,
				CASE WHEN favorite = 1 THEN 'true' ELSE 'false' END as enum_value,
				COUNT(*) as count
			FROM tags
			GROUP BY favorite
		)
	`

	rows, err := dbWrapper.Queryx(ctx, sql, limit, limit)
	if err != nil {
		return nil, fmt.Errorf("error executing facets query: %w", err)
	}
	defer rows.Close()

	// Parse results
	for rows.Next() {
		var facetType string
		var id int
		var label stdsql.NullString
		var enumValue stdsql.NullString
		var count int

		if err := rows.Scan(&facetType, &id, &label, &enumValue, &count); err != nil {
			return nil, fmt.Errorf("error scanning facet row: %w", err)
		}

		switch facetType {
		case "parent":
			result.Parents = append(result.Parents, models.FacetCount{
				ID:    strconv.Itoa(id),
				Label: label.String,
				Count: count,
			})
		case "child":
			result.Children = append(result.Children, models.FacetCount{
				ID:    strconv.Itoa(id),
				Label: label.String,
				Count: count,
			})
		case "favorite":
			if enumValue.Valid {
				result.Favorite = append(result.Favorite, models.BooleanFacetCount{
					Value: enumValue.String == "true",
					Count: count,
				})
			}
		}
	}

	return result, rows.Err()
}
