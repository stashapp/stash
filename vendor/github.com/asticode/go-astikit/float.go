package astikit

import (
	"bytes"
	"fmt"
	"strconv"
)

// Rational represents a rational
type Rational struct{ den, num int }

// NewRational creates a new rational
func NewRational(num, den int) *Rational {
	return &Rational{
		den: den,
		num: num,
	}
}

// Num returns the rational num
func (r *Rational) Num() int {
	return r.num
}

// Den returns the rational den
func (r *Rational) Den() int {
	return r.den
}

// ToFloat64 returns the rational as a float64
func (r *Rational) ToFloat64() float64 {
	return float64(r.num) / float64(r.den)
}

// MarshalText implements the TextMarshaler interface
func (r *Rational) MarshalText() (b []byte, err error) {
	b = []byte(fmt.Sprintf("%d/%d", r.num, r.den))
	return
}

// UnmarshalText implements the TextUnmarshaler interface
func (r *Rational) UnmarshalText(b []byte) (err error) {
	r.num = 0
	r.den = 1
	if len(b) == 0 {
		return
	}
	items := bytes.Split(b, []byte("/"))
	if r.num, err = strconv.Atoi(string(items[0])); err != nil {
		err = fmt.Errorf("astikit: atoi of %s failed: %w", string(items[0]), err)
		return
	}
	if len(items) > 1 {
		if r.den, err = strconv.Atoi(string(items[1])); err != nil {
			err = fmt.Errorf("astifloat: atoi of %s failed: %w", string(items[1]), err)
			return
		}
	}
	return
}
