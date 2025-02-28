//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type updateImageFunc func(ctx context.Context, id int, image []byte) error
type getImageFunc func(ctx context.Context, id int) ([]byte, error)

func testUpdateImage(t *testing.T, ctx context.Context, id int, updateFn updateImageFunc, getFn getImageFunc) error {
	image := []byte("image")
	err := updateFn(ctx, id, image)
	if err != nil {
		return fmt.Errorf("Error updating performer image: %s", err.Error())
	}

	// ensure image set
	storedImage, err := getFn(ctx, id)
	if err != nil {
		return fmt.Errorf("Error getting image: %s", err.Error())
	}
	assert.Equal(t, storedImage, image)

	// set nil image
	err = updateFn(ctx, id, nil)
	if err != nil {
		return fmt.Errorf("error setting nil image: %w", err)
	}

	// ensure image null
	storedImage, err = getFn(ctx, id)
	if err != nil {
		return fmt.Errorf("Error getting image: %s", err.Error())
	}
	assert.Nil(t, storedImage)

	return nil
}
