package api

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

func (r *savedFilterResolver) Filter(ctx context.Context, obj *models.SavedFilter) (string, error) {
	return "", nil
}

func (r *savedFilterResolver) LabelMapping(ctx context.Context, obj *models.SavedFilter) (*LabelMappingType, error) {
	mapping := &LabelMappingType{}
	if obj.ObjectFilter == nil {
		return mapping, nil
	}

	// Helper to extract IDs from a list of strings
	extractIDs := func(v interface{}) []int {
		var ids []int
		if list, ok := v.([]interface{}); ok {
			for _, item := range list {
				if strID, ok := item.(string); ok {
					if intID, err := strconv.Atoi(strID); err == nil {
						ids = append(ids, intID)
					}
				}
			}
		}
		return ids
	}

	// Helper to fetch and populate mapping
	populateMapping := func(criteriaKeys []string, fetchLabels func([]int) []*LabelMappingEntry) []*LabelMappingEntry {
		var allIDs []int

		for _, criteriaKey := range criteriaKeys {
			criterion, ok := obj.ObjectFilter[criteriaKey].(map[string]interface{})
			if !ok {
				continue
			}

			if val, ok := criterion["value"]; ok {
				allIDs = append(allIDs, extractIDs(val)...)
			}
			if excl, ok := criterion["excludes"]; ok {
				allIDs = append(allIDs, extractIDs(excl)...)
			}
		}

		if len(allIDs) > 0 {
			// deduplicate IDs
			idMap := make(map[int]bool)
			var dedupedIDs []int
			for _, id := range allIDs {
				if !idMap[id] {
					idMap[id] = true
					dedupedIDs = append(dedupedIDs, id)
				}
			}

			return fetchLabels(dedupedIDs)
		}

		return nil
	}

	err := r.withReadTxn(ctx, func(ctx context.Context) error {
		// Tags
		mapping.Tags = populateMapping([]string{"tags", "scene_tags", "performer_tags", "studio_tags", "parents", "children"}, func(ids []int) []*LabelMappingEntry {
			var res []*LabelMappingEntry
			tags, _ := r.repository.Tag.FindMany(ctx, ids)
			for _, t := range tags {
				res = append(res, &LabelMappingEntry{ID: strconv.Itoa(t.ID), Label: t.Name})
			}
			return res
		})

		// Performers
		mapping.Performers = populateMapping([]string{"performers"}, func(ids []int) []*LabelMappingEntry {
			var res []*LabelMappingEntry
			performers, _ := r.repository.Performer.FindMany(ctx, ids)
			for _, p := range performers {
				res = append(res, &LabelMappingEntry{ID: strconv.Itoa(p.ID), Label: p.Name})
			}
			return res
		})

		// Studios
		mapping.Studios = populateMapping([]string{"studios"}, func(ids []int) []*LabelMappingEntry {
			var res []*LabelMappingEntry
			studios, _ := r.repository.Studio.FindMany(ctx, ids)
			for _, s := range studios {
				res = append(res, &LabelMappingEntry{ID: strconv.Itoa(s.ID), Label: s.Name})
			}
			return res
		})

		// Groups
		mapping.Groups = populateMapping([]string{"groups", "containing_groups", "sub_groups"}, func(ids []int) []*LabelMappingEntry {
			var res []*LabelMappingEntry
			groups, _ := r.repository.Group.FindMany(ctx, ids)
			for _, g := range groups {
				res = append(res, &LabelMappingEntry{ID: strconv.Itoa(g.ID), Label: g.Name})
			}
			return res
		})

		// Galleries
		mapping.Galleries = populateMapping([]string{"galleries"}, func(ids []int) []*LabelMappingEntry {
			var res []*LabelMappingEntry
			galleries, _ := r.repository.Gallery.FindMany(ctx, ids)
			for _, g := range galleries {
				res = append(res, &LabelMappingEntry{ID: strconv.Itoa(g.ID), Label: g.Title})
			}
			return res
		})

		// Folders
		mapping.Folders = populateMapping([]string{"folders", "parent_folder"}, func(ids []int) []*LabelMappingEntry {
			var res []*LabelMappingEntry
			folderIDs := make([]models.FolderID, len(ids))
			for i, id := range ids {
				folderIDs[i] = models.FolderID(id)
			}
			folders, _ := r.repository.Folder.FindMany(ctx, folderIDs)
			for _, f := range folders {
				res = append(res, &LabelMappingEntry{ID: strconv.Itoa(int(f.ID)), Label: f.Path})
			}
			return res
		})

		// Scenes
		mapping.Scenes = populateMapping([]string{"scenes"}, func(ids []int) []*LabelMappingEntry {
			var res []*LabelMappingEntry
			scenes, _ := r.repository.Scene.FindMany(ctx, ids)
			for _, s := range scenes {
				label := s.Title
				if label == "" && s.Details != "" {
					label = s.Details
				}
				if label == "" {
					label = s.Checksum
				}
				res = append(res, &LabelMappingEntry{ID: strconv.Itoa(s.ID), Label: label})
			}
			return res
		})

		// Movies
		mapping.Movies = populateMapping([]string{"movies"}, func(ids []int) []*LabelMappingEntry {
			return nil
		})

		return nil
	})

	return mapping, err
}
