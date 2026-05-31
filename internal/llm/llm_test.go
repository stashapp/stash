package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegistryDefsFiltersWrites(t *testing.T) {
	reg := NewRegistry()
	reg.Register(&Tool{Name: "read1", Schema: json.RawMessage(`{"type":"object"}`)})
	reg.Register(&Tool{Name: "write1", Writes: true, Schema: json.RawMessage(`{"type":"object"}`)})

	if got := len(reg.Defs(true)); got != 2 {
		t.Fatalf("Defs(true): want 2, got %d", got)
	}
	defs := reg.Defs(false)
	if len(defs) != 1 || defs[0].Name != "read1" {
		t.Fatalf("Defs(false): want only read1, got %+v", defs)
	}

	// registration order is preserved
	if reg.Defs(true)[0].Name != "read1" {
		t.Fatalf("expected stable order, got %+v", reg.Defs(true))
	}
}

func TestClientCreateMessageToolUse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-api-key") != "test-key" {
			t.Errorf("missing x-api-key header")
		}
		if r.Header.Get("anthropic-version") == "" {
			t.Errorf("missing anthropic-version header")
		}
		var req messagesRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("decode request: %v", err)
		}
		if req.Model != "test-model" {
			t.Errorf("model not propagated: %q", req.Model)
		}
		if req.System != "sys" {
			t.Errorf("system not propagated: %q", req.System)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"x","model":"test-model","role":"assistant","stop_reason":"tool_use",` +
			`"content":[{"type":"text","text":"hi"},{"type":"tool_use","id":"t1","name":"find_scenes","input":{"query":"foo"}}],` +
			`"usage":{"input_tokens":5,"output_tokens":7}}`))
	}))
	defer ts.Close()

	c := NewClient("test-key", "test-model")
	c.baseURL = ts.URL

	resp, err := c.CreateMessage(context.Background(), "sys",
		[]Message{{Role: "user", Content: []ContentBlock{{Type: "text", Text: "hello"}}}},
		[]ToolDef{{Name: "find_scenes", InputSchema: json.RawMessage(`{"type":"object"}`)}},
	)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StopReason != "tool_use" {
		t.Fatalf("stop_reason: %q", resp.StopReason)
	}
	if len(resp.Content) != 2 {
		t.Fatalf("want 2 content blocks, got %d", len(resp.Content))
	}
	tu := resp.Content[1]
	if tu.Type != "tool_use" || tu.Name != "find_scenes" {
		t.Fatalf("bad tool_use block: %+v", tu)
	}
	if string(tu.Input) != `{"query":"foo"}` {
		t.Fatalf("tool input not preserved: %s", tu.Input)
	}
	if resp.Usage.OutputTokens != 7 {
		t.Fatalf("usage not decoded: %+v", resp.Usage)
	}
}

func TestClientCreateMessageAPIError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"type":"error","error":{"type":"invalid_request_error","message":"bad model"}}`))
	}))
	defer ts.Close()

	c := NewClient("k", "m")
	c.baseURL = ts.URL
	_, err := c.CreateMessage(context.Background(), "", nil, nil)
	if err == nil || !strings.Contains(err.Error(), "bad model") {
		t.Fatalf("want surfaced api error, got %v", err)
	}
}

// tool_result blocks must marshal with tool_use_id + content and omit the text/id
// fields so the Anthropic API accepts the follow-up message.
func TestToolResultBlockMarshal(t *testing.T) {
	b := ContentBlock{Type: "tool_result", ToolUseID: "t1", Content: `{"count":3}`}
	out, err := json.Marshal(b)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	for _, want := range []string{`"type":"tool_result"`, `"tool_use_id":"t1"`, `"content":"{\"count\":3}"`} {
		if !strings.Contains(s, want) {
			t.Fatalf("marshalled block missing %s: %s", want, s)
		}
	}
	if strings.Contains(s, `"text"`) || strings.Contains(s, `"name"`) {
		t.Fatalf("tool_result block should omit text/name: %s", s)
	}
}
