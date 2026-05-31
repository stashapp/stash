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
	// ErrNoAPIKey is returned when no Anthropic API key is configured.
	ErrNoAPIKey = errors.New("no Anthropic API key configured")
)

// Service drives conversations: it owns the tool registry and an in-memory
// conversation store. Provider key/model are read from config per request.
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
	apiKey      string
	model       string
	writePolicy string
}

func currentSettings() settings {
	c := config.GetInstance()
	return settings{
		enabled:     c.GetAssistantEnabled(),
		apiKey:      c.GetAnthropicAPIKey(),
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
	defs := s.registry.Defs(st.writePolicy != "readonly")
	tools := make([]string, 0, len(defs))
	for _, d := range defs {
		tools = append(tools, d.Name)
	}
	return Status{
		Enabled:     st.enabled,
		Configured:  st.apiKey != "",
		Model:       st.model,
		WritePolicy: st.writePolicy,
		Tools:       tools,
	}
}

// Chat runs one user turn through the agent loop, emitting SSE events as it goes.
// It returns the (possibly new) conversation id via the EventDone payload.
func (s *Service) Chat(ctx context.Context, convID, userMessage string, emit Emitter) error {
	st := currentSettings()
	if !st.enabled {
		return ErrDisabled
	}
	if st.apiKey == "" {
		return ErrNoAPIKey
	}
	client := NewClient(st.apiKey, st.model)

	if convID == "" {
		convID = newID()
	}
	msgs := s.convs.get(convID)
	msgs = append(msgs, Message{Role: "user", Content: []ContentBlock{{Type: "text", Text: userMessage}}})

	tools := s.registry.Defs(st.writePolicy != "readonly")

	for turn := 0; turn < maxTurns; turn++ {
		resp, err := client.CreateMessage(ctx, systemPrompt, msgs, tools)
		if err != nil {
			return err
		}
		msgs = append(msgs, Message{Role: "assistant", Content: resp.Content})

		var toolUses []ContentBlock
		for _, b := range resp.Content {
			switch b.Type {
			case "text":
				if b.Text != "" {
					_ = emit(EventText, map[string]string{"text": b.Text})
				}
			case "tool_use":
				toolUses = append(toolUses, b)
			}
		}

		if resp.StopReason != "tool_use" || len(toolUses) == 0 {
			s.convs.set(convID, msgs)
			return emit(EventDone, map[string]any{"conversationId": convID, "usage": resp.Usage})
		}

		results := make([]ContentBlock, 0, len(toolUses))
		for _, tu := range toolUses {
			_ = emit(EventToolCall, map[string]any{"id": tu.ID, "name": tu.Name, "input": tu.Input})
			out, isErr := s.runTool(ctx, tu, st.writePolicy, emit)
			_ = emit(EventToolResult, map[string]any{
				"id": tu.ID, "name": tu.Name, "is_error": isErr, "summary": truncate(out, 600),
			})
			results = append(results, ContentBlock{
				Type: "tool_result", ToolUseID: tu.ID, Content: out, IsError: isErr,
			})
		}
		msgs = append(msgs, Message{Role: "user", Content: results})
	}

	s.convs.set(convID, msgs)
	return emit(EventError, map[string]string{"error": "reached the maximum number of reasoning steps"})
}

func (s *Service) runTool(ctx context.Context, tu ContentBlock, writePolicy string, emit Emitter) (out string, isErr bool) {
	tool, ok := s.registry.Get(tu.Name)
	if !ok {
		return fmt.Sprintf("unknown tool %q", tu.Name), true
	}
	if tool.Writes && writePolicy == "ask" {
		_ = emit(EventConfirmRequired, map[string]any{"id": tu.ID, "name": tu.Name, "input": tu.Input})
		return "This write action requires user confirmation (write policy = ask). Do not retry it; " +
			"briefly tell the user what you will change and ask them to approve it in the UI.", false
	}
	res, err := tool.Run(ctx, tu.Input)
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
