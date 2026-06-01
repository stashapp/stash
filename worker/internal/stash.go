// Package internal holds the worker's implementation. Stash GraphQL client lives here.
package internal

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// StashClient is a minimal GraphQL client for the stash API. It uses stash's
// ApiKey header convention (not Authorization: Bearer) when a key is configured.
type StashClient struct {
	baseURL string // e.g. http://overwatch-stash:9999
	apiKey  string
	http    *http.Client
}

func NewStashClient(baseURL, apiKey string) *StashClient {
	return &StashClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
		http:    &http.Client{Timeout: 60 * time.Second},
	}
}

type graphqlReq struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

type graphqlResp struct {
	Data   json.RawMessage `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// do runs a GraphQL request and decodes the data field into out.
func (c *StashClient) do(ctx context.Context, query string, vars map[string]any, out any) error {
	body, err := json.Marshal(graphqlReq{Query: query, Variables: vars})
	if err != nil {
		return fmt.Errorf("marshal graphql request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/graphql", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("ApiKey", c.apiKey)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("stash request: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 32<<20))
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("stash HTTP %d: %s", resp.StatusCode, truncate(string(raw), 300))
	}

	var gr graphqlResp
	if err := json.Unmarshal(raw, &gr); err != nil {
		return fmt.Errorf("decode response envelope: %w", err)
	}
	if len(gr.Errors) > 0 {
		msgs := make([]string, 0, len(gr.Errors))
		for _, e := range gr.Errors {
			msgs = append(msgs, e.Message)
		}
		return fmt.Errorf("graphql errors: %s", strings.Join(msgs, "; "))
	}
	if out == nil {
		return nil
	}
	if err := json.Unmarshal(gr.Data, out); err != nil {
		return fmt.Errorf("decode response data: %w", err)
	}
	return nil
}

// ── config ─────────────────────────────────────────────────────────────────

// StashConfig is the subset of stash's configuration the worker needs.
// Field names map to the GraphQL query below; defaults match internal/manager/config/config.go.
type StashConfig struct {
	// VideoFileNamingAlgorithm is the enum stash uses to derive generated file
	// names. Values: "OSHASH" (default) or "MD5". The worker MUST mirror this
	// to compute matching paths.
	VideoFileNamingAlgorithm string

	// Preview knobs — see ../../docs/llm/EXTERNAL-WORKERS.md §4.
	PreviewSegments        int
	PreviewSegmentDuration float64
	PreviewExcludeStart    string // raw form; "30" or "5%"
	PreviewExcludeEnd      string
	PreviewAudio           bool
}

const configQuery = `{
  configuration {
    general {
      videoFileNamingAlgorithm
      previewSegments
      previewSegmentDuration
      previewExcludeStart
      previewExcludeEnd
      previewAudio
    }
  }
}`

type configRespRaw struct {
	Configuration struct {
		General struct {
			VideoFileNamingAlgorithm string  `json:"videoFileNamingAlgorithm"`
			PreviewSegments          int     `json:"previewSegments"`
			PreviewSegmentDuration   float64 `json:"previewSegmentDuration"`
			PreviewExcludeStart      string  `json:"previewExcludeStart"`
			PreviewExcludeEnd        string  `json:"previewExcludeEnd"`
			PreviewAudio             bool    `json:"previewAudio"`
		} `json:"general"`
	} `json:"configuration"`
}

// FetchConfig pulls the preview-related config from stash.
func (c *StashClient) FetchConfig(ctx context.Context) (*StashConfig, error) {
	var raw configRespRaw
	if err := c.do(ctx, configQuery, nil, &raw); err != nil {
		return nil, fmt.Errorf("fetch config: %w", err)
	}
	g := raw.Configuration.General
	if g.PreviewSegments <= 0 {
		g.PreviewSegments = 12
	}
	if g.PreviewSegmentDuration <= 0 {
		g.PreviewSegmentDuration = 0.75
	}
	if g.VideoFileNamingAlgorithm == "" {
		g.VideoFileNamingAlgorithm = "OSHASH"
	}
	return &StashConfig{
		VideoFileNamingAlgorithm: g.VideoFileNamingAlgorithm,
		PreviewSegments:          g.PreviewSegments,
		PreviewSegmentDuration:   g.PreviewSegmentDuration,
		PreviewExcludeStart:      g.PreviewExcludeStart,
		PreviewExcludeEnd:        g.PreviewExcludeEnd,
		PreviewAudio:             g.PreviewAudio,
	}, nil
}

// ── scenes ─────────────────────────────────────────────────────────────────

// Scene is the subset of a stash Scene the worker cares about for preview gen.
type Scene struct {
	ID    string
	Title string
	Files []SceneFile // index 0 is the primary file
}

// SceneFile is the subset of a scene's File the worker cares about.
type SceneFile struct {
	ID          string
	Path        string  // stash-side path (will be translated to worker-side via prefix rewriter)
	Duration    float64 // seconds; from video metadata
	Fingerprints []Fingerprint
}

// Fingerprint is a file fingerprint (oshash, md5, phash, ...).
type Fingerprint struct {
	Type  string // "oshash", "md5", "phash"
	Value string
}

// PrimaryHash returns the fingerprint that matches stash's video_file_naming_algorithm,
// used to derive the generated/ filename for this scene.
func (s *Scene) PrimaryHash(algorithm string) (string, error) {
	if len(s.Files) == 0 {
		return "", errors.New("scene has no files")
	}
	want := strings.ToLower(algorithm) // stash stores fingerprint types lowercase
	for _, fp := range s.Files[0].Fingerprints {
		if strings.EqualFold(fp.Type, want) {
			return fp.Value, nil
		}
	}
	return "", fmt.Errorf("scene %s primary file has no %s fingerprint", s.ID, algorithm)
}

const scenesQuery = `query($filter: FindFilterType, $scene_filter: SceneFilterType) {
  findScenes(filter: $filter, scene_filter: $scene_filter) {
    count
    scenes {
      id
      title
      files {
        id
        path
        duration
        fingerprints { type value }
      }
    }
  }
}`

type scenesRespRaw struct {
	FindScenes struct {
		Count  int `json:"count"`
		Scenes []struct {
			ID    string `json:"id"`
			Title string `json:"title"`
			Files []struct {
				ID           string  `json:"id"`
				Path         string  `json:"path"`
				Duration     float64 `json:"duration"`
				Fingerprints []struct {
					Type  string `json:"type"`
					Value string `json:"value"`
				} `json:"fingerprints"`
			} `json:"files"`
		} `json:"scenes"`
	} `json:"findScenes"`
}

// FetchScenesPage pulls a single page of scenes ordered by id ASC. Page is 1-based.
// missingFilter (e.g. "cover") narrows to scenes stash flags as missing that
// artifact. Empty string = all scenes.
func (c *StashClient) FetchScenesPage(ctx context.Context, page, perPage int, missingFilter string) ([]Scene, int, error) {
	asc := "ASC"
	sort := "id"
	vars := map[string]any{
		"filter": map[string]any{
			"page":      page,
			"per_page":  perPage,
			"sort":      sort,
			"direction": asc,
		},
	}
	if missingFilter != "" {
		vars["scene_filter"] = map[string]any{"is_missing": missingFilter}
	}
	var raw scenesRespRaw
	if err := c.do(ctx, scenesQuery, vars, &raw); err != nil {
		return nil, 0, fmt.Errorf("fetch scenes page %d: %w", page, err)
	}
	out := make([]Scene, 0, len(raw.FindScenes.Scenes))
	for _, s := range raw.FindScenes.Scenes {
		files := make([]SceneFile, 0, len(s.Files))
		for _, f := range s.Files {
			fps := make([]Fingerprint, 0, len(f.Fingerprints))
			for _, fp := range f.Fingerprints {
				fps = append(fps, Fingerprint{Type: fp.Type, Value: fp.Value})
			}
			files = append(files, SceneFile{
				ID: f.ID, Path: f.Path, Duration: f.Duration, Fingerprints: fps,
			})
		}
		out = append(out, Scene{ID: s.ID, Title: s.Title, Files: files})
	}
	return out, raw.FindScenes.Count, nil
}

// ── mutations ──────────────────────────────────────────────────────────────

const sceneUpdateCoverMutation = `mutation($id: ID!, $cover: String!) {
  sceneUpdate(input: { id: $id, cover_image: $cover }) { id }
}`

type sceneUpdateRespRaw struct {
	SceneUpdate struct {
		ID string `json:"id"`
	} `json:"sceneUpdate"`
}

// SceneUpdateCover sets the scene's cover image. The PerformerUpdateInput-style
// `cover_image` field accepts a base64 data-URL ("data:image/jpeg;base64,...");
// stash unpacks it server-side and writes to its blob storage.
// See graphql/schema/types/scene.graphql:151 + repository_scene.go:60.
func (c *StashClient) SceneUpdateCover(ctx context.Context, sceneID string, jpeg []byte) error {
	if len(jpeg) == 0 {
		return errors.New("empty cover bytes")
	}
	dataURL := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(jpeg)
	var raw sceneUpdateRespRaw
	err := c.do(ctx, sceneUpdateCoverMutation, map[string]any{
		"id":    sceneID,
		"cover": dataURL,
	}, &raw)
	if err != nil {
		return fmt.Errorf("sceneUpdate cover for %s: %w", sceneID, err)
	}
	return nil
}

const fileSetFingerprintsMutation = `mutation($id: ID!, $fingerprints: [SetFingerprintsInput!]!) {
  fileSetFingerprints(input: { id: $id, fingerprints: $fingerprints })
}`

// SetFilePhash sets the phash fingerprint on a single file. The value is the
// uint64 hash formatted as lowercase hex — stash's FileSetFingerprints resolver
// parses it with strconv.ParseUint(value, 16, 64) and stores it as int64, and
// its Fingerprint.Value() reads it back as strconv.FormatUint(_, 16). So hex in,
// hex out — the worker matches stash's own round-trip exactly.
//
// "only supplied fingerprint types will be modified" (per the schema), so this
// leaves oshash/md5 untouched.
func (c *StashClient) SetFilePhash(ctx context.Context, fileID string, hash uint64) error {
	if fileID == "" {
		return errors.New("empty file id")
	}
	vars := map[string]any{
		"id": fileID,
		"fingerprints": []map[string]any{
			{"type": "phash", "value": strconv.FormatUint(hash, 16)},
		},
	}
	if err := c.do(ctx, fileSetFingerprintsMutation, vars, nil); err != nil {
		return fmt.Errorf("fileSetFingerprints phash for file %s: %w", fileID, err)
	}
	return nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
