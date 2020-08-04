package scraper

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	jsoniter "github.com/json-iterator/go"
	"github.com/stashapp/stash/pkg/logger"
	"golang.org/x/net/html/charset"
	"golang.org/x/net/publicsuffix"
)

func loadURL(url string, scraperConfig config, globalConfig GlobalConfig) (io.Reader, error) {
	driverOptions := scraperConfig.DriverOptions
	if driverOptions != nil && driverOptions.UseCDP {
		// get the page using chrome dp
		return urlFromCDP(url, *driverOptions, globalConfig)
	}

	// get the page using http.Client
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, er := cookiejar.New(&options)
	if er != nil {
		return nil, er
	}

	client := &http.Client{
		Timeout: scrapeGetTimeout,
		// defaultCheckRedirect code with max changed from 10 to 20
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 20 {
				return errors.New("stopped after 20 redirects")
			}
			return nil
		},
		Jar: jar,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	userAgent := globalConfig.UserAgent
	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewReader(body)

	return charset.NewReader(bodyReader, resp.Header.Get("Content-Type"))
}

// func urlFromCDP uses chrome cdp and DOM to load and process the url
// if remote is set as true in the scraperConfig  it will try to use localhost:9222
// else it will look for google-chrome in path
func urlFromCDP(url string, driverOptions scraperDriverOptions, globalConfig GlobalConfig) (io.Reader, error) {
	const defaultSleep = 2

	if !driverOptions.UseCDP {
		return nil, fmt.Errorf("Url shouldn't be feetched through CDP")
	}

	sleep := defaultSleep

	if driverOptions.Sleep != 0 {
		sleep = driverOptions.Sleep
	}

	sleepDuration := time.Duration(sleep) * time.Second
	act := context.Background()

	// if scraperCDPPath is a remote address, then allocate accordingly
	if globalConfig.CDPPath != "" {
		var cancelAct context.CancelFunc

		if globalConfig.isCDPPathHTTP() || globalConfig.isCDPPathWS() {
			remote := globalConfig.CDPPath

			// if CDPPath is http(s) then we need to get the websocket URL
			if globalConfig.isCDPPathHTTP() {
				var err error
				remote, err = getRemoteCDPWSAddress(remote)
				if err != nil {
					return nil, err
				}
			}

			act, cancelAct = chromedp.NewRemoteAllocator(context.Background(), remote)
		} else {
			// user a temporary user directory for chrome
			dir, err := ioutil.TempDir("", "stash-chromedp")
			if err != nil {
				return nil, err
			}
			defer os.RemoveAll(dir)

			opts := append(chromedp.DefaultExecAllocatorOptions[:],
				chromedp.UserDataDir(dir),
				chromedp.ExecPath(globalConfig.CDPPath),
			)
			act, cancelAct = chromedp.NewExecAllocator(act, opts...)
		}

		defer cancelAct()
	}

	ctx, cancel := chromedp.NewContext(act)
	defer cancel()

	var res string
	err := chromedp.Run(ctx,
		network.Enable(),
		chromedp.Navigate(url),
		chromedp.Sleep(sleepDuration),
		chromedp.ActionFunc(func(ctx context.Context) error {
			node, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			res, err = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			return err
		}),
	)
	if err != nil {
		return nil, err
	}

	return strings.NewReader(res), nil
}

// getRemoteCDPWSAddress returns the complete remote address that is required to access the cdp instance
func getRemoteCDPWSAddress(address string) (string, error) {
	resp, err := http.Get(address)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	remote := result["webSocketDebuggerUrl"].(string)
	logger.Debugf("Remote cdp instance found %s", remote)
	return remote, err
}

func cdpNetwork(enable bool) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		if enable {
			network.Enable().Do(ctx)
		} else {
			network.Disable().Do(ctx)
		}
		return nil
	})
}
