package manager

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
)

// Each generate option must queue its task and count it under its own total (#7341).
// Overwrite is on and the scene has no hash, so no existence checks touch the disk.
func TestGenerateJobQueuesTaskPerOption(t *testing.T) {
	config.InitializeEmpty()

	const sceneID = 1
	videoFile := &models.VideoFile{BaseFile: &models.BaseFile{ID: 1}, Width: 1920, Height: 1080, Interactive: true}
	scene := &models.Scene{
		ID:    sceneID,
		Path:  "scene.mp4",
		Files: models.NewRelatedVideoFiles([]*models.VideoFile{videoFile}),
	}
	image := &models.Image{
		Files: models.NewRelatedFiles([]models.File{
			videoFile, // video primary enables clip previews and thumbnails
			&models.ImageFile{BaseFile: &models.BaseFile{ID: 2}, Width: 1920, Height: 1080},
		}),
	}

	tests := []struct {
		name  string
		input GenerateMetadataInput
		want  totalsGenerate
	}{
		{"covers", GenerateMetadataInput{Covers: true}, totalsGenerate{covers: 1, tasks: 1}},
		{"sprites", GenerateMetadataInput{Sprites: true}, totalsGenerate{sprites: 1, tasks: 1}},
		{"previews", GenerateMetadataInput{Previews: true}, totalsGenerate{previews: 1, tasks: 1}},
		{"previews with image previews", GenerateMetadataInput{Previews: true, ImagePreviews: true}, totalsGenerate{previews: 1, imagePreviews: 1, tasks: 1}},
		{"markers", GenerateMetadataInput{Markers: true}, totalsGenerate{markers: 1, tasks: 1}},
		{"marker image previews", GenerateMetadataInput{MarkerImagePreviews: true}, totalsGenerate{markers: 1, tasks: 1}},
		{"marker screenshots", GenerateMetadataInput{MarkerScreenshots: true}, totalsGenerate{markers: 1, tasks: 1}},
		{"transcodes", GenerateMetadataInput{Transcodes: true, ForceTranscodes: true}, totalsGenerate{transcodes: 1, tasks: 1}},
		{"phashes", GenerateMetadataInput{Phashes: true}, totalsGenerate{phashes: 1, tasks: 1}},
		{"interactive heatmaps", GenerateMetadataInput{InteractiveHeatmapsSpeeds: true}, totalsGenerate{interactiveHeatmapSpeeds: 1, tasks: 1}},
		{"image thumbnails", GenerateMetadataInput{ImageThumbnails: true}, totalsGenerate{imageThumbnails: 1, tasks: 1}},
		{"clip previews", GenerateMetadataInput{ClipPreviews: true}, totalsGenerate{clipPreviews: 1, tasks: 1}},
		{"image phashes", GenerateMetadataInput{ImagePhashes: true}, totalsGenerate{imagePhashes: 1, tasks: 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := mocks.NewDatabase()
			db.SceneMarker.On("FindBySceneID", mock.Anything, sceneID).Return([]*models.SceneMarker{{ID: 1, SceneID: sceneID}}, nil)

			j := &GenerateJob{repository: db.Repository(), input: tt.input, overwrite: true}
			queue := make(chan Task, 10)

			j.queueSceneJobs(context.Background(), nil, scene, queue)
			j.queueImageJob(nil, image, queue)

			assert.Equal(t, tt.want, j.totals)
			assert.Len(t, queue, tt.want.tasks)

			// the marker count ignores these flags, so check they reach the task
			close(queue)
			for task := range queue {
				if m, ok := task.(*GenerateMarkersTask); ok {
					assert.Equal(t, tt.input.Markers, m.VideoPreview, "marker VideoPreview")
					assert.Equal(t, tt.input.MarkerImagePreviews, m.ImagePreview, "marker ImagePreview")
					assert.Equal(t, tt.input.MarkerScreenshots, m.Screenshot, "marker Screenshot")
				}
			}
		})
	}
}
