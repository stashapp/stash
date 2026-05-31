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
	if len(defs) != 1 || defs[0].Function.Name != "read1" {
		t.Fatalf("Defs(false): want only read1, got %+v", defs)
	}
	// OpenAI function shape
	if defs[0].Type != "function" || string(defs[0].Function.Parameters) != `{"type":"object"}` {
		t.Fatalf("unexpected function def: %+v", defs[0])
	}
	// stable registration order
	if got := reg.defNames(true); got[0] != "read1" || got[1] != "write1" {
		t.Fatalf("expected stable order, got %v", got)
	}
}

func TestRunToolSafelyRecoversPanic(t *testing.T) {
	panicTool := &Tool{
		Name: "boom",
		Run: func(ctx context.Context, _ json.RawMessage) (string, error) {
			var p *int
			_ = *p // nil deref, like querying an uninitialized DB
			return "unreachable", nil
		},
	}
	// must not panic; must surface an error
	out, err := runToolSafely(context.Background(), panicTool, json.RawMessage(`{}`))
	if err == nil {
		t.Fatalf("expected error from panicking tool, got out=%q", out)
	}
	if !strings.Contains(err.Error(), "boom") {
		t.Fatalf("error should name the tool: %v", err)
	}

	// and via the service runTool path: a panic becomes is_error, not a crash
	s := &Service{registry: NewRegistry(), convs: newConvStore()}
	s.registry.Register(panicTool)
	tc := ToolCall{ID: "c1", Type: "function", Function: FunctionCall{Name: "boom", Arguments: "{}"}}
	gotOut, isErr := s.runTool(context.Background(), tc, "auto", func(string, any) error { return nil })
	if !isErr || !strings.Contains(gotOut, "boom") {
		t.Fatalf("runTool should report a graceful error, got isErr=%v out=%q", isErr, gotOut)
	}
}

func TestClientChatCompletionToolCall(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		if r.Header.Get("authorization") != "Bearer test-key" {
			t.Errorf("missing/incorrect auth header: %q", r.Header.Get("authorization"))
		}
		var req chatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("decode request: %v", err)
		}
		if req.Model != "test-model" {
			t.Errorf("model not propagated: %q", req.Model)
		}
		if len(req.Tools) != 1 || req.Tools[0].Function.Name != "find_scenes" || req.ToolChoice != "auto" {
			t.Errorf("tools not propagated: %+v choice=%q", req.Tools, req.ToolChoice)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"x","model":"test-model","choices":[{"index":0,` +
			`"message":{"role":"assistant","content":"on it","tool_calls":[{"id":"call_1","type":"function",` +
			`"function":{"name":"find_scenes","arguments":"{\"query\":\"foo\"}"}}]},"finish_reason":"tool_calls"}],` +
			`"usage":{"prompt_tokens":5,"completion_tokens":7,"total_tokens":12}}`))
	}))
	defer ts.Close()

	c := NewClient(ts.URL+"/v1", "test-key", "test-model")
	resp, err := c.CreateChatCompletion(context.Background(),
		[]Message{{Role: "user", Content: "hello"}},
		[]ToolDef{{Type: "function", Function: FunctionDef{Name: "find_scenes", Parameters: json.RawMessage(`{"type":"object"}`)}}},
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Choices) != 1 || resp.Choices[0].FinishReason != "tool_calls" {
		t.Fatalf("bad choices: %+v", resp.Choices)
	}
	tcs := resp.Choices[0].Message.ToolCalls
	if len(tcs) != 1 || tcs[0].Function.Name != "find_scenes" || tcs[0].Function.Arguments != `{"query":"foo"}` {
		t.Fatalf("tool_call not decoded: %+v", tcs)
	}
	if resp.Usage.CompletionTokens != 7 {
		t.Fatalf("usage not decoded: %+v", resp.Usage)
	}
}

func TestClientBaseURLTrailingSlash(t *testing.T) {
	c := NewClient("http://gw:4000/v1/", "", "m")
	if c.baseURL != "http://gw:4000/v1" {
		t.Fatalf("trailing slash not trimmed: %q", c.baseURL)
	}
}

func TestClientChatCompletionAPIError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":{"type":"invalid_request_error","message":"bad model"}}`))
	}))
	defer ts.Close()

	c := NewClient(ts.URL, "", "m")
	_, err := c.CreateChatCompletion(context.Background(), nil, nil)
	if err == nil || !strings.Contains(err.Error(), "bad model") {
		t.Fatalf("want surfaced api error, got %v", err)
	}
}

// A role:"tool" result message must marshal with tool_call_id + content and omit
// tool_calls so the gateway accepts the follow-up message.
func TestToolResultMessageMarshal(t *testing.T) {
	m := Message{Role: "tool", ToolCallID: "call_1", Content: `{"count":3}`}
	out, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	for _, want := range []string{`"role":"tool"`, `"tool_call_id":"call_1"`, `"content":"{\"count\":3}"`} {
		if !strings.Contains(s, want) {
			t.Fatalf("marshalled tool message missing %s: %s", want, s)
		}
	}
	if strings.Contains(s, "tool_calls") {
		t.Fatalf("tool message should omit tool_calls: %s", s)
	}
}
