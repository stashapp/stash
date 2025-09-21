//go:build db_integration
// +build db_integration

package database_test

import (
	"context"
	"os"
	"testing"

	"github.com/stashapp/stash/pkg/postgres"
	"github.com/stashapp/stash/pkg/sqlite"
)

func IsPostgresTest() *string {
	if val, ok := os.LookupEnv("PGSQL_TEST"); ok {
		return &val
	}
	return nil
}

func TestAnonymiser_Anonymise(t *testing.T) {
	f, err := os.CreateTemp("", "*.sqlite")
	if err != nil {
		t.Errorf("Could not create temporary file: %v", err)
		return
	}

	f.Close()
	defer os.Remove(f.Name())

	// use existing database
	var anonymiser *sqlite.Anonymiser
	if val := IsPostgresTest(); val != nil {
		anonymiser, err = postgres.NewAnonymiser(db.(*postgres.Database), f.Name())
	} else {
		anonymiser, err = sqlite.NewAnonymiser(db.(*sqlite.Database), f.Name())
	}

	if err != nil {
		t.Errorf("Could not create anonymiser: %v", err)
		return
	}

	if err := anonymiser.Anonymise(context.Background()); err != nil {
		t.Errorf("Could not anonymise: %v", err)
		return
	}

	t.Logf("Anonymised database written to %s", f.Name())

	// TODO - ensure anonymous
}
