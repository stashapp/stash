package sqlite

import (
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/stashapp/stash/pkg/models"
	"gopkg.in/guregu/null.v4"
	"gopkg.in/guregu/null.v4/zero"
)

type updateRecord struct {
	exp.Record
}

func (r *updateRecord) set(destField string, v interface{}) {
	r.Record[destField] = v
}

func (r *updateRecord) setString(destField string, v models.OptionalString) {
	if v.Set {
		if v.Null {
			panic("null value not allowed in optional string")
		}
		r.set(destField, v.Value)
	}
}

func (r *updateRecord) setNullString(destField string, v models.OptionalString) {
	if v.Set {
		r.set(destField, zero.StringFromPtr(v.Ptr()))
	}
}

func (r *updateRecord) setBool(destField string, v models.OptionalBool) {
	if v.Set {
		if v.Null {
			panic("null value not allowed in optional bool")
		}
		r.set(destField, v.Value)
	}
}

func (r *updateRecord) setInt(destField string, v models.OptionalInt) {
	if v.Set {
		if v.Null {
			panic("null value not allowed in optional int")
		}
		r.set(destField, v.Value)
	}
}

func (r *updateRecord) setNullInt(destField string, v models.OptionalInt) {
	if v.Set {
		r.set(destField, intFromPtr(v.Ptr()))
	}
}

// func (r *updateRecord) setInt64(destField string, v models.OptionalInt64) {
// 	if v.Set {
// 		if v.Null {
// 			panic("null value not allowed in optional int64")
// 		}
// 		r.set(destField, v.Value)
// 	}
// }

// func (r *updateRecord) setNullInt64(destField string, v models.OptionalInt64) {
// 	if v.Set {
// 		r.set(destField, null.IntFromPtr(v.Ptr()))
// 	}
// }

func (r *updateRecord) setFloat64(destField string, v models.OptionalFloat64) {
	if v.Set {
		if v.Null {
			panic("null value not allowed in optional float64")
		}
		r.set(destField, v.Value)
	}
}

func (r *updateRecord) setNullFloat64(destField string, v models.OptionalFloat64) {
	if v.Set {
		r.set(destField, null.FloatFromPtr(v.Ptr()))
	}
}

func (r *updateRecord) setTimestamp(destField string, v models.OptionalTime) {
	if v.Set {
		if v.Null {
			panic("null value not allowed in optional time")
		}
		r.set(destField, Timestamp{Timestamp: v.Value})
	}
}

//nolint:golint,unused
func (r *updateRecord) setNullTimestamp(destField string, v models.OptionalTime) {
	if v.Set {
		r.set(destField, NullTimestampFromTimePtr(v.Ptr()))
	}
}

func (r *updateRecord) setNullDate(destField string, v models.OptionalDate) {
	if v.Set {
		r.set(destField, NullDateFromDatePtr(v.Ptr()))
	}
}
