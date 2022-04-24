package models

import (
	"fmt"
	"strconv"
)

type MoviesScenes struct {
	MovieID int `json:"movie_id"`
	// SceneID    int  `json:"scene_id"`
	SceneIndex *int `json:"scene_index"`
}

func (s MoviesScenes) SceneMovieInput() *SceneMovieInput {
	return &SceneMovieInput{
		MovieID:    strconv.Itoa(s.MovieID),
		SceneIndex: s.SceneIndex,
	}
}

func (s MoviesScenes) Equal(o MoviesScenes) bool {
	return o.MovieID == s.MovieID && ((o.SceneIndex == nil && s.SceneIndex == nil) ||
		(o.SceneIndex != nil && s.SceneIndex != nil && *o.SceneIndex == *s.SceneIndex))
}

type UpdateMovieIDs struct {
	Movies []MoviesScenes         `json:"movies"`
	Mode   RelationshipUpdateMode `json:"mode"`
}

func (u *UpdateMovieIDs) SceneMovieInputs() []*SceneMovieInput {
	if u == nil {
		return nil
	}

	ret := make([]*SceneMovieInput, len(u.Movies))
	for _, id := range u.Movies {
		ret = append(ret, id.SceneMovieInput())
	}

	return ret
}

func UpdateMovieIDsFromInput(i []*SceneMovieInput) (*UpdateMovieIDs, error) {
	ret := &UpdateMovieIDs{
		Mode: RelationshipUpdateModeSet,
	}

	for _, v := range i {
		mID, err := strconv.Atoi(v.MovieID)
		if err != nil {
			return nil, fmt.Errorf("invalid movie ID: %s", v.MovieID)
		}

		ret.Movies = append(ret.Movies, MoviesScenes{
			MovieID:    mID,
			SceneIndex: v.SceneIndex,
		})
	}

	return ret, nil
}
