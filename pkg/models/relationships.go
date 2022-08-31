package models

import (
	"context"

	"github.com/stashapp/stash/pkg/file"
)

type SceneIDLoader interface {
	GetSceneIDs(ctx context.Context, relatedID int) ([]int, error)
}

type GalleryIDLoader interface {
	GetGalleryIDs(ctx context.Context, relatedID int) ([]int, error)
}

type PerformerIDLoader interface {
	GetPerformerIDs(ctx context.Context, relatedID int) ([]int, error)
}

type TagIDLoader interface {
	GetTagIDs(ctx context.Context, relatedID int) ([]int, error)
}

type SceneMovieLoader interface {
	GetMovies(ctx context.Context, id int) ([]MoviesScenes, error)
}

type StashIDLoader interface {
	GetStashIDs(ctx context.Context, relatedID int) ([]StashID, error)
}

type ImageFileLoader interface {
	GetFiles(ctx context.Context, relatedID int) ([]*file.ImageFile, error)
}

// RelatedIDs represents a list of related IDs.
// TODO - this can be made generic
type RelatedIDs struct {
	list []int
}

// NewRelatedIDs returns a loaded RelatedIDs object with the provided IDs.
// Loaded will return true when called on the returned object if the provided slice is not nil.
func NewRelatedIDs(ids []int) RelatedIDs {
	return RelatedIDs{
		list: ids,
	}
}

// Loaded returns true if the related IDs have been loaded.
func (r RelatedIDs) Loaded() bool {
	return r.list != nil
}

func (r RelatedIDs) mustLoaded() {
	if !r.Loaded() {
		panic("list has not been loaded")
	}
}

// List returns the related IDs. Panics if the relationship has not been loaded.
func (r RelatedIDs) List() []int {
	r.mustLoaded()

	return r.list
}

// Add adds the provided ids to the list. Panics if the relationship has not been loaded.
func (r *RelatedIDs) Add(ids ...int) {
	r.mustLoaded()

	r.list = append(r.list, ids...)
}

func (r *RelatedIDs) load(fn func() ([]int, error)) error {
	if r.Loaded() {
		return nil
	}

	ids, err := fn()
	if err != nil {
		return err
	}

	if ids == nil {
		ids = []int{}
	}

	r.list = ids

	return nil
}

// RelatedMovies represents a list of related Movies.
type RelatedMovies struct {
	list []MoviesScenes
}

// NewRelatedMovies returns a loaded RelatedMovies object with the provided movies.
// Loaded will return true when called on the returned object if the provided slice is not nil.
func NewRelatedMovies(list []MoviesScenes) RelatedMovies {
	return RelatedMovies{
		list: list,
	}
}

// Loaded returns true if the relationship has been loaded.
func (r RelatedMovies) Loaded() bool {
	return r.list != nil
}

func (r RelatedMovies) mustLoaded() {
	if !r.Loaded() {
		panic("list has not been loaded")
	}
}

// List returns the related Movies. Panics if the relationship has not been loaded.
func (r RelatedMovies) List() []MoviesScenes {
	r.mustLoaded()

	return r.list
}

// Add adds the provided ids to the list. Panics if the relationship has not been loaded.
func (r *RelatedMovies) Add(movies ...MoviesScenes) {
	r.mustLoaded()

	r.list = append(r.list, movies...)
}

func (r *RelatedMovies) load(fn func() ([]MoviesScenes, error)) error {
	if r.Loaded() {
		return nil
	}

	ids, err := fn()
	if err != nil {
		return err
	}

	if ids == nil {
		ids = []MoviesScenes{}
	}

	r.list = ids

	return nil
}

type RelatedStashIDs struct {
	list []StashID
}

// NewRelatedStashIDs returns a RelatedStashIDs object with the provided ids.
// Loaded will return true when called on the returned object if the provided slice is not nil.
func NewRelatedStashIDs(list []StashID) RelatedStashIDs {
	return RelatedStashIDs{
		list: list,
	}
}

func (r RelatedStashIDs) mustLoaded() {
	if !r.Loaded() {
		panic("list has not been loaded")
	}
}

// Loaded returns true if the relationship has been loaded.
func (r RelatedStashIDs) Loaded() bool {
	return r.list != nil
}

// List returns the related Stash IDs. Panics if the relationship has not been loaded.
func (r RelatedStashIDs) List() []StashID {
	r.mustLoaded()

	return r.list
}

func (r *RelatedStashIDs) load(fn func() ([]StashID, error)) error {
	if r.Loaded() {
		return nil
	}

	ids, err := fn()
	if err != nil {
		return err
	}

	if ids == nil {
		ids = []StashID{}
	}

	r.list = ids

	return nil
}

type RelatedImageFiles struct {
	primaryFile   *file.ImageFile
	files         []*file.ImageFile
	primaryLoaded bool
}

func NewRelatedImageFiles(files []*file.ImageFile) RelatedImageFiles {
	ret := RelatedImageFiles{
		files:         files,
		primaryLoaded: true,
	}

	if len(files) > 0 {
		ret.primaryFile = files[0]
	}

	return ret
}

// Loaded returns true if the relationship has been loaded.
func (r RelatedImageFiles) Loaded() bool {
	return r.files != nil
}

// Loaded returns true if the primary file relationship has been loaded.
func (r RelatedImageFiles) PrimaryLoaded() bool {
	return r.primaryLoaded
}

// List returns the related files. Panics if the relationship has not been loaded.
func (r RelatedImageFiles) List() []*file.ImageFile {
	if !r.Loaded() {
		panic("relationship has not been loaded")
	}

	return r.files
}

// Primary returns the primary file. Panics if the relationship has not been loaded.
func (r RelatedImageFiles) Primary() *file.ImageFile {
	if !r.PrimaryLoaded() {
		panic("relationship has not been loaded")
	}

	return r.primaryFile
}

func (r *RelatedImageFiles) load(fn func() ([]*file.ImageFile, error)) error {
	if r.Loaded() {
		return nil
	}

	var err error
	r.files, err = fn()
	if err != nil {
		return err
	}

	if len(r.files) > 0 {
		r.primaryFile = r.files[0]
	}

	r.primaryLoaded = true

	return nil
}

func (r *RelatedImageFiles) loadPrimary(fn func() (*file.ImageFile, error)) error {
	if r.PrimaryLoaded() {
		return nil
	}

	var err error
	r.primaryFile, err = fn()
	if err != nil {
		return err
	}

	r.primaryLoaded = true

	return nil
}
