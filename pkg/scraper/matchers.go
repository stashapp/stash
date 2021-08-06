package scraper

import (
	"strconv"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/tag"
)

// MatchScrapedScenePerformer matches the provided performer with the
// performers in the database and sets the ID field if one is found.
func MatchScrapedScenePerformer(qb models.PerformerReader, p *models.ScrapedPerformer) error {
	performers, err := qb.FindByNames([]string{p.Name}, true)

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

// MatchScrapedSceneStudio matches the provided studio with the studios
// in the database and sets the ID field if one is found.
func MatchScrapedSceneStudio(qb models.StudioReader, s *models.ScrapedStudio) error {
	studio, err := qb.FindByName(s.Name, true)

	if err != nil {
		return err
	}

	if studio == nil {
		// ignore - cannot match
		return nil
	}

	id := strconv.Itoa(studio.ID)
	s.StoredID = &id
	return nil
}

// MatchScrapedSceneMovie matches the provided movie with the movies
// in the database and sets the ID field if one is found.
func MatchScrapedSceneMovie(qb models.MovieReader, m *models.ScrapedMovie) error {
	movies, err := qb.FindByNames([]string{m.Name}, true)

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

// MatchScrapedSceneTag matches the provided tag with the tags
// in the database and sets the ID field if one is found.
func MatchScrapedSceneTag(qb models.TagReader, s *models.ScrapedTag) error {
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
