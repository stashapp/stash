package scraper

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
)

// Timeout for the scrape http request. Includes transfer time. May want to make this
// configurable at some point.
const scrapeGetTimeout = time.Second * 30

type commonXPathConfig map[string]string

func (c commonXPathConfig) applyCommon(src string) string {
	ret := src
	for commonKey, commonVal := range c {
		if strings.Contains(ret, commonKey) {
			ret = strings.Replace(ret, commonKey, commonVal, -1)
		}
	}

	return ret
}

type xpathScraperConfig map[string]interface{}

func createXPathScraperConfig(src map[interface{}]interface{}) xpathScraperConfig {
	ret := make(xpathScraperConfig)

	if src != nil {
		for k, v := range src {
			keyStr, isStr := k.(string)
			if isStr {
				ret[keyStr] = v
			}
		}
	}

	return ret
}

type xpathRegexConfig map[interface{}]interface{}
type xpathRegexConfigs []xpathRegexConfig

func (c xpathRegexConfig) apply(value string) string {
	regex := ""
	with := ""

	if regexI, _ := c["regex"]; regexI != nil {
		regex, _ = regexI.(string)
	}
	if withI, _ := c["with"]; withI != nil {
		with, _ = withI.(string)
	}

	if regex != "" {
		re, err := regexp.Compile(regex)
		if err != nil {
			logger.Warnf("Error compiling regex '%s': %s", regex, err.Error())
			return value
		}

		ret := re.ReplaceAllString(value, with)

		logger.Debugf(`Replace: '%s' with '%s'`, regex, with)
		logger.Debugf("Before: %s", value)
		logger.Debugf("After: %s", ret)
		return ret
	}

	return value
}

func (c xpathRegexConfigs) apply(value string) string {
	// apply regex in order
	for _, config := range c {
		value = config.apply(value)
	}

	// remove whitespace again
	value = commonPostProcess(value)

	return value
}

type xpathScraperAttrConfig map[interface{}]interface{}

func (c xpathScraperAttrConfig) getString(key string) string {
	ret, _ := c[key]

	if ret == nil {
		return ""
	}

	asStr, _ := ret.(string)
	return asStr
}

func (c xpathScraperAttrConfig) getSelector() string {
	const selectorKey = "selector"
	return c.getString(selectorKey)
}

func (c xpathScraperAttrConfig) getConcat() string {
	const concatKey = "concat"
	return c.getString(concatKey)
}

func (c xpathScraperAttrConfig) hasConcat() bool {
	return c.getConcat() != ""
}

func (c xpathScraperAttrConfig) getParseDate() string {
	const parseDateKey = "parseDate"
	return c.getString(parseDateKey)
}

func (c xpathScraperAttrConfig) getReplace() xpathRegexConfigs {
	const replaceKey = "replace"
	val, _ := c[replaceKey]

	var ret xpathRegexConfigs
	if val == nil {
		return ret
	}

	asSlice, _ := val.([]interface{})

	for _, v := range asSlice {
		asMap, _ := v.(map[interface{}]interface{})
		ret = append(ret, xpathRegexConfig(asMap))
	}

	return ret
}

func (c xpathScraperAttrConfig) getSubScraper() xpathScraperAttrConfig {
	const subScraperKey = "subScraper"
	val, _ := c[subScraperKey]

	if val == nil {
		return nil
	}

	asMap, _ := val.(map[interface{}]interface{})
	if asMap != nil {
		return xpathScraperAttrConfig(asMap)
	}

	return nil
}

func (c xpathScraperAttrConfig) concatenateResults(nodes []*html.Node) string {
	separator := c.getConcat()
	result := []string{}

	for _, elem := range nodes {
		text := NodeText(elem)
		text = commonPostProcess(text)

		result = append(result, text)
	}

	return strings.Join(result, separator)
}

func (c xpathScraperAttrConfig) parseDate(value string) string {
	parseDate := c.getParseDate()

	if parseDate == "" {
		return value
	}

	// try to parse the date using the pattern
	// if it fails, then just fall back to the original value
	parsedValue, err := time.Parse(parseDate, value)
	if err != nil {
		logger.Warnf("Error parsing date string '%s' using format '%s': %s", value, parseDate, err.Error())
		return value
	}

	// convert it into our date format
	const internalDateFormat = "2006-01-02"
	return parsedValue.Format(internalDateFormat)
}

func (c xpathScraperAttrConfig) replaceRegex(value string) string {
	replace := c.getReplace()
	return replace.apply(value)
}

func (c xpathScraperAttrConfig) applySubScraper(value string) string {
	subScraper := c.getSubScraper()

	if subScraper == nil {
		return value
	}

	logger.Debugf("Sub-scraping for: %s", value)
	doc, err := loadURL(value, nil)

	if err != nil {
		logger.Warnf("Error getting URL '%s' for sub-scraper: %s", value, err.Error())
		return ""
	}

	found := runXPathQuery(doc, subScraper.getSelector(), nil)

	if len(found) > 0 {
		// check if we're concatenating the results into a single result
		var result string
		if subScraper.hasConcat() {
			result = subScraper.concatenateResults(found)
		} else {
			result = NodeText(found[0])
			result = commonPostProcess(result)
		}

		result = subScraper.postProcess(result)
		return result
	}

	return ""
}

func (c xpathScraperAttrConfig) postProcess(value string) string {
	// perform regex replacements first
	value = c.replaceRegex(value)
	value = c.applySubScraper(value)
	value = c.parseDate(value)

	return value
}

func commonPostProcess(value string) string {
	value = strings.TrimSpace(value)

	// remove multiple whitespace and end lines
	re := regexp.MustCompile("\n")
	value = re.ReplaceAllString(value, "")
	re = regexp.MustCompile("  +")
	value = re.ReplaceAllString(value, " ")

	return value
}

func runXPathQuery(doc *html.Node, xpath string, common commonXPathConfig) []*html.Node {
	// apply common
	if common != nil {
		xpath = common.applyCommon(xpath)
	}

	found, err := htmlquery.QueryAll(doc, xpath)
	if err != nil {
		logger.Warnf("Error parsing xpath expression '%s': %s", xpath, err.Error())
		return nil
	}

	return found
}

func (s xpathScraperConfig) process(doc *html.Node, common commonXPathConfig) xPathResults {
	var ret xPathResults

	for k, value := range s {
		switch v := value.(type) {
		case string:
			found := runXPathQuery(doc, v, common)

			if len(found) > 0 {
				for i, elem := range found {
					text := NodeText(elem)
					text = commonPostProcess(text)

					ret = ret.setKey(i, k, text)
				}
			}
		case map[interface{}]interface{}:
			attrConfig := xpathScraperAttrConfig(v)

			found := runXPathQuery(doc, attrConfig.getSelector(), common)

			if len(found) > 0 {
				// check if we're concatenating the results into a single result
				if attrConfig.hasConcat() {
					result := attrConfig.concatenateResults(found)
					result = attrConfig.postProcess(result)
					const i = 0
					ret = ret.setKey(i, k, result)
				} else {
					for i, elem := range found {
						text := NodeText(elem)
						text = commonPostProcess(text)
						text = attrConfig.postProcess(text)

						ret = ret.setKey(i, k, text)
					}
				}
			}
		}
	}

	return ret
}

type xpathScrapers map[string]*xpathScraper

type xpathScraper struct {
	Common    commonXPathConfig  `yaml:"common"`
	Scene     xpathScraperConfig `yaml:"scene"`
	Performer xpathScraperConfig `yaml:"performer"`
}

const (
	XPathScraperConfigSceneTags       = "Tags"
	XPathScraperConfigScenePerformers = "Performers"
	XPathScraperConfigSceneStudio     = "Studio"
	XPathScraperConfigSceneMovies     = "Movies"
)

func (s xpathScraper) GetSceneSimple() xpathScraperConfig {
	// exclude the complex sub-configs
	ret := make(xpathScraperConfig)
	mapped := s.Scene

	if mapped != nil {
		for k, v := range mapped {
			if k != XPathScraperConfigSceneTags && k != XPathScraperConfigScenePerformers && k != XPathScraperConfigSceneStudio && k != XPathScraperConfigSceneMovies {
				ret[k] = v
			}
		}
	}

	return ret
}

func (s xpathScraper) getSceneSubMap(key string) xpathScraperConfig {
	var ret map[interface{}]interface{}
	mapped := s.Scene

	if mapped != nil {
		v, ok := mapped[key]
		if ok {
			ret, _ = v.(map[interface{}]interface{})
		}
	}

	if ret != nil {
		return createXPathScraperConfig(ret)
	}

	return nil
}

func (s xpathScraper) GetScenePerformers() xpathScraperConfig {
	return s.getSceneSubMap(XPathScraperConfigScenePerformers)
}

func (s xpathScraper) GetSceneTags() xpathScraperConfig {
	return s.getSceneSubMap(XPathScraperConfigSceneTags)
}

func (s xpathScraper) GetSceneStudio() xpathScraperConfig {
	return s.getSceneSubMap(XPathScraperConfigSceneStudio)
}

func (s xpathScraper) GetSceneMovies() xpathScraperConfig {
	return s.getSceneSubMap(XPathScraperConfigSceneMovies)
}

func (s xpathScraper) scrapePerformer(doc *html.Node) (*models.ScrapedPerformer, error) {
	var ret models.ScrapedPerformer

	performerMap := s.Performer
	if performerMap == nil {
		return nil, nil
	}

	results := performerMap.process(doc, s.Common)
	if len(results) > 0 {
		results[0].apply(&ret)
	}

	return &ret, nil
}

func (s xpathScraper) scrapePerformers(doc *html.Node) ([]*models.ScrapedPerformer, error) {
	var ret []*models.ScrapedPerformer

	performerMap := s.Performer
	if performerMap == nil {
		return nil, nil
	}

	results := performerMap.process(doc, s.Common)
	for _, r := range results {
		var p models.ScrapedPerformer
		r.apply(&p)
		ret = append(ret, &p)
	}

	return ret, nil
}

func (s xpathScraper) scrapeScene(doc *html.Node) (*models.ScrapedScene, error) {
	var ret models.ScrapedScene

	sceneMap := s.GetSceneSimple()
	if sceneMap == nil {
		return nil, nil
	}

	scenePerformersMap := s.GetScenePerformers()
	sceneTagsMap := s.GetSceneTags()
	sceneStudioMap := s.GetSceneStudio()
	sceneMoviesMap := s.GetSceneMovies()

	logger.Debug(`Processing scene:`)
	results := sceneMap.process(doc, s.Common)
	if len(results) > 0 {
		results[0].apply(&ret)

		// now apply the performers and tags
		if scenePerformersMap != nil {
			logger.Debug(`Processing scene performers:`)
			performerResults := scenePerformersMap.process(doc, s.Common)

			for _, p := range performerResults {
				performer := &models.ScrapedScenePerformer{}
				p.apply(performer)
				ret.Performers = append(ret.Performers, performer)
			}
		}

		if sceneTagsMap != nil {
			logger.Debug(`Processing scene tags:`)
			tagResults := sceneTagsMap.process(doc, s.Common)

			for _, p := range tagResults {
				tag := &models.ScrapedSceneTag{}
				p.apply(tag)
				ret.Tags = append(ret.Tags, tag)
			}
		}

		if sceneStudioMap != nil {
			logger.Debug(`Processing scene studio:`)
			studioResults := sceneStudioMap.process(doc, s.Common)

			if len(studioResults) > 0 {
				studio := &models.ScrapedSceneStudio{}
				studioResults[0].apply(studio)
				ret.Studio = studio
			}
		}

		if sceneMoviesMap != nil {
			logger.Debug(`Processing scene movies:`)
			movieResults := sceneMoviesMap.process(doc, s.Common)

			for _, p := range movieResults {
				movie := &models.ScrapedSceneMovie{}
				p.apply(movie)
				ret.Movies = append(ret.Movies, movie)
			}

		}
	}

	return &ret, nil
}

type xPathResult map[string]string
type xPathResults []xPathResult

func (r xPathResult) apply(dest interface{}) {
	destVal := reflect.ValueOf(dest)

	// dest should be a pointer
	destVal = destVal.Elem()

	for key, value := range r {
		field := destVal.FieldByName(key)

		if field.IsValid() {
			var reflectValue reflect.Value
			if field.Kind() == reflect.Ptr {
				// need to copy the value, otherwise everything is set to the
				// same pointer
				localValue := value
				reflectValue = reflect.ValueOf(&localValue)
			} else {
				reflectValue = reflect.ValueOf(value)
			}

			field.Set(reflectValue)
		} else {
			logger.Errorf("Field %s does not exist in %T", key, dest)
		}
	}
}

func (r xPathResults) setKey(index int, key string, value string) xPathResults {
	if index >= len(r) {
		r = append(r, make(xPathResult))
	}

	logger.Debugf(`[%d][%s] = %s`, index, key, value)
	r[index][key] = value
	return r
}

func loadURL(url string, c *scraperConfig) (*html.Node, error) {
	client := &http.Client{
		Timeout: scrapeGetTimeout,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	userAgent := config.GetScraperUserAgent()
	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}

	ret, err := html.Parse(r)

	if err == nil && c != nil && c.DebugOptions != nil && c.DebugOptions.PrintHTML {
		var b bytes.Buffer
		html.Render(&b, ret)
		logger.Infof("loadURL (%s) response: \n%s", url, b.String())
	}

	return ret, err
}

func scrapePerformerURLXpath(c scraperTypeConfig, url string) (*models.ScrapedPerformer, error) {
	scraper := c.scraperConfig.XPathScrapers[c.Scraper]

	if scraper == nil {
		return nil, errors.New("xpath scraper with name " + c.Scraper + " not found in config")
	}

	doc, err := loadURL(url, c.scraperConfig)

	if err != nil {
		return nil, err
	}

	return scraper.scrapePerformer(doc)
}

func scrapeSceneURLXPath(c scraperTypeConfig, url string) (*models.ScrapedScene, error) {
	scraper := c.scraperConfig.XPathScrapers[c.Scraper]

	if scraper == nil {
		return nil, errors.New("xpath scraper with name " + c.Scraper + " not found in config")
	}

	doc, err := loadURL(url, c.scraperConfig)

	if err != nil {
		return nil, err
	}

	return scraper.scrapeScene(doc)
}

func scrapePerformerNamesXPath(c scraperTypeConfig, name string) ([]*models.ScrapedPerformer, error) {
	scraper := c.scraperConfig.XPathScrapers[c.Scraper]

	if scraper == nil {
		return nil, errors.New("xpath scraper with name " + c.Scraper + " not found in config")
	}

	const placeholder = "{}"

	// replace the placeholder string with the URL-escaped name
	escapedName := url.QueryEscape(name)

	u := c.QueryURL
	u = strings.Replace(u, placeholder, escapedName, -1)

	doc, err := loadURL(u, c.scraperConfig)

	if err != nil {
		return nil, err
	}

	return scraper.scrapePerformers(doc)
}

func NodeText(n *html.Node) string {
	if n != nil && n.Type == html.CommentNode {
		return htmlquery.OutputHTML(n, true)
	}
	return htmlquery.InnerText(n)
}
