package api

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin/hook"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
)

// used to refetch manga after hooks run
func (r *mutationResolver) getManga(ctx context.Context, id int) (ret *models.Manga, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Manga.Find(ctx, id)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *mutationResolver) MangaCreate(ctx context.Context, input MangaCreateInput) (*models.Manga, error) {
	// name must be provided
	if input.Title == "" {
		return nil, errors.New("title must not be empty")
	}

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	// Populate a new manga from the input
	newManga := models.CreateMangaInput{
		Manga: &models.Manga{},
	}
	*newManga.Manga = models.NewManga()

	newManga.Title = strings.TrimSpace(input.Title)
	newManga.URL = translator.string(input.URL)
	newManga.Details = translator.string(input.Details)
	newManga.Rating = input.Rating100
	newManga.Organized = translator.bool(input.Organized)

	var err error

	newManga.Date, err = translator.datePtr(input.Date)
	if err != nil {
		return nil, fmt.Errorf("converting date: %w", err)
	}
	newManga.StudioID, err = translator.intPtrFromString(input.StudioID)
	if err != nil {
		return nil, fmt.Errorf("converting studio id: %w", err)
	}

	newManga.PerformerIDs, err = translator.relatedIds(input.PerformerIds)
	if err != nil {
		return nil, fmt.Errorf("converting performer ids: %w", err)
	}
	newManga.TagIDs, err = translator.relatedIds(input.TagIds)
	if err != nil {
		return nil, fmt.Errorf("converting tag ids: %w", err)
	}

	// Start the transaction and save the manga
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Manga
		if err := qb.Create(ctx, &newManga); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	r.hookExecutor.ExecutePostHooks(ctx, newManga.ID, hook.MangaCreatePost, input, nil)
	return r.getManga(ctx, newManga.ID)
}

func (r *mutationResolver) MangaUpdate(ctx context.Context, input models.MangaUpdateInput) (ret *models.Manga, err error) {
	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	// Start the transaction and save the manga
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.mangaUpdate(ctx, input, translator)
		return err
	}); err != nil {
		return nil, err
	}

	// execute post hooks outside txn
	r.hookExecutor.ExecutePostHooks(ctx, ret.ID, hook.MangaUpdatePost, input, translator.getFields())
	return r.getManga(ctx, ret.ID)
}

func (r *mutationResolver) MangasUpdate(ctx context.Context, input []*models.MangaUpdateInput) (ret []*models.Manga, err error) {
	inputMaps := getUpdateInputMaps(ctx)

	// Start the transaction and save the mangas
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		for i, manga := range input {
			translator := changesetTranslator{
				inputMap: inputMaps[i],
			}

			thisManga, err := r.mangaUpdate(ctx, *manga, translator)
			if err != nil {
				return err
			}

			ret = append(ret, thisManga)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// execute post hooks outside txn
	var newRet []*models.Manga
	for i, manga := range ret {
		translator := changesetTranslator{
			inputMap: inputMaps[i],
		}

		r.hookExecutor.ExecutePostHooks(ctx, manga.ID, hook.MangaUpdatePost, input, translator.getFields())

		manga, err = r.getManga(ctx, manga.ID)
		if err != nil {
			return nil, err
		}

		newRet = append(newRet, manga)
	}

	return newRet, nil
}

func (r *mutationResolver) mangaUpdate(ctx context.Context, input models.MangaUpdateInput, translator changesetTranslator) (*models.Manga, error) {
	mangaID, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, fmt.Errorf("converting id: %w", err)
	}

	qb := r.repository.Manga

	originalManga, err := qb.Find(ctx, mangaID)
	if err != nil {
		return nil, err
	}

	if originalManga == nil {
		return nil, fmt.Errorf("manga with id %d not found", mangaID)
	}

	// Populate manga from the input
	updatedManga := models.NewMangaPartial()

	if input.Title != nil {
		updatedManga.Title = models.NewOptionalString(*input.Title)
	}

	updatedManga.URL = translator.optionalString(input.URL, "url")
	updatedManga.Details = translator.optionalString(input.Details, "details")
	updatedManga.Rating = translator.optionalInt(input.Rating100, "rating100")
	updatedManga.Organized = translator.optionalBool(input.Organized, "organized")

	updatedManga.Date, err = translator.optionalDate(input.Date, "date")
	if err != nil {
		return nil, fmt.Errorf("converting date: %w", err)
	}
	updatedManga.StudioID, err = translator.optionalIntFromString(input.StudioID, "studio_id")
	if err != nil {
		return nil, fmt.Errorf("converting studio id: %w", err)
	}

	updatedManga.PerformerIDs, err = translator.updateIds(input.PerformerIds, "performer_ids")
	if err != nil {
		return nil, fmt.Errorf("converting performer ids: %w", err)
	}
	updatedManga.TagIDs, err = translator.updateIds(input.TagIds, "tag_ids")
	if err != nil {
		return nil, fmt.Errorf("converting tag ids: %w", err)
	}

	manga, err := qb.UpdatePartial(ctx, mangaID, updatedManga)
	if err != nil {
		return nil, err
	}

	return manga, nil
}

func (r *mutationResolver) BulkMangaUpdate(ctx context.Context, input BulkMangaUpdateInput) ([]*models.Manga, error) {
	mangaIDs, err := stringslice.StringSliceToIntSlice(input.Ids)
	if err != nil {
		return nil, fmt.Errorf("converting ids: %w", err)
	}

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	// Populate manga from the input
	updatedManga := models.NewMangaPartial()

	updatedManga.URL = translator.optionalString(input.URL, "url")
	updatedManga.Details = translator.optionalString(input.Details, "details")
	updatedManga.Rating = translator.optionalInt(input.Rating100, "rating100")
	updatedManga.Organized = translator.optionalBool(input.Organized, "organized")

	updatedManga.Date, err = translator.optionalDate(input.Date, "date")
	if err != nil {
		return nil, fmt.Errorf("converting date: %w", err)
	}
	updatedManga.StudioID, err = translator.optionalIntFromString(input.StudioID, "studio_id")
	if err != nil {
		return nil, fmt.Errorf("converting studio id: %w", err)
	}

	updatedManga.PerformerIDs, err = translator.updateIdsBulk(input.PerformerIds, "performer_ids")
	if err != nil {
		return nil, fmt.Errorf("converting performer ids: %w", err)
	}
	updatedManga.TagIDs, err = translator.updateIdsBulk(input.TagIds, "tag_ids")
	if err != nil {
		return nil, fmt.Errorf("converting tag ids: %w", err)
	}

	ret := []*models.Manga{}

	// Start the transaction and save the mangas
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Manga

		for _, mangaID := range mangaIDs {
			manga, err := qb.UpdatePartial(ctx, mangaID, updatedManga)
			if err != nil {
				return err
			}

			ret = append(ret, manga)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// execute post hooks outside of txn
	var newRet []*models.Manga
	for _, manga := range ret {
		r.hookExecutor.ExecutePostHooks(ctx, manga.ID, hook.MangaUpdatePost, input, translator.getFields())

		manga, err := r.getManga(ctx, manga.ID)
		if err != nil {
			return nil, err
		}

		newRet = append(newRet, manga)
	}

	return newRet, nil
}

func (r *mutationResolver) MangaDestroy(ctx context.Context, input MangaDestroyInput) (bool, error) {
	mangaIDs, err := stringslice.StringSliceToIntSlice(input.Ids)
	if err != nil {
		return false, fmt.Errorf("converting ids: %w", err)
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Manga

		for _, id := range mangaIDs {
			manga, err := qb.Find(ctx, id)
			if err != nil {
				return err
			}

			if manga == nil {
				return fmt.Errorf("manga with id %d not found", id)
			}

			if err := qb.Destroy(ctx, id); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return false, err
	}

	for _, mangaID := range mangaIDs {
		r.hookExecutor.ExecutePostHooks(ctx, mangaID, hook.MangaDestroyPost, input, nil)
	}

	return true, nil
}
