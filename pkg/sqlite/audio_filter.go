package sqlite

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

type audioFilterHandler struct {
	audioFilter *models.AudioFilterType
}

func (qb *audioFilterHandler) validate() error {
	audioFilter := qb.audioFilter
	if audioFilter == nil {
		return nil
	}

	if err := validateFilterCombination(audioFilter.OperatorFilter); err != nil {
		return err
	}

	if subFilter := audioFilter.SubFilter(); subFilter != nil {
		sqb := &audioFilterHandler{audioFilter: subFilter}
		if err := sqb.validate(); err != nil {
			return err
		}
	}

	return nil
}

func (qb *audioFilterHandler) handle(ctx context.Context, f *filterBuilder) {
	audioFilter := qb.audioFilter
	if audioFilter == nil {
		return
	}

	if err := qb.validate(); err != nil {
		f.setError(err)
		return
	}

	f.handleCriterion(ctx, qb.criterionHandler())

	sf := audioFilter.SubFilter()
	if sf != nil {
		sub := &audioFilterHandler{sf}
		handleSubFilter(ctx, sub, f, audioFilter.OperatorFilter)
	}
}

func (qb *audioFilterHandler) criterionHandler() criterionHandler {
	audioFilter := qb.audioFilter
	return compoundHandler{
		intCriterionHandler(audioFilter.ID, "audios.id", nil),
		pathCriterionHandler(audioFilter.Path, "folders.path", "files.basename", qb.addFoldersTable),
		qb.fileCountCriterionHandler(audioFilter.FileCount),
		stringCriterionHandler(audioFilter.Title, "audios.title"),
		stringCriterionHandler(audioFilter.Code, "audios.code"),
		stringCriterionHandler(audioFilter.Details, "audios.details"),

		criterionHandlerFunc(func(ctx context.Context, f *filterBuilder) {
			if audioFilter.Checksum != nil {
				joinType := joinTypeInner
				if audioFilter.Checksum.Modifier == models.CriterionModifierIsNull {
					joinType = joinTypeLeft
				}
				qb.addAudioFilesTable(f, joinType)
				f.addJoin(joinType, fingerprintTable, "fingerprints_md5", "audios_files.file_id = fingerprints_md5.file_id AND fingerprints_md5.type = 'md5'")
			}

			stringCriterionHandler(audioFilter.Checksum, "fingerprints_md5.fingerprint")(ctx, f)
		}),

		intCriterionHandler(audioFilter.Rating100, "audios.rating", nil),
		qb.oCountCriterionHandler(audioFilter.OCounter),
		boolCriterionHandler(audioFilter.Organized, "audios.organized", nil),

		floatIntCriterionHandler(audioFilter.Duration, "audio_files.duration", qb.addAudioFilesTable),
		intCriterionHandler(audioFilter.SampleRate, "audio_files.sample_rate", qb.addAudioFilesTable),
		intCriterionHandler(audioFilter.Bitrate, "audio_files.bit_rate", qb.addAudioFilesTable),
		qb.codecCriterionHandler(audioFilter.AudioCodec, "audio_files.audio_codec", qb.addAudioFilesTable),

		qb.isMissingCriterionHandler(audioFilter.IsMissing),
		qb.urlsCriterionHandler(audioFilter.URL),

		floatIntCriterionHandler(audioFilter.ResumeTime, "audios.resume_time", nil),
		floatIntCriterionHandler(audioFilter.PlayDuration, "audios.play_duration", nil),
		qb.playCountCriterionHandler(audioFilter.PlayCount),
		criterionHandlerFunc(func(ctx context.Context, f *filterBuilder) {
			if audioFilter.LastPlayedAt != nil {
				f.addLeftJoin(
					fmt.Sprintf("(SELECT %s, MAX(%s) as last_played_at FROM %s GROUP BY %s)", audioIDColumn, audioPlayDateColumn, audiosPlayDatesTable, audioIDColumn),
					"audio_last_view",
					fmt.Sprintf("audio_last_view.%s = audios.id", audioIDColumn),
				)
				h := timestampCriterionHandler{audioFilter.LastPlayedAt, "IFNULL(last_played_at, datetime(0))", nil}
				h.handle(ctx, f)
			}
		}),

		qb.tagsCriterionHandler(audioFilter.Tags),
		qb.tagCountCriterionHandler(audioFilter.TagCount),
		qb.performersCriterionHandler(audioFilter.Performers),
		qb.performerCountCriterionHandler(audioFilter.PerformerCount),
		studioCriterionHandler(audioTable, audioFilter.Studios),

		qb.groupsCriterionHandler(audioFilter.Groups),

		qb.performerTagsCriterionHandler(audioFilter.PerformerTags),
		qb.performerFavoriteCriterionHandler(audioFilter.PerformerFavorite),
		qb.performerAgeCriterionHandler(audioFilter.PerformerAge),
		&dateCriterionHandler{audioFilter.Date, "audios.date", nil},
		&timestampCriterionHandler{audioFilter.CreatedAt, "audios.created_at", nil},
		&timestampCriterionHandler{audioFilter.UpdatedAt, "audios.updated_at", nil},

		&customFieldsFilterHandler{
			table: audiosCustomFieldsTable.GetTable(),
			fkCol: audioIDColumn,
			c:     audioFilter.CustomFields,
			idCol: "audios.id",
		},

		&relatedFilterHandler{
			relatedIDCol:   "performers_join.performer_id",
			relatedRepo:    performerRepository.repository,
			relatedHandler: &performerFilterHandler{audioFilter.PerformersFilter},
			joinFn: func(f *filterBuilder) {
				audioRepository.performers.innerJoin(f, "performers_join", "audios.id")
			},
		},

		&relatedFilterHandler{
			relatedIDCol:   "audios.studio_id",
			relatedRepo:    studioRepository.repository,
			relatedHandler: &studioFilterHandler{audioFilter.StudiosFilter},
		},

		&relatedFilterHandler{
			relatedIDCol:   "audio_tag.tag_id",
			relatedRepo:    tagRepository.repository,
			relatedHandler: &tagFilterHandler{audioFilter.TagsFilter},
			joinFn: func(f *filterBuilder) {
				audioRepository.tags.innerJoin(f, "audio_tag", "audios.id")
			},
		},

		&relatedFilterHandler{
			relatedIDCol:   "groups_audios.group_id",
			relatedRepo:    groupRepository.repository,
			relatedHandler: &groupFilterHandler{audioFilter.GroupsFilter},
			joinFn: func(f *filterBuilder) {
				audioRepository.groups.innerJoin(f, "", "audios.id")
			},
		},

		&relatedFilterHandler{
			relatedIDCol: "files.id",
			relatedRepo:  fileRepository.repository,
			relatedHandler: &fileFilterHandler{
				fileFilter: audioFilter.FilesFilter,
				isRelated:  true,
			},
			joinFn: func(f *filterBuilder) {
				qb.addFilesTable(f, joinTypeInner)
				qb.addFoldersTable(f, joinTypeInner)
			},
			// don't use a subquery; join directly
			directJoin: true,
		},
	}
}

func (qb *audioFilterHandler) addAudioFilesTable(f *filterBuilder, joinType joinType) {
	f.addJoin(joinType, audiosFilesTable, "", "audios_files.audio_id = audios.id")
}

func (qb *audioFilterHandler) addFilesTable(f *filterBuilder, joinType joinType) {
	qb.addAudioFilesTable(f, joinType)
	f.addJoin(joinType, fileTable, "", "audios_files.file_id = files.id")
}

func (qb *audioFilterHandler) addFoldersTable(f *filterBuilder, joinType joinType) {
	qb.addFilesTable(f, joinType)
	f.addJoin(joinType, folderTable, "", "files.parent_folder_id = folders.id")
}

func (qb *audioFilterHandler) playCountCriterionHandler(count *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: audioTable,
		joinTable:    audiosPlayDatesTable,
		primaryFK:    audioIDColumn,
	}

	return h.handler(count)
}

func (qb *audioFilterHandler) oCountCriterionHandler(count *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: audioTable,
		joinTable:    audiosODatesTable,
		primaryFK:    audioIDColumn,
	}

	return h.handler(count)
}

func (qb *audioFilterHandler) fileCountCriterionHandler(fileCount *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: audioTable,
		joinTable:    audiosFilesTable,
		primaryFK:    audioIDColumn,
	}

	return h.handler(fileCount)
}

func (qb *audioFilterHandler) codecCriterionHandler(codec *models.StringCriterionInput, codecColumn string, addJoinFn func(f *filterBuilder, joinType joinType)) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if codec != nil {
			if addJoinFn != nil {
				joinType := joinTypeInner
				if codec.Modifier == models.CriterionModifierIsNull {
					joinType = joinTypeLeft
				}
				addJoinFn(f, joinType)
			}

			stringCriterionHandler(codec, codecColumn)(ctx, f)
		}
	}
}

func (qb *audioFilterHandler) isMissingCriterionHandler(isMissing *string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if isMissing != nil && *isMissing != "" {
			switch *isMissing {
			case "url":
				audiosURLsTableMgr.leftJoin(f, "", "audios.id")
				f.addWhere("audio_urls.url IS NULL")
			case "studio":
				f.addWhere("audios.studio_id IS NULL")
			case "movie", "group":
				audioRepository.groups.leftJoin(f, "groups_join", "audios.id")
				f.addWhere("groups_join.audio_id IS NULL")
			case "performers":
				audioRepository.performers.leftJoin(f, "performers_join", "audios.id")
				f.addWhere("performers_join.audio_id IS NULL")
			case "date":
				f.addWhere(`audios.date IS NULL OR audios.date IS ""`)
			case "tags":
				audioRepository.tags.leftJoin(f, "tags_join", "audios.id")
				f.addWhere("tags_join.audio_id IS NULL")
			default:
				if err := validateIsMissing(*isMissing, []string{
					"title", "code", "details", "rating",
				}); err != nil {
					f.setError(err)
					return
				}
				f.addWhere("(audios." + *isMissing + " IS NULL OR TRIM(audios." + *isMissing + ") = '')")
			}
		}
	}
}

func (qb *audioFilterHandler) urlsCriterionHandler(url *models.StringCriterionInput) criterionHandlerFunc {
	h := stringListCriterionHandlerBuilder{
		primaryTable: audioTable,
		primaryFK:    audioIDColumn,
		joinTable:    audiosURLsTable,
		stringColumn: audioURLColumn,
		addJoinTable: func(f *filterBuilder, joinType joinType) {
			audiosURLsTableMgr.join(f, joinType, "", "audios.id")
		},
	}

	return h.handler(url)
}

func (qb *audioFilterHandler) tagsCriterionHandler(tags *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	h := joinedHierarchicalMultiCriterionHandlerBuilder{
		primaryTable: audioTable,
		foreignTable: tagTable,
		foreignFK:    "tag_id",

		relationsTable: "tags_relations",
		joinAs:         "audio_tag",
		joinTable:      audiosTagsTable,
		primaryFK:      audioIDColumn,
	}

	return h.handler(tags)
}

func (qb *audioFilterHandler) tagCountCriterionHandler(tagCount *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: audioTable,
		joinTable:    audiosTagsTable,
		primaryFK:    audioIDColumn,
	}

	return h.handler(tagCount)
}

func (qb *audioFilterHandler) performersCriterionHandler(performers *models.MultiCriterionInput) criterionHandlerFunc {
	h := joinedMultiCriterionHandlerBuilder{
		primaryTable: audioTable,
		joinTable:    performersAudiosTable,
		joinAs:       "performers_join",
		primaryFK:    audioIDColumn,
		foreignFK:    performerIDColumn,

		addJoinTable: func(f *filterBuilder, joinType joinType) {
			audioRepository.performers.join(f, joinType, "performers_join", "audios.id")
		},
	}

	return h.handler(performers)
}

func (qb *audioFilterHandler) performerCountCriterionHandler(performerCount *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: audioTable,
		joinTable:    performersAudiosTable,
		primaryFK:    audioIDColumn,
	}

	return h.handler(performerCount)
}

func (qb *audioFilterHandler) performerFavoriteCriterionHandler(performerfavorite *bool) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if performerfavorite != nil {
			f.addLeftJoin("performers_audios", "", "audios.id = performers_audios.audio_id")

			if *performerfavorite {
				// contains at least one favorite
				f.addLeftJoin("performers", "", "performers.id = performers_audios.performer_id")
				f.addWhere("performers.favorite = 1")
			} else {
				// contains zero favorites
				f.addLeftJoin(`(SELECT performers_audios.audio_id as id FROM performers_audios
JOIN performers ON performers.id = performers_audios.performer_id
GROUP BY performers_audios.audio_id HAVING SUM(performers.favorite) = 0)`, "nofaves", "audios.id = nofaves.id")
				f.addWhere("performers_audios.audio_id IS NULL OR nofaves.id IS NOT NULL")
			}
		}
	}
}

func (qb *audioFilterHandler) performerAgeCriterionHandler(performerAge *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if performerAge != nil {
			f.addInnerJoin("performers_audios", "", "audios.id = performers_audios.audio_id")
			f.addInnerJoin("performers", "", "performers_audios.performer_id = performers.id")

			f.addWhere("audios.date != '' AND performers.birthdate != ''")
			f.addWhere("audios.date IS NOT NULL AND performers.birthdate IS NOT NULL")

			ageCalc := "cast(strftime('%Y.%m%d', audios.date) - strftime('%Y.%m%d', performers.birthdate) as int)"
			whereClause, args := getIntWhereClause(ageCalc, performerAge.Modifier, performerAge.Value, performerAge.Value2)
			f.addWhere(whereClause, args...)
		}
	}
}

func (qb *audioFilterHandler) groupsCriterionHandler(groups *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	h := joinedHierarchicalMultiCriterionHandlerBuilder{
		primaryTable: audioTable,
		foreignTable: groupTable,
		foreignFK:    "group_id",

		relationsTable: groupRelationsTable,
		parentFK:       "containing_id",
		childFK:        "sub_id",
		joinAs:         "audio_group",
		joinTable:      groupsAudiosTable,
		primaryFK:      audioIDColumn,
	}

	return h.handler(groups)
}

func (qb *audioFilterHandler) performerTagsCriterionHandler(tags *models.HierarchicalMultiCriterionInput) criterionHandler {
	return &joinedPerformerTagsHandler{
		criterion:      tags,
		primaryTable:   audioTable,
		joinTable:      performersAudiosTable,
		joinPrimaryKey: audioIDColumn,
	}
}
