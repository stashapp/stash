package sqlite

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type sceneFilterHandler struct {
	sceneFilter *models.SceneFilterType
}

func (qb *sceneFilterHandler) validate() error {
	sceneFilter := qb.sceneFilter
	if sceneFilter == nil {
		return nil
	}

	if err := validateFilterCombination(sceneFilter.OperatorFilter); err != nil {
		return err
	}

	if subFilter := sceneFilter.SubFilter(); subFilter != nil {
		sqb := &sceneFilterHandler{sceneFilter: subFilter}
		if err := sqb.validate(); err != nil {
			return err
		}
	}

	return nil
}

func (qb *sceneFilterHandler) handle(ctx context.Context, f *filterBuilder) {
	sceneFilter := qb.sceneFilter
	if sceneFilter == nil {
		return
	}

	if err := qb.validate(); err != nil {
		f.setError(err)
		return
	}

	sf := sceneFilter.SubFilter()
	if sf != nil {
		sub := &sceneFilterHandler{sf}
		handleSubFilter(ctx, sub, f, sceneFilter.OperatorFilter)
	}

	f.handleCriterion(ctx, qb.criterionHandler())
}

func (qb *sceneFilterHandler) criterionHandler() criterionHandler {
	sceneFilter := qb.sceneFilter
	return compoundHandler{
		intCriterionHandler(sceneFilter.ID, "scenes.id", nil),
		pathCriterionHandler(sceneFilter.Path, "folders.path", "files.basename", qb.addFoldersTable),
		qb.fileCountCriterionHandler(sceneFilter.FileCount),
		stringCriterionHandler(sceneFilter.Title, "scenes.title"),
		stringCriterionHandler(sceneFilter.Code, "scenes.code"),
		stringCriterionHandler(sceneFilter.Details, "scenes.details"),
		stringCriterionHandler(sceneFilter.Director, "scenes.director"),
		criterionHandlerFunc(func(ctx context.Context, f *filterBuilder) {
			if sceneFilter.Oshash != nil {
				qb.addSceneFilesTable(f)
				f.addLeftJoin(fingerprintTable, "fingerprints_oshash", "scenes_files.file_id = fingerprints_oshash.file_id AND fingerprints_oshash.type = 'oshash'")
			}

			stringCriterionHandler(sceneFilter.Oshash, "fingerprints_oshash.fingerprint")(ctx, f)
		}),

		criterionHandlerFunc(func(ctx context.Context, f *filterBuilder) {
			if sceneFilter.Checksum != nil {
				qb.addSceneFilesTable(f)
				f.addLeftJoin(fingerprintTable, "fingerprints_md5", "scenes_files.file_id = fingerprints_md5.file_id AND fingerprints_md5.type = 'md5'")
			}

			stringCriterionHandler(sceneFilter.Checksum, "fingerprints_md5.fingerprint")(ctx, f)
		}),

		criterionHandlerFunc(func(ctx context.Context, f *filterBuilder) {
			if sceneFilter.Phash != nil {
				// backwards compatibility
				qb.phashDistanceCriterionHandler(&models.PhashDistanceCriterionInput{
					Value:    sceneFilter.Phash.Value,
					Modifier: sceneFilter.Phash.Modifier,
				})(ctx, f)
			}
		}),

		qb.phashDistanceCriterionHandler(sceneFilter.PhashDistance),

		intCriterionHandler(sceneFilter.Rating100, "scenes.rating", nil),
		qb.oCountCriterionHandler(sceneFilter.OCounter),
		boolCriterionHandler(sceneFilter.Organized, "scenes.organized", nil),

		floatIntCriterionHandler(sceneFilter.Duration, "video_files.duration", qb.addVideoFilesTable),
		resolutionCriterionHandler(sceneFilter.Resolution, "video_files.height", "video_files.width", qb.addVideoFilesTable),
		orientationCriterionHandler(sceneFilter.Orientation, "video_files.height", "video_files.width", qb.addVideoFilesTable),
		floatIntCriterionHandler(sceneFilter.Framerate, "ROUND(video_files.frame_rate)", qb.addVideoFilesTable),
		intCriterionHandler(sceneFilter.Bitrate, "video_files.bit_rate", qb.addVideoFilesTable),
		qb.codecCriterionHandler(sceneFilter.VideoCodec, "video_files.video_codec", qb.addVideoFilesTable),
		qb.codecCriterionHandler(sceneFilter.AudioCodec, "video_files.audio_codec", qb.addVideoFilesTable),

		qb.hasMarkersCriterionHandler(sceneFilter.HasMarkers),
		qb.isMissingCriterionHandler(sceneFilter.IsMissing),
		qb.urlsCriterionHandler(sceneFilter.URL),

		criterionHandlerFunc(func(ctx context.Context, f *filterBuilder) {
			if sceneFilter.StashID != nil {
				sceneRepository.stashIDs.join(f, "scene_stash_ids", "scenes.id")
				stringCriterionHandler(sceneFilter.StashID, "scene_stash_ids.stash_id")(ctx, f)
			}
		}),

		&stashIDCriterionHandler{
			c:                 sceneFilter.StashIDEndpoint,
			stashIDRepository: &sceneRepository.stashIDs,
			stashIDTableAs:    "scene_stash_ids",
			parentIDCol:       "scenes.id",
		},

		boolCriterionHandler(sceneFilter.Interactive, "video_files.interactive", qb.addVideoFilesTable),
		intCriterionHandler(sceneFilter.InteractiveSpeed, "video_files.interactive_speed", qb.addVideoFilesTable),

		qb.captionCriterionHandler(sceneFilter.Captions),

		floatIntCriterionHandler(sceneFilter.ResumeTime, "scenes.resume_time", nil),
		floatIntCriterionHandler(sceneFilter.PlayDuration, "scenes.play_duration", nil),
		qb.playCountCriterionHandler(sceneFilter.PlayCount),
		criterionHandlerFunc(func(ctx context.Context, f *filterBuilder) {
			if sceneFilter.LastPlayedAt != nil {
				f.addLeftJoin(
					fmt.Sprintf("(SELECT %s, MAX(%s) as last_played_at FROM %s GROUP BY %s)", sceneIDColumn, sceneViewDateColumn, scenesViewDatesTable, sceneIDColumn),
					"scene_last_view",
					fmt.Sprintf("scene_last_view.%s = scenes.id", sceneIDColumn),
				)
				h := timestampCriterionHandler{sceneFilter.LastPlayedAt, "IFNULL(last_played_at, datetime(0))", nil}
				h.handle(ctx, f)
			}
		}),

		qb.tagsCriterionHandler(sceneFilter.Tags),
		qb.tagCountCriterionHandler(sceneFilter.TagCount),
		qb.performersCriterionHandler(sceneFilter.Performers),
		qb.performerCountCriterionHandler(sceneFilter.PerformerCount),
		studioCriterionHandler(sceneTable, sceneFilter.Studios),

		qb.groupsCriterionHandler(sceneFilter.Groups),
		qb.groupsCriterionHandler(sceneFilter.Movies),

		qb.galleriesCriterionHandler(sceneFilter.Galleries),
		qb.performerTagsCriterionHandler(sceneFilter.PerformerTags),
		qb.performerFavoriteCriterionHandler(sceneFilter.PerformerFavorite),
		qb.performerAgeCriterionHandler(sceneFilter.PerformerAge),
		qb.phashDuplicatedCriterionHandler(sceneFilter.Duplicated, qb.addSceneFilesTable),
		&dateCriterionHandler{sceneFilter.Date, "scenes.date", nil},
		&timestampCriterionHandler{sceneFilter.CreatedAt, "scenes.created_at", nil},
		&timestampCriterionHandler{sceneFilter.UpdatedAt, "scenes.updated_at", nil},

		&relatedFilterHandler{
			relatedIDCol:   "scenes_galleries.gallery_id",
			relatedRepo:    galleryRepository.repository,
			relatedHandler: &galleryFilterHandler{sceneFilter.GalleriesFilter},
			joinFn: func(f *filterBuilder) {
				sceneRepository.galleries.innerJoin(f, "", "scenes.id")
			},
		},

		&relatedFilterHandler{
			relatedIDCol:   "performers_join.performer_id",
			relatedRepo:    performerRepository.repository,
			relatedHandler: &performerFilterHandler{sceneFilter.PerformersFilter},
			joinFn: func(f *filterBuilder) {
				sceneRepository.performers.innerJoin(f, "performers_join", "scenes.id")
			},
		},

		&relatedFilterHandler{
			relatedIDCol:   "scenes.studio_id",
			relatedRepo:    studioRepository.repository,
			relatedHandler: &studioFilterHandler{sceneFilter.StudiosFilter},
		},

		&relatedFilterHandler{
			relatedIDCol:   "scene_tag.tag_id",
			relatedRepo:    tagRepository.repository,
			relatedHandler: &tagFilterHandler{sceneFilter.TagsFilter},
			joinFn: func(f *filterBuilder) {
				sceneRepository.tags.innerJoin(f, "scene_tag", "scenes.id")
			},
		},

		&relatedFilterHandler{
			relatedIDCol:   "movies_scenes.movie_id",
			relatedRepo:    groupRepository.repository,
			relatedHandler: &groupFilterHandler{sceneFilter.MoviesFilter},
			joinFn: func(f *filterBuilder) {
				sceneRepository.groups.innerJoin(f, "", "scenes.id")
			},
		},

		&relatedFilterHandler{
			relatedIDCol:   "scene_markers.id",
			relatedRepo:    sceneMarkerRepository.repository,
			relatedHandler: &sceneMarkerFilterHandler{sceneFilter.MarkersFilter},
			joinFn: func(f *filterBuilder) {
				f.addInnerJoin("scene_markers", "", "scenes.id")
			},
		},
	}
}

func (qb *sceneFilterHandler) addSceneFilesTable(f *filterBuilder) {
	f.addLeftJoin(scenesFilesTable, "", "scenes_files.scene_id = scenes.id")
}

func (qb *sceneFilterHandler) addFilesTable(f *filterBuilder) {
	qb.addSceneFilesTable(f)
	f.addLeftJoin(fileTable, "", "scenes_files.file_id = files.id")
}

func (qb *sceneFilterHandler) addFoldersTable(f *filterBuilder) {
	qb.addFilesTable(f)
	f.addLeftJoin(folderTable, "", "files.parent_folder_id = folders.id")
}

func (qb *sceneFilterHandler) addVideoFilesTable(f *filterBuilder) {
	qb.addSceneFilesTable(f)
	f.addLeftJoin(videoFileTable, "", "video_files.file_id = scenes_files.file_id")
}

func (qb *sceneFilterHandler) playCountCriterionHandler(count *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: sceneTable,
		joinTable:    scenesViewDatesTable,
		primaryFK:    sceneIDColumn,
	}

	return h.handler(count)
}

func (qb *sceneFilterHandler) oCountCriterionHandler(count *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: sceneTable,
		joinTable:    scenesODatesTable,
		primaryFK:    sceneIDColumn,
	}

	return h.handler(count)
}

func (qb *sceneFilterHandler) fileCountCriterionHandler(fileCount *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: sceneTable,
		joinTable:    scenesFilesTable,
		primaryFK:    sceneIDColumn,
	}

	return h.handler(fileCount)
}

func (qb *sceneFilterHandler) phashDuplicatedCriterionHandler(duplicatedFilter *models.PHashDuplicationCriterionInput, addJoinFn func(f *filterBuilder)) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		// TODO: Wishlist item: Implement Distance matching
		if duplicatedFilter != nil {
			if addJoinFn != nil {
				addJoinFn(f)
			}

			var v string
			if *duplicatedFilter.Duplicated {
				v = ">"
			} else {
				v = "="
			}

			f.addInnerJoin("(SELECT file_id FROM files_fingerprints INNER JOIN (SELECT fingerprint FROM files_fingerprints WHERE type = 'phash' GROUP BY fingerprint HAVING COUNT (fingerprint) "+v+" 1) dupes on files_fingerprints.fingerprint = dupes.fingerprint)", "scph", "scenes_files.file_id = scph.file_id")
		}
	}
}

func (qb *sceneFilterHandler) codecCriterionHandler(codec *models.StringCriterionInput, codecColumn string, addJoinFn func(f *filterBuilder)) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if codec != nil {
			if addJoinFn != nil {
				addJoinFn(f)
			}

			stringCriterionHandler(codec, codecColumn)(ctx, f)
		}
	}
}

func (qb *sceneFilterHandler) hasMarkersCriterionHandler(hasMarkers *string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if hasMarkers != nil {
			f.addLeftJoin("scene_markers", "", "scene_markers.scene_id = scenes.id")
			if *hasMarkers == "true" {
				f.addHaving("count(scene_markers.scene_id) > 0")
			} else {
				f.addWhere("scene_markers.id IS NULL")
			}
		}
	}
}

func (qb *sceneFilterHandler) isMissingCriterionHandler(isMissing *string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if isMissing != nil && *isMissing != "" {
			switch *isMissing {
			case "url":
				scenesURLsTableMgr.join(f, "", "scenes.id")
				f.addWhere("scene_urls.url IS NULL")
			case "galleries":
				sceneRepository.galleries.join(f, "galleries_join", "scenes.id")
				f.addWhere("galleries_join.scene_id IS NULL")
			case "studio":
				f.addWhere("scenes.studio_id IS NULL")
			case "movie":
				sceneRepository.groups.join(f, "movies_join", "scenes.id")
				f.addWhere("movies_join.scene_id IS NULL")
			case "performers":
				sceneRepository.performers.join(f, "performers_join", "scenes.id")
				f.addWhere("performers_join.scene_id IS NULL")
			case "date":
				f.addWhere(`scenes.date IS NULL OR scenes.date IS ""`)
			case "tags":
				sceneRepository.tags.join(f, "tags_join", "scenes.id")
				f.addWhere("tags_join.scene_id IS NULL")
			case "stash_id":
				sceneRepository.stashIDs.join(f, "scene_stash_ids", "scenes.id")
				f.addWhere("scene_stash_ids.scene_id IS NULL")
			case "phash":
				qb.addSceneFilesTable(f)
				f.addLeftJoin(fingerprintTable, "fingerprints_phash", "scenes_files.file_id = fingerprints_phash.file_id AND fingerprints_phash.type = 'phash'")
				f.addWhere("fingerprints_phash.fingerprint IS NULL")
			case "cover":
				f.addWhere("scenes.cover_blob IS NULL")
			default:
				f.addWhere("(scenes." + *isMissing + " IS NULL OR TRIM(scenes." + *isMissing + ") = '')")
			}
		}
	}
}

func (qb *sceneFilterHandler) urlsCriterionHandler(url *models.StringCriterionInput) criterionHandlerFunc {
	h := stringListCriterionHandlerBuilder{
		primaryTable: sceneTable,
		primaryFK:    sceneIDColumn,
		joinTable:    scenesURLsTable,
		stringColumn: sceneURLColumn,
		addJoinTable: func(f *filterBuilder) {
			scenesURLsTableMgr.join(f, "", "scenes.id")
		},
	}

	return h.handler(url)
}

func (qb *sceneFilterHandler) getMultiCriterionHandlerBuilder(foreignTable, joinTable, foreignFK string, addJoinsFunc func(f *filterBuilder)) multiCriterionHandlerBuilder {
	return multiCriterionHandlerBuilder{
		primaryTable: sceneTable,
		foreignTable: foreignTable,
		joinTable:    joinTable,
		primaryFK:    sceneIDColumn,
		foreignFK:    foreignFK,
		addJoinsFunc: addJoinsFunc,
	}
}

func (qb *sceneFilterHandler) captionCriterionHandler(captions *models.StringCriterionInput) criterionHandlerFunc {
	h := stringListCriterionHandlerBuilder{
		primaryTable: sceneTable,
		primaryFK:    sceneIDColumn,
		joinTable:    videoCaptionsTable,
		stringColumn: captionCodeColumn,
		addJoinTable: func(f *filterBuilder) {
			qb.addSceneFilesTable(f)
			f.addLeftJoin(videoCaptionsTable, "", "video_captions.file_id = scenes_files.file_id")
		},
		excludeHandler: func(f *filterBuilder, criterion *models.StringCriterionInput) {
			excludeClause := `scenes.id NOT IN (
				SELECT scenes_files.scene_id from scenes_files 
				INNER JOIN video_captions on video_captions.file_id = scenes_files.file_id 
				WHERE video_captions.language_code LIKE ?
			)`
			f.addWhere(excludeClause, criterion.Value)

			// TODO - should we also exclude null values?
		},
	}

	return h.handler(captions)
}

func (qb *sceneFilterHandler) tagsCriterionHandler(tags *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	h := joinedHierarchicalMultiCriterionHandlerBuilder{
		primaryTable: sceneTable,
		foreignTable: tagTable,
		foreignFK:    "tag_id",

		relationsTable: "tags_relations",
		joinAs:         "scene_tag",
		joinTable:      scenesTagsTable,
		primaryFK:      sceneIDColumn,
	}

	return h.handler(tags)
}

func (qb *sceneFilterHandler) tagCountCriterionHandler(tagCount *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: sceneTable,
		joinTable:    scenesTagsTable,
		primaryFK:    sceneIDColumn,
	}

	return h.handler(tagCount)
}

func (qb *sceneFilterHandler) performersCriterionHandler(performers *models.MultiCriterionInput) criterionHandlerFunc {
	h := joinedMultiCriterionHandlerBuilder{
		primaryTable: sceneTable,
		joinTable:    performersScenesTable,
		joinAs:       "performers_join",
		primaryFK:    sceneIDColumn,
		foreignFK:    performerIDColumn,

		addJoinTable: func(f *filterBuilder) {
			sceneRepository.performers.join(f, "performers_join", "scenes.id")
		},
	}

	return h.handler(performers)
}

func (qb *sceneFilterHandler) performerCountCriterionHandler(performerCount *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: sceneTable,
		joinTable:    performersScenesTable,
		primaryFK:    sceneIDColumn,
	}

	return h.handler(performerCount)
}

func (qb *sceneFilterHandler) performerFavoriteCriterionHandler(performerfavorite *bool) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if performerfavorite != nil {
			f.addLeftJoin("performers_scenes", "", "scenes.id = performers_scenes.scene_id")

			if *performerfavorite {
				// contains at least one favorite
				f.addLeftJoin("performers", "", "performers.id = performers_scenes.performer_id")
				f.addWhere("performers.favorite = 1")
			} else {
				// contains zero favorites
				f.addLeftJoin(`(SELECT performers_scenes.scene_id as id FROM performers_scenes
JOIN performers ON performers.id = performers_scenes.performer_id
GROUP BY performers_scenes.scene_id HAVING SUM(performers.favorite) = 0)`, "nofaves", "scenes.id = nofaves.id")
				f.addWhere("performers_scenes.scene_id IS NULL OR nofaves.id IS NOT NULL")
			}
		}
	}
}

func (qb *sceneFilterHandler) performerAgeCriterionHandler(performerAge *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if performerAge != nil {
			f.addInnerJoin("performers_scenes", "", "scenes.id = performers_scenes.scene_id")
			f.addInnerJoin("performers", "", "performers_scenes.performer_id = performers.id")

			f.addWhere("scenes.date != '' AND performers.birthdate != ''")
			f.addWhere("scenes.date IS NOT NULL AND performers.birthdate IS NOT NULL")

			ageCalc := "cast(strftime('%Y.%m%d', scenes.date) - strftime('%Y.%m%d', performers.birthdate) as int)"
			whereClause, args := getIntWhereClause(ageCalc, performerAge.Modifier, performerAge.Value, performerAge.Value2)
			f.addWhere(whereClause, args...)
		}
	}
}

func (qb *sceneFilterHandler) groupsCriterionHandler(movies *models.MultiCriterionInput) criterionHandlerFunc {
	addJoinsFunc := func(f *filterBuilder) {
		sceneRepository.groups.join(f, "", "scenes.id")
		f.addLeftJoin("movies", "", "movies_scenes.movie_id = movies.id")
	}
	h := qb.getMultiCriterionHandlerBuilder(groupTable, groupsScenesTable, "movie_id", addJoinsFunc)
	return h.handler(movies)
}

func (qb *sceneFilterHandler) galleriesCriterionHandler(galleries *models.MultiCriterionInput) criterionHandlerFunc {
	addJoinsFunc := func(f *filterBuilder) {
		sceneRepository.galleries.join(f, "", "scenes.id")
		f.addLeftJoin("galleries", "", "scenes_galleries.gallery_id = galleries.id")
	}
	h := qb.getMultiCriterionHandlerBuilder(galleryTable, scenesGalleriesTable, "gallery_id", addJoinsFunc)
	return h.handler(galleries)
}

func (qb *sceneFilterHandler) performerTagsCriterionHandler(tags *models.HierarchicalMultiCriterionInput) criterionHandler {
	return &joinedPerformerTagsHandler{
		criterion:      tags,
		primaryTable:   sceneTable,
		joinTable:      performersScenesTable,
		joinPrimaryKey: sceneIDColumn,
	}
}

func (qb *sceneFilterHandler) phashDistanceCriterionHandler(phashDistance *models.PhashDistanceCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if phashDistance != nil {
			qb.addSceneFilesTable(f)
			f.addLeftJoin(fingerprintTable, "fingerprints_phash", "scenes_files.file_id = fingerprints_phash.file_id AND fingerprints_phash.type = 'phash'")

			value, _ := utils.StringToPhash(phashDistance.Value)
			distance := 0
			if phashDistance.Distance != nil {
				distance = *phashDistance.Distance
			}

			if distance == 0 {
				// use the default handler
				intCriterionHandler(&models.IntCriterionInput{
					Value:    int(value),
					Modifier: phashDistance.Modifier,
				}, "fingerprints_phash.fingerprint", nil)(ctx, f)
			}

			switch {
			case phashDistance.Modifier == models.CriterionModifierEquals && distance > 0:
				// needed to avoid a type mismatch
				f.addWhere("typeof(fingerprints_phash.fingerprint) = 'integer'")
				f.addWhere("phash_distance(fingerprints_phash.fingerprint, ?) < ?", value, distance)
			case phashDistance.Modifier == models.CriterionModifierNotEquals && distance > 0:
				// needed to avoid a type mismatch
				f.addWhere("typeof(fingerprints_phash.fingerprint) = 'integer'")
				f.addWhere("phash_distance(fingerprints_phash.fingerprint, ?) > ?", value, distance)
			default:
				intCriterionHandler(&models.IntCriterionInput{
					Value:    int(value),
					Modifier: phashDistance.Modifier,
				}, "fingerprints_phash.fingerprint", nil)(ctx, f)
			}
		}
	}
}
