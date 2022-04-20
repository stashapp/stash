package sqlite

import (
	"database/sql"
	"time"
)

type nullString struct {
	sql.NullString
}

func (n nullString) stringPtr() *string {
	if !n.Valid {
		return nil
	}

	return &n.NullString.String
}

func newNullStringPtr(v *string) nullString {
	var vv string
	valid := v != nil && *v != ""
	if valid {
		vv = *v
	}

	return nullString{
		NullString: sql.NullString{
			String: vv,
			Valid:  valid,
		},
	}
}

// func newNullString(v string) nullString {
// 	return newNullStringPtr(&v)
// }

type nullInt64 struct {
	sql.NullInt64
}

func (n nullInt64) int64Ptr() *int64 {
	if !n.Valid {
		return nil
	}

	return &n.NullInt64.Int64
}

// func (n nullInt64) int32Ptr() *int32 {
// 	if !n.Valid {
// 		return nil
// 	}

// 	v := int32(n.NullInt64.Int64)

// 	return &v
// }

// func (n nullInt64) intPtr() *int {
// 	if !n.Valid {
// 		return nil
// 	}

// 	v := int(n.NullInt64.Int64)

// 	return &v
// }

// func newNullInt64(i int64) nullInt64 {
// 	return nullInt64{
// 		NullInt64: sql.NullInt64{
// 			Int64: i,
// 			Valid: i != 0,
// 		},
// 	}
// }

func newNullInt64Ptr(i *int64) nullInt64 {
	ret := nullInt64{
		NullInt64: sql.NullInt64{
			Valid: i != nil,
		},
	}

	if ret.Valid {
		ret.Int64 = *i
	}

	return ret
}

// type nullInt32 struct {
// 	sql.NullInt32
// }

// func (n nullInt32) int32Ptr() *int32 {
// 	if !n.Valid {
// 		return nil
// 	}

// 	return &n.NullInt32.Int32
// }

// func newNullInt32(i int32) nullInt32 {
// 	return nullInt32{
// 		NullInt32: sql.NullInt32{
// 			Int32: i,
// 			Valid: true,
// 		},
// 	}
// }

// func newNullInt32Ptr(i *int32) nullInt32 {
// 	ret := nullInt32{
// 		NullInt32: sql.NullInt32{
// 			Valid: i != nil,
// 		},
// 	}

// 	if ret.Valid {
// 		ret.Int32 = *i
// 	}

// 	return ret
// }

type nullInt struct {
	sql.NullInt64
}

func (n nullInt) intPtr() *int {
	if !n.Valid {
		return nil
	}

	v := int(n.NullInt64.Int64)

	return &v
}

// func (n nullInt) int() int {
// 	return int(n.NullInt64.Int64)
// }

func newNullIntPtr(i *int) nullInt {
	ret := nullInt{
		NullInt64: sql.NullInt64{
			Valid: i != nil,
		},
	}

	if ret.Valid {
		ret.Int64 = int64(*i)
	}

	return ret
}

// func newNullInt(i int) nullInt {
// 	return newNullIntPtr(&i)
// }

type nullTime struct {
	sql.NullTime
}

func (n nullTime) timePtr() *time.Time {
	if !n.Valid {
		return nil
	}

	return &n.NullTime.Time
}

func newNullTime(v *time.Time) nullTime {
	var vv time.Time
	if v != nil {
		vv = *v
	}

	return nullTime{
		NullTime: sql.NullTime{
			Time:  vv,
			Valid: v != nil,
		},
	}
}

// type nullFloat64 struct {
// 	sql.NullFloat64
// }
