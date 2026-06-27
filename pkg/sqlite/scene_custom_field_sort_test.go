package sqlite

import (
	"strings"
	"testing"
)

func TestValidateCustomFieldSortName(t *testing.T) {
	valid := []string{
		"smobile_rec_score",
		"label",
		"a",
		strings.Repeat("a", 64),
		"field_123",
		"A1_b",
	}
	for _, name := range valid {
		if err := validateCustomFieldSortName(name); err != nil {
			t.Errorf("expected %q to be valid, got: %v", name, err)
		}
	}

	invalid := []struct {
		name string
		desc string
	}{
		{"", "empty"},
		{strings.Repeat("a", 65), "65 chars (over limit)"},
		{"a;drop", "semicolon"},
		{"field name", "space"},
		{"field-name", "hyphen"},
		{"field.name", "dot"},
		{"field'name", "single quote"},
		{"field\"name", "double quote"},
	}
	for _, tc := range invalid {
		if err := validateCustomFieldSortName(tc.name); err == nil {
			t.Errorf("expected %q (%s) to be invalid but got no error", tc.name, tc.desc)
		}
	}
}

func TestValidateCustomFieldSortCast(t *testing.T) {
	valid := []string{"REAL", "INTEGER", "TEXT", "NUMERIC"}
	for _, cast := range valid {
		if err := validateCustomFieldSortCast(cast); err != nil {
			t.Errorf("expected %q to be valid, got: %v", cast, err)
		}
	}

	invalid := []string{"FLOAT", "DOUBLE", "VARCHAR", "BLOB", "INT", "real", "integer"}
	for _, cast := range invalid {
		if err := validateCustomFieldSortCast(cast); err == nil {
			t.Errorf("expected %q to be invalid but got no error", cast)
		}
	}
}

func TestSetCustomFieldSortValidation(t *testing.T) {
	qb := &SceneStore{} // zero-value safe: setCustomFieldSort must not require a live DB handle

	rejectCases := []struct {
		spec string
		desc string
	}{
		{"smobile_rec_score", "missing cast and direction"},
		{"smobile_rec_score:REAL", "missing direction"},
		{":REAL:DESC", "empty field name"},
		{"smobile_rec_score:FLOAT:DESC", "invalid cast"},
		{"smobile_rec_score:REAL:UP", "invalid direction"},
		{"smobile_rec_score:REAL:DOWN", "invalid direction DOWN"},
		{"a;drop:REAL:DESC", "dangerous field name"},
		{strings.Repeat("a", 65) + ":REAL:DESC", "field name too long"},
	}
	for _, tc := range rejectCases {
		qRej := &queryBuilder{}
		if err := qb.setCustomFieldSort(qRej, tc.spec); err == nil {
			t.Errorf("expected spec %q (%s) to be rejected, but got no error", tc.spec, tc.desc)
		}
	}

	acceptCases := []struct {
		spec string
		desc string
	}{
		{"smobile_rec_score:REAL:DESC", "numeric score descending"},
		{"label:TEXT:ASC", "text label ascending"},
		{"count:INTEGER:ASC", "integer count ascending"},
		{"score:NUMERIC:DESC", "numeric type descending"},
	}
	for _, tc := range acceptCases {
		q2 := &queryBuilder{}
		if err := qb.setCustomFieldSort(q2, tc.spec); err != nil {
			t.Errorf("expected spec %q (%s) to be accepted, got: %v", tc.spec, tc.desc, err)
		}
	}
}
