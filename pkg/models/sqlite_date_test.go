package models

import (
	"database/sql/driver"
	"reflect"
	"testing"
)

func TestSQLiteDate_Value(t *testing.T) {
	tests := []struct {
		name    string
		tr      string
		want    driver.Value
		wantErr bool
	}{
		{
			"empty string",
			"",
			"",
			false,
		},
		{
			"whitespace",
			" ",
			"",
			false,
		},
		{
			"RFC3339",
			"2021-11-22T17:11:55+11:00",
			"2021-11-22",
			false,
		},
		{
			"date",
			"2021-11-22",
			"2021-11-22",
			false,
		},
		{
			"date and time",
			"2021-11-22 17:12:05",
			"2021-11-22",
			false,
		},
		{
			"date, time and zone",
			"2021-11-22 17:33:05 AEST",
			"2021-11-22",
			false,
		},
		{
			"whitespaced date",
			"  2021-11-22 ",
			"2021-11-22",
			false,
		},
		{
			"invalid",
			"foo",
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := SQLiteDate{
				String: tt.tr,
				Valid:  true,
			}
			got, err := d.Value()
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
