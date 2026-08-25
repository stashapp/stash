package audio

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/stashapp/stash/pkg/models"
)

var ErrEmptyUpdater = errors.New("no fields have been set")

// UpdateSet is used to update a audio and its relationships.
type UpdateSet struct {
	ID int

	Partial models.AudioPartial

	// in future these could be moved into a separate struct and reused
	// for a Creator struct
}

// IsEmpty returns true if there is nothing to update.
func (u *UpdateSet) IsEmpty() bool {
	withoutID := u.Partial

	return withoutID == models.AudioPartial{}
}

// Update updates a audio by updating the fields in the Partial field, then
// updates non-nil relationships. Returns an error if there is no work to
// be done.
func (u *UpdateSet) Update(ctx context.Context, qb models.AudioUpdater) (*models.Audio, error) {
	if u.IsEmpty() {
		return nil, ErrEmptyUpdater
	}

	partial := u.Partial
	updatedAt := time.Now()
	partial.UpdatedAt = models.NewOptionalTime(updatedAt)

	ret, err := qb.UpdatePartial(ctx, u.ID, partial)
	if err != nil {
		return nil, fmt.Errorf("error updating audio: %w", err)
	}

	return ret, nil
}

// UpdateInput converts the UpdateSet into AudioUpdateInput for hook firing purposes.
func (u UpdateSet) UpdateInput() models.AudioUpdateInput {
	// ensure the partial ID is set
	ret := u.Partial.UpdateInput(u.ID)

	return ret
}

func AddPerformer(ctx context.Context, qb models.AudioUpdater, o *models.Audio, performerID int) error {
	audioPartial := models.NewAudioPartial()
	audioPartial.PerformerIDs = &models.UpdateIDs{
		IDs:  []int{performerID},
		Mode: models.RelationshipUpdateModeAdd,
	}
	_, err := qb.UpdatePartial(ctx, o.ID, audioPartial)
	return err
}

func AddTag(ctx context.Context, qb models.AudioUpdater, o *models.Audio, tagID int) error {
	audioPartial := models.NewAudioPartial()
	audioPartial.TagIDs = &models.UpdateIDs{
		IDs:  []int{tagID},
		Mode: models.RelationshipUpdateModeAdd,
	}
	_, err := qb.UpdatePartial(ctx, o.ID, audioPartial)
	return err
}

func (s *Service) AssignFile(ctx context.Context, audioID int, fileID models.FileID) error {
	// ensure file isn't a primary file and that it is an audio file
	f, err := s.File.Find(ctx, fileID)
	if err != nil {
		return err
	}

	if len(f) == 0 {
		return fmt.Errorf("file with id %d not found", fileID)
	}

	ff := f[0]
	if _, ok := ff.(*models.AudioFile); !ok {
		return fmt.Errorf("%s is not an audio file", ff.Base().Path)
	}

	isPrimary, err := s.File.IsPrimary(ctx, fileID)
	if err != nil {
		return err
	}

	if isPrimary {
		return errors.New("cannot reassign primary file")
	}

	return s.Repository.AssignFiles(ctx, audioID, []models.FileID{fileID})
}
