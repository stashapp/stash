package manager

import (
	"github.com/stashapp/stash/pkg/logger"
	"regexp"
	"strings"
)

func excludeFiles(files []string, patterns []string) ([]string, int) {
	if patterns == nil {
		logger.Infof("No exclude patterns in config")
		return files, 0
	} else {
		var results []string
		var exclCount int

		fileRegexps := generateRegexps(patterns)

		if len(fileRegexps) == 0 {
			logger.Infof("Excluded 0 files from scan")
			return files, 0
		}

		for i := 0; i < len(files); i++ {
			if matchFileSimple(files[i], fileRegexps) {
				logger.Infof("File matched pattern. Excluding:\"%s\"", files[i])
				exclCount++
			} else {

				//if pattern doesn't match add file to list
				results = append(results, files[i])
			}
		}
		logger.Infof("Excluded %d file(s) from scan", exclCount)
		return results, exclCount
	}
}

func matchFile(file string, patterns []string) bool {
	if patterns == nil {
		logger.Infof("No exclude patterns in config.")

	} else {
		fileRegexps := generateRegexps(patterns)

		if len(fileRegexps) == 0 {
			return false
		}

		for _, regPattern := range fileRegexps {
			if regPattern.Match([]byte(strings.ToLower(file))) {
				return true
			}

		}
	}

	return false
}

func generateRegexps(patterns []string) []*regexp.Regexp {

	var fileRegexps []*regexp.Regexp

	for _, pattern := range patterns {
		reg, err := regexp.Compile(strings.ToLower(pattern))
		if err != nil {
			logger.Errorf("Exclude :%v", err)
		} else {
			fileRegexps = append(fileRegexps, reg)
		}
	}

	if len(fileRegexps) == 0 {
		return nil
	} else {
		return fileRegexps
	}

}

func matchFileSimple(file string, regExps []*regexp.Regexp) bool {
	for _, regPattern := range regExps {
		if regPattern.Match([]byte(strings.ToLower(file))) {
			return true
		}
	}
	return false
}
