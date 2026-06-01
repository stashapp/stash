package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	gqlparser "github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

// Generic GraphQL tools. These give the assistant direct, schema-validated access
// to stash's OWN GraphQL API — the same surface its web UI uses. Almost every
// curation operation already exists there (tagsMerge, bulkSceneUpdate,
// performerMerge, …), so the assistant can do things we never hand-coded a Go
// tool for, and new capabilities never require a rebuild/redeploy.
//
// Safety: graphql_query refuses mutation/subscription operations, so reads stay
// un-gated. graphql_mutate is marked Writes:true and therefore runs through the
// write-policy confirm gate (the user approves the exact GraphQL in the UI before
// it executes).

const maxToolResultBytes = 16000

// RegisterGraphQLTools adds the introspection + query + mutate tools. No-op if
// the GraphQL executor wasn't wired (deps.GraphQL == nil).
func RegisterGraphQLTools(reg *Registry, deps Deps) {
	if deps.GraphQL == nil || deps.Schema == nil {
		return
	}
	reg.Register(deps.graphqlSchemaTool())
	reg.Register(deps.graphqlQueryTool())
	reg.Register(deps.graphqlMutateTool())
}

// ── execution helper ─────────────────────────────────────────────────────────

// operationKinds parses + validates a document against the schema and returns the
// set of operation types it contains ("query"/"mutation"/"subscription"). A parse
// or validation error is returned verbatim so the model can correct its GraphQL.
func (d Deps) operationKinds(query string) (map[ast.Operation]bool, error) {
	doc, errs := gqlparser.LoadQuery(d.Schema, query)
	if errs != nil {
		return nil, fmt.Errorf("invalid GraphQL: %s", errs.Error())
	}
	kinds := map[ast.Operation]bool{}
	for _, op := range doc.Operations {
		kinds[op.Operation] = true
	}
	return kinds, nil
}

// execGraphQL runs the operation and renders a compact, model-friendly result.
// GraphQL-level errors in the response are surfaced as a tool error.
func (d Deps) execGraphQL(ctx context.Context, query string, vars map[string]any) (string, error) {
	raw, err := d.GraphQL(ctx, query, vars)
	if err != nil {
		return "", fmt.Errorf("executing GraphQL: %w", err)
	}
	var resp struct {
		Data   json.RawMessage `json:"data"`
		Errors []struct {
			Message string `json:"message"`
			Path    []any  `json:"path,omitempty"`
		} `json:"errors,omitempty"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		// Not the shape we expected — hand back the raw body, truncated.
		return capResult(string(raw)), nil
	}
	if len(resp.Errors) > 0 {
		msgs := make([]string, 0, len(resp.Errors))
		for _, e := range resp.Errors {
			if len(e.Path) > 0 {
				msgs = append(msgs, fmt.Sprintf("%s (at %v)", e.Message, e.Path))
			} else {
				msgs = append(msgs, e.Message)
			}
		}
		return "", fmt.Errorf("GraphQL errors: %s", strings.Join(msgs, "; "))
	}
	if len(resp.Data) == 0 {
		return "{}", nil
	}
	return capResult(string(resp.Data)), nil
}

func capResult(s string) string {
	if len(s) <= maxToolResultBytes {
		return s
	}
	return s[:maxToolResultBytes] + fmt.Sprintf("\n…(truncated, %d bytes total — narrow the selection set or paginate with per_page)", len(s))
}

// ── graphql_query (read) ─────────────────────────────────────────────────────

func (d Deps) graphqlQueryTool() *Tool {
	return &Tool{
		Name: "graphql_query",
		Description: "Run a READ-ONLY GraphQL query against stash's own API (the same schema its web UI uses). " +
			"Use graphql_schema first to discover the exact query name, arguments, and return fields. " +
			"Mutations are rejected here — use graphql_mutate for changes. Always select only the fields you need.",
		Schema: json.RawMessage(`{
			"type":"object",
			"properties":{
				"query":{"type":"string","description":"A GraphQL query document, e.g. 'query($f:FindFilterType){findTags(filter:$f){count tags{id name}}}'."},
				"variables":{"type":"object","description":"Optional variables object for the query."}
			},
			"required":["query"]
		}`),
		Writes: false,
		Run: func(ctx context.Context, input json.RawMessage) (string, error) {
			var in struct {
				Query     string         `json:"query"`
				Variables map[string]any `json:"variables"`
			}
			if err := json.Unmarshal(input, &in); err != nil {
				return "", fmt.Errorf("bad input: %w", err)
			}
			if strings.TrimSpace(in.Query) == "" {
				return "", fmt.Errorf("query is required")
			}
			kinds, err := d.operationKinds(in.Query)
			if err != nil {
				return "", err
			}
			if kinds[ast.Mutation] || kinds[ast.Subscription] {
				return "", fmt.Errorf("graphql_query is read-only; this document contains a mutation or subscription — use graphql_mutate instead")
			}
			return d.execGraphQL(ctx, in.Query, in.Variables)
		},
	}
}

// ── graphql_mutate (write, confirm-gated) ────────────────────────────────────

func (d Deps) graphqlMutateTool() *Tool {
	return &Tool{
		Name: "graphql_mutate",
		Description: "Run a GraphQL MUTATION against stash to change the library (e.g. tagsMerge, bulkSceneUpdate, " +
			"tagUpdate, performerMerge, sceneDestroy). Discover the exact mutation + input shape with graphql_schema " +
			"first, and resolve names to ids with graphql_query before mutating. This is a WRITE and requires user " +
			"confirmation; state plainly what it will change and how many items it affects before calling it.",
		Schema: json.RawMessage(`{
			"type":"object",
			"properties":{
				"query":{"type":"string","description":"A GraphQL mutation document, e.g. 'mutation($i:TagsMergeInput!){tagsMerge(input:$i){id name}}'."},
				"variables":{"type":"object","description":"Variables object for the mutation."}
			},
			"required":["query"]
		}`),
		Writes: true,
		Run: func(ctx context.Context, input json.RawMessage) (string, error) {
			var in struct {
				Query     string         `json:"query"`
				Variables map[string]any `json:"variables"`
			}
			if err := json.Unmarshal(input, &in); err != nil {
				return "", fmt.Errorf("bad input: %w", err)
			}
			if strings.TrimSpace(in.Query) == "" {
				return "", fmt.Errorf("query is required")
			}
			kinds, err := d.operationKinds(in.Query)
			if err != nil {
				return "", err
			}
			if !kinds[ast.Mutation] {
				return "", fmt.Errorf("graphql_mutate expects a mutation; this document has none — use graphql_query for reads")
			}
			return d.execGraphQL(ctx, in.Query, in.Variables)
		},
	}
}

// ── graphql_schema (read, introspection) ─────────────────────────────────────

func (d Deps) graphqlSchemaTool() *Tool {
	return &Tool{
		Name: "graphql_schema",
		Description: "Introspect stash's GraphQL schema so you can craft correct queries/mutations. " +
			"Call with {type:\"Name\"} to see a type's full definition (fields/args, or input-object fields, or enum values); " +
			"with {section:\"mutations\"|\"queries\"} to list root operation signatures; with {search:\"tag\"} to filter by name; " +
			"or with no args to list all query and mutation root field names.",
		Schema: json.RawMessage(`{
			"type":"object",
			"properties":{
				"type":{"type":"string","description":"A GraphQL type name to describe in full (e.g. 'Scene', 'TagsMergeInput', 'SceneFilterType')."},
				"section":{"type":"string","enum":["queries","mutations","subscriptions"],"description":"List root field signatures for this operation type."},
				"search":{"type":"string","description":"Case-insensitive substring filter on names."}
			}
		}`),
		Writes: false,
		Run: func(ctx context.Context, input json.RawMessage) (string, error) {
			var in struct {
				Type    string `json:"type"`
				Section string `json:"section"`
				Search  string `json:"search"`
			}
			if len(input) > 0 {
				_ = json.Unmarshal(input, &in)
			}
			search := strings.ToLower(strings.TrimSpace(in.Search))

			if t := strings.TrimSpace(in.Type); t != "" {
				def, ok := d.Schema.Types[t]
				if !ok {
					return "", fmt.Errorf("no such type %q (names are case-sensitive)", t)
				}
				return capResult(describeType(def)), nil
			}

			if sec := strings.TrimSpace(in.Section); sec != "" {
				var root *ast.Definition
				switch sec {
				case "queries":
					root = d.Schema.Query
				case "mutations":
					root = d.Schema.Mutation
				case "subscriptions":
					root = d.Schema.Subscription
				}
				if root == nil {
					return "", fmt.Errorf("no %s in schema", sec)
				}
				return capResult(listFieldSignatures(root, search)), nil
			}

			// Default: compact index of root query + mutation field names.
			var b strings.Builder
			b.WriteString("Query fields:\n")
			b.WriteString(listFieldNames(d.Schema.Query, search))
			b.WriteString("\nMutation fields:\n")
			b.WriteString(listFieldNames(d.Schema.Mutation, search))
			b.WriteString("\n(Call graphql_schema with section=\"mutations\" for signatures, or type=\"Name\" for a full definition.)")
			return capResult(b.String()), nil
		},
	}
}

func listFieldNames(def *ast.Definition, search string) string {
	if def == nil {
		return "  (none)\n"
	}
	names := make([]string, 0, len(def.Fields))
	for _, f := range def.Fields {
		if strings.HasPrefix(f.Name, "__") {
			continue
		}
		if search != "" && !strings.Contains(strings.ToLower(f.Name), search) {
			continue
		}
		names = append(names, f.Name)
	}
	sort.Strings(names)
	if len(names) == 0 {
		return "  (no matches)\n"
	}
	return "  " + strings.Join(names, ", ") + "\n"
}

func listFieldSignatures(def *ast.Definition, search string) string {
	var b strings.Builder
	sigs := make([]string, 0, len(def.Fields))
	for _, f := range def.Fields {
		if strings.HasPrefix(f.Name, "__") {
			continue
		}
		if search != "" && !strings.Contains(strings.ToLower(f.Name), search) {
			continue
		}
		sigs = append(sigs, fieldSignature(f))
	}
	sort.Strings(sigs)
	for _, s := range sigs {
		b.WriteString(s)
		b.WriteString("\n")
	}
	if b.Len() == 0 {
		return "(no matching fields)"
	}
	return b.String()
}

func fieldSignature(f *ast.FieldDefinition) string {
	var args []string
	for _, a := range f.Arguments {
		args = append(args, fmt.Sprintf("%s: %s", a.Name, a.Type.String()))
	}
	if len(args) > 0 {
		return fmt.Sprintf("%s(%s): %s", f.Name, strings.Join(args, ", "), f.Type.String())
	}
	return fmt.Sprintf("%s: %s", f.Name, f.Type.String())
}

func describeType(def *ast.Definition) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s %s", strings.ToLower(string(def.Kind)), def.Name)
	if len(def.Interfaces) > 0 {
		fmt.Fprintf(&b, " implements %s", strings.Join(def.Interfaces, " & "))
	}
	b.WriteString(" {\n")
	switch def.Kind {
	case ast.Object, ast.Interface, ast.InputObject:
		for _, f := range def.Fields {
			if strings.HasPrefix(f.Name, "__") {
				continue
			}
			b.WriteString("  ")
			b.WriteString(fieldSignature(f))
			b.WriteString("\n")
		}
	case ast.Enum:
		for _, v := range def.EnumValues {
			fmt.Fprintf(&b, "  %s\n", v.Name)
		}
	case ast.Union:
		fmt.Fprintf(&b, "  = %s\n", strings.Join(def.Types, " | "))
	case ast.Scalar:
		b.WriteString("  (scalar)\n")
	}
	b.WriteString("}")
	return b.String()
}
