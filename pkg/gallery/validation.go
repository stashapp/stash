package gallery

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil"
)

type ContentsChangedError struct {
	Gallery *models.Gallery
}

func (e *ContentsChangedError) Error() string {
	typ := "zip-based"
	if e.Gallery.FolderID != nil {
		typ = "folder-based"
	}

	return fmt.Sprintf("cannot change contents of %s gallery %q", typ, e.Gallery.GetTitle())
}

// validateContentChange returns an error if a gallery cannot have its contents changed.
// Only manually created galleries can have images changed.
func validateContentChange(g *models.Gallery) error {
	if g.FolderID != nil || g.PrimaryFileID != nil {
		return &ContentsChangedError{
			Gallery: g,
		}
	}

	return nil
}

func (s *Service) ValidateImageGalleryChange(ctx context.Context, i *models.Image, updateIDs models.UpdateIDs) error {
	// determine what is changing
	var changedIDs []int

	switch updateIDs.Mode {
	case models.RelationshipUpdateModeAdd, models.RelationshipUpdateModeRemove:
		changedIDs = updateIDs.IDs
	case models.RelationshipUpdateModeSet:
		// get the difference between the two lists
		changedIDs = sliceutil.NotIntersect(i.GalleryIDs.List(), updateIDs.IDs)
	}

	galleries, err := s.Repository.FindMany(ctx, changedIDs)
	if err != nil {
		return err
	}

	for _, g := range galleries {
		if err := validateContentChange(g); err != nil {
			return fmt.Errorf("changing galleries of image %q: %w", i.GetTitle(), err)
		}
	}

	return nil
}
