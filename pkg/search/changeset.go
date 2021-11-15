package search

import (
	"fmt"
	"strings"

	"github.com/stashapp/stash/pkg/event"
)

// changeSet is a rollup structure for changes. These are handed off to
// a batch processor when it requests them.
type changeSet struct {
	performers map[int]struct{}
	scenes     map[int]struct{}
	studios    map[int]struct{}
	tags       map[int]struct{}
}

// newChangemap creates a new initialized empty changeMap.
func newChangeSet() *changeSet {
	return &changeSet{
		performers: make(map[int]struct{}),
		scenes:     make(map[int]struct{}),
		studios:    make(map[int]struct{}),
		tags:       make(map[int]struct{}),
	}
}

// track records the given change to the changeMap.
func (s *changeSet) track(e event.Change) {
	switch e.Type {
	case event.Performer:
		s.performers[e.ID] = struct{}{}
	case event.Scene:
		s.scenes[e.ID] = struct{}{}
	case event.Studio:
		s.studios[e.ID] = struct{}{}
	case event.Tag:
		s.tags[e.ID] = struct{}{}
	default:
		// Ignore changes we don't currently track
	}
}

// cutSceneIds returns a slice of sceneIds from the changeSet.
// The limit argument sets an upper bound on the number of elements
// in the slice. Returns if there's any more elements to cut as well.
func (s *changeSet) cutSceneIds(limit int) ([]int, bool) {
	var ret []int
	for k := range s.scenes {
		if limit == 0 {
			return ret, true
		}

		ret = append(ret, k)
		delete(s.scenes, k)
		limit--
	}

	return ret, false
}

func (s *changeSet) performerIds() []int {
	var ret []int
	for k := range s.performers {
		ret = append(ret, k)
	}

	return ret
}

func (s *changeSet) tagIds() []int {
	var ret []int
	for k := range s.tags {
		ret = append(ret, k)
	}

	return ret
}

func (s *changeSet) studioIds() []int {
	var ret []int
	for k := range s.studios {
		ret = append(ret, k)
	}

	return ret
}

// hasContent returns true if there are changes to process.
func (s *changeSet) hasContent() bool {
	return len(s.scenes) > 0 || len(s.performers) > 0 || len(s.tags) > 0 || len(s.studios) > 0
}

// String implements the Stringer interface for changeMaps.
func (s *changeSet) String() string {
	var elems []string
	if len(s.scenes) > 0 {
		elems = append(elems, fmt.Sprintf("(%d scenes)", len(s.scenes)))
	}
	if len(s.performers) > 0 {
		elems = append(elems, fmt.Sprintf("(%d performers)", len(s.performers)))
	}
	if len(s.tags) > 0 {
		elems = append(elems, fmt.Sprintf("(%d tags)", len(s.tags)))
	}
	if len(s.studios) > 0 {
		elems = append(elems, fmt.Sprintf("(%d studios)", len(s.studios)))
	}

	if len(elems) == 0 {
		return "empty changeset"
	}

	return strings.Join(elems, ", ")
}
