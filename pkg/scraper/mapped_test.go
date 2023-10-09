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

	c := &config{}
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

type dimensionToMetricTest struct {
	in  string
	out string
}

var dimensionToMetricTests = []dimensionToMetricTest{
	// Not actual measurements
	{"", "0"},
	{"a", "0"},
	// Inches
	{"9", "23"},
	{"9in", "23"},
	{"9 in", "23"},
	{"9.5", "24"},
	{"9.5in", "24"},
	{"9.5 in", "24"},
	// Feet and inches
	{"5 feet", "152"},
	{"5ft1", "155"},
	{"5 feet 2 inches", "157"},
	{"5'3", "160"},
	{"5'4\"", "163"},
	{"5ft5.99", "168"},
	{"5foot 7in", "170"},
	{"5 foot 8 inches", "173"},
	// Already metric
	{"174", "174"},
	{"175cm", "175"},
	{"176 cm", "176"},
	{"177.0", "177"},
	{"178.2", "178"},
	{"179.99", "179"},
	{"1.8 m", "180"},
	{"1.81m", "181"},
}

func TestDimensionToMetric(t *testing.T) {
	pp := postProcessDimensionToMetric(true)

	q := &xpathQuery{}

	for _, test := range dimensionToMetricTests {
		t.Logf("Testing %s", test.in)
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
