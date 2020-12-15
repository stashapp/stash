package scraper

import (
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

// MatchScrapedScenePerformer matches the provided performer with the
// performers in the database and sets the ID field if one is found.
func MatchScrapedScenePerformer(qb models.PerformerReader, p *models.ScrapedScenePerformer) error {
	performers, err := qb.FindByNames([]string{p.Name}, true)

	if err != nil {
		return err
	}

	if len(performers) != 1 {
		// ignore - cannot match
		return nil
	}

	id := strconv.Itoa(performers[0].ID)
	p.ID = &id
	return nil
}

// MatchScrapedSceneStudio matches the provided studio with the studios
// in the database and sets the ID field if one is found.
func MatchScrapedSceneStudio(qb models.StudioReader, s *models.ScrapedSceneStudio) error {
	studio, err := qb.FindByName(s.Name, true)

	if err != nil {
		return err
	}

	if studio == nil {
		// ignore - cannot match
		return nil
	}

	id := strconv.Itoa(studio.ID)
	s.ID = &id
	return nil
}

// MatchScrapedSceneMovie matches the provided movie with the movies
// in the database and sets the ID field if one is found.
func MatchScrapedSceneMovie(qb models.MovieReader, m *models.ScrapedSceneMovie) error {
	movies, err := qb.FindByNames([]string{m.Name}, true)

	if err != nil {
		return err
	}

	if len(movies) != 1 {
		// ignore - cannot match
		return nil
	}

	id := strconv.Itoa(movies[0].ID)
	m.ID = &id
	return nil
}

// MatchScrapedSceneTag matches the provided tag with the tags
// in the database and sets the ID field if one is found.
func MatchScrapedSceneTag(qb models.TagReader, s *models.ScrapedSceneTag) error {
	tag, err := qb.FindByName(s.Name, true)

	if err != nil {
		return err
	}

	if tag == nil {
		// ignore - cannot match
		return nil
	}

	id := strconv.Itoa(tag.ID)
	s.ID = &id
	return nil
}
