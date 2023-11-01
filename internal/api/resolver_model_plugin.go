package api

import (
	"context"

	"github.com/stashapp/stash/pkg/plugin"
)

type pluginURLBuilder struct {
	BaseURL string
	Plugin  *plugin.Plugin
}

func (b pluginURLBuilder) javascript() *string {
	if len(b.Plugin.UI.Javascript) == 0 {
		return nil
	}

	ret := b.BaseURL + "/plugin/" + b.Plugin.ID + "/javascript"
	return &ret
}

func (b pluginURLBuilder) css() *string {
	if len(b.Plugin.UI.CSS) == 0 {
		return nil
	}

	ret := b.BaseURL + "/plugin/" + b.Plugin.ID + "/css"
	return &ret
}

func (b *pluginURLBuilder) paths() *PluginPaths {
	return &PluginPaths{
		Javascript: b.javascript(),
		CSS:        b.css(),
	}
}

func (r *pluginResolver) Paths(ctx context.Context, obj *plugin.Plugin) (*PluginPaths, error) {
	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)

	b := pluginURLBuilder{
		BaseURL: baseURL,
		Plugin:  obj,
	}

	return b.paths(), nil
}
