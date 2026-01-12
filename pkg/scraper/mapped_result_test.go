package scraper

import (
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stretchr/testify/assert"
)

// Test string method
func TestMappedResultString(t *testing.T) {
	tests := []struct {
		name          string
		data          mappedResult
		key           string
		expectedValue string
		expectedOk    bool
	}{
		{
			name:          "valid string",
			data:          mappedResult{"name": "test"},
			key:           "name",
			expectedValue: "test",
			expectedOk:    true,
		},
		{
			name:          "missing key",
			data:          mappedResult{},
			key:           "missing",
			expectedValue: "",
			expectedOk:    false,
		},
		{
			name:          "wrong type still returns ok true but empty value",
			data:          mappedResult{"num": 123},
			key:           "num",
			expectedValue: "",
			expectedOk:    true, // logs error but returns ok=true
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			val, ok := test.data.string(test.key)
			assert.Equal(t, test.expectedValue, val)
			assert.Equal(t, test.expectedOk, ok)
		})
	}
}

// Test mustString method
func TestMappedResultMustString(t *testing.T) {
	tests := []struct {
		name          string
		data          mappedResult
		key           string
		expectedValue string
	}{
		{
			name:          "valid string",
			data:          mappedResult{"name": "test"},
			key:           "name",
			expectedValue: "test",
		},
		{
			name:          "missing key returns empty string",
			data:          mappedResult{},
			key:           "missing",
			expectedValue: "",
		},
		{
			name:          "wrong type returns empty string",
			data:          mappedResult{"num": 123},
			key:           "num",
			expectedValue: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			val := test.data.mustString(test.key)
			assert.Equal(t, test.expectedValue, val)
		})
	}
}

// Test stringPtr method
func TestMappedResultStringPtr(t *testing.T) {
	tests := []struct {
		name          string
		data          mappedResult
		key           string
		expectedValue *string
	}{
		{
			name:          "valid string",
			data:          mappedResult{"name": "test"},
			key:           "name",
			expectedValue: strPtr("test"),
		},
		{
			name:          "missing key returns nil",
			data:          mappedResult{},
			key:           "missing",
			expectedValue: nil,
		},
		{
			name:          "wrong type returns non-nil pointer to empty string",
			data:          mappedResult{"num": 123},
			key:           "num",
			expectedValue: strPtr(""), // string() returns empty string but ok=true
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			val := test.data.stringPtr(test.key)
			if test.expectedValue == nil {
				assert.Nil(t, val)
			} else {
				assert.NotNil(t, val)
				assert.Equal(t, *test.expectedValue, *val)
			}
		})
	}
}

// Test stringSlice method
func TestMappedResultStringSlice(t *testing.T) {
	tests := []struct {
		name          string
		data          mappedResult
		key           string
		expectedValue []string
	}{
		{
			name:          "valid slice",
			data:          mappedResult{"tags": []string{"a", "b", "c"}},
			key:           "tags",
			expectedValue: []string{"a", "b", "c"},
		},
		{
			name:          "missing key returns nil",
			data:          mappedResult{},
			key:           "missing",
			expectedValue: nil,
		},
		{
			name:          "single value converted to slice",
			data:          mappedResult{"tags": "not a slice"},
			key:           "tags",
			expectedValue: []string{"not a slice"},
		},
		{
			name:          "wrong type returns nil",
			data:          mappedResult{"tags": 123},
			key:           "tags",
			expectedValue: nil,
		},
		{
			name:          "empty slice",
			data:          mappedResult{"tags": []string{}},
			key:           "tags",
			expectedValue: []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			val := test.data.stringSlice(test.key)
			assert.Equal(t, test.expectedValue, val)
		})
	}
}

// Test IntPtr method
func TestMappedResultIntPtr(t *testing.T) {
	tests := []struct {
		name          string
		data          mappedResult
		key           string
		expectedValue *int
	}{
		{
			name:          "valid int",
			data:          mappedResult{"duration": 120},
			key:           "duration",
			expectedValue: intPtr(120),
		},
		{
			name:          "missing key returns nil",
			data:          mappedResult{},
			key:           "missing",
			expectedValue: nil,
		},
		{
			name:          "wrong type returns nil",
			data:          mappedResult{"duration": "120"},
			key:           "duration",
			expectedValue: nil,
		},
		{
			name:          "zero value",
			data:          mappedResult{"duration": 0},
			key:           "duration",
			expectedValue: intPtr(0),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			val := test.data.IntPtr(test.key)
			assert.Equal(t, test.expectedValue, val)
		})
	}
}

// Test setSingleValue method
func TestMappedResultsSetSingleValue(t *testing.T) {
	tests := []struct {
		name           string
		initialResults mappedResults
		index          int
		key            string
		value          string
		expectedLen    int
		shouldPanic    bool
	}{
		{
			name:           "append to empty",
			initialResults: mappedResults{},
			index:          0,
			key:            "name",
			value:          "test",
			expectedLen:    1,
			shouldPanic:    false,
		},
		{
			name:           "set in existing",
			initialResults: mappedResults{mappedResult{}},
			index:          0,
			key:            "name",
			value:          "test",
			expectedLen:    1,
			shouldPanic:    false,
		},
		{
			name:           "append to existing",
			initialResults: mappedResults{mappedResult{}},
			index:          1,
			key:            "name",
			value:          "test",
			expectedLen:    2,
			shouldPanic:    false,
		},
		{
			name:           "sparse index causes panic",
			initialResults: mappedResults{mappedResult{}},
			index:          5,
			key:            "name",
			value:          "test",
			expectedLen:    6,
			shouldPanic:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.shouldPanic {
				assert.Panics(t, func() {
					test.initialResults.setSingleValue(test.index, test.key, test.value)
				})
			} else {
				results := test.initialResults.setSingleValue(test.index, test.key, test.value)
				assert.Equal(t, test.expectedLen, len(results))
				assert.Equal(t, test.value, results[test.index][test.key])
			}
		})
	}
}

// Test setMultiValue method
func TestMappedResultsSetMultiValue(t *testing.T) {
	tests := []struct {
		name           string
		initialResults mappedResults
		index          int
		key            string
		value          []string
		expectedLen    int
	}{
		{
			name:           "append to empty",
			initialResults: mappedResults{},
			index:          0,
			key:            "tags",
			value:          []string{"a", "b"},
			expectedLen:    1,
		},
		{
			name:           "set in existing",
			initialResults: mappedResults{mappedResult{}},
			index:          0,
			key:            "tags",
			value:          []string{"a", "b"},
			expectedLen:    1,
		},
		{
			name:           "append to existing",
			initialResults: mappedResults{mappedResult{}},
			index:          1,
			key:            "tags",
			value:          []string{"x", "y"},
			expectedLen:    2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			results := test.initialResults.setMultiValue(test.index, test.key, test.value)
			assert.Equal(t, test.expectedLen, len(results))
			assert.Equal(t, test.value, results[test.index][test.key])
		})
	}
}

// Test scrapedTag method
func TestMappedResultScrapedTag(t *testing.T) {
	tests := []struct {
		name         string
		data         mappedResult
		expectedName string
	}{
		{
			name:         "valid tag",
			data:         mappedResult{"Name": "Action"},
			expectedName: "Action",
		},
		{
			name:         "missing name",
			data:         mappedResult{},
			expectedName: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tag := test.data.scrapedTag()
			assert.NotNil(t, tag)
			assert.Equal(t, test.expectedName, tag.Name)
		})
	}
}

// Test scrapedTags method
func TestMappedResultsScrapedTags(t *testing.T) {
	tests := []struct {
		name          string
		data          mappedResults
		expectedCount int
		expectedNames []string
	}{
		{
			name:          "empty results",
			data:          mappedResults{},
			expectedCount: 0,
		},
		{
			name: "single tag",
			data: mappedResults{
				mappedResult{"Name": "Action"},
			},
			expectedCount: 1,
			expectedNames: []string{"Action"},
		},
		{
			name: "multiple tags",
			data: mappedResults{
				mappedResult{"Name": "Action"},
				mappedResult{"Name": "Drama"},
				mappedResult{"Name": "Comedy"},
			},
			expectedCount: 3,
			expectedNames: []string{"Action", "Drama", "Comedy"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tags := test.data.scrapedTags()
			if test.expectedCount == 0 {
				assert.Nil(t, tags)
			} else {
				assert.NotNil(t, tags)
				assert.Equal(t, test.expectedCount, len(tags))
				for i, expectedName := range test.expectedNames {
					assert.Equal(t, expectedName, tags[i].Name)
				}
			}
		})
	}
}

// Test scrapedPerformer method
func TestMappedResultScrapedPerformer(t *testing.T) {
	tests := []struct {
		name     string
		data     mappedResult
		validate func(t *testing.T, p *models.ScrapedPerformer)
	}{
		{
			name: "full performer",
			data: mappedResult{
				"Name":           "Jane Doe",
				"Disambiguation": "Actress",
				"Gender":         "Female",
				"URL":            "https://example.com/jane",
				"URLs":           []string{"url1", "url2"},
				"Twitter":        "@jane",
				"Birthdate":      "1990-01-01",
				"Ethnicity":      "Caucasian",
				"Country":        "USA",
				"EyeColor":       "Blue",
				"Height":         "5'6\"",
				"Measurements":   "36-24-36",
				"FakeTits":       "No",
				"PenisLength":    "N/A",
				"Circumcised":    "N/A",
				"CareerLength":   "10 years",
				"Tattoos":        "Yes",
				"Piercings":      "Yes",
				"Aliases":        "Jane Smith",
				"Image":          "image.jpg",
				"Images":         []string{"img1", "img2"},
				"Details":        "Some details",
				"DeathDate":      "N/A",
				"HairColor":      "Blonde",
				"Weight":         "130 lbs",
			},
			validate: func(t *testing.T, p *models.ScrapedPerformer) {
				assert.NotNil(t, p)
				assert.Equal(t, "Jane Doe", *p.Name)
				assert.Equal(t, "Actress", *p.Disambiguation)
				assert.Equal(t, "Female", *p.Gender)
				assert.Equal(t, "https://example.com/jane", *p.URL)
				assert.Equal(t, []string{"url1", "url2"}, p.URLs)
				assert.Equal(t, "@jane", *p.Twitter)
				assert.Equal(t, "Blonde", *p.HairColor)
				assert.Equal(t, "130 lbs", *p.Weight)
			},
		},
		{
			name: "minimal performer",
			data: mappedResult{},
			validate: func(t *testing.T, p *models.ScrapedPerformer) {
				assert.NotNil(t, p)
				assert.Nil(t, p.Name)
				assert.Nil(t, p.Gender)
				assert.Empty(t, p.URLs)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			performer := test.data.scrapedPerformer()
			test.validate(t, performer)
		})
	}
}

// Test scrapedPerformers method
func TestMappedResultsScrapedPerformers(t *testing.T) {
	tests := []struct {
		name          string
		data          mappedResults
		expectedCount int
	}{
		{
			name:          "empty results",
			data:          mappedResults{},
			expectedCount: 0,
		},
		{
			name: "single performer",
			data: mappedResults{
				mappedResult{"Name": "Jane Doe"},
			},
			expectedCount: 1,
		},
		{
			name: "multiple performers",
			data: mappedResults{
				mappedResult{"Name": "Jane Doe"},
				mappedResult{"Name": "John Doe"},
				mappedResult{"Name": "Alice"},
			},
			expectedCount: 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			performers := test.data.scrapedPerformers()
			if test.expectedCount == 0 {
				assert.Nil(t, performers)
			} else {
				assert.NotNil(t, performers)
				assert.Equal(t, test.expectedCount, len(performers))
			}
		})
	}
}

// Test scrapedScene method
func TestMappedResultScrapedScene(t *testing.T) {
	tests := []struct {
		name     string
		data     mappedResult
		validate func(t *testing.T, s *models.ScrapedScene)
	}{
		{
			name: "full scene",
			data: mappedResult{
				"Title":    "Scene Title",
				"Code":     "CODE123",
				"Details":  "Scene details",
				"Director": "John Smith",
				"URL":      "https://example.com/scene",
				"URLs":     []string{"url1", "url2"},
				"Date":     "2020-01-01",
				"Image":    "scene.jpg",
				"Duration": 3600,
			},
			validate: func(t *testing.T, s *models.ScrapedScene) {
				assert.NotNil(t, s)
				assert.Equal(t, "Scene Title", *s.Title)
				assert.Equal(t, "CODE123", *s.Code)
				assert.Equal(t, "Scene details", *s.Details)
				assert.Equal(t, "John Smith", *s.Director)
				assert.Equal(t, "https://example.com/scene", *s.URL)
				assert.Equal(t, []string{"url1", "url2"}, s.URLs)
				assert.Equal(t, "2020-01-01", *s.Date)
				assert.Equal(t, "scene.jpg", *s.Image)
				assert.Equal(t, 3600, *s.Duration)
			},
		},
		{
			name: "minimal scene",
			data: mappedResult{},
			validate: func(t *testing.T, s *models.ScrapedScene) {
				assert.NotNil(t, s)
				assert.Nil(t, s.Title)
				assert.Nil(t, s.Duration)
				assert.Empty(t, s.URLs)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			scene := test.data.scrapedScene()
			test.validate(t, scene)
		})
	}
}

// Test scrapedImage method
func TestMappedResultScrapedImage(t *testing.T) {
	tests := []struct {
		name     string
		data     mappedResult
		validate func(t *testing.T, i *models.ScrapedImage)
	}{
		{
			name: "full image",
			data: mappedResult{
				"Title":        "Image Title",
				"Code":         "IMG123",
				"Details":      "Image details",
				"Photographer": "Jane Photographer",
				"URLs":         []string{"url1", "url2"},
				"Date":         "2020-06-15",
			},
			validate: func(t *testing.T, i *models.ScrapedImage) {
				assert.NotNil(t, i)
				assert.Equal(t, "Image Title", *i.Title)
				assert.Equal(t, "IMG123", *i.Code)
				assert.Equal(t, "Image details", *i.Details)
				assert.Equal(t, "Jane Photographer", *i.Photographer)
				assert.Equal(t, []string{"url1", "url2"}, i.URLs)
				assert.Equal(t, "2020-06-15", *i.Date)
			},
		},
		{
			name: "minimal image",
			data: mappedResult{},
			validate: func(t *testing.T, i *models.ScrapedImage) {
				assert.NotNil(t, i)
				assert.Nil(t, i.Title)
				assert.Empty(t, i.URLs)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			image := test.data.scrapedImage()
			test.validate(t, image)
		})
	}
}

// Test scrapedGallery method
func TestMappedResultScrapedGallery(t *testing.T) {
	tests := []struct {
		name     string
		data     mappedResult
		validate func(t *testing.T, g *models.ScrapedGallery)
	}{
		{
			name: "full gallery",
			data: mappedResult{
				"Title":        "Gallery Title",
				"Code":         "GAL123",
				"Details":      "Gallery details",
				"Photographer": "Jane Photographer",
				"URL":          "https://example.com/gallery",
				"URLs":         []string{"url1", "url2"},
				"Date":         "2020-07-20",
			},
			validate: func(t *testing.T, g *models.ScrapedGallery) {
				assert.NotNil(t, g)
				assert.Equal(t, "Gallery Title", *g.Title)
				assert.Equal(t, "GAL123", *g.Code)
				assert.Equal(t, "Gallery details", *g.Details)
				assert.Equal(t, "Jane Photographer", *g.Photographer)
				assert.Equal(t, "https://example.com/gallery", *g.URL)
				assert.Equal(t, []string{"url1", "url2"}, g.URLs)
				assert.Equal(t, "2020-07-20", *g.Date)
			},
		},
		{
			name: "minimal gallery",
			data: mappedResult{},
			validate: func(t *testing.T, g *models.ScrapedGallery) {
				assert.NotNil(t, g)
				assert.Nil(t, g.Title)
				assert.Empty(t, g.URLs)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gallery := test.data.scrapedGallery()
			test.validate(t, gallery)
		})
	}
}

// Test scrapedStudio method
func TestMappedResultScrapedStudio(t *testing.T) {
	tests := []struct {
		name     string
		data     mappedResult
		validate func(t *testing.T, st *models.ScrapedStudio)
	}{
		{
			name: "full studio",
			data: mappedResult{
				"Name":    "Studio Name",
				"URL":     "https://example.com/studio",
				"URLs":    []string{"url1", "url2"},
				"Image":   "studio.jpg",
				"Details": "Studio details",
				"Aliases": "Studio Alias",
			},
			validate: func(t *testing.T, st *models.ScrapedStudio) {
				assert.NotNil(t, st)
				assert.Equal(t, "Studio Name", st.Name)
				assert.Equal(t, "https://example.com/studio", *st.URL)
				assert.Equal(t, []string{"url1", "url2"}, st.URLs)
				assert.Equal(t, "studio.jpg", *st.Image)
				assert.Equal(t, "Studio details", *st.Details)
				assert.Equal(t, "Studio Alias", *st.Aliases)
			},
		},
		{
			name: "minimal studio",
			data: mappedResult{},
			validate: func(t *testing.T, st *models.ScrapedStudio) {
				assert.NotNil(t, st)
				assert.Equal(t, "", st.Name) // mustString returns empty string
				assert.Nil(t, st.URL)
				assert.Empty(t, st.URLs)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			studio := test.data.scrapedStudio()
			test.validate(t, studio)
		})
	}
}

// Test scrapedMovie method
func TestMappedResultScrapedMovie(t *testing.T) {
	tests := []struct {
		name     string
		data     mappedResult
		validate func(t *testing.T, m *models.ScrapedMovie)
	}{
		{
			name: "full movie",
			data: mappedResult{
				"Name":       "Movie Title",
				"Aliases":    "Movie Alias",
				"URLs":       []string{"url1", "url2"},
				"Duration":   "120 minutes",
				"Date":       "2020-05-10",
				"Director":   "John Director",
				"Synopsis":   "Movie synopsis",
				"FrontImage": "front.jpg",
				"BackImage":  "back.jpg",
			},
			validate: func(t *testing.T, m *models.ScrapedMovie) {
				assert.NotNil(t, m)
				assert.Equal(t, "Movie Title", *m.Name)
				assert.Equal(t, "Movie Alias", *m.Aliases)
				assert.Equal(t, []string{"url1", "url2"}, m.URLs)
				assert.Equal(t, "120 minutes", *m.Duration)
				assert.Equal(t, "2020-05-10", *m.Date)
				assert.Equal(t, "John Director", *m.Director)
				assert.Equal(t, "Movie synopsis", *m.Synopsis)
				assert.Equal(t, "front.jpg", *m.FrontImage)
				assert.Equal(t, "back.jpg", *m.BackImage)
			},
		},
		{
			name: "minimal movie",
			data: mappedResult{},
			validate: func(t *testing.T, m *models.ScrapedMovie) {
				assert.NotNil(t, m)
				assert.Nil(t, m.Name)
				assert.Empty(t, m.URLs)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			movie := test.data.scrapedMovie()
			test.validate(t, movie)
		})
	}
}

// Test scrapedMovies method
func TestMappedResultsScrapedMovies(t *testing.T) {
	tests := []struct {
		name          string
		data          mappedResults
		expectedCount int
	}{
		{
			name:          "empty results",
			data:          mappedResults{},
			expectedCount: 0,
		},
		{
			name: "single movie",
			data: mappedResults{
				mappedResult{"Name": "Movie 1"},
			},
			expectedCount: 1,
		},
		{
			name: "multiple movies",
			data: mappedResults{
				mappedResult{"Name": "Movie 1"},
				mappedResult{"Name": "Movie 2"},
				mappedResult{"Name": "Movie 3"},
			},
			expectedCount: 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			movies := test.data.scrapedMovies()
			if test.expectedCount == 0 {
				assert.Nil(t, movies)
			} else {
				assert.NotNil(t, movies)
				assert.Equal(t, test.expectedCount, len(movies))
			}
		})
	}
}

// Test scrapedGroup method
func TestMappedResultScrapedGroup(t *testing.T) {
	tests := []struct {
		name     string
		data     mappedResult
		validate func(t *testing.T, g *models.ScrapedGroup)
	}{
		{
			name: "full group",
			data: mappedResult{
				"Name":       "Group Title",
				"Aliases":    "Group Alias",
				"URL":        "https://example.com/group",
				"URLs":       []string{"url1", "url2"},
				"Duration":   "240 minutes",
				"Date":       "2020-08-15",
				"Director":   "Jane Director",
				"Synopsis":   "Group synopsis",
				"FrontImage": "front.jpg",
				"BackImage":  "back.jpg",
			},
			validate: func(t *testing.T, g *models.ScrapedGroup) {
				assert.NotNil(t, g)
				assert.Equal(t, "Group Title", *g.Name)
				assert.Equal(t, "Group Alias", *g.Aliases)
				assert.Equal(t, "https://example.com/group", *g.URL)
				assert.Equal(t, []string{"url1", "url2"}, g.URLs)
				assert.Equal(t, "240 minutes", *g.Duration)
				assert.Equal(t, "2020-08-15", *g.Date)
				assert.Equal(t, "Jane Director", *g.Director)
				assert.Equal(t, "Group synopsis", *g.Synopsis)
				assert.Equal(t, "front.jpg", *g.FrontImage)
				assert.Equal(t, "back.jpg", *g.BackImage)
			},
		},
		{
			name: "minimal group",
			data: mappedResult{},
			validate: func(t *testing.T, g *models.ScrapedGroup) {
				assert.NotNil(t, g)
				assert.Nil(t, g.Name)
				assert.Empty(t, g.URLs)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			group := test.data.scrapedGroup()
			test.validate(t, group)
		})
	}
}

// Test scrapedGroups method
func TestMappedResultsScrapedGroups(t *testing.T) {
	tests := []struct {
		name          string
		data          mappedResults
		expectedCount int
	}{
		{
			name:          "empty results",
			data:          mappedResults{},
			expectedCount: 0,
		},
		{
			name: "single group",
			data: mappedResults{
				mappedResult{"Name": "Group 1"},
			},
			expectedCount: 1,
		},
		{
			name: "multiple groups",
			data: mappedResults{
				mappedResult{"Name": "Group 1"},
				mappedResult{"Name": "Group 2"},
				mappedResult{"Name": "Group 3"},
			},
			expectedCount: 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			groups := test.data.scrapedGroups()
			if test.expectedCount == 0 {
				assert.Nil(t, groups)
			} else {
				assert.NotNil(t, groups)
				assert.Equal(t, test.expectedCount, len(groups))
			}
		})
	}
}

// Helper functions
func strPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
