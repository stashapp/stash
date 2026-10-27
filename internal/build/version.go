// Package build provides the version information for the application.
package build

import (
	"regexp"
)

var version string
var buildstamp string
var githash string
var officialBuild string

// updateRepo is the GitHub "owner/repo" that the built-in update check queries.
// It can be overridden at build time via -ldflags, allowing forks to point the
// update check at their own releases. It defaults to the upstream repository.
var updateRepo string

const defaultUpdateRepo = "stashapp/stash"

func Version() (string, string, string) {
	return version, githash, buildstamp
}

// UpdateRepo returns the GitHub "owner/repo" used by the update check.
func UpdateRepo() string {
	if updateRepo == "" {
		return defaultUpdateRepo
	}
	return updateRepo
}

func VersionString() string {
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
	if IsOfficial() {
		versionString += " - Official Build"
	} else {
		versionString += " - Unofficial Build"
	}
	if buildstamp != "" {
		versionString += " - " + buildstamp
	}
	return versionString
}

func IsOfficial() bool {
	return officialBuild == "true"
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
