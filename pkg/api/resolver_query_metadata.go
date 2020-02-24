package api

import (
	"context"

	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) MetadataScan(ctx context.Context, input models.ScanMetadataInput) (string, error) {
	manager.GetInstance().Scan(input.UseFileMetadata)
	return "todo", nil
}

func (r *queryResolver) MetadataImport(ctx context.Context) (string, error) {
	manager.GetInstance().Import()
	return "todo", nil
}

func (r *queryResolver) MetadataExport(ctx context.Context) (string, error) {
	manager.GetInstance().Export()
	return "todo", nil
}

func (r *queryResolver) MetadataGenerate(ctx context.Context, input models.GenerateMetadataInput) (string, error) {
	manager.GetInstance().Generate(input.Sprites, input.Previews, input.Markers, input.Transcodes)
	return "todo", nil
}

func (r *queryResolver) MetadataAutoTag(ctx context.Context, input models.AutoTagMetadataInput) (string, error) {
	manager.GetInstance().AutoTag(input.Performers, input.Studios, input.Tags)
	return "todo", nil
}

func (r *queryResolver) MetadataClean(ctx context.Context) (string, error) {
	manager.GetInstance().Clean()
	return "todo", nil
}

func (r *queryResolver) JobStatus(ctx context.Context) (*models.MetadataUpdateStatus, error) {
	status := manager.GetInstance().Status
	ret := models.MetadataUpdateStatus{
		Progress: status.Progress,
		Status:   status.Status.String(),
		Message:  "",
	}

	return &ret, nil
}

func (r *queryResolver) StopJob(ctx context.Context) (bool, error) {
	return manager.GetInstance().Status.Stop(), nil
}
