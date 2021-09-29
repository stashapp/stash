package api

import (
	"context"
	"crypto/tls"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"net/url"
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
	"github.com/rs/cors"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/session"
	"github.com/stashapp/stash/pkg/utils"
)

var version string
var buildstamp string
var githash string

const loginEndPoint = "/login"

func allowUnauthenticated(r *http.Request) bool {
	return strings.HasPrefix(r.URL.Path, loginEndPoint) || r.URL.Path == "/css"
}

func authenticateHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := config.GetInstance()

			if c.GetSecurityTripwireAccessedFromPublicInternet() {
				w.WriteHeader(http.StatusForbidden)
				_, err := w.Write([]byte("Stash is exposed to the public internet without authentication, and is not serving any more content to protect your privacy. " +
					"More information and fixes are available at https://github.com/stashapp/stash/wiki/Authentication-Required-When-Accessing-Stash-From-the-Internet"))
				if err != nil {
					logger.Error(err)
				}
				return
			}

			userID, err := manager.GetInstance().SessionStore.Authenticate(w, r)
			if err != nil {
				if err != session.ErrUnauthorized {
					w.WriteHeader(http.StatusInternalServerError)
					_, err = w.Write([]byte(err.Error()))
					if err != nil {
						logger.Error(err)
					}
					return
				}

				// unauthorized error
				w.Header().Add("WWW-Authenticate", `FormBased`)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := r.Context()

			if c.HasCredentials() {
				// authentication is required
				if userID == "" && !allowUnauthenticated(r) {
					//authentication was not received, redirect
					// if graphql was requested, we just return a forbidden error
					if r.URL.Path == "/graphql" {
						w.Header().Add("WWW-Authenticate", `FormBased`)
						w.WriteHeader(http.StatusUnauthorized)
						return
					}

					// otherwise redirect to the login page
					u := url.URL{
						Path: "/login",
					}
					q := u.Query()
					q.Set(returnURLParam, r.URL.Path)
					u.RawQuery = q.Encode()
					http.Redirect(w, r, u.String(), http.StatusFound)
					return
				}
			} else {
				//authentication is not required
				//security fix: traffic from the public internet with no auth is disallowed
				if !c.GetDangerousAllowPublicWithoutAuth() && !c.IsNewSystem() {
					requestIPString := r.RemoteAddr[0:strings.LastIndex(r.RemoteAddr, ":")]
					requestIP := net.ParseIP(requestIPString)
					_, cgNatAddrSpace, _ := net.ParseCIDR("100.64.0.0/10")

					if r.Header.Get("X-FORWARDED-FOR") != "" {
						// Requst was proxied
						trustedProxies := c.GetTrustedProxies()
						if trustedProxies == "" {
							//validate proxies against local network only
							if !(requestIP.IsPrivate() || requestIP.IsLoopback() || cgNatAddrSpace.Contains(requestIP)) {
								securityActivateTripwireAccessedFromInternetWithoutAuth(c, w)
								return
							} else {
								// Safe to validate X-Forwarded-For
								proxyChain := strings.Split(r.Header.Get("X-FORWARDED-FOR"), ", ")
								for i := range proxyChain {
									ip := net.ParseIP(proxyChain[i])
									if !(ip.IsPrivate() || ip.IsLoopback() || cgNatAddrSpace.Contains(ip)) {
										securityActivateTripwireAccessedFromInternetWithoutAuth(c, w)
										return
									}
								}
							}
						} else {
							//validate proxies against trusted proxies list
							trustedProxies := strings.Split(trustedProxies, ", ")
							if isIPTrustedProxy(requestIP, trustedProxies) {
								// Safe to validate X-Forwarded-For
								proxyChain := strings.Split(r.Header.Get("X-FORWARDED-FOR"), ", ")
								// validate backwards, as only the last one is not attacker-controlled
								for i := len(proxyChain) - 1; i >= 0; i-- {
									ip := net.ParseIP(proxyChain[i])
									if i == 0 {
										//last entry is originating device, check if from the public internet
										if !(ip.IsPrivate() || ip.IsLoopback() || cgNatAddrSpace.Contains(ip)) {
											securityActivateTripwireAccessedFromInternetWithoutAuth(c, w)
											return
										}
									} else if !isIPTrustedProxy(ip, trustedProxies) {
										logger.Warn([]byte("Rejected request from untrusted proxy in chain:" + ip.String()))
										w.WriteHeader(http.StatusForbidden)
										return
									}
								}
							} else {
								// Proxy not on safe proxy list
								logger.Warn([]byte("Rejected request from untrusted proxy:" + requestIP.String()))
								w.WriteHeader(http.StatusForbidden)
								return
							}
						}
					} else {
						// request was not proxied
						if !(requestIP.IsPrivate() || requestIP.IsLoopback() || cgNatAddrSpace.Contains(requestIP)) {
							securityActivateTripwireAccessedFromInternetWithoutAuth(c, w)
							return
						}
					}
				}
			}

			ctx = session.SetCurrentUserID(ctx, userID)

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

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
			_, _ = w.Write(getLoginPage(loginUIBox))
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

			prefix := ""
			if r.Header.Get("X-Forwarded-Prefix") != "" {
				prefix = strings.TrimRight(r.Header.Get("X-Forwarded-Prefix"), "/")
			}

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
				panic(error.Error(err))
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
		panic(fmt.Errorf("error loading TLS config: %s", err.Error()))
	}

	server := &http.Server{
		Addr:      address,
		Handler:   r,
		TLSConfig: tlsConfig,
	}

	go func() {
		printVersion()
		printLatestVersion()
		logger.Infof("stash is listening on " + address)

		if tlsConfig != nil {
			logger.Infof("stash is running at https://" + displayAddress + "/")
			logger.Error(server.ListenAndServeTLS("", ""))
		} else {
			logger.Infof("stash is running at http://" + displayAddress + "/")
			logger.Error(server.ListenAndServe())
		}
	}()
}

func printVersion() {
	versionString := githash
	if version != "" {
		versionString = version + " (" + versionString + ")"
	}
	fmt.Printf("stash version: %s - %s\n", versionString, buildstamp)
}
func isIPTrustedProxy(ip net.IP, trustedProxies []string) bool {
	for _, v := range trustedProxies {
		if ip.Equal(net.ParseIP(v)) {
			return true
		}
	}
	return false
}

func securityActivateTripwireAccessedFromInternetWithoutAuth(c *config.Instance, w http.ResponseWriter) {
	logger.Error("Stash has been accessed from the internet, without authentication. \n" +
		"This is extremely dangerous! The whole world can see your stash page and browse your files! \n" +
		"You probably forwarded a port from your router. At the very least, add a password to stash in the settings. \n" +
		"Stash will not start again until you edit config.yml and change security_tripwire_accessed_from_public_internet to false. \n" +
		"More information is available at https://github.com/stashapp/stash/wiki/Authentication-Required-When-Accessing-Stash-From-the-Internet \n" +
		"Stash is not answering any other requests to protect your privacy.")
	c.Set(config.SecurityTripwireAccessedFromPublicInternet, true)
	err := c.Write()
	if err != nil {
		logger.Error(err)
	}
	w.WriteHeader(http.StatusForbidden)
	_, err = w.Write([]byte("You have attempted to access Stash over the internet, and authentication is not enabled. " +
		"This is extremely dangerous! The whole world can see your your stash page and browse your files! " +
		"Stash is not answering any other requests to protect your privacy. " +
		"Please read the log entry or visit https://github.com/stashapp/stash/wiki/Authentication-Required-When-Accessing-Stash-From-the-Internet "))
	if err != nil {
		logger.Error(err)
	}
	err = manager.GetInstance().Shutdown()
	if err != nil {
		logger.Error(err)
	}
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
		return nil, fmt.Errorf("error parsing key pair: %s", err.Error())
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
		prefix := ""
		if r.Header.Get("X-Forwarded-Prefix") != "" {
			prefix = strings.TrimRight(r.Header.Get("X-Forwarded-Prefix"), "/")
		}

		port := ""
		forwardedPort := r.Header.Get("X-Forwarded-Port")
		if forwardedPort != "" && forwardedPort != "80" && forwardedPort != "8080" {
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
