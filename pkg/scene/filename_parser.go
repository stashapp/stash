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

type sceneHolder struct {
	scene      *models.Scene
	result     *models.Scene
	yyyy       string
	mm         string
	dd         string
	performers []string
	groups     []string
	studio     string
	tags       []string
}

func newSceneHolder(scene *models.Scene) *sceneHolder {
	sceneCopy := models.Scene{
		ID:    scene.ID,
		Files: scene.Files,
		// Checksum: scene.Checksum,
		// Path:     scene.Path,
	}
	ret := sceneHolder{
		scene:  scene,
		result: &sceneCopy,
	}

	return &ret
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

func (h *sceneHolder) setDate(field *parserField, value string) {
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
	if newDate != nil && h.scene.Date != nil && *h.scene.Date != *newDate {
		h.result.Date = newDate
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

func (h *sceneHolder) setField(field parserField, value interface{}) {
	if field.isFullDateField {
		h.setDate(&field, value.(string))
		return
	}

	switch field.field {
	case "title":
		v := value.(string)
		h.result.Title = v
	case "date":
		h.result.Date = parseDate(value.(string))
	case "rating":
		rating, _ := strconv.Atoi(value.(string))
		if validateRating(rating) {
			// convert to 1-100 scale
			rating = models.Rating5To100(rating)
			h.result.Rating = &rating
		}
	case "rating100":
		rating, _ := strconv.Atoi(value.(string))
		if validateRating100(rating) {
			h.result.Rating = &rating
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

func (h *sceneHolder) postParse() {
	// set the date if the components are set
	if h.yyyy != "" && h.mm != "" && h.dd != "" {
		fullDate := h.yyyy + "-" + h.mm + "-" + h.dd
		h.setField(validFields["date"], fullDate)
	}
}

func (m parseMapper) parse(scene *models.Scene) *sceneHolder {

	// #302 - if the pattern includes a path separator, then include the entire
	// scene path in the match. Otherwise, use the default behaviour of just
	// the file's basename
	// must be double \ because of the regex escaping
	filename := filepath.Base(scene.Path)
	if strings.Contains(m.regexString, `\\`) || strings.Contains(m.regexString, "/") {
		filename = scene.Path
	}

	result := m.regex.FindStringSubmatch(filename)

	if len(result) == 0 {
		return nil
	}

	initParserFields()

	sceneHolder := newSceneHolder(scene)

	for index, match := range result {
		if index == 0 {
			// skip entire match
			continue
		}

		field := m.fields[index-1]
		parserField, found := validFields[field]
		if found {
			sceneHolder.setField(parserField, match)
		}
	}

	sceneHolder.postParse()

	return sceneHolder
}

type FilenameParser struct {
	Pattern        string
	ParserInput    models.SceneParserInput
	Filter         *models.FindFilterType
	whitespaceRE   *regexp.Regexp
	repository     FilenameParserRepository
	performerCache map[string]*models.Performer
	studioCache    map[string]*models.Studio
	groupCache     map[string]*models.Group
	tagCache       map[string]*models.Tag
}

func NewFilenameParser(filter *models.FindFilterType, config models.SceneParserInput, repo FilenameParserRepository) *FilenameParser {
	p := &FilenameParser{
		Pattern:     *filter.Q,
		ParserInput: config,
		Filter:      filter,
		repository:  repo,
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
	if p.ParserInput.WhitespaceCharacters != nil {
		wsChars = *p.ParserInput.WhitespaceCharacters
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
	Performer PerformerNamesFinder
	Studio    models.StudioQueryer
	Group     GroupNameFinder
	Tag       models.TagQueryer
}

func NewFilenameParserRepository(repo models.Repository) FilenameParserRepository {
	return FilenameParserRepository{
		Scene:     repo.Scene,
		Performer: repo.Performer,
		Studio:    repo.Studio,
		Group:     repo.Group,
		Tag:       repo.Tag,
	}
}

func (p *FilenameParser) Parse(ctx context.Context) ([]*models.SceneParserResult, int, error) {
	// perform the query to find the scenes
	mapper, err := newParseMapper(p.Pattern, p.ParserInput.IgnoreWords)

	if err != nil {
		return nil, 0, err
	}

	sceneFilter := &models.SceneFilterType{
		Path: &models.StringCriterionInput{
			Modifier: models.CriterionModifierMatchesRegex,
			Value:    "(?i)" + mapper.regexString,
		},
	}

	if p.ParserInput.IgnoreOrganized != nil && *p.ParserInput.IgnoreOrganized {
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

func (p *FilenameParser) parseScenes(ctx context.Context, scenes []*models.Scene, mapper *parseMapper) []*models.SceneParserResult {
	var ret []*models.SceneParserResult
	for _, scene := range scenes {
		sceneHolder := mapper.parse(scene)

		if sceneHolder != nil {
			r := &models.SceneParserResult{
				Scene: scene,
			}
			p.setParserResult(ctx, *sceneHolder, r)

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

func (p *FilenameParser) queryTag(ctx context.Context, qb models.TagQueryer, tagName string) *models.Tag {
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

func (p *FilenameParser) setPerformers(ctx context.Context, qb PerformerNamesFinder, h sceneHolder, result *models.SceneParserResult) {
	// query for each performer
	performersSet := make(map[int]bool)
	for _, performerName := range h.performers {
		if performerName != "" {
			performer := p.queryPerformer(ctx, qb, performerName)
			if performer != nil {
				if _, found := performersSet[performer.ID]; !found {
					result.PerformerIds = append(result.PerformerIds, strconv.Itoa(performer.ID))
					performersSet[performer.ID] = true
				}
			}
		}
	}
}

func (p *FilenameParser) setTags(ctx context.Context, qb models.TagQueryer, h sceneHolder, result *models.SceneParserResult) {
	// query for each performer
	tagsSet := make(map[int]bool)
	for _, tagName := range h.tags {
		if tagName != "" {
			tag := p.queryTag(ctx, qb, tagName)
			if tag != nil {
				if _, found := tagsSet[tag.ID]; !found {
					result.TagIds = append(result.TagIds, strconv.Itoa(tag.ID))
					tagsSet[tag.ID] = true
				}
			}
		}
	}
}

func (p *FilenameParser) setStudio(ctx context.Context, qb models.StudioQueryer, h sceneHolder, result *models.SceneParserResult) {
	// query for each performer
	if h.studio != "" {
		studio := p.queryStudio(ctx, qb, h.studio)
		if studio != nil {
			studioID := strconv.Itoa(studio.ID)
			result.StudioID = &studioID
		}
	}
}

func (p *FilenameParser) setGroups(ctx context.Context, qb GroupNameFinder, h sceneHolder, result *models.SceneParserResult) {
	// query for each group
	groupsSet := make(map[int]bool)
	for _, groupName := range h.groups {
		if groupName != "" {
			group := p.queryGroup(ctx, qb, groupName)
			if group != nil {
				if _, found := groupsSet[group.ID]; !found {
					result.Movies = append(result.Movies, &models.SceneMovieID{
						MovieID: strconv.Itoa(group.ID),
					})
					groupsSet[group.ID] = true
				}
			}
		}
	}
}

func (p *FilenameParser) setParserResult(ctx context.Context, h sceneHolder, result *models.SceneParserResult) {
	if h.result.Title != "" {
		title := h.result.Title
		title = p.replaceWhitespaceCharacters(title)

		if p.ParserInput.CapitalizeTitle != nil && *p.ParserInput.CapitalizeTitle {
			title = capitalizeTitleRE.ReplaceAllStringFunc(title, strings.ToUpper)
		}

		result.Title = &title
	}

	if h.result.Date != nil {
		dateStr := h.result.Date.String()
		result.Date = &dateStr
	}

	if h.result.Rating != nil {
		result.Rating = h.result.Rating
	}

	r := p.repository

	if len(h.performers) > 0 {
		p.setPerformers(ctx, r.Performer, h, result)
	}
	if len(h.tags) > 0 {
		p.setTags(ctx, r.Tag, h, result)
	}
	p.setStudio(ctx, r.Studio, h, result)

	if len(h.groups) > 0 {
		p.setGroups(ctx, r.Group, h, result)
	}
}
