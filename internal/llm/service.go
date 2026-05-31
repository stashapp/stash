package llm

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/stashapp/stash/internal/manager/config"
)

const maxTurns = 12

const systemPrompt = `You are the built-in assistant for Stash, a self-hosted media library server.
You help the user explore and curate their library by calling the provided tools.

Guidelines:
- Use the tools to answer questions about the library; never invent ids, titles, or counts.
- Resolve names to ids first (find_scenes / find_tags / find_performers / find_studios) before any
  bulk action, and tell the user how many items a change will affect before doing it.
- Ratings are on a 1-100 scale. Dates are ISO (YYYY-MM-DD).
- Be concise. Prefer a short summary over dumping raw rows.
- Some write actions may require the user to confirm in the UI. If a tool result says confirmation is
  required, do not retry — explain what you intend to change and ask the user to approve it.`

var (
	// ErrDisabled is returned when the assistant is turned off in settings.
	ErrDisabled = errors.New("assistant is disabled")
	// ErrNotConfigured is returned when no gateway base URL is configured.
	ErrNotConfigured = errors.New("assistant gateway is not configured (set assistant_base_url)")
)

// Service drives conversations: it owns the tool registry and an in-memory
// conversation store. Gateway base URL / key / model are read from config per
// request, so settings changes take effect without a restart.
type Service struct {
	registry *Registry
	convs    *convStore
}

// NewService builds the service and registers the Phase 1 library tools.
func NewService(deps Deps) *Service {
	reg := NewRegistry()
	RegisterLibraryTools(reg, deps)
	return &Service{registry: reg, convs: newConvStore()}
}

type settings struct {
	enabled     bool
	baseURL     string
	apiKey      string
	model       string
	writePolicy string
}

func currentSettings() settings {
	c := config.GetInstance()
	return settings{
		enabled:     c.GetAssistantEnabled(),
		baseURL:     c.GetAssistantBaseURL(),
		apiKey:      c.GetAssistantAPIKey(),
		model:       c.GetAssistantModel(),
		writePolicy: c.GetAssistantWritePolicy(),
	}
}

// Status reports assistant availability for the UI (never exposes the key).
type Status struct {
	Enabled     bool     `json:"enabled"`
	Configured  bool     `json:"configured"`
	Model       string   `json:"model"`
	WritePolicy string   `json:"write_policy"`
	Tools       []string `json:"tools"`
}

func (s *Service) Status() Status {
	st := currentSettings()
	return Status{
		Enabled:     st.enabled,
		Configured:  st.baseURL != "",
		Model:       st.model,
		WritePolicy: st.writePolicy,
		Tools:       s.registry.defNames(st.writePolicy != "readonly"),
	}
}

// Chat runs one user turn through the agent loop, emitting SSE events as it goes.
func (s *Service) Chat(ctx context.Context, convID, userMessage string, emit Emitter) error {
	st := currentSettings()
	if !st.enabled {
		return ErrDisabled
	}
	if st.baseURL == "" {
		return ErrNotConfigured
	}
	client := NewClient(st.baseURL, st.apiKey, st.model)

	if convID == "" {
		convID = newID()
	}
	msgs := s.convs.get(convID)
	if len(msgs) == 0 {
		msgs = append(msgs, Message{Role: "system", Content: systemPrompt})
	}
	msgs = append(msgs, Message{Role: "user", Content: userMessage})

	tools := s.registry.Defs(st.writePolicy != "readonly")

	for turn := 0; turn < maxTurns; turn++ {
		resp, err := client.CreateChatCompletion(ctx, msgs, tools)
		if err != nil {
			return err
		}
		if len(resp.Choices) == 0 {
			return emit(EventError, map[string]string{"error": "gateway returned no choices"})
		}
		choice := resp.Choices[0]
		msg := choice.Message
		msg.Role = "assistant"
		msgs = append(msgs, msg)

		if msg.Content != "" {
			_ = emit(EventText, map[string]string{"text": msg.Content})
		}

		if choice.FinishReason != "tool_calls" || len(msg.ToolCalls) == 0 {
			s.convs.set(convID, msgs)
			return emit(EventDone, map[string]any{"conversationId": convID, "usage": resp.Usage})
		}

		for _, tc := range msg.ToolCalls {
			args := json.RawMessage(tc.Function.Arguments)
			_ = emit(EventToolCall, map[string]any{"id": tc.ID, "name": tc.Function.Name, "input": args})
			out, isErr := s.runTool(ctx, tc, st.writePolicy, emit)
			_ = emit(EventToolResult, map[string]any{
				"id": tc.ID, "name": tc.Function.Name, "is_error": isErr, "summary": truncate(out, 600),
			})
			msgs = append(msgs, Message{Role: "tool", ToolCallID: tc.ID, Content: out})
		}
	}

	s.convs.set(convID, msgs)
	return emit(EventError, map[string]string{"error": "reached the maximum number of reasoning steps"})
}

func (s *Service) runTool(ctx context.Context, tc ToolCall, writePolicy string, emit Emitter) (out string, isErr bool) {
	tool, ok := s.registry.Get(tc.Function.Name)
	if !ok {
		return fmt.Sprintf("unknown tool %q", tc.Function.Name), true
	}
	args := tc.Function.Arguments
	if args == "" {
		args = "{}"
	}
	if tool.Writes && writePolicy == "ask" {
		_ = emit(EventConfirmRequired, map[string]any{"id": tc.ID, "name": tc.Function.Name, "input": json.RawMessage(args)})
		return "This write action requires user confirmation (write policy = ask). Do not retry it; " +
			"briefly tell the user what you will change and ask them to approve it in the UI.", false
	}
	res, err := tool.Run(ctx, json.RawMessage(args))
	if err != nil {
		return err.Error(), true
	}
	return res, false
}

// ExecuteConfirmed runs a single tool directly, bypassing the ask gate. It is used
// by the /llm/confirm endpoint after the user has approved a proposed write.
func (s *Service) ExecuteConfirmed(ctx context.Context, name string, input json.RawMessage) (string, error) {
	st := currentSettings()
	if !st.enabled {
		return "", ErrDisabled
	}
	if st.writePolicy == "readonly" {
		return "", errors.New("write policy is readonly")
	}
	tool, ok := s.registry.Get(name)
	if !ok {
		return "", fmt.Errorf("unknown tool %q", name)
	}
	if len(input) == 0 {
		input = json.RawMessage("{}")
	}
	return tool.Run(ctx, input)
}

// ── conversation store ───────────────────────────────────────────────────────

type convStore struct {
	mu    sync.Mutex
	convs map[string][]Message
}

func newConvStore() *convStore {
	return &convStore{convs: map[string][]Message{}}
}

func (c *convStore) get(id string) []Message {
	if id == "" {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	src := c.convs[id]
	out := make([]Message, len(src))
	copy(out, src)
	return out
}

func (c *convStore) set(id string, msgs []Message) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.convs[id] = msgs
}

func newID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
