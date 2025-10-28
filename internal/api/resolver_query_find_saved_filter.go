package api

import (
	"context"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

func (r *queryResolver) FindSavedFilter(ctx context.Context, id string) (ret *models.SavedFilter, err error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.SavedFilter.Find(ctx, idInt)
		return err
	}); err != nil {
		return nil, err
	}
	return ret, err
}

func (r *queryResolver) FindSavedFilters(ctx context.Context, mode *models.FilterMode) (ret []*models.SavedFilter, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		if mode != nil {
			ret, err = r.repository.SavedFilter.FindByMode(ctx, *mode)
		} else {
			ret, err = r.repository.SavedFilter.All(ctx)
		}
		return err
	}); err != nil {
		return nil, err
	}
	return ret, err
}

func (r *queryResolver) FindDefaultFilter(ctx context.Context, mode models.FilterMode) (ret *models.SavedFilter, err error) {
	// deprecated - read from the config in the meantime
	config := config.GetInstance()

	uiConfig := config.GetUIConfiguration()
	if uiConfig == nil {
		return nil, nil
	}

	m := utils.NestedMap(uiConfig)
	filterRaw, _ := m.Get("defaultFilters." + strings.ToLower(mode.String()))

	if filterRaw == nil {
		return nil, nil
	}

	ret = &models.SavedFilter{}
	d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName:          "json",
		WeaklyTypedInput: true,
		Result:           ret,
	})

	if err != nil {
		return nil, err
	}

	if err := d.Decode(filterRaw); err != nil {
		return nil, err
	}

	return ret, nil
}
