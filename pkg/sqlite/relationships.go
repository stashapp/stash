package sqlite

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

type idRelationshipStore struct {
	joinTable *joinTable
}

func (s *idRelationshipStore) createRelationships(ctx context.Context, id int, fkIDs models.RelatedIDs) error {
	if fkIDs.Loaded() {
		if err := s.joinTable.insertJoins(ctx, id, fkIDs.List()); err != nil {
			return err
		}
	}

	return nil
}

func (s *idRelationshipStore) modifyRelationships(ctx context.Context, id int, fkIDs *models.UpdateIDs) error {
	if fkIDs != nil {
		if err := s.joinTable.modifyJoins(ctx, id, fkIDs.IDs, fkIDs.Mode); err != nil {
			return err
		}
	}

	return nil
}

func (s *idRelationshipStore) replaceRelationships(ctx context.Context, id int, fkIDs models.RelatedIDs) error {
	if fkIDs.Loaded() {
		if err := s.joinTable.replaceJoins(ctx, id, fkIDs.List()); err != nil {
			return err
		}
	}

	return nil
}
