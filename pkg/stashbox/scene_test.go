package stashbox

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/Yamashou/gqlgenc/clientv2"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/stashbox/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// unsupportedBatchError builds the error a stash-box without batch support
// returns: a GraphQL validation error for the unknown submitFingerprints field.
func unsupportedBatchError() error {
	return &clientv2.ErrorResponse{
		GqlErrors: &gqlerror.List{
			&gqlerror.Error{Message: `Cannot query field "submitFingerprints" on type "Mutation".`},
		},
	}
}

func gqlError(msg string) error {
	return &clientv2.ErrorResponse{
		GqlErrors: &gqlerror.List{&gqlerror.Error{Message: msg}},
	}
}

func TestIsBatchUnsupportedError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{
			"cannot query field",
			gqlError(`Cannot query field "submitFingerprints" on type "Mutation".`),
			true,
		},
		{
			"unknown field",
			gqlError(`Unknown field submitFingerprints`),
			true,
		},
		{
			// a validation error for the legacy mutation must not be mistaken
			// for the batch mutation being unsupported
			"legacy field error",
			gqlError(`Cannot query field "submitFingerprint" on type "Mutation".`),
			false,
		},
		{
			// a transport/network error must not be treated as "unsupported";
			// it should abort rather than silently fall back
			"network error",
			&clientv2.ErrorResponse{NetworkError: &clientv2.HTTPError{Code: 504, Message: "gateway timeout"}},
			false,
		},
		{
			// a plain (non-GraphQL) error must not trigger fallback
			"plain error",
			errors.New("dial tcp: i/o timeout"),
			false,
		},
		{
			// a per-fingerprint error should not trigger fallback
			"scene not found",
			gqlError("scene not found"),
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isBatchUnsupportedError(tt.err); got != tt.want {
				t.Errorf("isBatchUnsupportedError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollectFingerprints(t *testing.T) {
	const endpoint = "https://stashdb.org/graphql"
	const otherEndpoint = "https://other.org/graphql"

	makeScene := func(stashIDs []models.StashID, files []*models.VideoFile) *models.Scene {
		s := &models.Scene{}
		s.StashIDs = models.NewRelatedStashIDs(stashIDs)
		s.Files = models.NewRelatedVideoFiles(files)
		return s
	}

	makeFile := func(duration float64, fps models.Fingerprints) *models.VideoFile {
		return &models.VideoFile{
			BaseFile: &models.BaseFile{
				Fingerprints: fps,
			},
			Duration: duration,
		}
	}

	c := Client{box: models.StashBox{Endpoint: endpoint}}

	t.Run("matching stash id with multiple fingerprint types", func(t *testing.T) {
		scenes := []*models.Scene{
			makeScene(
				[]models.StashID{{Endpoint: endpoint, StashID: "scene-1"}},
				[]*models.VideoFile{makeFile(120, models.Fingerprints{
					{Type: models.FingerprintTypeMD5, Fingerprint: "deadbeef"},
					{Type: models.FingerprintTypePhash, Fingerprint: int64(123456)},
				})},
			),
		}

		got := c.collectFingerprints(scenes)
		if len(got) != 2 {
			t.Fatalf("expected 2 fingerprints, got %d", len(got))
		}
		for _, fp := range got {
			if fp.SceneID != "scene-1" {
				t.Errorf("expected scene id scene-1, got %s", fp.SceneID)
			}
			if fp.Fingerprint.Duration != 120 {
				t.Errorf("expected duration 120, got %d", fp.Fingerprint.Duration)
			}
		}
	})

	t.Run("scene without matching endpoint is skipped", func(t *testing.T) {
		scenes := []*models.Scene{
			makeScene(
				[]models.StashID{{Endpoint: otherEndpoint, StashID: "scene-1"}},
				[]*models.VideoFile{makeFile(120, models.Fingerprints{
					{Type: models.FingerprintTypeMD5, Fingerprint: "deadbeef"},
				})},
			),
		}

		if got := c.collectFingerprints(scenes); len(got) != 0 {
			t.Fatalf("expected 0 fingerprints, got %d", len(got))
		}
	})

	t.Run("zero duration files are skipped", func(t *testing.T) {
		scenes := []*models.Scene{
			makeScene(
				[]models.StashID{{Endpoint: endpoint, StashID: "scene-1"}},
				[]*models.VideoFile{makeFile(0, models.Fingerprints{
					{Type: models.FingerprintTypeMD5, Fingerprint: "deadbeef"},
				})},
			),
		}

		if got := c.collectFingerprints(scenes); len(got) != 0 {
			t.Fatalf("expected 0 fingerprints, got %d", len(got))
		}
	})

	t.Run("unknown fingerprint type is skipped", func(t *testing.T) {
		scenes := []*models.Scene{
			makeScene(
				[]models.StashID{{Endpoint: endpoint, StashID: "scene-1"}},
				[]*models.VideoFile{makeFile(120, models.Fingerprints{
					{Type: "unknown", Fingerprint: "deadbeef"},
				})},
			),
		}

		if got := c.collectFingerprints(scenes); len(got) != 0 {
			t.Fatalf("expected 0 fingerprints, got %d", len(got))
		}
	})
}

// fakeSubmitter implements fingerprintSubmitter for testing the batch/fallback
// submission logic without a network client.
type fakeSubmitter struct {
	// batchErr, if set, is returned by every SubmitFingerprints call
	batchErr error
	// batchItemErrors maps a fingerprint hash to a server-reported error in batch mode
	batchItemErrors map[string]string
	// singleErr maps a fingerprint hash to an error from SubmitFingerprint (legacy mode)
	singleErr map[string]error

	batchCalls  int
	singleCalls int
	batchSizes  []int
}

func (f *fakeSubmitter) SubmitFingerprints(ctx context.Context, input []*graphql.FingerprintBatchSubmission, interceptors ...clientv2.RequestInterceptor) (*graphql.SubmitFingerprints, error) {
	f.batchCalls++
	f.batchSizes = append(f.batchSizes, len(input))
	if f.batchErr != nil {
		return nil, f.batchErr
	}

	resp := &graphql.SubmitFingerprints{}
	for _, in := range input {
		r := &graphql.SubmitFingerprints_SubmitFingerprints{Hash: in.Hash, SceneID: in.SceneID}
		if msg, ok := f.batchItemErrors[in.Hash]; ok {
			m := msg
			r.Error = &m
		}
		resp.SubmitFingerprints = append(resp.SubmitFingerprints, r)
	}
	return resp, nil
}

func (f *fakeSubmitter) SubmitFingerprint(ctx context.Context, input graphql.FingerprintSubmission, interceptors ...clientv2.RequestInterceptor) (*graphql.SubmitFingerprint, error) {
	f.singleCalls++
	if err := f.singleErr[input.Fingerprint.Hash]; err != nil {
		return nil, err
	}
	return &graphql.SubmitFingerprint{}, nil
}

type fakeProgress struct {
	total     int
	processed int
}

func (p *fakeProgress) SetTotal(t int) { p.total = t }
func (p *fakeProgress) Increment()     { p.processed++ }

func makeTestFingerprints(n int) []graphql.FingerprintSubmission {
	out := make([]graphql.FingerprintSubmission, n)
	for i := range out {
		out[i] = graphql.FingerprintSubmission{
			SceneID: "scene",
			Fingerprint: &graphql.FingerprintInput{
				Hash:      "hash-" + strconv.Itoa(i),
				Algorithm: graphql.FingerprintAlgorithmMd5,
				Duration:  100,
			},
		}
	}
	return out
}

func TestSubmitFingerprints(t *testing.T) {
	t.Run("batches in chunks and reports progress", func(t *testing.T) {
		f := &fakeSubmitter{}
		p := &fakeProgress{}

		// 120 fingerprints should be split into 50 + 50 + 20
		result, err := submitFingerprints(context.Background(), f, makeTestFingerprints(120), p)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if f.batchCalls != 3 || f.singleCalls != 0 {
			t.Errorf("expected 3 batch calls and 0 single calls, got %d/%d", f.batchCalls, f.singleCalls)
		}
		if want := []int{50, 50, 20}; !equalInts(f.batchSizes, want) {
			t.Errorf("expected batch sizes %v, got %v", want, f.batchSizes)
		}
		if result.Succeeded != 120 || result.Failed != 0 || result.Total != 120 {
			t.Errorf("unexpected result %+v", result)
		}
		if p.total != 120 || p.processed != 120 {
			t.Errorf("expected progress total/processed 120/120, got %d/%d", p.total, p.processed)
		}
	})

	t.Run("records per-fingerprint batch errors without aborting", func(t *testing.T) {
		f := &fakeSubmitter{batchItemErrors: map[string]string{"hash-1": "scene not found"}}

		result, err := submitFingerprints(context.Background(), f, makeTestFingerprints(3), nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Succeeded != 2 || result.Failed != 1 {
			t.Errorf("expected 2 succeeded / 1 failed, got %+v", result)
		}
	})

	t.Run("falls back to legacy submission when batch is unsupported", func(t *testing.T) {
		f := &fakeSubmitter{
			batchErr: unsupportedBatchError(),
		}

		// 60 fingerprints: the first batch probes (and fails), then all 60 are
		// submitted individually
		result, err := submitFingerprints(context.Background(), f, makeTestFingerprints(60), nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if f.batchCalls != 1 {
			t.Errorf("expected the batch mutation to be attempted once, got %d", f.batchCalls)
		}
		if f.singleCalls != 60 {
			t.Errorf("expected 60 single submissions, got %d", f.singleCalls)
		}
		if result.Succeeded != 60 {
			t.Errorf("expected 60 succeeded, got %+v", result)
		}
	})

	t.Run("aborts on a hard batch error", func(t *testing.T) {
		f := &fakeSubmitter{batchErr: errors.New("502 bad gateway")}

		_, err := submitFingerprints(context.Background(), f, makeTestFingerprints(50), nil)
		if err == nil {
			t.Fatal("expected an error to be returned")
		}
		if f.singleCalls != 0 {
			t.Errorf("expected no fallback on a hard error, got %d single calls", f.singleCalls)
		}
	})

	t.Run("legacy submission continues past individual failures", func(t *testing.T) {
		f := &fakeSubmitter{
			batchErr:  unsupportedBatchError(),
			singleErr: map[string]error{"hash-1": errors.New("rejected")},
		}

		result, err := submitFingerprints(context.Background(), f, makeTestFingerprints(3), nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Succeeded != 2 || result.Failed != 1 {
			t.Errorf("expected 2 succeeded / 1 failed, got %+v", result)
		}
	})

	t.Run("aborts legacy submission on context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		f := &fakeSubmitter{
			batchErr:  unsupportedBatchError(),
			singleErr: map[string]error{"hash-0": context.Canceled},
		}

		_, err := submitFingerprints(ctx, f, makeTestFingerprints(5), nil)
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
		if f.singleCalls != 1 {
			t.Errorf("expected submission to stop after the first item, got %d", f.singleCalls)
		}
	})
}

func equalInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// fakeStashBox records the requests it receives and emulates either a stash-box
// that supports the batch mutation or an older one that doesn't.
type fakeStashBox struct {
	batchSupported bool

	mu          sync.Mutex
	batchCalls  int
	singleCalls int
	batchSizes  []int
}

func writeRaw(w http.ResponseWriter, body string) {
	_, _ = io.WriteString(w, body)
}

func (s *fakeStashBox) handler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)

	var req struct {
		OperationName string `json:"operationName"`
		Variables     struct {
			Input json.RawMessage `json:"input"`
		} `json:"variables"`
	}
	_ = json.Unmarshal(body, &req)

	w.Header().Set("Content-Type", "application/json")

	switch req.OperationName {
	case "SubmitFingerprints":
		var input []json.RawMessage
		_ = json.Unmarshal(req.Variables.Input, &input)

		s.mu.Lock()
		s.batchCalls++
		s.batchSizes = append(s.batchSizes, len(input))
		s.mu.Unlock()

		if !s.batchSupported {
			// emulate the GraphQL validation error an older stash-box returns
			writeRaw(w, `{"errors":[{"message":"Cannot query field \"submitFingerprints\" on type \"Mutation\"."}]}`)
			return
		}

		// return one success result per submitted fingerprint
		item := `{"scene_id":"","hash":"","error":null}`
		list := strings.TrimSuffix(strings.Repeat(item+",", len(input)), ",")
		writeRaw(w, `{"data":{"submitFingerprints":[`+list+`]}}`)

	case "SubmitFingerprint":
		s.mu.Lock()
		s.singleCalls++
		s.mu.Unlock()
		writeRaw(w, `{"data":{"submitFingerprint":true}}`)

	default:
		http.Error(w, "unexpected operation", http.StatusBadRequest)
	}
}

// sceneWithFingerprints builds a scene matched to the given endpoint with one
// file carrying md5/oshash/phash fingerprints (3 fingerprints).
func sceneWithFingerprints(endpoint string) *models.Scene {
	s := &models.Scene{}
	s.StashIDs = models.NewRelatedStashIDs([]models.StashID{{Endpoint: endpoint, StashID: "scene-1"}})
	s.Files = models.NewRelatedVideoFiles([]*models.VideoFile{{
		BaseFile: &models.BaseFile{
			Fingerprints: models.Fingerprints{
				{Type: models.FingerprintTypeMD5, Fingerprint: "abc123"},
				{Type: models.FingerprintTypeOshash, Fingerprint: "def456"},
				{Type: models.FingerprintTypePhash, Fingerprint: int64(789)},
			},
		},
		Duration: 100,
	}})
	return s
}

// TestSubmitFingerprintsOverHTTP drives the real stash-box client (request
// serialization, response/error parsing) against an in-process fake server, so
// the network path is exercised without any real server or fingerprint data.
func TestSubmitFingerprintsOverHTTP(t *testing.T) {
	t.Run("batch-capable server", func(t *testing.T) {
		fake := &fakeStashBox{batchSupported: true}
		server := httptest.NewServer(http.HandlerFunc(fake.handler))
		defer server.Close()

		client := NewClient(models.StashBox{Endpoint: server.URL, MaxRequestsPerMinute: 60000})

		result, err := client.SubmitFingerprints(context.Background(), []*models.Scene{sceneWithFingerprints(server.URL)}, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Succeeded != 3 || result.Failed != 0 {
			t.Errorf("expected 3 succeeded / 0 failed, got %+v", result)
		}
		if fake.batchCalls != 1 || fake.singleCalls != 0 {
			t.Errorf("expected 1 batch call and 0 single calls, got %d/%d", fake.batchCalls, fake.singleCalls)
		}
		if want := []int{3}; !equalInts(fake.batchSizes, want) {
			t.Errorf("expected batch sizes %v, got %v", want, fake.batchSizes)
		}
	})

	t.Run("server without batch support falls back to legacy", func(t *testing.T) {
		fake := &fakeStashBox{batchSupported: false}
		server := httptest.NewServer(http.HandlerFunc(fake.handler))
		defer server.Close()

		client := NewClient(models.StashBox{Endpoint: server.URL, MaxRequestsPerMinute: 60000})

		// this proves the real clientv2 error parsing flows into
		// isBatchUnsupportedError and triggers the legacy fallback end-to-end
		result, err := client.SubmitFingerprints(context.Background(), []*models.Scene{sceneWithFingerprints(server.URL)}, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Succeeded != 3 {
			t.Errorf("expected 3 succeeded, got %+v", result)
		}
		if fake.batchCalls != 1 {
			t.Errorf("expected the batch mutation to be probed once, got %d", fake.batchCalls)
		}
		if fake.singleCalls != 3 {
			t.Errorf("expected 3 legacy submissions, got %d", fake.singleCalls)
		}
	})
}
