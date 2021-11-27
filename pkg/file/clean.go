package file

import (
	"errors"
	"fmt"
	"io/fs"
	"regexp"
	"strings"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type CleanConfiguration interface {
	GetStashFromPath(pathToCheck string) *models.StashConfig
	GetGeneratedPath() string
	GetVideoExtensions() []string
	GetGalleryExtensions() []string
	GetImageExtensions() []string

	GetExcludes() []string
	GetImageExcludes() []string
}

type Cleaner struct {
	FileReader models.FileReader
	Config     CleanConfiguration
}

// GetFilesToClean returns the subset of the provided files that should be removed.
// Returns nil and an error if an error occurs.
func (c *Cleaner) GetFilesToClean(fileIDs []int, relatedObj interface{}) ([]int, error) {
	var ret []int

	files, err := c.FileReader.Find(fileIDs)
	if err != nil {
		return nil, fmt.Errorf("error getting files: %w", err)
	}

	for _, f := range files {
		if c.shouldClean(f, relatedObj) {
			ret = append(ret, f.ID)
		}
	}

	return ret, nil
}

func (c *Cleaner) shouldClean(f *models.File, relatedObj interface{}) bool {
	// #1102 - clean anything in generated path
	generatedPath := c.Config.GetGeneratedPath()
	pathStash := c.Config.GetStashFromPath(f.Path)
	if pathStash == nil || utils.IsPathInDir(generatedPath, f.Path) {
		logger.Infof("File not found. Marking to clean: %s", f.Path)
		return true
	}

	// ensure exists
	if _, err := Info(f); errors.Is(err, fs.ErrNotExist) {
		logger.Infof("File not found. Marking to clean: %s", f.Path)
		return true
	}

	var extensions []string
	var excludes []string
	isZip := false
	var relatedType string

	switch relatedObj.(type) {
	case *models.Scene:
		if pathStash.ExcludeVideo {
			logger.Infof("File in stash library that excludes video. Marking to clean: %s", f.Path)
			return true
		}

		relatedType = "video"
		extensions = c.Config.GetVideoExtensions()
		excludes = c.Config.GetExcludes()

	case *models.Image:
		if pathStash.ExcludeImage {
			logger.Infof("File in stash library that excludes images. Marking to clean: %s", f.Path)
			return true
		}

		relatedType = "image"
		extensions = c.Config.GetImageExtensions()
		excludes = c.Config.GetImageExcludes()

	case *models.Gallery:
		if pathStash.ExcludeImage {
			logger.Infof("File in stash library that excludes images. Marking to clean: %s", f.Path)
			return true
		}

		relatedType = "gallery"
		isZip = true
		extensions = c.Config.GetGalleryExtensions()
		excludes = c.Config.GetImageExcludes()

	default:
		panic("unexpected relatedObj type")
	}

	if !utils.MatchExtension(f.Path, extensions) {
		logger.Infof("File extension does not match %s extensions. Marking to clean: %s", relatedType, f.Path)
		return true
	}

	if matchFile(f.Path, excludes) {
		logger.Infof("File matched exclude regex. Marking to clean: %s", f.Path)
		return true
	}

	if isZip {
		imageExtensions := c.Config.GetImageExtensions()

		n, err := CountImagesInZip(f.Path, imageExtensions)
		if err != nil {
			logger.Errorf("Error counting images in zip %q: %v", f.Path, err)
			return false
		}

		if n == 0 {
			logger.Infof("Gallery has 0 images. Marking to clean: %s", f.Path)
			return true
		}
	}

	return false
}

func matchFileRegex(file string, fileRegexps []*regexp.Regexp) bool {
	for _, regPattern := range fileRegexps {
		if regPattern.MatchString(strings.ToLower(file)) {
			return true
		}
	}
	return false
}

func matchFile(file string, patterns []string) bool {
	if patterns != nil {
		fileRegexps := generateRegexps(patterns)

		return matchFileRegex(file, fileRegexps)
	}

	return false
}

func generateRegexps(patterns []string) []*regexp.Regexp {
	var fileRegexps []*regexp.Regexp

	for _, pattern := range patterns {
		reg, err := regexp.Compile(strings.ToLower(pattern))
		if err != nil {
			logger.Errorf("Exclude :%v", err)
		} else {
			fileRegexps = append(fileRegexps, reg)
		}
	}

	if len(fileRegexps) == 0 {
		return nil
	} else {
		return fileRegexps
	}

}
