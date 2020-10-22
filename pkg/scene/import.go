package scene

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/stashapp/stash/pkg/manager/jsonschema"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type Importer struct {
	ReaderWriter        models.SceneReaderWriter
	StudioWriter        models.StudioReaderWriter
	GalleryWriter       models.GalleryReaderWriter
	PerformerWriter     models.PerformerReaderWriter
	MovieWriter         models.MovieReaderWriter
	TagWriter           models.TagReaderWriter
	JoinWriter          models.JoinReaderWriter
	Input               jsonschema.Scene
	Path                string
	MissingRefBehaviour models.ImportMissingRefEnum
	FileNamingAlgorithm models.HashAlgorithm

	ID             int
	scene          models.Scene
	gallery        *models.Gallery
	performers     []*models.Performer
	movies         []models.MoviesScenes
	tags           []*models.Tag
	coverImageData []byte
}

func (i *Importer) PreImport() error {
	i.scene = i.sceneJSONToScene(i.Input)

	if err := i.populateStudio(); err != nil {
		return err
	}

	if err := i.populateGallery(); err != nil {
		return err
	}

	if err := i.populatePerformers(); err != nil {
		return err
	}

	if err := i.populateTags(); err != nil {
		return err
	}

	if err := i.populateMovies(); err != nil {
		return err
	}

	var err error
	if len(i.Input.Cover) > 0 {
		_, i.coverImageData, err = utils.ProcessBase64Image(i.Input.Cover)
		if err != nil {
			return fmt.Errorf("invalid cover image: %s", err.Error())
		}
	}

	return nil
}

func (i *Importer) sceneJSONToScene(sceneJSON jsonschema.Scene) models.Scene {
	newScene := models.Scene{
		Checksum: sql.NullString{String: sceneJSON.Checksum, Valid: sceneJSON.Checksum != ""},
		OSHash:   sql.NullString{String: sceneJSON.OSHash, Valid: sceneJSON.OSHash != ""},
		Path:     i.Path,
	}

	if sceneJSON.Title != "" {
		newScene.Title = sql.NullString{String: sceneJSON.Title, Valid: true}
	}
	if sceneJSON.Details != "" {
		newScene.Details = sql.NullString{String: sceneJSON.Details, Valid: true}
	}
	if sceneJSON.URL != "" {
		newScene.URL = sql.NullString{String: sceneJSON.URL, Valid: true}
	}
	if sceneJSON.Date != "" {
		newScene.Date = models.SQLiteDate{String: sceneJSON.Date, Valid: true}
	}
	if sceneJSON.Rating != 0 {
		newScene.Rating = sql.NullInt64{Int64: int64(sceneJSON.Rating), Valid: true}
	}

	newScene.OCounter = sceneJSON.OCounter
	newScene.CreatedAt = models.SQLiteTimestamp{Timestamp: sceneJSON.CreatedAt.GetTime()}
	newScene.UpdatedAt = models.SQLiteTimestamp{Timestamp: sceneJSON.UpdatedAt.GetTime()}

	if sceneJSON.File != nil {
		if sceneJSON.File.Size != "" {
			newScene.Size = sql.NullString{String: sceneJSON.File.Size, Valid: true}
		}
		if sceneJSON.File.Duration != "" {
			duration, _ := strconv.ParseFloat(sceneJSON.File.Duration, 64)
			newScene.Duration = sql.NullFloat64{Float64: duration, Valid: true}
		}
		if sceneJSON.File.VideoCodec != "" {
			newScene.VideoCodec = sql.NullString{String: sceneJSON.File.VideoCodec, Valid: true}
		}
		if sceneJSON.File.AudioCodec != "" {
			newScene.AudioCodec = sql.NullString{String: sceneJSON.File.AudioCodec, Valid: true}
		}
		if sceneJSON.File.Format != "" {
			newScene.Format = sql.NullString{String: sceneJSON.File.Format, Valid: true}
		}
		if sceneJSON.File.Width != 0 {
			newScene.Width = sql.NullInt64{Int64: int64(sceneJSON.File.Width), Valid: true}
		}
		if sceneJSON.File.Height != 0 {
			newScene.Height = sql.NullInt64{Int64: int64(sceneJSON.File.Height), Valid: true}
		}
		if sceneJSON.File.Framerate != "" {
			framerate, _ := strconv.ParseFloat(sceneJSON.File.Framerate, 64)
			newScene.Framerate = sql.NullFloat64{Float64: framerate, Valid: true}
		}
		if sceneJSON.File.Bitrate != 0 {
			newScene.Bitrate = sql.NullInt64{Int64: int64(sceneJSON.File.Bitrate), Valid: true}
		}
	}

	return newScene
}

func (i *Importer) populateStudio() error {
	if i.Input.Studio != "" {
		studio, err := i.StudioWriter.FindByName(i.Input.Studio, false)
		if err != nil {
			return fmt.Errorf("error finding studio by name: %s", err.Error())
		}

		if studio == nil {
			if i.MissingRefBehaviour == models.ImportMissingRefEnumFail {
				return fmt.Errorf("scene studio '%s' not found", i.Input.Studio)
			}

			if i.MissingRefBehaviour == models.ImportMissingRefEnumIgnore {
				return nil
			}

			if i.MissingRefBehaviour == models.ImportMissingRefEnumCreate {
				studioID, err := i.createStudio(i.Input.Studio)
				if err != nil {
					return err
				}
				i.scene.StudioID = sql.NullInt64{
					Int64: int64(studioID),
					Valid: true,
				}
			}
		} else {
			i.scene.StudioID = sql.NullInt64{Int64: int64(studio.ID), Valid: true}
		}
	}

	return nil
}

func (i *Importer) createStudio(name string) (int, error) {
	newStudio := *models.NewStudio(name)

	created, err := i.StudioWriter.Create(newStudio)
	if err != nil {
		return 0, err
	}

	return created.ID, nil
}

func (i *Importer) populateGallery() error {
	if i.Input.Gallery != "" {
		gallery, err := i.GalleryWriter.FindByChecksum(i.Input.Gallery)
		if err != nil {
			return fmt.Errorf("error finding gallery: %s", err.Error())
		}

		if gallery == nil {
			if i.MissingRefBehaviour == models.ImportMissingRefEnumFail {
				return fmt.Errorf("scene gallery '%s' not found", i.Input.Studio)
			}

			// we don't create galleries - just ignore
			if i.MissingRefBehaviour == models.ImportMissingRefEnumIgnore || i.MissingRefBehaviour == models.ImportMissingRefEnumCreate {
				return nil
			}
		} else {
			i.gallery = gallery
		}
	}

	return nil
}

func (i *Importer) populatePerformers() error {
	if len(i.Input.Performers) > 0 {
		names := i.Input.Performers
		performers, err := i.PerformerWriter.FindByNames(names, false)
		if err != nil {
			return err
		}

		var pluckedNames []string
		for _, performer := range performers {
			if !performer.Name.Valid {
				continue
			}
			pluckedNames = append(pluckedNames, performer.Name.String)
		}

		missingPerformers := utils.StrFilter(names, func(name string) bool {
			return !utils.StrInclude(pluckedNames, name)
		})

		if len(missingPerformers) > 0 {
			if i.MissingRefBehaviour == models.ImportMissingRefEnumFail {
				return fmt.Errorf("scene performers [%s] not found", strings.Join(missingPerformers, ", "))
			}

			if i.MissingRefBehaviour == models.ImportMissingRefEnumCreate {
				createdPerformers, err := i.createPerformers(missingPerformers)
				if err != nil {
					return fmt.Errorf("error creating scene performers: %s", err.Error())
				}

				performers = append(performers, createdPerformers...)
			}

			// ignore if MissingRefBehaviour set to Ignore
		}

		i.performers = performers
	}

	return nil
}

func (i *Importer) createPerformers(names []string) ([]*models.Performer, error) {
	var ret []*models.Performer
	for _, name := range names {
		newPerformer := *models.NewPerformer(name)

		created, err := i.PerformerWriter.Create(newPerformer)
		if err != nil {
			return nil, err
		}

		ret = append(ret, created)
	}

	return ret, nil
}

func (i *Importer) populateMovies() error {
	if len(i.Input.Movies) > 0 {
		for _, inputMovie := range i.Input.Movies {
			movie, err := i.MovieWriter.FindByName(inputMovie.MovieName, false)
			if err != nil {
				return fmt.Errorf("error finding scene movie: %s", err.Error())
			}

			if movie == nil {
				if i.MissingRefBehaviour == models.ImportMissingRefEnumFail {
					return fmt.Errorf("scene movie [%s] not found", inputMovie.MovieName)
				}

				if i.MissingRefBehaviour == models.ImportMissingRefEnumCreate {
					movie, err = i.createMovie(inputMovie.MovieName)
					if err != nil {
						return fmt.Errorf("error creating scene movie: %s", err.Error())
					}
				}

				// ignore if MissingRefBehaviour set to Ignore
				if i.MissingRefBehaviour == models.ImportMissingRefEnumIgnore {
					continue
				}
			}

			toAdd := models.MoviesScenes{
				MovieID: movie.ID,
			}

			if inputMovie.SceneIndex != 0 {
				toAdd.SceneIndex = sql.NullInt64{
					Int64: int64(inputMovie.SceneIndex),
					Valid: true,
				}
			}

			i.movies = append(i.movies, toAdd)
		}
	}

	return nil
}

func (i *Importer) createMovie(name string) (*models.Movie, error) {
	newMovie := *models.NewMovie(name)

	created, err := i.MovieWriter.Create(newMovie)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (i *Importer) populateTags() error {
	if len(i.Input.Tags) > 0 {

		tags, err := importTags(i.TagWriter, i.Input.Tags, i.MissingRefBehaviour)
		if err != nil {
			return err
		}

		i.tags = tags
	}

	return nil
}

func (i *Importer) PostImport(id int) error {
	if len(i.coverImageData) > 0 {
		if err := i.ReaderWriter.UpdateSceneCover(id, i.coverImageData); err != nil {
			return fmt.Errorf("error setting scene images: %s", err.Error())
		}
	}

	if i.gallery != nil {
		i.gallery.SceneID = sql.NullInt64{Int64: int64(id), Valid: true}
		_, err := i.GalleryWriter.Update(*i.gallery)
		if err != nil {
			return fmt.Errorf("failed to update gallery: %s", err.Error())
		}
	}

	if len(i.performers) > 0 {
		var performerJoins []models.PerformersScenes
		for _, performer := range i.performers {
			join := models.PerformersScenes{
				PerformerID: performer.ID,
				SceneID:     id,
			}
			performerJoins = append(performerJoins, join)
		}
		if err := i.JoinWriter.UpdatePerformersScenes(id, performerJoins); err != nil {
			return fmt.Errorf("failed to associate performers: %s", err.Error())
		}
	}

	if len(i.movies) > 0 {
		for index := range i.movies {
			i.movies[index].SceneID = id
		}
		if err := i.JoinWriter.UpdateMoviesScenes(id, i.movies); err != nil {
			return fmt.Errorf("failed to associate movies: %s", err.Error())
		}
	}

	if len(i.tags) > 0 {
		var tagJoins []models.ScenesTags
		for _, tag := range i.tags {
			join := models.ScenesTags{
				SceneID: id,
				TagID:   tag.ID,
			}
			tagJoins = append(tagJoins, join)
		}
		if err := i.JoinWriter.UpdateScenesTags(id, tagJoins); err != nil {
			return fmt.Errorf("failed to associate tags: %s", err.Error())
		}
	}

	return nil
}

func (i *Importer) Name() string {
	return i.Path
}

func (i *Importer) FindExistingID() (*int, error) {
	var existing *models.Scene
	var err error
	if i.FileNamingAlgorithm == models.HashAlgorithmMd5 {
		existing, err = i.ReaderWriter.FindByChecksum(i.Input.Checksum)
	} else if i.FileNamingAlgorithm == models.HashAlgorithmOshash {
		existing, err = i.ReaderWriter.FindByOSHash(i.Input.OSHash)
	} else {
		panic("unknown file naming algorithm")
	}

	if err != nil {
		return nil, err
	}

	if existing != nil {
		id := existing.ID
		return &id, nil
	}

	return nil, nil
}

func (i *Importer) Create() (*int, error) {
	created, err := i.ReaderWriter.Create(i.scene)
	if err != nil {
		return nil, fmt.Errorf("error creating scene: %s", err.Error())
	}

	id := created.ID
	i.ID = id
	return &id, nil
}

func (i *Importer) Update(id int) error {
	scene := i.scene
	scene.ID = id
	i.ID = id
	_, err := i.ReaderWriter.UpdateFull(scene)
	if err != nil {
		return fmt.Errorf("error updating existing scene: %s", err.Error())
	}

	return nil
}

func importTags(tagWriter models.TagReaderWriter, names []string, missingRefBehaviour models.ImportMissingRefEnum) ([]*models.Tag, error) {
	tags, err := tagWriter.FindByNames(names, false)
	if err != nil {
		return nil, err
	}

	var pluckedNames []string
	for _, tag := range tags {
		pluckedNames = append(pluckedNames, tag.Name)
	}

	missingTags := utils.StrFilter(names, func(name string) bool {
		return !utils.StrInclude(pluckedNames, name)
	})

	if len(missingTags) > 0 {
		if missingRefBehaviour == models.ImportMissingRefEnumFail {
			return nil, fmt.Errorf("tags [%s] not found", strings.Join(missingTags, ", "))
		}

		if missingRefBehaviour == models.ImportMissingRefEnumCreate {
			createdTags, err := createTags(tagWriter, missingTags)
			if err != nil {
				return nil, fmt.Errorf("error creating tags: %s", err.Error())
			}

			tags = append(tags, createdTags...)
		}

		// ignore if MissingRefBehaviour set to Ignore
	}

	return tags, nil
}

func createTags(tagWriter models.TagWriter, names []string) ([]*models.Tag, error) {
	var ret []*models.Tag
	for _, name := range names {
		newTag := *models.NewTag(name)

		created, err := tagWriter.Create(newTag)
		if err != nil {
			return nil, err
		}

		ret = append(ret, created)
	}

	return ret, nil
}
