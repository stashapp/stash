package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIfDateHigherThanAnotherDate(t *testing.T) {
	type test struct {
		fromDate string
		toDate   string
		expected bool
	}
	tests := []test{
		{fromDate: "1985-01-01", toDate: "2021-04-12", expected: false},
		{fromDate: "1900-02-02", toDate: "2000-02-02", expected: false},
		{fromDate: "2021-01-01", toDate: "1923-04-12", expected: true},
		{fromDate: "2020-01-01", toDate: "1982-04-12", expected: true},
	}

	assert := assert.New(t)
	for i, tc := range tests {
		result := IfDateHigherThanAnotherDate(tc.fromDate, tc.toDate)
		assert.Equal(tc.expected, result, "[%d] expected: %t fromDate: %s; toDate: %s", i, tc.expected, tc.fromDate, tc.toDate)
	}
}
