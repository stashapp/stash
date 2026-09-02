package scraper

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/enetx/g"
	"github.com/enetx/surf"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/net/html/charset"

	"github.com/stashapp/stash/pkg/logger"
)

const scrapeDefaultSleep = time.Second * 2

var proxyAuthRE = regexp.MustCompile(`^(https?:\/\/)(([\P{Cc}]+):([\P{Cc}]+)@)?(([a-zA-Z0-9][a-zA-Z0-9.-]*)(:[0-9]{1,5})?)`)

func loadURL(ctx context.Context, loadURL string, client *http.Client, def Definition, globalConfig GlobalConfig) (io.Reader, error) {
	driverOptions := def.DriverOptions
	if driverOptions != nil {
		switch {
		case driverOptions.UseCDP:
			return urlFromCDP(ctx, loadURL, *driverOptions, globalConfig)
		case driverOptions.UseSurf:
			return urlFromSurf(ctx, loadURL, *driverOptions, def, globalConfig)
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, loadURL, nil)
	if err != nil {
		return nil, err
	}

	jar, err := def.jar()
	if err != nil {
		return nil, fmt.Errorf("error creating cookie jar: %w", err)
	}

	u, err := url.Parse(loadURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing url %s: %w", loadURL, err)
	}

	// Fetch relevant cookies from the jar for url u and add them to the request
	cookies := jar.Cookies(u)
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	userAgent := globalConfig.GetScraperUserAgent()
	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}

	if driverOptions != nil { // setting the Headers after the UA allows us to override it from inside the scraper
		for _, h := range driverOptions.Headers {
			if h.Key != "" {
				req.Header.Set(h.Key, h.Value)
				logger.Debugf("[scraper] adding header <%s:%s>", h.Key, h.Value)
			}
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("http error %d:%s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewReader(body)
	printCookies(jar, def, "Jar cookies found for scraper urls")
	return charset.NewReader(bodyReader, resp.Header.Get("Content-Type"))
}

// func urlFromSurf uses enetx/surf with TLS browser emulation to bypass fingerprint-based blocking.
// this is a step down from CDP but faster and more lightweight and can succeed where CDP might fail
func urlFromSurf(ctx context.Context, loadURL string, driverOptions scraperDriverOptions, def Definition, globalConfig GlobalConfig) (io.Reader, error) {
	// get cookies
	jar, err := def.jar()
	if err != nil {
		return nil, fmt.Errorf("error creating cookie jar: %w", err)
	}

	// create client
	builder := surf.NewClient().
		Builder().
		Impersonate().Chrome()

	// get global proxy
	proxyURL := globalConfig.GetProxy()
	if proxyURL != "" {
		builder = builder.Proxy(g.String(proxyURL))
	}

	client := builder.
		Build().
		Unwrap().
		Std() // use stdlib adapter

	// pass in jar directly
	client.Jar = jar

	// try to mimic normal method above
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, loadURL, nil)
	if err != nil {
		return nil, err
	}
	// remove User-Agent header. This undermines TLS fingerprinting
	// because of GREASE (RFC 8701)
	// older fingerprints + UAs are still sustainable.
	for _, h := range driverOptions.Headers {
		if h.Key != "" {
			if strings.ToLower(h.Key) == "user-agent" {
				continue
			}
			req.Header.Set(h.Key, h.Value)
			logger.Debugf("[scraper] adding header <%s:%s>", h.Key, h.Value)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("http error %d:%s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewReader(body)
	printCookies(jar, def, "Jar cookies found for scraper urls")
	return charset.NewReader(bodyReader, resp.Header.Get("Content-Type"))
}

// func urlFromCDP uses chrome cdp and DOM to load and process the url
// if remote is set as true in the scraperConfig  it will try to use localhost:9222
// else it will look for google-chrome in path
func urlFromCDP(ctx context.Context, urlCDP string, driverOptions scraperDriverOptions, globalConfig GlobalConfig) (io.Reader, error) {
	if !driverOptions.UseCDP {
		return nil, fmt.Errorf("url shouldn't be fetched through CDP")
	}

	sleepDuration := scrapeDefaultSleep

	if driverOptions.Sleep > 0 {
		sleepDuration = time.Duration(driverOptions.Sleep) * time.Second
	}

	// if scraperCDPPath is a remote address, then allocate accordingly
	cdpPath := globalConfig.GetScraperCDPPath()
	if cdpPath != "" {
		var cancelAct context.CancelFunc

		if isCDPPathHTTP(globalConfig) || isCDPPathWS(globalConfig) {
			remote := cdpPath

			// -------------------------------------------------------------------
			// #1023
			// when chromium is listening over RDP it only accepts requests
			// with host headers that are either IPs or `localhost`
			cdpURL, err := url.Parse(remote)
			if err != nil {
				return nil, fmt.Errorf("failed to parse CDP Path: %v", err)
			}
			hostname := cdpURL.Hostname()
			if hostname != "localhost" {
				if net.ParseIP(hostname) == nil { // not an IP
					addr, err := net.LookupIP(hostname)
					if err != nil || len(addr) == 0 { // can not resolve to IP
						return nil, fmt.Errorf("CDP: hostname <%s> can not be resolved", hostname)
					}
					if len(addr[0]) == 0 { // nil IP
						return nil, fmt.Errorf("CDP: hostname <%s> resolved to nil", hostname)
					}
					// addr is a valid IP
					// replace the host part of the cdpURL with the IP
					cdpURL.Host = strings.Replace(cdpURL.Host, hostname, addr[0].String(), 1)
					// use that for remote
					remote = cdpURL.String()
				}
			}
			// --------------------------------------------------------------------

			// if CDPPath is http(s) then we need to get the websocket URL
			if isCDPPathHTTP(globalConfig) {
				var err error
				remote, err = getRemoteCDPWSAddress(ctx, remote)
				if err != nil {
					return nil, err
				}
			}

			ctx, cancelAct = chromedp.NewRemoteAllocator(ctx, remote)
		} else {
			// use a temporary user directory for chrome
			dir, err := os.MkdirTemp("", "stash-chromedp")
			if err != nil {
				return nil, err
			}
			defer os.RemoveAll(dir)

			opts := append(chromedp.DefaultExecAllocatorOptions[:],
				chromedp.UserDataDir(dir),
				chromedp.ExecPath(cdpPath),
			)
			if globalConfig.GetProxy() != "" {
				url, _, _ := splitProxyAuth(globalConfig.GetProxy())
				opts = append(opts, chromedp.ProxyServer(url))
			}

			ctx, cancelAct = chromedp.NewExecAllocator(ctx, opts...)
		}

		defer cancelAct()
	}

	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	// add a fixed timeout for the http request
	ctx, cancel = context.WithTimeout(ctx, scrapeGetTimeout)
	defer cancel()

	var res string
	headers := cdpHeaders(driverOptions)

	// Track the main document response so that, if it turns out to be JSON,
	// we can pull the raw response body via CDP's Network domain instead of
	// reading the rendered DOM. Chrome's built-in JSON viewer wraps raw JSON
	// responses in an HTML pretty-printer when navigated to directly, so
	// OuterHTML on such a page returns HTML containing the JSON rather than
	// the JSON itself.
	var jsonDoc jsonDocumentTracker
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		if ev, ok := ev.(*network.EventResponseReceived); ok {
			if ev.Type == network.ResourceTypeDocument {
				jsonDoc.markDocument(ev.RequestID, ev.Response.MimeType)
			}
		}
	})

	if proxyUsesAuth(globalConfig.GetProxy()) {
		_, user, pass := splitProxyAuth(globalConfig.GetProxy())

		// Based on https://github.com/chromedp/examples/blob/master/proxy/main.go
		lctx, lcancel := context.WithCancel(ctx)
		chromedp.ListenTarget(lctx, func(ev interface{}) {
			switch ev := ev.(type) {
			case *fetch.EventRequestPaused:
				go func() {
					_ = chromedp.Run(ctx, fetch.ContinueRequest(ev.RequestID))
				}()
			case *fetch.EventAuthRequired:
				if ev.AuthChallenge.Source == fetch.AuthChallengeSourceProxy {
					go func() {
						_ = chromedp.Run(ctx,
							fetch.ContinueWithAuth(ev.RequestID, &fetch.AuthChallengeResponse{
								Response: fetch.AuthChallengeResponseResponseProvideCredentials,
								Username: user,
								Password: pass,
							}),
							// Chrome will remember the credential for the current instance,
							// so we can disable the fetch domain once credential is provided.
							// Please file an issue if Chrome does not work in this way.
							fetch.Disable(),
						)
						// and cancel the event handler too.
						lcancel()
					}()
				}
			}
		})
	}

	err := chromedp.Run(ctx,
		network.Enable(),
		setCDPCookies(driverOptions),
		printCDPCookies(driverOptions, "Cookies found"),
		network.SetExtraHTTPHeaders(network.Headers(headers)),
		chromedp.Navigate(urlCDP),
		chromedp.Sleep(sleepDuration),
		setCDPClicks(driverOptions),
		chromedp.ActionFunc(func(ctx context.Context) error {
			requestID, isJSON := jsonDoc.mainDocument()

			if isJSON {
				body, err := network.GetResponseBody(requestID).Do(ctx)
				if err != nil {
					// Fall back to OuterHTML if the response body is no
					// longer available (e.g. evicted from Chrome's cache).
					// The scraper will most likely still fail downstream on
					// this fallback, since OuterHTML returns Chrome's
					// HTML-wrapped view of the JSON rather than the JSON
					// itself - but it's a better failure mode than losing
					// the response entirely.
					logger.Warnf("[scraper] could not get raw response body for JSON document, falling back to OuterHTML (the scrape will likely still fail as a result): %v", err)
					return chromedp.OuterHTML("html", &res, chromedp.ByQuery).Do(ctx)
				}
				res = string(body)
				return nil
			}
			return chromedp.OuterHTML("html", &res, chromedp.ByQuery).Do(ctx)
		}),
		printCDPCookies(driverOptions, "Cookies set"),
	)

	if err != nil {
		return nil, err
	}

	return strings.NewReader(res), nil
}

// click all xpaths listed in the scraper config
func setCDPClicks(driverOptions scraperDriverOptions) chromedp.Tasks {
	var tasks chromedp.Tasks
	for _, click := range driverOptions.Clicks { // for each click element find the node from the xpath and add a click action
		if click.XPath != "" {
			xpath := click.XPath
			waitDuration := scrapeDefaultSleep
			if click.Sleep > 0 {
				waitDuration = time.Duration(click.Sleep) * time.Second
			}

			action := chromedp.ActionFunc(func(ctx context.Context) error {
				var nodes []*cdp.Node
				if err := chromedp.Nodes(xpath, &nodes, chromedp.AtLeast(0)).Do(ctx); err != nil {
					logger.Debugf("Error %s looking for click xpath %s.\n", err, xpath)
					return err
				}
				if len(nodes) == 0 {
					logger.Debugf("Click xpath %s not found in page.\n", xpath)
					return nil
				}
				logger.Debugf("Clicking %s\n", xpath)
				return chromedp.MouseClickNode(nodes[0]).Do(ctx)
			})

			tasks = append(tasks, action)
			tasks = append(tasks, chromedp.Sleep(waitDuration))
		}

	}
	return tasks
}

// getRemoteCDPWSAddress returns the complete remote address that is required to access the cdp instance
func getRemoteCDPWSAddress(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	remote := result["webSocketDebuggerUrl"].(string)
	logger.Debugf("Remote cdp instance found %s", remote)
	return remote, err
}

// isJSONMimeType returns true if the given response MIME type indicates
// JSON content (e.g. "application/json", "application/ld+json",
// "text/json; charset=utf-8"), for deciding whether a CDP-fetched document
// should be read via its raw network response body rather than the
// rendered DOM.
func isJSONMimeType(mimeType string) bool {
	mimeType, _, _ = strings.Cut(mimeType, ";")
	mimeType = strings.TrimSpace(strings.ToLower(mimeType))

	return mimeType == "application/json" ||
		strings.HasSuffix(mimeType, "+json") ||
		mimeType == "text/json"
}

// jsonDocumentTracker records whether the CDP-loaded page's main document
// turned out to be JSON, and if so, which network request to fetch its raw
// body from.
//
// It only ever considers the *first* Document-type response it sees, via
// markDocument. Network.responseReceived fires with resourceType Document
// for iframe navigations too, not just the top-level one, so a naive
// "first Document that happens to be JSON" rule could latch onto an
// iframe's JSON response instead of the main page's HTML - misidentifying
// an HTML scrape as a JSON one. Redirects don't complicate this: a
// pre-redirect response never gets its own responseReceived event (that
// history arrives via requestWillBeSent's redirectResponse on the same
// request instead), so the first Document response reliably corresponds to
// the top-level navigation. A click-triggered navigation (via
// setCDPClicks, which runs after the initial Navigate) is not treated as a
// new "first" document either, so if a click leads to a JSON page, this
// still reports the original navigation's result.
//
// It's written from chromedp's event-processing goroutine (via
// markDocument, called from a ListenTarget callback) and read from the
// action sequence passed to chromedp.Run (via mainDocument), so access is
// synchronized with a mutex.
type jsonDocumentTracker struct {
	mutex     sync.Mutex
	requestID network.RequestID
	isJSON    bool
	recorded  bool
}

// markDocument records requestID as the main document if no document has
// been recorded yet. Subsequent calls after the first are no-ops.
func (t *jsonDocumentTracker) markDocument(requestID network.RequestID, mimeType string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.recorded {
		return
	}
	t.recorded = true
	t.requestID = requestID
	t.isJSON = isJSONMimeType(mimeType)
}

// mainDocument returns the main document's request ID and whether it was
// JSON.
func (t *jsonDocumentTracker) mainDocument() (network.RequestID, bool) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	return t.requestID, t.isJSON
}

func cdpHeaders(driverOptions scraperDriverOptions) map[string]interface{} {
	headers := map[string]interface{}{}
	if driverOptions.Headers != nil {
		for _, h := range driverOptions.Headers {
			if h.Key != "" {
				headers[h.Key] = h.Value
				logger.Debugf("[scraper] adding header <%s:%s>", h.Key, h.Value)
			}
		}
	}
	return headers
}

func proxyUsesAuth(proxyUrl string) bool {
	if proxyUrl == "" {
		return false
	}
	matches := proxyAuthRE.FindAllStringSubmatch(proxyUrl, -1)
	if matches != nil {
		split := matches[0]
		return len(split) == 0 || (len(split) > 5 && split[3] != "")
	}

	return false
}

func splitProxyAuth(proxyUrl string) (string, string, string) {
	if proxyUrl == "" {
		return "", "", ""
	}
	matches := proxyAuthRE.FindAllStringSubmatch(proxyUrl, -1)

	if matches != nil && len(matches[0]) > 5 {
		split := matches[0]
		return split[1] + split[5], split[3], split[4]
	}

	return proxyUrl, "", ""
}
