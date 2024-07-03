package models

import (
	"fmt"
	"strconv"
)

type GroupsScenes struct {
	GroupID int `json:"movie_id"`
	// SceneID    int  `json:"scene_id"`
	SceneIndex *int `json:"scene_index"`
}

func (s GroupsScenes) SceneMovieInput() SceneMovieInput {
	return SceneMovieInput{
		MovieID:    strconv.Itoa(s.GroupID),
		SceneIndex: s.SceneIndex,
	}
}

func (s GroupsScenes) Equal(o GroupsScenes) bool {
	return o.GroupID == s.GroupID && ((o.SceneIndex == nil && s.SceneIndex == nil) ||
		(o.SceneIndex != nil && s.SceneIndex != nil && *o.SceneIndex == *s.SceneIndex))
}

type UpdateGroupIDs struct {
	Groups []GroupsScenes         `json:"movies"`
	Mode   RelationshipUpdateMode `json:"mode"`
}

func (u *UpdateGroupIDs) SceneMovieInputs() []SceneMovieInput {
	if u == nil {
		return nil
	}

	ret := make([]SceneMovieInput, len(u.Groups))
	for _, id := range u.Groups {
		ret = append(ret, id.SceneMovieInput())
	}

	return ret
}

func (u *UpdateGroupIDs) AddUnique(v GroupsScenes) {
	for _, vv := range u.Groups {
		if vv.GroupID == v.GroupID {
			return
		}
	}

	u.Groups = append(u.Groups, v)
}

func GroupsScenesFromInput(input []SceneMovieInput) ([]GroupsScenes, error) {
	ret := make([]GroupsScenes, len(input))

	for i, v := range input {
		mID, err := strconv.Atoi(v.MovieID)
		if err != nil {
			return nil, fmt.Errorf("invalid movie ID: %s", v.MovieID)
		}

		ret[i] = GroupsScenes{
			GroupID:    mID,
			SceneIndex: v.SceneIndex,
		}
	}

	return ret, nil
}
