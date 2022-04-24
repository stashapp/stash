package sqlite

import (
	"time"

	"github.com/doug-martin/goqu/v9/exp"
	"github.com/stashapp/stash/pkg/models"
)

type updateRecord struct {
	exp.Record
}

func (r *updateRecord) set(destField string, v interface{}) {
	r.Record[destField] = v
}

func (r *updateRecord) setString(destField string, v *string) {
	if v != nil {
		r.set(destField, *v)
	}
}

// TODO - rename to setString
func (r *updateRecord) setOptionalString(destField string, v models.OptionalString) {
	if v.Set {
		if v.Null {
			panic("null value not allowed in optional string")
		}
		r.set(destField, v.Value)
	}
}

// TODO - rename to setNullString
func (r *updateRecord) setOptionalNullString(destField string, v models.OptionalString) {
	if v.Set {
		r.set(destField, newNullStringPtr(v.Ptr()))
	}
}

func (r *updateRecord) setBool(destField string, v *bool) {
	if v != nil {
		r.set(destField, *v)
	}
}

// TODO - rename to setBool
func (r *updateRecord) setOptionalBool(destField string, v models.OptionalBool) {
	if v.Set {
		if v.Null {
			panic("null value not allowed in optional int")
		}
		r.set(destField, v.Value)
	}
}

// TODO - rename to setInt
func (r *updateRecord) setOptionalInt(destField string, v models.OptionalInt) {
	if v.Set {
		if v.Null {
			panic("null value not allowed in optional int")
		}
		r.set(destField, v.Value)
	}
}

// TODO - rename to setNullInt
func (r *updateRecord) setOptionalNullInt(destField string, v models.OptionalInt) {
	if v.Set {
		r.set(destField, newNullIntPtr(v.Ptr()))
	}
}

func (r *updateRecord) setNullStringPtr(destField string, v **string) {
	if v != nil {
		r.set(destField, newNullStringPtr(*v))
	}
}

func (r *updateRecord) setNullString(destField string, v *string) {
	if v != nil {
		r.set(destField, newNullString(*v))
	}
}

func (r *updateRecord) setNullIntPtr(destField string, v **int) {
	if v != nil {
		r.set(destField, newNullIntPtr(*v))
	}
}

// TODO - rename to setInt64
// func (r *updateRecord) setOptionalInt64(destField string, v models.OptionalInt64) {
// 	if v.Set {
// 		if v.Null {
// 			panic("null value not allowed in optional int64")
// 		}
// 		r.set(destField, v.Value)
// 	}
// }

// TODO - rename to setNullInt64
func (r *updateRecord) setOptionalNullInt64(destField string, v models.OptionalInt64) {
	if v.Set {
		r.set(destField, newNullInt64Ptr(v.Ptr()))
	}
}

// TODO - rename to setFloat64
// func (r *updateRecord) setOptionalFloat64(destField string, v models.OptionalFloat64) {
// 	if v.Set {
// 		if v.Null {
// 			panic("null value not allowed in optional float64")
// 		}
// 		r.set(destField, v.Value)
// 	}
// }

// TODO - rename to setNullFloat64
func (r *updateRecord) setOptionalNullFloat64(destField string, v models.OptionalFloat64) {
	if v.Set {
		r.set(destField, newNullFloat64Ptr(v.Ptr()))
	}
}

func (r *updateRecord) setNullTimePtr(destField string, v **time.Time) {
	if v != nil {
		r.set(destField, newNullTime(*v))
	}
}

// TODO - rename to setTime
func (r *updateRecord) setOptionalTime(destField string, v models.OptionalTime) {
	if v.Set {
		if v.Null {
			panic("null value not allowed in optional time")
		}
		r.set(destField, v.Value)
	}
}

// TODO - rename to setNullTime
func (r *updateRecord) setOptionalNullTime(destField string, v models.OptionalTime) {
	if v.Set {
		r.set(destField, newNullTime(v.Ptr()))
	}
}

func (r *updateRecord) setTime(destField string, v *time.Time) {
	if v != nil {
		r.set(destField, *v)
	}
}

func (r *updateRecord) setOptionalSQLiteDate(destField string, v models.OptionalDate) {
	if v.Set {
		if v.Null {
			r.set(destField, models.SQLiteDate{})
		}

		r.set(destField, models.SQLiteDate{
			String: v.Value.String(),
			Valid:  true,
		})
	}
}
