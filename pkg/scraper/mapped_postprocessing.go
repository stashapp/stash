package scraper

import (
	"context"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/javascript"
	"github.com/stashapp/stash/pkg/logger"
)

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
	Apply(ctx context.Context, value string, q mappedQuery) string
}

type postProcessParseDate string

func (p *postProcessParseDate) Apply(ctx context.Context, value string, q mappedQuery) string {
	parseDate := string(*p)

	const internalDateFormat = "2006-01-02"

	valueLower := strings.ToLower(value)
	if valueLower == "today" || valueLower == "yesterday" { // handle today, yesterday
		dt := time.Now()
		if valueLower == "yesterday" { // subtract 1 day from now
			dt = dt.AddDate(0, 0, -1)
		}
		return dt.Format(internalDateFormat)
	}

	if parseDate == "" {
		return value
	}

	if parseDate == "unix" {
		// try to parse the date using unix timestamp format
		// if it fails, then just fall back to the original value
		timeAsInt, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			logger.Warnf("Error parsing date string '%s' using unix timestamp format : %s", value, err.Error())
			return value
		}
		parsedValue := time.Unix(timeAsInt, 0)

		return parsedValue.Format(internalDateFormat)
	}

	// try to parse the date using the pattern
	// if it fails, then just fall back to the original value
	parsedValue, err := time.Parse(parseDate, value)
	if err != nil {
		logger.Warnf("Error parsing date string '%s' using format '%s': %s", value, parseDate, err.Error())
		return value
	}

	// convert it into our date format
	return parsedValue.Format(internalDateFormat)
}

type postProcessSubtractDays bool

func (p *postProcessSubtractDays) Apply(ctx context.Context, value string, q mappedQuery) string {
	const internalDateFormat = "2006-01-02"

	i, err := strconv.Atoi(value)
	if err != nil {
		logger.Warnf("Error parsing day string %s: %s", value, err)
		return value
	}

	dt := time.Now()
	dt = dt.AddDate(0, 0, -i)
	return dt.Format(internalDateFormat)
}

type postProcessReplace mappedRegexConfigs

func (c *postProcessReplace) Apply(ctx context.Context, value string, q mappedQuery) string {
	replace := mappedRegexConfigs(*c)
	return replace.apply(value)
}

type postProcessSubScraper mappedScraperAttrConfig

func (p *postProcessSubScraper) Apply(ctx context.Context, value string, q mappedQuery) string {
	subScrapeConfig := mappedScraperAttrConfig(*p)

	logger.Debugf("Sub-scraping for: %s", value)
	ss := q.subScrape(ctx, value)

	if ss != nil {
		found, err := ss.runQuery(subScrapeConfig.Selector)
		if err != nil {
			logger.Warnf("subscrape for '%v': %v", value, err)
		}

		if len(found) > 0 {
			// check if we're concatenating the results into a single result
			var result string
			if subScrapeConfig.hasConcat() {
				result = subScrapeConfig.concatenateResults(found)
			} else {
				result = found[0]
			}

			result = subScrapeConfig.postProcess(ctx, result, ss)
			return result
		}
	}

	return ""
}

type postProcessMap map[string]string

func (p *postProcessMap) Apply(ctx context.Context, value string, q mappedQuery) string {
	// return the mapped value if present
	m := *p
	mapped, ok := m[value]

	if ok {
		return mapped
	}

	return value
}

type postProcessFeetToCm bool

func (p *postProcessFeetToCm) Apply(ctx context.Context, value string, q mappedQuery) string {
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

type postProcessLbToKg bool

func (p *postProcessLbToKg) Apply(ctx context.Context, value string, q mappedQuery) string {
	const lb_in_kg = 0.45359237
	w, err := strconv.ParseFloat(value, 64)
	if err == nil {
		w *= lb_in_kg
		value = strconv.Itoa(int(math.Round(w)))
	}
	return value
}

type postProcessJavascript string

func (p *postProcessJavascript) Apply(ctx context.Context, value string, q mappedQuery) string {
	vm := javascript.NewVM()
	if err := vm.Set("value", value); err != nil {
		logger.Warnf("javascript failed to set value: %v", err)
		return value
	}

	log := &javascript.Log{
		Logger:       logger.Logger,
		Prefix:       "",
		ProgressChan: make(chan float64),
	}

	if err := log.AddToVM("log", vm); err != nil {
		logger.Logger.Errorf("error adding log API: %w", err)
	}

	util := &javascript.Util{}
	if err := util.AddToVM("util", vm); err != nil {
		logger.Logger.Errorf("error adding util API: %w", err)
	}

	script, err := javascript.CompileScript("", "(function() { "+string(*p)+"})()")
	if err != nil {
		logger.Warnf("javascript failed to compile: %v", err)
		return value
	}

	output, err := vm.RunProgram(script)
	if err != nil {
		logger.Warnf("javascript failed to run: %v", err)
		return value
	}

	// assume output is string
	return output.String()
}

type mappedPostProcessAction struct {
	ParseDate    string                   `yaml:"parseDate"`
	SubtractDays bool                     `yaml:"subtractDays"`
	Replace      mappedRegexConfigs       `yaml:"replace"`
	SubScraper   *mappedScraperAttrConfig `yaml:"subScraper"`
	Map          map[string]string        `yaml:"map"`
	FeetToCm     bool                     `yaml:"feetToCm"`
	LbToKg       bool                     `yaml:"lbToKg"`
	Javascript   string                   `yaml:"javascript"`
}

func (a mappedPostProcessAction) ToPostProcessAction() (postProcessAction, error) {
	var found string
	var ret postProcessAction

	ensureOnly := func(field string) error {
		if found != "" {
			return fmt.Errorf("post-process actions must have a single field, found %s and %s", found, field)
		}
		found = field
		return nil
	}

	if a.ParseDate != "" {
		found = "parseDate"
		action := postProcessParseDate(a.ParseDate)
		ret = &action
	}
	if len(a.Replace) > 0 {
		if err := ensureOnly("replace"); err != nil {
			return nil, err
		}
		action := postProcessReplace(a.Replace)
		ret = &action
	}
	if a.SubScraper != nil {
		if err := ensureOnly("subScraper"); err != nil {
			return nil, err
		}
		action := postProcessSubScraper(*a.SubScraper)
		ret = &action
	}
	if a.Map != nil {
		if err := ensureOnly("map"); err != nil {
			return nil, err
		}
		action := postProcessMap(a.Map)
		ret = &action
	}
	if a.FeetToCm {
		if err := ensureOnly("feetToCm"); err != nil {
			return nil, err
		}
		action := postProcessFeetToCm(a.FeetToCm)
		ret = &action
	}
	if a.LbToKg {
		if err := ensureOnly("lbToKg"); err != nil {
			return nil, err
		}
		action := postProcessLbToKg(a.LbToKg)
		ret = &action
	}
	if a.SubtractDays {
		if err := ensureOnly("subtractDays"); err != nil {
			return nil, err
		}
		action := postProcessSubtractDays(a.SubtractDays)
		ret = &action
	}
	if a.Javascript != "" {
		if err := ensureOnly("javascript"); err != nil {
			return nil, err
		}
		action := postProcessJavascript(a.Javascript)
		ret = &action
	}

	if ret == nil {
		return nil, errors.New("invalid post-process action")
	}

	return ret, nil
}
