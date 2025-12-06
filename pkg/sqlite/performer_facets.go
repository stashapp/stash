package sqlite

import (
	"context"
	stdsql "database/sql"
	"fmt"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

// GetFacets returns aggregated facet counts for performers matching the given filter.
// This uses a single query with CTE to avoid re-executing the filter multiple times.
func (qb *PerformerStore) GetFacets(ctx context.Context, performerFilter *models.PerformerFilterType, limit int) (*models.PerformerFacets, error) {
	result := &models.PerformerFacets{
		Tags:        []models.FacetCount{},
		Studios:     []models.FacetCount{},
		Genders:     []models.GenderFacetCount{},
		Countries:   []models.FacetCount{},
		Circumcised: []models.CircumcisedFacetCount{},
		Favorite:    []models.BooleanFacetCount{},
		Ratings:     []models.RatingFacetCount{},
	}

	query, err := qb.makeQuery(ctx, performerFilter, nil)
	if err != nil {
		return nil, fmt.Errorf("error building base query: %w", err)
	}

	baseSQL := query.toSQL(false)
	args := append([]interface{}{}, query.args...)

	// Single query using CTE - filter executes only ONCE
	sql := fmt.Sprintf(`
		WITH filtered_performers AS (%s)
		
		SELECT * FROM (
			SELECT 'tag' as facet_type, t.id, t.name as label, NULL as enum_value, COUNT(DISTINCT pt.performer_id) as count
			FROM filtered_performers fp
			INNER JOIN performers_tags pt ON fp.id = pt.performer_id
			INNER JOIN tags t ON pt.tag_id = t.id
			GROUP BY t.id
			ORDER BY count DESC
			LIMIT ?
		)
		
		UNION ALL
		
		SELECT * FROM (
			SELECT 'studio' as facet_type, s.id, s.name as label, NULL as enum_value, COUNT(DISTINCT ps.performer_id) as count
			FROM filtered_performers fp
			INNER JOIN performers_scenes ps ON fp.id = ps.performer_id
			INNER JOIN scenes sc ON ps.scene_id = sc.id
			INNER JOIN studios s ON sc.studio_id = s.id
			WHERE sc.studio_id IS NOT NULL
			GROUP BY s.id
			ORDER BY count DESC
			LIMIT ?
		)
		
		UNION ALL
		
		SELECT * FROM (
			SELECT 'gender' as facet_type, 0 as id, '' as label, p.gender as enum_value, COUNT(*) as count
			FROM filtered_performers fp
			INNER JOIN performers p ON fp.id = p.id
			WHERE p.gender IS NOT NULL AND p.gender != ''
			GROUP BY p.gender
			ORDER BY count DESC
		)
		
		UNION ALL
		
		SELECT * FROM (
			SELECT 'country' as facet_type, 0 as id, p.country as label, p.country as enum_value, COUNT(*) as count
			FROM filtered_performers fp
			INNER JOIN performers p ON fp.id = p.id
			WHERE p.country IS NOT NULL AND p.country != ''
			GROUP BY p.country
			ORDER BY count DESC
			LIMIT ?
		)
		
		UNION ALL
		
		SELECT * FROM (
			SELECT 'favorite' as facet_type, 0 as id, '' as label,
				CASE WHEN p.favorite = 1 THEN 'true' ELSE 'false' END as enum_value,
				COUNT(*) as count
			FROM filtered_performers fp
			INNER JOIN performers p ON fp.id = p.id
			GROUP BY p.favorite
		)
		
		UNION ALL
		
		SELECT * FROM (
			SELECT 'circumcised' as facet_type, 0 as id, '' as label,
				p.circumcised as enum_value,
				COUNT(*) as count
			FROM filtered_performers fp
			INNER JOIN performers p ON fp.id = p.id
			WHERE p.circumcised IS NOT NULL AND p.circumcised != ''
			GROUP BY p.circumcised
			ORDER BY count DESC
		)
		
		UNION ALL
		
		SELECT * FROM (
			SELECT 'rating' as facet_type, 0 as id, '' as label,
				CAST(p.rating AS TEXT) as enum_value,
				COUNT(*) as count
			FROM filtered_performers fp
			INNER JOIN performers p ON fp.id = p.id
			WHERE p.rating IS NOT NULL
			GROUP BY p.rating
			ORDER BY p.rating DESC
		)
	`, baseSQL)

	// Add limit args for tag, studio, country
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
		case "studio":
			result.Studios = append(result.Studios, models.FacetCount{
				ID:    strconv.Itoa(id),
				Label: label.String,
				Count: count,
			})
		case "gender":
			if enumValue.Valid {
				gender := models.GenderEnum(enumValue.String)
				if gender.IsValid() {
					result.Genders = append(result.Genders, models.GenderFacetCount{
						Gender: gender,
						Count:  count,
					})
				}
			}
		case "country":
			if label.Valid && label.String != "" {
				result.Countries = append(result.Countries, models.FacetCount{
					ID:    label.String,
					Label: label.String,
					Count: count,
				})
			}
		case "favorite":
			if enumValue.Valid {
				result.Favorite = append(result.Favorite, models.BooleanFacetCount{
					Value: enumValue.String == "true",
					Count: count,
				})
			}
		case "circumcised":
			if enumValue.Valid {
				circumcised := models.CircumisedEnum(enumValue.String)
				if circumcised.IsValid() {
					result.Circumcised = append(result.Circumcised, models.CircumcisedFacetCount{
						Value: circumcised,
						Count: count,
					})
				}
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
