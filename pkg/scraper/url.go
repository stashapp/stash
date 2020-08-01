package scraper

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
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
		return urlFromCDP(url, *driverOptions)
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

	return charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
}

// func urlFromCDP uses chrome cdp and DOM to load and process the url
// if remote is set as true in the scraperConfig  it will try to use localhost:9222
// else it will look for google-chrome in path
func urlFromCDP(url string, driverOptions scraperDriverOptions) (io.Reader, error) {
	if !driverOptions.UseCDP {
		return nil, fmt.Errorf("Url shouldn't be feetched through CDP")
	}

	remote := false
	sleep := 2

	if driverOptions.Remote {
		remote = true
	}

	if driverOptions.Sleep != 0 {
		sleep = driverOptions.Sleep
	}

	sleepDuration := time.Duration(sleep) * time.Second
	act := context.Background()

	if remote {
		var cancelAct context.CancelFunc
		remote, errCDP := getRemoteCDP()
		if errCDP != nil {
			return nil, errCDP
		}
		act, cancelAct = chromedp.NewRemoteAllocator(context.Background(), remote)
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

// func getRemoteCDP returns the complete remote address that is required to access the cdp instance
func getRemoteCDP() (string, error) {
	resp, err := http.Get("http://localhost:9222/json/version")
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
