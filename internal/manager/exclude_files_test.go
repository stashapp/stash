package manager

import (
	"fmt"
	"testing"

	"github.com/stashapp/stash/pkg/logger"
)

var excludeTestFilenames = []string{
	"/stash/videos/filename.mp4",
	"/stash/videos/new filename.mp4",
	"filename sample.mp4",
	"/stash/videos/exclude/not wanted.webm",
	"/stash/videos/exclude/not wanted2.webm",
	"/somewhere/trash/not wanted.wmv",
	"/disk2/stash/videos/exclude/!!wanted!!.avi",
	"/disk2/stash/videos/xcl/not wanted.avi",
	"/stash/videos/partial.file.001.webm",
	"/stash/videos/partial.file.002.webm",
	"/stash/videos/partial.file.003.webm",
	"/stash/videos/sample file sample.mkv",
	"/stash/videos/.ckRVp1/.still_encoding.mp4",
	"c:\\stash\\videos\\exclude\\filename  windows.mp4",
	"c:\\stash\\videos\\filename  windows.mp4",
	"\\\\network\\videos\\filename  windows network.mp4",
	"\\\\network\\share\\windows network wanted.mp4",
	"\\\\network\\share\\windows network wanted sample.mp4",
	"\\\\network\\private\\windows.network.skip.mp4",
	"/stash/videos/a5.mp4",
	"/stash/videos/mIxEdCaSe.mp4"}

var excludeTests = []struct {
	testPattern []string
	expected    int
}{
	{[]string{"sample\\.mp4$", "trash", "\\.[\\d]{3}\\.webm$"}, 6}, // generic
	{[]string{"no_match\\.mp4"}, 0},                                // no match
	{[]string{"^/stash/videos/exclude/", "/videos/xcl/"}, 3},       // linux
	{[]string{"/\\.[[:word:]]+/"}, 1},                              // linux hidden dirs (handbrake unraid issue?)
	{[]string{"c:\\\\stash\\\\videos\\\\exclude"}, 1},              // windows
	{[]string{"\\/[/invalid"}, 0},                                  // invalid pattern
	{[]string{"\\/[/invalid", "sample\\.[[:alnum:]]+$"}, 3},        // invalid pattern but continue
	{[]string{"^\\\\\\\\network"}, 4},                              // windows net share
	{[]string{"\\\\private\\\\"}, 1},                               // windows net share
	{[]string{"\\\\private\\\\", "sample\\.mp4"}, 3},               // windows net share
	{[]string{"\\D\\d\\.mp4"}, 1},                                  // validates that \D doesn't get converted to lowercase \d
	{[]string{"mixedcase\\.mp4"}, 1},                               // validates we can match the mixed case file
	{[]string{"MIXEDCASE\\.mp4"}, 1},                               // validates we can match the mixed case file
	{[]string{"(?i)MIXEDCASE\\.mp4"}, 1},                           // validates we can match the mixed case file without adding another (?i) to it
}

func TestExcludeFiles(t *testing.T) {
	for _, test := range excludeTests {
		err := runExclude(excludeTestFilenames, test.testPattern, test.expected)
		if err != nil {
			t.Error(err)
		}
	}
}

func runExclude(filenames []string, patterns []string, expCount int) error {

	files, count := excludeFiles(filenames, patterns)

	if count != expCount {
		return fmt.Errorf("Was expecting %d, found %d", expCount, count)
	}
	if len(files) != len(filenames)-expCount {
		return fmt.Errorf("Returned list should have %d files, not %d ", len(filenames)-expCount, len(files))
	}

	return nil
}

func TestMatchFile(t *testing.T) {
	for _, test := range excludeTests {
		err := runMatch(excludeTestFilenames, test.testPattern, test.expected)
		if err != nil {
			t.Error(err)
		}
	}
}

func runMatch(filenames []string, patterns []string, expCount int) error {
	count := 0
	for _, file := range filenames {
		if matchFile(file, patterns) {
			logger.Infof("File \"%s\" matched pattern\n", file)
			count++
		}
	}
	if count != expCount {
		return fmt.Errorf("Was expecting %d, found %d", expCount, count)
	}

	return nil
}
