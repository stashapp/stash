package scraper

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/antchfx/htmlquery"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

// adapted from https://www.freeones.com/html/m_links/bio_Mia_Malkova.php
const htmlDoc1 = `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en" dir="ltr">
	<head>
		<title>Freeones:  Mia Malkova Biography</title>
	</head>
	<body data-babe="Mia Malkova">
		<div class="ContentBlock Block1">
			<div class="ContentBlockBody" style="padding: 0px;">
				<table id="biographyTable" border="0" cellspacing="0" cellpadding="0" width="100%">
					<tbody>
						<tr>
							<td class="paramname">
								<div><b>Babe Name:</b></div>
							</td>
							<td class="paramvalue">
								<a href="/html/m_links/Mia_Malkova/">Mia Malkova</a>&nbsp;
							</td>
						</tr>
						<tr>
							<td class="paramname">
								<div><b>Profession:</b></div>
							</td>
							<td class="paramvalue">Porn Star
							</td>
						</tr>
						<tr>
							<td class="paramname">
								<b>Ethnicity:</b>
							</td>
							<td class="paramvalue">
								Caucasian&nbsp;
							</td>
						</tr>
						<tr>
							<td class="paramname">
								<b>Country of Origin:</b>
							</td>
							<td class="paramvalue">
								
								<span class="country-us">
								
									United States
								<span>
							</span></span></td>
						</tr>
						<tr>
							<td class="paramname">
								<b>Date of Birth:</b>
							</td>
							<td class="paramvalue">
								July 1, 1992 (27 years old)&nbsp;
							</td>
						</tr>
						<tr>
							<td class="paramname">
								<b>Aliases:</b>
							</td>
							<td class="paramvalue">
								Mia Bliss, Madison Clover, Madison Swan, Mia Mountain, Jessica&nbsp;
							</td>
						</tr>
						<tr>
							<td class="paramname">
								<b>Eye Color:</b>
							</td>
							<td class="paramvalue">
								Hazel&nbsp;
							</td>
						</tr>
						<tr>
							<td class="paramname">
								<b>Hair Color:</b>
							</td>
							<td class="paramvalue">
								Blonde&nbsp;
							</td>
						</tr>
						<tr>
							<td class="paramname">
								<b>Height:</b>
							</td>
							<td class="paramvalue">
								5ft7
							</td>
						</tr>
						<tr>
							<td class="paramname">
								<b>Measurements:</b>
							</td>
							<td class="paramvalue">
								34C-26-36
							</td>
						</tr>
						<tr>
							<td class="paramname">
								<b>Fake boobs:</b>
							</td>
							<td class="paramvalue">
								No&nbsp;
							</td>
						</tr>
						<tr>
							<td class="paramname">
								<b>Career Start And End</b>
							</td>
							<td class="paramvalue">
								2012 - 2019
								(7 Years In The Business)
							</td>
						</tr>
						<tr>
							<td class="paramname">
								<b>Tattoos:</b>
							</td>
							<td class="paramvalue">
								None&nbsp;
							</td>
						</tr>
						<tr>
							<td class="paramname">
								<b>Piercings:</b>
							</td>
							<td class="paramvalue">
								<!-- None -->;
							</td>
						</tr>
						<tr>
							<td class="paramname">
								<div><b>Social Network Links:</b></div>
							</td>
							<td class="paramvalue">
								<ul id="socialmedia">
									<li class="twitter"><a href="https://twitter.com/MiaMalkova" target="_blank" alt="Mia Malkova Twitter" title="Mia Malkova Twitter">Twitter</a></li>
									<li class="facebook"><a href="https://www.facebook.com/MiaMalcove" target="_blank" alt="Mia Malkova Facebook" title="Mia Malkova Facebook">Facebook</a></li>
									<li class="youtube"><a href="https://www.youtube.com/channel/UCEPR0sZKa_ScMoyhemfB7nA" target="_blank" alt="Mia Malkova YouTube" title="Mia Malkova YouTube">YouTube</a></li>
									<li class="instagram"><a href="https://www.instagram.com/mia_malkova/" target="_blank" alt="Mia Malkova Instagram" title="Mia Malkova Instagram">Instagram</a></li>
								</ul>
							</td>
						</tr>
					</tbody>
				</table>
			</div>
		</div>
	</body>
</html>
`

func makeCommonXPath(attr string) string {
	return `//table[@id="biographyTable"]//tr/td[@class="paramname"]//b[text() = '` + attr + `']/ancestor::tr/td[@class="paramvalue"]`
}

func makeSimpleAttrConfig(str string) mappedScraperAttrConfig {
	return mappedScraperAttrConfig{
		Selector: str,
	}
}

func makeReplaceRegex(regex string, with string) mappedRegexConfig {
	ret := mappedRegexConfig{
		Regex: regex,
		With:  with,
	}

	return ret
}

func makeXPathConfig() mappedPerformerScraperConfig {
	config := mappedPerformerScraperConfig{
		mappedConfig: make(mappedConfig),
	}

	config.mappedConfig["Name"] = makeSimpleAttrConfig(makeCommonXPath("Babe Name:") + `/a`)
	config.mappedConfig["Ethnicity"] = makeSimpleAttrConfig(makeCommonXPath("Ethnicity:"))
	config.mappedConfig["Aliases"] = makeSimpleAttrConfig(makeCommonXPath("Aliases:"))
	config.mappedConfig["EyeColor"] = makeSimpleAttrConfig(makeCommonXPath("Eye Color:"))
	config.mappedConfig["Measurements"] = makeSimpleAttrConfig(makeCommonXPath("Measurements:"))
	config.mappedConfig["FakeTits"] = makeSimpleAttrConfig(makeCommonXPath("Fake boobs:"))
	config.mappedConfig["Tattoos"] = makeSimpleAttrConfig(makeCommonXPath("Tattoos:"))
	config.mappedConfig["Piercings"] = makeSimpleAttrConfig(makeCommonXPath("Piercings:") + "/comment()")

	// special handling for birthdate
	birthdateAttrConfig := makeSimpleAttrConfig(makeCommonXPath("Date of Birth:"))

	var birthdateReplace mappedRegexConfigs
	// make this leave the trailing space to test existing scrapers that do so
	birthdateReplace = append(birthdateReplace, makeReplaceRegex(`\(.* years old\)`, ""))

	birthdateReplaceAction := postProcessReplace(birthdateReplace)
	birthdateParseDate := postProcessParseDate("January 2, 2006") // "July 1, 1992 (27 years old)&nbsp;"
	birthdateAttrConfig.postProcessActions = []postProcessAction{
		&birthdateReplaceAction,
		&birthdateParseDate,
	}
	config.mappedConfig["Birthdate"] = birthdateAttrConfig

	// special handling for career length
	// no colon in attribute header
	careerLengthAttrConfig := makeSimpleAttrConfig(makeCommonXPath("Career Start And End"))

	var careerLengthReplace mappedRegexConfigs
	careerLengthReplace = append(careerLengthReplace, makeReplaceRegex(`\s+\(.*\)`, ""))
	careerLengthReplaceAction := postProcessReplace(careerLengthReplace)
	careerLengthAttrConfig.postProcessActions = []postProcessAction{
		&careerLengthReplaceAction,
	}

	config.mappedConfig["CareerLength"] = careerLengthAttrConfig

	// use map post-process action for gender
	genderConfig := makeSimpleAttrConfig(makeCommonXPath("Profession:"))
	genderMapAction := make(postProcessMap)
	genderMapAction["Porn Star"] = "Female"
	genderConfig.postProcessActions = []postProcessAction{
		&genderMapAction,
	}

	config.mappedConfig["Gender"] = genderConfig

	// use fixed for height
	config.mappedConfig["Country"] = mappedScraperAttrConfig{
		Fixed: "United States",
	}

	heightConfig := makeSimpleAttrConfig(makeCommonXPath("Height:"))
	heightConvAction := postProcessFeetToCm(true)
	heightConfig.postProcessActions = []postProcessAction{
		&heightConvAction,
	}
	config.mappedConfig["Height"] = heightConfig

	return config
}

func verifyField(t *testing.T, expected string, actual *string, field string) {
	t.Helper()

	if actual == nil || *actual != expected {
		if actual == nil {
			t.Errorf("Expected %s to be set to %s, instead got nil", field, expected)
		} else {
			t.Errorf("Expected %s to be set to %s, instead got %s", field, expected, *actual)
		}
	}
}

func TestScrapePerformerXPath(t *testing.T) {
	reader := strings.NewReader(htmlDoc1)
	doc, err := htmlquery.Parse(reader)

	if err != nil {
		t.Errorf("Error loading document: %s", err.Error())
		return
	}

	xpathConfig := makeXPathConfig()

	scraper := mappedScraper{
		Performer: &xpathConfig,
	}

	q := &xpathQuery{
		doc: doc,
	}

	performer, err := scraper.scrapePerformer(q)

	if err != nil {
		t.Errorf("Error scraping performer: %s", err.Error())
		return
	}

	const performerName = "Mia Malkova"
	const ethnicity = "Caucasian"
	const country = "United States"
	const birthdate = "1992-07-01"
	const aliases = "Mia Bliss, Madison Clover, Madison Swan, Mia Mountain, Jessica"
	const eyeColor = "Hazel"
	const measurements = "34C-26-36"
	const fakeTits = "No"
	const careerLength = "2012 - 2019"
	const tattoos = "None"
	const piercings = "<!-- None -->"
	const gender = "Female"
	const height = "170"

	verifyField(t, performerName, performer.Name, "Name")
	verifyField(t, gender, performer.Gender, "Gender")
	verifyField(t, ethnicity, performer.Ethnicity, "Ethnicity")
	verifyField(t, country, performer.Country, "Country")

	verifyField(t, birthdate, performer.Birthdate, "Birthdate")

	verifyField(t, aliases, performer.Aliases, "Aliases")
	verifyField(t, eyeColor, performer.EyeColor, "EyeColor")
	verifyField(t, measurements, performer.Measurements, "Measurements")
	verifyField(t, fakeTits, performer.FakeTits, "FakeTits")

	verifyField(t, careerLength, performer.CareerLength, "CareerLength")

	verifyField(t, tattoos, performer.Tattoos, "Tattoos")
	verifyField(t, piercings, performer.Piercings, "Piercings")
	verifyField(t, height, performer.Height, "Height")
}

func TestConcatXPath(t *testing.T) {
	const firstName = "FirstName"
	const lastName = "LastName"
	const eyeColor = "EyeColor"
	const separator = " "
	const testDoc = `
	<html>
	<div>` + firstName + `</div>
	<div>` + lastName + `</div>
	<span>` + eyeColor + `</span>
	</html>
	`

	reader := strings.NewReader(testDoc)
	doc, err := htmlquery.Parse(reader)

	if err != nil {
		t.Errorf("Error loading document: %s", err.Error())
		return
	}

	xpathConfig := make(mappedConfig)
	nameAttrConfig := mappedScraperAttrConfig{
		Selector: "//div",
		Concat:   separator,
	}
	xpathConfig["Name"] = nameAttrConfig
	xpathConfig["EyeColor"] = makeSimpleAttrConfig("//span")

	scraper := mappedScraper{
		Performer: &mappedPerformerScraperConfig{
			mappedConfig: xpathConfig,
		},
	}

	q := &xpathQuery{
		doc: doc,
	}

	performer, err := scraper.scrapePerformer(q)

	if err != nil {
		t.Errorf("Error scraping performer: %s", err.Error())
		return
	}

	const performerName = firstName + separator + lastName

	verifyField(t, performerName, performer.Name, "Name")
	verifyField(t, eyeColor, performer.EyeColor, "EyeColor")
}

const sceneHTML = `
<!DOCTYPE html>

<head>
    <title>Test Video - Pornhub.com</title>
    <meta property="og:title" content="Test Video" />
    <script type="application/ld+json">
		{
			"name": "Test Video",
			"uploadDate": "2019-10-13T00:33:51+00:00",
			"author" : "Mia Malkova"
		}
	</script>
</head>

<body class="logged-out">
    <div class="container  ">
        <div id="main-container" class="clearfix">
            <div id="vpContentContainer">
                <div id="hd-leftColVideoPage">
                    <div class="video-wrapper">
						<div class="title-container">
							<i class="isMe tooltipTrig" data-title="Video of verified member"></i>
                            <h1 class="title">
                                <span class="inlineFree">Test Video</span>
                            </h1>
                        </div>

                        <div class="video-actions-container">
                            <div class="video-actions-tabs">
                                <div class="video-action-tab about-tab active">
                                    <div class="video-detailed-info">
                                        <div class="video-info-row">
                                            From:&nbsp;
                                            <div class="usernameWrap clearfix" data-type="channel">
                                                <a rel="" href="/channels/sis-loves-me" class="bolded">Sis Loves Me</a>
                                            </div>
                                        </div>

                                        <div class="video-info-row">
                                            <div class="pornstarsWrapper">
                                                Pornstars:&nbsp;
                                                <a class="pstar-list-btn js-mxp" data-mxptype="Pornstar"
                                                    data-mxptext="Alex D" href="/pornstar/alex-d">Alex D
                                                </a>
                                                , <a class="pstar-list-btn js-mxp" data-mxptype="Pornstar"
                                                    data-mxptext="Mia Malkova" href="/pornstar/mia-malkova">
                                                </a>
                                                , <a class="pstar-list-btn js-mxp" data-mxptype="Pornstar"
                                                    data-mxptext="Riley Reid" href="/pornstar/riley-reid">Riley Reid
                                                </a>
                                                <div class="tooltipTrig suggestBtn" data-title="Add a pornstar">
                                                    <a class="add-btn-small add-pornstar-btn-2">+
                                                        <span>Suggest</span></a>
                                                </div>
                                            </div>
                                        </div>

                                        <div class="video-info-row showLess">
                                            <div class="categoriesWrapper">
                                                Categories:&nbsp;
                                                <a href="/video?c=3"
                                                    onclick="ga('send', 'event', 'Watch Page', 'click', 'Category');">Amateur</a>,
                                                <a href="/categories/babe"
                                                    onclick="ga('send', 'event', 'Watch Page', 'click', 'Category');">Babe</a>,
                                                <a href="/video?c=13"
                                                    onclick="ga('send', 'event', 'Watch Page', 'click', 'Category');">Blowjob</a>,
                                                <a href="/video?c=115"
                                                    onclick="ga('send', 'event', 'Watch Page', 'click', 'Category');">Exclusive</a>,
                                                <a href="/hd"
                                                    onclick="ga('send', 'event', 'Watch Page', 'click', 'Category');">HD
                                                    Porn</a>, <a href="/categories/pornstar"
                                                    onclick="ga('send', 'event', 'Watch Page', 'click', 'Category');">Pornstar</a>,
                                                <a href="/video?c=24"
                                                    onclick="ga('send', 'event', 'Watch Page', 'click', 'Category');">Public</a>,
                                                <a href="/video?c=131"
                                                    onclick="ga('send', 'event', 'Watch Page', 'click', 'Category');">Pussy
                                                    Licking</a>, <a href="/video?c=65"
                                                    onclick="ga('send', 'event', 'Watch Page', 'click', 'Category');">Threesome</a>,
                                                <a href="/video?c=139"
                                                    onclick="ga('send', 'event', 'Watch Page', 'click', 'Category');">Verified
                                                    Models</a>
                                                <div class="tooltipTrig suggestBtn" data-title="Suggest Categories">
                                                    <a id="categoryLink" class="add-btn-small ">+
                                                        <span>Suggest</span></a>
                                                </div>
                                            </div>
                                        </div>

                                        <div class="video-info-row showLess">
                                            <div class="tagsWrapper">
                                                Tags:&nbsp;
                                                <a href="/video/search?search=3some">3some</a>, <a
                                                    href="/video?c=9">blonde</a>, <a href="/video?c=59">small tits</a>,
                                                <a href="/video/search?search=butt">butt</a>, <a
                                                    href="/video/search?search=natural+tits">natural tits</a>, <a
                                                    href="/video/search?search=petite">petite</a>, <a
                                                    href="/video?c=24">public</a>, <a
                                                    href="/video/search?search=outside">outside</a>, <a
                                                    href="/video/search?search=car">car</a>, <a
                                                    href="/video/search?search=garage">garage</a>, <a
                                                    href="/video?c=65">threesome</a>, <a
                                                    href="/video/search?search=bgg">bgg</a>, <a
                                                    href="/video/search?search=girlfrien+d">girlfrien d</a>, <a
                                                    href="/video/search?search=parking">parking</a>, <a
                                                    href="/video/search?search=sex">sex</a>, <a
                                                    href="/video/search?search=gagging">gagging</a>, <a
                                                    href="/video?c=13">blowjob</a>, <a
                                                    href="/video/search?search=bj">bj</a>, <a
                                                    href="/video/search?search=double">double</a>, <a
                                                    href="/video/search?search=ass">ass</a>
                                                <div class="tooltipTrig suggestBtn" data-title="Suggest Tags">
                                                    <a id="tagLink" class="add-btn-small">+ <span>Suggest</span></a>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</body>
</html>`

func makeSceneXPathConfig() mappedScraper {
	common := make(commonMappedConfig)

	common["$performerElem"] = `//div[@class="pornstarsWrapper"]/a[@data-mxptype="Pornstar"]`
	common["$studioElem"] = `//div[@data-type="channel"]/a`

	config := mappedSceneScraperConfig{
		mappedConfig: make(mappedConfig),
	}

	config.mappedConfig["Title"] = makeSimpleAttrConfig(`//meta[@property="og:title"]/@content`)
	// this needs post-processing
	config.mappedConfig["Date"] = makeSimpleAttrConfig(`//script[@type="application/ld+json"]`)

	tagConfig := make(mappedConfig)
	tagConfig["Name"] = makeSimpleAttrConfig(`//div[@class="categoriesWrapper"]//a[not(@class="add-btn-small ")]`)
	config.Tags = tagConfig

	performerConfig := make(mappedConfig)
	performerConfig["Name"] = makeSimpleAttrConfig(`$performerElem/@data-mxptext`)
	performerConfig["URL"] = makeSimpleAttrConfig(`$performerElem/@href`)
	config.Performers.mappedConfig = performerConfig

	studioConfig := make(mappedConfig)
	studioConfig["Name"] = makeSimpleAttrConfig(`$studioElem`)
	studioConfig["URL"] = makeSimpleAttrConfig(`$studioElem/@href`)
	config.Studio = studioConfig

	const sep = " "
	moviesNameConfig := mappedScraperAttrConfig{
		Selector: `//i[@class="isMe tooltipTrig"]/@data-title`,
		Split:    sep,
	}
	moviesConfig := make(mappedConfig)
	moviesConfig["Name"] = moviesNameConfig
	config.Movies = moviesConfig

	scraper := mappedScraper{
		Scene:  &config,
		Common: common,
	}

	return scraper
}

func verifyTags(t *testing.T, expectedTagNames []string, actualTags []*models.ScrapedSceneTag) {
	t.Helper()

	i := 0
	for i < len(expectedTagNames) || i < len(actualTags) {
		expectedTag := ""
		actualTag := ""
		if i < len(expectedTagNames) {
			expectedTag = expectedTagNames[i]
		}
		if i < len(actualTags) {
			actualTag = actualTags[i].Name
		}

		if expectedTag != actualTag {
			t.Errorf("Expected tag %s, got %s", expectedTag, actualTag)
		}
		i++
	}
}

func verifyMovies(t *testing.T, expectedMovieNames []string, actualMovies []*models.ScrapedSceneMovie) {
	t.Helper()

	i := 0
	for i < len(expectedMovieNames) || i < len(actualMovies) {
		expectedMovie := ""
		actualMovie := ""
		if i < len(expectedMovieNames) {
			expectedMovie = expectedMovieNames[i]
		}
		if i < len(actualMovies) {
			actualMovie = actualMovies[i].Name
		}

		if expectedMovie != actualMovie {
			t.Errorf("Expected movie %s, got %s", expectedMovie, actualMovie)
		}
		i++
	}
}

func verifyPerformers(t *testing.T, expectedNames []string, expectedURLs []string, actualPerformers []*models.ScrapedScenePerformer) {
	t.Helper()

	i := 0
	for i < len(expectedNames) || i < len(actualPerformers) {
		expectedName := ""
		actualName := ""
		expectedURL := ""
		actualURL := ""
		if i < len(expectedNames) {
			expectedName = expectedNames[i]
		}
		if i < len(expectedURLs) {
			expectedURL = expectedURLs[i]
		}
		if i < len(actualPerformers) {
			actualName = actualPerformers[i].Name
			if actualPerformers[i].URL != nil {
				actualURL = *actualPerformers[i].URL
			}
		}

		if expectedName != actualName {
			t.Errorf("Expected performer name %s, got %s", expectedName, actualName)
		}
		if expectedURL != actualURL {
			t.Errorf("Expected perfromer URL %s, got %s", expectedName, actualName)
		}
		i++
	}
}

func TestApplySceneXPathConfig(t *testing.T) {
	reader := strings.NewReader(sceneHTML)
	doc, err := htmlquery.Parse(reader)

	if err != nil {
		t.Errorf("Error loading document: %s", err.Error())
		return
	}

	scraper := makeSceneXPathConfig()

	q := &xpathQuery{
		doc: doc,
	}
	scene, err := scraper.scrapeScene(q)

	if err != nil {
		t.Errorf("Error scraping scene: %s", err.Error())
		return
	}

	const title = "Test Video"

	verifyField(t, title, scene.Title, "Title")

	// verify tags
	expectedTags := []string{
		"Amateur",
		"Babe",
		"Blowjob",
		"Exclusive",
		"HD Porn",
		"Pornstar",
		"Public",
		"Pussy Licking",
		"Threesome",
		"Verified Models",
	}
	verifyTags(t, expectedTags, scene.Tags)

	// verify movies
	expectedMovies := []string{
		"Video",
		"of",
		"verified",
		"member",
	}
	verifyMovies(t, expectedMovies, scene.Movies)

	expectedPerformerNames := []string{
		"Alex D",
		"Mia Malkova",
		"Riley Reid",
	}

	expectedPerformerURLs := []string{
		"/pornstar/alex-d",
		"/pornstar/mia-malkova",
		"/pornstar/riley-reid",
	}

	verifyPerformers(t, expectedPerformerNames, expectedPerformerURLs, scene.Performers)

	const expectedStudioName = "Sis Loves Me"
	const expectedStudioURL = "/channels/sis-loves-me"

	verifyField(t, expectedStudioName, &scene.Studio.Name, "Studio.Name")
	verifyField(t, expectedStudioURL, scene.Studio.URL, "Studio.URL")
}

func TestLoadXPathScraperFromYAML(t *testing.T) {
	const yamlStr = `name: Test
performerByURL:
  - action: scrapeXPath
    url: 
      - test.com
    scraper: performerScraper
xPathScrapers:
  performerScraper:
    performer:
      name: //h1[@itemprop="name"]
  sceneScraper:
    scene:
      Title:
        selector: //title
        postProcess:
          - parseDate: January 2, 2006
      Tags:
        Name: //tags  
      Movies:
        Name: //movies  
      Performers:
        Name: //performers  
      Studio:
        Name: //studio
`

	c := &config{}
	err := yaml.Unmarshal([]byte(yamlStr), &c)

	if err != nil {
		t.Errorf("Error loading yaml: %s", err.Error())
		return
	}

	// ensure fields are filled in correctly
	sceneScraper := c.XPathScrapers["sceneScraper"]
	sceneConfig := sceneScraper.Scene

	assert.Equal(t, "//title", sceneConfig.mappedConfig["Title"].Selector)
	assert.Equal(t, "//tags", sceneConfig.Tags["Name"].Selector)
	assert.Equal(t, "//movies", sceneConfig.Movies["Name"].Selector)
	assert.Equal(t, "//performers", sceneConfig.Performers.mappedConfig["Name"].Selector)
	assert.Equal(t, "//studio", sceneConfig.Studio["Name"].Selector)

	postProcess := sceneConfig.mappedConfig["Title"].postProcessActions
	parseDate := postProcess[0].(*postProcessParseDate)
	assert.Equal(t, "January 2, 2006", string(*parseDate))
}

func TestLoadInvalidXPath(t *testing.T) {
	config := make(mappedConfig)

	config["Name"] = makeSimpleAttrConfig(`//a[id=']/span`)

	reader := strings.NewReader(htmlDoc1)
	doc, err := htmlquery.Parse(reader)

	if err != nil {
		t.Errorf("Error loading document: %s", err.Error())
		return
	}

	q := &xpathQuery{
		doc: doc,
	}

	config.process(q, nil)
}

func TestSubScrape(t *testing.T) {
	retHTML := `
	<div>
		<a href="/getName">A link</a>
	</div>
	`

	ssHTML := `
	<span>The name</span>
	`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/getName" {
			fmt.Fprint(w, ssHTML)
		} else {
			fmt.Fprint(w, retHTML)
		}
	}))
	defer ts.Close()

	yamlStr := `name: Test
performerByURL:
  - action: scrapeXPath
    url: 
      - ` + ts.URL + `
    scraper: performerScraper
xPathScrapers:
  performerScraper:
    performer:
      Name: 
        selector: //div/a/@href
        postProcess:
          - replace:
              - regex: ^
                with: ` + ts.URL + `
          - subScraper:
              selector: //span
`

	c := &config{}
	err := yaml.Unmarshal([]byte(yamlStr), &c)

	if err != nil {
		t.Errorf("Error loading yaml: %s", err.Error())
		return
	}

	globalConfig := GlobalConfig{}

	performer, err := c.ScrapePerformerURL(ts.URL, nil, globalConfig)

	if err != nil {
		t.Errorf("Error scraping performer: %s", err.Error())
		return
	}

	verifyField(t, "The name", performer.Name, "Name")
}
