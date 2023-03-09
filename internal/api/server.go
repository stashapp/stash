package api

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"regexp"
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
	"github.com/stashapp/stash/internal/api/loaders"
	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/plugin"
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

	dataloaders := loaders.Middleware{
		DatabaseProvider: txnManager,
		Repository:       txnManager,
	}

	r.Use(dataloaders.Middleware)

	pluginCache := manager.GetInstance().PluginCache
	sceneService := manager.GetInstance().SceneService
	imageService := manager.GetInstance().ImageService
	galleryService := manager.GetInstance().GalleryService
	resolver := &Resolver{
		txnManager:     txnManager,
		repository:     txnManager,
		sceneService:   sceneService,
		imageService:   imageService,
		galleryService: galleryService,
		hookExecutor:   pluginCache,
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
	// also requires the dataloader middleware
	gqlHandler := visitedPluginHandler(dataloaders.Middleware(http.HandlerFunc(gqlHandlerFunc)))
	manager.GetInstance().PluginCache.RegisterGQLHandler(gqlHandler)

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
		fileFinder:        txnManager.File,
		captionFinder:     txnManager.File,
		sceneMarkerFinder: txnManager.SceneMarker,
		tagFinder:         txnManager.Tag,
	}.Routes())
	r.Mount("/image", imageRoutes{
		txnManager:  txnManager,
		imageFinder: txnManager.Image,
		fileFinder:  txnManager.File,
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

	r.HandleFunc("/css", cssHandler(c, pluginCache))
	r.HandleFunc("/javascript", javascriptHandler(c, pluginCache))
	r.HandleFunc("/customlocales", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if c.GetCustomLocalesEnabled() {
			// search for custom-locales.json in current directory, then $HOME/.stash
			fn := c.GetCustomLocalesPath()
			exists, _ := fsutil.FileExists(fn)
			if exists {
				http.ServeFile(w, r, fn)
				return
			}
		}
		_, _ = w.Write([]byte("{}"))
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
		r.Mount("/custom", customRoutes{
			servedFolders: customServedFolders,
		}.Routes())
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
		// disable http/2 support by default
		// when http/2 is enabled, we are unable to hijack and close
		// the connection/request. This is necessary to stop running
		// streams when deleting a scene file.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
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

func copyFile(w io.Writer, path string) (time.Time, error) {
	f, err := os.Open(path)
	if err != nil {
		return time.Time{}, err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return time.Time{}, err
	}

	_, err = io.Copy(w, f)

	return info.ModTime(), err
}

func serveFiles(w http.ResponseWriter, r *http.Request, name string, paths []string) {
	buffer := bytes.Buffer{}

	latestModTime := time.Time{}

	for _, path := range paths {
		modTime, err := copyFile(&buffer, path)
		if err != nil {
			logger.Errorf("error serving file %s: %v", path, err)
		} else {
			if modTime.After(latestModTime) {
				latestModTime = modTime
			}
			buffer.Write([]byte("\n"))
		}
	}

	// Always revalidate with server
	w.Header().Set("Cache-Control", "no-cache")

	bufferReader := bytes.NewReader(buffer.Bytes())
	http.ServeContent(w, r, name, latestModTime, bufferReader)
}

func cssHandler(c *config.Instance, pluginCache *plugin.Cache) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// concatenate with plugin css files
		w.Header().Set("Content-Type", "text/css")

		// add plugin css files first
		var paths []string

		for _, p := range pluginCache.ListPlugins() {
			paths = append(paths, p.UI.CSS...)
		}

		if c.GetCSSEnabled() {
			// search for custom.css in current directory, then $HOME/.stash
			fn := c.GetCSSPath()
			exists, _ := fsutil.FileExists(fn)
			if exists {
				paths = append(paths, fn)
			}
		}

		serveFiles(w, r, "custom.css", paths)
	}
}

func javascriptHandler(c *config.Instance, pluginCache *plugin.Cache) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/javascript")

		// add plugin javascript files first
		var paths []string

		for _, p := range pluginCache.ListPlugins() {
			paths = append(paths, p.UI.Javascript...)
		}

		if c.GetJavascriptEnabled() {
			// search for custom.js in current directory, then $HOME/.stash
			fn := c.GetJavascriptPath()
			exists, _ := fsutil.FileExists(fn)
			if exists {
				paths = append(paths, fn)
			}
		}

		serveFiles(w, r, "custom.js", paths)
	}
}

func printVersion() {
	var versionString string
	switch {
	case version != "":
		if githash != "" && !IsDevelop() {
			versionString = version + " (" + githash + ")"
		} else {
			versionString = version
		}
	case githash != "":
		versionString = githash
	default:
		versionString = "unknown"
	}
	if config.IsOfficialBuild() {
		versionString += " - Official Build"
	} else {
		versionString += " - Unofficial Build"
	}
	if buildstamp != "" {
		versionString += " - " + buildstamp
	}
	logger.Infof("stash version: %s\n", versionString)
}

func GetVersion() (string, string, string) {
	return version, githash, buildstamp
}

func IsDevelop() bool {
	if githash == "" {
		return false
	}

	// if the version is suffixed with -x-xxxx, then we are running a development build
	develop := false
	re := regexp.MustCompile(`-\d+-g\w+$`)
	if re.MatchString(version) {
		develop = true
	}
	return develop
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
		return nil, fmt.Errorf("error reading SSL certificate file %s: %v", certFile, err)
	}

	key, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("error reading SSL key file %s: %v", keyFile, err)
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

		cspDirectives := "default-src data: 'self' 'unsafe-inline';" + connectableOrigins + "img-src data: *; script-src 'self' https://cdn.jsdelivr.net 'unsafe-inline' 'unsafe-eval'; style-src 'self' https://cdn.jsdelivr.net 'unsafe-inline'; style-src-elem 'self' https://cdn.jsdelivr.net 'unsafe-inline'; media-src 'self' blob:; child-src 'none'; worker-src blob:; object-src 'none'; form-action 'self'"

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

		scheme := "http"
		if strings.Compare("https", r.URL.Scheme) == 0 || r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
			scheme = "https"
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
