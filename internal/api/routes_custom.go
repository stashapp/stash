package api

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/stashapp/stash/internal/manager/config"
)

type customRoutes struct {
	servedFolders config.URLMap
}

func getCustomRoutes(servedFolders config.URLMap) chi.Router {
	return customRoutes{servedFolders: servedFolders}.Routes()
}

func (rs customRoutes) Routes() chi.Router {
	r := chi.NewRouter()

	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.Replace(r.URL.Path, "/custom", "", 1)

		// http.FileServer redirects to / if the path ends with index.html
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/index.html")

		// map the path to the applicable filesystem location
		var dir string
		r.URL.Path, dir = rs.servedFolders.GetFilesystemLocation(r.URL.Path)
		if dir != "" {
			http.FileServer(http.Dir(dir)).ServeHTTP(w, r)
		} else {
			http.NotFound(w, r)
		}
	})

	return r
}
