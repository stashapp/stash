package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/txn"
)

// Deps is the slice of stash's data layer the library tools need. It is wired from
// the repository in internal/api (getLLMRoutes).
type Deps struct {
	TxnManager models.TxnManager
	Scene      models.SceneReaderWriter
	Performer  models.PerformerReaderWriter
	Studio     models.StudioReaderWriter
	Tag        models.TagReaderWriter
}

// RegisterLibraryTools registers the Phase 1 read + write tools against deps.
func RegisterLibraryTools(reg *Registry, deps Deps) {
	reg.Register(deps.libraryStatsTool())
	reg.Register(deps.findScenesTool())
	reg.Register(deps.findPerformersTool())
	reg.Register(deps.findStudiosTool())
	reg.Register(deps.findTagsTool())
	reg.Register(deps.createTagTool())
	reg.Register(deps.addTagsToScenesTool())
	reg.Register(deps.setScenesOrganizedTool())
}

// ── helpers ──────────────────────────────────────────────────────────────────

func intsToStrings(ids []int) []string {
	out := make([]string, len(ids))
	for i, id := range ids {
		out[i] = strconv.Itoa(id)
	}
	return out
}

func dateStr(d *models.Date) string {
	if d == nil {
		return ""
	}
	return d.String()
}

func resultJSON(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func buildFindFilter(query string, page, perPage int) *models.FindFilterType {
	ff := &models.FindFilterType{}
	if query != "" {
		q := query
		ff.Q = &q
	}
	if page < 1 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 25
	}
	if perPage > 100 {
		perPage = 100
	}
	ff.Page = &page
	ff.PerPage = &perPage
	return ff
}

// ── read tools ───────────────────────────────────────────────────────────────

func (d Deps) libraryStatsTool() *Tool {
	return &Tool{
		Name:        "library_stats",
		Description: "Return overall counts for the media library: scenes, performers, studios, and tags.",
		Schema:      json.RawMessage(`{"type":"object","properties":{},"additionalProperties":false}`),
		Run: func(ctx context.Context, _ json.RawMessage) (string, error) {
			var scenes, performers, studios, tags int
			err := txn.WithReadTxn(ctx, d.TxnManager, func(ctx context.Context) error {
				var e error
				if scenes, e = d.Scene.Count(ctx); e != nil {
					return e
				}
				if performers, e = d.Performer.Count(ctx); e != nil {
					return e
				}
				if studios, e = d.Studio.Count(ctx); e != nil {
					return e
				}
				if tags, e = d.Tag.Count(ctx); e != nil {
					return e
				}
				return nil
			})
			if err != nil {
				return "", err
			}
			return resultJSON(map[string]int{
				"scenes": scenes, "performers": performers, "studios": studios, "tags": tags,
			})
		},
	}
}

type findScenesInput struct {
	Query      string   `json:"query"`
	Organized  *bool    `json:"organized"`
	Tags       []string `json:"tags"`
	Studio     string   `json:"studio"`
	Performers []string `json:"performers"`
	Page       int      `json:"page"`
	PerPage    int      `json:"per_page"`
}

type sceneOut struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Date      string `json:"date,omitempty"`
	Rating    *int   `json:"rating100,omitempty"`
	Organized bool   `json:"organized"`
	StudioID  *int   `json:"studio_id,omitempty"`
}

func (d Deps) findScenesTool() *Tool {
	return &Tool{
		Name: "find_scenes",
		Description: "Search scenes. Optional filters: free-text query (matches title/path/details), " +
			"organized flag, tag names (must have all), a studio name, and performer names (must have all). " +
			"Returns a compact list plus the total match count. Use this before any bulk action to identify scene ids.",
		Schema: json.RawMessage(`{
			"type":"object",
			"properties":{
				"query":{"type":"string","description":"free-text search"},
				"organized":{"type":"boolean"},
				"tags":{"type":"array","items":{"type":"string"},"description":"tag names; scene must have ALL"},
				"studio":{"type":"string","description":"studio name"},
				"performers":{"type":"array","items":{"type":"string"},"description":"performer names; scene must have ALL"},
				"page":{"type":"integer","minimum":1},
				"per_page":{"type":"integer","minimum":1,"maximum":100}
			},
			"additionalProperties":false
		}`),
		Run: func(ctx context.Context, raw json.RawMessage) (string, error) {
			var in findScenesInput
			if err := json.Unmarshal(raw, &in); err != nil {
				return "", fmt.Errorf("invalid input: %w", err)
			}
			ff := buildFindFilter(in.Query, in.Page, in.PerPage)

			var out []sceneOut
			var total int
			err := txn.WithReadTxn(ctx, d.TxnManager, func(ctx context.Context) error {
				sf := &models.SceneFilterType{}
				if in.Organized != nil {
					sf.Organized = in.Organized
				}
				if len(in.Tags) > 0 {
					tags, e := d.Tag.FindByNames(ctx, in.Tags, true)
					if e != nil {
						return e
					}
					ids := make([]int, 0, len(tags))
					for _, t := range tags {
						ids = append(ids, t.ID)
					}
					sf.Tags = &models.HierarchicalMultiCriterionInput{
						Value: intsToStrings(ids), Modifier: models.CriterionModifierIncludesAll,
					}
				}
				if in.Studio != "" {
					st, e := d.Studio.FindByName(ctx, in.Studio, true)
					if e != nil {
						return e
					}
					if st != nil {
						sf.Studios = &models.HierarchicalMultiCriterionInput{
							Value: []string{strconv.Itoa(st.ID)}, Modifier: models.CriterionModifierIncludes,
						}
					}
				}
				if len(in.Performers) > 0 {
					ps, e := d.Performer.FindByNames(ctx, in.Performers, true)
					if e != nil {
						return e
					}
					ids := make([]int, 0, len(ps))
					for _, p := range ps {
						ids = append(ids, p.ID)
					}
					sf.Performers = &models.MultiCriterionInput{
						Value: intsToStrings(ids), Modifier: models.CriterionModifierIncludesAll,
					}
				}

				res, e := d.Scene.Query(ctx, models.SceneQueryOptions{
					QueryOptions: models.QueryOptions{FindFilter: ff, Count: true},
					SceneFilter:  sf,
				})
				if e != nil {
					return e
				}
				total = res.Count
				scenes, e := res.Resolve(ctx)
				if e != nil {
					return e
				}
				for _, s := range scenes {
					out = append(out, sceneOut{
						ID: s.ID, Title: s.Title, Date: dateStr(s.Date),
						Rating: s.Rating, Organized: s.Organized, StudioID: s.StudioID,
					})
				}
				return nil
			})
			if err != nil {
				return "", err
			}
			return resultJSON(map[string]any{"count": total, "scenes": out})
		},
	}
}

type findByNameInput struct {
	Query   string `json:"query"`
	PerPage int    `json:"per_page"`
}

type namedOut struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (d Deps) findPerformersTool() *Tool {
	return &Tool{
		Name:        "find_performers",
		Description: "Search performers by name (free-text query). Returns id + name for each match and the total count.",
		Schema:      json.RawMessage(`{"type":"object","properties":{"query":{"type":"string"},"per_page":{"type":"integer","minimum":1,"maximum":100}},"additionalProperties":false}`),
		Run: func(ctx context.Context, raw json.RawMessage) (string, error) {
			var in findByNameInput
			if err := json.Unmarshal(raw, &in); err != nil {
				return "", fmt.Errorf("invalid input: %w", err)
			}
			ff := buildFindFilter(in.Query, 1, in.PerPage)
			var out []namedOut
			var total int
			err := txn.WithReadTxn(ctx, d.TxnManager, func(ctx context.Context) error {
				ps, t, e := d.Performer.Query(ctx, nil, ff)
				if e != nil {
					return e
				}
				total = t
				for _, p := range ps {
					out = append(out, namedOut{ID: p.ID, Name: p.Name})
				}
				return nil
			})
			if err != nil {
				return "", err
			}
			return resultJSON(map[string]any{"count": total, "performers": out})
		},
	}
}

func (d Deps) findStudiosTool() *Tool {
	return &Tool{
		Name:        "find_studios",
		Description: "Search studios by name (free-text query). Returns id + name for each match and the total count.",
		Schema:      json.RawMessage(`{"type":"object","properties":{"query":{"type":"string"},"per_page":{"type":"integer","minimum":1,"maximum":100}},"additionalProperties":false}`),
		Run: func(ctx context.Context, raw json.RawMessage) (string, error) {
			var in findByNameInput
			if err := json.Unmarshal(raw, &in); err != nil {
				return "", fmt.Errorf("invalid input: %w", err)
			}
			ff := buildFindFilter(in.Query, 1, in.PerPage)
			var out []namedOut
			var total int
			err := txn.WithReadTxn(ctx, d.TxnManager, func(ctx context.Context) error {
				ss, t, e := d.Studio.Query(ctx, nil, ff)
				if e != nil {
					return e
				}
				total = t
				for _, s := range ss {
					out = append(out, namedOut{ID: s.ID, Name: s.Name})
				}
				return nil
			})
			if err != nil {
				return "", err
			}
			return resultJSON(map[string]any{"count": total, "studios": out})
		},
	}
}

func (d Deps) findTagsTool() *Tool {
	return &Tool{
		Name:        "find_tags",
		Description: "Search tags by name (free-text query). Returns id + name for each match and the total count.",
		Schema:      json.RawMessage(`{"type":"object","properties":{"query":{"type":"string"},"per_page":{"type":"integer","minimum":1,"maximum":100}},"additionalProperties":false}`),
		Run: func(ctx context.Context, raw json.RawMessage) (string, error) {
			var in findByNameInput
			if err := json.Unmarshal(raw, &in); err != nil {
				return "", fmt.Errorf("invalid input: %w", err)
			}
			ff := buildFindFilter(in.Query, 1, in.PerPage)
			var out []namedOut
			var total int
			err := txn.WithReadTxn(ctx, d.TxnManager, func(ctx context.Context) error {
				ts, t, e := d.Tag.Query(ctx, nil, ff)
				if e != nil {
					return e
				}
				total = t
				for _, tg := range ts {
					out = append(out, namedOut{ID: tg.ID, Name: tg.Name})
				}
				return nil
			})
			if err != nil {
				return "", err
			}
			return resultJSON(map[string]any{"count": total, "tags": out})
		},
	}
}

// ── write tools ──────────────────────────────────────────────────────────────

type createTagInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (d Deps) createTagTool() *Tool {
	return &Tool{
		Name:        "create_tag",
		Description: "Create a new tag. Returns the new tag's id. Check find_tags first to avoid duplicates.",
		Writes:      true,
		Schema:      json.RawMessage(`{"type":"object","properties":{"name":{"type":"string"},"description":{"type":"string"}},"required":["name"],"additionalProperties":false}`),
		Run: func(ctx context.Context, raw json.RawMessage) (string, error) {
			var in createTagInput
			if err := json.Unmarshal(raw, &in); err != nil {
				return "", fmt.Errorf("invalid input: %w", err)
			}
			if in.Name == "" {
				return "", fmt.Errorf("name is required")
			}
			newTag := &models.Tag{Name: in.Name, Description: in.Description}
			err := txn.WithTxn(ctx, d.TxnManager, func(ctx context.Context) error {
				return d.Tag.Create(ctx, &models.CreateTagInput{Tag: newTag})
			})
			if err != nil {
				return "", err
			}
			return resultJSON(map[string]any{"id": newTag.ID, "name": newTag.Name})
		},
	}
}

type sceneTagsInput struct {
	SceneIDs []int `json:"scene_ids"`
	TagIDs   []int `json:"tag_ids"`
}

func (d Deps) addTagsToScenesTool() *Tool {
	return &Tool{
		Name:        "add_tags_to_scenes",
		Description: "Add the given tag ids to each of the given scene ids (additive; existing tags are kept). Resolve names to ids with find_tags / find_scenes first.",
		Writes:      true,
		Schema:      json.RawMessage(`{"type":"object","properties":{"scene_ids":{"type":"array","items":{"type":"integer"}},"tag_ids":{"type":"array","items":{"type":"integer"}}},"required":["scene_ids","tag_ids"],"additionalProperties":false}`),
		Run: func(ctx context.Context, raw json.RawMessage) (string, error) {
			var in sceneTagsInput
			if err := json.Unmarshal(raw, &in); err != nil {
				return "", fmt.Errorf("invalid input: %w", err)
			}
			if len(in.SceneIDs) == 0 || len(in.TagIDs) == 0 {
				return "", fmt.Errorf("scene_ids and tag_ids are both required")
			}
			updated := 0
			err := txn.WithTxn(ctx, d.TxnManager, func(ctx context.Context) error {
				for _, sceneID := range in.SceneIDs {
					_, e := d.Scene.UpdatePartial(ctx, sceneID, models.ScenePartial{
						TagIDs: &models.UpdateIDs{IDs: in.TagIDs, Mode: models.RelationshipUpdateModeAdd},
					})
					if e != nil {
						return fmt.Errorf("scene %d: %w", sceneID, e)
					}
					updated++
				}
				return nil
			})
			if err != nil {
				return "", err
			}
			return resultJSON(map[string]any{"updated_scenes": updated, "tag_ids": in.TagIDs})
		},
	}
}

type sceneOrganizedInput struct {
	SceneIDs  []int `json:"scene_ids"`
	Organized bool  `json:"organized"`
}

func (d Deps) setScenesOrganizedTool() *Tool {
	return &Tool{
		Name:        "set_scenes_organized",
		Description: "Set the 'organized' flag on each of the given scene ids.",
		Writes:      true,
		Schema:      json.RawMessage(`{"type":"object","properties":{"scene_ids":{"type":"array","items":{"type":"integer"}},"organized":{"type":"boolean"}},"required":["scene_ids","organized"],"additionalProperties":false}`),
		Run: func(ctx context.Context, raw json.RawMessage) (string, error) {
			var in sceneOrganizedInput
			if err := json.Unmarshal(raw, &in); err != nil {
				return "", fmt.Errorf("invalid input: %w", err)
			}
			if len(in.SceneIDs) == 0 {
				return "", fmt.Errorf("scene_ids is required")
			}
			updated := 0
			err := txn.WithTxn(ctx, d.TxnManager, func(ctx context.Context) error {
				for _, sceneID := range in.SceneIDs {
					_, e := d.Scene.UpdatePartial(ctx, sceneID, models.ScenePartial{
						Organized: models.NewOptionalBool(in.Organized),
					})
					if e != nil {
						return fmt.Errorf("scene %d: %w", sceneID, e)
					}
					updated++
				}
				return nil
			})
			if err != nil {
				return "", err
			}
			return resultJSON(map[string]any{"updated_scenes": updated, "organized": in.Organized})
		},
	}
}
