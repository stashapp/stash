package sqlite

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

type mangaFilterHandler struct {
	mangaFilter *models.MangaFilterType
}

func (qb *mangaFilterHandler) validate() error {
	mangaFilter := qb.mangaFilter
	if mangaFilter == nil {
		return nil
	}

	if err := validateFilterCombination(mangaFilter.OperatorFilter); err != nil {
		return err
	}

	if subFilter := mangaFilter.SubFilter(); subFilter != nil {
		sqb := &mangaFilterHandler{mangaFilter: subFilter}
		if err := sqb.validate(); err != nil {
			return err
		}
	}

	return nil
}

func (qb *mangaFilterHandler) handle(ctx context.Context, f *filterBuilder) {
	mangaFilter := qb.mangaFilter
	if mangaFilter == nil {
		return
	}

	if err := qb.validate(); err != nil {
		f.setError(err)
		return
	}

	f.handleCriterion(ctx, qb.criterionHandler())

	sf := mangaFilter.SubFilter()
	if sf != nil {
		sub := &mangaFilterHandler{sf}
		handleSubFilter(ctx, sub, f, mangaFilter.OperatorFilter)
	}
}

func (qb *mangaFilterHandler) criterionHandler() criterionHandler {
	filter := qb.mangaFilter
	return compoundHandler{
		intCriterionHandler(filter.ID, "mangas.id", nil),
		stringCriterionHandler(filter.Title, "mangas.title"),
		stringCriterionHandler(filter.Details, "mangas.details"),
		stringCriterionHandler(filter.URL, "mangas.url"),

		qb.tagsCriterionHandler(filter.Tags),
		qb.tagCountCriterionHandler(filter.TagCount),
		qb.performersCriterionHandler(filter.Performers),
		qb.performerCountCriterionHandler(filter.PerformerCount),

		intCriterionHandler(filter.Rating100, "mangas.rating", nil),
		boolCriterionHandler(filter.Organized, "mangas.organized", nil),
		studioCriterionHandler(mangaTable, filter.Studios),

		&dateCriterionHandler{filter.Date, "mangas.date", nil},
		&timestampCriterionHandler{filter.CreatedAt, "mangas.created_at", nil},
		&timestampCriterionHandler{filter.UpdatedAt, "mangas.updated_at", nil},

		qb.performerTagsCriterionHandler(filter.PerformerTags),
		qb.performerFavoriteCriterionHandler(filter.PerformerFavorite),
		qb.performerAgeCriterionHandler(filter.PerformerAge),

		&relatedFilterHandler{
			relatedIDCol:   "performers_join.performer_id",
			relatedRepo:    performerRepository.repository,
			relatedHandler: &performerFilterHandler{filter.PerformersFilter},
			joinFn: func(f *filterBuilder) {
				mangaRepository.performers.innerJoin(f, "performers_join", "mangas.id")
			},
		},

		&relatedFilterHandler{
			relatedIDCol:   "manga_tag.tag_id",
			relatedRepo:    tagRepository.repository,
			relatedHandler: &tagFilterHandler{filter.TagsFilter},
			joinFn: func(f *filterBuilder) {
				mangaRepository.tags.innerJoin(f, "manga_tag", "mangas.id")
			},
		},

		&relatedFilterHandler{
			relatedIDCol:   "mangas.studio_id",
			relatedRepo:    studioRepository.repository,
			relatedHandler: &studioFilterHandler{filter.StudiosFilter},
		},
	}
}

func (qb *mangaFilterHandler) tagsCriterionHandler(tags *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	h := joinedHierarchicalMultiCriterionHandlerBuilder{
		primaryTable: mangaTable,
		foreignTable: tagTable,
		foreignFK:    "tag_id",

		relationsTable: "tags_relations",
		joinAs:         "manga_tag",
		joinTable:      mangasTagsTable,
		primaryFK:      mangaIDColumn,
	}

	return h.handler(tags)
}

func (qb *mangaFilterHandler) tagCountCriterionHandler(tagCount *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: mangaTable,
		joinTable:    mangasTagsTable,
		primaryFK:    mangaIDColumn,
	}

	return h.handler(tagCount)
}

func (qb *mangaFilterHandler) performersCriterionHandler(performers *models.MultiCriterionInput) criterionHandlerFunc {
	h := joinedMultiCriterionHandlerBuilder{
		primaryTable: mangaTable,
		joinTable:    performersMangasTable,
		joinAs:       "performers_join",
		primaryFK:    mangaIDColumn,
		foreignFK:    performerIDColumn,

		addJoinTable: func(f *filterBuilder, joinType joinType) {
			mangaRepository.performers.join(f, joinType, "performers_join", "mangas.id")
		},
	}

	return h.handler(performers)
}

func (qb *mangaFilterHandler) performerCountCriterionHandler(performerCount *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: mangaTable,
		joinTable:    performersMangasTable,
		primaryFK:    mangaIDColumn,
	}

	return h.handler(performerCount)
}

func (qb *mangaFilterHandler) performerTagsCriterionHandler(tags *models.HierarchicalMultiCriterionInput) criterionHandler {
	return &joinedPerformerTagsHandler{
		criterion:      tags,
		primaryTable:   mangaTable,
		joinTable:      performersMangasTable,
		joinPrimaryKey: mangaIDColumn,
	}
}

func (qb *mangaFilterHandler) performerFavoriteCriterionHandler(performerfavorite *bool) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if performerfavorite != nil {
			f.addLeftJoin("performers_mangas", "", "mangas.id = performers_mangas.manga_id")

			if *performerfavorite {
				// contains at least one favorite
				f.addLeftJoin("performers", "", "performers.id = performers_mangas.performer_id")
				f.addWhere("performers.favorite = 1")
			} else {
				// contains zero favorites
				f.addLeftJoin(`(SELECT performers_mangas.manga_id as id FROM performers_mangas 
JOIN performers ON performers.id = performers_mangas.performer_id
GROUP BY performers_mangas.manga_id HAVING SUM(performers.favorite) = 0)`, "nofaves", "mangas.id = nofaves.id")
				f.addWhere("performers_mangas.manga_id IS NULL OR nofaves.id IS NOT NULL")
			}
		}
	}
}

func (qb *mangaFilterHandler) performerAgeCriterionHandler(performerAge *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if performerAge != nil {
			f.addInnerJoin("performers_mangas", "", "mangas.id = performers_mangas.manga_id")
			f.addInnerJoin("performers", "", "performers_mangas.performer_id = performers.id")

			f.addWhere("mangas.date != '' AND performers.birthdate != ''")
			f.addWhere("mangas.date IS NOT NULL AND performers.birthdate IS NOT NULL")

			ageCalc := "cast(strftime('%Y.%m%d', mangas.date) - strftime('%Y.%m%d', performers.birthdate) as int)"
			whereClause, args := getIntWhereClause(ageCalc, performerAge.Modifier, performerAge.Value, performerAge.Value2)
			f.addWhere(whereClause, args...)
		}
	}
}
