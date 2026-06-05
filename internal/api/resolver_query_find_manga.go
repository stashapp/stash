package api

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) FindManga(ctx context.Context, id string) (ret *models.Manga, err error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Manga.Find(ctx, idInt)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) FindMangas(ctx context.Context, mangaFilter *models.MangaFilterType, filter *models.FindFilterType, ids []string) (ret *FindMangasResultType, err error) {
	idInts, err := handleIDList(ids, "ids")
	if err != nil {
		return nil, err
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var mangas []*models.Manga
		var err error
		var total int

		if len(idInts) > 0 {
			mangas, err = r.repository.Manga.FindMany(ctx, idInts)
			total = len(mangas)
		} else {
			mangas, total, err = r.repository.Manga.Query(ctx, mangaFilter, filter)
		}

		if err != nil {
			return err
		}

		ret = &FindMangasResultType{
			Count:  total,
			Mangas: mangas,
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) AllMangas(ctx context.Context) (ret []*models.Manga, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Manga.All(ctx)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
