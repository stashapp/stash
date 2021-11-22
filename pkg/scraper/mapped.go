package scraper

import (
	"context"
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
	"github.com/stashapp/stash/pkg/utils"
	"gopkg.in/yaml.v2"
)

type queryer interface {
	query(selector string) ([]string, error)
	subScrape(ctx context.Context, value string) queryer
}

// commonMappedConfig defines a set of common replacements
type commonMappedConfig map[string]string

// replace uses the common map to perform substitutions in s.
func (c commonMappedConfig) replace(s string) string {
	if c == nil {
		return s
	}

	for commonKey, commonVal := range c {
		s = strings.ReplaceAll(s, commonKey, commonVal)
	}

	return s
}

type mappedConfig map[string]mappedScraperAttrConfig

func (s mappedConfig) process(ctx context.Context, q queryer, common commonMappedConfig) mappedResults {
	var ret mappedResults

	for k, attrConfig := range s {
		if attrConfig.Fixed != "" {
			// TODO - not sure if this needs to set _all_ indexes for the key
			const i = 0
			ret = ret.setKey(i, k, attrConfig.Fixed)
		} else {
			selector := attrConfig.Selector
			selector = common.replace(selector)

			found, err := q.query(selector)
			if err != nil {
				logger.Warnf("key '%v': %v", k, err)
			}

			if len(found) > 0 {
				result := s.postProcess(ctx, q, attrConfig, found)
				for i, text := range result {
					ret = ret.setKey(i, k, text)
				}
			}
		}
	}

	return ret
}

func (s mappedConfig) postProcess(ctx context.Context, q queryer, attrConfig mappedScraperAttrConfig, found []string) []string {
	// check if we're concatenating the results into a single result
	var ret []string
	if attrConfig.Concat != "" {
		result := strings.Join(found, attrConfig.Concat)
		result = attrConfig.postProcess(ctx, result, q)
		if attrConfig.Split != "" {
			results := fieldSplit(result, attrConfig.Split)
			results = cleanResults(results)
			return results
		}

		ret = []string{result}
	} else {
		for _, text := range found {
			text = attrConfig.postProcess(ctx, text, q)
			if attrConfig.Split != "" {
				return fieldSplit(text, attrConfig.Split)
			}

			ret = append(ret, text)
		}
		ret = cleanResults(ret)
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

type mappedRegexConfigs []mappedRegexConfig

func (cs mappedRegexConfigs) replace(s string) string {
	// apply regex in order, by chaining functions
	for _, c := range cs {
		f := regexFunc(c)
		s = f(s)
	}

	return s
}

type mappedRegexConfig struct {
	Regex string `yaml:"regex"`
	With  string `yaml:"with"`
}

// regexFunc turns a regex configuration block into a function
func regexFunc(c mappedRegexConfig) func(string) string {
	return func(s string) string {
		if c.Regex == "" {
			return s
		}

		re, err := regexp.Compile(c.Regex)
		if err != nil {
			logger.Warnf("Error compiling regex '%s': %s", c.Regex, err.Error())
			return s
		}

		ret := re.ReplaceAllString(s, c.With)

		// trim leading and trailing whitespace
		// this is done to maintain backwards compatibility with existing
		// scrapers
		ret = strings.TrimSpace(ret)

		logger.Debugf(`Replace: '%s' with '%s'`, c.Regex, c.With)
		logger.Debugf("Before: %s", s)
		logger.Debugf("After: %s", ret)

		return ret
	}
}

// handlerFunc is the type of post-process handlers
type handlerFunc func(ctx context.Context, value string, q queryer) string

// Apply runs the given handlerfunc
func (f handlerFunc) Apply(ctx context.Context, value string, q queryer) string {
	return f(ctx, value, q)
}

// done is a handlerFunc which returns its value. This is a no-op func
// which can be "plugged into" a handler chain to terminate it.
func done() handlerFunc {
	return func(ctx context.Context, value string, q queryer) string {
		return value
	}
}

func postProcessParseDate(layout string, next handlerFunc) handlerFunc {
	const internalDateFormat = "2006-01-02"

	return func(ctx context.Context, value string, q queryer) string {
		valueLower := strings.ToLower(value)
		if valueLower == "today" || valueLower == "yesterday" { // handle today, yesterday
			dt := time.Now()
			if valueLower == "yesterday" { // subtract 1 day from now
				dt = dt.AddDate(0, 0, -1)
			}
			return next(ctx, dt.Format(internalDateFormat), q)
		}

		if layout == "" {
			// Nothing to do, next please!
			return next(ctx, value, q)
		}

		// try to parse the date using the pattern
		// if it fails, then just fall back to the original value
		parsedValue, err := time.Parse(layout, value)
		if err != nil {
			logger.Warnf("Error parsing date string '%s' using format '%s': %s", value, layout, err.Error())
			return next(ctx, value, q)
		}

		// convert it into our date format
		return next(ctx, parsedValue.Format(internalDateFormat), q)
	}
}

func postProcessSubtractDays(next handlerFunc) handlerFunc {
	return func(ctx context.Context, value string, q queryer) string {
		const internalDateFormat = "2006-01-02"

		i, err := strconv.Atoi(value)
		if err != nil {
			logger.Warnf("Error parsing day string %s: %s", value, err)
			return next(ctx, value, q)
		}

		dt := time.Now()
		dt = dt.AddDate(0, 0, -i)
		return next(ctx, dt.Format(internalDateFormat), q)
	}
}

func postProcessReplace(cs mappedRegexConfigs, next handlerFunc) handlerFunc {
	return func(ctx context.Context, value string, q queryer) string {
		value = cs.replace(value)
		return next(ctx, value, q)
	}
}

func postProcessSubScraper(sub mappedScraperAttrConfig, next handlerFunc) handlerFunc {
	return func(ctx context.Context, value string, q queryer) string {
		logger.Debugf("Sub-scraping for: %s", value)
		ss := q.subScrape(ctx, value)

		if ss != nil {
			found, err := ss.query(sub.Selector)
			if err != nil {
				logger.Warnf("subscrape for '%v': %v", value, err)
			}

			if len(found) > 0 {
				// check if we're concatenating the results into a single result
				var result string
				if sub.Concat != "" {
					result = strings.Join(found, sub.Concat)
				} else {
					result = found[0]
				}

				result = sub.postProcess(ctx, result, ss)
				return next(ctx, result, q)
			}
		}

		return next(ctx, "", q)
	}
}

func postProcessMap(m map[string]string, next handlerFunc) handlerFunc {
	return func(ctx context.Context, value string, q queryer) string {
		// return the mapped value if present
		mapped, ok := m[value]

		if ok {
			return next(ctx, mapped, q)
		}

		return next(ctx, value, q)
	}
}

func postProcessFeetToCm(next handlerFunc) handlerFunc {
	const foot_in_cm = 30.48
	const inch_in_cm = 2.54

	reg := regexp.MustCompile("[0-9]+")

	return func(ctx context.Context, value string, q queryer) string {
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
		return next(ctx, strconv.Itoa(int(math.Round(centimeters))), q)
	}
}

func postProcessLbToKg(next handlerFunc) handlerFunc {
	const lb_in_kg = 0.45359237

	return func(ctx context.Context, value string, q queryer) string {
		w, err := strconv.ParseFloat(value, 64)
		if err == nil {
			w *= lb_in_kg
			value = strconv.Itoa(int(math.Round(w)))
		}
		return next(ctx, value, q)

	}
}

type mappedPostProcessAction struct {
	ParseDate    string                   `yaml:"parseDate"`
	SubtractDays bool                     `yaml:"subtractDays"`
	Replace      mappedRegexConfigs       `yaml:"replace"`
	SubScraper   *mappedScraperAttrConfig `yaml:"subScraper"`
	Map          map[string]string        `yaml:"map"`
	FeetToCm     bool                     `yaml:"feetToCm"`
	LbToKg       bool                     `yaml:"lbToKg"`
}

var ErrInvalid = errors.New("invalid post-process action")

func (a mappedPostProcessAction) parseHandlerFunc(next handlerFunc) (handlerFunc, error) {
	var found string
	var ret handlerFunc

	if a.ParseDate != "" {
		found = "parseDate"
		action := postProcessParseDate(a.ParseDate, next)
		ret = action
	}
	if len(a.Replace) > 0 {
		if found != "" {
			return nil, fmt.Errorf("post-process actions must have a single field, found %s and %s", found, "replace")
		}
		found = "replace"
		action := postProcessReplace(a.Replace, next)
		ret = action
	}
	if a.SubScraper != nil {
		if found != "" {
			return nil, fmt.Errorf("post-process actions must have a single field, found %s and %s", found, "subScraper")
		}
		found = "subScraper"
		action := postProcessSubScraper(*a.SubScraper, next)
		ret = action
	}
	if a.Map != nil {
		if found != "" {
			return nil, fmt.Errorf("post-process actions must have a single field, found %s and %s", found, "map")
		}
		found = "map"
		action := postProcessMap(a.Map, next)
		ret = action
	}
	if a.FeetToCm {
		if found != "" {
			return nil, fmt.Errorf("post-process actions must have a single field, found %s and %s", found, "feetToCm")
		}
		found = "feetToCm"
		action := postProcessFeetToCm(next)
		ret = action
	}
	if a.LbToKg {
		if found != "" {
			return nil, fmt.Errorf("post-process actions must have a single field, found %s and %s", found, "lbToKg")
		}
		found = "lbToKg"
		action := postProcessLbToKg(next)
		ret = action
	}
	if a.SubtractDays {
		if found != "" {
			return nil, fmt.Errorf("post-process actions must have a single field, found %s and %s", found, "subtractDays")
		}
		// found = "subtractDays"
		action := postProcessSubtractDays(next)
		ret = action
	}

	if ret == nil {
		return nil, ErrInvalid
	}

	return ret, nil
}

type mappedScraperAttrConfig struct {
	Selector    string                    `yaml:"selector"`
	Fixed       string                    `yaml:"fixed"`
	PostProcess []mappedPostProcessAction `yaml:"postProcess"`
	Concat      string                    `yaml:"concat"`
	Split       string                    `yaml:"split"`

	postProcessActions handlerFunc

	// Deprecated: use PostProcess instead
	ParseDate  string                   `yaml:"parseDate"`
	Replace    mappedRegexConfigs       `yaml:"replace"`
	SubScraper *mappedScraperAttrConfig `yaml:"subScraper"`
}

type _mappedScraperAttrConfig mappedScraperAttrConfig

func (c *mappedScraperAttrConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// try unmarshalling into a string first
	if err := unmarshal(&c.Selector); err != nil {
		// if it's a type error then we try to unmarshall to the full object
		var typeErr *yaml.TypeError
		if !errors.As(err, &typeErr) {
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

var ErrDeprecatedActions = errors.New("cannot include postProcess and (parseDate, replace, subScraper) deprecated fields")

// convertDeprecatedPostProcessActions handles the old deprecated fields
func (c *mappedScraperAttrConfig) convertDeprecatedPostProcessActions() error {
	// convert old deprecated fields if present in same order as they
	// used to be executed. Build up a chain in reverse
	action := done()

	if c.ParseDate != "" {
		action = postProcessParseDate(c.ParseDate, action)
		c.ParseDate = ""
	}

	if c.SubScraper != nil {
		action = postProcessSubScraper(*c.SubScraper, action)
		c.SubScraper = nil
	}

	if len(c.Replace) > 0 {
		action = postProcessReplace(c.Replace, action)
		c.Replace = nil
	}

	c.postProcessActions = action

	return nil
}

func (c *mappedScraperAttrConfig) convertPostProcessActions() error {
	// No post-process actions means we are in the old legacy case
	if len(c.PostProcess) == 0 {
		return c.convertDeprecatedPostProcessActions()
	}

	// ensure we don't have the old deprecated fields and the new post process field
	if c.ParseDate != "" || len(c.Replace) > 0 || c.SubScraper != nil {
		return ErrDeprecatedActions
	}

	// convert xpathPostProcessAction actions to postProcessActions
	// process in reverse order to make the action chain run forward
	// in the right order.

	action := done() // Final action just outputs the chain
	var err error
	for i := len(c.PostProcess) - 1; i >= 0; i-- {
		action, err = c.PostProcess[i].parseHandlerFunc(action)
		if err != nil {
			return err
		}
	}

	c.postProcessActions = action
	c.PostProcess = nil

	return nil
}

// cleanResults removes duplicate entries and empty values
func cleanResults(nodes []string) []string {
	cleaned := utils.StrUnique(nodes)      // remove duplicate values
	cleaned = utils.StrDelete(cleaned, "") // remove empty values
	return cleaned
}

// fieldSplit splits the value around sep and returns non-empty fields
func fieldSplit(value, sep string) []string {
	var res []string
	for _, str := range strings.Split(value, sep) {
		if str != "" {
			res = append(res, str)
		}
	}

	return res
}

func (c mappedScraperAttrConfig) postProcess(ctx context.Context, value string, q queryer) string {
	if c.postProcessActions == nil {
		return value
	}

	return c.postProcessActions.Apply(ctx, value, q)
}

type mappedResult map[string]string

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

type mappedResults []mappedResult

func (r mappedResults) setKey(index int, key string, value string) mappedResults {
	if index >= len(r) {
		r = append(r, make(mappedResult))
	}

	logger.Debugf(`[%d][%s] = %s`, index, key, value)
	r[index][key] = value
	return r
}

type mappedScraper struct {
	Common    commonMappedConfig            `yaml:"common"`
	Scene     *mappedSceneScraperConfig     `yaml:"scene"`
	Gallery   *mappedGalleryScraperConfig   `yaml:"gallery"`
	Performer *mappedPerformerScraperConfig `yaml:"performer"`
	Movie     *mappedMovieScraperConfig     `yaml:"movie"`
}

func (s mappedScraper) scrapePerformer(ctx context.Context, q queryer) (*models.ScrapedPerformer, error) {
	var ret models.ScrapedPerformer

	performerMap := s.Performer
	if performerMap == nil {
		return nil, nil
	}

	performerTagsMap := performerMap.Tags

	results := performerMap.process(ctx, q, s.Common)
	if len(results) > 0 {
		results[0].apply(&ret)

		// now apply the tags
		if performerTagsMap != nil {
			logger.Debug(`Processing performer tags:`)
			tagResults := performerTagsMap.process(ctx, q, s.Common)

			for _, p := range tagResults {
				tag := &models.ScrapedTag{}
				p.apply(tag)
				ret.Tags = append(ret.Tags, tag)
			}
		}
	}

	return &ret, nil
}

func (s mappedScraper) scrapePerformers(ctx context.Context, q queryer) ([]*models.ScrapedPerformer, error) {
	var ret []*models.ScrapedPerformer

	performerMap := s.Performer
	if performerMap == nil {
		return nil, nil
	}

	results := performerMap.process(ctx, q, s.Common)
	for _, r := range results {
		var p models.ScrapedPerformer
		r.apply(&p)
		ret = append(ret, &p)
	}

	return ret, nil
}

func (s mappedScraper) processScene(ctx context.Context, q queryer, r mappedResult) *models.ScrapedScene {
	var ret models.ScrapedScene

	sceneScraperConfig := s.Scene

	scenePerformersMap := sceneScraperConfig.Performers
	sceneTagsMap := sceneScraperConfig.Tags
	sceneStudioMap := sceneScraperConfig.Studio
	sceneMoviesMap := sceneScraperConfig.Movies

	scenePerformerTagsMap := scenePerformersMap.Tags

	r.apply(&ret)

	// process performer tags once
	var performerTagResults mappedResults
	if scenePerformerTagsMap != nil {
		performerTagResults = scenePerformerTagsMap.process(ctx, q, s.Common)
	}

	// now apply the performers and tags
	if scenePerformersMap.mappedConfig != nil {
		logger.Debug(`Processing scene performers:`)
		performerResults := scenePerformersMap.process(ctx, q, s.Common)

		for _, p := range performerResults {
			performer := &models.ScrapedPerformer{}
			p.apply(performer)

			for _, p := range performerTagResults {
				tag := &models.ScrapedTag{}
				p.apply(tag)
				ret.Tags = append(ret.Tags, tag)
			}

			ret.Performers = append(ret.Performers, performer)
		}
	}

	if sceneTagsMap != nil {
		logger.Debug(`Processing scene tags:`)
		tagResults := sceneTagsMap.process(ctx, q, s.Common)

		for _, p := range tagResults {
			tag := &models.ScrapedTag{}
			p.apply(tag)
			ret.Tags = append(ret.Tags, tag)
		}
	}

	if sceneStudioMap != nil {
		logger.Debug(`Processing scene studio:`)
		studioResults := sceneStudioMap.process(ctx, q, s.Common)

		if len(studioResults) > 0 {
			studio := &models.ScrapedStudio{}
			studioResults[0].apply(studio)
			ret.Studio = studio
		}
	}

	if sceneMoviesMap != nil {
		logger.Debug(`Processing scene movies:`)
		movieResults := sceneMoviesMap.process(ctx, q, s.Common)

		for _, p := range movieResults {
			movie := &models.ScrapedMovie{}
			p.apply(movie)
			ret.Movies = append(ret.Movies, movie)
		}
	}

	return &ret
}

func (s mappedScraper) scrapeScenes(ctx context.Context, q queryer) ([]*models.ScrapedScene, error) {
	var ret []*models.ScrapedScene

	sceneScraperConfig := s.Scene
	sceneMap := sceneScraperConfig.mappedConfig
	if sceneMap == nil {
		return nil, nil
	}

	logger.Debug(`Processing scenes:`)
	results := sceneMap.process(ctx, q, s.Common)
	for _, r := range results {
		logger.Debug(`Processing scene:`)
		ret = append(ret, s.processScene(ctx, q, r))
	}

	return ret, nil
}

func (s mappedScraper) scrapeScene(ctx context.Context, q queryer) (*models.ScrapedScene, error) {
	var ret models.ScrapedScene

	sceneScraperConfig := s.Scene
	sceneMap := sceneScraperConfig.mappedConfig
	if sceneMap == nil {
		return nil, nil
	}

	logger.Debug(`Processing scene:`)
	results := sceneMap.process(ctx, q, s.Common)
	if len(results) > 0 {
		ss := s.processScene(ctx, q, results[0])
		ret = *ss
	}

	return &ret, nil
}

func (s mappedScraper) scrapeGallery(ctx context.Context, q queryer) (*models.ScrapedGallery, error) {
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
	results := galleryMap.process(ctx, q, s.Common)
	if len(results) > 0 {
		results[0].apply(&ret)

		// now apply the performers and tags
		if galleryPerformersMap != nil {
			logger.Debug(`Processing gallery performers:`)
			performerResults := galleryPerformersMap.process(ctx, q, s.Common)

			for _, p := range performerResults {
				performer := &models.ScrapedPerformer{}
				p.apply(performer)
				ret.Performers = append(ret.Performers, performer)
			}
		}

		if galleryTagsMap != nil {
			logger.Debug(`Processing gallery tags:`)
			tagResults := galleryTagsMap.process(ctx, q, s.Common)

			for _, p := range tagResults {
				tag := &models.ScrapedTag{}
				p.apply(tag)
				ret.Tags = append(ret.Tags, tag)
			}
		}

		if galleryStudioMap != nil {
			logger.Debug(`Processing gallery studio:`)
			studioResults := galleryStudioMap.process(ctx, q, s.Common)

			if len(studioResults) > 0 {
				studio := &models.ScrapedStudio{}
				studioResults[0].apply(studio)
				ret.Studio = studio
			}
		}
	}

	return &ret, nil
}

func (s mappedScraper) scrapeMovie(ctx context.Context, q queryer) (*models.ScrapedMovie, error) {
	var ret models.ScrapedMovie

	movieScraperConfig := s.Movie
	movieMap := movieScraperConfig.mappedConfig
	if movieMap == nil {
		return nil, nil
	}

	movieStudioMap := movieScraperConfig.Studio

	results := movieMap.process(ctx, q, s.Common)
	if len(results) > 0 {
		results[0].apply(&ret)

		if movieStudioMap != nil {
			logger.Debug(`Processing movie studio:`)
			studioResults := movieStudioMap.process(ctx, q, s.Common)

			if len(studioResults) > 0 {
				studio := &models.ScrapedStudio{}
				studioResults[0].apply(studio)
				ret.Studio = studio
			}
		}
	}

	return &ret, nil
}
