package api

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/99designs/gqlgen/handler"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/manager/paths"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

var uiBox *packr.Box

//var legacyUiBox *packr.Box
var setupUIBox *packr.Box

func Start() {
	uiBox = packr.New("UI Box", "../../ui/v2/build")
	//legacyUiBox = packr.New("UI Box", "../../ui/v1/dist/stash-frontend")
	setupUIBox = packr.New("Setup UI Box", "../../ui/setup")

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
	websocketUpgrader := handler.WebsocketUpgrader(websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	})
	gqlHandler := handler.GraphQL(models.NewExecutableSchema(models.Config{Resolvers: &Resolver{}}), recoverFunc, requestMiddleware, websocketUpgrader)

	r.Handle("/graphql", gqlHandler)
	r.Handle("/playground", handler.Playground("GraphQL playground", "/graphql"))

	r.Mount("/gallery", galleryRoutes{}.Routes())
	r.Mount("/performer", performerRoutes{}.Routes())
	r.Mount("/scene", sceneRoutes{}.Routes())
	r.Mount("/studio", studioRoutes{}.Routes())

	r.HandleFunc("/css", func(w http.ResponseWriter, r *http.Request) {
		if !config.GetCSSEnabled() {
			return
		}
		
		// search for custom.css in current directory, then $HOME/.stash
		fn := config.GetCSSPath()
		exists, _ := utils.FileExists(fn)
		if !exists {
			return
		}

		http.ServeFile(w, r, fn)
	})

	// Serve the setup UI
	r.HandleFunc("/setup*", func(w http.ResponseWriter, r *http.Request) {
		ext := path.Ext(r.URL.Path)
		if ext == ".html" || ext == "" {
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
		generated := filepath.Clean(r.Form.Get("generated"))
		metadata := filepath.Clean(r.Form.Get("metadata"))
		cache := filepath.Clean(r.Form.Get("cache"))
		//downloads := filepath.Clean(r.Form.Get("downloads")) // TODO
		downloads := filepath.Join(metadata, "downloads")

		exists, _ := utils.DirExists(stash)
		if !exists || stash == "." {
			http.Error(w, fmt.Sprintf("the stash path either doesn't exist, or is not a directory <%s>.  Go back and try again.", stash), 500)
			return
		}

		exists, _ = utils.DirExists(generated)
		if !exists || generated == "." {
			http.Error(w, fmt.Sprintf("the generated path either doesn't exist, or is not a directory <%s>.  Go back and try again.", generated), 500)
			return
		}

		exists, _ = utils.DirExists(metadata)
		if !exists || metadata == "." {
			http.Error(w, fmt.Sprintf("the metadata path either doesn't exist, or is not a directory <%s>  Go back and try again.", metadata), 500)
			return
		}

		exists, _ = utils.DirExists(cache)
		if !exists || cache == "." {
			http.Error(w, fmt.Sprintf("the cache path either doesn't exist, or is not a directory <%s>  Go back and try again.", cache), 500)
			return
		}

		_ = os.Mkdir(downloads, 0755)

		config.Set(config.Stash, stash)
		config.Set(config.Generated, generated)
		config.Set(config.Metadata, metadata)
		config.Set(config.Cache, cache)
		config.Set(config.Downloads, downloads)
		if err := config.Write(); err != nil {
			http.Error(w, fmt.Sprintf("there was an error saving the config file: %s", err), 500)
			return
		}

		http.Redirect(w, r, "/", 301)
	})

	// Serve the web app
	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		ext := path.Ext(r.URL.Path)
		if ext == ".html" || ext == "" {
			data := uiBox.Bytes("index.html")
			_, _ = w.Write(data)
		} else {
			http.FileServer(uiBox).ServeHTTP(w, r)
		}
	})

	address := config.GetHost() + ":" + strconv.Itoa(config.GetPort())
	if tlsConfig := makeTLSConfig(); tlsConfig != nil {
		httpsServer := &http.Server{
			Addr:      address,
			Handler:   r,
			TLSConfig: tlsConfig,
		}

		go func() {
			logger.Infof("stash is running on HTTPS at https://" + address + "/")
			logger.Fatal(httpsServer.ListenAndServeTLS("", ""))
		}()
	} else {
		server := &http.Server{
			Addr:    address,
			Handler: r,
		}

		go func() {
			logger.Infof("stash is running on HTTP at http://" + address + "/")
			logger.Fatal(server.ListenAndServe())
		}()
	}
}

func makeTLSConfig() *tls.Config {
	cert, err := ioutil.ReadFile(paths.GetSSLCert())
	if err != nil {
		return nil
	}

	key, err := ioutil.ReadFile(paths.GetSSLKey())
	if err != nil {
		return nil
	}

	certs := make([]tls.Certificate, 1)
	certs[0], err = tls.X509KeyPair(cert, key)
	if err != nil {
		return nil
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
		if strings.Compare("https", r.URL.Scheme) == 0 || r.Proto == "HTTP/2.0" || r.Header.Get("X-Forwarded-Proto") == "https" {
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
		if !config.IsValid() && shouldRedirect {
			if !strings.HasPrefix(r.URL.Path, "/setup") {
				http.Redirect(w, r, "/setup", 301)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
