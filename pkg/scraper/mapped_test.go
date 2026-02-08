package scraper

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestInvalidPostProcessAction(t *testing.T) {
	yamlStr := `name: Test
performerByURL:
  - action: scrapeXPath
    scraper: performerScraper
xPathScrapers:
  performerScraper:
    performer:
      Name:
        selector: //div/a/@href
        postProcess:
          - parseDate: Jan 2, 2006
          - anything
`

	c := &Definition{}
	err := yaml.Unmarshal([]byte(yamlStr), &c)

	if err == nil {
		t.Error("expected error unmarshalling with invalid post-process action")
		return
	}
}

type feetToCMTest struct {
	in  string
	out string
}

var feetToCMTests = []feetToCMTest{
	{"", "0"},
	{"a", "0"},
	{"6", "183"},
	{"6 feet", "183"},
	{"6ft0", "183"},
	{"6ft2", "188"},
	{"6'2\"", "188"},
	{"6.2", "188"},
	{"6ft2.99", "188"},
	{"text6other2", "188"},
}

func TestFeetToCM(t *testing.T) {
	pp := postProcessFeetToCm(true)

	q := &xpathQuery{}

	for _, test := range feetToCMTests {
		assert.Equal(t, test.out, pp.Apply(context.Background(), test.in, q))
	}
}

func Test_postProcessParseDate_Apply(t *testing.T) {
	const internalDateFormat = "2006-01-02"

	unixDate := time.Date(2021, 9, 4, 1, 2, 3, 4, time.Local)

	tests := []struct {
		name  string
		arg   postProcessParseDate
		value string
		want  string
	}{
		{
			"simple",
			"2006=01=02",
			"2001=03=23",
			"2001-03-23",
		},
		{
			"today",
			"",
			"today",
			time.Now().Format(internalDateFormat),
		},
		{
			"yesterday",
			"",
			"yesterday",
			time.Now().Add(-24 * time.Hour).Format(internalDateFormat),
		},
		{
			"unix",
			"unix",
			strconv.FormatInt(unixDate.Unix(), 10),
			unixDate.Format(internalDateFormat),
		},
		{
			"invalid",
			"invalid",
			"2001=03=23",
			"2001=03=23",
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arg.Apply(ctx, tt.value, nil); got != tt.want {
				t.Errorf("postProcessParseDate.Apply() = %v, want %v", got, tt.want)
			}
		})
	}
}
