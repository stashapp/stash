package manager

import (
	"fmt"
	"testing"
)

func TestExcludeFiles(t *testing.T) {

	filenames := []string{
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
		"/stash/videos/sample file.mkv",
		"/stash/videos/.ckRVp1/.still_encoding.mp4",
		"c:\\stash\\videos\\exclude\\filename  windows.mp4",
		"c:\\stash\\videos\\filename  windows.mp4",
		"c:\\stash\\videos\\filename  windows sample.mp4"}

	var excludeTests = []struct {
		testPattern []string
		expected    int
	}{
		{[]string{"sample\\.mp4$", "trash", "\\.[\\d]{3}\\.webm$"}, 6}, //generic
		{[]string{"no_match\\.mp4"}, 0},                                //no match
		{[]string{"^/stash/videos/exclude/", "/videos/xcl/"}, 3},       //linux
		{[]string{"/\\.[[:word:]]+/"}, 1},                              //linux hidden dirs (handbrake unraid issue?)
		{[]string{"c:\\\\stash\\\\videos\\\\exclude"}, 1},              //windows
		{[]string{"\\/[/invalid"}, 0},                                  //invalid pattern
	}
	for _, test := range excludeTests {
		err := runExclude(filenames, test.testPattern, test.expected)
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
