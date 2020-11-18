package scraper

import (
	"context"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"

	"github.com/stashapp/stash/pkg/logger"
)

// set cookies for the native http client
func setCookies(jar *cookiejar.Jar, scraperConfig config) {
	driverOptions := scraperConfig.DriverOptions
	if driverOptions != nil && !driverOptions.UseCDP {
		var foundURLs []*url.URL

		for _, ckURL := range driverOptions.Cookies { // go through all cookies
			url, err := url.Parse(ckURL.CookieURL) // CookieURL must be valid, include schema
			if err != nil {
				logger.Warnf("Skipping jar cookies for cookieURL %s. Error %s", ckURL.CookieURL, err)
			} else {
				var httpCookies []*http.Cookie
				var httpCookie *http.Cookie

				for _, cookie := range ckURL.Cookies {
					httpCookie = &http.Cookie{
						Name:   cookie.Name,
						Value:  cookie.Value,
						Path:   cookie.Path,
						Domain: cookie.Domain,
					}

					httpCookies = append(httpCookies, httpCookie)
				}
				jar.SetCookies(url, httpCookies) // jar.SetCookies only sets cookies with the domain matching the URL

				if jar.Cookies(url) == nil {
					logger.Warnf("Setting jar cookies for %s failed", url.String())
				} else {

					foundURLs = append(foundURLs, url)
				}
			}

		}
	}
}

// print all cookies from the jar of the native http client
func printCookies(jar *cookiejar.Jar, scraperConfig config) {
	driverOptions := scraperConfig.DriverOptions
	if driverOptions != nil && !driverOptions.UseCDP {
		var foundURLs []*url.URL

		for _, ckURL := range driverOptions.Cookies { // go through all cookies
			url, err := url.Parse(ckURL.CookieURL) // CookieURL must be valid, include schema
			if err == nil {
				foundURLs = append(foundURLs, url)
			}
		}
		printJarCookies(jar, foundURLs)
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
					success, err := network.SetCookie(cookie.Name, cookie.Value).
						WithExpires(&expr).
						WithDomain(cookie.Domain).
						WithPath(cookie.Path).
						WithHTTPOnly(false).
						WithSecure(false).
						Do(ctx)
					if err != nil {
						return err
					}
					if !success {
						return fmt.Errorf("could not set chrome cookie %s", cookie.Name)
					}

				}
			}
			return nil
		}),
	}
}

// print cookies whose domain is included in the scraper  config
func printCDPCookies(driverOptions scraperDriverOptions) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		chromeCookies, err := network.GetAllCookies().Do(ctx)
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
			logger.Debugf("Chrome cookies found for scraper")
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
