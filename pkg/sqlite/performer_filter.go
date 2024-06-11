package sqlite

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type performerFilterHandler struct {
	performerFilter *models.PerformerFilterType
}

func (qb *performerFilterHandler) validate() error {
	filter := qb.performerFilter
	if filter == nil {
		return nil
	}

	if err := validateFilterCombination(filter.OperatorFilter); err != nil {
		return err
	}

	if subFilter := filter.SubFilter(); subFilter != nil {
		sqb := &performerFilterHandler{performerFilter: subFilter}
		if err := sqb.validate(); err != nil {
			return err
		}
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

func (qb *performerFilterHandler) handle(ctx context.Context, f *filterBuilder) {
	filter := qb.performerFilter
	if filter == nil {
		return
	}

	if err := qb.validate(); err != nil {
		f.setError(err)
		return
	}

	sf := filter.SubFilter()
	if sf != nil {
		sub := &performerFilterHandler{sf}
		handleSubFilter(ctx, sub, f, filter.OperatorFilter)
	}

	f.handleCriterion(ctx, qb.criterionHandler())
}

func (qb *performerFilterHandler) criterionHandler() criterionHandler {
	filter := qb.performerFilter
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
		qb.urlsCriterionHandler(filter.URL),
		intCriterionHandler(filter.Weight, tableName+".weight", nil),
		criterionHandlerFunc(func(ctx context.Context, f *filterBuilder) {
			if filter.StashID != nil {
				performerRepository.stashIDs.join(f, "performer_stash_ids", "performers.id")
				stringCriterionHandler(filter.StashID, "performer_stash_ids.stash_id")(ctx, f)
			}
		}),
		&stashIDCriterionHandler{
			c:                 filter.StashIDEndpoint,
			stashIDRepository: &performerRepository.stashIDs,
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
		&dateCriterionHandler{filter.Birthdate, tableName + ".birthdate", nil},
		&dateCriterionHandler{filter.DeathDate, tableName + ".death_date", nil},
		&timestampCriterionHandler{filter.CreatedAt, tableName + ".created_at", nil},
		&timestampCriterionHandler{filter.UpdatedAt, tableName + ".updated_at", nil},

		&relatedFilterHandler{
			relatedIDCol:   "performers_scenes.scene_id",
			relatedRepo:    sceneRepository.repository,
			relatedHandler: &sceneFilterHandler{filter.ScenesFilter},
			joinFn: func(f *filterBuilder) {
				performerRepository.scenes.innerJoin(f, "", "performers.id")
			},
		},

		&relatedFilterHandler{
			relatedIDCol:   "performers_images.image_id",
			relatedRepo:    imageRepository.repository,
			relatedHandler: &imageFilterHandler{filter.ImagesFilter},
			joinFn: func(f *filterBuilder) {
				performerRepository.images.innerJoin(f, "", "performers.id")
			},
		},

		&relatedFilterHandler{
			relatedIDCol:   "performers_galleries.gallery_id",
			relatedRepo:    galleryRepository.repository,
			relatedHandler: &galleryFilterHandler{filter.GalleriesFilter},
			joinFn: func(f *filterBuilder) {
				performerRepository.galleries.innerJoin(f, "", "performers.id")
			},
		},

		&relatedFilterHandler{
			relatedIDCol:   "performer_tag.tag_id",
			relatedRepo:    tagRepository.repository,
			relatedHandler: &tagFilterHandler{filter.TagsFilter},
			joinFn: func(f *filterBuilder) {
				performerRepository.tags.innerJoin(f, "performer_tag", "performers.id")
			},
		},
	}
}

// TODO - we need to provide a whitelist of possible values
func (qb *performerFilterHandler) performerIsMissingCriterionHandler(isMissing *string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if isMissing != nil && *isMissing != "" {
			switch *isMissing {
			case "url":
				performersURLsTableMgr.join(f, "", "performers.id")
				f.addWhere("performer_urls.url IS NULL")
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

func (qb *performerFilterHandler) performerAgeFilterCriterionHandler(age *models.IntCriterionInput) criterionHandlerFunc {
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

func (qb *performerFilterHandler) urlsCriterionHandler(url *models.StringCriterionInput) criterionHandlerFunc {
	h := stringListCriterionHandlerBuilder{
		primaryTable: performerTable,
		primaryFK:    performerIDColumn,
		joinTable:    performerURLsTable,
		stringColumn: performerURLColumn,
		addJoinTable: func(f *filterBuilder) {
			performersURLsTableMgr.join(f, "", "performers.id")
		},
	}

	return h.handler(url)
}

func (qb *performerFilterHandler) aliasCriterionHandler(alias *models.StringCriterionInput) criterionHandlerFunc {
	h := stringListCriterionHandlerBuilder{
		primaryTable: performerTable,
		primaryFK:    performerIDColumn,
		joinTable:    performersAliasesTable,
		stringColumn: performerAliasColumn,
		addJoinTable: func(f *filterBuilder) {
			performersAliasesTableMgr.join(f, "", "performers.id")
		},
	}

	return h.handler(alias)
}

func (qb *performerFilterHandler) tagsCriterionHandler(tags *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	h := joinedHierarchicalMultiCriterionHandlerBuilder{
		primaryTable: performerTable,
		foreignTable: tagTable,
		foreignFK:    "tag_id",

		relationsTable: "tags_relations",
		joinAs:         "performer_tag",
		joinTable:      performersTagsTable,
		primaryFK:      performerIDColumn,
	}

	return h.handler(tags)
}

func (qb *performerFilterHandler) tagCountCriterionHandler(count *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: performerTable,
		joinTable:    performersTagsTable,
		primaryFK:    performerIDColumn,
	}

	return h.handler(count)
}

func (qb *performerFilterHandler) sceneCountCriterionHandler(count *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: performerTable,
		joinTable:    performersScenesTable,
		primaryFK:    performerIDColumn,
	}

	return h.handler(count)
}

func (qb *performerFilterHandler) imageCountCriterionHandler(count *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: performerTable,
		joinTable:    performersImagesTable,
		primaryFK:    performerIDColumn,
	}

	return h.handler(count)
}

func (qb *performerFilterHandler) galleryCountCriterionHandler(count *models.IntCriterionInput) criterionHandlerFunc {
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

func (qb *performerFilterHandler) oCounterCriterionHandler(count *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if count == nil {
			return
		}

		lhs := "(" + selectPerformerOCountSQL + ")"
		clause, args := getIntCriterionWhereClause(lhs, *count)

		f.addWhere(clause, args...)
	}
}

func (qb *performerFilterHandler) playCounterCriterionHandler(count *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if count == nil {
			return
		}

		lhs := "(" + selectPerformerPlayCountSQL + ")"
		clause, args := getIntCriterionWhereClause(lhs, *count)

		f.addWhere(clause, args...)
	}
}

func (qb *performerFilterHandler) studiosCriterionHandler(studios *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
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
			valuesClause, err := getHierarchicalValues(ctx, studios.Value, studioTable, "", "parent_id", "child_id", studios.Depth)
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

func (qb *performerFilterHandler) appearsWithCriterionHandler(performers *models.MultiCriterionInput) criterionHandlerFunc {
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
