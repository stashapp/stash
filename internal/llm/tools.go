package llm

import (
	"context"
	"encoding/json"
)

// Tool is a single capability exposed to the model. Run executes it and returns a
// compact, model-friendly result string (usually JSON). Writes marks tools that
// mutate the library — these are gated by the write policy (see service.go).
type Tool struct {
	Name        string
	Description string
	Schema      json.RawMessage
	Writes      bool
	Run         func(ctx context.Context, input json.RawMessage) (string, error)
}

// Registry holds the available tools in a stable order.
type Registry struct {
	tools map[string]*Tool
	order []string
}

func NewRegistry() *Registry {
	return &Registry{tools: map[string]*Tool{}}
}

func (r *Registry) Register(t *Tool) {
	if _, exists := r.tools[t.Name]; !exists {
		r.order = append(r.order, t.Name)
	}
	r.tools[t.Name] = t
}

func (r *Registry) Get(name string) (*Tool, bool) {
	t, ok := r.tools[name]
	return t, ok
}

// Defs returns the tool definitions to advertise to the model. When includeWrites
// is false, write tools are omitted entirely (readonly policy).
func (r *Registry) Defs(includeWrites bool) []ToolDef {
	defs := make([]ToolDef, 0, len(r.order))
	for _, name := range r.order {
		t := r.tools[name]
		if t.Writes && !includeWrites {
			continue
		}
		defs = append(defs, ToolDef{
			Name:        t.Name,
			Description: t.Description,
			InputSchema: t.Schema,
		})
	}
	return defs
}

// Event types streamed to the UI over SSE.
const (
	EventText            = "text"             // assistant text chunk
	EventToolCall        = "tool_call"        // model invoked a tool
	EventToolResult      = "tool_result"      // tool finished
	EventConfirmRequired = "confirm_required" // a write needs user confirmation (ask policy)
	EventError           = "error"            // fatal error for this turn
	EventDone            = "done"             // turn complete
)

// Emitter sends a single SSE event. The route layer provides the implementation.
type Emitter func(event string, data any) error
