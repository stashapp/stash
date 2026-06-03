// Package loaders contains the dataloaders used by the resolver in [api].
// They are generated with `make generate-dataloaders`.
// The dataloaders are used to batch requests to the database.

//go:generate go run github.com/vektah/dataloaden SceneLoader int *github.com/stashapp/stash/pkg/models.Scene
//go:generate go run github.com/vektah/dataloaden AudioLoader int *github.com/stashapp/stash/pkg/models.Audio
//go:generate go run github.com/vektah/dataloaden GalleryLoader int *github.com/stashapp/stash/pkg/models.Gallery
//go:generate go run github.com/vektah/dataloaden ImageLoader int *github.com/stashapp/stash/pkg/models.Image
//go:generate go run github.com/vektah/dataloaden PerformerLoader int *github.com/stashapp/stash/pkg/models.Performer
//go:generate go run github.com/vektah/dataloaden StudioLoader int *github.com/stashapp/stash/pkg/models.Studio
//go:generate go run github.com/vektah/dataloaden TagLoader int *github.com/stashapp/stash/pkg/models.Tag
//go:generate go run github.com/vektah/dataloaden GroupLoader int *github.com/stashapp/stash/pkg/models.Group
//go:generate go run github.com/vektah/dataloaden FileLoader github.com/stashapp/stash/pkg/models.FileID github.com/stashapp/stash/pkg/models.File
//go:generate go run github.com/vektah/dataloaden FolderLoader github.com/stashapp/stash/pkg/models.FolderID *github.com/stashapp/stash/pkg/models.Folder
//go:generate go run github.com/vektah/dataloaden FolderRelatedFolderIDsLoader github.com/stashapp/stash/pkg/models.FolderID []github.com/stashapp/stash/pkg/models.FolderID
//go:generate go run github.com/vektah/dataloaden RelatedFileIDsLoader int []github.com/stashapp/stash/pkg/models.FileID
//go:generate go run github.com/vektah/dataloaden FileIDsRelatedIDsLoader github.com/stashapp/stash/pkg/models.FileID []int
//go:generate go run github.com/vektah/dataloaden CustomFieldsLoader int github.com/stashapp/stash/pkg/models.CustomFieldMap
//go:generate go run github.com/vektah/dataloaden SceneOCountLoader int int
//go:generate go run github.com/vektah/dataloaden ScenePlayCountLoader int int
//go:generate go run github.com/vektah/dataloaden SceneOHistoryLoader int []time.Time
//go:generate go run github.com/vektah/dataloaden ScenePlayHistoryLoader int []time.Time
//go:generate go run github.com/vektah/dataloaden SceneLastPlayedLoader int *time.Time
//go:generate go run github.com/vektah/dataloaden AudioOCountLoader int int
//go:generate go run github.com/vektah/dataloaden AudioPlayCountLoader int int
//go:generate go run github.com/vektah/dataloaden AudioOHistoryLoader int []time.Time
//go:generate go run github.com/vektah/dataloaden AudioPlayHistoryLoader int []time.Time
//go:generate go run github.com/vektah/dataloaden AudioLastPlayedLoader int *time.Time
package loaders

import (
	"context"
	"net/http"
	"time"

	"github.com/stashapp/stash/pkg/models"
)

type contextKey struct{ name string }

var (
	loadersCtxKey = &contextKey{"loaders"}
)

const (
	wait     = 1 * time.Millisecond
	maxBatch = 100
)

type Loaders struct {
	SceneByID         *SceneLoader
	SceneIDsByFileID  *FileIDsRelatedIDsLoader
	SceneFiles        *RelatedFileIDsLoader
	ScenePlayCount    *ScenePlayCountLoader
	SceneOCount       *SceneOCountLoader
	ScenePlayHistory  *ScenePlayHistoryLoader
	SceneOHistory     *SceneOHistoryLoader
	SceneLastPlayed   *SceneLastPlayedLoader
	SceneCustomFields *CustomFieldsLoader

	AudioByID         *AudioLoader
	AudioIDsByFileID  *FileIDsRelatedIDsLoader
	AudioFiles        *RelatedFileIDsLoader
	AudioPlayCount    *AudioPlayCountLoader
	AudioOCount       *AudioOCountLoader
	AudioPlayHistory  *AudioPlayHistoryLoader
	AudioOHistory     *AudioOHistoryLoader
	AudioLastPlayed   *AudioLastPlayedLoader
	AudioCustomFields *CustomFieldsLoader

	ImageFiles   *RelatedFileIDsLoader
	GalleryFiles *RelatedFileIDsLoader

	GalleryByID         *GalleryLoader
	GalleryIDsByFileID  *FileIDsRelatedIDsLoader
	GalleryCustomFields *CustomFieldsLoader
	ImageByID           *ImageLoader
	ImageIDsByFileID    *FileIDsRelatedIDsLoader
	ImageCustomFields   *CustomFieldsLoader

	PerformerByID         *PerformerLoader
	PerformerCustomFields *CustomFieldsLoader

	StudioByID         *StudioLoader
	StudioCustomFields *CustomFieldsLoader

	TagByID         *TagLoader
	TagCustomFields *CustomFieldsLoader

	GroupByID         *GroupLoader
	GroupCustomFields *CustomFieldsLoader

	FileByID *FileLoader

	FolderByID            *FolderLoader
	FolderParentFolderIDs *FolderRelatedFolderIDsLoader
	FolderSubFolderIDs    *FolderRelatedFolderIDsLoader
}

type Middleware struct {
	Repository models.Repository
}

func (m Middleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ldrs := Loaders{
			SceneByID: &SceneLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchScenes(ctx),
			},
			SceneIDsByFileID: &FileIDsRelatedIDsLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchSceneIDsByFileID(ctx),
			},
			AudioByID: &AudioLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchAudios(ctx),
			},
			AudioIDsByFileID: &FileIDsRelatedIDsLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchAudioIDsByFileID(ctx),
			},
			GalleryByID: &GalleryLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchGalleries(ctx),
			},
			GalleryIDsByFileID: &FileIDsRelatedIDsLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchGalleryIDsByFileID(ctx),
			},
			GalleryCustomFields: &CustomFieldsLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchGalleryCustomFields(ctx),
			},
			ImageByID: &ImageLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchImages(ctx),
			},
			ImageIDsByFileID: &FileIDsRelatedIDsLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchImageIDsByFileID(ctx),
			},
			ImageCustomFields: &CustomFieldsLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchImageCustomFields(ctx),
			},
			PerformerByID: &PerformerLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchPerformers(ctx),
			},
			PerformerCustomFields: &CustomFieldsLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchPerformerCustomFields(ctx),
			},
			StudioCustomFields: &CustomFieldsLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchStudioCustomFields(ctx),
			},
			SceneCustomFields: &CustomFieldsLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchSceneCustomFields(ctx),
			},
			AudioCustomFields: &CustomFieldsLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchAudioCustomFields(ctx),
			},
			StudioByID: &StudioLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchStudios(ctx),
			},
			TagByID: &TagLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchTags(ctx),
			},
			TagCustomFields: &CustomFieldsLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchTagCustomFields(ctx),
			},
			GroupByID: &GroupLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchGroups(ctx),
			},
			GroupCustomFields: &CustomFieldsLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchGroupCustomFields(ctx),
			},
			FileByID: &FileLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchFiles(ctx),
			},
			FolderByID: &FolderLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchFolders(ctx),
			},
			FolderParentFolderIDs: &FolderRelatedFolderIDsLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchFoldersParentFolderIDs(ctx),
			},
			FolderSubFolderIDs: &FolderRelatedFolderIDsLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchFoldersSubFolderIDs(ctx),
			},
			SceneFiles: &RelatedFileIDsLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchScenesFileIDs(ctx),
			},
			AudioFiles: &RelatedFileIDsLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchAudiosFileIDs(ctx),
			},
			ImageFiles: &RelatedFileIDsLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchImagesFileIDs(ctx),
			},
			GalleryFiles: &RelatedFileIDsLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchGalleriesFileIDs(ctx),
			},
			ScenePlayCount: &ScenePlayCountLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchScenesPlayCount(ctx),
			},
			SceneOCount: &SceneOCountLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchScenesOCount(ctx),
			},
			ScenePlayHistory: &ScenePlayHistoryLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchScenesPlayHistory(ctx),
			},
			SceneLastPlayed: &SceneLastPlayedLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchScenesLastPlayed(ctx),
			},
			SceneOHistory: &SceneOHistoryLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchScenesOHistory(ctx),
			},
			// Audio
			AudioPlayCount: &AudioPlayCountLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchAudiosPlayCount(ctx),
			},
			AudioOCount: &AudioOCountLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchAudiosOCount(ctx),
			},
			AudioPlayHistory: &AudioPlayHistoryLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchAudiosPlayHistory(ctx),
			},
			AudioLastPlayed: &AudioLastPlayedLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchAudiosLastPlayed(ctx),
			},
			AudioOHistory: &AudioOHistoryLoader{
				wait:     wait,
				maxBatch: maxBatch,
				fetch:    m.fetchAudiosOHistory(ctx),
			},
		}

		newCtx := context.WithValue(r.Context(), loadersCtxKey, ldrs)
		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}

func From(ctx context.Context) Loaders {
	return ctx.Value(loadersCtxKey).(Loaders)
}

func toErrorSlice(err error) []error {
	if err != nil {
		return []error{err}
	}

	return nil
}

func (m Middleware) fetchScenes(ctx context.Context) func(keys []int) ([]*models.Scene, []error) {
	return func(keys []int) (ret []*models.Scene, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Scene.FindMany(ctx, keys)
			return err
		})
		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchSceneIDsByFileID(ctx context.Context) func(keys []models.FileID) ([][]int, []error) {
	return func(keys []models.FileID) (ret [][]int, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Scene.GetManyIDsByFileIDs(ctx, keys)
			return err
		})
		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchSceneCustomFields(ctx context.Context) func(keys []int) ([]models.CustomFieldMap, []error) {
	return func(keys []int) (ret []models.CustomFieldMap, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Scene.GetCustomFieldsBulk(ctx, keys)
			return err
		})

		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchAudios(ctx context.Context) func(keys []int) ([]*models.Audio, []error) {
	return func(keys []int) (ret []*models.Audio, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Audio.FindMany(ctx, keys)
			return err
		})
		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchAudioIDsByFileID(ctx context.Context) func(keys []models.FileID) ([][]int, []error) {
	return func(keys []models.FileID) (ret [][]int, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Audio.GetManyIDsByFileIDs(ctx, keys)
			return err
		})
		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchAudioCustomFields(ctx context.Context) func(keys []int) ([]models.CustomFieldMap, []error) {
	return func(keys []int) (ret []models.CustomFieldMap, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Audio.GetCustomFieldsBulk(ctx, keys)
			return err
		})

		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchImages(ctx context.Context) func(keys []int) ([]*models.Image, []error) {
	return func(keys []int) (ret []*models.Image, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Image.FindMany(ctx, keys)
			return err
		})

		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchImageIDsByFileID(ctx context.Context) func(keys []models.FileID) ([][]int, []error) {
	return func(keys []models.FileID) (ret [][]int, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Image.GetManyIDsByFileIDs(ctx, keys)
			return err
		})
		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchImageCustomFields(ctx context.Context) func(keys []int) ([]models.CustomFieldMap, []error) {
	return func(keys []int) (ret []models.CustomFieldMap, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Image.GetCustomFieldsBulk(ctx, keys)
			return err
		})

		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchGalleries(ctx context.Context) func(keys []int) ([]*models.Gallery, []error) {
	return func(keys []int) (ret []*models.Gallery, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Gallery.FindMany(ctx, keys)
			return err
		})

		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchGalleryIDsByFileID(ctx context.Context) func(keys []models.FileID) ([][]int, []error) {
	return func(keys []models.FileID) (ret [][]int, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Gallery.GetManyIDsByFileIDs(ctx, keys)
			return err
		})
		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchPerformers(ctx context.Context) func(keys []int) ([]*models.Performer, []error) {
	return func(keys []int) (ret []*models.Performer, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Performer.FindMany(ctx, keys)
			return err
		})

		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchPerformerCustomFields(ctx context.Context) func(keys []int) ([]models.CustomFieldMap, []error) {
	return func(keys []int) (ret []models.CustomFieldMap, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Performer.GetCustomFieldsBulk(ctx, keys)
			return err
		})

		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchStudios(ctx context.Context) func(keys []int) ([]*models.Studio, []error) {
	return func(keys []int) (ret []*models.Studio, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Studio.FindMany(ctx, keys)
			return err
		})
		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchStudioCustomFields(ctx context.Context) func(keys []int) ([]models.CustomFieldMap, []error) {
	return func(keys []int) (ret []models.CustomFieldMap, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Studio.GetCustomFieldsBulk(ctx, keys)
			return err
		})

		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchTags(ctx context.Context) func(keys []int) ([]*models.Tag, []error) {
	return func(keys []int) (ret []*models.Tag, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Tag.FindMany(ctx, keys)
			return err
		})
		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchTagCustomFields(ctx context.Context) func(keys []int) ([]models.CustomFieldMap, []error) {
	return func(keys []int) (ret []models.CustomFieldMap, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Tag.GetCustomFieldsBulk(ctx, keys)
			return err
		})

		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchGroupCustomFields(ctx context.Context) func(keys []int) ([]models.CustomFieldMap, []error) {
	return func(keys []int) (ret []models.CustomFieldMap, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Group.GetCustomFieldsBulk(ctx, keys)
			return err
		})

		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchGalleryCustomFields(ctx context.Context) func(keys []int) ([]models.CustomFieldMap, []error) {
	return func(keys []int) (ret []models.CustomFieldMap, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Gallery.GetCustomFieldsBulk(ctx, keys)
			return err
		})

		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchGroups(ctx context.Context) func(keys []int) ([]*models.Group, []error) {
	return func(keys []int) (ret []*models.Group, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Group.FindMany(ctx, keys)
			return err
		})
		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchFiles(ctx context.Context) func(keys []models.FileID) ([]models.File, []error) {
	return func(keys []models.FileID) (ret []models.File, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.File.Find(ctx, keys...)
			return err
		})
		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchFolders(ctx context.Context) func(keys []models.FolderID) ([]*models.Folder, []error) {
	return func(keys []models.FolderID) (ret []*models.Folder, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Folder.FindMany(ctx, keys)
			return err
		})
		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchFoldersParentFolderIDs(ctx context.Context) func(keys []models.FolderID) ([][]models.FolderID, []error) {
	return func(keys []models.FolderID) (ret [][]models.FolderID, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Folder.GetManyParentFolderIDs(ctx, keys)
			return err
		})
		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchFoldersSubFolderIDs(ctx context.Context) func(keys []models.FolderID) ([][]models.FolderID, []error) {
	return func(keys []models.FolderID) (ret [][]models.FolderID, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Folder.GetManySubFolderIDs(ctx, keys)
			return err
		})
		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchScenesFileIDs(ctx context.Context) func(keys []int) ([][]models.FileID, []error) {
	return func(keys []int) (ret [][]models.FileID, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Scene.GetManyFileIDs(ctx, keys)
			return err
		})
		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchAudiosFileIDs(ctx context.Context) func(keys []int) ([][]models.FileID, []error) {
	return func(keys []int) (ret [][]models.FileID, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Audio.GetManyFileIDs(ctx, keys)
			return err
		})
		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchImagesFileIDs(ctx context.Context) func(keys []int) ([][]models.FileID, []error) {
	return func(keys []int) (ret [][]models.FileID, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Image.GetManyFileIDs(ctx, keys)
			return err
		})
		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchGalleriesFileIDs(ctx context.Context) func(keys []int) ([][]models.FileID, []error) {
	return func(keys []int) (ret [][]models.FileID, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Gallery.GetManyFileIDs(ctx, keys)
			return err
		})
		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchScenesOCount(ctx context.Context) func(keys []int) ([]int, []error) {
	return func(keys []int) (ret []int, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Scene.GetManyOCount(ctx, keys)
			return err
		})
		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchScenesPlayCount(ctx context.Context) func(keys []int) ([]int, []error) {
	return func(keys []int) (ret []int, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Scene.GetManyViewCount(ctx, keys)
			return err
		})
		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchScenesOHistory(ctx context.Context) func(keys []int) ([][]time.Time, []error) {
	return func(keys []int) (ret [][]time.Time, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Scene.GetManyODates(ctx, keys)
			return err
		})
		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchScenesPlayHistory(ctx context.Context) func(keys []int) ([][]time.Time, []error) {
	return func(keys []int) (ret [][]time.Time, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Scene.GetManyViewDates(ctx, keys)
			return err
		})
		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchScenesLastPlayed(ctx context.Context) func(keys []int) ([]*time.Time, []error) {
	return func(keys []int) (ret []*time.Time, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Scene.GetManyLastViewed(ctx, keys)
			return err
		})
		return ret, toErrorSlice(err)
	}
}

// Audio
func (m Middleware) fetchAudiosOCount(ctx context.Context) func(keys []int) ([]int, []error) {
	return func(keys []int) (ret []int, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Audio.GetManyOCount(ctx, keys)
			return err
		})
		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchAudiosPlayCount(ctx context.Context) func(keys []int) ([]int, []error) {
	return func(keys []int) (ret []int, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Audio.GetManyViewCount(ctx, keys)
			return err
		})
		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchAudiosOHistory(ctx context.Context) func(keys []int) ([][]time.Time, []error) {
	return func(keys []int) (ret [][]time.Time, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Audio.GetManyODates(ctx, keys)
			return err
		})
		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchAudiosPlayHistory(ctx context.Context) func(keys []int) ([][]time.Time, []error) {
	return func(keys []int) (ret [][]time.Time, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Audio.GetManyViewDates(ctx, keys)
			return err
		})
		return ret, toErrorSlice(err)
	}
}

func (m Middleware) fetchAudiosLastPlayed(ctx context.Context) func(keys []int) ([]*time.Time, []error) {
	return func(keys []int) (ret []*time.Time, errs []error) {
		err := m.Repository.WithDB(ctx, func(ctx context.Context) error {
			var err error
			ret, err = m.Repository.Audio.GetManyLastViewed(ctx, keys)
			return err
		})
		return ret, toErrorSlice(err)
	}
}
