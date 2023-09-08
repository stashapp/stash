package mocks

import (
	"github.com/stretchr/testify/assert"

	"github.com/stashapp/stash/pkg/models"
)

// asserts that got == expected
// ignores expected.UpdatedAt, but ensures that got.UpdatedAt is set and not null
func AssertGalleryPartial(t assert.TestingT, got, expected models.GalleryPartial) bool {
	// updated at should be set and not null
	if !got.UpdatedAt.Set || got.UpdatedAt.Null {
		return false
	}
	// else ignore the exact value
	got.UpdatedAt = models.OptionalTime{}

	return assert.Equal(t, got, expected)
}

// asserts that got == expected
// ignores expected.UpdatedAt, but ensures that got.UpdatedAt is set and not null
func AssertImagePartial(t assert.TestingT, got, expected models.ImagePartial) bool {
	// updated at should be set and not null
	if !got.UpdatedAt.Set || got.UpdatedAt.Null {
		return false
	}
	// else ignore the exact value
	got.UpdatedAt = models.OptionalTime{}

	return assert.Equal(t, got, expected)
}

// asserts that got == expected
// ignores expected.UpdatedAt, but ensures that got.UpdatedAt is set and not null
func AssertPerformerPartial(t assert.TestingT, got, expected models.PerformerPartial) bool {
	// updated at should be set and not null
	if !got.UpdatedAt.Set || got.UpdatedAt.Null {
		return false
	}
	// else ignore the exact value
	got.UpdatedAt = models.OptionalTime{}

	return assert.Equal(t, got, expected)
}

// asserts that got == expected
// ignores expected.UpdatedAt, but ensures that got.UpdatedAt is set and not null
func AssertScenePartial(t assert.TestingT, got, expected models.ScenePartial) bool {
	// updated at should be set and not null
	if !got.UpdatedAt.Set || got.UpdatedAt.Null {
		return false
	}
	// else ignore the exact value
	got.UpdatedAt = models.OptionalTime{}

	return assert.Equal(t, got, expected)
}

// asserts that got == expected
// ignores expected.UpdatedAt, but ensures that got.UpdatedAt is set and not null
func AssertStudioPartial(t assert.TestingT, got, expected models.StudioPartial) bool {
	// updated at should be set and not null
	if !got.UpdatedAt.Set || got.UpdatedAt.Null {
		return false
	}
	// else ignore the exact value
	got.UpdatedAt = models.OptionalTime{}

	return assert.Equal(t, got, expected)
}

// asserts that got == expected
// ignores expected.UpdatedAt, but ensures that got.UpdatedAt is set and not null
func AssertTagPartial(t assert.TestingT, got, expected models.TagPartial) bool {
	// updated at should be set and not null
	if !got.UpdatedAt.Set || got.UpdatedAt.Null {
		return false
	}
	// else ignore the exact value
	got.UpdatedAt = models.OptionalTime{}

	return assert.Equal(t, got, expected)
}
