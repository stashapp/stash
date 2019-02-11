package api

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/handler"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gobuffalo/packr/v2"
	"github.com/rs/cors"
	"github.com/stashapp/stash/logger"
	"github.com/stashapp/stash/manager"
	"github.com/stashapp/stash/manager/jsonschema"
	"github.com/stashapp/stash/models"
	"github.com/stashapp/stash/utils"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime/debug"
	"strings"
)

const httpPort = "9998"
const httpsPort = "9999"

var certsBox *packr.Box
var uiBox *packr.Box
var setupUIBox *packr.Box

func Start() {
	//port := os.Getenv("PORT")
	//if port == "" {
	//	port = defaultPort
	//}

	certsBox = packr.New("Cert Box", "../certs")
	uiBox = packr.New("UI Box", "../ui/v1/dist/stash-frontend")
	setupUIBox = packr.New("Setup UI Box", "../ui/setup")

	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.DefaultCompress)
	r.Use(middleware.StripSlashes)
	r.Use(cors.AllowAll().Handler)
	r.Use(BaseURLMiddleware)
	r.Use(ConfigCheckMiddleware)

	recoverFunc := handler.RecoverFunc(func(ctx context.Context, err interface{}) error {
		logger.Error(err)
		debug.PrintStack()

		message := fmt.Sprintf("Internal system error. Error <%v>", err)
		return errors.New(message)
	})
	requestMiddleware := handler.RequestMiddleware(func(ctx context.Context, next func(ctx context.Context) []byte) []byte {
		//api.GetRequestContext(ctx).Variables[]
		return next(ctx)
	})
	gqlHandler := handler.GraphQL(models.NewExecutableSchema(models.Config{Resolvers: &Resolver{}}), recoverFunc, requestMiddleware)

	// https://stash.server:9999/certs/server.crt
	r.Handle("/certs/*", http.FileServer(certsBox))

	r.Handle("/graphql", gqlHandler)
	r.Handle("/playground", handler.Playground("GraphQL playground", "/graphql"))

	r.Mount("/gallery", galleryRoutes{}.Routes())
	r.Mount("/performer", performerRoutes{}.Routes())
	r.Mount("/scene", sceneRoutes{}.Routes())
	r.Mount("/studio", studioRoutes{}.Routes())

	// Serve the setup UI
	r.HandleFunc("/setup*", func(w http.ResponseWriter, r *http.Request) {
		ext := path.Ext(r.URL.Path)
		if ext == ".html" || ext == ""  {
			data := setupUIBox.Bytes("index.html")
			_, _ = w.Write(data)
		} else {
			r.URL.Path = strings.Replace(r.URL.Path, "/setup", "", 1)
			http.FileServer(setupUIBox).ServeHTTP(w, r)
		}
	})
	r.Post("/init", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, fmt.Sprintf("error: %s", err), 500)
		}
		stash := filepath.Clean(r.Form.Get("stash"))
		metadata := filepath.Clean(r.Form.Get("metadata"))
		cache := filepath.Clean(r.Form.Get("cache"))
		//downloads := filepath.Clean(r.Form.Get("downloads")) // TODO
		downloads := filepath.Join(metadata, "downloads")

		exists, _ := utils.FileExists(stash)
		fileInfo, _ := os.Stat(stash)
		if !exists || !fileInfo.IsDir() {
			http.Error(w, fmt.Sprintf("the stash path either doesn't exist, or is not a directory <%s>.  Go back and try again.", stash), 500)
			return
		}

		exists, _ = utils.FileExists(metadata)
		fileInfo, _ = os.Stat(metadata)
		if !exists || !fileInfo.IsDir() {
			http.Error(w, fmt.Sprintf("the metadata path either doesn't exist, or is not a directory <%s>  Go back and try again.", metadata), 500)
			return
		}

		exists, _ = utils.FileExists(cache)
		fileInfo, _ = os.Stat(cache)
		if !exists || !fileInfo.IsDir() {
			http.Error(w, fmt.Sprintf("the cache path either doesn't exist, or is not a directory <%s>  Go back and try again.", cache), 500)
			return
		}

		_ = os.Mkdir(downloads, 0755)

		config := &jsonschema.Config{
			Stash: stash,
			Metadata: metadata,
			Cache: cache,
			Downloads: downloads,
		}
		if err := manager.GetInstance().SaveConfig(config); err != nil {
			http.Error(w, fmt.Sprintf("there was an error saving the config file: %s", err), 500)
			return
		}

		http.Redirect(w, r, "/", 301)
	})

	// Serve the angular app
	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		ext := path.Ext(r.URL.Path)
		if ext == ".html" || ext == ""  {
			data := uiBox.Bytes("index.html")
			_, _ = w.Write(data)
		} else {
			http.FileServer(uiBox).ServeHTTP(w, r)
		}
	})

	httpsServer := &http.Server{
		Addr: ":"+httpsPort,
		Handler: r,
		TLSConfig: makeTLSConfig(),
	}
	server := &http.Server{
		Addr: ":"+httpPort,
		Handler: r,
	}

	go func() {
		logger.Infof("stash is running on HTTP at http://localhost:9998/")
		logger.Fatal(server.ListenAndServe())
	}()

	go func() {
		logger.Infof("stash is running on HTTPS at https://localhost:9999/")
		logger.Fatal(httpsServer.ListenAndServeTLS("", ""))
	}()
}

func makeTLSConfig() *tls.Config {
	cert, err := certsBox.Find("server.crt")
	key, err := certsBox.Find("server.key")

	certs := make([]tls.Certificate, 1)
	certs[0], err = tls.X509KeyPair(cert, key)
	if err != nil {
		panic(err)
	}
	tlsConfig := &tls.Config{
		Certificates: certs,
	}

	return tlsConfig
}


type contextKey struct {
	name string
}
var (
	BaseURLCtxKey = &contextKey{"BaseURL"}
)
func BaseURLMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var scheme string
		if strings.Compare("https", r.URL.Scheme) == 0 || r.Proto == "HTTP/2.0" {
			scheme = "https"
		} else {
			scheme = "http"
		}
		baseURL := scheme + "://" + r.Host

		r = r.WithContext(context.WithValue(ctx, BaseURLCtxKey, baseURL))

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func ConfigCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ext := path.Ext(r.URL.Path)
		shouldRedirect := ext == "" && r.Method == "GET" && r.URL.Path != "/init"
		if !manager.HasValidConfig() && shouldRedirect {
			if !strings.HasPrefix(r.URL.Path, "/setup") {
				http.Redirect(w, r, "/setup", 301)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}