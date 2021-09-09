package scraper

import (
	"strconv"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/studio"
	"github.com/stashapp/stash/pkg/tag"
)

// MatchScrapedPerformer matches the provided performer with the
// performers in the database and sets the ID field if one is found.
func MatchScrapedPerformer(qb models.PerformerReader, p *models.ScrapedPerformer) error {
	if p.Name == nil {
		return nil
	}

	performers, err := qb.FindByNames([]string{*p.Name}, true)

	if err != nil {
		return err
	}

	if len(performers) != 1 {
		// ignore - cannot match
		return nil
	}

	id := strconv.Itoa(performers[0].ID)
	p.StoredID = &id
	return nil
}

// MatchScrapedStudio matches the provided studio with the studios
// in the database and sets the ID field if one is found.
func MatchScrapedStudio(qb models.StudioReader, s *models.ScrapedStudio) error {
	st, err := studio.ByName(qb, s.Name)

	if err != nil {
		return err
	}

	if st == nil {
		// try matching by alias
		st, err = studio.ByAlias(qb, s.Name)
		if err != nil {
			return err
		}
	}

	if st == nil {
		// ignore - cannot match
		return nil
	}

	id := strconv.Itoa(st.ID)
	s.StoredID = &id
	return nil
}

// MatchScrapedMovie matches the provided movie with the movies
// in the database and sets the ID field if one is found.
func MatchScrapedMovie(qb models.MovieReader, m *models.ScrapedMovie) error {
	if m.Name == nil {
		return nil
	}

	movies, err := qb.FindByNames([]string{*m.Name}, true)

	if err != nil {
		return err
	}

	if len(movies) != 1 {
		// ignore - cannot match
		return nil
	}

	id := strconv.Itoa(movies[0].ID)
	m.StoredID = &id
	return nil
}

// MatchScrapedTag matches the provided tag with the tags
// in the database and sets the ID field if one is found.
func MatchScrapedTag(qb models.TagReader, s *models.ScrapedTag) error {
	t, err := tag.ByName(qb, s.Name)

	if err != nil {
		return err
	}

	if t == nil {
		// try matching by alias
		t, err = tag.ByAlias(qb, s.Name)
		if err != nil {
			return err
		}
	}

	if t == nil {
		// ignore - cannot match
		return nil
	}

	id := strconv.Itoa(t.ID)
	s.StoredID = &id
	return nil
}
