package models

import (
	"context"
)

type SceneIDLoader interface {
	GetSceneIDs(ctx context.Context, relatedID int) ([]int, error)
}

type ImageIDLoader interface {
	GetImageIDs(ctx context.Context, relatedID int) ([]int, error)
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

type TagRelationLoader interface {
	GetParentIDs(ctx context.Context, relatedID int) ([]int, error)
	GetChildIDs(ctx context.Context, relatedID int) ([]int, error)
}

type FileIDLoader interface {
	GetManyFileIDs(ctx context.Context, ids []int) ([][]FileID, error)
}

type SceneGroupLoader interface {
	GetGroups(ctx context.Context, id int) ([]GroupsScenes, error)
}

type StashIDLoader interface {
	GetStashIDs(ctx context.Context, relatedID int) ([]StashID, error)
}

type VideoFileLoader interface {
	GetFiles(ctx context.Context, relatedID int) ([]*VideoFile, error)
}

type FileLoader interface {
	GetFiles(ctx context.Context, relatedID int) ([]File, error)
}

type AliasLoader interface {
	GetAliases(ctx context.Context, relatedID int) ([]string, error)
}

type URLLoader interface {
	GetURLs(ctx context.Context, relatedID int) ([]string, error)
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

// RelatedGroups represents a list of related Groups.
type RelatedGroups struct {
	list []GroupsScenes
}

// NewRelatedGroups returns a loaded RelateGroups object with the provided groups.
// Loaded will return true when called on the returned object if the provided slice is not nil.
func NewRelatedGroups(list []GroupsScenes) RelatedGroups {
	return RelatedGroups{
		list: list,
	}
}

// Loaded returns true if the relationship has been loaded.
func (r RelatedGroups) Loaded() bool {
	return r.list != nil
}

func (r RelatedGroups) mustLoaded() {
	if !r.Loaded() {
		panic("list has not been loaded")
	}
}

// List returns the related Groups. Panics if the relationship has not been loaded.
func (r RelatedGroups) List() []GroupsScenes {
	r.mustLoaded()

	return r.list
}

// Add adds the provided ids to the list. Panics if the relationship has not been loaded.
func (r *RelatedGroups) Add(groups ...GroupsScenes) {
	r.mustLoaded()

	r.list = append(r.list, groups...)
}

// ForID returns the GroupsScenes object for the given group ID. Returns nil if not found.
func (r *RelatedGroups) ForID(id int) *GroupsScenes {
	r.mustLoaded()

	for _, v := range r.list {
		if v.GroupID == id {
			return &v
		}
	}

	return nil
}

func (r *RelatedGroups) load(fn func() ([]GroupsScenes, error)) error {
	if r.Loaded() {
		return nil
	}

	ids, err := fn()
	if err != nil {
		return err
	}

	if ids == nil {
		ids = []GroupsScenes{}
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

// ForID returns the StashID object for the given endpoint. Returns nil if not found.
func (r *RelatedStashIDs) ForEndpoint(endpoint string) *StashID {
	r.mustLoaded()

	for _, v := range r.list {
		if v.Endpoint == endpoint {
			return &v
		}
	}

	return nil
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

type RelatedVideoFiles struct {
	primaryFile   *VideoFile
	files         []*VideoFile
	primaryLoaded bool
}

func NewRelatedVideoFiles(files []*VideoFile) RelatedVideoFiles {
	ret := RelatedVideoFiles{
		files:         files,
		primaryLoaded: true,
	}

	if len(files) > 0 {
		ret.primaryFile = files[0]
	}

	return ret
}

func (r *RelatedVideoFiles) SetPrimary(f *VideoFile) {
	r.primaryFile = f
	r.primaryLoaded = true
}

func (r *RelatedVideoFiles) Set(f []*VideoFile) {
	r.files = f
	if len(r.files) > 0 {
		r.primaryFile = r.files[0]
	}

	r.primaryLoaded = true
}

// Loaded returns true if the relationship has been loaded.
func (r RelatedVideoFiles) Loaded() bool {
	return r.files != nil
}

// Loaded returns true if the primary file relationship has been loaded.
func (r RelatedVideoFiles) PrimaryLoaded() bool {
	return r.primaryLoaded
}

// List returns the related files. Panics if the relationship has not been loaded.
func (r RelatedVideoFiles) List() []*VideoFile {
	if !r.Loaded() {
		panic("relationship has not been loaded")
	}

	return r.files
}

// Primary returns the primary file. Panics if the relationship has not been loaded.
func (r RelatedVideoFiles) Primary() *VideoFile {
	if !r.PrimaryLoaded() {
		panic("relationship has not been loaded")
	}

	return r.primaryFile
}

func (r *RelatedVideoFiles) load(fn func() ([]*VideoFile, error)) error {
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

func (r *RelatedVideoFiles) loadPrimary(fn func() (*VideoFile, error)) error {
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

type RelatedFiles struct {
	primaryFile   File
	files         []File
	primaryLoaded bool
}

func NewRelatedFiles(files []File) RelatedFiles {
	ret := RelatedFiles{
		files:         files,
		primaryLoaded: true,
	}

	if len(files) > 0 {
		ret.primaryFile = files[0]
	}

	return ret
}

// Loaded returns true if the relationship has been loaded.
func (r RelatedFiles) Loaded() bool {
	return r.files != nil
}

// Loaded returns true if the primary file relationship has been loaded.
func (r RelatedFiles) PrimaryLoaded() bool {
	return r.primaryLoaded
}

// List returns the related files. Panics if the relationship has not been loaded.
func (r RelatedFiles) List() []File {
	if !r.Loaded() {
		panic("relationship has not been loaded")
	}

	return r.files
}

// Primary returns the primary file. Panics if the relationship has not been loaded.
func (r RelatedFiles) Primary() File {
	if !r.PrimaryLoaded() {
		panic("relationship has not been loaded")
	}

	return r.primaryFile
}

func (r *RelatedFiles) load(fn func() ([]File, error)) error {
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

func (r *RelatedFiles) loadPrimary(fn func() (File, error)) error {
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

// RelatedStrings represents a list of related strings.
// TODO - this can be made generic
type RelatedStrings struct {
	list []string
}

// NewRelatedStrings returns a loaded RelatedStrings object with the provided values.
// Loaded will return true when called on the returned object if the provided slice is not nil.
func NewRelatedStrings(values []string) RelatedStrings {
	return RelatedStrings{
		list: values,
	}
}

// Loaded returns true if the related IDs have been loaded.
func (r RelatedStrings) Loaded() bool {
	return r.list != nil
}

func (r RelatedStrings) mustLoaded() {
	if !r.Loaded() {
		panic("list has not been loaded")
	}
}

// List returns the related values. Panics if the relationship has not been loaded.
func (r RelatedStrings) List() []string {
	r.mustLoaded()

	return r.list
}

// Add adds the provided values to the list. Panics if the relationship has not been loaded.
func (r *RelatedStrings) Add(values ...string) {
	r.mustLoaded()

	r.list = append(r.list, values...)
}

func (r *RelatedStrings) load(fn func() ([]string, error)) error {
	if r.Loaded() {
		return nil
	}

	values, err := fn()
	if err != nil {
		return err
	}

	if values == nil {
		values = []string{}
	}

	r.list = values

	return nil
}
