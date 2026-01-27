package sqlite

import (
	"database/sql/driver"
	"time"

	"github.com/stashapp/stash/pkg/models"
	"gopkg.in/guregu/null.v4"
)

const sqliteDateLayout = "2006-01-02"

// Date represents a date stored as "YYYY-MM-DD"
type Date struct {
	Date time.Time
}

// Scan implements the Scanner interface.
func (d *Date) Scan(value interface{}) error {
	d.Date = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (d Date) Value() (driver.Value, error) {
	return d.Date.Format(sqliteDateLayout), nil
}

// NullDate represents a nullable date stored as "YYYY-MM-DD"
type NullDate struct {
	Date  time.Time
	Valid bool
}

// Scan implements the Scanner interface.
func (d *NullDate) Scan(value interface{}) error {
	var ok bool
	d.Date, ok = value.(time.Time)
	if !ok {
		d.Date = time.Time{}
		d.Valid = false
		return nil
	}

	d.Valid = true
	return nil
}

// Value implements the driver Valuer interface.
func (d NullDate) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}

	return d.Date.Format(sqliteDateLayout), nil
}

func (d *NullDate) DatePtr(precision null.Int) *models.Date {
	if d == nil || !d.Valid {
		return nil
	}

	return &models.Date{Time: d.Date, Precision: models.DatePrecision(precision.Int64)}
}

func NullDateFromDatePtr(d *models.Date) NullDate {
	if d == nil {
		return NullDate{Valid: false}
	}
	return NullDate{Date: d.Time, Valid: true}
}

func datePrecisionFromDatePtr(d *models.Date) null.Int {
	if d == nil {
		// default to day precision
		return null.Int{}
	}
	return null.IntFrom(int64(d.Precision))
}
