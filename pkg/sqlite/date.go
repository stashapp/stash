package sqlite

import (
	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/models"
	"gopkg.in/guregu/null.v4"
)

// Date represents a date stored as "YYYY-MM-DD"
type Date = database.Date

// NullDate represents a nullable date stored as "YYYY-MM-DD"
type NullDate = database.NullDate

var NullDateFromDatePtr = database.NullDateFromDatePtr

func datePrecisionFromDatePtr(d *models.Date) null.Int {
	return database.DatePrecisionFromDatePtr(d)
}
