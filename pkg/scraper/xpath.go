package scraper

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

const (
	scraperKeyCommon    = "common"
	scraperKeyPerformer = "performer"
	scraperKeyScene     = "scene"
)

type scraperKey string

func (e scraperKey) IsValid() bool {
	switch e {
	case scraperKeyCommon, scraperKeyPerformer, scraperKeyScene:
		return true
	}
	return false
}

type xpathScraper map[scraperKey]interface{}

func (s xpathScraper) Validate() error {
	for k, v := range s {
		if !k.IsValid() {
			return fmt.Errorf("%s is not a valid scraper key", string(k))
		}

		_, ok := v.(map[interface{}]interface{})

		if !ok {
			return fmt.Errorf("%s is of incorrect type: %T", string(k), v)
		}

		// TODO - ensure all keys are prefixed with $
	}

	return nil
}

func (s xpathScraper) GetCommon() map[string]string {
	var ret map[string]string
	v, ok := s[scraperKeyCommon]
	if ok {
		ret, _ = v.(map[string]string)
	}

	return ret
}

func (s xpathScraper) GetPerformer() map[interface{}]interface{} {
	var ret map[interface{}]interface{}
	v, ok := s[scraperKeyPerformer]
	if ok {
		ret, _ = v.(map[interface{}]interface{})
	}

	return ret
}

func (s xpathScraper) GetRawScene() map[interface{}]interface{} {
	var ret map[interface{}]interface{}
	v, ok := s[scraperKeyScene]
	if ok {
		ret, _ = v.(map[interface{}]interface{})
	}

	return ret
}

func (s xpathScraper) GetScene() map[interface{}]interface{} {
	ret := make(map[interface{}]interface{})
	mapped := s.GetRawScene()

	if mapped != nil {
		for k, v := range mapped {
			if k != "Tags" && k != "Performers" {
				ret[k] = v
			}
		}
	}

	return ret
}

func (s xpathScraper) GetScenePerformers() map[interface{}]interface{} {
	var ret map[interface{}]interface{}
	mapped := s.GetRawScene()

	if mapped != nil {
		v, ok := mapped["Performers"]
		if ok {
			ret, _ = v.(map[interface{}]interface{})
		}
	}

	return ret
}

func (s xpathScraper) GetSceneTags() map[interface{}]interface{} {
	var ret map[interface{}]interface{}
	mapped := s.GetRawScene()

	if mapped != nil {
		v, ok := mapped["Tags"]
		if ok {
			ret, _ = v.(map[interface{}]interface{})
		}
	}

	return ret
}

func (s xpathScraper) GetCommonElements(doc *html.Node) map[string]interface{} {
	ret := make(map[string]interface{})

	common := s.GetCommon()
	if common == nil {
		return ret
	}

	for k, v := range common {
		elements := htmlquery.Find(doc, v)
		ret[k] = elements
	}

	return ret
}

func (s xpathScraper) scrapePerformer(doc *html.Node) (*models.ScrapedPerformer, error) {
	// parse common first
	commonMap := s.GetCommon()
	var ret models.ScrapedPerformer

	performerMap := s.GetPerformer()
	if performerMap == nil {
		return nil, nil
	}

	err := applyCommon(performerMap, commonMap)
	if err != nil {
		return nil, err
	}

	results := s.processXPathConfig(doc, performerMap)
	if len(results) > 0 {
		results[0].apply(&ret)
	}

	return &ret, nil
}

func (s xpathScraper) scrapeScene(doc *html.Node) (*models.ScrapedScene, error) {
	// parse common first
	commonMap := s.GetCommon()
	var ret models.ScrapedScene

	sceneMap := s.GetScene()
	if sceneMap == nil {
		return nil, nil
	}

	scenePerformersMap := s.GetScenePerformers()
	sceneTagsMap := s.GetSceneTags()

	err := applyCommon(sceneMap, commonMap)
	if err != nil {
		return nil, err
	}

	results := s.processXPathConfig(doc, sceneMap)
	if len(results) > 0 {
		results[0].apply(&ret)

		// now apply the performers and tags
		if scenePerformersMap != nil {
			performerResults := s.processXPathConfig(doc, scenePerformersMap)

			for _, p := range performerResults {
				performer := &models.ScrapedScenePerformer{}
				p.apply(performer)
				ret.Performers = append(ret.Performers, performer)
			}
		}

		if sceneTagsMap != nil {
			tagResults := s.processXPathConfig(doc, sceneTagsMap)

			for _, p := range tagResults {
				tag := &models.ScrapedSceneTag{}
				p.apply(tag)
				ret.Tags = append(ret.Tags, tag)
			}
		}
	}

	return &ret, nil
}

type xPathResult map[string]*html.Node

func (s xpathScraper) processXPathConfig(doc *html.Node, config map[interface{}]interface{}) []xPathResult {
	var ret []xPathResult

	for k, v := range config {
		key, keyIsStr := k.(string)
		if !keyIsStr {
			logger.Errorf("Key type not string: %T", k)
			// this should be an error or something
			continue
		}

		asStr, isStr := v.(string)
		if isStr {
			found := htmlquery.Find(doc, asStr)
			if len(found) > 0 {
				for i, elem := range found {
					if i >= len(ret) {
						ret = append(ret, make(xPathResult))
					}

					ret[i][key] = elem
				}
			}
		}
		// TODO - handle map type
	}

	return ret
}

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

func (s xpathScraper) applyXPathConfig(doc *html.Node, dest interface{}, config map[interface{}]interface{}) {
	destVal := reflect.ValueOf(dest)

	// dest should be a pointer
	destVal = destVal.Elem()

	for k, v := range config {
		key, keyIsStr := k.(string)
		if !keyIsStr {
			logger.Errorf("Key type not string: %T", k)
			// this should be an error or something
			continue
		}

		asStr, isStr := v.(string)
		if isStr {
			found := htmlquery.FindOne(doc, asStr)
			if found != nil {
				field := destVal.FieldByName(key)

				if field.IsValid() {
					value := htmlquery.InnerText(found)
					value = strings.TrimSpace(value)

					reflectValue := reflect.ValueOf(&value)

					field.Set(reflectValue)
				} else {
					logger.Errorf("Field %s does not exist in %T", key, dest)
				}
			}
		}
	}
}

type xpathScrapers map[string]*xpathScraper

func (s xpathScrapers) Validate() error {
	for _, v := range s {
		err := v.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

func applyCommon(dest map[interface{}]interface{}, common map[string]string) error {
	for commonKey, commonVal := range common {
		for destKey, destVal := range dest {
			valStr, ok := destVal.(string)
			if ok {
				if strings.Contains(valStr, commonKey) {
					dest[destKey] = strings.ReplaceAll(valStr, commonKey, commonVal)
				}
			} else {
				// hopefully this means it is a map
				// TODO - handle this
			}
		}
	}

	return nil
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
