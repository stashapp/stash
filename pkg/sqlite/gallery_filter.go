package sqlite

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/stashapp/stash/pkg/models"
)

func (qb *GalleryStore) validateFilter(galleryFilter *models.GalleryFilterType) error {
	if galleryFilter == nil {
		return nil
	}

	const and = "AND"
	const or = "OR"
	const not = "NOT"

	if galleryFilter.And != nil {
		if galleryFilter.Or != nil {
			return illegalFilterCombination(and, or)
		}
		if galleryFilter.Not != nil {
			return illegalFilterCombination(and, not)
		}

		return qb.validateFilter(galleryFilter.And)
	}

	if galleryFilter.Or != nil {
		if galleryFilter.Not != nil {
			return illegalFilterCombination(or, not)
		}

		return qb.validateFilter(galleryFilter.Or)
	}

	if galleryFilter.Not != nil {
		return qb.validateFilter(galleryFilter.Not)
	}

	return nil
}

func (qb *GalleryStore) makeFilter(ctx context.Context, galleryFilter *models.GalleryFilterType) *filterBuilder {
	if galleryFilter == nil {
		return nil
	}

	query := &filterBuilder{}

	if galleryFilter.And != nil {
		query.and(qb.makeFilter(ctx, galleryFilter.And))
	}
	if galleryFilter.Or != nil {
		query.or(qb.makeFilter(ctx, galleryFilter.Or))
	}
	if galleryFilter.Not != nil {
		query.not(qb.makeFilter(ctx, galleryFilter.Not))
	}

	query.handleCriterion(ctx, qb.criterionHandler(galleryFilter))

	return query
}

func (qb *GalleryStore) criterionHandler(filter *models.GalleryFilterType) criterionHandler {
	return compoundHandler{
		intCriterionHandler(filter.ID, "galleries.id", nil),
		stringCriterionHandler(filter.Title, "galleries.title"),
		stringCriterionHandler(filter.Code, "galleries.code"),
		stringCriterionHandler(filter.Details, "galleries.details"),
		stringCriterionHandler(filter.Photographer, "galleries.photographer"),

		criterionHandlerFunc(func(ctx context.Context, f *filterBuilder) {
			if filter.Checksum != nil {
				qb.addGalleriesFilesTable(f)
				f.addLeftJoin(fingerprintTable, "fingerprints_md5", "galleries_files.file_id = fingerprints_md5.file_id AND fingerprints_md5.type = 'md5'")
			}

			stringCriterionHandler(filter.Checksum, "fingerprints_md5.fingerprint")(ctx, f)
		}),

		criterionHandlerFunc(func(ctx context.Context, f *filterBuilder) {
			if filter.IsZip != nil {
				qb.addGalleriesFilesTable(f)
				if *filter.IsZip {

					f.addWhere("galleries_files.file_id IS NOT NULL")
				} else {
					f.addWhere("galleries_files.file_id IS NULL")
				}
			}
		}),

		qb.pathCriterionHandler(filter.Path),
		qb.fileCountCriterionHandler(filter.FileCount),
		intCriterionHandler(filter.Rating100, "galleries.rating", nil),
		qb.urlsCriterionHandler(filter.URL),
		boolCriterionHandler(filter.Organized, "galleries.organized", nil),
		qb.missingCriterionHandler(filter.IsMissing),
		qb.tagsCriterionHandler(filter.Tags),
		qb.tagCountCriterionHandler(filter.TagCount),
		qb.performersCriterionHandler(filter.Performers),
		qb.performerCountCriterionHandler(filter.PerformerCount),
		qb.scenesCriterionHandler(filter.Scenes),
		qb.hasChaptersCriterionHandler(filter.HasChapters),
		studioCriterionHandler(galleryTable, filter.Studios),
		qb.performerTagsCriterionHandler(filter.PerformerTags),
		qb.averageResolutionCriterionHandler(filter.AverageResolution),
		qb.imageCountCriterionHandler(filter.ImageCount),
		qb.performerFavoriteCriterionHandler(filter.PerformerFavorite),
		qb.performerAgeCriterionHandler(filter.PerformerAge),
		dateCriterionHandler(filter.Date, "galleries.date"),
		timestampCriterionHandler(filter.CreatedAt, "galleries.created_at"),
		timestampCriterionHandler(filter.UpdatedAt, "galleries.updated_at"),
	}
}

func (qb *GalleryStore) urlsCriterionHandler(url *models.StringCriterionInput) criterionHandlerFunc {
	h := stringListCriterionHandlerBuilder{
		joinTable:    galleriesURLsTable,
		stringColumn: galleriesURLColumn,
		addJoinTable: func(f *filterBuilder) {
			galleriesURLsTableMgr.join(f, "", "galleries.id")
		},
	}

	return h.handler(url)
}

func (qb *GalleryStore) getMultiCriterionHandlerBuilder(foreignTable, joinTable, foreignFK string, addJoinsFunc func(f *filterBuilder)) multiCriterionHandlerBuilder {
	return multiCriterionHandlerBuilder{
		primaryTable: galleryTable,
		foreignTable: foreignTable,
		joinTable:    joinTable,
		primaryFK:    galleryIDColumn,
		foreignFK:    foreignFK,
		addJoinsFunc: addJoinsFunc,
	}
}

func (qb *GalleryStore) pathCriterionHandler(c *models.StringCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if c != nil {
			qb.addFoldersTable(f)
			f.addLeftJoin(folderTable, "gallery_folder", "galleries.folder_id = gallery_folder.id")

			const pathColumn = "folders.path"
			const basenameColumn = "files.basename"
			const folderPathColumn = "gallery_folder.path"

			addWildcards := true
			not := false

			if modifier := c.Modifier; c.Modifier.IsValid() {
				switch modifier {
				case models.CriterionModifierIncludes:
					clause := getPathSearchClauseMany(pathColumn, basenameColumn, c.Value, addWildcards, not)
					clause2 := getStringSearchClause([]string{folderPathColumn}, c.Value, false)
					f.whereClauses = append(f.whereClauses, orClauses(clause, clause2))
				case models.CriterionModifierExcludes:
					not = true
					clause := getPathSearchClauseMany(pathColumn, basenameColumn, c.Value, addWildcards, not)
					clause2 := getStringSearchClause([]string{folderPathColumn}, c.Value, true)
					f.whereClauses = append(f.whereClauses, orClauses(clause, clause2))
				case models.CriterionModifierEquals:
					addWildcards = false
					clause := getPathSearchClause(pathColumn, basenameColumn, c.Value, addWildcards, not)
					clause2 := makeClause(folderPathColumn+" LIKE ?", c.Value)
					f.whereClauses = append(f.whereClauses, orClauses(clause, clause2))
				case models.CriterionModifierNotEquals:
					addWildcards = false
					not = true
					clause := getPathSearchClause(pathColumn, basenameColumn, c.Value, addWildcards, not)
					clause2 := makeClause(folderPathColumn+" NOT LIKE ?", c.Value)
					f.whereClauses = append(f.whereClauses, orClauses(clause, clause2))
				case models.CriterionModifierMatchesRegex:
					if _, err := regexp.Compile(c.Value); err != nil {
						f.setError(err)
						return
					}
					filepathColumn := fmt.Sprintf("%s || '%s' || %s", pathColumn, string(filepath.Separator), basenameColumn)
					clause := makeClause(fmt.Sprintf("%s IS NOT NULL AND %s IS NOT NULL AND %s regexp ?", pathColumn, basenameColumn, filepathColumn), c.Value)
					clause2 := makeClause(fmt.Sprintf("%s IS NOT NULL AND %[1]s regexp ?", folderPathColumn), c.Value)
					f.whereClauses = append(f.whereClauses, orClauses(clause, clause2))
				case models.CriterionModifierNotMatchesRegex:
					if _, err := regexp.Compile(c.Value); err != nil {
						f.setError(err)
						return
					}
					filepathColumn := fmt.Sprintf("%s || '%s' || %s", pathColumn, string(filepath.Separator), basenameColumn)
					f.addWhere(fmt.Sprintf("%s IS NULL OR %s IS NULL OR %s NOT regexp ?", pathColumn, basenameColumn, filepathColumn), c.Value)
					f.addWhere(fmt.Sprintf("%s IS NULL OR %[1]s NOT regexp ?", folderPathColumn), c.Value)
				case models.CriterionModifierIsNull:
					f.addWhere(fmt.Sprintf("%s IS NULL OR TRIM(%[1]s) = '' OR %s IS NULL OR TRIM(%[2]s) = ''", pathColumn, basenameColumn))
					f.addWhere(fmt.Sprintf("%s IS NULL OR TRIM(%[1]s) = ''", folderPathColumn))
				case models.CriterionModifierNotNull:
					clause := makeClause(fmt.Sprintf("%s IS NOT NULL AND TRIM(%[1]s) != '' AND %s IS NOT NULL AND TRIM(%[2]s) != ''", pathColumn, basenameColumn))
					clause2 := makeClause(fmt.Sprintf("%s IS NOT NULL AND TRIM(%[1]s) != ''", folderPathColumn))
					f.whereClauses = append(f.whereClauses, orClauses(clause, clause2))
				default:
					panic("unsupported string filter modifier")
				}
			}
		}
	}
}

func (qb *GalleryStore) fileCountCriterionHandler(fileCount *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: galleryTable,
		joinTable:    galleriesFilesTable,
		primaryFK:    galleryIDColumn,
	}

	return h.handler(fileCount)
}

func (qb *GalleryStore) missingCriterionHandler(isMissing *string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if isMissing != nil && *isMissing != "" {
			switch *isMissing {
			case "url":
				galleriesURLsTableMgr.join(f, "", "galleries.id")
				f.addWhere("gallery_urls.url IS NULL")
			case "scenes":
				f.addLeftJoin("scenes_galleries", "scenes_join", "scenes_join.gallery_id = galleries.id")
				f.addWhere("scenes_join.gallery_id IS NULL")
			case "studio":
				f.addWhere("galleries.studio_id IS NULL")
			case "performers":
				qb.performersRepository().join(f, "performers_join", "galleries.id")
				f.addWhere("performers_join.gallery_id IS NULL")
			case "date":
				f.addWhere("galleries.date IS NULL OR galleries.date IS \"\"")
			case "tags":
				qb.tagsRepository().join(f, "tags_join", "galleries.id")
				f.addWhere("tags_join.gallery_id IS NULL")
			default:
				f.addWhere("(galleries." + *isMissing + " IS NULL OR TRIM(galleries." + *isMissing + ") = '')")
			}
		}
	}
}

func (qb *GalleryStore) tagsCriterionHandler(tags *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	h := joinedHierarchicalMultiCriterionHandlerBuilder{
		tx: qb.tx,

		primaryTable: galleryTable,
		foreignTable: tagTable,
		foreignFK:    "tag_id",

		relationsTable: "tags_relations",
		joinAs:         "image_tag",
		joinTable:      galleriesTagsTable,
		primaryFK:      galleryIDColumn,
	}

	return h.handler(tags)
}

func (qb *GalleryStore) tagCountCriterionHandler(tagCount *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: galleryTable,
		joinTable:    galleriesTagsTable,
		primaryFK:    galleryIDColumn,
	}

	return h.handler(tagCount)
}

func (qb *GalleryStore) scenesCriterionHandler(scenes *models.MultiCriterionInput) criterionHandlerFunc {
	addJoinsFunc := func(f *filterBuilder) {
		qb.scenesRepository().join(f, "", "galleries.id")
		f.addLeftJoin("scenes", "", "scenes_galleries.scene_id = scenes.id")
	}
	h := qb.getMultiCriterionHandlerBuilder(sceneTable, galleriesScenesTable, "scene_id", addJoinsFunc)
	return h.handler(scenes)
}

func (qb *GalleryStore) performersCriterionHandler(performers *models.MultiCriterionInput) criterionHandlerFunc {
	h := joinedMultiCriterionHandlerBuilder{
		primaryTable: galleryTable,
		joinTable:    performersGalleriesTable,
		joinAs:       "performers_join",
		primaryFK:    galleryIDColumn,
		foreignFK:    performerIDColumn,

		addJoinTable: func(f *filterBuilder) {
			qb.performersRepository().join(f, "performers_join", "galleries.id")
		},
	}

	return h.handler(performers)
}

func (qb *GalleryStore) performerCountCriterionHandler(performerCount *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: galleryTable,
		joinTable:    performersGalleriesTable,
		primaryFK:    galleryIDColumn,
	}

	return h.handler(performerCount)
}

func (qb *GalleryStore) imageCountCriterionHandler(imageCount *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: galleryTable,
		joinTable:    galleriesImagesTable,
		primaryFK:    galleryIDColumn,
	}

	return h.handler(imageCount)
}

func (qb *GalleryStore) hasChaptersCriterionHandler(hasChapters *string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if hasChapters != nil {
			f.addLeftJoin("galleries_chapters", "", "galleries_chapters.gallery_id = galleries.id")
			if *hasChapters == "true" {
				f.addHaving("count(galleries_chapters.gallery_id) > 0")
			} else {
				f.addWhere("galleries_chapters.id IS NULL")
			}
		}
	}
}

func (qb *GalleryStore) performerTagsCriterionHandler(tags *models.HierarchicalMultiCriterionInput) criterionHandler {
	return &joinedPerformerTagsHandler{
		criterion:      tags,
		primaryTable:   galleryTable,
		joinTable:      performersGalleriesTable,
		joinPrimaryKey: galleryIDColumn,
	}
}

func (qb *GalleryStore) performerFavoriteCriterionHandler(performerfavorite *bool) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if performerfavorite != nil {
			f.addLeftJoin("performers_galleries", "", "galleries.id = performers_galleries.gallery_id")

			if *performerfavorite {
				// contains at least one favorite
				f.addLeftJoin("performers", "", "performers.id = performers_galleries.performer_id")
				f.addWhere("performers.favorite = 1")
			} else {
				// contains zero favorites
				f.addLeftJoin(`(SELECT performers_galleries.gallery_id as id FROM performers_galleries 
JOIN performers ON performers.id = performers_galleries.performer_id
GROUP BY performers_galleries.gallery_id HAVING SUM(performers.favorite) = 0)`, "nofaves", "galleries.id = nofaves.id")
				f.addWhere("performers_galleries.gallery_id IS NULL OR nofaves.id IS NOT NULL")
			}
		}
	}
}

func (qb *GalleryStore) performerAgeCriterionHandler(performerAge *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if performerAge != nil {
			f.addInnerJoin("performers_galleries", "", "galleries.id = performers_galleries.gallery_id")
			f.addInnerJoin("performers", "", "performers_galleries.performer_id = performers.id")

			f.addWhere("galleries.date != '' AND performers.birthdate != ''")
			f.addWhere("galleries.date IS NOT NULL AND performers.birthdate IS NOT NULL")

			ageCalc := "cast(strftime('%Y.%m%d', galleries.date) - strftime('%Y.%m%d', performers.birthdate) as int)"
			whereClause, args := getIntWhereClause(ageCalc, performerAge.Modifier, performerAge.Value, performerAge.Value2)
			f.addWhere(whereClause, args...)
		}
	}
}

func (qb *GalleryStore) averageResolutionCriterionHandler(resolution *models.ResolutionCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if resolution != nil && resolution.Value.IsValid() {
			qb.imagesRepository().join(f, "images_join", "galleries.id")
			f.addLeftJoin("images", "", "images_join.image_id = images.id")
			f.addLeftJoin("images_files", "", "images.id = images_files.image_id")
			f.addLeftJoin("image_files", "", "images_files.file_id = image_files.file_id")

			min := resolution.Value.GetMinResolution()
			max := resolution.Value.GetMaxResolution()

			const widthHeight = "avg(MIN(image_files.width, image_files.height))"

			switch resolution.Modifier {
			case models.CriterionModifierEquals:
				f.addHaving(fmt.Sprintf("%s BETWEEN %d AND %d", widthHeight, min, max))
			case models.CriterionModifierNotEquals:
				f.addHaving(fmt.Sprintf("%s NOT BETWEEN %d AND %d", widthHeight, min, max))
			case models.CriterionModifierLessThan:
				f.addHaving(fmt.Sprintf("%s < %d", widthHeight, min))
			case models.CriterionModifierGreaterThan:
				f.addHaving(fmt.Sprintf("%s > %d", widthHeight, max))
			}
		}
	}
}
