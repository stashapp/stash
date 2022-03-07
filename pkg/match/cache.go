package match

import "github.com/stashapp/stash/pkg/models"

const singleFirstCharacterRegex = `^[\p{L}][.\-_ ]`

// Cache is used to cache queries that should not change across an autotag process.
type Cache struct {
	singleCharPerformers []*models.Performer
	singleCharStudios    []*models.Studio
	singleCharTags       []*models.Tag
}

// getSingleLetterPerformers returns all performers with names that start with single character words.
// The autotag query splits the words into two-character words to query
// against. This means that performers with single-letter words in their names could potentially
// be missed.
// This query is expensive, so it's queried once and cached, if the cache if provided.
func getSingleLetterPerformers(c *Cache, reader models.PerformerReader) ([]*models.Performer, error) {
	if c == nil {
		c = &Cache{}
	}

	if c.singleCharPerformers == nil {
		pp := -1
		performers, _, err := reader.Query(&models.PerformerFilterType{
			Name: &models.StringCriterionInput{
				Value:    singleFirstCharacterRegex,
				Modifier: models.CriterionModifierMatchesRegex,
			},
		}, &models.FindFilterType{
			PerPage: &pp,
		})

		if err != nil {
			return nil, err
		}

		if len(performers) == 0 {
			// make singleWordPerformers not nil
			c.singleCharPerformers = make([]*models.Performer, 0)
		} else {
			c.singleCharPerformers = performers
		}
	}

	return c.singleCharPerformers, nil
}

// getSingleLetterStudios returns all performers with names that start with single character words.
// See getSingleLetterPerformers for details.
func getSingleLetterStudios(c *Cache, reader models.StudioReader) ([]*models.Studio, error) {
	if c == nil {
		c = &Cache{}
	}

	if c.singleCharPerformers == nil {
		pp := -1
		studios, _, err := reader.Query(&models.StudioFilterType{
			Name: &models.StringCriterionInput{
				Value:    singleFirstCharacterRegex,
				Modifier: models.CriterionModifierMatchesRegex,
			},
		}, &models.FindFilterType{
			PerPage: &pp,
		})

		if err != nil {
			return nil, err
		}

		if len(studios) == 0 {
			// make singleWordStudios not nil
			c.singleCharStudios = make([]*models.Studio, 0)
		} else {
			c.singleCharStudios = studios
		}
	}

	return c.singleCharStudios, nil
}

// getSingleLetterTags returns all performers with names that start with single character words.
// See getSingleLetterPerformers for details.
func getSingleLetterTags(c *Cache, reader models.TagReader) ([]*models.Tag, error) {
	if c == nil {
		c = &Cache{}
	}

	if c.singleCharPerformers == nil {
		pp := -1
		tags, _, err := reader.Query(&models.TagFilterType{
			Name: &models.StringCriterionInput{
				Value:    singleFirstCharacterRegex,
				Modifier: models.CriterionModifierMatchesRegex,
			},
			Or: &models.TagFilterType{
				Aliases: &models.StringCriterionInput{
					Value:    singleFirstCharacterRegex,
					Modifier: models.CriterionModifierMatchesRegex,
				},
			},
		}, &models.FindFilterType{
			PerPage: &pp,
		})

		if err != nil {
			return nil, err
		}

		if len(tags) == 0 {
			// make singleWordTags not nil
			c.singleCharTags = make([]*models.Tag, 0)
		} else {
			c.singleCharTags = tags
		}
	}

	return c.singleCharTags, nil
}
