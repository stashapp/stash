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
              - regex: (.+)\sidentifies.+
                with: $1
      URL: //link[@rel="alternate" and @hreflang="x-default"]/@href
      Twitter: //form//a[contains(@href,'twitter.com/')]/@href
      Instagram: //form//a[contains(@href,'instagram.com/')]/@href
      Birthdate:
        selector: //span[@data-test="link_span_dateOfBirth"]/text()
        postProcess:
          - parseDate: January 2, 2006
      Ethnicity:
        selector: //span[@data-test="link_span_ethnicity"]
        postProcess:
          - map:
              Asian: Asian
              Caucasian: White
              Black: Black
              Latin: Hispanic
      Country:
        selector: //a[@data-test="link_placeOfBirth"][contains(@href, 'country')]/span/text()
        postProcess:
          - map:
              United States: "USA"
      EyeColor: //span[text()='Eye Color:']/following-sibling::span/a/span/text()
      Height:
        selector: //span[text()='Height:']/following-sibling::span/a
        postProcess:
          - replace:
            - regex: \scm
              with: ""
          - map:
              Unknown: ""
      Measurements:
        selector: //span[(@data-test='link_span_bra') or (@data-test='link_span_waist') or (@data-test='link_span_hip')]
        concat: " - "
        postProcess:
          - replace:
              - regex: \sIn
                with: ""
          - map:
              Unknown: ""
      FakeTits:
        selector: //span[text()='Boobs:']/following-sibling::span/a
        postProcess:
          - map:
              Unknown: ""
              Fake: "Yes"
              Natural: "No"
      CareerLength:
        selector: //div[contains(@class,'timeline-horizontal')]//p[@class='m-0']
        concat: "-"
      Aliases:
        selector: //span[@data-test='link_span_aliases']/text()
        concat: ", "
      Tattoos:
        selector: //span[text()='Tattoo locations:']/following-sibling::span
        postProcess:
          - map:
              Unknown: ""
      Piercings:
        selector: //span[text()='Piercing locations:']/following-sibling::span
        postProcess:
          - map:
              Unknown: ""
      Image:
        selector: //div[contains(@class,'image-container')]//a/img/@src
      Gender:
        selector: //h1/*[1]/*[1]/text()Add commentMore actions
        postProcess:
          - replace:
            - regex: .+ identifies as (.+)
              with: $1
      DeathDate:
        selector: //div[contains(text(),'Passed away on')]
        postProcess:
          - replace:
              - regex: Passed away on (.+) at the age of \d+
                with: $1
          - parseDate: January 2, 2006
      HairColor: //span[@data-test="link_span_hair_color"]
      Weight:
        selector: //span[@data-test="link_span_weight"]
        postProcess:
          - replace:
            - regex: \skg
              with: ""

# Last Updated June 22, 2025
`

func getFreeonesScraper(globalConfig GlobalConfig) scraper {
	yml := freeonesScraperConfig

	c, err := loadConfigFromYAML(FreeonesScraperID, strings.NewReader(yml))
	if err != nil {
		logger.Fatalf("Error loading builtin freeones scraper: %s", err.Error())
	}

	return scraperFromDefinition(*c, globalConfig)
}
