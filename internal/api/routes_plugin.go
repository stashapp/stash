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

func (rs pluginRoutes) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/{pluginId}", func(r chi.Router) {
		r.Use(rs.PluginCtx)
		r.Get("/assets", rs.Assets)
		r.Get("/assets/*", rs.Assets)
		r.Get("/javascript", rs.Javascript)
		r.Get("/css", rs.CSS)
	})

	return r
}

func (rs pluginRoutes) Assets(w http.ResponseWriter, r *http.Request) {
	p := r.Context().Value(pluginKey).(*plugin.Plugin)

	if !p.Enabled {
		http.Error(w, "plugin disabled", http.StatusBadRequest)
		return
	}

	prefix := "/plugin/" + chi.URLParam(r, "pluginId") + "/assets"

	r.URL.Path = strings.Replace(r.URL.Path, prefix, "", 1)

	// http.FileServer redirects to / if the path ends with index.html
	r.URL.Path = strings.TrimSuffix(r.URL.Path, "/index.html")

	pluginDir := filepath.Dir(p.ConfigPath)

	// map the path to the applicable filesystem location
	var dir string
	r.URL.Path, dir = p.UI.Assets.GetFilesystemLocation(r.URL.Path)
	if dir == "" {
		http.NotFound(w, r)
		return
	}

	dir = filepath.Join(pluginDir, filepath.FromSlash(dir))

	// ensure directory is still within the plugin directory
	if !strings.HasPrefix(dir, pluginDir) {
		http.NotFound(w, r)
		return
	}

	http.FileServer(http.Dir(dir)).ServeHTTP(w, r)
}

func (rs pluginRoutes) Javascript(w http.ResponseWriter, r *http.Request) {
	p := r.Context().Value(pluginKey).(*plugin.Plugin)

	if !p.Enabled {
		http.Error(w, "plugin disabled", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/javascript")
	serveFiles(w, r, p.UI.Javascript)
}

func (rs pluginRoutes) CSS(w http.ResponseWriter, r *http.Request) {
	p := r.Context().Value(pluginKey).(*plugin.Plugin)

	if !p.Enabled {
		http.Error(w, "plugin disabled", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/css")
	serveFiles(w, r, p.UI.CSS)
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
