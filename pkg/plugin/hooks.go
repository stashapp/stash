package plugin

import (
	"github.com/stashapp/stash/pkg/plugin/common"
)

type HookTypeEnum string

const (
	SceneMarkerCreatePost HookTypeEnum = "SceneMarker.Create.Post"
	SceneMarkerUpdatePost HookTypeEnum = "SceneMarker.Update.Post"
)

var AllHookTypeEnum = []HookTypeEnum{
	SceneMarkerCreatePost,
	SceneMarkerUpdatePost,
}

func (e HookTypeEnum) IsValid() bool {
	switch e {
	case SceneMarkerCreatePost, SceneMarkerUpdatePost:
		return true
	}
	return false
}

func (e HookTypeEnum) String() string {
	return string(e)
}

func addHookContext(argsMap common.ArgsMap, hookType HookTypeEnum, input interface{}) {
	argsMap[common.HookContextKey] = common.HookContext{
		Type:  string(hookType),
		Input: input,
	}
}
