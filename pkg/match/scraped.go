package match

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/studio"
	"github.com/stashapp/stash/pkg/tag"
)

// ScrapedPerformer matches the provided performer with the
// performers in the database and sets the ID field if one is found.
func ScrapedPerformer(ctx context.Context, qb models.PerformerReader, p *models.ScrapedPerformer, stashBoxEndpoint *string) error {
	if p.StoredID != nil || p.Name == nil {
		return nil
	}

	// Check if a performer with the StashID already exists
	if stashBoxEndpoint != nil && p.RemoteSiteID != nil {
		performers, err := qb.FindByStashID(ctx, models.StashID{
			StashID:  *p.RemoteSiteID,
			Endpoint: *stashBoxEndpoint,
		})
		if err != nil {
			return err
		}
		if len(performers) > 0 {
			id := strconv.Itoa(performers[0].ID)
			p.StoredID = &id
			return nil
		}
	}

	performers, err := qb.FindByNames(ctx, []string{*p.Name}, true)

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

// ScrapedStudio matches the provided studio with the studios
// in the database and sets the ID field if one is found.
func ScrapedStudio(ctx context.Context, qb models.StudioReader, s *models.ScrapedStudio, stashBoxEndpoint *string) error {
	if s.StoredID != nil {
		return nil
	}

	// Check if a studio with the StashID already exists
	if stashBoxEndpoint != nil && s.RemoteSiteID != nil {
		studios, err := qb.FindByStashID(ctx, models.StashID{
			StashID:  *s.RemoteSiteID,
			Endpoint: *stashBoxEndpoint,
		})
		if err != nil {
			return err
		}
		if len(studios) > 0 {
			id := strconv.Itoa(studios[0].ID)
			s.StoredID = &id
			return nil
		}
	}

	st, err := studio.ByName(ctx, qb, s.Name)

	if err != nil {
		return err
	}

	if st == nil {
		// try matching by alias
		st, err = studio.ByAlias(ctx, qb, s.Name)
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

// ScrapedMovie matches the provided movie with the movies
// in the database and sets the ID field if one is found.
func ScrapedMovie(ctx context.Context, qb models.MovieReader, m *models.ScrapedMovie) error {
	if m.StoredID != nil || m.Name == nil {
		return nil
	}

	movies, err := qb.FindByNames(ctx, []string{*m.Name}, true)

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

// ScrapedTag matches the provided tag with the tags
// in the database and sets the ID field if one is found.
func ScrapedTag(ctx context.Context, qb models.TagReader, s *models.ScrapedTag) error {
	if s.StoredID != nil {
		return nil
	}

	t, err := tag.ByName(ctx, qb, s.Name)

	if err != nil {
		return err
	}

	if t == nil {
		// try matching by alias
		t, err = tag.ByAlias(ctx, qb, s.Name)
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
