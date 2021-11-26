package models

import (
	"database/sql/driver"
	"reflect"
	"testing"
)

func TestSQLiteDate_Value(t *testing.T) {
	tests := []struct {
		name    string
		tr      SQLiteDate
		want    driver.Value
		wantErr bool
	}{
		{
			"empty string",
			SQLiteDate{"", true},
			"",
			false,
		},
		{
			"whitespace",
			SQLiteDate{" ", true},
			"",
			false,
		},
		{
			"RFC3339",
			SQLiteDate{"2021-11-22T17:11:55+11:00", true},
			"2021-11-22",
			false,
		},
		{
			"date",
			SQLiteDate{"2021-11-22", true},
			"2021-11-22",
			false,
		},
		{
			"date and time",
			SQLiteDate{"2021-11-22 17:12:05", true},
			"2021-11-22",
			false,
		},
		{
			"date, time and zone",
			SQLiteDate{"2021-11-22 17:33:05 AEST", true},
			"2021-11-22",
			false,
		},
		{
			"whitespaced date",
			SQLiteDate{"  2021-11-22 ", true},
			"2021-11-22",
			false,
		},
		{
			"bad format",
			SQLiteDate{"foo", true},
			nil,
			true,
		},
		{
			"invalid",
			SQLiteDate{"null", false},
			nil,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tr.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLiteDate.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SQLiteDate.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
