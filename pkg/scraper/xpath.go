package scraper

import (
	"errors"
	"reflect"
	"regexp"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type commonXPathConfig map[string]string

func (c commonXPathConfig) applyCommon(src string) string {
	ret := src
	for commonKey, commonVal := range c {
		if strings.Contains(ret, commonKey) {
			ret = strings.ReplaceAll(ret, commonKey, commonVal)
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

func (s xpathScraperConfig) process(doc *html.Node, common commonXPathConfig) []xPathResult {
	var ret []xPathResult

	for k, v := range s {
		asStr, isStr := v.(string)

		if isStr {
			// apply common
			if common != nil {
				asStr = common.applyCommon(asStr)
			}

			found := htmlquery.Find(doc, asStr)
			if len(found) > 0 {
				for i, elem := range found {
					if i >= len(ret) {
						ret = append(ret, make(xPathResult))
					}

					ret[i][k] = elem
				}
			}
		}
		// TODO - handle map type
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
)

func (s xpathScraper) GetSceneSimple() xpathScraperConfig {
	// exclude the complex sub-configs
	ret := make(xpathScraperConfig)
	mapped := s.Scene

	if mapped != nil {
		for k, v := range mapped {
			if k != XPathScraperConfigSceneTags && k != XPathScraperConfigScenePerformers && k != XPathScraperConfigSceneStudio {
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

func (s xpathScraper) scrapeScene(doc *html.Node) (*models.ScrapedScene, error) {
	var ret models.ScrapedScene

	sceneMap := s.GetSceneSimple()
	if sceneMap == nil {
		return nil, nil
	}

	scenePerformersMap := s.GetScenePerformers()
	sceneTagsMap := s.GetSceneTags()
	sceneStudioMap := s.GetSceneStudio()

	results := sceneMap.process(doc, s.Common)
	if len(results) > 0 {
		results[0].apply(&ret)

		// now apply the performers and tags
		if scenePerformersMap != nil {
			performerResults := scenePerformersMap.process(doc, s.Common)

			for _, p := range performerResults {
				performer := &models.ScrapedScenePerformer{}
				p.apply(performer)
				ret.Performers = append(ret.Performers, performer)
			}
		}

		if sceneTagsMap != nil {
			tagResults := sceneTagsMap.process(doc, s.Common)

			for _, p := range tagResults {
				tag := &models.ScrapedSceneTag{}
				p.apply(tag)
				ret.Tags = append(ret.Tags, tag)
			}
		}

		if sceneStudioMap != nil {
			studioResults := sceneStudioMap.process(doc, s.Common)

			if len(studioResults) > 0 {
				studio := &models.ScrapedSceneStudio{}
				studioResults[0].apply(studio)
				ret.Studio = studio
			}
		}
	}

	return &ret, nil
}

type xPathResult map[string]*html.Node

func (r xPathResult) apply(dest interface{}) {
	destVal := reflect.ValueOf(dest)

	// dest should be a pointer
	destVal = destVal.Elem()

	for key, v := range r {
		field := destVal.FieldByName(key)

		if field.IsValid() {
			value := htmlquery.InnerText(v)
			value = strings.TrimSpace(value)

			// remove multiple whitespace and end lines
			re := regexp.MustCompile("\n")
			value = re.ReplaceAllString(value, "")
			re = regexp.MustCompile("  +")
			value = re.ReplaceAllString(value, " ")

			var reflectValue reflect.Value
			if field.Kind() == reflect.Ptr {
				reflectValue = reflect.ValueOf(&value)
			} else {
				reflectValue = reflect.ValueOf(value)
			}

			field.Set(reflectValue)
		} else {
			logger.Errorf("Field %s does not exist in %T", key, dest)
		}
	}
}

func scrapePerformerURLXpath(c scraperTypeConfig, url string) (*models.ScrapedPerformer, error) {
	scraper := c.scraperConfig.XPathScrapers[c.Scraper]

	if scraper == nil {
		return nil, errors.New("xpath scraper with name " + c.Scraper + " not found in config")
	}

	doc, err := htmlquery.LoadURL(url)

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

	doc, err := htmlquery.LoadURL(url)

	if err != nil {
		return nil, err
	}

	return scraper.scrapeScene(doc)
}
