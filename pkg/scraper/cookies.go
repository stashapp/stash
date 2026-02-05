package scraper

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"golang.org/x/net/publicsuffix"

	"github.com/stashapp/stash/pkg/logger"
)

// jar constructs a cookie jar from a configuration
func (c Definition) jar() (*cookiejar.Jar, error) {
	opts := c.DriverOptions
	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		return nil, err
	}

	if opts == nil || opts.UseCDP {
		return jar, nil
	}

	for i, ckURL := range opts.Cookies {
		url, err := url.Parse(ckURL.CookieURL) // CookieURL must be valid, include schema
		if err != nil {
			logger.Warnf("skipping cookie [%d] for cookieURL %s: %v", i, ckURL.CookieURL, err)
			continue
		}

		var httpCookies []*http.Cookie
		for _, cookie := range ckURL.Cookies {
			c := &http.Cookie{
				Name:   cookie.Name,
				Value:  getCookieValue(cookie),
				Path:   cookie.Path,
				Domain: cookie.Domain,
			}
			httpCookies = append(httpCookies, c)
		}

		jar.SetCookies(url, httpCookies)
		if jar.Cookies(url) == nil {
			logger.Warnf("setting jar cookies for %s failed", url.String())
		}
	}

	return jar, nil
}

func getCookieValue(cookie *scraperCookies) string {
	if cookie.ValueRandom > 0 {
		return randomSequence(cookie.ValueRandom)
	}
	return cookie.Value
}

var characters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func randomSequence(n int) string {
	b := make([]rune, n)
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = characters[rand.Intn(len(characters))]
	}
	return string(b)
}

// printCookies prints all cookies from the given cookie jar
func printCookies(jar *cookiejar.Jar, scraperConfig Definition, msg string) {
	driverOptions := scraperConfig.DriverOptions
	if driverOptions != nil && !driverOptions.UseCDP {
		var foundURLs []*url.URL

		for _, ckURL := range driverOptions.Cookies { // go through all cookies
			url, err := url.Parse(ckURL.CookieURL) // CookieURL must be valid, include schema
			if err == nil {
				foundURLs = append(foundURLs, url)
			}
		}
		if len(foundURLs) > 0 {
			logger.Debugf("%s\n", msg)
			printJarCookies(jar, foundURLs)

		}
	}
}

// print all cookies from the jar of the native http client for given urls
func printJarCookies(jar *cookiejar.Jar, urls []*url.URL) {
	for _, url := range urls {
		logger.Debugf("Jar cookies for %s", url.String())
		for i, cookie := range jar.Cookies(url) {
			logger.Debugf("[%d]: Name: \"%s\" Value: \"%s\"", i, cookie.Name, cookie.Value)
		}
	}
}

// set all cookies listed in the scraper config
func setCDPCookies(driverOptions scraperDriverOptions) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			// create cookie expiration
			expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))

			for _, ckURL := range driverOptions.Cookies {
				for _, cookie := range ckURL.Cookies {
					err := network.SetCookie(cookie.Name, getCookieValue(cookie)).
						WithExpires(&expr).
						WithDomain(cookie.Domain).
						WithPath(cookie.Path).
						WithHTTPOnly(false).
						WithSecure(false).
						Do(ctx)
					if err != nil {
						return fmt.Errorf("could not set chrome cookie %s: %s", cookie.Name, err)
					}
				}
			}
			return nil
		}),
	}
}

// print cookies whose domain is included in the scraper  config
func printCDPCookies(driverOptions scraperDriverOptions, msg string) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		chromeCookies, err := network.GetCookies().Do(ctx)
		if err != nil {
			return err
		}

		scraperDomains := make(map[string]struct{})
		for _, ckURL := range driverOptions.Cookies {
			for _, cookie := range ckURL.Cookies {
				scraperDomains[cookie.Domain] = struct{}{}
			}
		}

		if len(scraperDomains) > 0 { // only print the cookies if they are listed in the scraper
			logger.Debugf("%s\n", msg)
			for i, cookie := range chromeCookies {
				_, ok := scraperDomains[cookie.Domain]
				if ok {
					logger.Debugf("[%d]: Name: \"%s\" Value: \"%s\"  Domain: \"%s\"", i, cookie.Name, cookie.Value, cookie.Domain)
				}
			}
		}
		return nil
	})
}
