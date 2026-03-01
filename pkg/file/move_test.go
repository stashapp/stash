package file

import (
	"context"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCorrectSubFolderHierarchy(t *testing.T) {
	ctx := context.Background()

	parentID := models.FolderID(1)

	tests := []struct {
		name          string
		setup         func(rw *mocks.FolderReaderWriter)
		parent        *models.Folder
		oldParentPath string
		wantErr       bool
	}{
		{
			name: "updates child folder path",
			parent: &models.Folder{
				ID:   parentID,
				Path: "/new/parent",
			},
			oldParentPath: "/old/parent",
			setup: func(rw *mocks.FolderReaderWriter) {
				child := &models.Folder{
					ID:   2,
					Path: "/old/parent/child",
				}
				rw.On("FindByParentFolderID", mock.Anything, parentID).Return([]*models.Folder{child}, nil)
				rw.On("FindByPath", mock.Anything, "/new/parent/child", true).Return(nil, nil)
				rw.On("Update", mock.Anything, mock.MatchedBy(func(f *models.Folder) bool {
					return f.ID == 2 && f.Path == "/new/parent/child"
				})).Return(nil)
				rw.On("FindByParentFolderID", mock.Anything, models.FolderID(2)).Return([]*models.Folder{}, nil)
			},
		},
		{
			name: "skips when path already correct",
			parent: &models.Folder{
				ID:   parentID,
				Path: "/content",
			},
			oldParentPath: "/content",
			setup: func(rw *mocks.FolderReaderWriter) {
				child := &models.Folder{
					ID:   2,
					Path: "/content/child",
				}
				rw.On("FindByParentFolderID", mock.Anything, parentID).Return([]*models.Folder{child}, nil)
				rw.On("FindByParentFolderID", mock.Anything, models.FolderID(2)).Return([]*models.Folder{}, nil)
			},
		},
		{
			name: "skips mis-parented child",
			parent: &models.Folder{
				ID:   parentID,
				Path: "/new/parent",
			},
			oldParentPath: "/old/parent",
			setup: func(rw *mocks.FolderReaderWriter) {
				misparented := &models.Folder{
					ID:   2,
					Path: "/old/parent/sorted/yes",
				}
				rw.On("FindByParentFolderID", mock.Anything, parentID).Return([]*models.Folder{misparented}, nil)
			},
		},
		{
			name: "skips when target path already exists",
			parent: &models.Folder{
				ID:   parentID,
				Path: "/new/parent",
			},
			oldParentPath: "/old/parent",
			setup: func(rw *mocks.FolderReaderWriter) {
				child := &models.Folder{
					ID:   2,
					Path: "/old/parent/child",
				}
				preExisting := &models.Folder{
					ID:   5,
					Path: "/new/parent/child",
				}
				rw.On("FindByParentFolderID", mock.Anything, parentID).Return([]*models.Folder{child}, nil)
				rw.On("FindByPath", mock.Anything, "/new/parent/child", true).Return(preExisting, nil)
			},
		},
		{
			name: "recurses into sub-folders",
			parent: &models.Folder{
				ID:   parentID,
				Path: "/new/root",
			},
			oldParentPath: "/old/root",
			setup: func(rw *mocks.FolderReaderWriter) {
				child := &models.Folder{
					ID:   2,
					Path: "/old/root/sorted",
				}
				grandchild := &models.Folder{
					ID:   3,
					Path: "/old/root/sorted/deep",
				}
				rw.On("FindByParentFolderID", mock.Anything, parentID).Return([]*models.Folder{child}, nil)
				rw.On("FindByPath", mock.Anything, "/new/root/sorted", true).Return(nil, nil)
				rw.On("Update", mock.Anything, mock.MatchedBy(func(f *models.Folder) bool {
					return f.ID == 2 && f.Path == "/new/root/sorted"
				})).Return(nil)

				rw.On("FindByParentFolderID", mock.Anything, models.FolderID(2)).Return([]*models.Folder{grandchild}, nil)
				rw.On("FindByPath", mock.Anything, "/new/root/sorted/deep", true).Return(nil, nil)
				rw.On("Update", mock.Anything, mock.MatchedBy(func(f *models.Folder) bool {
					return f.ID == 3 && f.Path == "/new/root/sorted/deep"
				})).Return(nil)

				rw.On("FindByParentFolderID", mock.Anything, models.FolderID(3)).Return([]*models.Folder{}, nil)
			},
		},
		{
			name: "issue 6427: duplicate basename at different hierarchy levels",
			parent: &models.Folder{
				ID:   parentID,
				Path: "/new/content",
			},
			oldParentPath: "/old/content",
			setup: func(rw *mocks.FolderReaderWriter) {
				childYes := &models.Folder{
					ID:   2,
					Path: "/old/content/yes",
				}
				childSorted := &models.Folder{
					ID:   3,
					Path: "/old/content/sorted",
				}
				grandchildYes := &models.Folder{
					ID:   4,
					Path: "/old/content/sorted/yes",
				}

				rw.On("FindByParentFolderID", mock.Anything, parentID).Return([]*models.Folder{childYes, childSorted}, nil)

				rw.On("FindByPath", mock.Anything, "/new/content/yes", true).Return(nil, nil)
				rw.On("Update", mock.Anything, mock.MatchedBy(func(f *models.Folder) bool {
					return f.ID == 2 && f.Path == "/new/content/yes"
				})).Return(nil)
				rw.On("FindByParentFolderID", mock.Anything, models.FolderID(2)).Return([]*models.Folder{}, nil)

				rw.On("FindByPath", mock.Anything, "/new/content/sorted", true).Return(nil, nil)
				rw.On("Update", mock.Anything, mock.MatchedBy(func(f *models.Folder) bool {
					return f.ID == 3 && f.Path == "/new/content/sorted"
				})).Return(nil)

				rw.On("FindByParentFolderID", mock.Anything, models.FolderID(3)).Return([]*models.Folder{grandchildYes}, nil)
				rw.On("FindByPath", mock.Anything, "/new/content/sorted/yes", true).Return(nil, nil)
				rw.On("Update", mock.Anything, mock.MatchedBy(func(f *models.Folder) bool {
					return f.ID == 4 && f.Path == "/new/content/sorted/yes"
				})).Return(nil)
				rw.On("FindByParentFolderID", mock.Anything, models.FolderID(4)).Return([]*models.Folder{}, nil)
			},
		},
		{
			name: "issue 6427: conflict from incorrect parent_folder_id",
			parent: &models.Folder{
				ID:   parentID,
				Path: "/new/content",
			},
			oldParentPath: "/old/content",
			setup: func(rw *mocks.FolderReaderWriter) {
				childYes := &models.Folder{
					ID:   2,
					Path: "/old/content/yes",
				}
				misparentedYes := &models.Folder{
					ID:   4,
					Path: "/old/content/sorted/yes",
				}

				rw.On("FindByParentFolderID", mock.Anything, parentID).Return([]*models.Folder{childYes, misparentedYes}, nil)

				rw.On("FindByPath", mock.Anything, "/new/content/yes", true).Return(nil, nil)
				rw.On("Update", mock.Anything, mock.MatchedBy(func(f *models.Folder) bool {
					return f.ID == 2 && f.Path == "/new/content/yes"
				})).Return(nil)
				rw.On("FindByParentFolderID", mock.Anything, models.FolderID(2)).Return([]*models.Folder{}, nil)

				// misparentedYes is skipped — filepath.Dir("/old/content/sorted/yes")
				// is "/old/content/sorted", not "/old/content"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := &mocks.FolderReaderWriter{}
			tt.setup(rw)

			err := correctSubFolderHierarchy(ctx, rw, tt.parent, tt.oldParentPath)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			rw.AssertExpectations(t)
		})
	}
}
