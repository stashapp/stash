package scraper

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"gopkg.in/yaml.v2"
)

type mappedQuery interface {
	runQuery(selector string) []string
	subScrape(value string) mappedQuery
}

type commonMappedConfig map[string]string

type mappedConfig map[string]mappedScraperAttrConfig

func (s mappedConfig) applyCommon(c commonMappedConfig, src string) string {
	if c == nil {
		return src
	}

	ret := src
	for commonKey, commonVal := range c {
		if strings.Contains(ret, commonKey) {
			ret = strings.Replace(ret, commonKey, commonVal, -1)
		}
	}

	return ret
}

func (s mappedConfig) process(q mappedQuery, common commonMappedConfig) mappedResults {
	var ret mappedResults

	for k, attrConfig := range s {

		if attrConfig.Fixed != "" {
			// TODO - not sure if this needs to set _all_ indexes for the key
			const i = 0
			ret = ret.setKey(i, k, attrConfig.Fixed)
		} else {
			selector := attrConfig.Selector
			selector = s.applyCommon(common, selector)

			found := q.runQuery(selector)

			if len(found) > 0 {
				result := s.postProcess(q, attrConfig, found)
				for i, text := range result {
					ret = ret.setKey(i, k, text)
				}
			}
		}
	}

	return ret
}

func (s mappedConfig) postProcess(q mappedQuery, attrConfig mappedScraperAttrConfig, found []string) []string {
	// check if we're concatenating the results into a single result
	var ret []string
	if attrConfig.hasConcat() {
		result := attrConfig.concatenateResults(found)
		result = attrConfig.postProcess(result, q)
		if attrConfig.hasSplit() {
			return attrConfig.splitString(result)
		}

		ret = []string{result}
	} else {
		for _, text := range found {
			text = attrConfig.postProcess(text, q)
			if attrConfig.hasSplit() {
				return attrConfig.splitString(text)
			}

			ret = append(ret, text)
		}
	}

	return ret
}

type mappedSceneScraperConfig struct {
	mappedConfig

	Tags       mappedConfig                 `yaml:"Tags"`
	Performers mappedPerformerScraperConfig `yaml:"Performers"`
	Studio     mappedConfig                 `yaml:"Studio"`
	Movies     mappedConfig                 `yaml:"Movies"`
}
type _mappedSceneScraperConfig mappedSceneScraperConfig

const (
	mappedScraperConfigSceneTags       = "Tags"
	mappedScraperConfigScenePerformers = "Performers"
	mappedScraperConfigSceneStudio     = "Studio"
	mappedScraperConfigSceneMovies     = "Movies"
)

func (s *mappedSceneScraperConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// HACK - unmarshal to map first, then remove known scene sub-fields, then
	// remarshal to yaml and pass that down to the base map
	parentMap := make(map[string]interface{})
	if err := unmarshal(parentMap); err != nil {
		return err
	}

	// move the known sub-fields to a separate map
	thisMap := make(map[string]interface{})

	thisMap[mappedScraperConfigSceneTags] = parentMap[mappedScraperConfigSceneTags]
	thisMap[mappedScraperConfigScenePerformers] = parentMap[mappedScraperConfigScenePerformers]
	thisMap[mappedScraperConfigSceneStudio] = parentMap[mappedScraperConfigSceneStudio]
	thisMap[mappedScraperConfigSceneMovies] = parentMap[mappedScraperConfigSceneMovies]

	delete(parentMap, mappedScraperConfigSceneTags)
	delete(parentMap, mappedScraperConfigScenePerformers)
	delete(parentMap, mappedScraperConfigSceneStudio)
	delete(parentMap, mappedScraperConfigSceneMovies)

	// re-unmarshal the sub-fields
	yml, err := yaml.Marshal(thisMap)
	if err != nil {
		return err
	}

	// needs to be a different type to prevent infinite recursion
	c := _mappedSceneScraperConfig{}
	if err := yaml.Unmarshal(yml, &c); err != nil {
		return err
	}

	*s = mappedSceneScraperConfig(c)

	yml, err = yaml.Marshal(parentMap)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(yml, &s.mappedConfig); err != nil {
		return err
	}

	return nil
}

type mappedGalleryScraperConfig struct {
	mappedConfig

	Tags       mappedConfig `yaml:"Tags"`
	Performers mappedConfig `yaml:"Performers"`
	Studio     mappedConfig `yaml:"Studio"`
}
type _mappedGalleryScraperConfig mappedGalleryScraperConfig

func (s *mappedGalleryScraperConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// HACK - unmarshal to map first, then remove known scene sub-fields, then
	// remarshal to yaml and pass that down to the base map
	parentMap := make(map[string]interface{})
	if err := unmarshal(parentMap); err != nil {
		return err
	}

	// move the known sub-fields to a separate map
	thisMap := make(map[string]interface{})

	thisMap[mappedScraperConfigSceneTags] = parentMap[mappedScraperConfigSceneTags]
	thisMap[mappedScraperConfigScenePerformers] = parentMap[mappedScraperConfigScenePerformers]
	thisMap[mappedScraperConfigSceneStudio] = parentMap[mappedScraperConfigSceneStudio]

	delete(parentMap, mappedScraperConfigSceneTags)
	delete(parentMap, mappedScraperConfigScenePerformers)
	delete(parentMap, mappedScraperConfigSceneStudio)

	// re-unmarshal the sub-fields
	yml, err := yaml.Marshal(thisMap)
	if err != nil {
		return err
	}

	// needs to be a different type to prevent infinite recursion
	c := _mappedGalleryScraperConfig{}
	if err := yaml.Unmarshal(yml, &c); err != nil {
		return err
	}

	*s = mappedGalleryScraperConfig(c)

	yml, err = yaml.Marshal(parentMap)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(yml, &s.mappedConfig); err != nil {
		return err
	}

	return nil
}

type mappedPerformerScraperConfig struct {
	mappedConfig

	Tags mappedConfig `yaml:"Tags"`
}
type _mappedPerformerScraperConfig mappedPerformerScraperConfig

const (
	mappedScraperConfigPerformerTags = "Tags"
)

func (s *mappedPerformerScraperConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// HACK - unmarshal to map first, then remove known scene sub-fields, then
	// remarshal to yaml and pass that down to the base map
	parentMap := make(map[string]interface{})
	if err := unmarshal(parentMap); err != nil {
		return err
	}

	// move the known sub-fields to a separate map
	thisMap := make(map[string]interface{})

	thisMap[mappedScraperConfigPerformerTags] = parentMap[mappedScraperConfigPerformerTags]

	delete(parentMap, mappedScraperConfigPerformerTags)

	// re-unmarshal the sub-fields
	yml, err := yaml.Marshal(thisMap)
	if err != nil {
		return err
	}

	// needs to be a different type to prevent infinite recursion
	c := _mappedPerformerScraperConfig{}
	if err := yaml.Unmarshal(yml, &c); err != nil {
		return err
	}

	*s = mappedPerformerScraperConfig(c)

	yml, err = yaml.Marshal(parentMap)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(yml, &s.mappedConfig); err != nil {
		return err
	}

	return nil
}

type mappedMovieScraperConfig struct {
	mappedConfig

	Studio mappedConfig `yaml:"Studio"`
}
type _mappedMovieScraperConfig mappedMovieScraperConfig

const (
	mappedScraperConfigMovieStudio = "Studio"
)

func (s *mappedMovieScraperConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// HACK - unmarshal to map first, then remove known movie sub-fields, then
	// remarshal to yaml and pass that down to the base map
	parentMap := make(map[string]interface{})
	if err := unmarshal(parentMap); err != nil {
		return err
	}

	// move the known sub-fields to a separate map
	thisMap := make(map[string]interface{})

	thisMap[mappedScraperConfigMovieStudio] = parentMap[mappedScraperConfigMovieStudio]

	delete(parentMap, mappedScraperConfigMovieStudio)

	// re-unmarshal the sub-fields
	yml, err := yaml.Marshal(thisMap)
	if err != nil {
		return err
	}

	// needs to be a different type to prevent infinite recursion
	c := _mappedMovieScraperConfig{}
	if err := yaml.Unmarshal(yml, &c); err != nil {
		return err
	}

	*s = mappedMovieScraperConfig(c)

	yml, err = yaml.Marshal(parentMap)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(yml, &s.mappedConfig); err != nil {
		return err
	}

	return nil
}

type mappedRegexConfig struct {
	Regex string `yaml:"regex"`
	With  string `yaml:"with"`
}

type mappedRegexConfigs []mappedRegexConfig

func (c mappedRegexConfig) apply(value string) string {
	if c.Regex != "" {
		re, err := regexp.Compile(c.Regex)
		if err != nil {
			logger.Warnf("Error compiling regex '%s': %s", c.Regex, err.Error())
			return value
		}

		ret := re.ReplaceAllString(value, c.With)

		// trim leading and trailing whitespace
		// this is done to maintain backwards compatibility with existing
		// scrapers
		ret = strings.TrimSpace(ret)

		logger.Debugf(`Replace: '%s' with '%s'`, c.Regex, c.With)
		logger.Debugf("Before: %s", value)
		logger.Debugf("After: %s", ret)
		return ret
	}

	return value
}

func (c mappedRegexConfigs) apply(value string) string {
	// apply regex in order
	for _, config := range c {
		value = config.apply(value)
	}

	return value
}

type postProcessAction interface {
	Apply(value string, q mappedQuery) string
}

type postProcessParseDate string

func (p *postProcessParseDate) Apply(value string, q mappedQuery) string {
	parseDate := string(*p)

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

type postProcessReplace mappedRegexConfigs

func (c *postProcessReplace) Apply(value string, q mappedQuery) string {
	replace := mappedRegexConfigs(*c)
	return replace.apply(value)
}

type postProcessSubScraper mappedScraperAttrConfig

func (p *postProcessSubScraper) Apply(value string, q mappedQuery) string {
	subScrapeConfig := mappedScraperAttrConfig(*p)

	logger.Debugf("Sub-scraping for: %s", value)
	ss := q.subScrape(value)

	if ss != nil {
		found := ss.runQuery(subScrapeConfig.Selector)

		if len(found) > 0 {
			// check if we're concatenating the results into a single result
			var result string
			if subScrapeConfig.hasConcat() {
				result = subScrapeConfig.concatenateResults(found)
			} else {
				result = found[0]
			}

			result = subScrapeConfig.postProcess(result, ss)
			return result
		}
	}

	return ""
}

type postProcessMap map[string]string

func (p *postProcessMap) Apply(value string, q mappedQuery) string {
	// return the mapped value if present
	m := *p
	mapped, ok := m[value]

	if ok {
		return mapped
	}

	return value
}

type postProcessFeetToCm bool

func (p *postProcessFeetToCm) Apply(value string, q mappedQuery) string {
	const foot_in_cm = 30.48
	const inch_in_cm = 2.54

	reg := regexp.MustCompile("[0-9]+")
	filtered := reg.FindAllString(value, -1)

	var feet float64
	var inches float64
	if len(filtered) > 0 {
		feet, _ = strconv.ParseFloat(filtered[0], 64)
	}
	if len(filtered) > 1 {
		inches, _ = strconv.ParseFloat(filtered[1], 64)
	}

	var centimeters = feet*foot_in_cm + inches*inch_in_cm

	// Return rounded integer string
	return strconv.Itoa(int(math.Round(centimeters)))
}

type mappedPostProcessAction struct {
	ParseDate  string                   `yaml:"parseDate"`
	Replace    mappedRegexConfigs       `yaml:"replace"`
	SubScraper *mappedScraperAttrConfig `yaml:"subScraper"`
	Map        map[string]string        `yaml:"map"`
	FeetToCm   bool                     `yaml:"feetToCm"`
}

func (a mappedPostProcessAction) ToPostProcessAction() (postProcessAction, error) {
	var found string
	var ret postProcessAction

	if a.ParseDate != "" {
		found = "parseDate"
		action := postProcessParseDate(a.ParseDate)
		ret = &action
	}
	if len(a.Replace) > 0 {
		if found != "" {
			return nil, fmt.Errorf("post-process actions must have a single field, found %s and %s", found, "replace")
		}
		found = "replace"
		action := postProcessReplace(a.Replace)
		ret = &action
	}
	if a.SubScraper != nil {
		if found != "" {
			return nil, fmt.Errorf("post-process actions must have a single field, found %s and %s", found, "subScraper")
		}
		found = "subScraper"
		action := postProcessSubScraper(*a.SubScraper)
		ret = &action
	}
	if a.Map != nil {
		if found != "" {
			return nil, fmt.Errorf("post-process actions must have a single field, found %s and %s", found, "map")
		}
		found = "map"
		action := postProcessMap(a.Map)
		ret = &action
	}
	if a.FeetToCm {
		if found != "" {
			return nil, fmt.Errorf("post-process actions must have a single field, found %s and %s", found, "feetToCm")
		}
		found = "feetToCm"
		action := postProcessFeetToCm(a.FeetToCm)
		ret = &action
	}

	if ret == nil {
		return nil, errors.New("invalid post-process action")
	}

	return ret, nil
}

type mappedScraperAttrConfig struct {
	Selector    string                    `yaml:"selector"`
	Fixed       string                    `yaml:"fixed"`
	PostProcess []mappedPostProcessAction `yaml:"postProcess"`
	Concat      string                    `yaml:"concat"`
	Split       string                    `yaml:"split"`

	postProcessActions []postProcessAction

	// deprecated: use PostProcess instead
	ParseDate  string                   `yaml:"parseDate"`
	Replace    mappedRegexConfigs       `yaml:"replace"`
	SubScraper *mappedScraperAttrConfig `yaml:"subScraper"`
}

type _mappedScraperAttrConfig mappedScraperAttrConfig

func (c *mappedScraperAttrConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// try unmarshalling into a string first
	if err := unmarshal(&c.Selector); err != nil {
		// if it's a type error then we try to unmarshall to the full object
		if _, ok := err.(*yaml.TypeError); !ok {
			return err
		}

		// unmarshall to full object
		// need it as a separate object
		t := _mappedScraperAttrConfig{}
		if err = unmarshal(&t); err != nil {
			return err
		}

		*c = mappedScraperAttrConfig(t)
	}

	return c.convertPostProcessActions()
}

func (c *mappedScraperAttrConfig) convertPostProcessActions() error {
	// ensure we don't have the old deprecated fields and the new post process field
	if len(c.PostProcess) > 0 {
		if c.ParseDate != "" || len(c.Replace) > 0 || c.SubScraper != nil {
			return errors.New("cannot include postProcess and (parseDate, replace, subScraper) deprecated fields")
		}

		// convert xpathPostProcessAction actions to postProcessActions
		for _, a := range c.PostProcess {
			action, err := a.ToPostProcessAction()
			if err != nil {
				return err
			}
			c.postProcessActions = append(c.postProcessActions, action)
		}

		c.PostProcess = nil
	} else {
		// convert old deprecated fields if present
		// in same order as they used to be executed
		if len(c.Replace) > 0 {
			action := postProcessReplace(c.Replace)
			c.postProcessActions = append(c.postProcessActions, &action)
			c.Replace = nil
		}

		if c.SubScraper != nil {
			action := postProcessSubScraper(*c.SubScraper)
			c.postProcessActions = append(c.postProcessActions, &action)
			c.SubScraper = nil
		}

		if c.ParseDate != "" {
			action := postProcessParseDate(c.ParseDate)
			c.postProcessActions = append(c.postProcessActions, &action)
			c.ParseDate = ""
		}
	}

	return nil
}

func (c mappedScraperAttrConfig) hasConcat() bool {
	return c.Concat != ""
}

func (c mappedScraperAttrConfig) hasSplit() bool {
	return c.Split != ""
}

func (c mappedScraperAttrConfig) concatenateResults(nodes []string) string {
	separator := c.Concat
	result := []string{}

	for _, text := range nodes {
		result = append(result, text)
	}

	return strings.Join(result, separator)
}

func (c mappedScraperAttrConfig) splitString(value string) []string {
	separator := c.Split
	var res []string

	if separator == "" {
		return []string{value}
	}

	for _, str := range strings.Split(value, separator) {
		if str != "" {
			res = append(res, str)
		}
	}

	return res
}

func (c mappedScraperAttrConfig) postProcess(value string, q mappedQuery) string {
	for _, action := range c.postProcessActions {
		value = action.Apply(value, q)
	}

	return value
}

type mappedScrapers map[string]*mappedScraper

type mappedScraper struct {
	Common    commonMappedConfig            `yaml:"common"`
	Scene     *mappedSceneScraperConfig     `yaml:"scene"`
	Gallery   *mappedGalleryScraperConfig   `yaml:"gallery"`
	Performer *mappedPerformerScraperConfig `yaml:"performer"`
	Movie     *mappedMovieScraperConfig     `yaml:"movie"`
}

type mappedResult map[string]string
type mappedResults []mappedResult

func (r mappedResult) apply(dest interface{}) {
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

func (r mappedResults) setKey(index int, key string, value string) mappedResults {
	if index >= len(r) {
		r = append(r, make(mappedResult))
	}

	logger.Debugf(`[%d][%s] = %s`, index, key, value)
	r[index][key] = value
	return r
}

func (s mappedScraper) scrapePerformer(q mappedQuery) (*models.ScrapedPerformer, error) {
	var ret models.ScrapedPerformer

	performerMap := s.Performer
	if performerMap == nil {
		return nil, nil
	}

	performerTagsMap := performerMap.Tags

	results := performerMap.process(q, s.Common)
	if len(results) > 0 {
		results[0].apply(&ret)

		// now apply the tags
		if performerTagsMap != nil {
			logger.Debug(`Processing performer tags:`)
			tagResults := performerTagsMap.process(q, s.Common)

			for _, p := range tagResults {
				tag := &models.ScrapedSceneTag{}
				p.apply(tag)
				ret.Tags = append(ret.Tags, tag)
			}
		}
	}

	return &ret, nil
}

func (s mappedScraper) scrapePerformers(q mappedQuery) ([]*models.ScrapedPerformer, error) {
	var ret []*models.ScrapedPerformer

	performerMap := s.Performer
	if performerMap == nil {
		return nil, nil
	}

	results := performerMap.process(q, s.Common)
	for _, r := range results {
		var p models.ScrapedPerformer
		r.apply(&p)
		ret = append(ret, &p)
	}

	return ret, nil
}

func (s mappedScraper) scrapeScene(q mappedQuery) (*models.ScrapedScene, error) {
	var ret models.ScrapedScene

	sceneScraperConfig := s.Scene
	sceneMap := sceneScraperConfig.mappedConfig
	if sceneMap == nil {
		return nil, nil
	}

	scenePerformersMap := sceneScraperConfig.Performers
	sceneTagsMap := sceneScraperConfig.Tags
	sceneStudioMap := sceneScraperConfig.Studio
	sceneMoviesMap := sceneScraperConfig.Movies

	scenePerformerTagsMap := scenePerformersMap.Tags

	logger.Debug(`Processing scene:`)
	results := sceneMap.process(q, s.Common)
	if len(results) > 0 {
		results[0].apply(&ret)

		// process performer tags once
		var performerTagResults mappedResults
		if scenePerformerTagsMap != nil {
			performerTagResults = scenePerformerTagsMap.process(q, s.Common)
		}

		// now apply the performers and tags
		if scenePerformersMap.mappedConfig != nil {
			logger.Debug(`Processing scene performers:`)
			performerResults := scenePerformersMap.process(q, s.Common)

			for _, p := range performerResults {
				performer := &models.ScrapedScenePerformer{}
				p.apply(performer)

				for _, p := range performerTagResults {
					tag := &models.ScrapedSceneTag{}
					p.apply(tag)
					ret.Tags = append(ret.Tags, tag)
				}

				ret.Performers = append(ret.Performers, performer)
			}
		}

		if sceneTagsMap != nil {
			logger.Debug(`Processing scene tags:`)
			tagResults := sceneTagsMap.process(q, s.Common)

			for _, p := range tagResults {
				tag := &models.ScrapedSceneTag{}
				p.apply(tag)
				ret.Tags = append(ret.Tags, tag)
			}
		}

		if sceneStudioMap != nil {
			logger.Debug(`Processing scene studio:`)
			studioResults := sceneStudioMap.process(q, s.Common)

			if len(studioResults) > 0 {
				studio := &models.ScrapedSceneStudio{}
				studioResults[0].apply(studio)
				ret.Studio = studio
			}
		}

		if sceneMoviesMap != nil {
			logger.Debug(`Processing scene movies:`)
			movieResults := sceneMoviesMap.process(q, s.Common)

			for _, p := range movieResults {
				movie := &models.ScrapedSceneMovie{}
				p.apply(movie)
				ret.Movies = append(ret.Movies, movie)
			}

		}
	}

	return &ret, nil
}

func (s mappedScraper) scrapeGallery(q mappedQuery) (*models.ScrapedGallery, error) {
	var ret models.ScrapedGallery

	galleryScraperConfig := s.Gallery
	galleryMap := galleryScraperConfig.mappedConfig
	if galleryMap == nil {
		return nil, nil
	}

	galleryPerformersMap := galleryScraperConfig.Performers
	galleryTagsMap := galleryScraperConfig.Tags
	galleryStudioMap := galleryScraperConfig.Studio

	logger.Debug(`Processing gallery:`)
	results := galleryMap.process(q, s.Common)
	if len(results) > 0 {
		results[0].apply(&ret)

		// now apply the performers and tags
		if galleryPerformersMap != nil {
			logger.Debug(`Processing gallery performers:`)
			performerResults := galleryPerformersMap.process(q, s.Common)

			for _, p := range performerResults {
				performer := &models.ScrapedScenePerformer{}
				p.apply(performer)
				ret.Performers = append(ret.Performers, performer)
			}
		}

		if galleryTagsMap != nil {
			logger.Debug(`Processing gallery tags:`)
			tagResults := galleryTagsMap.process(q, s.Common)

			for _, p := range tagResults {
				tag := &models.ScrapedSceneTag{}
				p.apply(tag)
				ret.Tags = append(ret.Tags, tag)
			}
		}

		if galleryStudioMap != nil {
			logger.Debug(`Processing gallery studio:`)
			studioResults := galleryStudioMap.process(q, s.Common)

			if len(studioResults) > 0 {
				studio := &models.ScrapedSceneStudio{}
				studioResults[0].apply(studio)
				ret.Studio = studio
			}
		}
	}

	return &ret, nil
}

func (s mappedScraper) scrapeMovie(q mappedQuery) (*models.ScrapedMovie, error) {
	var ret models.ScrapedMovie

	movieScraperConfig := s.Movie
	movieMap := movieScraperConfig.mappedConfig
	if movieMap == nil {
		return nil, nil
	}

	movieStudioMap := movieScraperConfig.Studio

	results := movieMap.process(q, s.Common)
	if len(results) > 0 {
		results[0].apply(&ret)

		if movieStudioMap != nil {
			logger.Debug(`Processing movie studio:`)
			studioResults := movieStudioMap.process(q, s.Common)

			if len(studioResults) > 0 {
				studio := &models.ScrapedMovieStudio{}
				studioResults[0].apply(studio)
				ret.Studio = studio
			}
		}
	}

	return &ret, nil
}
