package sqlite

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type imageQueryBuilder struct {
	imageFilter *models.ImageFilterType
}

func (qb *imageQueryBuilder) validateFilter(imageFilter *models.ImageFilterType) error {
	if imageFilter == nil {
		return nil
	}

	if err := validateFilterCombination(imageFilter.And, imageFilter.Or, imageFilter.Not); err != nil {
		return err
	}

	if subFilter := utils.FirstNotNil(imageFilter.And, imageFilter.Or, imageFilter.Not); subFilter != nil {
		if err := qb.validateFilter(subFilter); err != nil {
			return err
		}
	}

	// if err := galleryStore.validateFilter(imageFilter.GalleriesFilter); err != nil {
	// 	return err
	// }
	// if err := performerStore.validateFilter(imageFilter.PerformersFilter); err != nil {
	// 	return err
	// }
	// if err := studioStore.validateFilter(imageFilter.StudiosFilter); err != nil {
	// 	return err
	// }
	// if err := tagStore.validateFilter(imageFilter.TagsFilter); err != nil {
	// 	return err
	// }

	return nil
}

func (qb *imageQueryBuilder) makeFilter(ctx context.Context, imageFilter *models.ImageFilterType) *filterBuilder {
	if imageFilter == nil {
		return nil
	}

	query := &filterBuilder{}

	if imageFilter.And != nil {
		query.and(qb.makeFilter(ctx, imageFilter.And))
	}
	if imageFilter.Or != nil {
		query.or(qb.makeFilter(ctx, imageFilter.Or))
	}
	if imageFilter.Not != nil {
		query.not(qb.makeFilter(ctx, imageFilter.Not))
	}

	query.handleCriterion(ctx, qb.criterionHandler(imageFilter))

	return query
}

func (qb *imageQueryBuilder) criterionHandler(imageFilter *models.ImageFilterType) criterionHandler {
	return compoundHandler{
		intCriterionHandler(imageFilter.ID, "images.id", nil),
		criterionHandlerFunc(func(ctx context.Context, f *filterBuilder) {
			if imageFilter.Checksum != nil {
				imageRepository.addImagesFilesTable(f)
				f.addInnerJoin(fingerprintTable, "fingerprints_md5", "images_files.file_id = fingerprints_md5.file_id AND fingerprints_md5.type = 'md5'")
			}

			stringCriterionHandler(imageFilter.Checksum, "fingerprints_md5.fingerprint")(ctx, f)
		}),
		stringCriterionHandler(imageFilter.Title, "images.title"),
		stringCriterionHandler(imageFilter.Code, "images.code"),
		stringCriterionHandler(imageFilter.Details, "images.details"),
		stringCriterionHandler(imageFilter.Photographer, "images.photographer"),

		pathCriterionHandler(imageFilter.Path, "folders.path", "files.basename", imageRepository.addFoldersTable),
		qb.fileCountCriterionHandler(imageFilter.FileCount),
		intCriterionHandler(imageFilter.Rating100, "images.rating", nil),
		intCriterionHandler(imageFilter.OCounter, "images.o_counter", nil),
		boolCriterionHandler(imageFilter.Organized, "images.organized", nil),
		dateCriterionHandler(imageFilter.Date, "images.date"),
		qb.urlsCriterionHandler(imageFilter.URL),

		resolutionCriterionHandler(imageFilter.Resolution, "image_files.height", "image_files.width", imageRepository.addImageFilesTable),
		orientationCriterionHandler(imageFilter.Orientation, "image_files.height", "image_files.width", imageRepository.addImageFilesTable),
		qb.missingCriterionHandler(imageFilter.IsMissing),

		qb.tagsCriterionHandler(imageFilter.Tags),
		qb.tagCountCriterionHandler(imageFilter.TagCount),
		qb.galleriesCriterionHandler(imageFilter.Galleries),
		qb.performersCriterionHandler(imageFilter.Performers),
		qb.performerCountCriterionHandler(imageFilter.PerformerCount),
		studioCriterionHandler(imageTable, imageFilter.Studios),
		qb.performerTagsCriterionHandler(imageFilter.PerformerTags),
		qb.performerFavoriteCriterionHandler(imageFilter.PerformerFavorite),
		qb.performerAgeCriterionHandler(imageFilter.PerformerAge),
		timestampCriterionHandler(imageFilter.CreatedAt, "images.created_at"),
		timestampCriterionHandler(imageFilter.UpdatedAt, "images.updated_at"),

		// &relatedFilterHandler{
		// 	relatedIDCol: "galleries_images.gallery_id",
		// 	relatedStore: &galleryRepository,
		// 	makeFilterFn: func(ctx context.Context) *filterBuilder {
		// 		return galleryStore.makeFilter(ctx, imageFilter.GalleriesFilter)
		// 	},
		// 	joinFn: func(f *filterBuilder) {
		// 		f.addInnerJoin(galleriesImagesTable, "", "galleries_images.image_id = images.id")
		// 	},
		// },

		// &relatedFilterHandler{
		// 	relatedIDCol: "performers_join.performer_id",
		// 	relatedStore: &performerRepository,
		// 	makeFilterFn: func(ctx context.Context) *filterBuilder {
		// 		return performerStore.makeFilter(ctx, imageFilter.PerformersFilter)
		// 	},
		// 	joinFn: func(f *filterBuilder) {
		// 		qb.performersRepository().join(f, "performers_join", "images.id")
		// 	},
		// },

		// &relatedFilterHandler{
		// 	relatedIDCol: "images.studio_id",
		// 	relatedStore: studioStore,
		// 	makeFilterFn: func(ctx context.Context) *filterBuilder {
		// 		return studioStore.makeFilter(ctx, imageFilter.StudiosFilter)
		// 	},
		// },

		// &relatedFilterHandler{
		// 	relatedIDCol: "image_tag.tag_id",
		// 	relatedStore: tagStore,
		// 	makeFilterFn: func(ctx context.Context) *filterBuilder {
		// 		return tagStore.makeFilter(ctx, imageFilter.TagsFilter)
		// 	},
		// 	joinFn: func(f *filterBuilder) {
		// 		f.addInnerJoin(imagesTagsTable, "image_tag", "image_tag.image_id = images.id")
		// 	},
		// },
	}
}

func (qb *imageQueryBuilder) fileCountCriterionHandler(fileCount *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: imageTable,
		joinTable:    imagesFilesTable,
		primaryFK:    imageIDColumn,
	}

	return h.handler(fileCount)
}

func (qb *imageQueryBuilder) missingCriterionHandler(isMissing *string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if isMissing != nil && *isMissing != "" {
			switch *isMissing {
			case "studio":
				f.addWhere("images.studio_id IS NULL")
			case "performers":
				imageRepository.performers.join(f, "performers_join", "images.id")
				f.addWhere("performers_join.image_id IS NULL")
			case "galleries":
				imageRepository.galleries.join(f, "galleries_join", "images.id")
				f.addWhere("galleries_join.image_id IS NULL")
			case "tags":
				imageRepository.tags.join(f, "tags_join", "images.id")
				f.addWhere("tags_join.image_id IS NULL")
			default:
				f.addWhere("(images." + *isMissing + " IS NULL OR TRIM(images." + *isMissing + ") = '')")
			}
		}
	}
}

func (qb *imageQueryBuilder) urlsCriterionHandler(url *models.StringCriterionInput) criterionHandlerFunc {
	h := stringListCriterionHandlerBuilder{
		joinTable:    imagesURLsTable,
		stringColumn: imageURLColumn,
		addJoinTable: func(f *filterBuilder) {
			imagesURLsTableMgr.join(f, "", "images.id")
		},
	}

	return h.handler(url)
}

func (qb *imageQueryBuilder) getMultiCriterionHandlerBuilder(foreignTable, joinTable, foreignFK string, addJoinsFunc func(f *filterBuilder)) multiCriterionHandlerBuilder {
	return multiCriterionHandlerBuilder{
		primaryTable: imageTable,
		foreignTable: foreignTable,
		joinTable:    joinTable,
		primaryFK:    imageIDColumn,
		foreignFK:    foreignFK,
		addJoinsFunc: addJoinsFunc,
	}
}

func (qb *imageQueryBuilder) tagsCriterionHandler(tags *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	h := joinedHierarchicalMultiCriterionHandlerBuilder{
		primaryTable: imageTable,
		foreignTable: tagTable,
		foreignFK:    "tag_id",

		relationsTable: "tags_relations",
		joinAs:         "image_tag",
		joinTable:      imagesTagsTable,
		primaryFK:      imageIDColumn,
	}

	return h.handler(tags)
}

func (qb *imageQueryBuilder) tagCountCriterionHandler(tagCount *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: imageTable,
		joinTable:    imagesTagsTable,
		primaryFK:    imageIDColumn,
	}

	return h.handler(tagCount)
}

func (qb *imageQueryBuilder) galleriesCriterionHandler(galleries *models.MultiCriterionInput) criterionHandlerFunc {
	addJoinsFunc := func(f *filterBuilder) {
		if galleries.Modifier == models.CriterionModifierIncludes || galleries.Modifier == models.CriterionModifierIncludesAll {
			f.addInnerJoin(galleriesImagesTable, "", "galleries_images.image_id = images.id")
			f.addInnerJoin(galleryTable, "", "galleries_images.gallery_id = galleries.id")
		}
	}
	h := qb.getMultiCriterionHandlerBuilder(galleryTable, galleriesImagesTable, galleryIDColumn, addJoinsFunc)

	return h.handler(galleries)
}

func (qb *imageQueryBuilder) performersCriterionHandler(performers *models.MultiCriterionInput) criterionHandlerFunc {
	h := joinedMultiCriterionHandlerBuilder{
		primaryTable: imageTable,
		joinTable:    performersImagesTable,
		joinAs:       "performers_join",
		primaryFK:    imageIDColumn,
		foreignFK:    performerIDColumn,

		addJoinTable: func(f *filterBuilder) {
			imageRepository.performers.join(f, "performers_join", "images.id")
		},
	}

	return h.handler(performers)
}

func (qb *imageQueryBuilder) performerCountCriterionHandler(performerCount *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: imageTable,
		joinTable:    performersImagesTable,
		primaryFK:    imageIDColumn,
	}

	return h.handler(performerCount)
}

func (qb *imageQueryBuilder) performerFavoriteCriterionHandler(performerfavorite *bool) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if performerfavorite != nil {
			f.addLeftJoin("performers_images", "", "images.id = performers_images.image_id")

			if *performerfavorite {
				// contains at least one favorite
				f.addLeftJoin("performers", "", "performers.id = performers_images.performer_id")
				f.addWhere("performers.favorite = 1")
			} else {
				// contains zero favorites
				f.addLeftJoin(`(SELECT performers_images.image_id as id FROM performers_images
JOIN performers ON performers.id = performers_images.performer_id
GROUP BY performers_images.image_id HAVING SUM(performers.favorite) = 0)`, "nofaves", "images.id = nofaves.id")
				f.addWhere("performers_images.image_id IS NULL OR nofaves.id IS NOT NULL")
			}
		}
	}
}

func (qb *imageQueryBuilder) performerAgeCriterionHandler(performerAge *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if performerAge != nil {
			f.addInnerJoin("performers_images", "", "images.id = performers_images.image_id")
			f.addInnerJoin("performers", "", "performers_images.performer_id = performers.id")

			f.addWhere("images.date != '' AND performers.birthdate != ''")
			f.addWhere("images.date IS NOT NULL AND performers.birthdate IS NOT NULL")

			ageCalc := "cast(strftime('%Y.%m%d', images.date) - strftime('%Y.%m%d', performers.birthdate) as int)"
			whereClause, args := getIntWhereClause(ageCalc, performerAge.Modifier, performerAge.Value, performerAge.Value2)
			f.addWhere(whereClause, args...)
		}
	}
}

func (qb *imageQueryBuilder) performerTagsCriterionHandler(tags *models.HierarchicalMultiCriterionInput) criterionHandler {
	return &joinedPerformerTagsHandler{
		criterion:      tags,
		primaryTable:   imageTable,
		joinTable:      performersImagesTable,
		joinPrimaryKey: imageIDColumn,
	}
}
