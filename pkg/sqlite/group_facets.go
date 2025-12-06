package sqlite

import (
	"context"
	stdsql "database/sql"
	"fmt"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

// GetFacets returns aggregated facet counts for groups matching the given filter.
// This uses a single query with CTE to avoid re-executing the filter multiple times.
func (qb *GroupStore) GetFacets(ctx context.Context, groupFilter *models.GroupFilterType, limit int) (*models.GroupFacets, error) {
	result := &models.GroupFacets{
		Tags:       []models.FacetCount{},
		Performers: []models.FacetCount{},
		Studios:    []models.FacetCount{},
	}

	query, err := qb.makeQuery(ctx, groupFilter, nil)
	if err != nil {
		return nil, fmt.Errorf("error building base query: %w", err)
	}

	baseSQL := query.toSQL(false)
	args := append([]interface{}{}, query.args...)

	// Single query using CTE - filter executes only ONCE
	sql := fmt.Sprintf(`
		WITH filtered_groups AS (%s)
		
		SELECT * FROM (
			SELECT 'tag' as facet_type, t.id, t.name as label, COUNT(DISTINCT gt.group_id) as count
			FROM filtered_groups fg
			INNER JOIN groups_tags gt ON fg.id = gt.group_id
			INNER JOIN tags t ON gt.tag_id = t.id
			GROUP BY t.id
			ORDER BY count DESC
			LIMIT ?
		)
		
		UNION ALL
		
		SELECT * FROM (
			SELECT 'performer' as facet_type, p.id, p.name as label, COUNT(DISTINCT gs.group_id) as count
			FROM filtered_groups fg
			INNER JOIN groups_scenes gs ON fg.id = gs.group_id
			INNER JOIN performers_scenes ps ON gs.scene_id = ps.scene_id
			INNER JOIN performers p ON ps.performer_id = p.id
			GROUP BY p.id
			ORDER BY count DESC
			LIMIT ?
		)
		
		UNION ALL
		
		SELECT * FROM (
			SELECT 'studio' as facet_type, s.id, s.name as label, COUNT(DISTINCT g.id) as count
			FROM filtered_groups fg
			INNER JOIN groups g ON fg.id = g.id
			INNER JOIN studios s ON g.studio_id = s.id
			WHERE g.studio_id IS NOT NULL
			GROUP BY s.id
			ORDER BY count DESC
			LIMIT ?
		)
	`, baseSQL)

	// Add limit args for tag, performer, studio
	args = append(args, limit, limit, limit)

	rows, err := dbWrapper.Queryx(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing facets query: %w", err)
	}
	defer rows.Close()

	// Parse results
	for rows.Next() {
		var facetType string
		var id int
		var label stdsql.NullString
		var count int

		if err := rows.Scan(&facetType, &id, &label, &count); err != nil {
			return nil, fmt.Errorf("error scanning facet row: %w", err)
		}

		switch facetType {
		case "tag":
			result.Tags = append(result.Tags, models.FacetCount{
				ID:    strconv.Itoa(id),
				Label: label.String,
				Count: count,
			})
		case "performer":
			result.Performers = append(result.Performers, models.FacetCount{
				ID:    strconv.Itoa(id),
				Label: label.String,
				Count: count,
			})
		case "studio":
			result.Studios = append(result.Studios, models.FacetCount{
				ID:    strconv.Itoa(id),
				Label: label.String,
				Count: count,
			})
		}
	}

	return result, rows.Err()
}
