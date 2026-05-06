package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"

	"golang.org/x/sys/cpu"

	"github.com/stashapp/stash/internal/build"
	"github.com/stashapp/stash/pkg/logger"
)

// we use the github REST V3 API as no login is required
const apiReleases string = "https://api.github.com/repos/stashapp/stash/releases"
const apiTags string = "https://api.github.com/repos/stashapp/stash/tags"
const apiAcceptHeader string = "application/vnd.github.v3+json"
const developmentTag string = "latest_develop"
const defaultSHLength int = 8 // default length of SHA short hash returned by <git rev-parse --short HEAD>

var stashReleases = func() map[string]string {
	return map[string]string{
		"darwin/amd64":  "stash-macos",
		"darwin/arm64":  "stash-macos",
		"linux/amd64":   "stash-linux",
		"windows/amd64": "stash-win.exe",
		"linux/arm":     "stash-linux-arm32v6",
		"linux/arm64":   "stash-linux-arm64v8",
		"linux/armv7":   "stash-linux-arm32v7",
	}
}

// isMacOSBundle checks if the application is running from within a macOS .app bundle
func isMacOSBundle() bool {
	exec, err := os.Executable()
	return err == nil && strings.Contains(exec, "Stash.app/")
}

// getWantedRelease determines which release variant to download based on platform and bundle type
func getWantedRelease(platform string) string {
	release := stashReleases()[platform]

	// On macOS, check if running from .app bundle
	if runtime.GOOS == "darwin" && isMacOSBundle() {
		return "Stash.app.zip"
	}

	return release
}

type githubReleasesResponse struct {
	Url              string
	Assets_url       string
	Upload_url       string
	Html_url         string
	Id               int64
	Node_id          string
	Tag_name         string
	Target_commitish string
	Name             string
	Draft            bool
	Author           githubAuthor
	Prerelease       bool
	Created_at       string
	Published_at     string
	Assets           []githubAsset
	Tarball_url      string
	Zipball_url      string
	Body             string
}

type githubAuthor struct {
	Login               string
	Id                  int64
	Node_id             string
	Avatar_url          string
	Gravatar_id         string
	Url                 string
	Html_url            string
	Followers_url       string
	Following_url       string
	Gists_url           string
	Starred_url         string
	Subscriptions_url   string
	Organizations_url   string
	Repos_url           string
	Events_url          string
	Received_events_url string
	Type                string
	Site_admin          bool
}

type githubAsset struct {
	Url                  string
	Id                   int64
	Node_id              string
	Name                 string
	Label                string
	Uploader             githubAuthor
	Content_type         string
	State                string
	Size                 int64
	Download_count       int64
	Created_at           string
	Updated_at           string
	Browser_download_url string
}

type githubTagResponse struct {
	Name        string
	Zipball_url string
	Tarball_url string
	Commit      struct {
		Sha string
		Url string
	}
	Node_id string
}

type LatestRelease struct {
	Version   string
	Hash      string
	ShortHash string
	Date      string
	Url       string
}

func makeGithubRequest(ctx context.Context, url string, output interface{}) error {
	transport := &http.Transport{Proxy: http.ProxyFromEnvironment}

	client := &http.Client{
		Timeout:   3 * time.Second,
		Transport: transport,
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	req.Header.Add("Accept", apiAcceptHeader) // gh api recommendation , send header with api version
	logger.Debugf("Github API request: %s", url)
	response, err := client.Do(req)

	if err != nil {
		//lint:ignore ST1005 Github is a proper capitalized noun
		return fmt.Errorf("Github API request failed: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		//lint:ignore ST1005 Github is a proper capitalized noun
		return fmt.Errorf("Github API request failed: %s", response.Status)
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		//lint:ignore ST1005 Github is a proper capitalized noun
		return fmt.Errorf("Github API read response failed: %w", err)
	}

	err = json.Unmarshal(data, output)
	if err != nil {
		return fmt.Errorf("unmarshalling Github API response failed: %w", err)
	}

	return nil
}

// GetLatestRelease gets latest release information from github API
// If running a build from the "master" branch, then the latest full release
// is used, otherwise it uses the release that is tagged with "latest_develop"
// which is the latest pre-release build.
func GetLatestRelease(ctx context.Context) (*LatestRelease, error) {
	arch := runtime.GOARCH

	// https://en.wikipedia.org/wiki/Comparison_of_ARM_cores
	// armv6 doesn't support any of these features
	isARMv7 := cpu.ARM.HasNEON || cpu.ARM.HasVFPv3 || cpu.ARM.HasVFPv3D16 || cpu.ARM.HasVFPv4
	if arch == "arm" && isARMv7 {
		arch = "armv7"
	}

	platform := fmt.Sprintf("%s/%s", runtime.GOOS, arch)
	wantedRelease := getWantedRelease(platform)

	url := apiReleases
	if build.IsDevelop() {
		// get the release tagged with the development tag
		url += "/tags/" + developmentTag
	} else {
		// just get the latest full release
		url += "/latest"
	}

	var release githubReleasesResponse
	err := makeGithubRequest(ctx, url, &release)
	if err != nil {
		return nil, err
	}

	version := release.Name
	if release.Prerelease {
		// find version in prerelease name
		re := regexp.MustCompile(`v[\w-\.]+-\d+-g[0-9a-f]+`)
		if match := re.FindString(version); match != "" {
			version = match
		}
	}

	latestHash, err := getReleaseHash(ctx, release.Tag_name)
	if err != nil {
		return nil, err
	}

	var releaseDate string
	if publishedAt, err := time.Parse(time.RFC3339, release.Published_at); err == nil {
		releaseDate = publishedAt.Format("2006-01-02")
	}

	var releaseUrl string
	if wantedRelease != "" {
		for _, asset := range release.Assets {
			if asset.Name == wantedRelease {
				releaseUrl = asset.Browser_download_url
				break
			}
		}
	}

	_, githash, _ := build.Version()
	shLength := len(githash)
	if shLength == 0 {
		shLength = defaultSHLength
	}

	return &LatestRelease{
		Version:   version,
		Hash:      latestHash,
		ShortHash: latestHash[:shLength],
		Date:      releaseDate,
		Url:       releaseUrl,
	}, nil
}

func getReleaseHash(ctx context.Context, tagName string) (string, error) {
	// Start with a small page size if not searching for latest_develop
	perPage := 10
	if tagName == developmentTag {
		perPage = 100
	}

	// Limit to 5 pages, ie 500 tags - should be plenty
	for page := 1; page <= 5; {
		url := fmt.Sprintf("%s?per_page=%d&page=%d", apiTags, perPage, page)
		tags := []githubTagResponse{}
		err := makeGithubRequest(ctx, url, &tags)
		if err != nil {
			return "", err
		}

		for _, tag := range tags {
			if tag.Name == tagName {
				if len(tag.Commit.Sha) != 40 {
					return "", errors.New("invalid Github API response")
				}
				return tag.Commit.Sha, nil
			}
		}

		if len(tags) == 0 {
			break
		}

		// if not found in the first 10, search again on page 1 with the first 100
		if perPage == 10 {
			perPage = 100
		} else {
			page++
		}
	}

	return "", errors.New("invalid Github API response")
}

func printLatestVersion(ctx context.Context) {
	latestRelease, err := GetLatestRelease(ctx)
	if err != nil {
		logger.Errorf("Couldn't retrieve latest version: %v", err)
	} else {
		_, githash, _ := build.Version()
		switch {
		case githash == "":
			logger.Infof("Latest version: %s (%s)", latestRelease.Version, latestRelease.ShortHash)
		case githash == latestRelease.ShortHash:
			logger.Infof("Version %s (%s) is already the latest released", latestRelease.Version, latestRelease.ShortHash)
		default:
			logger.Infof("New version available: %s (%s)", latestRelease.Version, latestRelease.ShortHash)
		}
	}
}
