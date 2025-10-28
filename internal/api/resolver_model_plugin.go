package api

import (
	"context"

	"github.com/stashapp/stash/pkg/plugin"
)

type pluginURLBuilder struct {
	BaseURL string
	Plugin  *plugin.Plugin
}

func (b pluginURLBuilder) javascript() []string {
	ui := b.Plugin.UI
	if len(ui.Javascript) == 0 && len(ui.ExternalScript) == 0 {
		return nil
	}

	var ret []string

	ret = append(ret, ui.ExternalScript...)
	ret = append(ret, b.BaseURL+"/plugin/"+b.Plugin.ID+"/javascript")

	return ret
}

func (b pluginURLBuilder) css() []string {
	ui := b.Plugin.UI
	if len(ui.CSS) == 0 && len(ui.ExternalCSS) == 0 {
		return nil
	}

	var ret []string

	ret = append(ret, b.Plugin.UI.ExternalCSS...)
	ret = append(ret, b.BaseURL+"/plugin/"+b.Plugin.ID+"/css")
	return ret
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

func (r *pluginResolver) Requires(ctx context.Context, obj *plugin.Plugin) ([]string, error) {
	return obj.UI.Requires, nil
}
