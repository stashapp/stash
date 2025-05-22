package autotag

import (
	"context"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/stashapp/stash/pkg/models"
)

// Regular expressions to recognize various date patterns
var datePatterns = []*regexp.Regexp{
	// YYYY-MM-DD format (2020-11-12)
	regexp.MustCompile(`(\d{4})-(\d{1,2})-(\d{1,2})`),
	// YYYYMMDD format (20201112)
	regexp.MustCompile(`(\d{4})(\d{2})(\d{2})`),
	// DD.MM.YYYY format (12.11.2020)
	regexp.MustCompile(`(\d{1,2})\.(\d{1,2})\.(\d{4})`),
}

// extracts a date from a file path.
func ExtractDateFromPath(path string) *time.Time {
	// Extract filename only
	filename := filepath.Base(path)
	
	// Attempt each pattern
	for i, pattern := range datePatterns {
		matches := pattern.FindStringSubmatch(filename)
		if len(matches) >= 4 {
			var year, month, day int
			var err error
			
			// Implement date parsing logic based on the pattern
			switch i {
			case 0, 1: // YYYY-MM-DD or YYYYMMDD
				year, err = strconv.Atoi(matches[1])
				if err != nil {
					continue
				}
				month, err = strconv.Atoi(matches[2])
				if err != nil {
					continue
				}
				day, err = strconv.Atoi(matches[3])
				if err != nil {
					continue
				}
			case 2: // DD.MM.YYYY
				day, err = strconv.Atoi(matches[1])
				if err != nil {
					continue
				}
				month, err = strconv.Atoi(matches[2])
				if err != nil {
					continue
				}
				year, err = strconv.Atoi(matches[3])
				if err != nil {
					continue
				}
			}
			
			// Validate the date
			if year < 1900 || year > 2100 || month < 1 || month > 12 || day < 1 || day > 31 {
				continue
			}
			
			// Create date object
			date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
			return &date
		}
	}
	
	return nil
}

// extracts and sets the date from a scene's file path.
func SceneDate(ctx context.Context, s *models.Scene, rw models.SceneUpdater) error {
	// Skip if the date is already set
	if s.Date != nil {
		return nil
	}
	
	// Extract date from file path
	date := ExtractDateFromPath(s.Path)
	if date == nil {
		return nil // Date not found
	}
	
	// Update scene object
	partial := models.NewScenePartial()
	
	// Convert time.Time to models.Date
	dateModel := models.Date{Time: *date}
	partial.Date = models.NewOptionalDate(dateModel)
	
	// Update the database
	_, err := rw.UpdatePartial(ctx, s.ID, partial)
	return err
}
