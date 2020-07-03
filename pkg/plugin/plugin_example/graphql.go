// +build plugin_example
package main

import (
	"github.com/shurcooL/graphql"
)

type api struct {
	stopping bool
}

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
	Count  graphql.Int
	Scenes []Scene
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
