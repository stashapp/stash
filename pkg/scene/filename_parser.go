package scene

import (
	"context"
	"errors"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/studio"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/tag"
)

type parserField struct {
	field           string
	fieldRegex      *regexp.Regexp
	regex           string
	isFullDateField bool
	isCaptured      bool
}

func newParserField(field string, regex string, captured bool) parserField {
	ret := parserField{
		field:           field,
		isFullDateField: false,
		isCaptured:      captured,
	}

	ret.fieldRegex, _ = regexp.Compile(`\{` + ret.field + `\}`)

	regexStr := regex

	if captured {
		regexStr = "(" + regexStr + ")"
	}
	ret.regex = regexStr

	return ret
}

func newFullDateParserField(field string, regex string) parserField {
	ret := newParserField(field, regex, true)
	ret.isFullDateField = true
	return ret
}

func (f parserField) replaceInPattern(pattern string) string {
	return string(f.fieldRegex.ReplaceAllString(pattern, f.regex))
}

var validFields map[string]parserField
var escapeCharRE *regexp.Regexp
var capitalizeTitleRE *regexp.Regexp
var multiWSRE *regexp.Regexp
var delimiterRE *regexp.Regexp

func compileREs() {
	const escapeCharPattern = `([\-\.\(\)\[\]])`
	escapeCharRE = regexp.MustCompile(escapeCharPattern)

	const capitaliseTitlePattern = `(?:^| )\w`
	capitalizeTitleRE = regexp.MustCompile(capitaliseTitlePattern)

	const multiWSPattern = ` {2,}`
	multiWSRE = regexp.MustCompile(multiWSPattern)

	const delimiterPattern = `(?:\.|-|_)`
	delimiterRE = regexp.MustCompile(delimiterPattern)
}

func initParserFields() {
	if validFields != nil {
		return
	}

	ret := make(map[string]parserField)

	ret["title"] = newParserField("title", ".*", true)
	ret["ext"] = newParserField("ext", ".*$", false)

	ret["d"] = newParserField("d", `(?:\.|-|_)`, false)
	ret["rating"] = newParserField("rating", `\d`, true)
	ret["rating100"] = newParserField("rating100", `\d`, true)
	ret["performer"] = newParserField("performer", ".*", true)
	ret["studio"] = newParserField("studio", ".*", true)
	ret["movie"] = newParserField("movie", ".*", true)
	ret["tag"] = newParserField("tag", ".*", true)

	// date fields
	ret["date"] = newParserField("date", `\d{4}-\d{2}-\d{2}`, true)
	ret["yyyy"] = newParserField("yyyy", `\d{4}`, true)
	ret["yy"] = newParserField("yy", `\d{2}`, true)
	ret["mm"] = newParserField("mm", `\d{2}`, true)
	ret["mmm"] = newParserField("mmm", `\w{3}`, true)
	ret["dd"] = newParserField("dd", `\d{2}`, true)
	ret["yyyymmdd"] = newFullDateParserField("yyyymmdd", `\d{8}`)
	ret["yymmdd"] = newFullDateParserField("yymmdd", `\d{6}`)
	ret["ddmmyyyy"] = newFullDateParserField("ddmmyyyy", `\d{8}`)
	ret["ddmmyy"] = newFullDateParserField("ddmmyy", `\d{6}`)
	ret["mmddyyyy"] = newFullDateParserField("mmddyyyy", `\d{8}`)
	ret["mmddyy"] = newFullDateParserField("mmddyy", `\d{6}`)

	validFields = ret
}

func replacePatternWithRegex(pattern string, ignoreWords []string) string {
	initParserFields()

	for _, field := range validFields {
		pattern = field.replaceInPattern(pattern)
	}

	ignoreClause := getIgnoreClause(ignoreWords)
	ignoreField := newParserField("i", ignoreClause, false)
	pattern = ignoreField.replaceInPattern(pattern)

	return pattern
}

type parseMapper struct {
	fields      []string
	regexString string
	regex       *regexp.Regexp
}

func getIgnoreClause(ignoreFields []string) string {
	if len(ignoreFields) == 0 {
		return ""
	}

	var ignoreClauses []string

	for _, v := range ignoreFields {
		newVal := string(escapeCharRE.ReplaceAllString(v, `\$1`))
		newVal = strings.TrimSpace(newVal)
		newVal = "(?:" + newVal + ")"
		ignoreClauses = append(ignoreClauses, newVal)
	}

	return "(?:" + strings.Join(ignoreClauses, "|") + ")"
}

func newParseMapper(pattern string, ignoreFields []string) (*parseMapper, error) {
	ret := &parseMapper{}

	// escape control characters
	regex := escapeCharRE.ReplaceAllString(pattern, `\$1`)

	// replace {} with wildcard
	braceRE := regexp.MustCompile(`\{\}`)
	regex = braceRE.ReplaceAllString(regex, ".*")

	// replace all known fields with applicable regexes
	regex = replacePatternWithRegex(regex, ignoreFields)

	ret.regexString = regex

	// make case insensitive
	regex = "(?i)" + regex

	var err error

	ret.regex, err = regexp.Compile(regex)

	if err != nil {
		return nil, err
	}

	// find invalid fields
	invalidRE := regexp.MustCompile(`\{[A-Za-z]+\}`)
	foundInvalid := invalidRE.FindAllString(regex, -1)
	if len(foundInvalid) > 0 {
		return nil, errors.New("Invalid fields: " + strings.Join(foundInvalid, ", "))
	}

	fieldExtractor := regexp.MustCompile(`\{([A-Za-z]+)\}`)

	result := fieldExtractor.FindAllStringSubmatch(pattern, -1)

	var fields []string
	for _, v := range result {
		field := v[1]

		// only add to fields if it is captured
		parserField, found := validFields[field]
		if found && parserField.isCaptured {
			fields = append(fields, field)
		}
	}

	ret.fields = fields

	return ret, nil
}

// fileHolder accumulates the field values matched against a single filename. It is
// deliberately media-agnostic so that scenes and audios share the same parsing machinery.
type fileHolder struct {
	// the date already stored against the object being parsed. A parsed date is only
	// reported when it differs from this.
	existingDate *models.Date

	title      string
	date       *models.Date
	rating     *int
	performers []string
	groups     []string
	studio     string
	tags       []string

	yyyy string
	mm   string
	dd   string
}

func newFileHolder(existingDate *models.Date) *fileHolder {
	return &fileHolder{
		existingDate: existingDate,
	}
}

func validateRating(rating int) bool {
	return rating >= 1 && rating <= 5
}

func validateRating100(rating100 int) bool {
	return rating100 >= 1 && rating100 <= 100
}

// returns nil if invalid
func parseDate(dateStr string) *models.Date {
	splits := strings.Split(dateStr, "-")
	if len(splits) != 3 {
		return nil
	}

	year, _ := strconv.Atoi(splits[0])
	month, _ := strconv.Atoi(splits[1])
	d, _ := strconv.Atoi(splits[2])

	// assume year must be between 1900 and 2100
	if year < 1900 || year > 2100 {
		return nil
	}

	if month < 1 || month > 12 {
		return nil
	}

	// not checking individual months to ensure date is in the correct range
	if d < 1 || d > 31 {
		return nil
	}

	ret, err := models.ParseDate(dateStr)
	if err != nil {
		return nil
	}
	return &ret
}

func (h *fileHolder) setDate(field *parserField, value string) {
	yearIndex := 0
	yearLength := len(strings.Split(field.field, "y")) - 1
	dateIndex := 0
	monthIndex := 0

	switch field.field {
	case "yyyymmdd", "yymmdd":
		monthIndex = yearLength
		dateIndex = monthIndex + 2
	case "ddmmyyyy", "ddmmyy":
		monthIndex = 2
		yearIndex = monthIndex + 2
	case "mmddyyyy", "mmddyy":
		dateIndex = monthIndex + 2
		yearIndex = dateIndex + 2
	}

	yearValue := value[yearIndex : yearIndex+yearLength]
	monthValue := value[monthIndex : monthIndex+2]
	dateValue := value[dateIndex : dateIndex+2]

	fullDate := yearValue + "-" + monthValue + "-" + dateValue

	// ensure the date is valid
	// only set if new value is different from the old
	newDate := parseDate(fullDate)
	// NOTE - matches the long-standing scene parser behaviour: a full-date field is only
	// reported when the object already has a date set and the parsed date differs from it.
	if newDate != nil && h.existingDate != nil && *h.existingDate != *newDate {
		h.date = newDate
	}
}

func mmmToMonth(mmm string) string {
	format := "02-Jan-2006"
	dateStr := "01-" + mmm + "-2000"
	t, err := time.Parse(format, dateStr)

	if err != nil {
		return ""
	}

	// expect month in two-digit format
	format = "01-02-2006"
	return t.Format(format)[0:2]
}

func (h *fileHolder) setField(field parserField, value interface{}) {
	if field.isFullDateField {
		h.setDate(&field, value.(string))
		return
	}

	switch field.field {
	case "title":
		h.title = value.(string)
	case "date":
		h.date = parseDate(value.(string))
	case "rating":
		rating, _ := strconv.Atoi(value.(string))
		if validateRating(rating) {
			// convert to 1-100 scale
			rating = models.Rating5To100(rating)
			h.rating = &rating
		}
	case "rating100":
		rating, _ := strconv.Atoi(value.(string))
		if validateRating100(rating) {
			h.rating = &rating
		}
	case "performer":
		// add performer to list
		h.performers = append(h.performers, value.(string))
	case "studio":
		h.studio = value.(string)
	case "movie":
		h.groups = append(h.groups, value.(string))
	case "tag":
		h.tags = append(h.tags, value.(string))
	case "yyyy":
		h.yyyy = value.(string)
	case "yy":
		v := value.(string)
		v = "20" + v
		h.yyyy = v
	case "mmm":
		h.mm = mmmToMonth(value.(string))
	case "mm":
		h.mm = value.(string)
	case "dd":
		h.dd = value.(string)
	}
}

func (h *fileHolder) postParse() {
	// set the date if the components are set
	if h.yyyy != "" && h.mm != "" && h.dd != "" {
		fullDate := h.yyyy + "-" + h.mm + "-" + h.dd
		h.setField(validFields["date"], fullDate)
	}
}

func (m parseMapper) parse(path string, existingDate *models.Date) *fileHolder {

	// #302 - if the pattern includes a path separator, then include the entire
	// path in the match. Otherwise, use the default behaviour of just
	// the file's basename
	// must be double \ because of the regex escaping
	filename := filepath.Base(path)
	if strings.Contains(m.regexString, `\\`) || strings.Contains(m.regexString, "/") {
		filename = path
	}

	result := m.regex.FindStringSubmatch(filename)

	if len(result) == 0 {
		return nil
	}

	initParserFields()

	holder := newFileHolder(existingDate)

	for index, match := range result {
		if index == 0 {
			// skip entire match
			continue
		}

		field := m.fields[index-1]
		parserField, found := validFields[field]
		if found {
			holder.setField(parserField, match)
		}
	}

	holder.postParse()

	return holder
}

// parserOptions are the parsing options shared by the scene and audio parsers.
type parserOptions struct {
	IgnoreWords          []string
	WhitespaceCharacters *string
	CapitalizeTitle      *bool
	IgnoreOrganized      *bool
}

type FilenameParser struct {
	Pattern        string
	Options        parserOptions
	Filter         *models.FindFilterType
	whitespaceRE   *regexp.Regexp
	repository     FilenameParserRepository
	performerCache map[string]*models.Performer
	studioCache    map[string]*models.Studio
	groupCache     map[string]*models.Group
	tagCache       map[string]*models.Tag
}

func NewFilenameParser(filter *models.FindFilterType, config models.SceneParserInput, repo FilenameParserRepository) *FilenameParser {
	return newFilenameParser(filter, parserOptions{
		IgnoreWords:          config.IgnoreWords,
		WhitespaceCharacters: config.WhitespaceCharacters,
		CapitalizeTitle:      config.CapitalizeTitle,
		IgnoreOrganized:      config.IgnoreOrganized,
	}, repo)
}

// NewAudioFilenameParser returns a parser that matches audio filenames. It lives alongside the
// scene parser because both share the same field-matching machinery.
func NewAudioFilenameParser(filter *models.FindFilterType, config models.AudioParserInput, repo FilenameParserRepository) *FilenameParser {
	return newFilenameParser(filter, parserOptions{
		IgnoreWords:          config.IgnoreWords,
		WhitespaceCharacters: config.WhitespaceCharacters,
		CapitalizeTitle:      config.CapitalizeTitle,
		IgnoreOrganized:      config.IgnoreOrganized,
	}, repo)
}

func newFilenameParser(filter *models.FindFilterType, options parserOptions, repo FilenameParserRepository) *FilenameParser {
	p := &FilenameParser{
		Pattern:    *filter.Q,
		Options:    options,
		Filter:     filter,
		repository: repo,
	}

	p.performerCache = make(map[string]*models.Performer)
	p.studioCache = make(map[string]*models.Studio)
	p.groupCache = make(map[string]*models.Group)
	p.tagCache = make(map[string]*models.Tag)

	p.initWhiteSpaceRegex()

	return p
}

func (p *FilenameParser) initWhiteSpaceRegex() {
	compileREs()

	wsChars := ""
	if p.Options.WhitespaceCharacters != nil {
		wsChars = *p.Options.WhitespaceCharacters
		wsChars = strings.TrimSpace(wsChars)
	}

	if len(wsChars) > 0 {
		wsRegExp := escapeCharRE.ReplaceAllString(wsChars, `\$1`)
		wsRegExp = "[" + wsRegExp + "]"
		p.whitespaceRE = regexp.MustCompile(wsRegExp)
	}
}

type FilenameParserRepository struct {
	Scene     models.SceneQueryer
	Audio     models.AudioQueryer
	Performer PerformerNamesFinder
	Studio    models.StudioQueryer
	Group     GroupNameFinder
	Tag       models.TagNameFinder
}

func NewFilenameParserRepository(repo models.Repository) FilenameParserRepository {
	return FilenameParserRepository{
		Scene:     repo.Scene,
		Audio:     repo.Audio,
		Performer: repo.Performer,
		Studio:    repo.Studio,
		Group:     repo.Group,
		Tag:       repo.Tag,
	}
}

func (p *FilenameParser) Parse(ctx context.Context) ([]*models.SceneParserResult, int, error) {
	// perform the query to find the scenes
	mapper, err := newParseMapper(p.Pattern, p.Options.IgnoreWords)

	if err != nil {
		return nil, 0, err
	}

	sceneFilter := &models.SceneFilterType{
		Path: &models.StringCriterionInput{
			Modifier: models.CriterionModifierMatchesRegex,
			Value:    "(?i)" + mapper.regexString,
		},
	}

	if p.Options.IgnoreOrganized != nil && *p.Options.IgnoreOrganized {
		organized := false
		sceneFilter.Organized = &organized
	}

	p.Filter.Q = nil

	scenes, total, err := QueryWithCount(ctx, p.repository.Scene, sceneFilter, p.Filter)
	if err != nil {
		return nil, 0, err
	}

	ret := p.parseScenes(ctx, scenes, mapper)

	return ret, total, nil
}

// ParseAudios parses audio filenames using the same pattern syntax as the scene parser.
func (p *FilenameParser) ParseAudios(ctx context.Context) ([]*models.AudioParserResult, int, error) {
	mapper, err := newParseMapper(p.Pattern, p.Options.IgnoreWords)

	if err != nil {
		return nil, 0, err
	}

	audioFilter := &models.AudioFilterType{
		Path: &models.StringCriterionInput{
			Modifier: models.CriterionModifierMatchesRegex,
			Value:    "(?i)" + mapper.regexString,
		},
	}

	if p.Options.IgnoreOrganized != nil && *p.Options.IgnoreOrganized {
		organized := false
		audioFilter.Organized = &organized
	}

	p.Filter.Q = nil

	total, err := p.repository.Audio.QueryCount(ctx, audioFilter, p.Filter)
	if err != nil {
		return nil, 0, err
	}

	audios, err := audioQuery(ctx, p.repository.Audio, audioFilter, p.Filter)
	if err != nil {
		return nil, 0, err
	}

	ret := p.parseAudios(ctx, audios, mapper)

	return ret, total, nil
}

func (p *FilenameParser) parseScenes(ctx context.Context, scenes []*models.Scene, mapper *parseMapper) []*models.SceneParserResult {
	var ret []*models.SceneParserResult
	for _, scene := range scenes {
		holder := mapper.parse(scene.Path, scene.Date)

		if holder != nil {
			r := &models.SceneParserResult{
				Scene: scene,
			}

			f := p.resolveFields(ctx, *holder)
			r.Title = f.title
			r.Date = f.date
			r.Rating = f.rating
			r.PerformerIds = f.performerIDs
			r.TagIds = f.tagIDs
			r.StudioID = f.studioID
			for _, groupID := range f.groupIDs {
				r.Movies = append(r.Movies, &models.SceneMovieID{
					MovieID: groupID,
				})
			}

			ret = append(ret, r)
		}
	}

	return ret
}

func (p *FilenameParser) parseAudios(ctx context.Context, audios []*models.Audio, mapper *parseMapper) []*models.AudioParserResult {
	var ret []*models.AudioParserResult
	for _, audio := range audios {
		holder := mapper.parse(audio.Path, audio.Date)

		if holder != nil {
			r := &models.AudioParserResult{
				Audio: audio,
			}

			f := p.resolveFields(ctx, *holder)
			r.Title = f.title
			r.Date = f.date
			r.Rating100 = f.rating
			r.PerformerIds = f.performerIDs
			r.TagIds = f.tagIDs
			r.StudioID = f.studioID
			for _, groupID := range f.groupIDs {
				r.Groups = append(r.Groups, &models.AudioGroupID{
					GroupID: groupID,
				})
			}

			ret = append(ret, r)
		}
	}

	return ret
}

func (p FilenameParser) replaceWhitespaceCharacters(value string) string {
	if p.whitespaceRE != nil {
		value = p.whitespaceRE.ReplaceAllString(value, " ")
		// remove consecutive spaces
		value = multiWSRE.ReplaceAllString(value, " ")
	}

	return value
}

type PerformerNamesFinder interface {
	FindByNames(ctx context.Context, names []string, nocase bool) ([]*models.Performer, error)
}

func (p *FilenameParser) queryPerformer(ctx context.Context, qb PerformerNamesFinder, performerName string) *models.Performer {
	// massage the performer name
	performerName = delimiterRE.ReplaceAllString(performerName, " ")

	// check cache first
	if ret, found := p.performerCache[performerName]; found {
		return ret
	}

	// perform an exact match and grab the first
	performers, _ := qb.FindByNames(ctx, []string{performerName}, true)

	var ret *models.Performer
	if len(performers) > 0 {
		ret = performers[0]
	}

	// add result to cache
	p.performerCache[performerName] = ret

	return ret
}

func (p *FilenameParser) queryStudio(ctx context.Context, qb models.StudioQueryer, studioName string) *models.Studio {
	// massage the performer name
	studioName = delimiterRE.ReplaceAllString(studioName, " ")

	// check cache first
	if ret, found := p.studioCache[studioName]; found {
		return ret
	}

	ret, _ := studio.ByName(ctx, qb, studioName)

	// try to match on alias
	if ret == nil {
		ret, _ = studio.ByAlias(ctx, qb, studioName)
	}

	// add result to cache
	p.studioCache[studioName] = ret

	return ret
}

type GroupNameFinder interface {
	FindByName(ctx context.Context, name string, nocase bool) (*models.Group, error)
}

func (p *FilenameParser) queryGroup(ctx context.Context, qb GroupNameFinder, groupName string) *models.Group {
	// massage the group name
	groupName = delimiterRE.ReplaceAllString(groupName, " ")

	// check cache first
	if ret, found := p.groupCache[groupName]; found {
		return ret
	}

	ret, _ := qb.FindByName(ctx, groupName, true)

	// add result to cache
	p.groupCache[groupName] = ret

	return ret
}

func (p *FilenameParser) queryTag(ctx context.Context, qb models.TagNameFinder, tagName string) *models.Tag {
	// massage the tag name
	tagName = delimiterRE.ReplaceAllString(tagName, " ")

	// check cache first
	if ret, found := p.tagCache[tagName]; found {
		return ret
	}

	// match tag name exactly
	ret, _ := tag.ByName(ctx, qb, tagName)

	// try to match on alias
	if ret == nil {
		ret, _ = tag.ByAlias(ctx, qb, tagName)
	}

	// add result to cache
	p.tagCache[tagName] = ret

	return ret
}

func (p *FilenameParser) performerIDs(ctx context.Context, qb PerformerNamesFinder, h fileHolder) []string {
	// query for each performer
	var ret []string
	performersSet := make(map[int]bool)
	for _, performerName := range h.performers {
		if performerName != "" {
			performer := p.queryPerformer(ctx, qb, performerName)
			if performer != nil {
				if _, found := performersSet[performer.ID]; !found {
					ret = append(ret, strconv.Itoa(performer.ID))
					performersSet[performer.ID] = true
				}
			}
		}
	}

	return ret
}

func (p *FilenameParser) tagIDs(ctx context.Context, qb models.TagNameFinder, h fileHolder) []string {
	// query for each tag
	var ret []string
	tagsSet := make(map[int]bool)
	for _, tagName := range h.tags {
		if tagName != "" {
			tag := p.queryTag(ctx, qb, tagName)
			if tag != nil {
				if _, found := tagsSet[tag.ID]; !found {
					ret = append(ret, strconv.Itoa(tag.ID))
					tagsSet[tag.ID] = true
				}
			}
		}
	}

	return ret
}

func (p *FilenameParser) studioID(ctx context.Context, qb models.StudioQueryer, h fileHolder) *string {
	if h.studio == "" {
		return nil
	}

	studio := p.queryStudio(ctx, qb, h.studio)
	if studio == nil {
		return nil
	}

	studioID := strconv.Itoa(studio.ID)
	return &studioID
}

func (p *FilenameParser) groupIDs(ctx context.Context, qb GroupNameFinder, h fileHolder) []string {
	// query for each group
	var ret []string
	groupsSet := make(map[int]bool)
	for _, groupName := range h.groups {
		if groupName != "" {
			group := p.queryGroup(ctx, qb, groupName)
			if group != nil {
				if _, found := groupsSet[group.ID]; !found {
					ret = append(ret, strconv.Itoa(group.ID))
					groupsSet[group.ID] = true
				}
			}
		}
	}

	return ret
}

// resolvedFields is the media-agnostic outcome of parsing a single filename, with names
// already resolved to object ids.
type resolvedFields struct {
	title        *string
	date         *string
	rating       *int
	performerIDs []string
	tagIDs       []string
	studioID     *string
	groupIDs     []string
}

func (p *FilenameParser) resolveFields(ctx context.Context, h fileHolder) resolvedFields {
	var ret resolvedFields

	if h.title != "" {
		title := p.replaceWhitespaceCharacters(h.title)

		if p.Options.CapitalizeTitle != nil && *p.Options.CapitalizeTitle {
			title = capitalizeTitleRE.ReplaceAllStringFunc(title, strings.ToUpper)
		}

		ret.title = &title
	}

	if h.date != nil {
		dateStr := h.date.String()
		ret.date = &dateStr
	}

	ret.rating = h.rating

	r := p.repository

	if len(h.performers) > 0 {
		ret.performerIDs = p.performerIDs(ctx, r.Performer, h)
	}
	if len(h.tags) > 0 {
		ret.tagIDs = p.tagIDs(ctx, r.Tag, h)
	}
	ret.studioID = p.studioID(ctx, r.Studio, h)

	if len(h.groups) > 0 {
		ret.groupIDs = p.groupIDs(ctx, r.Group, h)
	}

	return ret
}

// audioQuery runs an audio query and resolves the results.
func audioQuery(ctx context.Context, qb models.AudioQueryer, audioFilter *models.AudioFilterType, findFilter *models.FindFilterType) ([]*models.Audio, error) {
	result, err := qb.Query(ctx, models.AudioQueryOptions{
		QueryOptions: models.QueryOptions{
			FindFilter: findFilter,
		},
		AudioFilter: audioFilter,
	})
	if err != nil {
		return nil, err
	}

	return result.Resolve(ctx)
}
