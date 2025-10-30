package sqlite

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

type folderFilterHandler struct {
	folderFilter *models.FolderFilterType
	table        sqlTable
	isRelated    bool
}

func (qb *folderFilterHandler) validate() error {
	folderFilter := qb.folderFilter
	if folderFilter == nil {
		return nil
	}

	if err := validateFilterCombination(folderFilter.OperatorFilter); err != nil {
		return err
	}

	if qb.isRelated && (folderFilter.GalleriesFilter != nil) {
		return fmt.Errorf("cannot use related filters inside a related filter")
	}

	if subFilter := folderFilter.SubFilter(); subFilter != nil {
		sqb := &folderFilterHandler{folderFilter: subFilter, isRelated: qb.isRelated}
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
		sub := &folderFilterHandler{folderFilter: sf, table: qb.table}
		handleSubFilter(ctx, sub, f, folderFilter.OperatorFilter)
	}

	f.handleCriterion(ctx, qb.criterionHandler())
}

func (qb *folderFilterHandler) criterionHandler() criterionHandler {
	if qb.table == "" {
		qb.table = folderTable
	}

	folderFilter := qb.folderFilter
	return compoundHandler{
		stringCriterionHandler(folderFilter.Path, qb.table.Col("path")),
		&timestampCriterionHandler{folderFilter.ModTime, qb.table.Col("mod_time"), nil},

		qb.parentFolderCriterionHandler(folderFilter.ParentFolder),
		qb.zipFileCriterionHandler(folderFilter.ZipFile),

		qb.galleryCountCriterionHandler(folderFilter.GalleryCount),

		&timestampCriterionHandler{folderFilter.CreatedAt, qb.table.Col("created_at"), nil},
		&timestampCriterionHandler{folderFilter.UpdatedAt, qb.table.Col("updated_at"), nil},

		&relatedFilterHandler{
			relatedIDCol:   qb.table.Col("id"),
			relatedRepo:    galleryRepository.repository,
			relatedHandler: &galleryFilterHandler{folderFilter.GalleriesFilter},
			joinFn: func(f *filterBuilder) {
				folderRepository.galleries.innerJoin(f, "", qb.table.Col("id"))
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

				f.addWhere(fmt.Sprintf("%s.zip_file_id IS %s NULL", qb.table.Name(), notClause))
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
				whereClause = fmt.Sprintf("%s.zip_file_id IN %s", qb.table.Name(), getInBinding(len(criterion.Value)))
			case models.CriterionModifierExcludes:
				whereClause = fmt.Sprintf("%s.zip_file_id NOT IN %s", qb.table.Name(), getInBinding(len(criterion.Value)))
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
			primaryTable: qb.table.Name(),
			foreignTable: qb.table.Name(),
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
