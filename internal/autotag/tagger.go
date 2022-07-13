// Package autotag provides methods to auto-tag scenes with performers,
// studios and tags.
//
// The autotag engine tags scenes with performers/studios/tags if the scene's
// path matches the performer/studio/tag name. A scene's path is considered
// a match if it contains the performer/studio/tag's full name, ignoring any
// '.', '-', '_' characters in the path.
//
// For example, for a performer "foo bar", the following paths would be
// considered a match: "foo bar.mp4", "foobar.mp4", "foo.bar.mp4",
// "foo-bar.mp4", "aaa.foo bar.bbb.mp4".
// The following would not be considered a match:
// "aafoo bar.mp4", "foo barbb.mp4", "foo/bar.mp4"
package autotag

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
)

type tagger struct {
	ID      int
	Type    string
	Name    string
	Path    string
	trimExt bool

	cache *match.Cache
}

type addLinkFunc func(subjectID, otherID int) (bool, error)
type addImageLinkFunc func(o *models.Image) (bool, error)
type addGalleryLinkFunc func(o *models.Gallery) (bool, error)
type addSceneLinkFunc func(o *models.Scene) (bool, error)

func (t *tagger) addError(otherType, otherName string, err error) error {
	return fmt.Errorf("error adding %s '%s' to %s '%s': %s", otherType, otherName, t.Type, t.Name, err.Error())
}

func (t *tagger) addLog(otherType, otherName string) {
	logger.Infof("Added %s '%s' to %s '%s'", otherType, otherName, t.Type, t.Name)
}

func (t *tagger) tagPerformers(ctx context.Context, performerReader match.PerformerAutoTagQueryer, addFunc addLinkFunc) error {
	others, err := match.PathToPerformers(ctx, t.Path, performerReader, t.cache, t.trimExt)
	if err != nil {
		return err
	}

	for _, p := range others {
		added, err := addFunc(t.ID, p.ID)

		if err != nil {
			return t.addError("performer", p.Name.String, err)
		}

		if added {
			t.addLog("performer", p.Name.String)
		}
	}

	return nil
}

func (t *tagger) tagStudios(ctx context.Context, studioReader match.StudioAutoTagQueryer, addFunc addLinkFunc) error {
	studio, err := match.PathToStudio(ctx, t.Path, studioReader, t.cache, t.trimExt)
	if err != nil {
		return err
	}

	if studio != nil {
		added, err := addFunc(t.ID, studio.ID)

		if err != nil {
			return t.addError("studio", studio.Name.String, err)
		}

		if added {
			t.addLog("studio", studio.Name.String)
		}
	}

	return nil
}

func (t *tagger) tagTags(ctx context.Context, tagReader match.TagAutoTagQueryer, addFunc addLinkFunc) error {
	others, err := match.PathToTags(ctx, t.Path, tagReader, t.cache, t.trimExt)
	if err != nil {
		return err
	}

	for _, p := range others {
		added, err := addFunc(t.ID, p.ID)

		if err != nil {
			return t.addError("tag", p.Name, err)
		}

		if added {
			t.addLog("tag", p.Name)
		}
	}

	return nil
}

func (t *tagger) tagScenes(ctx context.Context, paths []string, sceneReader scene.Queryer, addFunc addSceneLinkFunc) error {
	others, err := match.PathToScenes(ctx, t.Name, paths, sceneReader)
	if err != nil {
		return err
	}

	for _, p := range others {
		added, err := addFunc(p)

		if err != nil {
			return t.addError("scene", p.GetTitle(), err)
		}

		if added {
			t.addLog("scene", p.GetTitle())
		}
	}

	return nil
}

func (t *tagger) tagImages(ctx context.Context, paths []string, imageReader image.Queryer, addFunc addImageLinkFunc) error {
	others, err := match.PathToImages(ctx, t.Name, paths, imageReader)
	if err != nil {
		return err
	}

	for _, p := range others {
		added, err := addFunc(p)

		if err != nil {
			return t.addError("image", p.GetTitle(), err)
		}

		if added {
			t.addLog("image", p.GetTitle())
		}
	}

	return nil
}

func (t *tagger) tagGalleries(ctx context.Context, paths []string, galleryReader gallery.Queryer, addFunc addGalleryLinkFunc) error {
	others, err := match.PathToGalleries(ctx, t.Name, paths, galleryReader)
	if err != nil {
		return err
	}

	for _, p := range others {
		added, err := addFunc(p)

		if err != nil {
			return t.addError("gallery", p.GetTitle(), err)
		}

		if added {
			t.addLog("gallery", p.GetTitle())
		}
	}

	return nil
}
