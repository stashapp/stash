// +build integration

package models_test

import (
	"strings"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestMarkerFindBySceneMarkerID(t *testing.T) {
	tqb := models.NewTagQueryBuilder()

	markerID := markerIDs[markerIdxWithScene]

	tags, err := tqb.FindBySceneMarkerID(markerID, nil)

	if err != nil {
		t.Fatalf("Error finding tags: %s", err.Error())
	}

	assert.Len(t, tags, 1)
	assert.Equal(t, tagIDs[tagIdxWithMarker], tags[0].ID)

	tags, err = tqb.FindBySceneMarkerID(0, nil)

	if err != nil {
		t.Fatalf("Error finding tags: %s", err.Error())
	}

	assert.Len(t, tags, 0)
}

func TestTagFindByName(t *testing.T) {

	tqb := models.NewTagQueryBuilder()

	name := tagNames[tagIdxWithScene] // find a tag by name

	tag, err := tqb.FindByName(name, nil, false)

	if err != nil {
		t.Fatalf("Error finding tags: %s", err.Error())
	}

	assert.Equal(t, tagNames[tagIdxWithScene], tag.Name)

	name = tagNames[tagIdxWithDupName] // find a tag by name nocase

	tag, err = tqb.FindByName(name, nil, true)

	if err != nil {
		t.Fatalf("Error finding tags: %s", err.Error())
	}
	// tagIdxWithDupName and tagIdxWithScene should have similar names ( only diff should be Name vs NaMe)
	//tag.Name should match with tagIdxWithScene since its ID is before tagIdxWithDupName
	assert.Equal(t, tagNames[tagIdxWithScene], tag.Name)
	//tag.Name should match with tagIdxWithDupName if the check is not case sensitive
	assert.Equal(t, strings.ToLower(tagNames[tagIdxWithDupName]), strings.ToLower(tag.Name))

}

func TestTagFindByNames(t *testing.T) {
	var names []string

	tqb := models.NewTagQueryBuilder()

	names = append(names, tagNames[tagIdxWithScene]) // find tags by names

	tags, err := tqb.FindByNames(names, nil, false)
	if err != nil {
		t.Fatalf("Error finding tags: %s", err.Error())
	}
	assert.Len(t, tags, 1)
	assert.Equal(t, tagNames[tagIdxWithScene], tags[0].Name)

	tags, err = tqb.FindByNames(names, nil, true) // find tags by names nocase
	if err != nil {
		t.Fatalf("Error finding tags: %s", err.Error())
	}
	assert.Len(t, tags, 2) // tagIdxWithScene and tagIdxWithDupName
	assert.Equal(t, strings.ToLower(tagNames[tagIdxWithScene]), strings.ToLower(tags[0].Name))
	assert.Equal(t, strings.ToLower(tagNames[tagIdxWithScene]), strings.ToLower(tags[1].Name))

	names = append(names, tagNames[tagIdx1WithScene]) // find tags by names ( 2 names )

	tags, err = tqb.FindByNames(names, nil, false)
	if err != nil {
		t.Fatalf("Error finding tags: %s", err.Error())
	}
	assert.Len(t, tags, 2) // tagIdxWithScene and tagIdx1WithScene
	assert.Equal(t, tagNames[tagIdxWithScene], tags[0].Name)
	assert.Equal(t, tagNames[tagIdx1WithScene], tags[1].Name)

	tags, err = tqb.FindByNames(names, nil, true) // find tags by names ( 2 names nocase)
	if err != nil {
		t.Fatalf("Error finding tags: %s", err.Error())
	}
	assert.Len(t, tags, 4) // tagIdxWithScene and tagIdxWithDupName , tagIdx1WithScene and tagIdx1WithDupName
	assert.Equal(t, tagNames[tagIdxWithScene], tags[0].Name)
	assert.Equal(t, tagNames[tagIdx1WithScene], tags[1].Name)
	assert.Equal(t, tagNames[tagIdx1WithDupName], tags[2].Name)
	assert.Equal(t, tagNames[tagIdxWithDupName], tags[3].Name)

}

// TODO Create
// TODO Update
// TODO Destroy
// TODO Find
// TODO FindBySceneID
// TODO FindBySceneMarkerID
// TODO Count
// TODO All
// TODO AllSlim
// TODO Query
