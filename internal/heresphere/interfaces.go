package heresphere

import (
	"github.com/stashapp/stash/pkg/models"
)

type sceneFinder interface {
	models.SceneQueryer
	models.SceneGetter
}
