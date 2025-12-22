package scraper

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/sliceutil"
	"gopkg.in/yaml.v2"
)

type commonMappedConfig map[string]string

type mappedConfig map[string]mappedScraperAttrConfig

func (s mappedConfig) applyCommon(c commonMappedConfig, src string) string {
	if c == nil {
		return src
	}

	ret := src
	for commonKey, commonVal := range c {
		ret = strings.ReplaceAll(ret, commonKey, commonVal)
	}

	return ret
}

// extractHostname parses a URL string and returns the hostname.
// Returns empty string if the URL cannot be parsed.
func extractHostname(urlStr string) string {
	if urlStr == "" {
		return ""
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		logger.Warnf("Error parsing URL '%s': %s", urlStr, err.Error())
		return ""
	}

	return u.Hostname()
}

type isMultiFunc func(key string) bool

func (s mappedConfig) process(ctx context.Context, q mappedQuery, common commonMappedConfig, isMulti isMultiFunc) mappedResults {
	var ret mappedResults

	for k, attrConfig := range s {

		if attrConfig.Fixed != "" {
			// TODO - not sure if this needs to set _all_ indexes for the key
			const i = 0
			// Support {inputURL} and {inputHostname} placeholders in fixed values
			value := strings.ReplaceAll(attrConfig.Fixed, "{inputURL}", q.getURL())
			value = strings.ReplaceAll(value, "{inputHostname}", extractHostname(q.getURL()))
			ret = ret.setSingleValue(i, k, value)
		} else {
			selector := attrConfig.Selector
			selector = s.applyCommon(common, selector)
			// Support {inputURL} and {inputHostname} placeholders in selectors
			selector = strings.ReplaceAll(selector, "{inputURL}", q.getURL())
			selector = strings.ReplaceAll(selector, "{inputHostname}", extractHostname(q.getURL()))

			found, err := q.runQuery(selector)
			if err != nil {
				logger.Warnf("key '%v': %v", k, err)
			}

			if len(found) > 0 {
				result := s.postProcess(ctx, q, attrConfig, found)

				// HACK - if the key is URLs, then we need to set the value as a multi-value
				isMulti := isMulti != nil && isMulti(k)
				if isMulti {
					ret = ret.setMultiValue(0, k, result)
				} else {
					for i, text := range result {
						ret = ret.setSingleValue(i, k, text)
					}
				}
			}
		}
	}

	return ret
}

func (s mappedConfig) postProcess(ctx context.Context, q mappedQuery, attrConfig mappedScraperAttrConfig, found []string) []string {
	// check if we're concatenating the results into a single result
	var ret []string
	if attrConfig.hasConcat() {
		result := attrConfig.concatenateResults(found)
		result = attrConfig.postProcess(ctx, result, q)
		if attrConfig.hasSplit() {
			results := attrConfig.splitString(result)
			// skip cleaning when the query is used for searching
			if q.getType() == SearchQuery {
				return results
			}
			results = attrConfig.cleanResults(results)
			return results
		}

		ret = []string{result}
	} else {
		for _, text := range found {
			text = attrConfig.postProcess(ctx, text, q)
			if attrConfig.hasSplit() {
				return attrConfig.splitString(text)
			}

			ret = append(ret, text)
		}
		// skip cleaning when the query is used for searching
		if q.getType() == SearchQuery {
			return ret
		}
		ret = attrConfig.cleanResults(ret)

	}

	return ret
}

type mappedSceneScraperConfig struct {
	mappedConfig

	Tags       mappedConfig                 `yaml:"Tags"`
	Performers mappedPerformerScraperConfig `yaml:"Performers"`
	Studio     mappedConfig                 `yaml:"Studio"`
	Movies     mappedConfig                 `yaml:"Movies"`
	Groups     mappedConfig                 `yaml:"Groups"`
}
type _mappedSceneScraperConfig mappedSceneScraperConfig

const (
	mappedScraperConfigSceneTags       = "Tags"
	mappedScraperConfigScenePerformers = "Performers"
	mappedScraperConfigSceneStudio     = "Studio"
	mappedScraperConfigSceneMovies     = "Movies"
	mappedScraperConfigSceneGroups     = "Groups"
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
	thisMap[mappedScraperConfigSceneGroups] = parentMap[mappedScraperConfigSceneGroups]

	delete(parentMap, mappedScraperConfigSceneTags)
	delete(parentMap, mappedScraperConfigScenePerformers)
	delete(parentMap, mappedScraperConfigSceneStudio)
	delete(parentMap, mappedScraperConfigSceneMovies)
	delete(parentMap, mappedScraperConfigSceneGroups)

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

type mappedImageScraperConfig struct {
	mappedConfig

	Tags       mappedConfig `yaml:"Tags"`
	Performers mappedConfig `yaml:"Performers"`
	Studio     mappedConfig `yaml:"Studio"`
}
type _mappedImageScraperConfig mappedImageScraperConfig

func (s *mappedImageScraperConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
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
	c := _mappedImageScraperConfig{}
	if err := yaml.Unmarshal(yml, &c); err != nil {
		return err
	}

	*s = mappedImageScraperConfig(c)

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
	Tags   mappedConfig `yaml:"Tags"`
}
type _mappedMovieScraperConfig mappedMovieScraperConfig

const (
	mappedScraperConfigMovieStudio = "Studio"
	mappedScraperConfigMovieTags   = "Tags"
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

	thisMap[mappedScraperConfigMovieTags] = parentMap[mappedScraperConfigMovieTags]
	delete(parentMap, mappedScraperConfigMovieTags)

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

type mappedScraperAttrConfig struct {
	Selector    string                    `yaml:"selector"`
	Fixed       string                    `yaml:"fixed"`
	PostProcess []mappedPostProcessAction `yaml:"postProcess"`
	Concat      string                    `yaml:"concat"`
	Split       string                    `yaml:"split"`

	postProcessActions []postProcessAction

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
	return strings.Join(nodes, separator)
}

func (c mappedScraperAttrConfig) cleanResults(nodes []string) []string {
	cleaned := sliceutil.Unique(nodes)      // remove duplicate values
	cleaned = sliceutil.Delete(cleaned, "") // remove empty values
	return cleaned
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

func (c mappedScraperAttrConfig) postProcess(ctx context.Context, value string, q mappedQuery) string {
	for _, action := range c.postProcessActions {
		value = action.Apply(ctx, value, q)
	}

	return value
}
