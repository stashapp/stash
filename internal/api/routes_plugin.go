package api

import (
	"context"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/stashapp/stash/pkg/plugin"
)

type pluginRoutes struct {
	pluginCache *plugin.Cache
}

func getPluginRoutes(pluginCache *plugin.Cache) chi.Router {
	return pluginRoutes{
		pluginCache: pluginCache,
	}.Routes()
}

func (rs pluginRoutes) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/{pluginId}", func(r chi.Router) {
		r.Use(rs.PluginCtx)
		r.Get("/assets/*", rs.Assets)
	})

	return r
}

func (rs pluginRoutes) Assets(w http.ResponseWriter, r *http.Request) {
	p := r.Context().Value(pluginKey).(*plugin.Plugin)

	prefix := "/plugin/" + chi.URLParam(r, "pluginId") + "/assets"

	r.URL.Path = strings.Replace(r.URL.Path, prefix, "", 1)

	// http.FileServer redirects to / if the path ends with index.html
	r.URL.Path = strings.TrimSuffix(r.URL.Path, "/index.html")

	// map the path to the applicable filesystem location
	dir := filepath.Dir(p.ConfigPath)

	// ensure that the path is allowed to be served
	if !rs.canServe(p, r.URL.Path) {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	http.FileServer(http.Dir(dir)).ServeHTTP(w, r)
}

func (rs pluginRoutes) canServe(plugin *plugin.Plugin, path string) bool {
	if path == "" {
		path = "index.html"
	}

	path = strings.TrimPrefix(path, "/")

	for _, p := range plugin.UI.Assets {
		// ignore errors
		if matched, _ := filepath.Match(p, filepath.FromSlash(path)); matched {
			return true
		}
	}
	return false
}

func (rs pluginRoutes) PluginCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := rs.pluginCache.GetPlugin(chi.URLParam(r, "pluginId"))
		if p == nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), pluginKey, p)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
