package sqlite

import (
	"context"
	stdsql "database/sql"
	"fmt"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

// GetFacets returns aggregated facet counts for galleries matching the given filter.
// This uses a single query with CTE to avoid re-executing the filter multiple times.
func (qb *GalleryStore) GetFacets(ctx context.Context, galleryFilter *models.GalleryFilterType, limit int) (*models.GalleryFacets, error) {
	result := &models.GalleryFacets{
		Tags:       []models.FacetCount{},
		Performers: []models.FacetCount{},
		Studios:    []models.FacetCount{},
		Organized:  []models.BooleanFacetCount{},
		Ratings:    []models.RatingFacetCount{},
	}

	query, err := qb.makeQuery(ctx, galleryFilter, nil)
	if err != nil {
		return nil, fmt.Errorf("error building base query: %w", err)
	}

	baseSQL := query.toSQL(false)
	args := append([]interface{}{}, query.args...)

	// Single query using CTE - filter executes only ONCE
	sql := fmt.Sprintf(`
		WITH filtered_galleries AS (%s)
		
		SELECT * FROM (
			SELECT 'tag' as facet_type, t.id, t.name as label, NULL as enum_value, COUNT(DISTINCT gt.gallery_id) as count
			FROM filtered_galleries fg
			INNER JOIN galleries_tags gt ON fg.id = gt.gallery_id
			INNER JOIN tags t ON gt.tag_id = t.id
			GROUP BY t.id
			ORDER BY count DESC
			LIMIT ?
		)
		
		UNION ALL
		
		SELECT * FROM (
			SELECT 'performer' as facet_type, p.id, p.name as label, NULL as enum_value, COUNT(DISTINCT pg.gallery_id) as count
			FROM filtered_galleries fg
			INNER JOIN performers_galleries pg ON fg.id = pg.gallery_id
			INNER JOIN performers p ON pg.performer_id = p.id
			GROUP BY p.id
			ORDER BY count DESC
			LIMIT ?
		)
		
		UNION ALL
		
		SELECT * FROM (
			SELECT 'studio' as facet_type, s.id, s.name as label, NULL as enum_value, COUNT(DISTINCT g.id) as count
			FROM filtered_galleries fg
			INNER JOIN galleries g ON fg.id = g.id
			INNER JOIN studios s ON g.studio_id = s.id
			WHERE g.studio_id IS NOT NULL
			GROUP BY s.id
			ORDER BY count DESC
			LIMIT ?
		)
		
		UNION ALL
		
		SELECT * FROM (
			SELECT 'organized' as facet_type, 0 as id, '' as label,
				CASE WHEN g.organized = 1 THEN 'true' ELSE 'false' END as enum_value,
				COUNT(*) as count
			FROM filtered_galleries fg
			INNER JOIN galleries g ON fg.id = g.id
			GROUP BY g.organized
		)
		
		UNION ALL
		
		SELECT * FROM (
			SELECT 'rating' as facet_type, 0 as id, '' as label,
				CAST(g.rating AS TEXT) as enum_value,
				COUNT(*) as count
			FROM filtered_galleries fg
			INNER JOIN galleries g ON fg.id = g.id
			WHERE g.rating IS NOT NULL
			GROUP BY g.rating
			ORDER BY g.rating DESC
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
		var enumValue stdsql.NullString
		var count int

		if err := rows.Scan(&facetType, &id, &label, &enumValue, &count); err != nil {
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
		case "organized":
			if enumValue.Valid {
				result.Organized = append(result.Organized, models.BooleanFacetCount{
					Value: enumValue.String == "true",
					Count: count,
				})
			}
		case "rating":
			if enumValue.Valid {
				rating, err := strconv.Atoi(enumValue.String)
				if err == nil {
					result.Ratings = append(result.Ratings, models.RatingFacetCount{
						Rating: rating,
						Count:  count,
					})
				}
			}
		}
	}

	return result, rows.Err()
}
