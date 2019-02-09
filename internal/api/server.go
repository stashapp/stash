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
	"github.com/stashapp/stash/internal/logger"
	"github.com/stashapp/stash/internal/models"
	"net/http"
	"path"
	"runtime/debug"
	"strings"
)

const httpPort = "9998"
const httpsPort = "9999"

var certsBox *packr.Box
var uiBox *packr.Box

func Start() {
	//port := os.Getenv("PORT")
	//if port == "" {
	//	port = defaultPort
	//}

	certsBox = packr.New("Cert Box", "../../certs")
	uiBox = packr.New("UI Box", "../../ui/v1/dist/stash-frontend")

	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.DefaultCompress)
	r.Use(middleware.StripSlashes)
	r.Use(cors.AllowAll().Handler)
	r.Use(BaseURLMiddleware)

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

	logger.Infof("stash is running on HTTPS at https://localhost:9999/")
	logger.Fatal(httpsServer.ListenAndServeTLS("", ""))
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