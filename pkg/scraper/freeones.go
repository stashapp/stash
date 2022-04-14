package scraper

import (
	"strings"

	"github.com/stashapp/stash/pkg/logger"
)

// FreeonesScraperID is the scraper ID for the built-in Freeones scraper
const FreeonesScraperID = "builtin_freeones"

// 537: stolen from: https://github.com/stashapp/CommunityScrapers/blob/master/scrapers/FreeonesCommunity.yml
const freeonesScraperConfig = `
name: Freeones
performerByName:
  action: scrapeXPath
  queryURL: https://www.freeones.com/babes?q={}&v=teasers&s=relevance&l=96&m%5BcanPreviewFeatures%5D=0
  scraper: performerSearch
performerByURL:
  - action: scrapeXPath
    url:
      - freeones.xxx
      - freeones.com
    scraper: performerScraper

xPathScrapers:
  performerSearch:
    performer:
      Name: //div[@id="search-result"]//p[@data-test="subject-name"]/text()
      URL:
        selector: //div[@id="search-result"]//div[@data-test="teaser-subject"]/a/@href
        postProcess:
          - replace:
              - regex: ^
                with: https://www.freeones.com
              - regex: /feed$
                with: /bio

  performerScraper:
    performer:
      Name:
        selector: //h1
        postProcess:
          - replace:
              - regex: \sBio\s*$
                with: ""
      URL: //link[@rel="alternate" and @hreflang="x-default"]/@href
      Twitter: //a[not(starts-with(@href,'https://twitter.com/FreeOnes'))][contains(@href,'twitter.com/')]/@href
      Instagram: //a[contains(@href,'instagram.com/')]/@href
      Birthdate:
        selector: //span[contains(text(),'Born On')]
        postProcess:
          - replace:
              - regex: Born On
                with:
          - parseDate: January 2, 2006
      Ethnicity:
        selector: //a[@data-test="link_ethnicity"]/span/text()
        postProcess:
          - map:
              Asian: Asian
              Caucasian: White
              Black: Black
              Latin: Hispanic
      Country: //a[@data-test="link-country"]/span/text()
      EyeColor: //span[text()='Eye Color']/following-sibling::span/a
      Height:
        selector: //span[text()='Height']/following-sibling::span/a
        postProcess:
          - replace:
              - regex: \D+[\s\S]+
                with: ""
          - map:
              Unknown: ""
      Measurements:
        selector: //span[text()='Measurements']/following-sibling::span/span/a
        concat: " - "
        postProcess:
          - map:
              Unknown: ""
      FakeTits:
        selector: //span[text()='Boobs']/following-sibling::span/a
        postProcess:
          - map:
              Unknown: ""
              Fake: "Yes"
              Natural: "No"
      CareerLength:
        selector: //div[contains(@class,'timeline-horizontal')]//p[@class='m-0']
        concat: "-"
      Aliases: //p[@data-test='p_aliases']/text()
      Tattoos:
        selector: //span[text()='Tattoos']/following-sibling::span/span
        postProcess:
          - map:
              Unknown: ""
      Piercings:
        selector: //span[text()='Piercings']/following-sibling::span/span
        postProcess:
          - map:
              Unknown: ""
      Image:
        selector: //div[contains(@class,'image-container')]//a/img/@src
      Gender:
        fixed: "Female"
      Details: //div[@data-test="biography"]
      DeathDate:
        selector: //div[contains(text(),'Passed away on')]
        postProcess:
          - replace:
              - regex: Passed away on (.+) at the age of \d+
                with: $1
          - parseDate: January 2, 2006
      HairColor: //span[text()='Hair Color']/following-sibling::span/a
      Weight:
        selector: //span[text()='Weight']/following-sibling::span/a
        postProcess:
        - replace:
            - regex: \D+[\s\S]+
              with: ""

# Last updated April 13, 2021
`

func getFreeonesScraper(globalConfig GlobalConfig) scraper {
	yml := freeonesScraperConfig

	c, err := loadConfigFromYAML(FreeonesScraperID, strings.NewReader(yml))
	if err != nil {
		logger.Fatalf("Error loading builtin freeones scraper: %s", err.Error())
	}

	return newGroupScraper(*c, globalConfig)
}
