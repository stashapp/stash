//go:build plugin_example
// +build plugin_example

package common

import (
	"context"
	"errors"
	"fmt"

	graphql "github.com/hasura/go-graphql-client"
	"github.com/stashapp/stash/pkg/plugin/common/log"
)

const tagName = "Hawwwwt"

// graphql inputs and returns
type TagCreate struct {
	ID graphql.ID `graphql:"id"`
}

type TagCreateInput struct {
	Name graphql.String `graphql:"name" json:"name"`
}

type TagDestroyInput struct {
	ID graphql.ID `graphql:"id" json:"id"`
}

type FindScenesResultType struct {
	Count           graphql.Int
	DurationSeconds graphql.Float `graphql:"duration" json:"duration"`
	FilesizeBytes   graphql.Float `graphql:"filesize" json:"filesize"`
	Scenes          []Scene
}

type Tag struct {
	ID   graphql.ID     `graphql:"id"`
	Name graphql.String `graphql:"name"`
}

type Scene struct {
	ID   graphql.ID
	Tags []Tag
}

func (s Scene) getTagIds() []graphql.ID {
	ret := []graphql.ID{}

	for _, t := range s.Tags {
		ret = append(ret, t.ID)
	}

	return ret
}

type FindFilterType struct {
	PerPage *graphql.Int    `graphql:"per_page" json:"per_page"`
	Sort    *graphql.String `graphql:"sort" json:"sort"`
}

type SceneUpdate struct {
	ID graphql.ID `graphql:"id"`
}

type SceneUpdateInput struct {
	ID     graphql.ID   `graphql:"id" json:"id"`
	TagIds []graphql.ID `graphql:"tag_ids" json:"tag_ids"`
}

func getTagID(ctx context.Context, client *graphql.Client, create bool) (*graphql.ID, error) {
	log.Info("Checking if tag exists already")

	// see if tag exists already
	var q struct {
		AllTags []Tag `graphql:"allTags"`
	}

	err := client.Query(ctx, &q, nil)
	if err != nil {
		return nil, fmt.Errorf("Error getting tags: %s\n", err.Error())
	}

	for _, t := range q.AllTags {
		if t.Name == tagName {
			id := t.ID
			return &id, nil
		}
	}

	if !create {
		log.Info("Not found and not creating")
		return nil, nil
	}

	// create the tag
	var m struct {
		TagCreate TagCreate `graphql:"tagCreate(input: $s)"`
	}

	input := TagCreateInput{
		Name: tagName,
	}

	vars := map[string]interface{}{
		"s": input,
	}

	log.Info("Creating new tag")

	err = client.Mutate(ctx, &m, vars)
	if err != nil {
		return nil, fmt.Errorf("Error mutating scene: %s\n", err.Error())
	}

	return &m.TagCreate.ID, nil
}

func findRandomScene(ctx context.Context, client *graphql.Client) (*Scene, error) {
	// get a random scene
	var q struct {
		FindScenes FindScenesResultType `graphql:"findScenes(filter: $c)"`
	}

	pp := graphql.Int(1)
	sort := graphql.String("random")
	filterInput := &FindFilterType{
		PerPage: &pp,
		Sort:    &sort,
	}

	vars := map[string]interface{}{
		"c": filterInput,
	}

	log.Info("Finding a random scene")
	err := client.Query(ctx, &q, vars)
	if err != nil {
		return nil, fmt.Errorf("Error getting random scene: %s\n", err.Error())
	}

	if q.FindScenes.Count == 0 {
		return nil, nil
	}

	return &q.FindScenes.Scenes[0], nil
}

func addTagId(tagIds []graphql.ID, tagId graphql.ID) []graphql.ID {
	for _, t := range tagIds {
		if t == tagId {
			return tagIds
		}
	}

	tagIds = append(tagIds, tagId)
	return tagIds
}

func AddTag(ctx context.Context, client *graphql.Client) error {
	tagID, err := getTagID(ctx, client, true)

	if err != nil {
		return err
	}

	scene, err := findRandomScene(ctx, client)

	if err != nil {
		return err
	}

	if scene == nil {
		return errors.New("no scenes to add tag to")
	}

	var m struct {
		SceneUpdate SceneUpdate `graphql:"sceneUpdate(input: $s)"`
	}

	input := SceneUpdateInput{
		ID:     scene.ID,
		TagIds: scene.getTagIds(),
	}

	input.TagIds = addTagId(input.TagIds, *tagID)

	vars := map[string]interface{}{
		"s": input,
	}

	log.Infof("Adding tag to scene %v", scene.ID)
	err = client.Mutate(ctx, &m, vars)
	if err != nil {
		return fmt.Errorf("Error mutating scene: %v", err)
	}

	return nil
}

func RemoveTag(ctx context.Context, client *graphql.Client) error {
	tagID, err := getTagID(ctx, client, false)

	if err != nil {
		return err
	}

	if tagID == nil {
		log.Info("Tag does not exist. Nothing to remove")
		return nil
	}

	// destroy the tag
	var m struct {
		TagDestroy bool `graphql:"tagDestroy(input: $s)"`
	}

	input := TagDestroyInput{
		ID: *tagID,
	}

	vars := map[string]interface{}{
		"s": input,
	}

	log.Info("Destroying tag")

	err = client.Mutate(ctx, &m, vars)
	if err != nil {
		return fmt.Errorf("Error destroying tag: %v", err)
	}

	return nil
}
