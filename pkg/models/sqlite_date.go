package models

import (
	"database/sql/driver"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/utils"
	"time"
)

type SQLiteDate struct {
	String string
	Valid  bool
}

// Scan implements the Scanner interface.
func (t *SQLiteDate) Scan(value interface{}) error {
	dateTime, ok := value.(time.Time)
	if !ok {
		t.String = ""
		t.Valid = false
		return nil
	}

	t.String = dateTime.Format("2006-01-02")
	if t.String != "" && t.String != "0001-01-01" {
		t.Valid = true
	} else {
		t.Valid = false
	}
	return nil
}

// Value implements the driver Valuer interface.
func (t SQLiteDate) Value() (driver.Value, error) {
	result, err := utils.ParseDateStringAsFormat(t.String, "2006-01-02")
	if err != nil {
		logger.Debugf("sqlite date conversion error: %s", err.Error())
	}
	return result, nil
}
