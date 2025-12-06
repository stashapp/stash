package sqlite

import (
	"context"
	stdsql "database/sql"
	"fmt"
	"strconv"
	"sync"

	"github.com/stashapp/stash/pkg/models"
)

// GetFacets returns aggregated facet counts for scenes matching the given filter.
// This uses parallel queries for better performance on large databases.
// The options parameter controls which expensive facets to include (lazy loading).
func (qb *SceneStore) GetFacets(ctx context.Context, sceneFilter *models.SceneFilterType, limit int, options models.SceneFacetOptions) (*models.SceneFacets, error) {
	result := &models.SceneFacets{
		Tags:          []models.FacetCount{},
		Performers:    []models.FacetCount{},
		Studios:       []models.FacetCount{},
		Groups:        []models.FacetCount{},
		PerformerTags: []models.FacetCount{},
		Resolutions:   []models.ResolutionFacetCount{},
		Orientations:  []models.OrientationFacetCount{},
		Organized:     []models.BooleanFacetCount{},
		Interactive:   []models.BooleanFacetCount{},
		Ratings:       []models.RatingFacetCount{},
		Captions:      []models.CaptionFacetCount{},
	}

	query, err := qb.makeQuery(ctx, sceneFilter, nil)
	if err != nil {
		return nil, err
	}

	baseSQL := query.toSQL(false)
	baseArgs := append([]interface{}{}, query.args...)

	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error

	// Run core facets (fast) in main query
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := qb.getCoreFacets(ctx, baseSQL, baseArgs, limit, result, &mu); err != nil {
			mu.Lock()
			if firstErr == nil {
				firstErr = err
			}
			mu.Unlock()
		}
	}()

	// Run performer_tags facet in parallel only if requested (expensive - 3 joins)
	if options.IncludePerformerTags {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := qb.getPerformerTagsFacet(ctx, baseSQL, baseArgs, limit, result, &mu); err != nil {
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
			}
		}()
	}

	// Run captions facet in parallel only if requested (expensive - file joins)
	if options.IncludeCaptions {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := qb.getCaptionsFacet(ctx, baseSQL, baseArgs, result, &mu); err != nil {
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}

	return result, nil
}

// getCoreFacets fetches the main facets that are relatively fast
func (qb *SceneStore) getCoreFacets(ctx context.Context, baseSQL string, baseArgs []interface{}, limit int, result *models.SceneFacets, mu *sync.Mutex) error {
	args := append([]interface{}{}, baseArgs...)

	sql := fmt.Sprintf(`
		WITH filtered_scenes AS (%s)
		
		SELECT * FROM (
			SELECT 'tag' as facet_type, t.id, t.name as label, NULL as enum_value, COUNT(DISTINCT st.scene_id) as count
			FROM filtered_scenes fs
			INNER JOIN scenes_tags st ON fs.id = st.scene_id
			INNER JOIN tags t ON st.tag_id = t.id
			GROUP BY t.id
			ORDER BY count DESC
			LIMIT ?
		)
		
		UNION ALL
		
		SELECT * FROM (
			SELECT 'performer' as facet_type, p.id, p.name as label, NULL as enum_value, COUNT(DISTINCT ps.scene_id) as count
			FROM filtered_scenes fs
			INNER JOIN performers_scenes ps ON fs.id = ps.scene_id
			INNER JOIN performers p ON ps.performer_id = p.id
			GROUP BY p.id
			ORDER BY count DESC
			LIMIT ?
		)
		
		UNION ALL
		
		SELECT * FROM (
			SELECT 'studio' as facet_type, s.id, s.name as label, NULL as enum_value, COUNT(DISTINCT sc.id) as count
			FROM filtered_scenes fs
			INNER JOIN scenes sc ON fs.id = sc.id
			INNER JOIN studios s ON sc.studio_id = s.id
			WHERE sc.studio_id IS NOT NULL
			GROUP BY s.id
			ORDER BY count DESC
			LIMIT ?
		)
		
		UNION ALL
		
		SELECT * FROM (
			SELECT 'group' as facet_type, g.id, g.name as label, NULL as enum_value, COUNT(DISTINCT gs.scene_id) as count
			FROM filtered_scenes fs
			INNER JOIN groups_scenes gs ON fs.id = gs.scene_id
			INNER JOIN groups g ON gs.group_id = g.id
			GROUP BY g.id
			ORDER BY count DESC
			LIMIT ?
		)
		
		UNION ALL
		
		SELECT * FROM (
			SELECT 'resolution' as facet_type, 0 as id, '' as label, 
				CASE 
					WHEN vf.height >= 144 AND vf.height < 240 THEN 'VERY_LOW'
					WHEN vf.height >= 240 AND vf.height < 360 THEN 'LOW'
					WHEN vf.height >= 360 AND vf.height < 480 THEN 'R360P'
					WHEN vf.height >= 480 AND vf.height < 540 THEN 'STANDARD'
					WHEN vf.height >= 540 AND vf.height < 720 THEN 'WEB_HD'
					WHEN vf.height >= 720 AND vf.height < 1080 THEN 'STANDARD_HD'
					WHEN vf.height >= 1080 AND vf.height < 1440 THEN 'FULL_HD'
					WHEN vf.height >= 1440 AND vf.height < 1920 THEN 'QUAD_HD'
					WHEN vf.height >= 1920 AND vf.height < 2160 THEN 'VR_HD'
					WHEN vf.height >= 2160 AND vf.height < 2560 THEN 'FOUR_K'
					WHEN vf.height >= 2560 AND vf.height < 3000 THEN 'FIVE_K'
					WHEN vf.height >= 3000 AND vf.height < 3584 THEN 'SIX_K'
					WHEN vf.height >= 3584 AND vf.height < 3840 THEN 'SEVEN_K'
					WHEN vf.height >= 3840 AND vf.height < 6144 THEN 'EIGHT_K'
					WHEN vf.height >= 6144 THEN 'HUGE'
					ELSE 'UNKNOWN'
				END as enum_value,
				COUNT(DISTINCT sf.scene_id) as count
			FROM filtered_scenes fs
			INNER JOIN scenes_files sf ON fs.id = sf.scene_id AND sf."primary" = 1
			INNER JOIN video_files vf ON sf.file_id = vf.file_id
			WHERE vf.height IS NOT NULL
			GROUP BY enum_value
			HAVING enum_value != 'UNKNOWN'
		)
		
		UNION ALL
		
		SELECT * FROM (
			SELECT 'orientation' as facet_type, 0 as id, '' as label,
				CASE 
					WHEN vf.width > vf.height THEN 'LANDSCAPE'
					WHEN vf.width < vf.height THEN 'PORTRAIT'
					ELSE 'SQUARE'
				END as enum_value,
				COUNT(DISTINCT sf.scene_id) as count
			FROM filtered_scenes fs
			INNER JOIN scenes_files sf ON fs.id = sf.scene_id AND sf."primary" = 1
			INNER JOIN video_files vf ON sf.file_id = vf.file_id
			GROUP BY enum_value
		)
		
		UNION ALL
		
		SELECT * FROM (
			SELECT 'organized' as facet_type, 0 as id, '' as label,
				CASE WHEN sc.organized = 1 THEN 'true' ELSE 'false' END as enum_value,
				COUNT(*) as count
			FROM filtered_scenes fs
			INNER JOIN scenes sc ON fs.id = sc.id
			GROUP BY sc.organized
		)
		
		UNION ALL
		
		SELECT * FROM (
			SELECT 'interactive' as facet_type, 0 as id, '' as label,
				CASE WHEN vf.interactive = 1 THEN 'true' ELSE 'false' END as enum_value,
				COUNT(DISTINCT sf.scene_id) as count
			FROM filtered_scenes fs
			INNER JOIN scenes_files sf ON fs.id = sf.scene_id AND sf."primary" = 1
			INNER JOIN video_files vf ON sf.file_id = vf.file_id
			GROUP BY vf.interactive
		)
		
		UNION ALL
		
		SELECT * FROM (
			SELECT 'rating' as facet_type, 0 as id, '' as label,
				CAST(sc.rating AS TEXT) as enum_value,
				COUNT(*) as count
			FROM filtered_scenes fs
			INNER JOIN scenes sc ON fs.id = sc.id
			WHERE sc.rating IS NOT NULL
			GROUP BY sc.rating
			ORDER BY sc.rating DESC
		)
	`, baseSQL)

	args = append(args, limit, limit, limit, limit)

	rows, err := dbWrapper.Queryx(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("error executing core facets query: %w", err)
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
			return fmt.Errorf("error scanning facet row: %w", err)
		}

		mu.Lock()
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
		case "group":
			result.Groups = append(result.Groups, models.FacetCount{
				ID:    strconv.Itoa(id),
				Label: label.String,
				Count: count,
			})
		case "resolution":
			if enumValue.Valid {
				res := models.ResolutionEnum(enumValue.String)
				if res.IsValid() {
					result.Resolutions = append(result.Resolutions, models.ResolutionFacetCount{
						Resolution: res,
						Count:      count,
					})
				}
			}
		case "orientation":
			if enumValue.Valid {
				result.Orientations = append(result.Orientations, models.OrientationFacetCount{
					Orientation: models.OrientationEnum(enumValue.String),
					Count:       count,
				})
			}
		case "organized":
			if enumValue.Valid {
				result.Organized = append(result.Organized, models.BooleanFacetCount{
					Value: enumValue.String == "true",
					Count: count,
				})
			}
		case "interactive":
			if enumValue.Valid {
				result.Interactive = append(result.Interactive, models.BooleanFacetCount{
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
		mu.Unlock()
	}

	return rows.Err()
}

// getPerformerTagsFacet fetches performer tags facet (expensive - 3 joins)
func (qb *SceneStore) getPerformerTagsFacet(ctx context.Context, baseSQL string, baseArgs []interface{}, limit int, result *models.SceneFacets, mu *sync.Mutex) error {
	args := append([]interface{}{}, baseArgs...)

	sql := fmt.Sprintf(`
		WITH filtered_scenes AS (%s)
		SELECT t.id, t.name as label, COUNT(DISTINCT fs.id) as count
		FROM filtered_scenes fs
		INNER JOIN performers_scenes ps ON fs.id = ps.scene_id
		INNER JOIN performers_tags pt ON ps.performer_id = pt.performer_id
		INNER JOIN tags t ON pt.tag_id = t.id
		GROUP BY t.id
		ORDER BY count DESC
		LIMIT ?
	`, baseSQL)

	args = append(args, limit)

	rows, err := dbWrapper.Queryx(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("error executing performer_tags facet query: %w", err)
	}
	defer rows.Close()

	var performerTags []models.FacetCount
	for rows.Next() {
		var id int
		var label string
		var count int

		if err := rows.Scan(&id, &label, &count); err != nil {
			return fmt.Errorf("error scanning performer_tag row: %w", err)
		}

		performerTags = append(performerTags, models.FacetCount{
			ID:    strconv.Itoa(id),
			Label: label,
			Count: count,
		})
	}

	mu.Lock()
	result.PerformerTags = performerTags
	mu.Unlock()

	return rows.Err()
}

// getCaptionsFacet fetches captions facet (expensive - file joins)
func (qb *SceneStore) getCaptionsFacet(ctx context.Context, baseSQL string, baseArgs []interface{}, result *models.SceneFacets, mu *sync.Mutex) error {
	args := append([]interface{}{}, baseArgs...)

	sql := fmt.Sprintf(`
		WITH filtered_scenes AS (%s)
		SELECT vc.language_code, COUNT(DISTINCT sf.scene_id) as count
		FROM filtered_scenes fs
		INNER JOIN scenes_files sf ON fs.id = sf.scene_id AND sf."primary" = 1
		INNER JOIN video_captions vc ON sf.file_id = vc.file_id
		WHERE vc.language_code IS NOT NULL AND vc.language_code != ''
		GROUP BY vc.language_code
		ORDER BY count DESC
	`, baseSQL)

	rows, err := dbWrapper.Queryx(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("error executing captions facet query: %w", err)
	}
	defer rows.Close()

	var captions []models.CaptionFacetCount
	for rows.Next() {
		var languageCode string
		var count int

		if err := rows.Scan(&languageCode, &count); err != nil {
			return fmt.Errorf("error scanning caption row: %w", err)
		}

		captions = append(captions, models.CaptionFacetCount{
			Language: languageCode,
			Count:    count,
		})
	}

	mu.Lock()
	result.Captions = captions
	mu.Unlock()

	return rows.Err()
}
