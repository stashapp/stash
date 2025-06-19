package sqlite

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

type folderFilterHandler struct {
	folderFilter *models.FolderFilterType
}

func (qb *folderFilterHandler) validate() error {
	folderFilter := qb.folderFilter
	if folderFilter == nil {
		return nil
	}

	if err := validateFilterCombination(folderFilter.OperatorFilter); err != nil {
		return err
	}

	if subFilter := folderFilter.SubFilter(); subFilter != nil {
		sqb := &folderFilterHandler{folderFilter: subFilter}
		if err := sqb.validate(); err != nil {
			return err
		}
	}

	return nil
}

func (qb *folderFilterHandler) handle(ctx context.Context, f *filterBuilder) {
	folderFilter := qb.folderFilter
	if folderFilter == nil {
		return
	}

	if err := qb.validate(); err != nil {
		f.setError(err)
		return
	}

	sf := folderFilter.SubFilter()
	if sf != nil {
		sub := &folderFilterHandler{sf}
		handleSubFilter(ctx, sub, f, folderFilter.OperatorFilter)
	}

	f.handleCriterion(ctx, qb.criterionHandler())
}

func (qb *folderFilterHandler) criterionHandler() criterionHandler {
	folderFilter := qb.folderFilter
	return compoundHandler{
		stringCriterionHandler(folderFilter.Path, "folders.path"),
		&timestampCriterionHandler{folderFilter.ModTime, "folders.mod_time", nil},

		qb.parentFolderCriterionHandler(folderFilter.ParentFolder),

		qb.galleryCountCriterionHandler(folderFilter.GalleryCount),

		&timestampCriterionHandler{folderFilter.CreatedAt, "folders.created_at", nil},
		&timestampCriterionHandler{folderFilter.UpdatedAt, "folders.updated_at", nil},

		&relatedFilterHandler{
			relatedIDCol:   "galleries.id",
			relatedRepo:    galleryRepository.repository,
			relatedHandler: &galleryFilterHandler{folderFilter.GalleriesFilter},
			joinFn: func(f *filterBuilder) {
				folderRepository.galleries.innerJoin(f, "", "folders.id")
			},
		},
	}
}

func (qb *folderFilterHandler) parentFolderCriterionHandler(folder *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
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
			primaryTable: folderTable,
			foreignTable: folderTable,
			foreignFK:    "parent_folder_id",
			parentFK:     "parent_folder_id",
		}

		hh.handler(&folderCopy)(ctx, f)
	}
}

func (qb *folderFilterHandler) galleryCountCriterionHandler(galleryCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if galleryCount != nil {
			f.addLeftJoin("galleries", "", "galleries.folder_id = folders.id")
			clause, args := getIntCriterionWhereClause("count(distinct galleries.id)", *galleryCount)

			f.addHaving(clause, args...)
		}
	}
}
