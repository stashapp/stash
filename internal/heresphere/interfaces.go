package heresphere

import (
	"github.com/stashapp/stash/pkg/scene"
)

type sceneFinder interface {
	scene.Queryer
	scene.IDFinder
}
