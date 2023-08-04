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

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/txn"
)

type Tagger struct {
	TxnManager txn.Manager
	Cache      *match.Cache
}

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

func (t *tagger) tagPerformers(ctx context.Context, performerReader models.PerformerAutoTagQueryer, addFunc addLinkFunc) error {
	others, err := match.PathToPerformers(ctx, t.Path, performerReader, t.cache, t.trimExt)
	if err != nil {
		return err
	}

	for _, p := range others {
		added, err := addFunc(t.ID, p.ID)

		if err != nil {
			return t.addError("performer", p.Name, err)
		}

		if added {
			t.addLog("performer", p.Name)
		}
	}

	return nil
}

func (t *tagger) tagStudios(ctx context.Context, studioReader models.StudioAutoTagQueryer, addFunc addLinkFunc) error {
	studio, err := match.PathToStudio(ctx, t.Path, studioReader, t.cache, t.trimExt)
	if err != nil {
		return err
	}

	if studio != nil {
		added, err := addFunc(t.ID, studio.ID)

		if err != nil {
			return t.addError("studio", studio.Name, err)
		}

		if added {
			t.addLog("studio", studio.Name)
		}
	}

	return nil
}

func (t *tagger) tagTags(ctx context.Context, tagReader models.TagAutoTagQueryer, addFunc addLinkFunc) error {
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

func (t *tagger) tagScenes(ctx context.Context, paths []string, sceneReader models.SceneQueryer, addFunc addSceneLinkFunc) error {
	return match.PathToScenesFn(ctx, t.Name, paths, sceneReader, func(ctx context.Context, p *models.Scene) error {
		added, err := addFunc(p)

		if err != nil {
			return t.addError("scene", p.DisplayName(), err)
		}

		if added {
			t.addLog("scene", p.DisplayName())
		}

		return nil
	})
}

func (t *tagger) tagImages(ctx context.Context, paths []string, imageReader models.ImageQueryer, addFunc addImageLinkFunc) error {
	return match.PathToImagesFn(ctx, t.Name, paths, imageReader, func(ctx context.Context, p *models.Image) error {
		added, err := addFunc(p)

		if err != nil {
			return t.addError("image", p.DisplayName(), err)
		}

		if added {
			t.addLog("image", p.DisplayName())
		}

		return nil
	})
}

func (t *tagger) tagGalleries(ctx context.Context, paths []string, galleryReader models.GalleryQueryer, addFunc addGalleryLinkFunc) error {
	return match.PathToGalleriesFn(ctx, t.Name, paths, galleryReader, func(ctx context.Context, p *models.Gallery) error {
		added, err := addFunc(p)

		if err != nil {
			return t.addError("gallery", p.DisplayName(), err)
		}

		if added {
			t.addLog("gallery", p.DisplayName())
		}

		return nil
	})
}
