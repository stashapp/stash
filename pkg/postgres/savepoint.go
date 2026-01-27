package postgres

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/hash"
)

type SavepointAction func(ctx context.Context) error

const savePointPrefix = "savepoint_" // prefix for savepoint

// Encapsulates an action in a savepoint
// Its mostly used to rollback if an error occurred in postgres, as errors in postgres cancel the transaction.
func withSavepoint(ctx context.Context, action SavepointAction) error {
	tx, err := getTx(ctx)
	if err != nil {
		return err
	}

	// Generate savepoint
	rnd, err := hash.GenerateRandomKey(64)
	if err != nil {
		return err
	}

	// Sqlite needs some letters infront of the identifier
	rnd = savePointPrefix + rnd

	// Create a savepoint
	_, err = tx.Exec("SAVEPOINT " + rnd)
	if err != nil {
		return fmt.Errorf("failed to create savepoint: %w", err)
	}

	// Execute the action
	err = action(ctx)
	if err != nil {
		// Rollback to savepoint on error
		if _, rbErr := tx.Exec("ROLLBACK TO SAVEPOINT " + rnd); rbErr != nil {
			return fmt.Errorf("action failed and rollback to savepoint failed: %w", rbErr)
		}
		return fmt.Errorf("action failed: %w", err)
	}

	// Release the savepoint on success
	_, err = tx.Exec("RELEASE SAVEPOINT " + rnd)
	if err != nil {
		return fmt.Errorf("failed to release savepoint: %w", err)
	}

	return nil
}
