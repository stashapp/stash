package scraper

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/stashapp/stash/pkg/models"
)

const freeonesTimeout = 45 * time.Second

const freeonesScraperID = "builtin_freeones"
const freeonesName = "Freeones"

var freeonesURLs = []string{
	"freeones.com",
}

func GetFreeonesScraper() scraperConfig {
	return scraperConfig{
		ID:   freeonesScraperID,
		Name: "Freeones",
		PerformerByName: &performerByNameConfig{
			performScrape: GetPerformerNames,
		},
		PerformerByFragment: &performerByFragmentConfig{
			performScrape: GetPerformer,
		},
		PerformerByURL: []*scrapePerformerByURLConfig{
			&scrapePerformerByURLConfig{
				scrapeByURLConfig: scrapeByURLConfig{
					URL: freeonesURLs,
				},
				performScrape: GetPerformerURL,
			},
		},
	}
}

func GetPerformerNames(c scraperTypeConfig, q string) ([]*models.ScrapedPerformer, error) {
	// Request the HTML page.
	queryURL := "https://www.freeones.com/suggestions.php?q=" + url.PathEscape(q) + "&t=1"
	client := http.Client{
		Timeout: freeonesTimeout,
	}
	res, err := client.Get(queryURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	// Find the performers
	var performers []*models.ScrapedPerformer
	doc.Find(".suggestion").Each(func(i int, s *goquery.Selection) {
		name := strings.Trim(s.Text(), " ")
		p := models.ScrapedPerformer{
			Name: &name,
		}
		performers = append(performers, &p)
	})

	return performers, nil
}

func GetPerformerURL(c scraperTypeConfig, href string) (*models.ScrapedPerformer, error) {
	// if we're already in the bio page, just scrape it
	reg := regexp.MustCompile(`\/bio_.*\.php$`)
	if reg.MatchString(href) {
		return getPerformerBio(c, href)
	}

	// otherwise try to get the bio page from the url
	profileRE := regexp.MustCompile(`_links\/(.*?)\/$`)
	if profileRE.MatchString(href) {
		href = profileRE.ReplaceAllString(href, "_links/bio_$1.php")
		return getPerformerBio(c, href)
	}

	return nil, fmt.Errorf("Bio page not found in %s", href)
}

func getPerformerBio(c scraperTypeConfig, href string) (*models.ScrapedPerformer, error) {
	client := http.Client{
		Timeout: freeonesTimeout,
	}

	bioRes, err := client.Get(href)
	if err != nil {
		return nil, err
	}
	defer bioRes.Body.Close()
	if bioRes.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", bioRes.StatusCode, bioRes.Status)
	}

	// Load the HTML document
	bioDoc, err := goquery.NewDocumentFromReader(bioRes.Body)
	if err != nil {
		return nil, err
	}

	params := bioDoc.Find(".paramvalue")
	paramIndexes := getIndexes(bioDoc)

	result := models.ScrapedPerformer{}

	performerURL := bioRes.Request.URL.String()
	result.URL = &performerURL

	name := paramValue(params, paramIndexes["name"])
	result.Name = &name

	ethnicity := getEthnicity(paramValue(params, paramIndexes["ethnicity"]))
	result.Ethnicity = &ethnicity

	country := paramValue(params, paramIndexes["country"])
	result.Country = &country

	eyeColor := paramValue(params, paramIndexes["eye_color"])
	result.EyeColor = &eyeColor

	measurements := paramValue(params, paramIndexes["measurements"])
	result.Measurements = &measurements

	fakeTits := paramValue(params, paramIndexes["fake_tits"])
	result.FakeTits = &fakeTits

	careerLength := paramValue(params, paramIndexes["career_length"])
	careerRegex := regexp.MustCompile(`\([\s\S]*`)
	careerLength = careerRegex.ReplaceAllString(careerLength, "")
	careerLength = trim(careerLength)
	result.CareerLength = &careerLength

	tattoos := paramValue(params, paramIndexes["tattoos"])
	result.Tattoos = &tattoos

	piercings := paramValue(params, paramIndexes["piercings"])
	result.Piercings = &piercings

	aliases := paramValue(params, paramIndexes["aliases"])
	result.Aliases = &aliases

	birthdate := paramValue(params, paramIndexes["birthdate"])
	birthdateRegex := regexp.MustCompile(` \(\d* years old\)`)
	birthdate = birthdateRegex.ReplaceAllString(birthdate, "")
	birthdate = trim(birthdate)
	if birthdate != "Unknown" && len(birthdate) > 0 {
		t, _ := time.Parse("January _2, 2006", birthdate) // TODO
		formattedBirthdate := t.Format("2006-01-02")
		result.Birthdate = &formattedBirthdate
	}

	height := paramValue(params, paramIndexes["height"])
	heightRegex := regexp.MustCompile(`heightcm = "(.*)"\;`)
	heightMatches := heightRegex.FindStringSubmatch(height)
	if len(heightMatches) > 1 {
		result.Height = &heightMatches[1]
	}

	twitterElement := bioDoc.Find(".twitter a")
	twitterHref, _ := twitterElement.Attr("href")
	if twitterHref != "" {
		twitterURL, _ := url.Parse(twitterHref)
		twitterHandle := strings.Replace(twitterURL.Path, "/", "", -1)
		result.Twitter = &twitterHandle
	}

	instaElement := bioDoc.Find(".instagram a")
	instaHref, _ := instaElement.Attr("href")
	if instaHref != "" {
		instaURL, _ := url.Parse(instaHref)
		instaHandle := strings.Replace(instaURL.Path, "/", "", -1)
		result.Instagram = &instaHandle
	}

	return &result, nil
}

func GetPerformer(c scraperTypeConfig, scrapedPerformer models.ScrapedPerformerInput) (*models.ScrapedPerformer, error) {
	if scrapedPerformer.Name == nil {
		return nil, nil
	}

	performerName := *scrapedPerformer.Name
	queryURL := "https://www.freeones.com/search/?t=1&q=" + url.PathEscape(performerName) + "&view=thumbs"
	res, err := http.Get(queryURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	performerLink := doc.Find("div.Block3 a").FilterFunction(func(i int, s *goquery.Selection) bool {
		href, _ := s.Attr("href")
		if href == "/html/j_links/Jenna_Leigh_c/" || href == "/html/a_links/Alexa_Grace_c/" {
			return false
		}
		if strings.ToLower(s.Text()) == strings.ToLower(performerName) {
			return true
		}
		alias := s.ParentsFiltered(".babeNameBlock").Find(".babeAlias").First()
		if strings.Contains(strings.ToLower(alias.Text()), strings.ToLower(performerName)) {
			return true
		}
		return false
	})

	href, _ := performerLink.Attr("href")
	href = strings.TrimSuffix(href, "/")
	regex := regexp.MustCompile(`.+_links\/(.+)`)
	matches := regex.FindStringSubmatch(href)
	if len(matches) < 2 {
		return nil, fmt.Errorf("No matches found in %s", href)
	}

	href = strings.Replace(href, matches[1], "bio_"+matches[1]+".php", -1)
	href = "https://www.freeones.com" + href

	return getPerformerBio(c, href)

}

func getIndexes(doc *goquery.Document) map[string]int {
	var indexes = make(map[string]int)
	doc.Find(".paramname").Each(func(i int, s *goquery.Selection) {
		index := i + 1
		paramName := trim(s.Text())
		switch paramName {
		case "Babe Name:":
			indexes["name"] = index
		case "Ethnicity:":
			indexes["ethnicity"] = index
		case "Country of Origin:":
			indexes["country"] = index
		case "Date of Birth:":
			indexes["birthdate"] = index
		case "Eye Color:":
			indexes["eye_color"] = index
		case "Height:":
			indexes["height"] = index
		case "Measurements:":
			indexes["measurements"] = index
		case "Fake boobs:":
			indexes["fake_tits"] = index
		case "Career Start And End":
			indexes["career_length"] = index
		case "Tattoos:":
			indexes["tattoos"] = index
		case "Piercings:":
			indexes["piercings"] = index
		case "Aliases:":
			indexes["aliases"] = index
		}
	})
	return indexes
}

func getEthnicity(ethnicity string) string {
	switch ethnicity {
	case "Caucasian":
		return "white"
	case "Black":
		return "black"
	case "Latin":
		return "hispanic"
	case "Asian":
		return "asian"
	default:
		// #367 - unknown ethnicity shouldn't cause the entire operation to
		// fail. Just return the original string instead
		return ethnicity
	}
}

func paramValue(params *goquery.Selection, paramIndex int) string {
	i := paramIndex - 1
	if paramIndex <= 0 {
		return ""
	}
	node := params.Get(i).FirstChild
	content := trim(node.Data)
	if content != "" {
		return content
	}
	node = node.NextSibling
	if node == nil {
		return ""
	}
	return trim(node.FirstChild.Data)
}

// https://stackoverflow.com/questions/20305966/why-does-strip-not-remove-the-leading-whitespace
func trim(text string) string {
	// return text.replace(/\A\p{Space}*|\p{Space}*\z/, "");
	return strings.TrimSpace(text)
}
