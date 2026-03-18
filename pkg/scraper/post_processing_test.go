package scraper

import (
	"context"
	"testing"

	"github.com/stashapp/stash/pkg/models"
)

func TestPostScrapePerformerCareerLength(t *testing.T) {
	ctx := context.Background()
	const related = false

	strPtr := func(s string) *string {
		return &s
	}

	tests := []struct {
		name  string
		input models.ScrapedPerformer
		want  models.ScrapedPerformer
	}{
		{
			"start = 2000",
			models.ScrapedPerformer{
				CareerStart: strPtr("2000"),
			},
			models.ScrapedPerformer{
				CareerStart:  strPtr("2000"),
				CareerLength: strPtr("2000 -"),
			},
		},
		{
			"end = 2000",
			models.ScrapedPerformer{
				CareerEnd: strPtr("2000"),
			},
			models.ScrapedPerformer{
				CareerEnd:    strPtr("2000"),
				CareerLength: strPtr("- 2000"),
			},
		},
		{
			"start = 2000, end = 2020",
			models.ScrapedPerformer{
				CareerStart: strPtr("2000"),
				CareerEnd:   strPtr("2020"),
			},
			models.ScrapedPerformer{
				CareerStart:  strPtr("2000"),
				CareerEnd:    strPtr("2020"),
				CareerLength: strPtr("2000 - 2020"),
			},
		},
		{
			"length = 2000 -",
			models.ScrapedPerformer{
				CareerLength: strPtr("2000 -"),
			},
			models.ScrapedPerformer{
				CareerStart:  strPtr("2000"),
				CareerLength: strPtr("2000 -"),
			},
		},
		{
			"length = - 2010",
			models.ScrapedPerformer{
				CareerLength: strPtr("- 2010"),
			},
			models.ScrapedPerformer{
				CareerEnd:    strPtr("2010"),
				CareerLength: strPtr("- 2010"),
			},
		},
		{
			"length = 2000 - 2010",
			models.ScrapedPerformer{
				CareerLength: strPtr("2000 - 2010"),
			},
			models.ScrapedPerformer{
				CareerStart:  strPtr("2000"),
				CareerEnd:    strPtr("2010"),
				CareerLength: strPtr("2000 - 2010"),
			},
		},
		{
			"invalid start",
			models.ScrapedPerformer{
				CareerStart: strPtr("two thousand"),
			},
			models.ScrapedPerformer{
				CareerStart: strPtr("two thousand"),
			},
		},
		{
			"invalid end",
			models.ScrapedPerformer{
				CareerEnd: strPtr("two thousand"),
			},
			models.ScrapedPerformer{
				CareerEnd: strPtr("two thousand"),
			},
		},
		{
			"invalid career length",
			models.ScrapedPerformer{
				CareerLength: strPtr("1234 - 4567 - 9224"),
			},
			models.ScrapedPerformer{
				CareerLength: strPtr("1234 - 4567 - 9224"),
			},
		},
	}

	compareStrPtr := func(a, b *string) bool {
		if a == b {
			return true
		}
		if a == nil || b == nil {
			return false
		}
		return *a == *b
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &postScraper{}
			got, err := c.postScrapePerformer(ctx, tt.input, related)
			if err != nil {
				t.Fatalf("postScrapePerformer returned error: %v", err)
			}
			postScraped := got.(models.ScrapedPerformer)
			if !compareStrPtr(postScraped.CareerStart, tt.want.CareerStart) {
				t.Errorf("CareerStart = %v, want %v", postScraped.CareerStart, tt.want.CareerStart)
			}
			if !compareStrPtr(postScraped.CareerEnd, tt.want.CareerEnd) {
				t.Errorf("CareerEnd = %v, want %v", postScraped.CareerEnd, tt.want.CareerEnd)
			}
			if !compareStrPtr(postScraped.CareerLength, tt.want.CareerLength) {
				t.Errorf("CareerLength = %v, want %v", postScraped.CareerLength, tt.want.CareerLength)
			}
		})
	}
}
