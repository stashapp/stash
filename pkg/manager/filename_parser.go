package manager

import (
	"database/sql"
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/stashapp/stash/pkg/models"
)

type parserField struct {
	field           string
	fieldRegex      *regexp.Regexp
	regex           string
	isFullDateField bool
}

func newParserField(field string, regex string, captured bool) parserField {
	ret := parserField{
		field:           field,
		isFullDateField: false,
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
	return string(f.fieldRegex.ReplaceAll([]byte(pattern), []byte(f.regex)))
}

func (f parserField) getFieldPattern() string {
	return "{" + f.field + "}"
}

var validFields map[string]parserField

func initParserFields() {
	if validFields != nil {
		return
	}

	ret := make(map[string]parserField)

	ret["title"] = newParserField("title", ".*", true)
	ret["ext"] = newParserField("ext", ".*$", false)

	//I = new ParserField("i", undefined, "Matches any ignored word", false);

	ret["d"] = newParserField("d", `(?:\.|-|_)`, false)
	ret["performer"] = newParserField("performer", ".*", true)
	ret["studio"] = newParserField("studio", ".*", true)
	ret["tag"] = newParserField("tag", ".*", true)

	// date fields
	ret["date"] = newParserField("date", `\d{4}-\d{2}-\d{2}`, true)
	ret["yyyy"] = newParserField("yyyy", `\d{4}`, true)
	ret["yy"] = newParserField("yy", `\d{2}`, true)
	ret["mm"] = newParserField("mm", `\d{2}`, true)
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
	fields []string
	regex  *regexp.Regexp
}

func getIgnoreClause(ignoreFields []string) string {
	if len(ignoreFields) == 0 {
		return ""
	}

	var ignoreClauses []string

	regex := regexp.MustCompile(`([\-\.\(\)\[\]])`)
	for _, v := range ignoreFields {
		newVal := string(regex.ReplaceAllString(v, "$$1"))
		newVal = strings.TrimSpace(newVal)
		newVal = "(?:" + newVal + ")"
		ignoreClauses = append(ignoreClauses, newVal)
	}

	return "(?:" + strings.Join(ignoreClauses, "|") + ")"
}

func newParseMapper(pattern string, ignoreFields []string) (*parseMapper, error) {
	ret := &parseMapper{}

	// escape control characters
	escapeCharRE := regexp.MustCompile(`([\-\.\(\)\[\]])`)
	regex := escapeCharRE.ReplaceAllString(pattern, `$$1`)

	// replace {} with wildcard
	braceRE := regexp.MustCompile(`\{\}`)
	regex = braceRE.ReplaceAllString(regex, ".*")

	// replace all known fields with applicable regexes
	regex = replacePatternWithRegex(regex, ignoreFields)

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
		fields = append(fields, field)
	}

	ret.fields = fields

	return ret, nil
}

type sceneHolder struct {
	scene  *models.Scene
	result *models.Scene
	yyyy   string
	mm     string
	dd     string
}

func newSceneHolder(scene *models.Scene) *sceneHolder {
	sceneCopy := models.Scene{
		ID:       scene.ID,
		Checksum: scene.Checksum,
		Path:     scene.Path,
	}
	ret := sceneHolder{
		scene:  scene,
		result: &sceneCopy,
	}

	return &ret
}

func validateDate(dateStr string) bool {
	splits := strings.Split(dateStr, "-")
	if len(splits) != 3 {
		return false
	}

	year, _ := strconv.Atoi(splits[0])
	month, _ := strconv.Atoi(splits[1])
	d, _ := strconv.Atoi(splits[2])

	// assume year must be between 1900 and 2100
	if year < 1900 || year > 2100 {
		return false
	}

	if month < 1 || month > 12 {
		return false
	}

	// not checking individual months to ensure date is in the correct range
	if d < 1 || d > 31 {
		return false
	}

	return true
}

func (h *sceneHolder) setDate(field *parserField, value string) {
	yearIndex := 0
	yearLength := len(strings.Split(field.field, "y")) - 1
	dateIndex := 0
	monthIndex := 0

	switch field.field {
	case "yyyymmdd":
	case "yymmdd":
		monthIndex = yearLength
		dateIndex = monthIndex + 2
		break
	case "ddmmyyyy":
	case "ddmmyy":
		monthIndex = 2
		yearIndex = monthIndex + 2
		break
	case "mmddyyyy":
	case "mmddyy":
		dateIndex = monthIndex + 2
		yearIndex = dateIndex + 2
		break
	}

	yearValue := value[yearIndex : yearIndex+yearLength-1]
	monthValue := value[monthIndex : monthIndex+1]
	dateValue := value[dateIndex : dateIndex+1]

	fullDate := yearValue + "-" + monthValue + "-" + dateValue

	// ensure the date is valid
	// only set if new value is different from the old
	if validateDate(fullDate) && h.scene.Date.String != fullDate {
		h.result.Date = models.SQLiteDate{
			String: fullDate,
			Valid:  true,
		}
	}
}

func (h *sceneHolder) setField(field parserField, value interface{}) {
	if field.isFullDateField {
		h.setDate(&field, value.(string))
		return
	}

	switch field.field {
	case "title":
		h.result.Title = sql.NullString{
			String: value.(string),
			Valid:  true,
		}
		break
	case "date":
		if validateDate(value.(string)) {
			h.result.Date = models.SQLiteDate{
				String: value.(string),
				Valid:  true,
			}
		}
		break
	case "performer":
		// TODO
		break
	case "studio":
		// TODO
		break
	case "tag":
		// TODO
		break
	case "yyyy":
		h.yyyy = value.(string)
		break
	case "yy":
		v := value.(string)
		v = "20" + v
		h.yyyy = v
		break
	case "mm":
		h.mm = value.(string)
		break
	case "dd":
		h.dd = value.(string)
		break
	}
	// TODO - other fields
}

func (h *sceneHolder) postParse() {
	// set the date if the components are set
	if h.yyyy != "" && h.mm != "" && h.dd != "" {
		fullDate := h.yyyy + "-" + h.mm + "-" + h.dd
		h.setField(validFields["date"], fullDate)
	}
}

func (m parseMapper) parse(scene *models.Scene) *models.SceneParserResult {
	result := m.regex.FindString(scene.Path)

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

	return &models.SceneParserResult{
		Scene:        scene,
		ParserResult: sceneHolder.result,
	}
}

type SceneFilenameParser struct {
	Pattern     string
	ParserInput models.SceneParserInput
	Filter      *models.FindFilterType
}

func (p *SceneFilenameParser) Parse() ([]*models.SceneParserResult, int, error) {
	// perform the query to find the scenes
	mapper, err := newParseMapper(p.Pattern, p.ParserInput.IgnoreWords)

	if err != nil {
		return nil, 0, err
	}

	p.Filter.Q = &p.Pattern

	qb := models.NewSceneQueryBuilder()
	scenes, total := qb.QueryByPathRegex(p.Filter)

	var ret []*models.SceneParserResult
	for _, scene := range scenes {
		r := mapper.parse(scene)

		if r != nil {
			ret = append(ret, r)
		}
	}

	return ret, total, nil
}
