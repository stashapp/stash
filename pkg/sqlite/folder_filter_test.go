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

func TestFolderQuery(t *testing.T) {
	tests := []struct {
		name        string
		findFilter  *models.FindFilterType
		filter      *models.FolderFilterType
		includeIdxs []int
		includeIDs  []models.FolderID
		excludeIdxs []int
		wantErr     bool
	}{
		{
			name: "path",
			filter: &models.FolderFilterType{
				Path: &models.StringCriterionInput{
					Value:    getFolderPath(folderIdxWithSubFolder, nil),
					Modifier: models.CriterionModifierIncludes,
				},
			},
			includeIdxs: []int{folderIdxWithSubFolder, folderIdxWithParentFolder},
			excludeIdxs: []int{folderIdxInZip},
		},
		{
			name: "parent folder",
			filter: &models.FolderFilterType{
				ParentFolder: &models.HierarchicalMultiCriterionInput{
					Value: []string{
						strconv.Itoa(int(folderIDs[folderIdxWithSubFolder])),
					},
					Modifier: models.CriterionModifierIncludes,
				},
			},
			includeIdxs: []int{folderIdxWithParentFolder},
			excludeIdxs: []int{folderIdxWithSubFolder, folderIdxInZip},
		},
		{
			name: "zip file",
			filter: &models.FolderFilterType{
				ZipFile: &models.MultiCriterionInput{
					Value: []string{
						strconv.Itoa(int(fileIDs[fileIdxZip])),
					},
					Modifier: models.CriterionModifierIncludes,
				},
			},
			includeIdxs: []int{folderIdxInZip},
			excludeIdxs: []int{folderIdxForObjectFiles},
		},
		// TODO - add more tests for other folder filters
	}

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			results, err := db.Folder.Query(ctx, models.FolderQueryOptions{
				FolderFilter: tt.filter,
				QueryOptions: models.QueryOptions{
					FindFilter: tt.findFilter,
				},
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("SceneStore.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			include := indexesToIDPtrs(folderIDs, tt.includeIdxs)
			for _, id := range tt.includeIDs {
				v := id
				include = append(include, &v)
			}
			exclude := indexesToIDPtrs(folderIDs, tt.excludeIdxs)

			for _, i := range include {
				assert.Contains(results.IDs, models.FolderID(*i))
			}
			for _, e := range exclude {
				assert.NotContains(results.IDs, models.FolderID(*e))
			}
		})
	}
}
