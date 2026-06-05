package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

const (
	storyTable              = "stories"
	storyIDColumn           = "story_id"
	storiesTagsTable        = "stories_tags"
	storiesTagIDColumn      = "tag_id"
	performersStoriesTable  = "performers_stories"
	performersStoryIDColumn = "performer_id"
	storyURLsTable          = "story_urls"
	storyURLColumn          = "url"
	storyFrontImageBlobCol  = "front_image_blob"
	storyBackImageBlobCol   = "back_image_blob"
)

type storyRow struct {
	ID              int          `db:"id" goqu:"skipinsert"`
	Title           zero.String  `db:"title"`
	Author          zero.String  `db:"author"`
	URL             zero.String  `db:"url"`
	Date            NullDate     `db:"date"`
	Language        zero.String  `db:"language"`
	TagLine         zero.String  `db:"tag_line"`
	Details         zero.String  `db:"details"`
	StudioID        null.Int     `db:"studio_id"`
	Rating          null.Int     `db:"rating"`
	CreatedAt       Timestamp    `db:"created_at"`
	UpdatedAt       Timestamp    `db:"updated_at"`
	FrontImageBlob  zero.String  `db:"front_image_blob"`
	BackImageBlob   zero.String  `db:"back_image_blob"`
}

type StoryStore struct {
	blobJoinQueryBuilder *blobJoinQueryBuilder
	blobStore           *BlobStore
	tableMgr            *table
}

func NewStoryStore(blobStore *BlobStore) *StoryStore {
	return &StoryStore{
		blobStore: blobStore,
		tableMgr:  storyTableMgr,
	}
}

func (s *StoryStore) table() exp.IdentifierExpression {
	return s.tableMgr.table
}

func (s *StoryStore) Create(ctx context.Context, input *models.CreateStoryInput) error {
	story := input.Story
	r := &storyRow{}
	r.ID = story.ID
	r.Title = zero.StringFrom(story.Title)
	r.Author = zero.StringFrom(story.Author)
	r.URL = zero.StringFrom(story.URL)
	if story.Date != nil {
		r.Date = NullDateFromDatePtr(story.Date)
	}
	r.Language = zero.StringFrom(story.Language)
	r.TagLine = zero.StringFrom(story.TagLine)
	r.Details = zero.StringFrom(story.Details)
	r.StudioID = intFromPtr(story.StudioID)
	r.Rating = intFromPtr(story.Rating)
	r.CreatedAt = Timestamp{Timestamp: story.CreatedAt}
	r.UpdatedAt = Timestamp{Timestamp: story.UpdatedAt}

	q := dialect.Insert(s.table()).Rows(r).Returning(s.table().Col("id"))
	if err := querySimple(ctx, q, &story.ID); err != nil {
		return fmt.Errorf("creating story: %w", err)
	}

	// Update join tables
	if err := storiesTagsTableMgr.replace(ctx, story.ID, story.TagIDs.List()); err != nil {
		return fmt.Errorf("updating story tags: %w", err)
	}
	if err := performersStoriesTableMgr.replace(ctx, story.ID, story.PerformerIDs.List()); err != nil {
		return fmt.Errorf("updating story performers: %w", err)
	}
	if err := storyURLsTableMgr.replace(ctx, story.ID, story.URLs.List()); err != nil {
		return fmt.Errorf("updating story urls: %w", err)
	}

	return nil
}

func (s *StoryStore) Update(ctx context.Context, input *models.UpdateStoryInput) error {
	q := dialect.Update(s.table()).Set(
		goqu.Record{
			"title":     input.Title,
			"author":    input.Author,
			"url":       input.URL,
			"date":      input.Date,
			"language":  input.Language,
			"tag_line":  input.TagLine,
			"details":   input.Details,
			"studio_id": input.StudioID,
			"rating":    input.Rating,
		},
	).Where(s.table().Col("id").Eq(input.ID))

	if _, err := dbWrapper.Exec(ctx, q); err != nil {
		return fmt.Errorf("updating story: %w", err)
	}

	return nil
}

func (s *StoryStore) UpdatePartial(ctx context.Context, id int, partial models.StoryPartial) (*models.Story, error) {
	existing, err := s.Find(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, fmt.Errorf("story with id %d not found", id)
	}

	record := goqu.Record{}
	record["updated_at"] = Timestamp{Timestamp: time.Now()}
	if partial.Title.Set { record["title"] = zero.StringFrom(partial.Title.Value) }
	if partial.Author.Set { record["author"] = zero.StringFrom(partial.Author.Value) }
	if partial.URL.Set { record["url"] = zero.StringFrom(partial.URL.Value) }
	if partial.Date.Set { record["date"] = NullDateFromDatePtr(partial.Date.Value) }
	if partial.Language.Set { record["language"] = zero.StringFrom(partial.Language.Value) }
	if partial.TagLine.Set { record["tag_line"] = zero.StringFrom(partial.TagLine.Value) }
	if partial.Details.Set { record["details"] = zero.StringFrom(partial.Details.Value) }
	if partial.StudioID.Set { record["studio_id"] = intFromPtr(partial.StudioID.Value) }
	if partial.Rating.Set { record["rating"] = intFromPtr(partial.Rating.Value) }
	if partial.Organized.Set { record["organized"] = partial.Organized.Value }

	q := dialect.Update(s.table()).Set(record).Where(s.table().Col("id").Eq(id))
	if _, err := dbWrapper.Exec(ctx, q); err != nil {
		return nil, fmt.Errorf("updating story: %w", err)
	}

	if partial.TagIDs != nil {
		if err := storiesTagsTableMgr.replace(ctx, id, partial.TagIDs.Apply(existing.TagIDs.List())); err != nil {
			return nil, fmt.Errorf("updating story tags: %w", err)
		}
	}
	if partial.PerformerIDs != nil {
		if err := performersStoriesTableMgr.replace(ctx, id, partial.PerformerIDs.Apply(existing.PerformerIDs.List())); err != nil {
			return nil, fmt.Errorf("updating story performers: %w", err)
		}
	}
	if partial.URLs != nil {
		if err := storyURLsTableMgr.replace(ctx, id, partial.URLs.Apply(existing.URLs.List())); err != nil {
			return nil, fmt.Errorf("updating story urls: %w", err)
		}
	}

	return s.Find(ctx, id)
}

func (s *StoryStore) Destroy(ctx context.Context, id int) error {
	_, err := dbWrapper.Exec(ctx, "DELETE FROM stories WHERE id = ?", id)
	return err
}

func (s *StoryStore) Find(ctx context.Context, id int) (*models.Story, error) {
	q := dialect.From(s.table()).Where(s.table().Col("id").Eq(id))
	ret, err := s.find(ctx, q)
	if err != nil {
		return nil, err
	}
	if len(ret) == 0 {
		return nil, nil
	}
	return ret[0], nil
}

func (s *StoryStore) FindMany(ctx context.Context, ids []int) ([]*models.Story, error) {
	q := dialect.From(s.table()).Where(s.table().Col("id").In(ids))
	return s.find(ctx, q)
}

func (s *StoryStore) FindByPerformerID(ctx context.Context, performerID int) ([]*models.Story, error) {
	q := dialect.From(s.table()).
		Join(performersStoriesTableMgr.table, goqu.On(s.table().Col("id").Eq(performersStoriesTableMgr.table.Col(storyIDColumn)))).
		Where(performersStoriesTableMgr.table.Col(performersStoryIDColumn).Eq(performerID))
	return s.find(ctx, q)
}

func (s *StoryStore) FindByStudioID(ctx context.Context, studioID int) ([]*models.Story, error) {
	q := dialect.From(s.table()).Where(s.table().Col("studio_id").Eq(studioID))
	return s.find(ctx, q)
}

func (s *StoryStore) All(ctx context.Context) ([]*models.Story, error) {
	q := dialect.From(s.table())
	return s.find(ctx, q)
}

func (s *StoryStore) Query(ctx context.Context, options models.StoryFilterType, findFilter *models.FindFilterType) ([]*models.Story, int, error) {
	query := s.table()
	if q := query; true {
		_ = q
	}
	ret, count, err := s.executeQuery(ctx, options, findFilter)
	if err != nil {
		return nil, 0, err
	}
	return ret, count, nil
}

func (s *StoryStore) find(ctx context.Context, q *goqu.SelectDataset) ([]*models.Story, error) {
	var rows []storyRow
	if err := querySimple(ctx, q, &rows); err != nil {
		return nil, err
	}

	var ret []*models.Story
	for _, r := range rows {
		story := r.toStory()
		ret = append(ret, story)
	}

	return ret, nil
}

func (r *storyRow) toStory() *models.Story {
	return &models.Story{
		ID:        r.ID,
		Title:     r.Title.String,
		Author:    r.Author.String,
		URL:       r.URL.String,
		Language:  r.Language.String,
		TagLine:   r.TagLine.String,
		Details:   r.Details.String,
		StudioID:  nullIntPtr(r.StudioID),
		Rating:    nullIntPtr(r.Rating),
		CreatedAt: r.CreatedAt.Timestamp,
		UpdatedAt: r.UpdatedAt.Timestamp,
		Date:      r.Date.ToDatePtr(),
	}
}

func (s *StoryStore) GetFrontImage(ctx context.Context, storyID int) ([]byte, error) {
	return s.blobJoinQueryBuilder.GetImage(ctx, storyID, storyFrontImageBlobCol)
}

func (s *StoryStore) GetBackImage(ctx context.Context, storyID int) ([]byte, error) {
	return s.blobJoinQueryBuilder.GetImage(ctx, storyID, storyBackImageBlobCol)
}

func (s *StoryStore) HasFrontImage(ctx context.Context, storyID int) (bool, error) {
	return s.blobJoinQueryBuilder.HasImage(ctx, storyID, storyFrontImageBlobCol)
}

func (s *StoryStore) HasBackImage(ctx context.Context, storyID int) (bool, error) {
	return s.blobJoinQueryBuilder.HasImage(ctx, storyID, storyBackImageBlobCol)
}

func (s *StoryStore) UpdateFrontImage(ctx context.Context, storyID int, image []byte) error {
	return s.blobJoinQueryBuilder.UpdateImage(ctx, storyID, storyFrontImageBlobCol, image)
}

func (s *StoryStore) UpdateBackImage(ctx context.Context, storyID int, image []byte) error {
	return s.blobJoinQueryBuilder.UpdateImage(ctx, storyID, storyBackImageBlobCol, image)
}

func (s *StoryStore) GetTagIDs(ctx context.Context, storyID int) ([]int, error) {
	return storiesTagsTableMgr.get(ctx, storyID)
}

func (s *StoryStore) GetPerformerIDs(ctx context.Context, storyID int) ([]int, error) {
	return performersStoriesTableMgr.get(ctx, storyID)
}

func (s *StoryStore) GetURLs(ctx context.Context, storyID int) ([]string, error) {
	return storyURLsTableMgr.get(ctx, storyID)
}
