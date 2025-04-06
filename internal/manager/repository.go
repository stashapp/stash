package manager

import (
	"context"

	"github.com/stashapp/stash/pkg/group"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
)

type SceneService interface {
	Create(ctx context.Context, input *models.Scene, fileIDs []models.FileID, coverImage []byte) (*models.Scene, error)
	AssignFile(ctx context.Context, sceneID int, fileID models.FileID) error
	Merge(ctx context.Context, sourceIDs []int, destinationID int, fileDeleter *scene.FileDeleter, options scene.MergeOptions) error
	Destroy(ctx context.Context, scene *models.Scene, fileDeleter *scene.FileDeleter, deleteGenerated, deleteFile bool) error

	FindByIDs(ctx context.Context, ids []int, load ...scene.LoadRelationshipOption) ([]*models.Scene, error)
	sceneFingerprintGetter
}

type ImageService interface {
	Destroy(ctx context.Context, image *models.Image, fileDeleter *image.FileDeleter, deleteGenerated, deleteFile bool) error
	DestroyZipImages(ctx context.Context, zipFile models.File, fileDeleter *image.FileDeleter, deleteGenerated bool) ([]*models.Image, error)
}

type GalleryService interface {
	AddImages(ctx context.Context, g *models.Gallery, toAdd ...int) error
	RemoveImages(ctx context.Context, g *models.Gallery, toRemove ...int) error

	SetCover(ctx context.Context, g *models.Gallery, coverImageId int) error
	ResetCover(ctx context.Context, g *models.Gallery) error

	Destroy(ctx context.Context, i *models.Gallery, fileDeleter *image.FileDeleter, deleteGenerated, deleteFile bool) ([]*models.Image, error)

	ValidateImageGalleryChange(ctx context.Context, i *models.Image, updateIDs models.UpdateIDs) error

	Updated(ctx context.Context, galleryID int) error
}

type GroupService interface {
	Create(ctx context.Context, group *models.Group, frontimageData []byte, backimageData []byte) error
	UpdatePartial(ctx context.Context, id int, updatedGroup models.GroupPartial, frontImage group.ImageInput, backImage group.ImageInput) (*models.Group, error)

	AddSubGroups(ctx context.Context, groupID int, subGroups []models.GroupIDDescription, insertIndex *int) error
	RemoveSubGroups(ctx context.Context, groupID int, subGroupIDs []int) error
	ReorderSubGroups(ctx context.Context, groupID int, subGroupIDs []int, insertPointID int, insertAfter bool) error
}
