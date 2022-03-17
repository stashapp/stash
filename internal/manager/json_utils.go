package manager

import (
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/models/paths"
)

type jsonUtils struct {
	json paths.JSONPaths
}

func (jp *jsonUtils) getMappings() (*jsonschema.Mappings, error) {
	return jsonschema.LoadMappingsFile(jp.json.MappingsFile)
}

func (jp *jsonUtils) saveMappings(mappings *jsonschema.Mappings) error {
	return jsonschema.SaveMappingsFile(jp.json.MappingsFile, mappings)
}

func (jp *jsonUtils) getScraped() ([]jsonschema.ScrapedItem, error) {
	return jsonschema.LoadScrapedFile(jp.json.ScrapedFile)
}

func (jp *jsonUtils) saveScaped(scraped []jsonschema.ScrapedItem) error {
	return jsonschema.SaveScrapedFile(jp.json.ScrapedFile, scraped)
}

func (jp *jsonUtils) getPerformer(checksum string) (*jsonschema.Performer, error) {
	return jsonschema.LoadPerformerFile(jp.json.PerformerJSONPath(checksum))
}

func (jp *jsonUtils) savePerformer(checksum string, performer *jsonschema.Performer) error {
	return jsonschema.SavePerformerFile(jp.json.PerformerJSONPath(checksum), performer)
}

func (jp *jsonUtils) getStudio(checksum string) (*jsonschema.Studio, error) {
	return jsonschema.LoadStudioFile(jp.json.StudioJSONPath(checksum))
}

func (jp *jsonUtils) saveStudio(checksum string, studio *jsonschema.Studio) error {
	return jsonschema.SaveStudioFile(jp.json.StudioJSONPath(checksum), studio)
}

func (jp *jsonUtils) getTag(checksum string) (*jsonschema.Tag, error) {
	return jsonschema.LoadTagFile(jp.json.TagJSONPath(checksum))
}

func (jp *jsonUtils) saveTag(checksum string, tag *jsonschema.Tag) error {
	return jsonschema.SaveTagFile(jp.json.TagJSONPath(checksum), tag)
}

func (jp *jsonUtils) getMovie(checksum string) (*jsonschema.Movie, error) {
	return jsonschema.LoadMovieFile(jp.json.MovieJSONPath(checksum))
}

func (jp *jsonUtils) saveMovie(checksum string, movie *jsonschema.Movie) error {
	return jsonschema.SaveMovieFile(jp.json.MovieJSONPath(checksum), movie)
}

func (jp *jsonUtils) getScene(checksum string) (*jsonschema.Scene, error) {
	return jsonschema.LoadSceneFile(jp.json.SceneJSONPath(checksum))
}

func (jp *jsonUtils) saveScene(checksum string, scene *jsonschema.Scene) error {
	return jsonschema.SaveSceneFile(jp.json.SceneJSONPath(checksum), scene)
}

func (jp *jsonUtils) getImage(checksum string) (*jsonschema.Image, error) {
	return jsonschema.LoadImageFile(jp.json.ImageJSONPath(checksum))
}

func (jp *jsonUtils) saveImage(checksum string, image *jsonschema.Image) error {
	return jsonschema.SaveImageFile(jp.json.ImageJSONPath(checksum), image)
}

func (jp *jsonUtils) getGallery(checksum string) (*jsonschema.Gallery, error) {
	return jsonschema.LoadGalleryFile(jp.json.GalleryJSONPath(checksum))
}

func (jp *jsonUtils) saveGallery(checksum string, gallery *jsonschema.Gallery) error {
	return jsonschema.SaveGalleryFile(jp.json.GalleryJSONPath(checksum), gallery)
}
