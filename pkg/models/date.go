package models

import "time"

type Date struct {
	time.Time
}

const dateFormat = "2006-01-02"

func (d Date) String() string {
	return d.Format(dateFormat)
}

func NewDate(s string) Date {
	t, _ := time.Parse(dateFormat, s)
	return Date{t}
}
