package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"runtime"
	"time"

	"golang.org/x/sys/cpu"

	"github.com/stashapp/stash/pkg/logger"
)

//we use the github REST V3 API as no login is required
const apiReleases string = "https://api.github.com/repos/stashapp/stash/releases"
const apiTags string = "https://api.github.com/repos/stashapp/stash/tags"
const apiAcceptHeader string = "application/vnd.github.v3+json"
const developmentTag string = "latest_develop"
const defaultSHLength int = 7 // default length of SHA short hash returned by <git rev-parse --short HEAD>

// ErrNoVersion indicates that no version information has been embedded in the
// stash binary
var ErrNoVersion = errors.New("no stash version")

var stashReleases = func() map[string]string {
	return map[string]string{
		"darwin/amd64":  "stash-osx",
		"darwin/arm64":  "stash-osx-applesilicon",
		"linux/amd64":   "stash-linux",
		"windows/amd64": "stash-win.exe",
		"linux/arm":     "stash-pi",
		"linux/arm64":   "stash-linux-arm64v8",
		"linux/armv7":   "stash-linux-arm32v7",
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

func makeGithubRequest(url string, output interface{}) error {
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Accept", apiAcceptHeader) // gh api recommendation , send header with api version
	response, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("Github API request failed: %s", err)
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Github API request failed: %s", response.Status)
	}

	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("Github API read response failed: %s", err)
	}

	err = json.Unmarshal(data, output)
	if err != nil {
		return fmt.Errorf("Unmarshalling Github API response failed: %s", err)
	}

	return nil
}

// GetLatestVersion gets latest version (git commit hash) from github API
// If running a build from the "master" branch, then the latest full release
// is used, otherwise it uses the release that is tagged with "latest_develop"
// which is the latest pre-release build.
func GetLatestVersion(shortHash bool) (latestVersion string, latestRelease string, err error) {

	arch := runtime.GOARCH                                                                    // https://en.wikipedia.org/wiki/Comparison_of_ARM_cores
	isARMv7 := cpu.ARM.HasNEON || cpu.ARM.HasVFPv3 || cpu.ARM.HasVFPv3D16 || cpu.ARM.HasVFPv4 // armv6 doesn't support any of these features
	if arch == "arm" && isARMv7 {
		arch = "armv7"
	}

	platform := fmt.Sprintf("%s/%s", runtime.GOOS, arch)
	wantedRelease := stashReleases()[platform]

	version, _, _ := GetVersion()
	if version == "" {
		return "", "", ErrNoVersion
	}

	// if the version is suffixed with -x-xxxx, then we are running a development build
	usePreRelease := false
	re := regexp.MustCompile(`-\d+-g\w+$`)
	if re.MatchString(version) {
		usePreRelease = true
	}

	url := apiReleases
	if !usePreRelease {
		// just get the latest full release
		url += "/latest"
	} else {
		// get the release tagged with the development tag
		url += "/tags/" + developmentTag
	}

	release := githubReleasesResponse{}
	err = makeGithubRequest(url, &release)

	if err != nil {
		return "", "", err
	}

	if release.Prerelease == usePreRelease {
		latestVersion = getReleaseHash(release, shortHash, usePreRelease)

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

func getReleaseHash(release githubReleasesResponse, shortHash bool, usePreRelease bool) string {
	shaLength := len(release.Target_commitish)
	// the /latest API call doesn't return the hash in target_commitish
	// also add sanity check in case Target_commitish is not 40 characters
	if !usePreRelease || shaLength != 40 {
		return getShaFromTags(shortHash, release.Tag_name)
	}

	if shortHash {
		last := defaultSHLength                                // default length of git short hash
		_, gitShort, _ := GetVersion()                         // retrieve it to check actual length
		if len(gitShort) > last && len(gitShort) < shaLength { // sometimes short hash is longer
			last = len(gitShort)
		}
		return release.Target_commitish[0:last]
	}

	return release.Target_commitish
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

// get sha from the github api tags endpoint
// returns the sha1 hash/shorthash or "" if something's wrong
func getShaFromTags(shortHash bool, name string) string {
	url := apiTags
	tags := []githubTagResponse{}
	err := makeGithubRequest(url, &tags)

	if err != nil {
		logger.Errorf("Github Tags Api %v", err)
		return ""
	}
	_, gitShort, _ := GetVersion() // retrieve short hash to check actual length

	for _, tag := range tags {
		if tag.Name == name {
			shaLength := len(tag.Commit.Sha)
			if shaLength != 40 {
				return ""
			}
			if shortHash {
				last := defaultSHLength                                // default length of git short hash
				if len(gitShort) > last && len(gitShort) < shaLength { // sometimes short hash is longer
					last = len(gitShort)
				}
				return tag.Commit.Sha[0:last]
			}

			return tag.Commit.Sha
		}
	}

	return ""
}
