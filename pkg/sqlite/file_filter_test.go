//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestFileQuery(t *testing.T) {
	tests := []struct {
		name        string
		findFilter  *models.FindFilterType
		filter      *models.FileFilterType
		includeIdxs []int
		includeIDs  []models.FileID
		excludeIdxs []int
		wantErr     bool
	}{
		{
			name: "path",
			filter: &models.FileFilterType{
				Path: &models.StringCriterionInput{
					Value:    getPrefixedStringValue("file", fileIdxStartVideoFiles, "basename"),
					Modifier: models.CriterionModifierIncludes,
				},
			},
			includeIdxs: []int{fileIdxStartVideoFiles},
			excludeIdxs: []int{fileIdxStartImageFiles},
		},
		{
			name: "basename",
			filter: &models.FileFilterType{
				Basename: &models.StringCriterionInput{
					Value:    getPrefixedStringValue("file", fileIdxStartVideoFiles, "basename"),
					Modifier: models.CriterionModifierIncludes,
				},
			},
			includeIdxs: []int{fileIdxStartVideoFiles},
			excludeIdxs: []int{fileIdxStartImageFiles},
		},
		{
			name: "dir",
			filter: &models.FileFilterType{
				Path: &models.StringCriterionInput{
					Value:    folderPaths[folderIdxWithSceneFiles],
					Modifier: models.CriterionModifierIncludes,
				},
			},
			includeIDs:  []models.FileID{sceneFileIDs[sceneIdxWithGroup]},
			excludeIdxs: []int{fileIdxStartImageFiles},
		},
		{
			name: "parent folder",
			filter: &models.FileFilterType{
				ParentFolder: &models.HierarchicalMultiCriterionInput{
					Value: []string{
						strconv.Itoa(int(folderIDs[folderIdxWithSceneFiles])),
					},
					Modifier: models.CriterionModifierIncludes,
				},
			},
			includeIDs:  []models.FileID{sceneFileIDs[sceneIdxWithGroup]},
			excludeIdxs: []int{fileIdxStartImageFiles},
		},
		{
			name: "zip file",
			filter: &models.FileFilterType{
				ZipFile: &models.MultiCriterionInput{
					Value: []string{
						strconv.Itoa(int(fileIDs[fileIdxZip])),
					},
					Modifier: models.CriterionModifierIncludes,
				},
			},
			includeIDs:  []models.FileID{fileIDs[fileIdxInZip]},
			excludeIdxs: []int{fileIdxStartImageFiles},
		},
		// TODO - add more tests for other file filters
	}

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			results, err := db.File.Query(ctx, models.FileQueryOptions{
				FileFilter: tt.filter,
				QueryOptions: models.QueryOptions{
					FindFilter: tt.findFilter,
				},
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("SceneStore.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			include := indexesToIDPtrs(fileIDs, tt.includeIdxs)
			for _, id := range tt.includeIDs {
				v := id
				include = append(include, &v)
			}
			exclude := indexesToIDPtrs(fileIDs, tt.excludeIdxs)

			for _, i := range include {
				assert.Contains(results.IDs, models.FileID(*i))
			}
			for _, e := range exclude {
				assert.NotContains(results.IDs, models.FileID(*e))
			}
		})
	}
}
