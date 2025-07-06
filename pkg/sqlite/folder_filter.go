package sqlite

import (
	"context"
	"fmt"

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
		qb.zipFileCriterionHandler(folderFilter.ZipFile),

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

func (qb *folderFilterHandler) zipFileCriterionHandler(criterion *models.MultiCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if criterion != nil {
			if criterion.Modifier == models.CriterionModifierIsNull || criterion.Modifier == models.CriterionModifierNotNull {
				var notClause string
				if criterion.Modifier == models.CriterionModifierNotNull {
					notClause = "NOT"
				}

				f.addWhere(fmt.Sprintf("folders.zip_file_id IS %s NULL", notClause))
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
				whereClause = "folders.zip_file_id IN " + getInBinding(len(criterion.Value))
			case models.CriterionModifierExcludes:
				whereClause = "folders.zip_file_id NOT IN " + getInBinding(len(criterion.Value))
			}

			f.addWhere(whereClause, args...)
			f.addHaving(havingClause)
		}
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
