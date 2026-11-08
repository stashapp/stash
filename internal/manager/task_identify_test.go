package manager

import (
	"context"
	"strconv"
	"sync"
	"testing"

	"github.com/stashapp/stash/internal/identify"
	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stashapp/stash/pkg/scraper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockSceneScraper struct {
	results map[int][]*models.ScrapedScene
}

func (s mockSceneScraper) ScrapeScenes(ctx context.Context, sceneID int) ([]*models.ScrapedScene, error) {
	return s.results[sceneID], nil
}

type mockGalleryScraper struct {
	results map[int][]*models.ScrapedGallery
}

func (s mockGalleryScraper) ScrapeGalleries(ctx context.Context, galleryID int) ([]*models.ScrapedGallery, error) {
	return s.results[galleryID], nil
}

type mockHookExecutor struct{}

func (mockHookExecutor) ExecuteSceneUpdatePostHooks(ctx context.Context, input models.SceneUpdateInput, inputFields []string) {
}

func (mockHookExecutor) ExecuteGalleryUpdatePostHooks(ctx context.Context, input models.GalleryUpdateInput, inputFields []string) {
}

func TestIdentifyJob_Execute(t *testing.T) {
	testScraperID := "test-scraper"
	source := &identify.Source{
		Source: &scraper.Source{ScraperID: &testScraperID},
	}
	sceneSource := identify.ScraperSource{
		Name:    "test-scene-scraper",
		Scraper: mockSceneScraper{},
	}
	gallerySource := identify.GalleryScraperSource{
		Name:    "test-gallery-scraper",
		Scraper: mockGalleryScraper{},
	}

	tests := []struct {
		name            string
		sources         []*identify.Source
		sceneIDs        []string
		galleryIDs      []string
		sceneFIDs       []int
		galleryFIDs     []int
		wantErr         string
		cancelAfter     int
		noGallerySource bool
	}{
		{
			name: "no sources",
		},
		{
			name:      "scene IDs only",
			sources:   []*identify.Source{source},
			sceneIDs:  []string{"1", "2"},
			sceneFIDs: []int{1, 2},
		},
		{
			name:        "gallery IDs only",
			sources:     []*identify.Source{source},
			galleryIDs:  []string{"1", "2"},
			galleryFIDs: []int{1, 2},
		},
		{
			name:        "both IDs",
			sources:     []*identify.Source{source},
			sceneIDs:    []string{"1"},
			galleryIDs:  []string{"2"},
			sceneFIDs:   []int{1},
			galleryFIDs: []int{2},
		},
		{
			name:     "scene ID not found",
			sources:  []*identify.Source{source},
			sceneIDs: []string{"999"},
			wantErr:  "scene with id 999 not found",
		},
		{
			name:       "gallery ID not found",
			sources:    []*identify.Source{source},
			galleryIDs: []string{"999"},
			wantErr:    "gallery with id 999 not found",
		},
		{
			name:     "invalid scene ID",
			sources:  []*identify.Source{source},
			sceneIDs: []string{"bad"},
			wantErr:  "invalid scene IDs",
		},
		{
			name:       "invalid gallery ID",
			sources:    []*identify.Source{source},
			galleryIDs: []string{"bad"},
			wantErr:    "invalid gallery IDs",
		},
		{
			name:            "no IDs identify-all scenes only",
			sources:         []*identify.Source{source},
			sceneFIDs:       []int{10},
			noGallerySource: true,
		},
		{
			name:        "no IDs identify-all galleries and scenes",
			sources:     []*identify.Source{source},
			sceneFIDs:   []int{10},
			galleryFIDs: []int{20},
		},
		{
			name:        "cancellation during gallery IDs",
			sources:     []*identify.Source{source},
			galleryIDs:  []string{"1", "2", "3"},
			galleryFIDs: []int{1},
			cancelAfter: 1,
		},
		{
			name:        "cancellation during scene IDs",
			sources:     []*identify.Source{source},
			sceneIDs:    []string{"1", "2", "3"},
			sceneFIDs:   []int{1},
			cancelAfter: 1,
		},
		{
			name:       "both IDs with stash-box-only gallery source",
			sources:    []*identify.Source{source},
			sceneIDs:   []string{"1"},
			galleryIDs: []string{"2"},
			sceneFIDs:  []int{1},
			wantErr:    "no gallery sources available",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := mocks.NewDatabase()

			ctx := context.Background()
			var cancel context.CancelFunc
			if tt.cancelAfter > 0 {
				ctx, cancel = context.WithCancel(ctx)
			}

			for _, idStr := range tt.sceneIDs {
				id, _ := strconv.Atoi(idStr)
				if id == 999 {
					db.Scene.On("Find", mock.Anything, 999).Return(nil, nil)
				} else {
					db.Scene.On("Find", mock.Anything, id).Return(&models.Scene{ID: id, Path: "scene-" + idStr}, nil)
				}
			}

			for _, idStr := range tt.galleryIDs {
				id, _ := strconv.Atoi(idStr)
				if id == 999 {
					db.Gallery.On("Find", mock.Anything, 999).Return(nil, nil)
				} else {
					db.Gallery.On("Find", mock.Anything, id).Return(&models.Gallery{ID: id, Title: "gallery-" + idStr}, nil)
				}
			}

			hasIdentifyAll := tt.wantErr == "" && len(tt.sceneIDs) == 0 && len(tt.galleryIDs) == 0 && len(tt.sources) > 0
			if hasIdentifyAll && len(tt.sceneFIDs) > 0 {
				scenes := make([]*models.Scene, len(tt.sceneFIDs))
				for i, id := range tt.sceneFIDs {
					scenes[i] = &models.Scene{ID: id, Path: "scene-" + strconv.Itoa(id)}
				}
				db.Scene.On("Query", mock.Anything, mock.Anything).Return(mocks.SceneQueryResult(scenes, len(scenes)), nil)
			}
			if hasIdentifyAll && !tt.noGallerySource && len(tt.galleryFIDs) > 0 {
				galleries := make([]*models.Gallery, len(tt.galleryFIDs))
				for i, id := range tt.galleryFIDs {
					galleries[i] = &models.Gallery{ID: id, Title: "gallery-" + strconv.Itoa(id)}
				}
				db.Gallery.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(galleries, len(galleries), nil)
			}

			callCount := 0
			var mu sync.Mutex
			var identifiedScenes []int
			var identifiedGalleries []int

			identifyJob := &IdentifyJob{
				repository:              db.Repository(),
				postHookExecutor:        mockHookExecutor{},
				galleryPostHookExecutor: mockHookExecutor{},
				input: identify.Options{
					Sources:    tt.sources,
					SceneIDs:   tt.sceneIDs,
					GalleryIDs: tt.galleryIDs,
				},
			}
			identifyJob.sourcesFn = func() ([]identify.ScraperSource, error) {
				return []identify.ScraperSource{sceneSource}, nil
			}
			identifyJob.gallerySourcesFn = func() ([]identify.GalleryScraperSource, error) {
				if tt.noGallerySource || tt.wantErr == "no gallery sources available" {
					return nil, nil
				}
				return []identify.GalleryScraperSource{gallerySource}, nil
			}
			identifyJob.identifySceneFn = func(ctx context.Context, s *models.Scene, sources []identify.ScraperSource) {
				mu.Lock()
				defer mu.Unlock()
				callCount++
				if cancel != nil && callCount >= tt.cancelAfter {
					cancel()
				}
				identifiedScenes = append(identifiedScenes, s.ID)
			}
			identifyJob.identifyGalleryFn = func(ctx context.Context, g *models.Gallery, sources []identify.GalleryScraperSource) {
				mu.Lock()
				defer mu.Unlock()
				callCount++
				if cancel != nil && callCount >= tt.cancelAfter {
					cancel()
				}
				identifiedGalleries = append(identifiedGalleries, g.ID)
			}

			progress := &job.Progress{}
			err := identifyJob.Execute(ctx, progress)

			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			if tt.cancelAfter > 0 {
				assert.LessOrEqual(t, len(identifiedScenes)+len(identifiedGalleries), tt.cancelAfter)
			} else if tt.wantErr == "" {

				assert.ElementsMatch(t, tt.sceneFIDs, identifiedScenes)
				assert.ElementsMatch(t, tt.galleryFIDs, identifiedGalleries)
			}
		})
	}
}

func Test_supportsGalleryFragment(t *testing.T) {
	tests := []struct {
		name    string
		scraper *scraper.Scraper
		want    bool
	}{
		{
			name:    "nil scraper",
			scraper: nil,
			want:    false,
		},
		{
			name: "gallery spec is nil",
			scraper: &scraper.Scraper{
				ID:   "scene-only",
				Name: "Scene Only",
				Scene: &scraper.ScraperSpec{
					SupportedScrapes: []scraper.ScrapeType{scraper.ScrapeTypeFragment},
				},
			},
			want: false,
		},
		{
			name: "gallery spec does not support fragment scraping",
			scraper: &scraper.Scraper{
				ID:   "gallery-url-only",
				Name: "Gallery URL Only",
				Gallery: &scraper.ScraperSpec{
					SupportedScrapes: []scraper.ScrapeType{scraper.ScrapeTypeURL},
				},
			},
			want: false,
		},
		{
			name: "gallery spec supports fragment scraping",
			scraper: &scraper.Scraper{
				ID:   "gallery-fragment",
				Name: "Gallery Fragment",
				Gallery: &scraper.ScraperSpec{
					SupportedScrapes: []scraper.ScrapeType{scraper.ScrapeTypeFragment},
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, supportsGalleryFragment(tt.scraper))
		})
	}
}
