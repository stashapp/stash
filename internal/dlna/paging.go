package dlna

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
)

type scenePager struct {
	sceneFilter *models.SceneFilterType
	parentID    string
}

func (p *scenePager) getPageID(page int) string {
	return p.parentID + "/page/" + strconv.Itoa(page)
}

func (p *scenePager) getPages(ctx context.Context, r scene.Queryer, total int) ([]interface{}, error) {
	var objs []interface{}

	// get the first scene of each page to set an appropriate title
	pages := int(math.Ceil(float64(total) / float64(pageSize)))

	singlePageSize := 1
	sort := "title"
	findFilter := &models.FindFilterType{
		PerPage: &singlePageSize,
		Sort:    &sort,
	}

	for page := 1; page <= pages; page++ {
		// TODO - this is really slow. Not sure if there's a better way
		title := fmt.Sprintf("Page %d", page)
		if pages <= 10 || (page-1)%(pages/10) == 0 {
			thisPage := ((page - 1) * pageSize) + 1
			findFilter.Page = &thisPage
			scenes, err := scene.Query(ctx, r, p.sceneFilter, findFilter)
			if err != nil {
				return nil, err
			}

			sceneTitle := scenes[0].GetTitle()

			// use the first three letters as a prefix
			if len(sceneTitle) > 3 {
				sceneTitle = sceneTitle[0:3]
			}

			title += fmt.Sprintf(" (%s...)", sceneTitle)
		}

		objs = append(objs, makeStorageFolder(p.getPageID(page), title, p.parentID))
	}

	return objs, nil
}

func (p *scenePager) getPageVideos(ctx context.Context, r SceneFinder, page int, host string) ([]interface{}, error) {
	var objs []interface{}

	sort := "title"
	findFilter := &models.FindFilterType{
		PerPage: &pageSize,
		Page:    &page,
		Sort:    &sort,
	}

	scenes, err := scene.Query(ctx, r, p.sceneFilter, findFilter)
	if err != nil {
		return nil, err
	}

	for _, s := range scenes {
		objs = append(objs, sceneToContainer(s, p.parentID, host))
	}

	return objs, nil
}
