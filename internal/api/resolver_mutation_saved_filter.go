package api

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

func (r *mutationResolver) SaveFilter(ctx context.Context, input SaveFilterInput) (ret *models.SavedFilter, err error) {
	if strings.TrimSpace(input.Name) == "" {
		return nil, errors.New("name must be non-empty")
	}

	var id *int
	if input.ID != nil {
		idv, err := strconv.Atoi(*input.ID)
		if err != nil {
			return nil, fmt.Errorf("converting id: %w", err)
		}
		id = &idv
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.SavedFilter

		f := models.SavedFilter{
			Mode:         input.Mode,
			Name:         input.Name,
			FindFilter:   input.FindFilter,
			ObjectFilter: input.ObjectFilter,
			UIOptions:    input.UIOptions,
		}

		if id == nil {
			err = qb.Create(ctx, &f)
			ret = &f
		} else {
			f.ID = *id
			err = qb.Update(ctx, &f)
			ret = &f
		}

		return err
	}); err != nil {
		return nil, err
	}
	return ret, err
}

func (r *mutationResolver) DestroySavedFilter(ctx context.Context, input DestroyFilterInput) (bool, error) {
	id, err := strconv.Atoi(input.ID)
	if err != nil {
		return false, fmt.Errorf("converting id: %w", err)
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		return r.repository.SavedFilter.Destroy(ctx, id)
	}); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) SetDefaultFilter(ctx context.Context, input SetDefaultFilterInput) (bool, error) {
	// deprecated - write to the config in the meantime
	config := config.GetInstance()

	uiConfig := config.GetUIConfiguration()
	if uiConfig == nil {
		uiConfig = make(map[string]interface{})
	}

	m := utils.NestedMap(uiConfig)

	if input.FindFilter == nil && input.ObjectFilter == nil && input.UIOptions == nil {
		// clearing
		m.Delete("defaultFilters." + strings.ToLower(input.Mode.String()))
		config.SetUIConfiguration(m)

		if err := config.Write(); err != nil {
			return false, err
		}

		return true, nil
	}

	subMap := make(map[string]interface{})
	d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName:          "json",
		WeaklyTypedInput: true,
		Result:           &subMap,
	})

	if err != nil {
		return false, err
	}

	if err := d.Decode(input); err != nil {
		return false, err
	}

	m.Set("defaultFilters."+strings.ToLower(input.Mode.String()), subMap)

	config.SetUIConfiguration(m)

	if err := config.Write(); err != nil {
		return false, err
	}

	return true, nil
}
