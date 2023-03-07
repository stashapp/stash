package ptrs

func Strptr(v string) *string {
	p := new(string)
	*p = v
	return p
}

func Intptr(v int) *int {
	p := new(int)
	*p = v
	return p
}

func Int64ptr(v int64) *int64 {
	p := new(int64)
	*p = v
	return p
}

func Uintptr(v uint) *uint {
	p := new(uint)
	*p = v
	return p
}

func Uint32ptr(v uint32) *uint32 {
	p := new(uint32)
	*p = v
	return p
}

func Uint64ptr(v uint64) *uint64 {
	p := new(uint64)
	*p = v
	return p
}

func Boolptr(v bool) *bool {
	p := new(bool)
	*p = v
	return p
}

func Float64ptr(v float64) *float64 {
	p := new(float64)
	*p = v
	return p
}
