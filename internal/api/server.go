package api

import (
	"context"
	"crypto/tls"
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
	"github.com/vearutop/statigz"

	"github.com/go-chi/httplog"
	"github.com/rs/cors"
	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/ui"
)

var version string
var buildstamp string
var githash string

var uiBox = ui.UIBox
var loginUIBox = ui.LoginUIBox

func Start() error {
	initialiseImages()

	r := chi.NewRouter()

	r.Use(middleware.Heartbeat("/healthz"))
	r.Use(authenticateHandler())
	visitedPluginHandler := manager.GetInstance().SessionStore.VisitedPluginHandler()
	r.Use(visitedPluginHandler)

	r.Use(middleware.Recoverer)

	c := config.GetInstance()
	if c.GetLogAccess() {
		httpLogger := httplog.NewLogger("Stash", httplog.Options{
			Concise: true,
		})
		r.Use(httplog.RequestLogger(httpLogger))
	}
	r.Use(SecurityHeadersMiddleware)
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

	txnManager := manager.GetInstance().Repository
	pluginCache := manager.GetInstance().PluginCache
	resolver := &Resolver{
		txnManager:   txnManager,
		repository:   txnManager,
		hookExecutor: pluginCache,
	}

	gqlSrv := gqlHandler.New(NewExecutableSchema(Config{Resolvers: resolver}))
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
		txnManager:      txnManager,
		performerFinder: txnManager.Performer,
	}.Routes())
	r.Mount("/scene", sceneRoutes{
		txnManager:        txnManager,
		sceneFinder:       txnManager.Scene,
		sceneMarkerFinder: txnManager.SceneMarker,
		tagFinder:         txnManager.Tag,
	}.Routes())
	r.Mount("/image", imageRoutes{
		txnManager:  txnManager,
		imageFinder: txnManager.Image,
	}.Routes())
	r.Mount("/studio", studioRoutes{
		txnManager:   txnManager,
		studioFinder: txnManager.Studio,
	}.Routes())
	r.Mount("/movie", movieRoutes{
		txnManager:  txnManager,
		movieFinder: txnManager.Movie,
	}.Routes())
	r.Mount("/tag", tagRoutes{
		txnManager: txnManager,
		tagFinder:  txnManager.Tag,
	}.Routes())
	r.Mount("/downloads", downloadsRoutes{}.Routes())

	r.HandleFunc("/css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		if !c.GetCSSEnabled() {
			return
		}

		// search for custom.css in current directory, then $HOME/.stash
		fn := c.GetCSSPath()
		exists, _ := fsutil.FileExists(fn)
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
	static := statigz.FileServer(uiBox)

	// Serve the web app
	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		const uiRootDir = "v2.5/build"

		ext := path.Ext(r.URL.Path)

		if customUILocation != "" {
			if r.URL.Path == "index.html" || ext == "" {
				r.URL.Path = "/"
			}

			http.FileServer(http.Dir(customUILocation)).ServeHTTP(w, r)
			return
		}

		if ext == ".html" || ext == "" {
			themeColor := c.GetThemeColor()
			data, err := uiBox.ReadFile(uiRootDir + "/index.html")
			if err != nil {
				panic(err)
			}

			prefix := getProxyPrefix(r.Header)
			baseURLIndex := strings.ReplaceAll(string(data), "%COLOR%", themeColor)
			baseURLIndex = strings.ReplaceAll(baseURLIndex, "/%BASE_URL%", prefix)
			baseURLIndex = strings.Replace(baseURLIndex, "base href=\"/\"", fmt.Sprintf("base href=\"%s\"", prefix+"/"), 1)
			_, _ = w.Write([]byte(baseURLIndex))
		} else {
			isStatic, _ := path.Match("/static/*/*", r.URL.Path)
			if isStatic {
				w.Header().Add("Cache-Control", "max-age=604800000")
			}

			prefix := getProxyPrefix(r.Header)
			if prefix != "" {
				r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)
			}
			r.URL.Path = uiRootDir + r.URL.Path

			static.ServeHTTP(w, r)
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

	printVersion()
	go printLatestVersion(context.TODO())
	logger.Infof("stash is listening on " + address)
	if tlsConfig != nil {
		displayAddress = "https://" + displayAddress + "/"
	} else {
		displayAddress = "http://" + displayAddress + "/"
	}

	logger.Infof("stash is running at " + displayAddress)
	if tlsConfig != nil {
		err = server.ListenAndServeTLS("", "")
	} else {
		err = server.ListenAndServe()
	}

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func printVersion() {
	versionString := githash
	if config.IsOfficialBuild() {
		versionString += " - Official Build"
	} else {
		versionString += " - Unofficial Build"
	}
	if version != "" {
		versionString = version + " (" + versionString + ")"
	}
	fmt.Printf("stash version: %s - %s\n", versionString, buildstamp)
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

func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		c := config.GetInstance()
		connectableOrigins := "connect-src data: 'self'"

		// Workaround Safari bug https://bugs.webkit.org/show_bug.cgi?id=201591
		// Allows websocket requests to any origin
		connectableOrigins += " ws: wss:"

		// The graphql playground pulls its frontend from a cdn
		connectableOrigins += " https://cdn.jsdelivr.net "

		if !c.IsNewSystem() && c.GetHandyKey() != "" {
			connectableOrigins += " https://www.handyfeeling.com"
		}
		connectableOrigins += "; "

		cspDirectives := "default-src data: 'self' 'unsafe-inline';" + connectableOrigins + "img-src data: *; script-src 'self' https://cdn.jsdelivr.net 'unsafe-inline' 'unsafe-eval'; style-src 'self' https://cdn.jsdelivr.net 'unsafe-inline'; style-src-elem 'self' https://cdn.jsdelivr.net 'unsafe-inline'; media-src 'self' blob:; child-src 'none'; object-src 'none'; form-action 'self'"

		w.Header().Set("Referrer-Policy", "same-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-XSS-Protection", "1")
		w.Header().Set("Content-Security-Policy", cspDirectives)

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

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

		baseURL := scheme + "://" + r.Host + prefix

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
