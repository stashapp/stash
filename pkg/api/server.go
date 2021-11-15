package api

import (
	"context"
	"crypto/tls"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	gqlHandler "github.com/99designs/gqlgen/graphql/handler"
	gqlExtension "github.com/99designs/gqlgen/graphql/handler/extension"
	gqlLru "github.com/99designs/gqlgen/graphql/handler/lru"
	gqlTransport "github.com/99designs/gqlgen/graphql/handler/transport"
	gqlPlayground "github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/websocket"
	"github.com/pkg/browser"
	"github.com/rs/cors"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

var version string
var buildstamp string
var githash string
var officialBuild string

func Start(uiBox embed.FS, loginUIBox embed.FS) {
	initialiseImages()

	r := chi.NewRouter()

	r.Use(middleware.Heartbeat("/healthz"))
	r.Use(authenticateHandler())
	visitedPluginHandler := manager.GetInstance().SessionStore.VisitedPluginHandler()
	r.Use(visitedPluginHandler)

	r.Use(middleware.Recoverer)

	c := config.GetInstance()
	if c.GetLogAccess() {
		r.Use(middleware.Logger)
	}
	r.Use(middleware.DefaultCompress)
	r.Use(middleware.StripSlashes)
	r.Use(cors.AllowAll().Handler)
	r.Use(BaseURLMiddleware)

	recoverFunc := func(ctx context.Context, err interface{}) error {
		logger.Error(err)
		debug.PrintStack()

		message := fmt.Sprintf("Internal system error. Error <%v>", err)
		return errors.New(message)
	}

	txnManager := manager.GetInstance().TxnManager
	pluginCache := manager.GetInstance().PluginCache
	resolver := &Resolver{
		txnManager:   txnManager,
		hookExecutor: pluginCache,
	}

	gqlSrv := gqlHandler.New(models.NewExecutableSchema(models.Config{Resolvers: resolver}))
	gqlSrv.SetRecoverFunc(recoverFunc)
	gqlSrv.AddTransport(gqlTransport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		KeepAlivePingInterval: 10 * time.Second,
	})
	gqlSrv.AddTransport(gqlTransport.Options{})
	gqlSrv.AddTransport(gqlTransport.GET{})
	gqlSrv.AddTransport(gqlTransport.POST{})
	gqlSrv.AddTransport(gqlTransport.MultipartForm{
		MaxUploadSize: c.GetMaxUploadSize(),
	})

	gqlSrv.SetQueryCache(gqlLru.New(1000))
	gqlSrv.Use(gqlExtension.Introspection{})

	gqlHandlerFunc := func(w http.ResponseWriter, r *http.Request) {
		gqlSrv.ServeHTTP(w, r)
	}

	// register GQL handler with plugin cache
	// chain the visited plugin handler
	manager.GetInstance().PluginCache.RegisterGQLHandler(visitedPluginHandler(http.HandlerFunc(gqlHandlerFunc)))

	r.HandleFunc("/graphql", gqlHandlerFunc)
	r.HandleFunc("/playground", gqlPlayground.Handler("GraphQL playground", "/graphql"))

	// session handlers
	r.Post(loginEndPoint, handleLogin(loginUIBox))
	r.Get("/logout", handleLogout(loginUIBox))

	r.Get(loginEndPoint, getLoginHandler(loginUIBox))

	r.Mount("/performer", performerRoutes{
		txnManager: txnManager,
	}.Routes())
	r.Mount("/scene", sceneRoutes{
		txnManager: txnManager,
	}.Routes())
	r.Mount("/image", imageRoutes{
		txnManager: txnManager,
	}.Routes())
	r.Mount("/studio", studioRoutes{
		txnManager: txnManager,
	}.Routes())
	r.Mount("/movie", movieRoutes{
		txnManager: txnManager,
	}.Routes())
	r.Mount("/tag", tagRoutes{
		txnManager: txnManager,
	}.Routes())
	r.Mount("/downloads", downloadsRoutes{}.Routes())

	r.HandleFunc("/css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		if !c.GetCSSEnabled() {
			return
		}

		// search for custom.css in current directory, then $HOME/.stash
		fn := c.GetCSSPath()
		exists, _ := utils.FileExists(fn)
		if !exists {
			return
		}

		http.ServeFile(w, r, fn)
	})

	r.HandleFunc("/login*", func(w http.ResponseWriter, r *http.Request) {
		ext := path.Ext(r.URL.Path)
		if ext == ".html" || ext == "" {
			prefix := getProxyPrefix(r.Header)

			data := getLoginPage(loginUIBox)
			baseURLIndex := strings.Replace(string(data), "%BASE_URL%", prefix+"/", 2)
			_, _ = w.Write([]byte(baseURLIndex))
		} else {
			r.URL.Path = strings.Replace(r.URL.Path, loginEndPoint, "", 1)
			loginRoot, err := fs.Sub(loginUIBox, loginRootDir)
			if err != nil {
				panic(err)
			}
			http.FileServer(http.FS(loginRoot)).ServeHTTP(w, r)
		}
	})

	// Serve static folders
	customServedFolders := c.GetCustomServedFolders()
	if customServedFolders != nil {
		r.HandleFunc("/custom/*", func(w http.ResponseWriter, r *http.Request) {
			r.URL.Path = strings.Replace(r.URL.Path, "/custom", "", 1)

			// map the path to the applicable filesystem location
			var dir string
			r.URL.Path, dir = customServedFolders.GetFilesystemLocation(r.URL.Path)
			if dir != "" {
				http.FileServer(http.Dir(dir)).ServeHTTP(w, r)
			} else {
				http.NotFound(w, r)
			}
		})
	}

	customUILocation := c.GetCustomUILocation()

	// Serve the web app
	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		const uiRootDir = "ui/v2.5/build"

		ext := path.Ext(r.URL.Path)

		if customUILocation != "" {
			if r.URL.Path == "index.html" || ext == "" {
				r.URL.Path = "/"
			}

			http.FileServer(http.Dir(customUILocation)).ServeHTTP(w, r)
			return
		}

		if ext == ".html" || ext == "" {
			data, err := uiBox.ReadFile(uiRootDir + "/index.html")
			if err != nil {
				panic(err)
			}

			prefix := getProxyPrefix(r.Header)
			baseURLIndex := strings.Replace(string(data), "%BASE_URL%", prefix+"/", 2)
			baseURLIndex = strings.Replace(baseURLIndex, "base href=\"/\"", fmt.Sprintf("base href=\"%s\"", prefix+"/"), 2)
			_, _ = w.Write([]byte(baseURLIndex))
		} else {
			isStatic, _ := path.Match("/static/*/*", r.URL.Path)
			if isStatic {
				w.Header().Add("Cache-Control", "max-age=604800000")
			}
			uiRoot, err := fs.Sub(uiBox, uiRootDir)
			if err != nil {
				panic(err)
			}
			http.FileServer(http.FS(uiRoot)).ServeHTTP(w, r)
		}
	})

	displayHost := c.GetHost()
	if displayHost == "0.0.0.0" {
		displayHost = "localhost"
	}
	displayAddress := displayHost + ":" + strconv.Itoa(c.GetPort())

	address := c.GetHost() + ":" + strconv.Itoa(c.GetPort())
	tlsConfig, err := makeTLSConfig(c)
	if err != nil {
		// assume we don't want to start with a broken TLS configuration
		panic(fmt.Errorf("error loading TLS config: %v", err))
	}

	server := &http.Server{
		Addr:      address,
		Handler:   r,
		TLSConfig: tlsConfig,
	}

	go func() {
		printVersion()
		printLatestVersion(context.TODO())
		logger.Infof("stash is listening on " + address)
		if tlsConfig != nil {
			displayAddress = "https://" + displayAddress + "/"
		} else {
			displayAddress = "http://" + displayAddress + "/"
		}

		// This can be done before actually starting the server, as modern browsers will
		// automatically reload the page if a local port is closed at page load and then opened.
		if !c.GetNoBrowser() && manager.GetInstance().IsDesktop() {
			err = browser.OpenURL(displayAddress)
			if err != nil {
				logger.Error("Could not open browser: " + err.Error())
			}
		}

		if tlsConfig != nil {
			logger.Infof("stash is running at " + displayAddress)
			logger.Error(server.ListenAndServeTLS("", ""))
		} else {
			logger.Infof("stash is running at " + displayAddress)
			logger.Error(server.ListenAndServe())
		}
	}()
}

func printVersion() {
	versionString := githash
	if IsOfficialBuild() {
		versionString += " - Official Build"
	} else {
		versionString += " - Unofficial Build"
	}
	if version != "" {
		versionString = version + " (" + versionString + ")"
	}
	fmt.Printf("stash version: %s - %s\n", versionString, buildstamp)
}

func IsOfficialBuild() bool {
	return officialBuild == "true"
}

func GetVersion() (string, string, string) {
	return version, githash, buildstamp
}

func makeTLSConfig(c *config.Instance) (*tls.Config, error) {
	c.InitTLS()
	certFile, keyFile := c.GetTLSFiles()

	if certFile == "" && keyFile == "" {
		// assume http configuration
		return nil, nil
	}

	// ensure both files are present
	if certFile == "" {
		return nil, errors.New("SSL certificate file must be present if key file is present")
	}

	if keyFile == "" {
		return nil, errors.New("SSL key file must be present if certificate file is present")
	}

	cert, err := os.ReadFile(certFile)
	if err != nil {
		return nil, fmt.Errorf("error reading SSL certificate file %s: %s", certFile, err.Error())
	}

	key, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("error reading SSL key file %s: %s", keyFile, err.Error())
	}

	certs := make([]tls.Certificate, 1)
	certs[0], err = tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, fmt.Errorf("error parsing key pair: %v", err)
	}
	tlsConfig := &tls.Config{
		Certificates: certs,
	}

	return tlsConfig, nil
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
		prefix := getProxyPrefix(r.Header)

		port := ""
		forwardedPort := r.Header.Get("X-Forwarded-Port")
		if forwardedPort != "" && forwardedPort != "80" && forwardedPort != "8080" && forwardedPort != "443" && !strings.Contains(r.Host, ":") {
			port = ":" + forwardedPort
		}

		baseURL := scheme + "://" + r.Host + port + prefix

		externalHost := config.GetInstance().GetExternalHost()
		if externalHost != "" {
			baseURL = externalHost + prefix
		}

		r = r.WithContext(context.WithValue(ctx, BaseURLCtxKey, baseURL))

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func getProxyPrefix(headers http.Header) string {
	prefix := ""
	if headers.Get("X-Forwarded-Prefix") != "" {
		prefix = strings.TrimRight(headers.Get("X-Forwarded-Prefix"), "/")
	}

	return prefix
}
