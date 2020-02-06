package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"runtime"
	"time"

	"github.com/stashapp/stash/pkg/logger"
)

//we use the github REST V3 API as no login is required
const apiReleases string = "https://api.github.com/repos/stashapp/stash/releases"
const apiAcceptHeader string = "application/vnd.github.v3+json"
const developmentTag string = "latest_develop"

var stashReleases = func() map[string]string {
	return map[string]string{
		"windows/amd64": "stash-win.exe",
		"linux/amd64":   "stash-linux",
		"darwin/amd64":  "stash-osx",
		"linux/arm":     "stash-pi",
	}
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

// GetLatestVersion gets latest version (git commit hash) from github API
// If running a build from the "master" branch, then the latest full release
// is used, otherwise it uses the release that is tagged with "latest_develop"
// which is the latest pre-release build.
func GetLatestVersion(shortHash bool) (latestVersion string, latestRelease string, err error) {

	platform := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
	wantedRelease := stashReleases()[platform]

	version, _, _ := GetVersion()
	if version == "" {
		return "", "", fmt.Errorf("Stash doesn't have a version. Version check not supported.")
	}

	// if the version is suffixed with -x-xxxx, then we are running a development build
	usePreRelease := false
	re := regexp.MustCompile(`-\d+-g\w+$`)
	if re.MatchString(version) {
		usePreRelease = true
	}

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	url := apiReleases
	if !usePreRelease {
		// just get the latest full release
		url += "/latest"
	} else {
		// get the release tagged with the development tag
		url += "/tags/" + developmentTag
	}
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Accept", apiAcceptHeader) // gh api recommendation , send header with api version
	response, err := client.Do(req)

	release := githubReleasesResponse{}

	if response.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("Github API request failed: %s", response.Status)
	}

	if err != nil {
		return "", "", fmt.Errorf("Github API request failed: %s", err)
	}

	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", "", fmt.Errorf("Github API read response failed: %s", err)
	}
	err = json.Unmarshal(data, &release)
	if err != nil {
		return "", "", fmt.Errorf("Unmarshalling Github API response failed: %s", err)
	}

	if release.Prerelease == usePreRelease {
		if shortHash {
			latestVersion = release.Target_commitish[0:7] //shorthash is first 7 digits of git commit hash
		} else {
			latestVersion = release.Target_commitish
		}
		if wantedRelease != "" {
			for _, asset := range release.Assets {
				if asset.Name == wantedRelease {
					latestRelease = asset.Browser_download_url
					break
				}
			}
		}
	}

	if latestVersion == "" {
		return "", "", fmt.Errorf("No version found for \"%s\"", version)
	}
	return latestVersion, latestRelease, nil
}

func printLatestVersion() {
	_, githash, _ = GetVersion()
	latest, _, err := GetLatestVersion(true)
	if err != nil {
		logger.Errorf("Couldn't find latest version: %s", err)
	} else {
		if githash == latest {
			logger.Infof("Version: (%s) is already the latest released.", latest)
		} else {
			logger.Infof("New version: (%s) available.", latest)
		}
	}
}
