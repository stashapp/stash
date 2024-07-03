package manager

import (
	"path/filepath"

	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/models/paths"
)

type jsonUtils struct {
	json paths.JSONPaths
}

func (jp *jsonUtils) savePerformer(fn string, performer *jsonschema.Performer) error {
	return jsonschema.SavePerformerFile(filepath.Join(jp.json.Performers, fn), performer)
}

func (jp *jsonUtils) saveStudio(fn string, studio *jsonschema.Studio) error {
	return jsonschema.SaveStudioFile(filepath.Join(jp.json.Studios, fn), studio)
}

func (jp *jsonUtils) saveTag(fn string, tag *jsonschema.Tag) error {
	return jsonschema.SaveTagFile(filepath.Join(jp.json.Tags, fn), tag)
}

func (jp *jsonUtils) saveGroup(fn string, group *jsonschema.Group) error {
	return jsonschema.SaveGroupFile(filepath.Join(jp.json.Groups, fn), group)
}

func (jp *jsonUtils) saveScene(fn string, scene *jsonschema.Scene) error {
	return jsonschema.SaveSceneFile(filepath.Join(jp.json.Scenes, fn), scene)
}

func (jp *jsonUtils) saveImage(fn string, image *jsonschema.Image) error {
	return jsonschema.SaveImageFile(filepath.Join(jp.json.Images, fn), image)
}

func (jp *jsonUtils) saveGallery(fn string, gallery *jsonschema.Gallery) error {
	return jsonschema.SaveGalleryFile(filepath.Join(jp.json.Galleries, fn), gallery)
}

func (jp *jsonUtils) saveFile(fn string, file jsonschema.DirEntry) error {
	return jsonschema.SaveFileFile(filepath.Join(jp.json.Files, fn), file)
}
