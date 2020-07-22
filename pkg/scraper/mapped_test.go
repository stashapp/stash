package scraper

import (
	"testing"

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
