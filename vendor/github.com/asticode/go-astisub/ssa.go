package astisub

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/asticode/go-astikit"
)

// https://www.matroska.org/technical/specs/subtitles/ssa.html
// http://moodub.free.fr/video/ass-specs.doc
// https://en.wikipedia.org/wiki/SubStation_Alpha

// SSA alignment
const (
	ssaAlignmentCentered              = 2
	ssaAlignmentLeft                  = 1
	ssaAlignmentLeftJustifiedTopTitle = 5
	ssaAlignmentMidTitle              = 8
	ssaAlignmentRight                 = 3
	ssaAlignmentTopTitle              = 4
)

// SSA border styles
const (
	ssaBorderStyleOpaqueBox            = 3
	ssaBorderStyleOutlineAndDropShadow = 1
)

// SSA collisions
const (
	ssaCollisionsNormal  = "Normal"
	ssaCollisionsReverse = "Reverse"
)

// SSA event categories
const (
	ssaEventCategoryCommand  = "Command"
	ssaEventCategoryComment  = "Comment"
	ssaEventCategoryDialogue = "Dialogue"
	ssaEventCategoryMovie    = "Movie"
	ssaEventCategoryPicture  = "Picture"
	ssaEventCategorySound    = "Sound"
)

// SSA event format names
const (
	ssaEventFormatNameEffect  = "Effect"
	ssaEventFormatNameEnd     = "End"
	ssaEventFormatNameLayer   = "Layer"
	ssaEventFormatNameMarginL = "MarginL"
	ssaEventFormatNameMarginR = "MarginR"
	ssaEventFormatNameMarginV = "MarginV"
	ssaEventFormatNameMarked  = "Marked"
	ssaEventFormatNameName    = "Name"
	ssaEventFormatNameStart   = "Start"
	ssaEventFormatNameStyle   = "Style"
	ssaEventFormatNameText    = "Text"
)

// SSA script info names
const (
	ssaScriptInfoNameCollisions          = "Collisions"
	ssaScriptInfoNameOriginalEditing     = "Original Editing"
	ssaScriptInfoNameOriginalScript      = "Original Script"
	ssaScriptInfoNameOriginalTiming      = "Original Timing"
	ssaScriptInfoNameOriginalTranslation = "Original Translation"
	ssaScriptInfoNamePlayDepth           = "PlayDepth"
	ssaScriptInfoNamePlayResX            = "PlayResX"
	ssaScriptInfoNamePlayResY            = "PlayResY"
	ssaScriptInfoNameScriptType          = "ScriptType"
	ssaScriptInfoNameScriptUpdatedBy     = "Script Updated By"
	ssaScriptInfoNameSynchPoint          = "Synch Point"
	ssaScriptInfoNameTimer               = "Timer"
	ssaScriptInfoNameTitle               = "Title"
	ssaScriptInfoNameUpdateDetails       = "Update Details"
	ssaScriptInfoNameWrapStyle           = "WrapStyle"
)

// SSA section names
const (
	ssaSectionNameEvents     = "events"
	ssaSectionNameScriptInfo = "script.info"
	ssaSectionNameStyles     = "styles"
	ssaSectionNameUnknown    = "unknown"
)

// SSA style format names
const (
	ssaStyleFormatNameAlignment       = "Alignment"
	ssaStyleFormatNameAlphaLevel      = "AlphaLevel"
	ssaStyleFormatNameAngle           = "Angle"
	ssaStyleFormatNameBackColour      = "BackColour"
	ssaStyleFormatNameBold            = "Bold"
	ssaStyleFormatNameBorderStyle     = "BorderStyle"
	ssaStyleFormatNameEncoding        = "Encoding"
	ssaStyleFormatNameFontName        = "Fontname"
	ssaStyleFormatNameFontSize        = "Fontsize"
	ssaStyleFormatNameItalic          = "Italic"
	ssaStyleFormatNameMarginL         = "MarginL"
	ssaStyleFormatNameMarginR         = "MarginR"
	ssaStyleFormatNameMarginV         = "MarginV"
	ssaStyleFormatNameName            = "Name"
	ssaStyleFormatNameOutline         = "Outline"
	ssaStyleFormatNameOutlineColour   = "OutlineColour"
	ssaStyleFormatNamePrimaryColour   = "PrimaryColour"
	ssaStyleFormatNameScaleX          = "ScaleX"
	ssaStyleFormatNameScaleY          = "ScaleY"
	ssaStyleFormatNameSecondaryColour = "SecondaryColour"
	ssaStyleFormatNameShadow          = "Shadow"
	ssaStyleFormatNameSpacing         = "Spacing"
	ssaStyleFormatNameStrikeout       = "Strikeout"
	ssaStyleFormatNameTertiaryColour  = "TertiaryColour"
	ssaStyleFormatNameUnderline       = "Underline"
)

// SSA wrap style
const (
	ssaWrapStyleEndOfLineWordWrapping                   = "1"
	ssaWrapStyleNoWordWrapping                          = "2"
	ssaWrapStyleSmartWrapping                           = "0"
	ssaWrapStyleSmartWrappingWithLowerLinesGettingWider = "3"
)

// SSA regexp
var ssaRegexpEffect = regexp.MustCompile(`\{[^\{]+\}`)

// ReadFromSSA parses an .ssa content
func ReadFromSSA(i io.Reader) (o *Subtitles, err error) {
	o, err = ReadFromSSAWithOptions(i, defaultSSAOptions())
	return o, err
}

// ReadFromSSAWithOptions parses an .ssa content
func ReadFromSSAWithOptions(i io.Reader, opts SSAOptions) (o *Subtitles, err error) {
	// Init
	o = NewSubtitles()
	var scanner = bufio.NewScanner(i)
	var si = &ssaScriptInfo{}
	var ss = []*ssaStyle{}
	var es = []*ssaEvent{}

	// Scan
	var line, sectionName string
	var format map[int]string
	isFirstLine := true
	for scanner.Scan() {
		// Fetch line
		line = strings.TrimSpace(scanner.Text())

		// Remove BOM header
		if isFirstLine {
			line = strings.TrimPrefix(line, string(BytesBOM))
			isFirstLine = false
		}

		// Empty line
		if len(line) == 0 {
			continue
		}

		// Section name
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			switch strings.ToLower(line[1 : len(line)-1]) {
			case "events":
				sectionName = ssaSectionNameEvents
				format = make(map[int]string)
				continue
			case "script info":
				sectionName = ssaSectionNameScriptInfo
				continue
			case "v4 styles", "v4+ styles", "v4 styles+":
				sectionName = ssaSectionNameStyles
				format = make(map[int]string)
				continue
			default:
				if opts.OnUnknownSectionName != nil {
					opts.OnUnknownSectionName(line)
				}
				sectionName = ssaSectionNameUnknown
				continue
			}
		}

		// Unknown section
		if sectionName == ssaSectionNameUnknown {
			continue
		}

		// Comment
		if len(line) > 0 && line[0] == ';' {
			si.comments = append(si.comments, strings.TrimSpace(line[1:]))
			continue
		}

		// Split on ":"
		var split = strings.Split(line, ":")
		if len(split) < 2 || split[0] == "" {
			if opts.OnInvalidLine != nil {
				opts.OnInvalidLine(line)
			}
			continue
		}
		var header = strings.TrimSpace(split[0])
		var content = strings.TrimSpace(strings.Join(split[1:], ":"))

		// Switch on section name
		switch sectionName {
		case ssaSectionNameScriptInfo:
			if err = si.parse(header, content); err != nil {
				err = fmt.Errorf("astisub: parsing script info block failed: %w", err)
				return
			}
		case ssaSectionNameEvents, ssaSectionNameStyles:
			// Parse format
			if header == "Format" {
				for idx, item := range strings.Split(content, ",") {
					format[idx] = strings.TrimSpace(item)
				}
			} else {
				// No format provided
				if len(format) == 0 {
					err = fmt.Errorf("astisub: no %s format provided", sectionName)
					return
				}

				// Switch on section name
				switch sectionName {
				case ssaSectionNameEvents:
					var e *ssaEvent
					if e, err = newSSAEventFromString(header, content, format); err != nil {
						err = fmt.Errorf("astisub: building new ssa event failed: %w", err)
						return
					}
					es = append(es, e)
				case ssaSectionNameStyles:
					var s *ssaStyle
					if s, err = newSSAStyleFromString(content, format); err != nil {
						err = fmt.Errorf("astisub: building new ssa style failed: %w", err)
						return
					}
					ss = append(ss, s)
				}
			}
		}
	}

	// Set metadata
	o.Metadata = si.metadata()

	// Loop through styles
	for _, s := range ss {
		var st = s.style()
		o.Styles[st.ID] = st
	}

	// Loop through events
	for _, e := range es {
		// Only process dialogues
		if e.category == ssaEventCategoryDialogue {
			// Build item
			var item *Item
			if item, err = e.item(o.Styles); err != nil {
				return
			}

			// Append item
			o.Items = append(o.Items, item)
		}
	}
	return
}

// newColorFromSSAColor builds a new color based on an SSA color
func newColorFromSSAColor(i string) (_ *Color, _ error) {
	// Empty
	if len(i) == 0 {
		return
	}

	// Check whether input is decimal or hexadecimal
	var s = i
	var base = 10
	if strings.HasPrefix(i, "&H") {
		s = i[2:]
		base = 16
	}
	return newColorFromSSAString(s, base)
}

// newSSAColorFromColor builds a new SSA color based on a color
func newSSAColorFromColor(i *Color) string {
	return "&H" + i.SSAString()
}

// ssaScriptInfo represents an SSA script info block
type ssaScriptInfo struct {
	collisions          string
	comments            []string
	originalEditing     string
	originalScript      string
	originalTiming      string
	originalTranslation string
	playDepth           *int
	playResX, playResY  *int
	scriptType          string
	scriptUpdatedBy     string
	synchPoint          string
	timer               *float64
	title               string
	updateDetails       string
	wrapStyle           string
}

// newSSAScriptInfo builds an SSA script info block based on metadata
func newSSAScriptInfo(m *Metadata) (o *ssaScriptInfo) {
	// Init
	o = &ssaScriptInfo{}

	// Add metadata
	if m != nil {
		o.collisions = m.SSACollisions
		o.comments = m.Comments
		o.originalEditing = m.SSAOriginalEditing
		o.originalScript = m.SSAOriginalScript
		o.originalTiming = m.SSAOriginalTiming
		o.originalTranslation = m.SSAOriginalTranslation
		o.playDepth = m.SSAPlayDepth
		o.playResX = m.SSAPlayResX
		o.playResY = m.SSAPlayResY
		o.scriptType = m.SSAScriptType
		o.scriptUpdatedBy = m.SSAScriptUpdatedBy
		o.synchPoint = m.SSASynchPoint
		o.timer = m.SSATimer
		o.title = m.Title
		o.updateDetails = m.SSAUpdateDetails
		o.wrapStyle = m.SSAWrapStyle
	}
	return
}

// parse parses a script info header/content
func (b *ssaScriptInfo) parse(header, content string) (err error) {
	switch header {
	case ssaScriptInfoNameCollisions:
		b.collisions = content
	case ssaScriptInfoNameOriginalEditing:
		b.originalEditing = content
	case ssaScriptInfoNameOriginalScript:
		b.originalScript = content
	case ssaScriptInfoNameOriginalTiming:
		b.originalTiming = content
	case ssaScriptInfoNameOriginalTranslation:
		b.originalTranslation = content
	case ssaScriptInfoNameScriptType:
		b.scriptType = content
	case ssaScriptInfoNameScriptUpdatedBy:
		b.scriptUpdatedBy = content
	case ssaScriptInfoNameSynchPoint:
		b.synchPoint = content
	case ssaScriptInfoNameTitle:
		b.title = content
	case ssaScriptInfoNameUpdateDetails:
		b.updateDetails = content
	case ssaScriptInfoNameWrapStyle:
		b.wrapStyle = content
	// Int
	case ssaScriptInfoNamePlayResX, ssaScriptInfoNamePlayResY, ssaScriptInfoNamePlayDepth:
		var v int
		if v, err = strconv.Atoi(content); err != nil {
			err = fmt.Errorf("astisub: atoi of %s failed: %w", content, err)
		}
		switch header {
		case ssaScriptInfoNamePlayDepth:
			b.playDepth = astikit.IntPtr(v)
		case ssaScriptInfoNamePlayResX:
			b.playResX = astikit.IntPtr(v)
		case ssaScriptInfoNamePlayResY:
			b.playResY = astikit.IntPtr(v)
		}
	// Float
	case ssaScriptInfoNameTimer:
		var v float64
		if v, err = strconv.ParseFloat(strings.Replace(content, ",", ".", -1), 64); err != nil {
			err = fmt.Errorf("astisub: parseFloat of %s failed: %w", content, err)
		}
		b.timer = astikit.Float64Ptr(v)
	}
	return
}

// metadata returns the block as Metadata
func (b *ssaScriptInfo) metadata() *Metadata {
	return &Metadata{
		Comments:               b.comments,
		SSACollisions:          b.collisions,
		SSAOriginalEditing:     b.originalEditing,
		SSAOriginalScript:      b.originalScript,
		SSAOriginalTiming:      b.originalTiming,
		SSAOriginalTranslation: b.originalTranslation,
		SSAPlayDepth:           b.playDepth,
		SSAPlayResX:            b.playResX,
		SSAPlayResY:            b.playResY,
		SSAScriptType:          b.scriptType,
		SSAScriptUpdatedBy:     b.scriptUpdatedBy,
		SSASynchPoint:          b.synchPoint,
		SSATimer:               b.timer,
		SSAUpdateDetails:       b.updateDetails,
		SSAWrapStyle:           b.wrapStyle,
		Title:                  b.title,
	}
}

// bytes returns the block as bytes
func (b *ssaScriptInfo) bytes() (o []byte) {
	o = []byte("[Script Info]")
	o = append(o, bytesLineSeparator...)
	for _, c := range b.comments {
		o = appendStringToBytesWithNewLine(o, "; "+c)
	}
	if len(b.collisions) > 0 {
		o = appendStringToBytesWithNewLine(o, ssaScriptInfoNameCollisions+": "+b.collisions)
	}
	if len(b.originalEditing) > 0 {
		o = appendStringToBytesWithNewLine(o, ssaScriptInfoNameOriginalEditing+": "+b.originalEditing)
	}
	if len(b.originalScript) > 0 {
		o = appendStringToBytesWithNewLine(o, ssaScriptInfoNameOriginalScript+": "+b.originalScript)
	}
	if len(b.originalTiming) > 0 {
		o = appendStringToBytesWithNewLine(o, ssaScriptInfoNameOriginalTiming+": "+b.originalTiming)
	}
	if len(b.originalTranslation) > 0 {
		o = appendStringToBytesWithNewLine(o, ssaScriptInfoNameOriginalTranslation+": "+b.originalTranslation)
	}
	if b.playDepth != nil {
		o = appendStringToBytesWithNewLine(o, ssaScriptInfoNamePlayDepth+": "+strconv.Itoa(*b.playDepth))
	}
	if b.playResX != nil {
		o = appendStringToBytesWithNewLine(o, ssaScriptInfoNamePlayResX+": "+strconv.Itoa(*b.playResX))
	}
	if b.playResY != nil {
		o = appendStringToBytesWithNewLine(o, ssaScriptInfoNamePlayResY+": "+strconv.Itoa(*b.playResY))
	}
	if len(b.scriptType) > 0 {
		o = appendStringToBytesWithNewLine(o, ssaScriptInfoNameScriptType+": "+b.scriptType)
	}
	if len(b.scriptUpdatedBy) > 0 {
		o = appendStringToBytesWithNewLine(o, ssaScriptInfoNameScriptUpdatedBy+": "+b.scriptUpdatedBy)
	}
	if len(b.synchPoint) > 0 {
		o = appendStringToBytesWithNewLine(o, ssaScriptInfoNameSynchPoint+": "+b.synchPoint)
	}
	if b.timer != nil {
		o = appendStringToBytesWithNewLine(o, ssaScriptInfoNameTimer+": "+strings.Replace(strconv.FormatFloat(*b.timer, 'f', -1, 64), ".", ",", -1))
	}
	if len(b.title) > 0 {
		o = appendStringToBytesWithNewLine(o, ssaScriptInfoNameTitle+": "+b.title)
	}
	if len(b.updateDetails) > 0 {
		o = appendStringToBytesWithNewLine(o, ssaScriptInfoNameUpdateDetails+": "+b.updateDetails)
	}
	if len(b.wrapStyle) > 0 {
		o = appendStringToBytesWithNewLine(o, ssaScriptInfoNameWrapStyle+": "+b.wrapStyle)
	}
	return
}

// ssaStyle represents an SSA style
type ssaStyle struct {
	alignment       *int
	alphaLevel      *float64
	angle           *float64 // degrees
	backColour      *Color
	bold            *bool
	borderStyle     *int
	encoding        *int
	fontName        string
	fontSize        *float64
	italic          *bool
	outline         *float64 // pixels
	outlineColour   *Color
	marginLeft      *int // pixels
	marginRight     *int // pixels
	marginVertical  *int // pixels
	name            string
	primaryColour   *Color
	scaleX          *float64 // %
	scaleY          *float64 // %
	secondaryColour *Color
	shadow          *float64 // pixels
	spacing         *float64 // pixels
	strikeout       *bool
	underline       *bool
}

// newSSAStyleFromStyle returns an SSA style based on a Style
func newSSAStyleFromStyle(i Style) *ssaStyle {
	return &ssaStyle{
		alignment:       i.InlineStyle.SSAAlignment,
		alphaLevel:      i.InlineStyle.SSAAlphaLevel,
		angle:           i.InlineStyle.SSAAngle,
		backColour:      i.InlineStyle.SSABackColour,
		bold:            i.InlineStyle.SSABold,
		borderStyle:     i.InlineStyle.SSABorderStyle,
		encoding:        i.InlineStyle.SSAEncoding,
		fontName:        i.InlineStyle.SSAFontName,
		fontSize:        i.InlineStyle.SSAFontSize,
		italic:          i.InlineStyle.SSAItalic,
		outline:         i.InlineStyle.SSAOutline,
		outlineColour:   i.InlineStyle.SSAOutlineColour,
		marginLeft:      i.InlineStyle.SSAMarginLeft,
		marginRight:     i.InlineStyle.SSAMarginRight,
		marginVertical:  i.InlineStyle.SSAMarginVertical,
		name:            i.ID,
		primaryColour:   i.InlineStyle.SSAPrimaryColour,
		scaleX:          i.InlineStyle.SSAScaleX,
		scaleY:          i.InlineStyle.SSAScaleY,
		secondaryColour: i.InlineStyle.SSASecondaryColour,
		shadow:          i.InlineStyle.SSAShadow,
		spacing:         i.InlineStyle.SSASpacing,
		strikeout:       i.InlineStyle.SSAStrikeout,
		underline:       i.InlineStyle.SSAUnderline,
	}
}

// newSSAStyleFromString returns an SSA style based on an input string and a format
func newSSAStyleFromString(content string, format map[int]string) (s *ssaStyle, err error) {
	// Split content
	var items = strings.Split(content, ",")

	// Not enough items
	if len(items) < len(format) {
		err = fmt.Errorf("astisub: content has %d items whereas style format has %d items", len(items), len(format))
		return
	}

	// Loop through items
	s = &ssaStyle{}
	for idx, item := range items {
		// Index not found in format
		var attr string
		var ok bool
		if attr, ok = format[idx]; !ok {
			err = fmt.Errorf("astisub: index %d not found in style format %+v", idx, format)
			return
		}

		// Switch on attribute name
		switch attr {
		// Bool
		case ssaStyleFormatNameBold, ssaStyleFormatNameItalic, ssaStyleFormatNameStrikeout,
			ssaStyleFormatNameUnderline:
			var b = item == "-1"
			switch attr {
			case ssaStyleFormatNameBold:
				s.bold = astikit.BoolPtr(b)
			case ssaStyleFormatNameItalic:
				s.italic = astikit.BoolPtr(b)
			case ssaStyleFormatNameStrikeout:
				s.strikeout = astikit.BoolPtr(b)
			case ssaStyleFormatNameUnderline:
				s.underline = astikit.BoolPtr(b)
			}
		// Color
		case ssaStyleFormatNamePrimaryColour, ssaStyleFormatNameSecondaryColour,
			ssaStyleFormatNameTertiaryColour, ssaStyleFormatNameOutlineColour, ssaStyleFormatNameBackColour:
			// Build color
			var c *Color
			if c, err = newColorFromSSAColor(item); err != nil {
				err = fmt.Errorf("astisub: building new %s from ssa color %s failed: %w", attr, item, err)
				return
			}

			// Set color
			switch attr {
			case ssaStyleFormatNameBackColour:
				s.backColour = c
			case ssaStyleFormatNamePrimaryColour:
				s.primaryColour = c
			case ssaStyleFormatNameSecondaryColour:
				s.secondaryColour = c
			case ssaStyleFormatNameTertiaryColour, ssaStyleFormatNameOutlineColour:
				s.outlineColour = c
			}
		// Float
		case ssaStyleFormatNameAlphaLevel, ssaStyleFormatNameAngle, ssaStyleFormatNameFontSize,
			ssaStyleFormatNameScaleX, ssaStyleFormatNameScaleY,
			ssaStyleFormatNameOutline, ssaStyleFormatNameShadow, ssaStyleFormatNameSpacing:
			// Parse float
			var f float64
			if f, err = strconv.ParseFloat(item, 64); err != nil {
				err = fmt.Errorf("astisub: parsing float %s failed: %w", item, err)
				return
			}

			// Set float
			switch attr {
			case ssaStyleFormatNameAlphaLevel:
				s.alphaLevel = astikit.Float64Ptr(f)
			case ssaStyleFormatNameAngle:
				s.angle = astikit.Float64Ptr(f)
			case ssaStyleFormatNameFontSize:
				s.fontSize = astikit.Float64Ptr(f)
			case ssaStyleFormatNameScaleX:
				s.scaleX = astikit.Float64Ptr(f)
			case ssaStyleFormatNameScaleY:
				s.scaleY = astikit.Float64Ptr(f)
			case ssaStyleFormatNameOutline:
				s.outline = astikit.Float64Ptr(f)
			case ssaStyleFormatNameShadow:
				s.shadow = astikit.Float64Ptr(f)
			case ssaStyleFormatNameSpacing:
				s.spacing = astikit.Float64Ptr(f)
			}
		// Int
		case ssaStyleFormatNameAlignment, ssaStyleFormatNameBorderStyle, ssaStyleFormatNameEncoding,
			ssaStyleFormatNameMarginL, ssaStyleFormatNameMarginR, ssaStyleFormatNameMarginV:
			// Parse int
			var i int
			if i, err = strconv.Atoi(item); err != nil {
				err = fmt.Errorf("astisub: atoi of %s failed: %w", item, err)
				return
			}

			// Set int
			switch attr {
			case ssaStyleFormatNameAlignment:
				s.alignment = astikit.IntPtr(i)
			case ssaStyleFormatNameBorderStyle:
				s.borderStyle = astikit.IntPtr(i)
			case ssaStyleFormatNameEncoding:
				s.encoding = astikit.IntPtr(i)
			case ssaStyleFormatNameMarginL:
				s.marginLeft = astikit.IntPtr(i)
			case ssaStyleFormatNameMarginR:
				s.marginRight = astikit.IntPtr(i)
			case ssaStyleFormatNameMarginV:
				s.marginVertical = astikit.IntPtr(i)
			}
		// String
		case ssaStyleFormatNameFontName, ssaStyleFormatNameName:
			switch attr {
			case ssaStyleFormatNameFontName:
				s.fontName = item
			case ssaStyleFormatNameName:
				s.name = item
			}
		}
	}
	return
}

// ssaUpdateFormat updates an SSA format
func ssaUpdateFormat(n string, formatMap map[string]bool, format []string) []string {
	if _, ok := formatMap[n]; !ok {
		formatMap[n] = true
		format = append(format, n)
	}
	return format
}

// updateFormat updates the format based on the non empty fields
func (s ssaStyle) updateFormat(formatMap map[string]bool, format []string) []string {
	if s.alignment != nil {
		format = ssaUpdateFormat(ssaStyleFormatNameAlignment, formatMap, format)
	}
	if s.alphaLevel != nil {
		format = ssaUpdateFormat(ssaStyleFormatNameAlphaLevel, formatMap, format)
	}
	if s.angle != nil {
		format = ssaUpdateFormat(ssaStyleFormatNameAngle, formatMap, format)
	}
	if s.backColour != nil {
		format = ssaUpdateFormat(ssaStyleFormatNameBackColour, formatMap, format)
	}
	if s.bold != nil {
		format = ssaUpdateFormat(ssaStyleFormatNameBold, formatMap, format)
	}
	if s.borderStyle != nil {
		format = ssaUpdateFormat(ssaStyleFormatNameBorderStyle, formatMap, format)
	}
	if s.encoding != nil {
		format = ssaUpdateFormat(ssaStyleFormatNameEncoding, formatMap, format)
	}
	if len(s.fontName) > 0 {
		format = ssaUpdateFormat(ssaStyleFormatNameFontName, formatMap, format)
	}
	if s.fontSize != nil {
		format = ssaUpdateFormat(ssaStyleFormatNameFontSize, formatMap, format)
	}
	if s.italic != nil {
		format = ssaUpdateFormat(ssaStyleFormatNameItalic, formatMap, format)
	}
	if s.marginLeft != nil {
		format = ssaUpdateFormat(ssaStyleFormatNameMarginL, formatMap, format)
	}
	if s.marginRight != nil {
		format = ssaUpdateFormat(ssaStyleFormatNameMarginR, formatMap, format)
	}
	if s.marginVertical != nil {
		format = ssaUpdateFormat(ssaStyleFormatNameMarginV, formatMap, format)
	}
	if s.outline != nil {
		format = ssaUpdateFormat(ssaStyleFormatNameOutline, formatMap, format)
	}
	if s.outlineColour != nil {
		format = ssaUpdateFormat(ssaStyleFormatNameOutlineColour, formatMap, format)
	}
	if s.primaryColour != nil {
		format = ssaUpdateFormat(ssaStyleFormatNamePrimaryColour, formatMap, format)
	}
	if s.scaleX != nil {
		format = ssaUpdateFormat(ssaStyleFormatNameScaleX, formatMap, format)
	}
	if s.scaleY != nil {
		format = ssaUpdateFormat(ssaStyleFormatNameScaleY, formatMap, format)
	}
	if s.secondaryColour != nil {
		format = ssaUpdateFormat(ssaStyleFormatNameSecondaryColour, formatMap, format)
	}
	if s.shadow != nil {
		format = ssaUpdateFormat(ssaStyleFormatNameShadow, formatMap, format)
	}
	if s.spacing != nil {
		format = ssaUpdateFormat(ssaStyleFormatNameSpacing, formatMap, format)
	}
	if s.strikeout != nil {
		format = ssaUpdateFormat(ssaStyleFormatNameStrikeout, formatMap, format)
	}
	if s.underline != nil {
		format = ssaUpdateFormat(ssaStyleFormatNameUnderline, formatMap, format)
	}
	return format
}

// string returns the block as a string
func (s ssaStyle) string(format []string) string {
	var ss = []string{s.name}
	for _, attr := range format {
		var v string
		var found = true
		switch attr {
		// Bool
		case ssaStyleFormatNameBold, ssaStyleFormatNameItalic, ssaStyleFormatNameStrikeout,
			ssaStyleFormatNameUnderline:
			var b *bool
			switch attr {
			case ssaStyleFormatNameBold:
				b = s.bold
			case ssaStyleFormatNameItalic:
				b = s.italic
			case ssaStyleFormatNameStrikeout:
				b = s.strikeout
			case ssaStyleFormatNameUnderline:
				b = s.underline
			}
			if b != nil {
				v = "0"
				if *b {
					v = "1"
				}
			}
		// Color
		case ssaStyleFormatNamePrimaryColour, ssaStyleFormatNameSecondaryColour,
			ssaStyleFormatNameOutlineColour, ssaStyleFormatNameBackColour:
			var c *Color
			switch attr {
			case ssaStyleFormatNameBackColour:
				c = s.backColour
			case ssaStyleFormatNamePrimaryColour:
				c = s.primaryColour
			case ssaStyleFormatNameSecondaryColour:
				c = s.secondaryColour
			case ssaStyleFormatNameOutlineColour:
				c = s.outlineColour
			}
			if c != nil {
				v = newSSAColorFromColor(c)
			}
		// Float
		case ssaStyleFormatNameAlphaLevel, ssaStyleFormatNameAngle, ssaStyleFormatNameFontSize,
			ssaStyleFormatNameScaleX, ssaStyleFormatNameScaleY,
			ssaStyleFormatNameOutline, ssaStyleFormatNameShadow, ssaStyleFormatNameSpacing:
			var f *float64
			switch attr {
			case ssaStyleFormatNameAlphaLevel:
				f = s.alphaLevel
			case ssaStyleFormatNameAngle:
				f = s.angle
			case ssaStyleFormatNameFontSize:
				f = s.fontSize
			case ssaStyleFormatNameScaleX:
				f = s.scaleX
			case ssaStyleFormatNameScaleY:
				f = s.scaleY
			case ssaStyleFormatNameOutline:
				f = s.outline
			case ssaStyleFormatNameShadow:
				f = s.shadow
			case ssaStyleFormatNameSpacing:
				f = s.spacing
			}
			if f != nil {
				v = strconv.FormatFloat(*f, 'f', 3, 64)
			}
		// Int
		case ssaStyleFormatNameAlignment, ssaStyleFormatNameBorderStyle, ssaStyleFormatNameEncoding,
			ssaStyleFormatNameMarginL, ssaStyleFormatNameMarginR, ssaStyleFormatNameMarginV:
			var i *int
			switch attr {
			case ssaStyleFormatNameAlignment:
				i = s.alignment
			case ssaStyleFormatNameBorderStyle:
				i = s.borderStyle
			case ssaStyleFormatNameEncoding:
				i = s.encoding
			case ssaStyleFormatNameMarginL:
				i = s.marginLeft
			case ssaStyleFormatNameMarginR:
				i = s.marginRight
			case ssaStyleFormatNameMarginV:
				i = s.marginVertical
			}
			if i != nil {
				v = strconv.Itoa(*i)
			}
		// String
		case ssaStyleFormatNameFontName:
			switch attr {
			case ssaStyleFormatNameFontName:
				v = s.fontName
			}
		default:
			found = false
		}
		if found {
			ss = append(ss, v)
		}
	}
	return strings.Join(ss, ",")
}

// style converts ssaStyle to Style
func (s ssaStyle) style() (o *Style) {
	o = &Style{
		ID: s.name,
		InlineStyle: &StyleAttributes{
			SSAAlignment:       s.alignment,
			SSAAlphaLevel:      s.alphaLevel,
			SSAAngle:           s.angle,
			SSABackColour:      s.backColour,
			SSABold:            s.bold,
			SSABorderStyle:     s.borderStyle,
			SSAEncoding:        s.encoding,
			SSAFontName:        s.fontName,
			SSAFontSize:        s.fontSize,
			SSAItalic:          s.italic,
			SSAOutline:         s.outline,
			SSAOutlineColour:   s.outlineColour,
			SSAMarginLeft:      s.marginLeft,
			SSAMarginRight:     s.marginRight,
			SSAMarginVertical:  s.marginVertical,
			SSAPrimaryColour:   s.primaryColour,
			SSAScaleX:          s.scaleX,
			SSAScaleY:          s.scaleY,
			SSASecondaryColour: s.secondaryColour,
			SSAShadow:          s.shadow,
			SSASpacing:         s.spacing,
			SSAStrikeout:       s.strikeout,
			SSAUnderline:       s.underline,
		},
	}
	o.InlineStyle.propagateSSAAttributes()
	return
}

// ssaEvent represents an SSA event
type ssaEvent struct {
	category       string
	effect         string
	end            time.Duration
	layer          *int
	marked         *bool
	marginLeft     *int // pixels
	marginRight    *int // pixels
	marginVertical *int // pixels
	name           string
	start          time.Duration
	style          string
	text           string
}

// newSSAEventFromItem returns an SSA Event based on an input item
func newSSAEventFromItem(i Item) (e *ssaEvent) {
	// Init
	e = &ssaEvent{
		category: ssaEventCategoryDialogue,
		end:      i.EndAt,
		start:    i.StartAt,
	}

	// Style
	if i.Style != nil {
		e.style = i.Style.ID
	}

	// Inline style
	if i.InlineStyle != nil {
		e.effect = i.InlineStyle.SSAEffect
		e.layer = i.InlineStyle.SSALayer
		e.marginLeft = i.InlineStyle.SSAMarginLeft
		e.marginRight = i.InlineStyle.SSAMarginRight
		e.marginVertical = i.InlineStyle.SSAMarginVertical
		e.marked = i.InlineStyle.SSAMarked
	}

	// Text
	var lines []string
	for _, l := range i.Lines {
		var items []string
		for _, item := range l.Items {
			var s string
			if item.InlineStyle != nil && len(item.InlineStyle.SSAEffect) > 0 {
				s += item.InlineStyle.SSAEffect
			}
			s += item.Text
			items = append(items, s)
		}
		if len(l.VoiceName) > 0 {
			e.name = l.VoiceName
		}
		lines = append(lines, strings.Join(items, " "))
	}
	e.text = strings.Join(lines, "\\n")
	return
}

// newSSAEventFromString returns an SSA event based on an input string and a format
func newSSAEventFromString(header, content string, format map[int]string) (e *ssaEvent, err error) {
	// Split content
	var items = strings.Split(content, ",")

	// Not enough items
	if len(items) < len(format) {
		err = fmt.Errorf("astisub: content has %d items whereas style format has %d items", len(items), len(format))
		return
	}

	// Last item may contain commas, therefore we need to fix it
	items[len(format)-1] = strings.Join(items[len(format)-1:], ",")
	items = items[:len(format)]

	// Loop through items
	e = &ssaEvent{category: header}
	for idx, item := range items {
		// Index not found in format
		var attr string
		var ok bool
		if attr, ok = format[idx]; !ok {
			err = fmt.Errorf("astisub: index %d not found in event format %+v", idx, format)
			return
		}

		// Switch on attribute name
		switch attr {
		// Duration
		case ssaEventFormatNameStart, ssaEventFormatNameEnd:
			// Parse duration
			var d time.Duration
			if d, err = parseDurationSSA(item); err != nil {
				err = fmt.Errorf("astisub: parsing ssa duration %s failed: %w", item, err)
				return
			}

			// Set duration
			switch attr {
			case ssaEventFormatNameEnd:
				e.end = d
			case ssaEventFormatNameStart:
				e.start = d
			}
		// Int
		case ssaEventFormatNameLayer, ssaEventFormatNameMarginL, ssaEventFormatNameMarginR,
			ssaEventFormatNameMarginV:
			// Parse int
			var i int
			if i, err = strconv.Atoi(item); err != nil {
				err = fmt.Errorf("astisub: atoi of %s failed: %w", item, err)
				return
			}

			// Set int
			switch attr {
			case ssaEventFormatNameLayer:
				e.layer = astikit.IntPtr(i)
			case ssaEventFormatNameMarginL:
				e.marginLeft = astikit.IntPtr(i)
			case ssaEventFormatNameMarginR:
				e.marginRight = astikit.IntPtr(i)
			case ssaEventFormatNameMarginV:
				e.marginVertical = astikit.IntPtr(i)
			}
		// String
		case ssaEventFormatNameEffect, ssaEventFormatNameName, ssaEventFormatNameStyle, ssaEventFormatNameText:
			switch attr {
			case ssaEventFormatNameEffect:
				e.effect = item
			case ssaEventFormatNameName:
				e.name = item
			case ssaEventFormatNameStyle:
				// *Default is reserved
				// http://www.tcax.org/docs/ass-specs.htm
				if item == "*Default" {
					e.style = "Default"
				} else {
					e.style = item
				}
			case ssaEventFormatNameText:
				e.text = strings.TrimSpace(item)
			}
		// Marked
		case ssaEventFormatNameMarked:
			if item == "Marked=1" {
				e.marked = astikit.BoolPtr(true)
			} else {
				e.marked = astikit.BoolPtr(false)
			}
		}
	}
	return
}

// item converts an SSA event to an Item
func (e *ssaEvent) item(styles map[string]*Style) (i *Item, err error) {
	// Init item
	i = &Item{
		EndAt: e.end,
		InlineStyle: &StyleAttributes{
			SSAEffect:         e.effect,
			SSALayer:          e.layer,
			SSAMarginLeft:     e.marginLeft,
			SSAMarginRight:    e.marginRight,
			SSAMarginVertical: e.marginVertical,
			SSAMarked:         e.marked,
		},
		StartAt: e.start,
	}

	// Set style
	if len(e.style) > 0 {
		var ok bool
		if i.Style, ok = styles[e.style]; !ok {
			err = fmt.Errorf("astisub: style %s not found", e.style)
			return
		}
	}

	// Loop through lines
	for _, s := range strings.Split(e.text, "\\n") {
		// Init
		s = strings.TrimSpace(s)
		var l = Line{VoiceName: e.name}

		// Extract effects
		var matches = ssaRegexpEffect.FindAllStringIndex(s, -1)
		if len(matches) > 0 {
			// Loop through matches
			var lineItem *LineItem
			var previousEffectEndOffset int
			for _, idxs := range matches {
				if lineItem != nil {
					lineItem.Text = s[previousEffectEndOffset:idxs[0]]
					l.Items = append(l.Items, *lineItem)
				} else if idxs[0] > 0 {
					l.Items = append(l.Items, LineItem{Text: s[previousEffectEndOffset:idxs[0]]})
				}
				previousEffectEndOffset = idxs[1]
				lineItem = &LineItem{InlineStyle: &StyleAttributes{SSAEffect: s[idxs[0]:idxs[1]]}}
			}
			lineItem.Text = s[previousEffectEndOffset:]
			l.Items = append(l.Items, *lineItem)
		} else {
			l.Items = append(l.Items, LineItem{Text: s})
		}

		// Add line
		i.Lines = append(i.Lines, l)
	}
	return
}

// updateFormat updates the format based on the non empty fields
func (e ssaEvent) updateFormat(formatMap map[string]bool, format []string) []string {
	if len(e.effect) > 0 {
		format = ssaUpdateFormat(ssaEventFormatNameEffect, formatMap, format)
	}
	if e.layer != nil {
		format = ssaUpdateFormat(ssaEventFormatNameLayer, formatMap, format)
	}
	if e.marginLeft != nil {
		format = ssaUpdateFormat(ssaEventFormatNameMarginL, formatMap, format)
	}
	if e.marginRight != nil {
		format = ssaUpdateFormat(ssaEventFormatNameMarginR, formatMap, format)
	}
	if e.marginVertical != nil {
		format = ssaUpdateFormat(ssaEventFormatNameMarginV, formatMap, format)
	}
	if e.marked != nil {
		format = ssaUpdateFormat(ssaEventFormatNameMarked, formatMap, format)
	}
	if len(e.name) > 0 {
		format = ssaUpdateFormat(ssaEventFormatNameName, formatMap, format)
	}
	if len(e.style) > 0 {
		format = ssaUpdateFormat(ssaEventFormatNameStyle, formatMap, format)
	}
	return format
}

// formatDurationSSA formats an .ssa duration
func formatDurationSSA(i time.Duration) string {
	return formatDuration(i, ".", 2)
}

// string returns the block as a string
func (e *ssaEvent) string(format []string) string {
	var ss []string
	for _, attr := range format {
		var v string
		var found = true
		switch attr {
		// Duration
		case ssaEventFormatNameEnd, ssaEventFormatNameStart:
			switch attr {
			case ssaEventFormatNameEnd:
				v = formatDurationSSA(e.end)
			case ssaEventFormatNameStart:
				v = formatDurationSSA(e.start)
			}
		// Marked
		case ssaEventFormatNameMarked:
			if e.marked != nil {
				if *e.marked {
					v = "Marked=1"
				} else {
					v = "Marked=0"
				}
			}
		// Int
		case ssaEventFormatNameLayer, ssaEventFormatNameMarginL, ssaEventFormatNameMarginR,
			ssaEventFormatNameMarginV:
			var i *int
			switch attr {
			case ssaEventFormatNameLayer:
				i = e.layer
			case ssaEventFormatNameMarginL:
				i = e.marginLeft
			case ssaEventFormatNameMarginR:
				i = e.marginRight
			case ssaEventFormatNameMarginV:
				i = e.marginVertical
			}
			if i != nil {
				v = strconv.Itoa(*i)
			}
		// String
		case ssaEventFormatNameEffect, ssaEventFormatNameName, ssaEventFormatNameStyle, ssaEventFormatNameText:
			switch attr {
			case ssaEventFormatNameEffect:
				v = e.effect
			case ssaEventFormatNameName:
				v = e.name
			case ssaEventFormatNameStyle:
				v = e.style
			case ssaEventFormatNameText:
				v = e.text
			}
		default:
			found = false
		}
		if found {
			ss = append(ss, v)
		}
	}
	return strings.Join(ss, ",")
}

// parseDurationSSA parses an .ssa duration
func parseDurationSSA(i string) (time.Duration, error) {
	return parseDuration(i, ".", 3)
}

// WriteToSSA writes subtitles in .ssa format
func (s Subtitles) WriteToSSA(o io.Writer) (err error) {
	// Do not write anything if no subtitles
	if len(s.Items) == 0 {
		err = ErrNoSubtitlesToWrite
		return
	}

	// Write Script Info block
	var si = newSSAScriptInfo(s.Metadata)
	if _, err = o.Write(si.bytes()); err != nil {
		err = fmt.Errorf("astisub: writing script info block failed: %w", err)
		return
	}

	// Write Styles block
	if len(s.Styles) > 0 {
		// Header
		var b = []byte("\n[V4 Styles]\n")

		// Format
		var formatMap = make(map[string]bool)
		var format = []string{ssaStyleFormatNameName}
		var styles = make(map[string]*ssaStyle)
		var styleNames []string
		for _, s := range s.Styles {
			var ss = newSSAStyleFromStyle(*s)
			format = ss.updateFormat(formatMap, format)
			styles[ss.name] = ss
			styleNames = append(styleNames, ss.name)
		}
		b = append(b, []byte("Format: "+strings.Join(format, ", ")+"\n")...)

		// Styles
		sort.Strings(styleNames)
		for _, n := range styleNames {
			b = append(b, []byte("Style: "+styles[n].string(format)+"\n")...)
		}

		// Write
		if _, err = o.Write(b); err != nil {
			err = fmt.Errorf("astisub: writing styles block failed: %w", err)
			return
		}
	}

	// Write Events block
	if len(s.Items) > 0 {
		// Header
		var b = []byte("\n[Events]\n")

		// Format
		var formatMap = make(map[string]bool)
		var format = []string{
			ssaEventFormatNameStart,
			ssaEventFormatNameEnd,
		}
		var events []*ssaEvent
		for _, i := range s.Items {
			var e = newSSAEventFromItem(*i)
			format = e.updateFormat(formatMap, format)
			events = append(events, e)
		}
		format = append(format, ssaEventFormatNameText)
		b = append(b, []byte("Format: "+strings.Join(format, ", ")+"\n")...)

		// Styles
		for _, e := range events {
			b = append(b, []byte(ssaEventCategoryDialogue+": "+e.string(format)+"\n")...)
		}

		// Write
		if _, err = o.Write(b); err != nil {
			err = fmt.Errorf("astisub: writing events block failed: %w", err)
			return
		}
	}
	return
}

// SSAOptions
type SSAOptions struct {
	OnUnknownSectionName func(name string)
	OnInvalidLine        func(line string)
}

func defaultSSAOptions() SSAOptions {
	return SSAOptions{
		OnUnknownSectionName: func(name string) {
			log.Printf("astisub: unknown section: %s", name)
		},
		OnInvalidLine: func(line string) {
			log.Printf("astisub: not understood: '%s', ignoring", line)
		},
	}
}
