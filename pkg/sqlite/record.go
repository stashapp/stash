package sqlite

import (
	"time"

	"github.com/doug-martin/goqu/v9/exp"
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

func (r *updateRecord) setBool(destField string, v *bool) {
	if v != nil {
		r.set(destField, *v)
	}
}

func (r *updateRecord) setInt(destField string, v *int) {
	if v != nil {
		r.set(destField, *v)
	}
}

func (r *updateRecord) setNullStringPtr(destField string, v **string) {
	if v != nil {
		r.set(destField, newNullStringPtr(*v))
	}
}

func (r *updateRecord) setNullIntPtr(destField string, v **int) {
	if v != nil {
		r.set(destField, newNullIntPtr(*v))
	}
}

func (r *updateRecord) setNullInt64Ptr(destField string, v **int64) {
	if v != nil {
		r.set(destField, newNullInt64Ptr(*v))
	}
}

func (r *updateRecord) setNullTimePtr(destField string, v **time.Time) {
	if v != nil {
		r.set(destField, newNullTime(*v))
	}
}

func (r *updateRecord) setTime(destField string, v *time.Time) {
	if v != nil {
		r.set(destField, *v)
	}
}
