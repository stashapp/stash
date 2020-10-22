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
  queryURL: https://www.freeones.xxx/babes?q={}&v=teasers&s=relevance&l=96&m%5BcanPreviewFeatures%5D=0
  scraper: performerSearch
performerByURL:
  - action: scrapeXPath
    url:
      - freeones.xxx
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
              with: https://www.freeones.xxx
            - regex: $
              with: /profile

  performerScraper:
    performer:
      Name: 
        selector: //h1
        postProcess:
          - replace:
            - regex: \sBio\s*$
              with: ""
      URL:
        selector: //a[span[text()="Profile"]]/@href
        postProcess:
          - replace:
            - regex: ^
              with: https://www.freeones.xxx
      Twitter: //a[contains(@href,'twitter.com/')]/@href
      Instagram: //a[contains(@href,'instagram.com/')]/@href
      Birthdate:
        selector: //div[p[text()='Personal Information']]//span[contains(text(),'Born On')]
        postProcess:
          - replace:
            - regex: Born On
              with:
          - parseDate: January 2, 2006
      Ethnicity:
        selector: //div[p[text()='Ethnicity']]//a[@data-test="link_ethnicity"]
        postProcess:
          - map:
              Asian: asian
              Caucasian: white
              Black: black
              Latin: hispanic
      Country: //div[p[text()='Personal Information']]//a[@data-test="link-country"]
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
              Fake: Yes
              Natural: No
      CareerLength:
        selector: //div[p[text()='career']]//div[contains(@class,'timeline-horizontal')]//p[@class='m-0']
        concat: "-"
      Aliases: //p[text()='Aliases']/following-sibling::div/p
      Tattoos: //span[text()='Tattoos']/following-sibling::span/span
      Piercings: //span[text()='Piercings']/following-sibling::span/span
      Image:
        selector: //div[@class='profile-image-container']//a/img/@src
      Gender:
        fixed: Female
# Last updated October 21, 2020
`

func getFreeonesScraper() config {
	yml := freeonesScraperConfig

	scraper, err := loadScraperFromYAML(FreeonesScraperID, strings.NewReader(yml))
	if err != nil {
		logger.Fatalf("Error loading builtin freeones scraper: %s", err.Error())
	}

	return *scraper
}
