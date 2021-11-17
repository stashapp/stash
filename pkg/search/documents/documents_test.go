package documents

import (
	"database/sql"
	"testing"

	"github.com/stashapp/stash/pkg/models"
)

func TestNewScene(t *testing.T) {
	have := models.Scene{
		ID: 123,
		Title: sql.NullString{
			Valid:  true,
			String: "Dark Elves having fun",
		},
		Date: models.SQLiteDate{
			Valid:  true,
			String: "2021-10-02",
		},
	}

	got := NewScene(have, nil, nil, nil)

	if have.Title.String != got.Title {
		t.Errorf("[title] want: %v; got: %v", have.Title, got.Title)
	}

	if got.Year == nil {
		t.Errorf("nil year")
	}

	expectYear := 2021
	if got.Year != nil && *got.Year != expectYear {
		t.Errorf("[year] want: %v; got: %v", expectYear, *got.Year)
	}
}
