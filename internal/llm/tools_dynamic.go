package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
)

// Dynamic tools let the assistant DEFINE new, reusable tools at runtime — without
// a rebuild/redeploy. A dynamic tool is just a named GraphQL operation plus a
// JSON-Schema for its arguments: when called, the arguments are passed straight
// through as GraphQL variables. Definitions persist as JSON files in a writable
// directory (deps.ToolsDir, on a mounted volume — NOT the container image), so
// they survive restarts and are reloaded on boot.
//
// Safety: a dynamic tool whose GraphQL contains a mutation is force-marked
// Writes:true, so it runs through the same confirm gate as graphql_mutate. The
// meta-tools that manage definitions (define/list/delete) are not library
// mutations, so they are not gated — that is what makes self-extension fluid.

var dynamicToolNameRe = regexp.MustCompile(`^[a-z][a-z0-9_]{2,48}$`)

// DynamicToolSpec is the on-disk definition of an assistant-defined tool.
type DynamicToolSpec struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  json.RawMessage `json:"parameters"` // JSON Schema for the tool's input
	Query       string          `json:"query"`      // GraphQL document; the call's args become its variables
	Writes      bool            `json:"writes"`     // true if the query mutates (auto-detected on define)
}

// RegisterDynamicTools loads persisted specs and registers them, plus the
// define/list/delete meta-tools. No-op for the GraphQL-backed parts if the
// executor isn't wired.
func RegisterDynamicTools(reg *Registry, deps Deps) {
	if deps.GraphQL == nil || deps.Schema == nil || deps.ToolsDir == "" {
		return
	}
	for _, spec := range loadDynamicSpecs(deps.ToolsDir) {
		reg.Register(deps.makeDynamicTool(spec))
	}
	reg.Register(deps.defineToolTool(reg))
	reg.Register(deps.listDynamicToolsTool())
	reg.Register(deps.deleteDynamicToolTool(reg))
}

// makeDynamicTool turns a spec into a runnable Tool. Its handler binds the call
// arguments as GraphQL variables and executes the stored operation.
func (d Deps) makeDynamicTool(spec DynamicToolSpec) *Tool {
	params := spec.Parameters
	if len(params) == 0 {
		params = json.RawMessage(`{"type":"object"}`)
	}
	desc := spec.Description
	if desc == "" {
		desc = "Assistant-defined tool."
	}
	desc += " [assistant-defined]"
	query := spec.Query
	return &Tool{
		Name:        spec.Name,
		Description: desc,
		Schema:      params,
		Writes:      spec.Writes,
		Run: func(ctx context.Context, input json.RawMessage) (string, error) {
			var vars map[string]any
			if len(input) > 0 {
				if err := json.Unmarshal(input, &vars); err != nil {
					return "", fmt.Errorf("bad input for %q: %w", spec.Name, err)
				}
			}
			return d.execGraphQL(ctx, query, vars)
		},
	}
}

// ── define_tool ──────────────────────────────────────────────────────────────

func (d Deps) defineToolTool(reg *Registry) *Tool {
	return &Tool{
		Name: "define_tool",
		Description: "Create (or update) a REUSABLE tool for tasks you expect to repeat, so you don't have to " +
			"hand-write the GraphQL each time. A tool is a named GraphQL query/mutation plus a JSON-Schema for its " +
			"arguments; when the tool is called later, its arguments are passed straight through as GraphQL variables. " +
			"The definition persists across restarts. A tool whose GraphQL mutates is automatically write-gated. " +
			"It becomes callable on your next step. Prefer this for recurring operations; use graphql_query/graphql_mutate for one-offs.",
		Schema: json.RawMessage(`{
			"type":"object",
			"properties":{
				"name":{"type":"string","description":"Tool name: lowercase letters, digits, underscores; 3-49 chars (e.g. 'merge_tags')."},
				"description":{"type":"string","description":"What the tool does and when to use it. Be specific."},
				"parameters":{"type":"object","description":"A JSON Schema (object) describing the tool's arguments. These exact keys are passed as GraphQL variables."},
				"query":{"type":"string","description":"The GraphQL document. Reference the arguments as variables, e.g. 'mutation($i:TagsMergeInput!){tagsMerge(input:$i){id name}}'."}
			},
			"required":["name","description","query"]
		}`),
		Writes: false, // defining a capability isn't a library mutation; execution is what gets gated
		Run: func(ctx context.Context, input json.RawMessage) (string, error) {
			var in struct {
				Name        string          `json:"name"`
				Description string          `json:"description"`
				Parameters  json.RawMessage `json:"parameters"`
				Query       string          `json:"query"`
			}
			if err := json.Unmarshal(input, &in); err != nil {
				return "", fmt.Errorf("bad input: %w", err)
			}
			in.Name = strings.TrimSpace(in.Name)
			if !dynamicToolNameRe.MatchString(in.Name) {
				return "", fmt.Errorf("invalid name %q: use lowercase letters/digits/underscores, 3-49 chars", in.Name)
			}
			if strings.TrimSpace(in.Query) == "" {
				return "", fmt.Errorf("query is required")
			}
			// Don't allow shadowing a built-in tool (one that exists in the
			// registry but has no on-disk dynamic spec).
			if existing, ok := reg.Get(in.Name); ok {
				if _, err := os.Stat(filepath.Join(d.ToolsDir, in.Name+".json")); err != nil {
					return "", fmt.Errorf("%q is a built-in tool and can't be redefined", existing.Name)
				}
			}
			// Validate the GraphQL against the live schema and detect writes.
			kinds, err := d.operationKinds(in.Query)
			if err != nil {
				return "", err
			}
			writes := kinds[ast.Mutation] || kinds[ast.Subscription]

			params := in.Parameters
			if len(params) == 0 {
				params = json.RawMessage(`{"type":"object"}`)
			} else {
				var probe any
				if err := json.Unmarshal(params, &probe); err != nil {
					return "", fmt.Errorf("parameters must be a valid JSON Schema object: %w", err)
				}
			}

			spec := DynamicToolSpec{
				Name:        in.Name,
				Description: in.Description,
				Parameters:  params,
				Query:       in.Query,
				Writes:      writes,
			}
			if err := saveDynamicSpec(d.ToolsDir, spec); err != nil {
				return "", fmt.Errorf("persisting tool: %w", err)
			}
			// Hot-register so it's usable immediately (no restart).
			reg.Register(d.makeDynamicTool(spec))

			gate := "runs immediately when called"
			if writes {
				gate = "write — will require your confirmation when called"
			}
			out, _ := json.Marshal(map[string]any{
				"defined": spec.Name,
				"writes":  writes,
				"note":    fmt.Sprintf("Tool %q saved and now available (%s).", spec.Name, gate),
			})
			return string(out), nil
		},
	}
}

// ── list_dynamic_tools ───────────────────────────────────────────────────────

func (d Deps) listDynamicToolsTool() *Tool {
	return &Tool{
		Name:        "list_dynamic_tools",
		Description: "List the tools you (the assistant) have previously defined, with their descriptions and whether each is a write.",
		Schema:      json.RawMessage(`{"type":"object"}`),
		Writes:      false,
		Run: func(ctx context.Context, _ json.RawMessage) (string, error) {
			specs := loadDynamicSpecs(d.ToolsDir)
			type row struct {
				Name        string `json:"name"`
				Description string `json:"description"`
				Writes      bool   `json:"writes"`
			}
			rows := make([]row, 0, len(specs))
			for _, s := range specs {
				rows = append(rows, row{s.Name, s.Description, s.Writes})
			}
			out, _ := json.Marshal(map[string]any{"count": len(rows), "tools": rows})
			return string(out), nil
		},
	}
}

// ── delete_dynamic_tool ──────────────────────────────────────────────────────

func (d Deps) deleteDynamicToolTool(reg *Registry) *Tool {
	return &Tool{
		Name:        "delete_dynamic_tool",
		Description: "Delete a tool you previously defined with define_tool. Only assistant-defined tools can be removed.",
		Schema: json.RawMessage(`{
			"type":"object",
			"properties":{"name":{"type":"string","description":"Name of the assistant-defined tool to delete."}},
			"required":["name"]
		}`),
		Writes: false,
		Run: func(ctx context.Context, input json.RawMessage) (string, error) {
			var in struct {
				Name string `json:"name"`
			}
			if err := json.Unmarshal(input, &in); err != nil {
				return "", fmt.Errorf("bad input: %w", err)
			}
			in.Name = strings.TrimSpace(in.Name)
			path := filepath.Join(d.ToolsDir, in.Name+".json")
			if !dynamicToolNameRe.MatchString(in.Name) {
				return "", fmt.Errorf("invalid name %q", in.Name)
			}
			if _, err := os.Stat(path); err != nil {
				return "", fmt.Errorf("no assistant-defined tool named %q", in.Name)
			}
			if err := os.Remove(path); err != nil {
				return "", fmt.Errorf("removing tool: %w", err)
			}
			// Replace the live tool with a tombstone so it stops working this
			// session (the registry has no unregister; a clean reload happens on
			// next restart). Returning an error from the handler is the signal.
			reg.Register(&Tool{
				Name:        in.Name,
				Description: fmt.Sprintf("(deleted tool %q)", in.Name),
				Schema:      json.RawMessage(`{"type":"object"}`),
				Run: func(context.Context, json.RawMessage) (string, error) {
					return "", fmt.Errorf("tool %q was deleted", in.Name)
				},
			})
			out, _ := json.Marshal(map[string]any{"deleted": in.Name})
			return string(out), nil
		},
	}
}

// ── persistence ──────────────────────────────────────────────────────────────

func loadDynamicSpecs(dir string) []DynamicToolSpec {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil // dir missing == no tools yet
	}
	var specs []DynamicToolSpec
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		b, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		var s DynamicToolSpec
		if err := json.Unmarshal(b, &s); err != nil {
			continue
		}
		if !dynamicToolNameRe.MatchString(s.Name) {
			continue
		}
		specs = append(specs, s)
	}
	sort.Slice(specs, func(i, j int) bool { return specs[i].Name < specs[j].Name })
	return specs
}

func saveDynamicSpec(dir string, spec DynamicToolSpec) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return err
	}
	// Atomic-ish write.
	tmp := filepath.Join(dir, spec.Name+".json.tmp")
	final := filepath.Join(dir, spec.Name+".json")
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, final)
}
