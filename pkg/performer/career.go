package performer

import (
	"github.com/stashapp/stash/pkg/utils"
)

// ParseCareerLength parses a career_length string into start and end year integers.
// Deprecated: Use utils.ParseYearRangeString directly.
func ParseCareerLength(s string) (start *int, end *int, err error) {
	return utils.ParseYearRangeString(s)
}
