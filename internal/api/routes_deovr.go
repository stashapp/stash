package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

// DeoVRFile represents a single file in the DeoVR API format.
type DeoVRFile struct {
	FileTitle       string  `json:"file_title"`
	FileURL         string  `json:"file_url"`
	FileSize        int64   `json:"file_size"`
	IsAvailable     bool    `json:"is_available"`
	VideoProjection string  `json:"video_projection"`
	VideoDuration   float64 `json:"video_duration"`
	VideoWidth      int     `json:"video_width"`
	VideoHeight     int     `json:"video_height"`
	ThumbnailURL    string  `json:"thumbnail_url"`
}

// DeoVRFilesResponse is the top-level response for file listing.
type DeoVRFilesResponse struct {
	Files []DeoVRFile `json:"files"`
}

type deovrRoutes struct {
	routes
	sceneQueryer models.SceneQueryer
	sceneFinder  models.SceneGetter
}

func (rs deovrRoutes) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/files", rs.ListFiles)
	r.Get("/file/{fileId}", rs.GetFile)

	return r
}

// ListFiles handles GET /deovr/files?search=QUERY
func (rs deovrRoutes) ListFiles(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("search")

	var scenes []*models.Scene
	if err := rs.withReadTxn(r, func(ctx context.Context) error {
		perPage := 200
		findFilter := &models.FindFilterType{
			PerPage: &perPage,
			Sort:    strPtr("title"),
		}

		var sceneFilter *models.SceneFilterType
		if searchQuery != "" {
			sceneFilter = &models.SceneFilterType{
				Title: &models.StringCriterionInput{
					Value:    searchQuery,
					Modifier: models.CriterionModifierIncludes,
				},
			}
		}

		result, err := rs.sceneQueryer.Query(ctx, models.SceneQueryOptions{
			QueryOptions: models.QueryOptions{
				FindFilter: findFilter,
				Count:      false,
			},
			SceneFilter: sceneFilter,
		})
		if err != nil {
			return fmt.Errorf("querying scenes: %w", err)
		}

		scenes, err = result.Resolve(ctx)
		if err != nil {
			return fmt.Errorf("resolving scenes: %w", err)
		}
		logger.Debugf("[DeoVR] listed %d scenes", len(scenes))
		return nil
	}); err != nil {
		logger.Errorf("[DeoVR] error listing files: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	files := make([]DeoVRFile, 0, len(scenes))
	for _, s := range scenes {
		f := deovrSceneToFile(s, r)
		if f != nil {
			files = append(files, *f)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	resp := DeoVRFilesResponse{Files: files}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Errorf("[DeoVR] error encoding response: %v", err)
	}
}

// GetFile handles GET /deovr/file/{fileId}
func (rs deovrRoutes) GetFile(w http.ResponseWriter, r *http.Request) {
	fileIDStr := chi.URLParam(r, "fileId")
	fileID, err := strconv.Atoi(fileIDStr)
	if err != nil {
		http.Error(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	var s *models.Scene
	if err := rs.withReadTxn(r, func(ctx context.Context) error {
		var err error
		s, err = rs.sceneFinder.Find(ctx, fileID)
		return err
	}); err != nil {
		logger.Errorf("[DeoVR] error finding scene %d: %v", fileID, err)
		http.Error(w, "Scene not found", http.StatusNotFound)
		return
	}

	if s == nil {
		http.Error(w, "Scene not found", http.StatusNotFound)
		return
	}

	f := deovrSceneToFile(s, r)
	if f == nil {
		http.Error(w, "Scene has no video files", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(f); err != nil {
		logger.Errorf("[DeoVR] error encoding response: %v", err)
	}
}

// deovrSceneToFile converts a stash scene to a DeoVRFile.
func deovrSceneToFile(s *models.Scene, r *http.Request) *DeoVRFile {
	title := s.Title
	if title == "" {
		title = s.Path
	}

	var fileSize int64
	var width, height int
	var duration float64

	if s.Files != nil && len(s.Files) > 0 {
		f := s.Files[0]
		if f.BaseVideoFile != nil {
			fileSize = f.BaseVideoFile.Size
			width = f.Width
			height = f.Height
		}
		duration = f.Duration
	}

	projection := detectProjection(title)

	return &DeoVRFile{
		FileTitle:       title,
		FileURL:         fmt.Sprintf("%s://%s/scene/%d", scheme(r), r.Host, s.ID),
		FileSize:        fileSize,
		IsAvailable:     true,
		VideoProjection: projection,
		VideoDuration:   duration,
		VideoWidth:      width,
		VideoHeight:     height,
		ThumbnailURL:    fmt.Sprintf("%s://%s/scene/%d/screenshot", scheme(r), r.Host, s.ID),
	}
}

func detectProjection(title string) string {
	t := strings.ToLower(title)
	switch {
	case strings.Contains(t, "360") || strings.Contains(t, "equirectangular"):
		return "EQUIRECTANGULAR"
	case strings.Contains(t, "180") || strings.Contains(t, "fisheye") || strings.Contains(t, "vr"):
		return "FISHEYE"
	case strings.Contains(t, "cube") || strings.Contains(t, "cubemap"):
		return "CUBEMAP"
	default:
		return "FLAT"
	}
}

func scheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}
	return "http"
}

func strPtr(s string) *string {
	return &s
}
