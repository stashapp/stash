//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stretchr/testify/assert"
)

var (
	invalidFolderID = models.FolderID(invalidID)
	invalidFileID   = models.FileID(invalidID)
)

func Test_FolderStore_Create(t *testing.T) {
	var (
		path        = "path"
		fileModTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		createdAt   = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt   = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
	)

	tests := []struct {
		name      string
		newObject models.Folder
		wantErr   bool
	}{
		{
			"full",
			models.Folder{
				DirEntry: models.DirEntry{
					ZipFileID: &fileIDs[fileIdxZip],
					ZipFile:   makeZipFileWithID(fileIdxZip),
					ModTime:   fileModTime,
				},
				Path:      path,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			false,
		},
		{
			"invalid parent folder id",
			models.Folder{
				Path:           path,
				ParentFolderID: &invalidFolderID,
			},
			true,
		},
		{
			"invalid zip file id",
			models.Folder{
				DirEntry: models.DirEntry{
					ZipFileID: &invalidFileID,
				},
				Path: path,
			},
			true,
		},
	}

	qb := db.Folder

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			s := tt.newObject
			if err := qb.Create(ctx, &s); (err != nil) != tt.wantErr {
				t.Errorf("FolderStore.Create() error = %v, wantErr = %v", err, tt.wantErr)
			}

			if tt.wantErr {
				assert.Zero(s.ID)
				return
			}

			assert.NotZero(s.ID)

			copy := tt.newObject
			copy.ID = s.ID

			assert.Equal(copy, s)

			// ensure can find the folder
			found, err := qb.FindByPath(ctx, path, true)
			if err != nil {
				t.Errorf("FolderStore.Find() error = %v", err)
			}

			assert.Equal(copy, *found)
		})
	}
}

func Test_FolderStore_Update(t *testing.T) {
	var (
		path        = "path"
		fileModTime = time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		createdAt   = time.Date(2001, 1, 2, 3, 4, 5, 0, time.UTC)
		updatedAt   = time.Date(2002, 1, 2, 3, 4, 5, 0, time.UTC)
	)

	tests := []struct {
		name          string
		updatedObject *models.Folder
		wantErr       bool
	}{
		{
			"full",
			&models.Folder{
				ID: folderIDs[folderIdxWithParentFolder],
				DirEntry: models.DirEntry{
					ZipFileID: &fileIDs[fileIdxZip],
					ZipFile:   makeZipFileWithID(fileIdxZip),
					ModTime:   fileModTime,
				},
				Path:      path,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			false,
		},
		{
			"clear zip",
			&models.Folder{
				ID:   folderIDs[folderIdxInZip],
				Path: path,
			},
			false,
		},
		{
			"clear folder",
			&models.Folder{
				ID:   folderIDs[folderIdxWithParentFolder],
				Path: path,
			},
			false,
		},
		{
			"invalid parent folder id",
			&models.Folder{
				ID:             folderIDs[folderIdxWithParentFolder],
				Path:           path,
				ParentFolderID: &invalidFolderID,
			},
			true,
		},
		{
			"invalid zip file id",
			&models.Folder{
				ID: folderIDs[folderIdxWithParentFolder],
				DirEntry: models.DirEntry{
					ZipFileID: &invalidFileID,
				},
				Path: path,
			},
			true,
		},
	}

	qb := db.Folder
	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			copy := *tt.updatedObject

			if err := qb.Update(ctx, tt.updatedObject); (err != nil) != tt.wantErr {
				t.Errorf("FolderStore.Update() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			s, err := qb.FindByPath(ctx, path, true)
			if err != nil {
				t.Errorf("FolderStore.Find() error = %v", err)
			}

			assert.Equal(copy, *s)

			return
		})
	}
}

func makeFolderWithID(index int) *models.Folder {
	ret := makeFolder(index)
	ret.ID = folderIDs[index]

	return &ret
}

func Test_FolderStore_FindByPath(t *testing.T) {
	getPath := func(index int) string {
		return folderPaths[index]
	}

	tests := []struct {
		name    string
		path    string
		want    *models.Folder
		wantErr bool
	}{
		{
			"valid",
			getPath(folderIdxWithFiles),
			makeFolderWithID(folderIdxWithFiles),
			false,
		},
		{
			"invalid",
			"invalid path",
			nil,
			false,
		},
	}

	qb := db.Folder

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			got, err := qb.FindByPath(ctx, tt.path, true)
			if (err != nil) != tt.wantErr {
				t.Errorf("FolderStore.FindByPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FolderStore.FindByPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
