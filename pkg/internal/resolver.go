package internal

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	"context"
	"time"

	"github.com/stashapp/stash/pkg/models"
)

type Resolver struct{}

func (r *galleryResolver) Path(ctx context.Context, obj *models.Gallery) (*string, error) {
	panic("not implemented")
}

func (r *galleryResolver) Title(ctx context.Context, obj *models.Gallery) (*string, error) {
	panic("not implemented")
}

func (r *galleryResolver) URL(ctx context.Context, obj *models.Gallery) (*string, error) {
	panic("not implemented")
}

func (r *galleryResolver) Date(ctx context.Context, obj *models.Gallery) (*string, error) {
	panic("not implemented")
}

func (r *galleryResolver) Details(ctx context.Context, obj *models.Gallery) (*string, error) {
	panic("not implemented")
}

func (r *galleryResolver) Rating(ctx context.Context, obj *models.Gallery) (*int, error) {
	panic("not implemented")
}

func (r *galleryResolver) CreatedAt(ctx context.Context, obj *models.Gallery) (*time.Time, error) {
	panic("not implemented")
}

func (r *galleryResolver) UpdatedAt(ctx context.Context, obj *models.Gallery) (*time.Time, error) {
	panic("not implemented")
}

func (r *galleryResolver) FileModTime(ctx context.Context, obj *models.Gallery) (*time.Time, error) {
	panic("not implemented")
}

func (r *galleryResolver) Scenes(ctx context.Context, obj *models.Gallery) ([]*models.Scene, error) {
	panic("not implemented")
}

func (r *galleryResolver) Studio(ctx context.Context, obj *models.Gallery) (*models.Studio, error) {
	panic("not implemented")
}

func (r *galleryResolver) ImageCount(ctx context.Context, obj *models.Gallery) (int, error) {
	panic("not implemented")
}

func (r *galleryResolver) Tags(ctx context.Context, obj *models.Gallery) ([]*models.Tag, error) {
	panic("not implemented")
}

func (r *galleryResolver) Performers(ctx context.Context, obj *models.Gallery) ([]*models.Performer, error) {
	panic("not implemented")
}

func (r *galleryResolver) Images(ctx context.Context, obj *models.Gallery) ([]*models.Image, error) {
	panic("not implemented")
}

func (r *galleryResolver) Cover(ctx context.Context, obj *models.Gallery) (*models.Image, error) {
	panic("not implemented")
}

func (r *imageResolver) Title(ctx context.Context, obj *models.Image) (*string, error) {
	panic("not implemented")
}

func (r *imageResolver) Rating(ctx context.Context, obj *models.Image) (*int, error) {
	panic("not implemented")
}

func (r *imageResolver) CreatedAt(ctx context.Context, obj *models.Image) (*time.Time, error) {
	panic("not implemented")
}

func (r *imageResolver) UpdatedAt(ctx context.Context, obj *models.Image) (*time.Time, error) {
	panic("not implemented")
}

func (r *imageResolver) FileModTime(ctx context.Context, obj *models.Image) (*time.Time, error) {
	panic("not implemented")
}

func (r *imageResolver) File(ctx context.Context, obj *models.Image) (*models.ImageFileType, error) {
	panic("not implemented")
}

func (r *imageResolver) Paths(ctx context.Context, obj *models.Image) (*models.ImagePathsType, error) {
	panic("not implemented")
}

func (r *imageResolver) Galleries(ctx context.Context, obj *models.Image) ([]*models.Gallery, error) {
	panic("not implemented")
}

func (r *imageResolver) Studio(ctx context.Context, obj *models.Image) (*models.Studio, error) {
	panic("not implemented")
}

func (r *imageResolver) Tags(ctx context.Context, obj *models.Image) ([]*models.Tag, error) {
	panic("not implemented")
}

func (r *imageResolver) Performers(ctx context.Context, obj *models.Image) ([]*models.Performer, error) {
	panic("not implemented")
}

func (r *movieResolver) Name(ctx context.Context, obj *models.Movie) (string, error) {
	panic("not implemented")
}

func (r *movieResolver) Aliases(ctx context.Context, obj *models.Movie) (*string, error) {
	panic("not implemented")
}

func (r *movieResolver) Duration(ctx context.Context, obj *models.Movie) (*int, error) {
	panic("not implemented")
}

func (r *movieResolver) Date(ctx context.Context, obj *models.Movie) (*string, error) {
	panic("not implemented")
}

func (r *movieResolver) Rating(ctx context.Context, obj *models.Movie) (*int, error) {
	panic("not implemented")
}

func (r *movieResolver) Studio(ctx context.Context, obj *models.Movie) (*models.Studio, error) {
	panic("not implemented")
}

func (r *movieResolver) Director(ctx context.Context, obj *models.Movie) (*string, error) {
	panic("not implemented")
}

func (r *movieResolver) Synopsis(ctx context.Context, obj *models.Movie) (*string, error) {
	panic("not implemented")
}

func (r *movieResolver) URL(ctx context.Context, obj *models.Movie) (*string, error) {
	panic("not implemented")
}

func (r *movieResolver) CreatedAt(ctx context.Context, obj *models.Movie) (*time.Time, error) {
	panic("not implemented")
}

func (r *movieResolver) UpdatedAt(ctx context.Context, obj *models.Movie) (*time.Time, error) {
	panic("not implemented")
}

func (r *movieResolver) FrontImagePath(ctx context.Context, obj *models.Movie) (*string, error) {
	panic("not implemented")
}

func (r *movieResolver) BackImagePath(ctx context.Context, obj *models.Movie) (*string, error) {
	panic("not implemented")
}

func (r *movieResolver) SceneCount(ctx context.Context, obj *models.Movie) (*int, error) {
	panic("not implemented")
}

func (r *movieResolver) Scenes(ctx context.Context, obj *models.Movie) ([]*models.Scene, error) {
	panic("not implemented")
}

func (r *mutationResolver) Setup(ctx context.Context, input models.SetupInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) Migrate(ctx context.Context, input models.MigrateInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) SceneUpdate(ctx context.Context, input models.SceneUpdateInput) (*models.Scene, error) {
	panic("not implemented")
}

func (r *mutationResolver) BulkSceneUpdate(ctx context.Context, input models.BulkSceneUpdateInput) ([]*models.Scene, error) {
	panic("not implemented")
}

func (r *mutationResolver) SceneDestroy(ctx context.Context, input models.SceneDestroyInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) ScenesDestroy(ctx context.Context, input models.ScenesDestroyInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) ScenesUpdate(ctx context.Context, input []*models.SceneUpdateInput) ([]*models.Scene, error) {
	panic("not implemented")
}

func (r *mutationResolver) SceneIncrementO(ctx context.Context, id string) (int, error) {
	panic("not implemented")
}

func (r *mutationResolver) SceneDecrementO(ctx context.Context, id string) (int, error) {
	panic("not implemented")
}

func (r *mutationResolver) SceneResetO(ctx context.Context, id string) (int, error) {
	panic("not implemented")
}

func (r *mutationResolver) SceneGenerateScreenshot(ctx context.Context, id string, at *float64) (string, error) {
	panic("not implemented")
}

func (r *mutationResolver) SceneMarkerCreate(ctx context.Context, input models.SceneMarkerCreateInput) (*models.SceneMarker, error) {
	panic("not implemented")
}

func (r *mutationResolver) SceneMarkerUpdate(ctx context.Context, input models.SceneMarkerUpdateInput) (*models.SceneMarker, error) {
	panic("not implemented")
}

func (r *mutationResolver) SceneMarkerDestroy(ctx context.Context, id string) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) ImageUpdate(ctx context.Context, input models.ImageUpdateInput) (*models.Image, error) {
	panic("not implemented")
}

func (r *mutationResolver) BulkImageUpdate(ctx context.Context, input models.BulkImageUpdateInput) ([]*models.Image, error) {
	panic("not implemented")
}

func (r *mutationResolver) ImageDestroy(ctx context.Context, input models.ImageDestroyInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) ImagesDestroy(ctx context.Context, input models.ImagesDestroyInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) ImagesUpdate(ctx context.Context, input []*models.ImageUpdateInput) ([]*models.Image, error) {
	panic("not implemented")
}

func (r *mutationResolver) ImageIncrementO(ctx context.Context, id string) (int, error) {
	panic("not implemented")
}

func (r *mutationResolver) ImageDecrementO(ctx context.Context, id string) (int, error) {
	panic("not implemented")
}

func (r *mutationResolver) ImageResetO(ctx context.Context, id string) (int, error) {
	panic("not implemented")
}

func (r *mutationResolver) GalleryCreate(ctx context.Context, input models.GalleryCreateInput) (*models.Gallery, error) {
	panic("not implemented")
}

func (r *mutationResolver) GalleryUpdate(ctx context.Context, input models.GalleryUpdateInput) (*models.Gallery, error) {
	panic("not implemented")
}

func (r *mutationResolver) BulkGalleryUpdate(ctx context.Context, input models.BulkGalleryUpdateInput) ([]*models.Gallery, error) {
	panic("not implemented")
}

func (r *mutationResolver) GalleryDestroy(ctx context.Context, input models.GalleryDestroyInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) GalleriesUpdate(ctx context.Context, input []*models.GalleryUpdateInput) ([]*models.Gallery, error) {
	panic("not implemented")
}

func (r *mutationResolver) AddGalleryImages(ctx context.Context, input models.GalleryAddInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) RemoveGalleryImages(ctx context.Context, input models.GalleryRemoveInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) PerformerCreate(ctx context.Context, input models.PerformerCreateInput) (*models.Performer, error) {
	panic("not implemented")
}

func (r *mutationResolver) PerformerUpdate(ctx context.Context, input models.PerformerUpdateInput) (*models.Performer, error) {
	panic("not implemented")
}

func (r *mutationResolver) PerformerDestroy(ctx context.Context, input models.PerformerDestroyInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) PerformersDestroy(ctx context.Context, ids []string) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) BulkPerformerUpdate(ctx context.Context, input models.BulkPerformerUpdateInput) ([]*models.Performer, error) {
	panic("not implemented")
}

func (r *mutationResolver) StudioCreate(ctx context.Context, input models.StudioCreateInput) (*models.Studio, error) {
	panic("not implemented")
}

func (r *mutationResolver) StudioUpdate(ctx context.Context, input models.StudioUpdateInput) (*models.Studio, error) {
	panic("not implemented")
}

func (r *mutationResolver) StudioDestroy(ctx context.Context, input models.StudioDestroyInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) StudiosDestroy(ctx context.Context, ids []string) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) MovieCreate(ctx context.Context, input models.MovieCreateInput) (*models.Movie, error) {
	panic("not implemented")
}

func (r *mutationResolver) MovieUpdate(ctx context.Context, input models.MovieUpdateInput) (*models.Movie, error) {
	panic("not implemented")
}

func (r *mutationResolver) MovieDestroy(ctx context.Context, input models.MovieDestroyInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) MoviesDestroy(ctx context.Context, ids []string) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) BulkMovieUpdate(ctx context.Context, input models.BulkMovieUpdateInput) ([]*models.Movie, error) {
	panic("not implemented")
}

func (r *mutationResolver) TagCreate(ctx context.Context, input models.TagCreateInput) (*models.Tag, error) {
	panic("not implemented")
}

func (r *mutationResolver) TagUpdate(ctx context.Context, input models.TagUpdateInput) (*models.Tag, error) {
	panic("not implemented")
}

func (r *mutationResolver) TagDestroy(ctx context.Context, input models.TagDestroyInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) TagsDestroy(ctx context.Context, ids []string) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) TagsMerge(ctx context.Context, input models.TagsMergeInput) (*models.Tag, error) {
	panic("not implemented")
}

func (r *mutationResolver) SaveFilter(ctx context.Context, input models.SaveFilterInput) (*models.SavedFilter, error) {
	panic("not implemented")
}

func (r *mutationResolver) DestroySavedFilter(ctx context.Context, input models.DestroyFilterInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) SetDefaultFilter(ctx context.Context, input models.SetDefaultFilterInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) ConfigureGeneral(ctx context.Context, input models.ConfigGeneralInput) (*models.ConfigGeneralResult, error) {
	panic("not implemented")
}

func (r *mutationResolver) ConfigureInterface(ctx context.Context, input models.ConfigInterfaceInput) (*models.ConfigInterfaceResult, error) {
	panic("not implemented")
}

func (r *mutationResolver) ConfigureDlna(ctx context.Context, input models.ConfigDLNAInput) (*models.ConfigDLNAResult, error) {
	panic("not implemented")
}

func (r *mutationResolver) ConfigureScraping(ctx context.Context, input models.ConfigScrapingInput) (*models.ConfigScrapingResult, error) {
	panic("not implemented")
}

func (r *mutationResolver) ConfigureDefaults(ctx context.Context, input models.ConfigDefaultSettingsInput) (*models.ConfigDefaultSettingsResult, error) {
	panic("not implemented")
}

func (r *mutationResolver) GenerateAPIKey(ctx context.Context, input models.GenerateAPIKeyInput) (string, error) {
	panic("not implemented")
}

func (r *mutationResolver) ExportObjects(ctx context.Context, input models.ExportObjectsInput) (*string, error) {
	panic("not implemented")
}

func (r *mutationResolver) ImportObjects(ctx context.Context, input models.ImportObjectsInput) (string, error) {
	panic("not implemented")
}

func (r *mutationResolver) MetadataImport(ctx context.Context) (string, error) {
	panic("not implemented")
}

func (r *mutationResolver) MetadataExport(ctx context.Context) (string, error) {
	panic("not implemented")
}

func (r *mutationResolver) MetadataScan(ctx context.Context, input models.ScanMetadataInput) (string, error) {
	panic("not implemented")
}

func (r *mutationResolver) MetadataGenerate(ctx context.Context, input models.GenerateMetadataInput) (string, error) {
	panic("not implemented")
}

func (r *mutationResolver) MetadataAutoTag(ctx context.Context, input models.AutoTagMetadataInput) (string, error) {
	panic("not implemented")
}

func (r *mutationResolver) MetadataClean(ctx context.Context, input models.CleanMetadataInput) (string, error) {
	panic("not implemented")
}

func (r *mutationResolver) MetadataIdentify(ctx context.Context, input models.IdentifyMetadataInput) (string, error) {
	panic("not implemented")
}

func (r *mutationResolver) MigrateHashNaming(ctx context.Context) (string, error) {
	panic("not implemented")
}

func (r *mutationResolver) ReloadScrapers(ctx context.Context) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) RunPluginTask(ctx context.Context, pluginID string, taskName string, args []*models.PluginArgInput) (string, error) {
	panic("not implemented")
}

func (r *mutationResolver) ReloadPlugins(ctx context.Context) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) StopJob(ctx context.Context, jobID string) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) StopAllJobs(ctx context.Context) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) SubmitStashBoxFingerprints(ctx context.Context, input models.StashBoxFingerprintSubmissionInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) SubmitStashBoxSceneDraft(ctx context.Context, input models.StashBoxDraftSubmissionInput) (*string, error) {
	panic("not implemented")
}

func (r *mutationResolver) SubmitStashBoxPerformerDraft(ctx context.Context, input models.StashBoxDraftSubmissionInput) (*string, error) {
	panic("not implemented")
}

func (r *mutationResolver) BackupDatabase(ctx context.Context, input models.BackupDatabaseInput) (*string, error) {
	panic("not implemented")
}

func (r *mutationResolver) StashBoxBatchPerformerTag(ctx context.Context, input models.StashBoxBatchPerformerTagInput) (string, error) {
	panic("not implemented")
}

func (r *mutationResolver) EnableDlna(ctx context.Context, input models.EnableDLNAInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) DisableDlna(ctx context.Context, input models.DisableDLNAInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) AddTempDlnaip(ctx context.Context, input models.AddTempDLNAIPInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) RemoveTempDlnaip(ctx context.Context, input models.RemoveTempDLNAIPInput) (bool, error) {
	panic("not implemented")
}

func (r *performerResolver) Name(ctx context.Context, obj *models.Performer) (*string, error) {
	panic("not implemented")
}

func (r *performerResolver) URL(ctx context.Context, obj *models.Performer) (*string, error) {
	panic("not implemented")
}

func (r *performerResolver) Gender(ctx context.Context, obj *models.Performer) (*models.GenderEnum, error) {
	panic("not implemented")
}

func (r *performerResolver) Twitter(ctx context.Context, obj *models.Performer) (*string, error) {
	panic("not implemented")
}

func (r *performerResolver) Instagram(ctx context.Context, obj *models.Performer) (*string, error) {
	panic("not implemented")
}

func (r *performerResolver) Birthdate(ctx context.Context, obj *models.Performer) (*string, error) {
	panic("not implemented")
}

func (r *performerResolver) Ethnicity(ctx context.Context, obj *models.Performer) (*string, error) {
	panic("not implemented")
}

func (r *performerResolver) Country(ctx context.Context, obj *models.Performer) (*string, error) {
	panic("not implemented")
}

func (r *performerResolver) EyeColor(ctx context.Context, obj *models.Performer) (*string, error) {
	panic("not implemented")
}

func (r *performerResolver) Height(ctx context.Context, obj *models.Performer) (*string, error) {
	panic("not implemented")
}

func (r *performerResolver) Measurements(ctx context.Context, obj *models.Performer) (*string, error) {
	panic("not implemented")
}

func (r *performerResolver) FakeTits(ctx context.Context, obj *models.Performer) (*string, error) {
	panic("not implemented")
}

func (r *performerResolver) CareerLength(ctx context.Context, obj *models.Performer) (*string, error) {
	panic("not implemented")
}

func (r *performerResolver) Tattoos(ctx context.Context, obj *models.Performer) (*string, error) {
	panic("not implemented")
}

func (r *performerResolver) Piercings(ctx context.Context, obj *models.Performer) (*string, error) {
	panic("not implemented")
}

func (r *performerResolver) Aliases(ctx context.Context, obj *models.Performer) (*string, error) {
	panic("not implemented")
}

func (r *performerResolver) Favorite(ctx context.Context, obj *models.Performer) (bool, error) {
	panic("not implemented")
}

func (r *performerResolver) Tags(ctx context.Context, obj *models.Performer) ([]*models.Tag, error) {
	panic("not implemented")
}

func (r *performerResolver) ImagePath(ctx context.Context, obj *models.Performer) (*string, error) {
	panic("not implemented")
}

func (r *performerResolver) SceneCount(ctx context.Context, obj *models.Performer) (*int, error) {
	panic("not implemented")
}

func (r *performerResolver) ImageCount(ctx context.Context, obj *models.Performer) (*int, error) {
	panic("not implemented")
}

func (r *performerResolver) GalleryCount(ctx context.Context, obj *models.Performer) (*int, error) {
	panic("not implemented")
}

func (r *performerResolver) Scenes(ctx context.Context, obj *models.Performer) ([]*models.Scene, error) {
	panic("not implemented")
}

func (r *performerResolver) StashIds(ctx context.Context, obj *models.Performer) ([]*models.StashID, error) {
	panic("not implemented")
}

func (r *performerResolver) Rating(ctx context.Context, obj *models.Performer) (*int, error) {
	panic("not implemented")
}

func (r *performerResolver) Details(ctx context.Context, obj *models.Performer) (*string, error) {
	panic("not implemented")
}

func (r *performerResolver) DeathDate(ctx context.Context, obj *models.Performer) (*string, error) {
	panic("not implemented")
}

func (r *performerResolver) HairColor(ctx context.Context, obj *models.Performer) (*string, error) {
	panic("not implemented")
}

func (r *performerResolver) Weight(ctx context.Context, obj *models.Performer) (*int, error) {
	panic("not implemented")
}

func (r *performerResolver) CreatedAt(ctx context.Context, obj *models.Performer) (*time.Time, error) {
	panic("not implemented")
}

func (r *performerResolver) UpdatedAt(ctx context.Context, obj *models.Performer) (*time.Time, error) {
	panic("not implemented")
}

func (r *performerResolver) MovieCount(ctx context.Context, obj *models.Performer) (*int, error) {
	panic("not implemented")
}

func (r *performerResolver) Movies(ctx context.Context, obj *models.Performer) ([]*models.Movie, error) {
	panic("not implemented")
}

func (r *queryResolver) FindSavedFilters(ctx context.Context, mode models.FilterMode) ([]*models.SavedFilter, error) {
	panic("not implemented")
}

func (r *queryResolver) FindDefaultFilter(ctx context.Context, mode models.FilterMode) (*models.SavedFilter, error) {
	panic("not implemented")
}

func (r *queryResolver) FindScene(ctx context.Context, id *string, checksum *string) (*models.Scene, error) {
	panic("not implemented")
}

func (r *queryResolver) FindSceneByHash(ctx context.Context, input models.SceneHashInput) (*models.Scene, error) {
	panic("not implemented")
}

func (r *queryResolver) FindScenes(ctx context.Context, sceneFilter *models.SceneFilterType, sceneIds []int, filter *models.FindFilterType) (*models.FindScenesResultType, error) {
	panic("not implemented")
}

func (r *queryResolver) FindScenesByPathRegex(ctx context.Context, filter *models.FindFilterType) (*models.FindScenesResultType, error) {
	panic("not implemented")
}

func (r *queryResolver) FindDuplicateScenes(ctx context.Context, distance *int) ([][]*models.Scene, error) {
	panic("not implemented")
}

func (r *queryResolver) SceneStreams(ctx context.Context, id *string) ([]*models.SceneStreamEndpoint, error) {
	panic("not implemented")
}

func (r *queryResolver) ParseSceneFilenames(ctx context.Context, filter *models.FindFilterType, config models.SceneParserInput) (*models.SceneParserResultType, error) {
	panic("not implemented")
}

func (r *queryResolver) FindSceneMarkers(ctx context.Context, sceneMarkerFilter *models.SceneMarkerFilterType, filter *models.FindFilterType) (*models.FindSceneMarkersResultType, error) {
	panic("not implemented")
}

func (r *queryResolver) FindImage(ctx context.Context, id *string, checksum *string) (*models.Image, error) {
	panic("not implemented")
}

func (r *queryResolver) FindImages(ctx context.Context, imageFilter *models.ImageFilterType, imageIds []int, filter *models.FindFilterType) (*models.FindImagesResultType, error) {
	panic("not implemented")
}

func (r *queryResolver) FindPerformer(ctx context.Context, id string) (*models.Performer, error) {
	panic("not implemented")
}

func (r *queryResolver) FindPerformers(ctx context.Context, performerFilter *models.PerformerFilterType, filter *models.FindFilterType) (*models.FindPerformersResultType, error) {
	panic("not implemented")
}

func (r *queryResolver) FindStudio(ctx context.Context, id string) (*models.Studio, error) {
	panic("not implemented")
}

func (r *queryResolver) FindStudios(ctx context.Context, studioFilter *models.StudioFilterType, filter *models.FindFilterType) (*models.FindStudiosResultType, error) {
	panic("not implemented")
}

func (r *queryResolver) FindMovie(ctx context.Context, id string) (*models.Movie, error) {
	panic("not implemented")
}

func (r *queryResolver) FindMovies(ctx context.Context, movieFilter *models.MovieFilterType, filter *models.FindFilterType) (*models.FindMoviesResultType, error) {
	panic("not implemented")
}

func (r *queryResolver) FindGallery(ctx context.Context, id string) (*models.Gallery, error) {
	panic("not implemented")
}

func (r *queryResolver) FindGalleries(ctx context.Context, galleryFilter *models.GalleryFilterType, filter *models.FindFilterType) (*models.FindGalleriesResultType, error) {
	panic("not implemented")
}

func (r *queryResolver) FindTag(ctx context.Context, id string) (*models.Tag, error) {
	panic("not implemented")
}

func (r *queryResolver) FindTags(ctx context.Context, tagFilter *models.TagFilterType, filter *models.FindFilterType) (*models.FindTagsResultType, error) {
	panic("not implemented")
}

func (r *queryResolver) MarkerWall(ctx context.Context, q *string) ([]*models.SceneMarker, error) {
	panic("not implemented")
}

func (r *queryResolver) SceneWall(ctx context.Context, q *string) ([]*models.Scene, error) {
	panic("not implemented")
}

func (r *queryResolver) MarkerStrings(ctx context.Context, q *string, sort *string) ([]*models.MarkerStringsResultType, error) {
	panic("not implemented")
}

func (r *queryResolver) Stats(ctx context.Context) (*models.StatsResultType, error) {
	panic("not implemented")
}

func (r *queryResolver) SceneMarkerTags(ctx context.Context, sceneID string) ([]*models.SceneMarkerTag, error) {
	panic("not implemented")
}

func (r *queryResolver) Logs(ctx context.Context) ([]*models.LogEntry, error) {
	panic("not implemented")
}

func (r *queryResolver) ListScrapers(ctx context.Context, types []models.ScrapeContentType) ([]*models.Scraper, error) {
	panic("not implemented")
}

func (r *queryResolver) ListPerformerScrapers(ctx context.Context) ([]*models.Scraper, error) {
	panic("not implemented")
}

func (r *queryResolver) ListSceneScrapers(ctx context.Context) ([]*models.Scraper, error) {
	panic("not implemented")
}

func (r *queryResolver) ListGalleryScrapers(ctx context.Context) ([]*models.Scraper, error) {
	panic("not implemented")
}

func (r *queryResolver) ListMovieScrapers(ctx context.Context) ([]*models.Scraper, error) {
	panic("not implemented")
}

func (r *queryResolver) ScrapeSingleScene(ctx context.Context, source models.ScraperSourceInput, input models.ScrapeSingleSceneInput) ([]*models.ScrapedScene, error) {
	panic("not implemented")
}

func (r *queryResolver) ScrapeMultiScenes(ctx context.Context, source models.ScraperSourceInput, input models.ScrapeMultiScenesInput) ([][]*models.ScrapedScene, error) {
	panic("not implemented")
}

func (r *queryResolver) ScrapeSinglePerformer(ctx context.Context, source models.ScraperSourceInput, input models.ScrapeSinglePerformerInput) ([]*models.ScrapedPerformer, error) {
	panic("not implemented")
}

func (r *queryResolver) ScrapeMultiPerformers(ctx context.Context, source models.ScraperSourceInput, input models.ScrapeMultiPerformersInput) ([][]*models.ScrapedPerformer, error) {
	panic("not implemented")
}

func (r *queryResolver) ScrapeSingleGallery(ctx context.Context, source models.ScraperSourceInput, input models.ScrapeSingleGalleryInput) ([]*models.ScrapedGallery, error) {
	panic("not implemented")
}

func (r *queryResolver) ScrapeSingleMovie(ctx context.Context, source models.ScraperSourceInput, input models.ScrapeSingleMovieInput) ([]*models.ScrapedMovie, error) {
	panic("not implemented")
}

func (r *queryResolver) ScrapeURL(ctx context.Context, url string, ty models.ScrapeContentType) (models.ScrapedContent, error) {
	panic("not implemented")
}

func (r *queryResolver) ScrapePerformerURL(ctx context.Context, url string) (*models.ScrapedPerformer, error) {
	panic("not implemented")
}

func (r *queryResolver) ScrapeSceneURL(ctx context.Context, url string) (*models.ScrapedScene, error) {
	panic("not implemented")
}

func (r *queryResolver) ScrapeGalleryURL(ctx context.Context, url string) (*models.ScrapedGallery, error) {
	panic("not implemented")
}

func (r *queryResolver) ScrapeMovieURL(ctx context.Context, url string) (*models.ScrapedMovie, error) {
	panic("not implemented")
}

func (r *queryResolver) ScrapePerformerList(ctx context.Context, scraperID string, query string) ([]*models.ScrapedPerformer, error) {
	panic("not implemented")
}

func (r *queryResolver) ScrapePerformer(ctx context.Context, scraperID string, scrapedPerformer models.ScrapedPerformerInput) (*models.ScrapedPerformer, error) {
	panic("not implemented")
}

func (r *queryResolver) ScrapeScene(ctx context.Context, scraperID string, scene models.SceneUpdateInput) (*models.ScrapedScene, error) {
	panic("not implemented")
}

func (r *queryResolver) ScrapeGallery(ctx context.Context, scraperID string, gallery models.GalleryUpdateInput) (*models.ScrapedGallery, error) {
	panic("not implemented")
}

func (r *queryResolver) ScrapeFreeonesPerformerList(ctx context.Context, query string) ([]string, error) {
	panic("not implemented")
}

func (r *queryResolver) QueryStashBoxScene(ctx context.Context, input models.StashBoxSceneQueryInput) ([]*models.ScrapedScene, error) {
	panic("not implemented")
}

func (r *queryResolver) QueryStashBoxPerformer(ctx context.Context, input models.StashBoxPerformerQueryInput) ([]*models.StashBoxPerformerQueryResult, error) {
	panic("not implemented")
}

func (r *queryResolver) Plugins(ctx context.Context) ([]*models.Plugin, error) {
	panic("not implemented")
}

func (r *queryResolver) PluginTasks(ctx context.Context) ([]*models.PluginTask, error) {
	panic("not implemented")
}

func (r *queryResolver) Configuration(ctx context.Context) (*models.ConfigResult, error) {
	panic("not implemented")
}

func (r *queryResolver) Directory(ctx context.Context, path *string, locale *string) (*models.Directory, error) {
	panic("not implemented")
}

func (r *queryResolver) ValidateStashBoxCredentials(ctx context.Context, input models.StashBoxInput) (*models.StashBoxValidationResult, error) {
	panic("not implemented")
}

func (r *queryResolver) SystemStatus(ctx context.Context) (*models.SystemStatus, error) {
	panic("not implemented")
}

func (r *queryResolver) JobQueue(ctx context.Context) ([]*models.Job, error) {
	panic("not implemented")
}

func (r *queryResolver) FindJob(ctx context.Context, input models.FindJobInput) (*models.Job, error) {
	panic("not implemented")
}

func (r *queryResolver) DlnaStatus(ctx context.Context) (*models.DLNAStatus, error) {
	panic("not implemented")
}

func (r *queryResolver) AllPerformers(ctx context.Context) ([]*models.Performer, error) {
	panic("not implemented")
}

func (r *queryResolver) AllStudios(ctx context.Context) ([]*models.Studio, error) {
	panic("not implemented")
}

func (r *queryResolver) AllMovies(ctx context.Context) ([]*models.Movie, error) {
	panic("not implemented")
}

func (r *queryResolver) AllTags(ctx context.Context) ([]*models.Tag, error) {
	panic("not implemented")
}

func (r *queryResolver) Version(ctx context.Context) (*models.Version, error) {
	panic("not implemented")
}

func (r *queryResolver) Latestversion(ctx context.Context) (*models.ShortVersion, error) {
	panic("not implemented")
}

func (r *sceneResolver) Checksum(ctx context.Context, obj *models.Scene) (*string, error) {
	panic("not implemented")
}

func (r *sceneResolver) Oshash(ctx context.Context, obj *models.Scene) (*string, error) {
	panic("not implemented")
}

func (r *sceneResolver) Title(ctx context.Context, obj *models.Scene) (*string, error) {
	panic("not implemented")
}

func (r *sceneResolver) Details(ctx context.Context, obj *models.Scene) (*string, error) {
	panic("not implemented")
}

func (r *sceneResolver) URL(ctx context.Context, obj *models.Scene) (*string, error) {
	panic("not implemented")
}

func (r *sceneResolver) Date(ctx context.Context, obj *models.Scene) (*string, error) {
	panic("not implemented")
}

func (r *sceneResolver) Rating(ctx context.Context, obj *models.Scene) (*int, error) {
	panic("not implemented")
}

func (r *sceneResolver) Phash(ctx context.Context, obj *models.Scene) (*string, error) {
	panic("not implemented")
}

func (r *sceneResolver) InteractiveSpeed(ctx context.Context, obj *models.Scene) (*int, error) {
	panic("not implemented")
}

func (r *sceneResolver) CreatedAt(ctx context.Context, obj *models.Scene) (*time.Time, error) {
	panic("not implemented")
}

func (r *sceneResolver) UpdatedAt(ctx context.Context, obj *models.Scene) (*time.Time, error) {
	panic("not implemented")
}

func (r *sceneResolver) FileModTime(ctx context.Context, obj *models.Scene) (*time.Time, error) {
	panic("not implemented")
}

func (r *sceneResolver) File(ctx context.Context, obj *models.Scene) (*models.SceneFileType, error) {
	panic("not implemented")
}

func (r *sceneResolver) Paths(ctx context.Context, obj *models.Scene) (*models.ScenePathsType, error) {
	panic("not implemented")
}

func (r *sceneResolver) SceneMarkers(ctx context.Context, obj *models.Scene) ([]*models.SceneMarker, error) {
	panic("not implemented")
}

func (r *sceneResolver) Galleries(ctx context.Context, obj *models.Scene) ([]*models.Gallery, error) {
	panic("not implemented")
}

func (r *sceneResolver) Studio(ctx context.Context, obj *models.Scene) (*models.Studio, error) {
	panic("not implemented")
}

func (r *sceneResolver) Movies(ctx context.Context, obj *models.Scene) ([]*models.SceneMovie, error) {
	panic("not implemented")
}

func (r *sceneResolver) Tags(ctx context.Context, obj *models.Scene) ([]*models.Tag, error) {
	panic("not implemented")
}

func (r *sceneResolver) Performers(ctx context.Context, obj *models.Scene) ([]*models.Performer, error) {
	panic("not implemented")
}

func (r *sceneResolver) StashIds(ctx context.Context, obj *models.Scene) ([]*models.StashID, error) {
	panic("not implemented")
}

func (r *sceneMarkerResolver) Scene(ctx context.Context, obj *models.SceneMarker) (*models.Scene, error) {
	panic("not implemented")
}

func (r *sceneMarkerResolver) PrimaryTag(ctx context.Context, obj *models.SceneMarker) (*models.Tag, error) {
	panic("not implemented")
}

func (r *sceneMarkerResolver) Tags(ctx context.Context, obj *models.SceneMarker) ([]*models.Tag, error) {
	panic("not implemented")
}

func (r *sceneMarkerResolver) CreatedAt(ctx context.Context, obj *models.SceneMarker) (*time.Time, error) {
	panic("not implemented")
}

func (r *sceneMarkerResolver) UpdatedAt(ctx context.Context, obj *models.SceneMarker) (*time.Time, error) {
	panic("not implemented")
}

func (r *sceneMarkerResolver) Stream(ctx context.Context, obj *models.SceneMarker) (string, error) {
	panic("not implemented")
}

func (r *sceneMarkerResolver) Preview(ctx context.Context, obj *models.SceneMarker) (string, error) {
	panic("not implemented")
}

func (r *sceneMarkerResolver) Screenshot(ctx context.Context, obj *models.SceneMarker) (string, error) {
	panic("not implemented")
}

func (r *studioResolver) Name(ctx context.Context, obj *models.Studio) (string, error) {
	panic("not implemented")
}

func (r *studioResolver) URL(ctx context.Context, obj *models.Studio) (*string, error) {
	panic("not implemented")
}

func (r *studioResolver) ParentStudio(ctx context.Context, obj *models.Studio) (*models.Studio, error) {
	panic("not implemented")
}

func (r *studioResolver) ChildStudios(ctx context.Context, obj *models.Studio) ([]*models.Studio, error) {
	panic("not implemented")
}

func (r *studioResolver) Aliases(ctx context.Context, obj *models.Studio) ([]string, error) {
	panic("not implemented")
}

func (r *studioResolver) ImagePath(ctx context.Context, obj *models.Studio) (*string, error) {
	panic("not implemented")
}

func (r *studioResolver) SceneCount(ctx context.Context, obj *models.Studio) (*int, error) {
	panic("not implemented")
}

func (r *studioResolver) ImageCount(ctx context.Context, obj *models.Studio) (*int, error) {
	panic("not implemented")
}

func (r *studioResolver) GalleryCount(ctx context.Context, obj *models.Studio) (*int, error) {
	panic("not implemented")
}

func (r *studioResolver) StashIds(ctx context.Context, obj *models.Studio) ([]*models.StashID, error) {
	panic("not implemented")
}

func (r *studioResolver) Rating(ctx context.Context, obj *models.Studio) (*int, error) {
	panic("not implemented")
}

func (r *studioResolver) Details(ctx context.Context, obj *models.Studio) (*string, error) {
	panic("not implemented")
}

func (r *studioResolver) CreatedAt(ctx context.Context, obj *models.Studio) (*time.Time, error) {
	panic("not implemented")
}

func (r *studioResolver) UpdatedAt(ctx context.Context, obj *models.Studio) (*time.Time, error) {
	panic("not implemented")
}

func (r *studioResolver) MovieCount(ctx context.Context, obj *models.Studio) (*int, error) {
	panic("not implemented")
}

func (r *studioResolver) Movies(ctx context.Context, obj *models.Studio) ([]*models.Movie, error) {
	panic("not implemented")
}

func (r *subscriptionResolver) JobsSubscribe(ctx context.Context) (<-chan *models.JobStatusUpdate, error) {
	panic("not implemented")
}

func (r *subscriptionResolver) LoggingSubscribe(ctx context.Context) (<-chan []*models.LogEntry, error) {
	panic("not implemented")
}

func (r *subscriptionResolver) ScanCompleteSubscribe(ctx context.Context) (<-chan bool, error) {
	panic("not implemented")
}

func (r *tagResolver) Aliases(ctx context.Context, obj *models.Tag) ([]string, error) {
	panic("not implemented")
}

func (r *tagResolver) CreatedAt(ctx context.Context, obj *models.Tag) (*time.Time, error) {
	panic("not implemented")
}

func (r *tagResolver) UpdatedAt(ctx context.Context, obj *models.Tag) (*time.Time, error) {
	panic("not implemented")
}

func (r *tagResolver) ImagePath(ctx context.Context, obj *models.Tag) (*string, error) {
	panic("not implemented")
}

func (r *tagResolver) SceneCount(ctx context.Context, obj *models.Tag) (*int, error) {
	panic("not implemented")
}

func (r *tagResolver) SceneMarkerCount(ctx context.Context, obj *models.Tag) (*int, error) {
	panic("not implemented")
}

func (r *tagResolver) ImageCount(ctx context.Context, obj *models.Tag) (*int, error) {
	panic("not implemented")
}

func (r *tagResolver) GalleryCount(ctx context.Context, obj *models.Tag) (*int, error) {
	panic("not implemented")
}

func (r *tagResolver) PerformerCount(ctx context.Context, obj *models.Tag) (*int, error) {
	panic("not implemented")
}

func (r *tagResolver) Parents(ctx context.Context, obj *models.Tag) ([]*models.Tag, error) {
	panic("not implemented")
}

func (r *tagResolver) Children(ctx context.Context, obj *models.Tag) ([]*models.Tag, error) {
	panic("not implemented")
}

// Gallery returns models.GalleryResolver implementation.
func (r *Resolver) Gallery() models.GalleryResolver { return &galleryResolver{r} }

// Image returns models.ImageResolver implementation.
func (r *Resolver) Image() models.ImageResolver { return &imageResolver{r} }

// Movie returns models.MovieResolver implementation.
func (r *Resolver) Movie() models.MovieResolver { return &movieResolver{r} }

// Mutation returns models.MutationResolver implementation.
func (r *Resolver) Mutation() models.MutationResolver { return &mutationResolver{r} }

// Performer returns models.PerformerResolver implementation.
func (r *Resolver) Performer() models.PerformerResolver { return &performerResolver{r} }

// Query returns models.QueryResolver implementation.
func (r *Resolver) Query() models.QueryResolver { return &queryResolver{r} }

// Scene returns models.SceneResolver implementation.
func (r *Resolver) Scene() models.SceneResolver { return &sceneResolver{r} }

// SceneMarker returns models.SceneMarkerResolver implementation.
func (r *Resolver) SceneMarker() models.SceneMarkerResolver { return &sceneMarkerResolver{r} }

// Studio returns models.StudioResolver implementation.
func (r *Resolver) Studio() models.StudioResolver { return &studioResolver{r} }

// Subscription returns models.SubscriptionResolver implementation.
func (r *Resolver) Subscription() models.SubscriptionResolver { return &subscriptionResolver{r} }

// Tag returns models.TagResolver implementation.
func (r *Resolver) Tag() models.TagResolver { return &tagResolver{r} }

type galleryResolver struct{ *Resolver }
type imageResolver struct{ *Resolver }
type movieResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type performerResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type sceneResolver struct{ *Resolver }
type sceneMarkerResolver struct{ *Resolver }
type studioResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
type tagResolver struct{ *Resolver }
