package api

import (
	"context"

	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/models"
)

func (r *mutationResolver) MetadataScan(ctx context.Context, input models.ScanMetadataInput) (string, error) {
	manager.GetInstance().Scan(input.UseFileMetadata)
	return "todo", nil
}

func (r *mutationResolver) MetadataImport(ctx context.Context) (string, error) {
	manager.GetInstance().Import()
	return "todo", nil
}

func (r *mutationResolver) MetadataExport(ctx context.Context) (string, error) {
	manager.GetInstance().Export()
	return "todo", nil
}

func (r *mutationResolver) MetadataGenerate(ctx context.Context, input models.GenerateMetadataInput) (string, error) {
	manager.GetInstance().Generate(input)
	return "todo", nil
}

func (r *mutationResolver) MetadataAutoTag(ctx context.Context, input models.AutoTagMetadataInput) (string, error) {
	manager.GetInstance().AutoTag(input.Performers, input.Studios, input.Tags)
	return "todo", nil
}

func (r *mutationResolver) MetadataClean(ctx context.Context) (string, error) {
	manager.GetInstance().Clean()
	return "todo", nil
}

func (r *mutationResolver) MigrateHashNaming(ctx context.Context) (string, error) {
	manager.GetInstance().MigrateHash()
	return "todo", nil
}

func (r *mutationResolver) JobStatus(ctx context.Context) (*models.MetadataUpdateStatus, error) {
	status := manager.GetInstance().Status
	ret := models.MetadataUpdateStatus{
		Progress: status.Progress,
		Status:   status.Status.String(),
		Message:  "",
	}

	return &ret, nil
}

func (r *mutationResolver) StopJob(ctx context.Context) (bool, error) {
	return manager.GetInstance().Status.Stop(), nil
}
