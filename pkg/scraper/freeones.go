package scraper

import (
	"strings"

	"github.com/stashapp/stash/pkg/logger"
)

const freeonesScraperID = "builtin_freeones"

// 537: stolen from: https://github.com/stashapp/CommunityScrapers/blob/master/scrapers/NewFreeones.yml
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
      Name: //div[@id="search-result"]//a[@class=""]//div//p/text()
      URL:
        selector: //div[@id="search-result"]//a[@class=""]/@href
        # URL is a partial url, add the first part
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
        # URL is a partial url, add the first part
        replace:
          - regex: ^
            with: https://www.freeones.xxx
      Twitter: //div[p[text()='Follow On']]//div//a[@class='d-flex align-items-center justify-content-center mr-2 social-icons color-twitter']/@href 
      Instagram: //div[p[text()='Follow On']]//div//a[@class='d-flex align-items-center justify-content-center mr-2 social-icons color-telegram']/@href
      # need to add support for concatenating two elements or something
      Birthdate:
        selector: //div[p[text()='Personal Information']]//div//p[1]//a
        replace:
          - regex: Born On
            with:
          - regex: ","
            with:
        # reference date is: 2006/01/02
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
      Country: //div[p[text()='Personal Information']]//div//p[3]//a[last()]
      EyeColor: //div[p[text()='Eye Color']]//div//p//a//span
      Height: 
        selector: //div[p[text()='Height']]//div//p//a//span
        replace: 
          - regex: \D+[\s\S]+
            with: ""
      Measurements: //div[p[text()='Measurements']]//div[@class='p-3']//p
      FakeTits: //div[p[text()='Fake Boobs']]//div[@class='p-3']//p
      # nbsp; screws up the parsing, so use contains instead
      CareerLength: 
        selector: //div[p[text()='career']]//div//div[@class='timeline-horizontal mb-3']//div//p[@class='m-0']
        concat: "-"
        replace: 
          - regex: -\w+-\w+-\w+-\w+-\w+$
            with: ""
      Aliases: //div[p[text()='Aliases']]//div//p[@class='mb-0 text-center']
      Tattoos: //div[p[text()='Tattoos']]//div//p[@class='mb-0 text-center']
      Piercings: //div[p[text()='Piercings']]//div//p[@class='mb-0 text-center']
      Image:
        selector: //div[@class='profile-image-large']//a/img/@src
        # URL is a partial url, add the first part
`

func GetFreeonesScraper() scraperConfig {
	yml := freeonesScraperConfig

	scraper, err := loadScraperFromYAML(freeonesScraperID, strings.NewReader(yml))
	if err != nil {
		logger.Fatalf("Error loading builtin freeones scraper: %s", err.Error())
	}

	return *scraper
}
