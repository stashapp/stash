package scraper

import (
	"context"
	"testing"

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
