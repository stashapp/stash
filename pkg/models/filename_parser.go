package models

type SceneParserInput struct {
	IgnoreWords          []string `json:"ignoreWords"`
	WhitespaceCharacters *string  `json:"whitespaceCharacters"`
	CapitalizeTitle      *bool    `json:"capitalizeTitle"`
	IgnoreOrganized      *bool    `json:"ignoreOrganized"`
}

type SceneParserResult struct {
	Scene        *Scene          `json:"scene"`
	Title        *string         `json:"title"`
	Code         *string         `json:"code"`
	Details      *string         `json:"details"`
	Director     *string         `json:"director"`
	URL          *string         `json:"url"`
	Date         *string         `json:"date"`
	Rating       *int            `json:"rating"`
	Rating100    *int            `json:"rating100"`
	StudioID     *string         `json:"studio_id"`
	GalleryIds   []string        `json:"gallery_ids"`
	PerformerIds []string        `json:"performer_ids"`
	Movies       []*SceneMovieID `json:"movies"`
	TagIds       []string        `json:"tag_ids"`
}

type SceneMovieID struct {
	MovieID    string  `json:"movie_id"`
	SceneIndex *string `json:"scene_index"`
}

type AudioParserInput struct {
	IgnoreWords          []string `json:"ignoreWords"`
	WhitespaceCharacters *string  `json:"whitespaceCharacters"`
	CapitalizeTitle      *bool    `json:"capitalizeTitle"`
	IgnoreOrganized      *bool    `json:"ignoreOrganized"`
}

type AudioParserResult struct {
	Audio        *Audio          `json:"audio"`
	Title        *string         `json:"title"`
	Code         *string         `json:"code"`
	Details      *string         `json:"details"`
	URL          *string         `json:"url"`
	Date         *string         `json:"date"`
	Rating100    *int            `json:"rating100"`
	StudioID     *string         `json:"studio_id"`
	GalleryIds   []string        `json:"gallery_ids"`
	PerformerIds []string        `json:"performer_ids"`
	Groups       []*AudioGroupID `json:"groups"`
	TagIds       []string        `json:"tag_ids"`
}

type AudioGroupID struct {
	GroupID    string  `json:"group_id"`
	AudioIndex *string `json:"audio_index"`
}
