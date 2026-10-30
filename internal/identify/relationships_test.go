package identify

import (
	"strconv"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stretchr/testify/assert"
)

func Test_relationshipResolver_performers_skipSingleName(t *testing.T) {
	const performerID = 2
	createMissing := true
	singleName := "Single"
	storedID := strconv.Itoa(performerID)

	resolver := relationshipResolver{
		fieldOptions: map[string]*FieldOptions{
			"performers": {
				Strategy:      FieldStrategyMerge,
				CreateMissing: &createMissing,
			},
		},
		skipSingleNamePerformers: true,
	}

	tests := []struct {
		name        string
		existingIDs []int
		scraped     []*models.ScrapedPerformer
		want        []int
	}{
		{
			name:        "only skipped performer",
			existingIDs: nil,
			scraped:     []*models.ScrapedPerformer{{Name: &singleName}},
			want:        nil,
		},
		{
			name:        "skipped performer with another result",
			existingIDs: []int{},
			scraped: []*models.ScrapedPerformer{
				{Name: &singleName},
				{Name: &singleName, StoredID: &storedID},
			},
			want: []int{performerID},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolver.performers(testCtx, tt.existingIDs, tt.scraped, nil)

			assert.Equal(t, tt.want, got)
			assert.ErrorIs(t, err, ErrSkipSingleNamePerformer)
		})
	}
}
