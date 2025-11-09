package sqlite

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type fileFilterHandler struct {
	fileFilter *models.FileFilterType
	// if true, don't allow use of related filters
	isRelated bool
}

func (qb *fileFilterHandler) validate() error {
	fileFilter := qb.fileFilter
	if fileFilter == nil {
		return nil
	}

	if err := validateFilterCombination(fileFilter.OperatorFilter); err != nil {
		return err
	}

	if qb.isRelated && (fileFilter.ScenesFilter != nil || fileFilter.ImagesFilter != nil || fileFilter.GalleriesFilter != nil) {
		return fmt.Errorf("cannot use related filters inside a related filter")
	}

	if subFilter := fileFilter.SubFilter(); subFilter != nil {
		sqb := &fileFilterHandler{fileFilter: subFilter, isRelated: qb.isRelated}
		if err := sqb.validate(); err != nil {
			return err
		}
	}

	return nil
}

func (qb *fileFilterHandler) handle(ctx context.Context, f *filterBuilder) {
	fileFilter := qb.fileFilter
	if fileFilter == nil {
		return
	}

	if err := qb.validate(); err != nil {
		f.setError(err)
		return
	}

	sf := fileFilter.SubFilter()
	if sf != nil {
		sub := &fileFilterHandler{sf, qb.isRelated}
		handleSubFilter(ctx, sub, f, fileFilter.OperatorFilter)
	}

	f.handleCriterion(ctx, qb.criterionHandler())
}

func (qb *fileFilterHandler) criterionHandler() criterionHandler {
	fileFilter := qb.fileFilter
	return compoundHandler{
		&videoFileFilterHandler{
			filter: fileFilter.VideoFileFilter,
		},
		&imageFileFilterHandler{
			filter: fileFilter.ImageFileFilter,
		},

		pathCriterionHandler(fileFilter.Path, "folders.path", "files.basename", nil),
		stringCriterionHandler(fileFilter.Basename, "files.basename"),
		stringCriterionHandler(fileFilter.Dir, "folders.path"),
		&timestampCriterionHandler{fileFilter.ModTime, "files.mod_time", nil},

		qb.parentFolderCriterionHandler(fileFilter.ParentFolder),
		qb.zipFileCriterionHandler(fileFilter.ZipFile),

		qb.sceneCountCriterionHandler(fileFilter.SceneCount),
		qb.imageCountCriterionHandler(fileFilter.ImageCount),
		qb.galleryCountCriterionHandler(fileFilter.GalleryCount),

		qb.hashesCriterionHandler(fileFilter.Hashes),

		qb.phashDuplicatedCriterionHandler(fileFilter.Duplicated),
		&timestampCriterionHandler{fileFilter.CreatedAt, "files.created_at", nil},
		&timestampCriterionHandler{fileFilter.UpdatedAt, "files.updated_at", nil},

		&relatedFilterHandler{
			relatedIDCol:   "scenes_files.scene_id",
			relatedRepo:    sceneRepository.repository,
			relatedHandler: &sceneFilterHandler{fileFilter.ScenesFilter},
			joinFn: func(f *filterBuilder) {
				fileRepository.scenes.innerJoin(f, "", "files.id")
			},
		},
		&relatedFilterHandler{
			relatedIDCol:   "images_files.image_id",
			relatedRepo:    imageRepository.repository,
			relatedHandler: &imageFilterHandler{fileFilter.ImagesFilter},
			joinFn: func(f *filterBuilder) {
				fileRepository.images.innerJoin(f, "", "files.id")
			},
		},
		&relatedFilterHandler{
			relatedIDCol:   "galleries_files.gallery_id",
			relatedRepo:    galleryRepository.repository,
			relatedHandler: &galleryFilterHandler{fileFilter.GalleriesFilter},
			joinFn: func(f *filterBuilder) {
				fileRepository.galleries.innerJoin(f, "", "files.id")
			},
		},
	}
}

func (qb *fileFilterHandler) zipFileCriterionHandler(criterion *models.MultiCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if criterion != nil {
			if criterion.Modifier == models.CriterionModifierIsNull || criterion.Modifier == models.CriterionModifierNotNull {
				var notClause string
				if criterion.Modifier == models.CriterionModifierNotNull {
					notClause = "NOT"
				}

				f.addWhere(fmt.Sprintf("files.zip_file_id IS %s NULL", notClause))
				return
			}

			if len(criterion.Value) == 0 {
				return
			}

			var args []interface{}
			for _, tagID := range criterion.Value {
				args = append(args, tagID)
			}

			whereClause := ""
			havingClause := ""
			switch criterion.Modifier {
			case models.CriterionModifierIncludes:
				whereClause = "files.zip_file_id IN " + getInBinding(len(criterion.Value))
			case models.CriterionModifierExcludes:
				whereClause = "files.zip_file_id NOT IN " + getInBinding(len(criterion.Value))
			}

			f.addWhere(whereClause, args...)
			f.addHaving(havingClause)
		}
	}
}

func (qb *fileFilterHandler) parentFolderCriterionHandler(folder *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if folder == nil {
			return
		}

		folderCopy := *folder
		switch folderCopy.Modifier {
		case models.CriterionModifierEquals:
			folderCopy.Modifier = models.CriterionModifierIncludesAll
		case models.CriterionModifierNotEquals:
			folderCopy.Modifier = models.CriterionModifierExcludes
		}

		hh := hierarchicalMultiCriterionHandlerBuilder{
			primaryTable: fileTable,
			foreignTable: folderTable,
			foreignFK:    "parent_folder_id",
			parentFK:     "parent_folder_id",
		}

		hh.handler(&folderCopy)(ctx, f)
	}
}

func (qb *fileFilterHandler) sceneCountCriterionHandler(c *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: fileTable,
		joinTable:    scenesFilesTable,
		primaryFK:    fileIDColumn,
	}

	return h.handler(c)
}

func (qb *fileFilterHandler) imageCountCriterionHandler(c *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: fileTable,
		joinTable:    imagesFilesTable,
		primaryFK:    fileIDColumn,
	}

	return h.handler(c)
}

func (qb *fileFilterHandler) galleryCountCriterionHandler(c *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: fileTable,
		joinTable:    galleriesFilesTable,
		primaryFK:    fileIDColumn,
	}

	return h.handler(c)
}

func (qb *fileFilterHandler) phashDuplicatedCriterionHandler(duplicatedFilter *models.PHashDuplicationCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		// TODO: Wishlist item: Implement Distance matching
		if duplicatedFilter != nil {
			var v string
			if *duplicatedFilter.Duplicated {
				v = ">"
			} else {
				v = "="
			}

			f.addInnerJoin("(SELECT file_id FROM files_fingerprints INNER JOIN (SELECT fingerprint FROM files_fingerprints WHERE type = 'phash' GROUP BY fingerprint HAVING COUNT (fingerprint) "+v+" 1) dupes on files_fingerprints.fingerprint = dupes.fingerprint)", "scph", "files.id = scph.file_id")
		}
	}
}

func (qb *fileFilterHandler) hashesCriterionHandler(hashes []*models.FingerprintFilterInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		// TODO - this won't work for AND/OR combinations
		for i, hash := range hashes {
			t := fmt.Sprintf("file_fingerprints_%d", i)
			f.addLeftJoin(fingerprintTable, t, fmt.Sprintf("files.id = %s.file_id AND %s.type = ?", t, t), hash.Type)

			value, _ := utils.StringToPhash(hash.Value)
			distance := 0
			if hash.Distance != nil {
				distance = *hash.Distance
			}

			if distance > 0 {
				// needed to avoid a type mismatch
				f.addWhere(fmt.Sprintf("typeof(%s.fingerprint) = 'integer'", t))
				f.addWhere(fmt.Sprintf("phash_distance(%s.fingerprint, ?) < ?", t), value, distance)
			} else {
				// use the default handler
				intCriterionHandler(&models.IntCriterionInput{
					Value:    int(value),
					Modifier: models.CriterionModifierEquals,
				}, t+".fingerprint", nil)(ctx, f)
			}
		}
	}
}

type videoFileFilterHandler struct {
	filter *models.VideoFileFilterInput
}

func (qb *videoFileFilterHandler) handle(ctx context.Context, f *filterBuilder) {
	videoFileFilter := qb.filter
	if videoFileFilter == nil {
		return
	}
	f.handleCriterion(ctx, qb.criterionHandler())
}

func (qb *videoFileFilterHandler) criterionHandler() criterionHandler {
	videoFileFilter := qb.filter
	return compoundHandler{
		joinedStringCriterionHandler(videoFileFilter.Format, "video_files.format", qb.addVideoFilesTable),
		floatIntCriterionHandler(videoFileFilter.Duration, "video_files.duration", qb.addVideoFilesTable),
		resolutionCriterionHandler(videoFileFilter.Resolution, "video_files.height", "video_files.width", qb.addVideoFilesTable),
		orientationCriterionHandler(videoFileFilter.Orientation, "video_files.height", "video_files.width", qb.addVideoFilesTable),
		floatIntCriterionHandler(videoFileFilter.Framerate, "ROUND(video_files.frame_rate)", qb.addVideoFilesTable),
		intCriterionHandler(videoFileFilter.Bitrate, "video_files.bit_rate", qb.addVideoFilesTable),
		qb.codecCriterionHandler(videoFileFilter.VideoCodec, "video_files.video_codec", qb.addVideoFilesTable),
		qb.codecCriterionHandler(videoFileFilter.AudioCodec, "video_files.audio_codec", qb.addVideoFilesTable),

		boolCriterionHandler(videoFileFilter.Interactive, "video_files.interactive", qb.addVideoFilesTable),
		intCriterionHandler(videoFileFilter.InteractiveSpeed, "video_files.interactive_speed", qb.addVideoFilesTable),

		qb.captionCriterionHandler(videoFileFilter.Captions),
	}
}

func (qb *videoFileFilterHandler) addVideoFilesTable(f *filterBuilder) {
	f.addLeftJoin(videoFileTable, "", "video_files.file_id = files.id")
}

func (qb *videoFileFilterHandler) codecCriterionHandler(codec *models.StringCriterionInput, codecColumn string, addJoinFn func(f *filterBuilder)) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if codec != nil {
			if addJoinFn != nil {
				addJoinFn(f)
			}

			stringCriterionHandler(codec, codecColumn)(ctx, f)
		}
	}
}

func (qb *videoFileFilterHandler) captionCriterionHandler(captions *models.StringCriterionInput) criterionHandlerFunc {
	h := stringListCriterionHandlerBuilder{
		primaryTable: sceneTable,
		primaryFK:    sceneIDColumn,
		joinTable:    videoCaptionsTable,
		stringColumn: captionCodeColumn,
		addJoinTable: func(f *filterBuilder) {
			f.addLeftJoin(videoCaptionsTable, "", "video_captions.file_id = files.id")
		},
		excludeHandler: func(f *filterBuilder, criterion *models.StringCriterionInput) {
			excludeClause := `files.id NOT IN (
				SELECT files.id from files 
				INNER JOIN video_captions on video_captions.file_id = files.id 
				WHERE video_captions.language_code LIKE ?
			)`
			f.addWhere(excludeClause, criterion.Value)

			// TODO - should we also exclude null values?
		},
	}

	return h.handler(captions)
}

type imageFileFilterHandler struct {
	filter *models.ImageFileFilterInput
}

func (qb *imageFileFilterHandler) handle(ctx context.Context, f *filterBuilder) {
	ff := qb.filter
	if ff == nil {
		return
	}
	f.handleCriterion(ctx, qb.criterionHandler())
}

func (qb *imageFileFilterHandler) criterionHandler() criterionHandler {
	ff := qb.filter
	return compoundHandler{
		joinedStringCriterionHandler(ff.Format, "image_files.format", qb.addImageFilesTable),
		resolutionCriterionHandler(ff.Resolution, "image_files.height", "image_files.width", qb.addImageFilesTable),
		orientationCriterionHandler(ff.Orientation, "image_files.height", "image_files.width", qb.addImageFilesTable),
	}
}

func (qb *imageFileFilterHandler) addImageFilesTable(f *filterBuilder) {
	f.addLeftJoin(imageFileTable, "", "image_files.file_id = files.id")
}
