package scraper

import (
	"strings"

	"github.com/stashapp/stash/pkg/logger"
)

const freeonesScraperID = "builtin_freeones"

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
      - https://www.freeones.xxx
    scraper: performerScraper

xPathScrapers:
  performerSearch:
    performer:
      Name: //div[@id="search-result"]//p[@data-test="subject-name"]/text()
      URL:
        selector: //div[@id="search-result"]//div[@data-test="teaser-subject"]/a/@href
        replace:
          - regex: ^
            with: https://www.freeones.xxx
          - regex: $
            with: /profile

  performerScraper:
    performer:
      Name: //h1
      URL:
        selector: //a[span[text()="Profile"]]/@href
        replace:
          - regex: ^
            with: https://www.freeones.xxx
      Twitter: //div[p[text()='Follow On']]//div//a[@class='d-flex align-items-center justify-content-center m-2 social-icons color-twitter']/@href
      Instagram: //div[p[text()='Follow On']]//div//a[@class='d-flex align-items-center justify-content-center m-2 social-icons color-telegram']/@href
      Birthdate:
        selector: //div[p[text()='Personal Information']]//div//p/a/span[contains(text(),'Born On')]
        replace:
          - regex: Born On
            with:
          - regex: ","
            with:
        parseDate: January 2 2006
      Ethnicity:
        selector: //div[p[text()='Ethnicity']]//div//p[@class='mb-0 text-center']
        replace:
          - regex: Asian
            with: "asian"
          - regex: Caucasian
            with: "white"
          - regex: Black
            with: "black"
          - regex: Latin
            with: "hispanic"
      Country: //div[p[text()='Personal Information']]//div//p//a[@data-test="link-country"]
      EyeColor: //span[@data-test="link_span_eye_color"]
      Height:
        selector: //span[@data-test="link_span_height"]
        replace:
          - regex: \D+[\s\S]+
            with: ""
      Measurements:
        selector: //span[@data-test="p-measurements"]//a/span
        concat: " - "
        replace:
          - regex: Unknown
            with:
      FakeTits:
        selector: //span[@data-test='link_span_boobs']
        replace:
          - regex: Unknown
            with:
          - regex: Fake
            with: "Yes"
          - regex: Natural
            with: "No"
      CareerLength:
        selector: //div[p[text()='career']]//div//div[@class='timeline-horizontal mb-3']//div//p[@class='m-0']
        concat: "-"
        replace:
          - regex: -\w+-\w+-\w+-\w+-\w+$
            with: ""
      Aliases: //div[p[text()='Aliases']]//div//p[@class='mb-0 text-center']
      Tattoos: //span[@data-test="p_has_tattoos"]|//span[@cdata-test="p_has_tattoos"]
      Piercings: //span[@data-test="p_has_piercings"]
      Image:
        selector: //div[@class='profile-image-container']//a/img/@src
      Gender:
        selector: //meta[@name="language"]/@name
        replace:
          - regex: language
            with: "Female"
# Last updated June 15, 2020
`

func GetFreeonesScraper() scraperConfig {
	yml := freeonesScraperConfig

	scraper, err := loadScraperFromYAML(freeonesScraperID, strings.NewReader(yml))
	if err != nil {
		logger.Fatalf("Error loading builtin freeones scraper: %s", err.Error())
	}

	return *scraper
}
