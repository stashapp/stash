package sqlite

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

func (qb *PerformerStore) validateFilter(filter *models.PerformerFilterType) error {
	if filter == nil {
		return nil
	}

	const and = "AND"
	const or = "OR"
	const not = "NOT"

	if filter.And != nil {
		if filter.Or != nil {
			return illegalFilterCombination(and, or)
		}
		if filter.Not != nil {
			return illegalFilterCombination(and, not)
		}

		return qb.validateFilter(filter.And)
	}

	if filter.Or != nil {
		if filter.Not != nil {
			return illegalFilterCombination(or, not)
		}

		return qb.validateFilter(filter.Or)
	}

	if filter.Not != nil {
		return qb.validateFilter(filter.Not)
	}

	// if legacy height filter used, ensure only supported modifiers are used
	if filter.Height != nil {
		// treat as an int filter
		intCrit := &models.IntCriterionInput{
			Modifier: filter.Height.Modifier,
		}
		if !intCrit.ValidModifier() {
			return fmt.Errorf("invalid height modifier: %s", filter.Height.Modifier)
		}

		// ensure value is a valid number
		if _, err := strconv.Atoi(filter.Height.Value); err != nil {
			return fmt.Errorf("invalid height value: %s", filter.Height.Value)
		}
	}

	return nil
}

func (qb *PerformerStore) makeFilter(ctx context.Context, filter *models.PerformerFilterType) *filterBuilder {
	if filter == nil {
		return nil
	}

	query := &filterBuilder{}

	if filter.And != nil {
		query.and(qb.makeFilter(ctx, filter.And))
	}
	if filter.Or != nil {
		query.or(qb.makeFilter(ctx, filter.Or))
	}
	if filter.Not != nil {
		query.not(qb.makeFilter(ctx, filter.Not))
	}

	query.handleCriterion(ctx, qb.criterionHandler(filter))

	return query
}

func (qb *PerformerStore) criterionHandler(filter *models.PerformerFilterType) criterionHandler {
	const tableName = performerTable
	heightCmCrit := filter.HeightCm

	return compoundHandler{
		stringCriterionHandler(filter.Name, tableName+".name"),
		stringCriterionHandler(filter.Disambiguation, tableName+".disambiguation"),
		stringCriterionHandler(filter.Details, tableName+".details"),

		boolCriterionHandler(filter.FilterFavorites, tableName+".favorite", nil),
		boolCriterionHandler(filter.IgnoreAutoTag, tableName+".ignore_auto_tag", nil),

		yearFilterCriterionHandler(filter.BirthYear, tableName+".birthdate"),
		yearFilterCriterionHandler(filter.DeathYear, tableName+".death_date"),

		qb.performerAgeFilterCriterionHandler(filter.Age),

		criterionHandlerFunc(func(ctx context.Context, f *filterBuilder) {
			if gender := filter.Gender; gender != nil {
				genderCopy := *gender
				if genderCopy.Value.IsValid() && len(genderCopy.ValueList) == 0 {
					genderCopy.ValueList = []models.GenderEnum{genderCopy.Value}
				}

				v := utils.StringerSliceToStringSlice(genderCopy.ValueList)
				enumCriterionHandler(genderCopy.Modifier, v, tableName+".gender")(ctx, f)
			}
		}),

		qb.performerIsMissingCriterionHandler(filter.IsMissing),
		stringCriterionHandler(filter.Ethnicity, tableName+".ethnicity"),
		stringCriterionHandler(filter.Country, tableName+".country"),
		stringCriterionHandler(filter.EyeColor, tableName+".eye_color"),

		// special handler for legacy height filter
		criterionHandlerFunc(func(ctx context.Context, f *filterBuilder) {
			if heightCmCrit == nil && filter.Height != nil {
				heightCm, _ := strconv.Atoi(filter.Height.Value) // already validated
				heightCmCrit = &models.IntCriterionInput{
					Value:    heightCm,
					Modifier: filter.Height.Modifier,
				}
			}
		}),

		intCriterionHandler(heightCmCrit, tableName+".height", nil),

		stringCriterionHandler(filter.Measurements, tableName+".measurements"),
		stringCriterionHandler(filter.FakeTits, tableName+".fake_tits"),
		floatCriterionHandler(filter.PenisLength, tableName+".penis_length", nil),

		criterionHandlerFunc(func(ctx context.Context, f *filterBuilder) {
			if circumcised := filter.Circumcised; circumcised != nil {
				v := utils.StringerSliceToStringSlice(circumcised.Value)
				enumCriterionHandler(circumcised.Modifier, v, tableName+".circumcised")(ctx, f)
			}
		}),

		stringCriterionHandler(filter.CareerLength, tableName+".career_length"),
		stringCriterionHandler(filter.Tattoos, tableName+".tattoos"),
		stringCriterionHandler(filter.Piercings, tableName+".piercings"),
		intCriterionHandler(filter.Rating100, tableName+".rating", nil),
		stringCriterionHandler(filter.HairColor, tableName+".hair_color"),
		stringCriterionHandler(filter.URL, tableName+".url"),
		intCriterionHandler(filter.Weight, tableName+".weight", nil),
		criterionHandlerFunc(func(ctx context.Context, f *filterBuilder) {
			if filter.StashID != nil {
				qb.stashIDRepository().join(f, "performer_stash_ids", "performers.id")
				stringCriterionHandler(filter.StashID, "performer_stash_ids.stash_id")(ctx, f)
			}
		}),
		&stashIDCriterionHandler{
			c:                 filter.StashIDEndpoint,
			stashIDRepository: qb.stashIDRepository(),
			stashIDTableAs:    "performer_stash_ids",
			parentIDCol:       "performers.id",
		},

		qb.aliasCriterionHandler(filter.Aliases),

		qb.tagsCriterionHandler(filter.Tags),

		qb.studiosCriterionHandler(filter.Studios),

		qb.appearsWithCriterionHandler(filter.Performers),

		qb.tagCountCriterionHandler(filter.TagCount),
		qb.sceneCountCriterionHandler(filter.SceneCount),
		qb.imageCountCriterionHandler(filter.ImageCount),
		qb.galleryCountCriterionHandler(filter.GalleryCount),
		qb.playCounterCriterionHandler(filter.PlayCount),
		qb.oCounterCriterionHandler(filter.OCounter),
		dateCriterionHandler(filter.Birthdate, tableName+".birthdate"),
		dateCriterionHandler(filter.DeathDate, tableName+".death_date"),
		timestampCriterionHandler(filter.CreatedAt, tableName+".created_at"),
		timestampCriterionHandler(filter.UpdatedAt, tableName+".updated_at"),
	}
}

// TODO - we need to provide a whitelist of possible values
func (qb *PerformerStore) performerIsMissingCriterionHandler(isMissing *string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if isMissing != nil && *isMissing != "" {
			switch *isMissing {
			case "scenes": // Deprecated: use `scene_count == 0` filter instead
				f.addLeftJoin(performersScenesTable, "scenes_join", "scenes_join.performer_id = performers.id")
				f.addWhere("scenes_join.scene_id IS NULL")
			case "image":
				f.addWhere("performers.image_blob IS NULL")
			case "stash_id":
				performersStashIDsTableMgr.join(f, "performer_stash_ids", "performers.id")
				f.addWhere("performer_stash_ids.performer_id IS NULL")
			case "aliases":
				performersAliasesTableMgr.join(f, "", "performers.id")
				f.addWhere("performer_aliases.alias IS NULL")
			default:
				f.addWhere("(performers." + *isMissing + " IS NULL OR TRIM(performers." + *isMissing + ") = '')")
			}
		}
	}
}

func (qb *PerformerStore) performerAgeFilterCriterionHandler(age *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if age != nil && age.Modifier.IsValid() {
			clause, args := getIntCriterionWhereClause(
				"cast(IFNULL(strftime('%Y.%m%d', performers.death_date), strftime('%Y.%m%d', 'now')) - strftime('%Y.%m%d', performers.birthdate) as int)",
				*age,
			)
			f.addWhere(clause, args...)
		}
	}
}

func (qb *PerformerStore) aliasCriterionHandler(alias *models.StringCriterionInput) criterionHandlerFunc {
	h := stringListCriterionHandlerBuilder{
		joinTable:    performersAliasesTable,
		stringColumn: performerAliasColumn,
		addJoinTable: func(f *filterBuilder) {
			performersAliasesTableMgr.join(f, "", "performers.id")
		},
	}

	return h.handler(alias)
}

func (qb *PerformerStore) tagsCriterionHandler(tags *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	h := joinedHierarchicalMultiCriterionHandlerBuilder{
		tx: qb.tx,

		primaryTable: performerTable,
		foreignTable: tagTable,
		foreignFK:    "tag_id",

		relationsTable: "tags_relations",
		joinAs:         "image_tag",
		joinTable:      performersTagsTable,
		primaryFK:      performerIDColumn,
	}

	return h.handler(tags)
}

func (qb *PerformerStore) tagCountCriterionHandler(count *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: performerTable,
		joinTable:    performersTagsTable,
		primaryFK:    performerIDColumn,
	}

	return h.handler(count)
}

func (qb *PerformerStore) sceneCountCriterionHandler(count *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: performerTable,
		joinTable:    performersScenesTable,
		primaryFK:    performerIDColumn,
	}

	return h.handler(count)
}

func (qb *PerformerStore) imageCountCriterionHandler(count *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: performerTable,
		joinTable:    performersImagesTable,
		primaryFK:    performerIDColumn,
	}

	return h.handler(count)
}

func (qb *PerformerStore) galleryCountCriterionHandler(count *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: performerTable,
		joinTable:    performersGalleriesTable,
		primaryFK:    performerIDColumn,
	}

	return h.handler(count)
}

// used for sorting and filtering on performer o-count
var selectPerformerOCountSQL = utils.StrFormat(
	"SELECT SUM(o_counter) "+
		"FROM ("+
		"SELECT SUM(o_counter) as o_counter from {performers_images} s "+
		"LEFT JOIN {images} ON {images}.id = s.{images_id} "+
		"WHERE s.{performer_id} = {performers}.id "+
		"UNION ALL "+
		"SELECT COUNT({scenes_o_dates}.{o_date}) as o_counter from {performers_scenes} s "+
		"LEFT JOIN {scenes} ON {scenes}.id = s.{scene_id} "+
		"LEFT JOIN {scenes_o_dates} ON {scenes_o_dates}.{scene_id} = {scenes}.id "+
		"WHERE s.{performer_id} = {performers}.id "+
		")",
	map[string]interface{}{
		"performers_images": performersImagesTable,
		"images":            imageTable,
		"performer_id":      performerIDColumn,
		"images_id":         imageIDColumn,
		"performers":        performerTable,
		"performers_scenes": performersScenesTable,
		"scenes":            sceneTable,
		"scene_id":          sceneIDColumn,
		"scenes_o_dates":    scenesODatesTable,
		"o_date":            sceneODateColumn,
	},
)

// used for sorting and filtering play count on performer view count
var selectPerformerPlayCountSQL = utils.StrFormat(
	"SELECT COUNT(DISTINCT {view_date}) FROM ("+
		"SELECT {view_date} FROM {performers_scenes} s "+
		"LEFT JOIN {scenes} ON {scenes}.id = s.{scene_id} "+
		"LEFT JOIN {scenes_view_dates} ON {scenes_view_dates}.{scene_id} = {scenes}.id "+
		"WHERE s.{performer_id} = {performers}.id"+
		")",
	map[string]interface{}{
		"performer_id":      performerIDColumn,
		"performers":        performerTable,
		"performers_scenes": performersScenesTable,
		"scenes":            sceneTable,
		"scene_id":          sceneIDColumn,
		"scenes_view_dates": scenesViewDatesTable,
		"view_date":         sceneViewDateColumn,
	},
)

func (qb *PerformerStore) oCounterCriterionHandler(count *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if count == nil {
			return
		}

		lhs := "(" + selectPerformerOCountSQL + ")"
		clause, args := getIntCriterionWhereClause(lhs, *count)

		f.addWhere(clause, args...)
	}
}

func (qb *PerformerStore) playCounterCriterionHandler(count *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if count == nil {
			return
		}

		lhs := "(" + selectPerformerPlayCountSQL + ")"
		clause, args := getIntCriterionWhereClause(lhs, *count)

		f.addWhere(clause, args...)
	}
}

func (qb *PerformerStore) studiosCriterionHandler(studios *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if studios != nil {
			formatMaps := []utils.StrFormatMap{
				{
					"primaryTable": sceneTable,
					"joinTable":    performersScenesTable,
					"primaryFK":    sceneIDColumn,
				},
				{
					"primaryTable": imageTable,
					"joinTable":    performersImagesTable,
					"primaryFK":    imageIDColumn,
				},
				{
					"primaryTable": galleryTable,
					"joinTable":    performersGalleriesTable,
					"primaryFK":    galleryIDColumn,
				},
			}

			if studios.Modifier == models.CriterionModifierIsNull || studios.Modifier == models.CriterionModifierNotNull {
				var notClause string
				if studios.Modifier == models.CriterionModifierNotNull {
					notClause = "NOT"
				}

				var conditions []string
				for _, c := range formatMaps {
					f.addLeftJoin(c["joinTable"].(string), "", fmt.Sprintf("%s.performer_id = performers.id", c["joinTable"]))
					f.addLeftJoin(c["primaryTable"].(string), "", fmt.Sprintf("%s.%s = %s.id", c["joinTable"], c["primaryFK"], c["primaryTable"]))

					conditions = append(conditions, fmt.Sprintf("%s.studio_id IS NULL", c["primaryTable"]))
				}

				f.addWhere(fmt.Sprintf("%s (%s)", notClause, strings.Join(conditions, " AND ")))
				return
			}

			if len(studios.Value) == 0 {
				return
			}

			var clauseCondition string

			switch studios.Modifier {
			case models.CriterionModifierIncludes:
				// return performers who appear in scenes/images/galleries with any of the given studios
				clauseCondition = "NOT"
			case models.CriterionModifierExcludes:
				// exclude performers who appear in scenes/images/galleries with any of the given studios
				clauseCondition = ""
			default:
				return
			}

			const derivedPerformerStudioTable = "performer_studio"
			valuesClause, err := getHierarchicalValues(ctx, qb.tx, studios.Value, studioTable, "", "parent_id", "child_id", studios.Depth)
			if err != nil {
				f.setError(err)
				return
			}
			f.addWith("studio(root_id, item_id) AS (" + valuesClause + ")")

			templStr := `SELECT performer_id FROM {primaryTable}
	INNER JOIN {joinTable} ON {primaryTable}.id = {joinTable}.{primaryFK}
	INNER JOIN studio ON {primaryTable}.studio_id = studio.item_id`

			var unions []string
			for _, c := range formatMaps {
				unions = append(unions, utils.StrFormat(templStr, c))
			}

			f.addWith(fmt.Sprintf("%s AS (%s)", derivedPerformerStudioTable, strings.Join(unions, " UNION ")))

			f.addLeftJoin(derivedPerformerStudioTable, "", fmt.Sprintf("performers.id = %s.performer_id", derivedPerformerStudioTable))
			f.addWhere(fmt.Sprintf("%s.performer_id IS %s NULL", derivedPerformerStudioTable, clauseCondition))
		}
	}
}

func (qb *PerformerStore) appearsWithCriterionHandler(performers *models.MultiCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if performers != nil {
			formatMaps := []utils.StrFormatMap{
				{
					"primaryTable": performersScenesTable,
					"joinTable":    performersScenesTable,
					"primaryFK":    sceneIDColumn,
				},
				{
					"primaryTable": performersImagesTable,
					"joinTable":    performersImagesTable,
					"primaryFK":    imageIDColumn,
				},
				{
					"primaryTable": performersGalleriesTable,
					"joinTable":    performersGalleriesTable,
					"primaryFK":    galleryIDColumn,
				},
			}

			if len(performers.Value) == '0' {
				return
			}

			const derivedPerformerPerformersTable = "performer_performers"

			valuesClause := strings.Join(performers.Value, "),(")

			f.addWith("performer(id) AS (VALUES(" + valuesClause + "))")

			templStr := `SELECT {primaryTable}2.performer_id FROM {primaryTable}
			INNER JOIN {primaryTable} AS {primaryTable}2 ON {primaryTable}.{primaryFK} = {primaryTable}2.{primaryFK}
			INNER JOIN performer ON {primaryTable}.performer_id = performer.id
			WHERE {primaryTable}2.performer_id != performer.id`

			if performers.Modifier == models.CriterionModifierIncludesAll && len(performers.Value) > 1 {
				templStr += `
							GROUP BY {primaryTable}2.performer_id
							HAVING(count(distinct {primaryTable}.performer_id) IS ` + strconv.Itoa(len(performers.Value)) + `)`
			}

			var unions []string
			for _, c := range formatMaps {
				unions = append(unions, utils.StrFormat(templStr, c))
			}

			f.addWith(fmt.Sprintf("%s AS (%s)", derivedPerformerPerformersTable, strings.Join(unions, " UNION ")))

			f.addInnerJoin(derivedPerformerPerformersTable, "", fmt.Sprintf("performers.id = %s.performer_id", derivedPerformerPerformersTable))
		}
	}
}
