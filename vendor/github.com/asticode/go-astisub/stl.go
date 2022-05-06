package astisub

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/asticode/go-astikit"
	"golang.org/x/text/unicode/norm"
)

// https://tech.ebu.ch/docs/tech/tech3264.pdf
// https://github.com/yanncoupin/stl2srt/blob/master/to_srt.py

// STL block sizes
const (
	stlBlockSizeGSI = 1024
	stlBlockSizeTTI = 128
)

// STL character code table number
const (
	stlCharacterCodeTableNumberLatin         uint16 = 12336
	stlCharacterCodeTableNumberLatinCyrillic uint16 = 12337
	stlCharacterCodeTableNumberLatinArabic   uint16 = 12338
	stlCharacterCodeTableNumberLatinGreek    uint16 = 12339
	stlCharacterCodeTableNumberLatinHebrew   uint16 = 12340
)

// STL character code tables
// TODO Add missing tables
var (
	stlCharacterCodeTables = map[uint16]*astikit.BiMap{
		stlCharacterCodeTableNumberLatin: astikit.NewBiMap().
			Set(0x20, " ").Set(0x21, "!").Set(0x22, "\"").Set(0x23, "#").
			Set(0x24, "¤").Set(0x25, "%").Set(0x26, "&").Set(0x27, "'").
			Set(0x28, "(").Set(0x29, ")").Set(0x2a, "*").Set(0x2b, "+").
			Set(0x2c, ",").Set(0x2d, "-").Set(0x2e, ".").Set(0x2f, "/").
			Set(0x30, "0").Set(0x31, "1").Set(0x32, "2").Set(0x33, "3").
			Set(0x34, "4").Set(0x35, "5").Set(0x36, "6").Set(0x37, "7").
			Set(0x38, "8").Set(0x39, "9").Set(0x3a, ":").Set(0x3b, ";").
			Set(0x3c, "<").Set(0x3d, "=").Set(0x3e, ">").Set(0x3f, "?").
			Set(0x40, "@").Set(0x41, "A").Set(0x42, "B").Set(0x43, "C").
			Set(0x44, "D").Set(0x45, "E").Set(0x46, "F").Set(0x47, "G").
			Set(0x48, "H").Set(0x49, "I").Set(0x4a, "J").Set(0x4b, "K").
			Set(0x4c, "L").Set(0x4d, "M").Set(0x4e, "N").Set(0x4f, "O").
			Set(0x50, "P").Set(0x51, "Q").Set(0x52, "R").Set(0x53, "S").
			Set(0x54, "T").Set(0x55, "U").Set(0x56, "V").Set(0x57, "W").
			Set(0x58, "X").Set(0x59, "Y").Set(0x5a, "Z").Set(0x5b, "[").
			Set(0x5c, "\\").Set(0x5d, "]").Set(0x5e, "^").Set(0x5f, "_").
			Set(0x60, "`").Set(0x61, "a").Set(0x62, "b").Set(0x63, "c").
			Set(0x64, "d").Set(0x65, "e").Set(0x66, "f").Set(0x67, "g").
			Set(0x68, "h").Set(0x69, "i").Set(0x6a, "j").Set(0x6b, "k").
			Set(0x6c, "l").Set(0x6d, "m").Set(0x6e, "n").Set(0x6f, "o").
			Set(0x70, "p").Set(0x71, "q").Set(0x72, "r").Set(0x73, "s").
			Set(0x74, "t").Set(0x75, "u").Set(0x76, "v").Set(0x77, "w").
			Set(0x78, "x").Set(0x79, "y").Set(0x7a, "z").Set(0x7b, "{").
			Set(0x7c, "|").Set(0x7d, "}").Set(0x7e, "~").
			Set(0xa0, string([]byte{0xC2, 0xA0})).Set(0xa1, "¡").Set(0xa2, "¢").
			Set(0xa3, "£").Set(0xa4, "$").Set(0xa5, "¥").Set(0xa7, "§").
			Set(0xa9, "‘").Set(0xaa, "“").Set(0xab, "«").Set(0xac, "←").
			Set(0xad, "↑").Set(0xae, "→").Set(0xaf, "↓").
			Set(0xb0, "°").Set(0xb1, "±").Set(0xb2, "²").Set(0xb3, "³").
			Set(0xb4, "×").Set(0xb5, "µ").Set(0xb6, "¶").Set(0xb7, "·").
			Set(0xb8, "÷").Set(0xb9, "’").Set(0xba, "”").Set(0xbb, "»").
			Set(0xbc, "¼").Set(0xbd, "½").Set(0xbe, "¾").Set(0xbf, "¿").
			Set(0xc1, string([]byte{0xCC, 0x80})).Set(0xc2, string([]byte{0xCC, 0x81})).
			Set(0xc3, string([]byte{0xCC, 0x82})).Set(0xc4, string([]byte{0xCC, 0x83})).
			Set(0xc5, string([]byte{0xCC, 0x84})).Set(0xc6, string([]byte{0xCC, 0x86})).
			Set(0xc7, string([]byte{0xCC, 0x87})).Set(0xc8, string([]byte{0xCC, 0x88})).
			Set(0xca, string([]byte{0xCC, 0x8A})).Set(0xcb, string([]byte{0xCC, 0xA7})).
			Set(0xcd, string([]byte{0xCC, 0x8B})).Set(0xce, string([]byte{0xCC, 0xA8})).
			Set(0xcf, string([]byte{0xCC, 0x8C})).
			Set(0xd0, "―").Set(0xd1, "¹").Set(0xd2, "®").Set(0xd3, "©").
			Set(0xd4, "™").Set(0xd5, "♪").Set(0xd6, "¬").Set(0xd7, "¦").
			Set(0xdc, "⅛").Set(0xdd, "⅜").Set(0xde, "⅝").Set(0xdf, "⅞").
			Set(0xe0, "Ω").Set(0xe1, "Æ").Set(0xe2, "Đ").Set(0xe3, "ª").
			Set(0xe4, "Ħ").Set(0xe6, "Ĳ").Set(0xe7, "Ŀ").Set(0xe8, "Ł").
			Set(0xe9, "Ø").Set(0xea, "Œ").Set(0xeb, "º").Set(0xec, "Þ").
			Set(0xed, "Ŧ").Set(0xee, "Ŋ").Set(0xef, "ŉ").
			Set(0xf0, "ĸ").Set(0xf1, "æ").Set(0xf2, "đ").Set(0xf3, "ð").
			Set(0xf4, "ħ").Set(0xf5, "ı").Set(0xf6, "ĳ").Set(0xf7, "ŀ").
			Set(0xf8, "ł").Set(0xf9, "ø").Set(0xfa, "œ").Set(0xfb, "ß").
			Set(0xfc, "þ").Set(0xfd, "ŧ").Set(0xfe, "ŋ").Set(0xff, string([]byte{0xC2, 0xAD})),
	}
)

// STL code page numbers
const (
	stlCodePageNumberCanadaFrench uint32 = 3683891
	stlCodePageNumberMultilingual uint32 = 3683632
	stlCodePageNumberNordic       uint32 = 3683893
	stlCodePageNumberPortugal     uint32 = 3683888
	stlCodePageNumberUnitedStates uint32 = 3420983
)

// STL comment flag
const (
	stlCommentFlagTextContainsSubtitleData                       = '\x00'
	stlCommentFlagTextContainsCommentsNotIntendedForTransmission = '\x01'
)

// STL country codes
const (
	stlCountryCodeChinese = "CHN"
	stlCountryCodeFrance  = "FRA"
	stlCountryCodeJapan   = "JPN"
	stlCountryCodeNorway  = "NOR"
)

// STL cumulative status
const (
	stlCumulativeStatusFirstSubtitleOfACumulativeSet        = '\x01'
	stlCumulativeStatusIntermediateSubtitleOfACumulativeSet = '\x02'
	stlCumulativeStatusLastSubtitleOfACumulativeSet         = '\x03'
	stlCumulativeStatusSubtitleNotPartOfACumulativeSet      = '\x00'
)

// STL display standard code
const (
	stlDisplayStandardCodeOpenSubtitling = "0"
	stlDisplayStandardCodeLevel1Teletext = "1"
	stlDisplayStandardCodeLevel2Teletext = "2"
)

// STL framerate mapping
var stlFramerateMapping = astikit.NewBiMap().
	Set("STL25.01", 25).
	Set("STL30.01", 30)

// STL justification code
const (
	stlJustificationCodeCentredText           = '\x02'
	stlJustificationCodeLeftJustifiedText     = '\x01'
	stlJustificationCodeRightJustifiedText    = '\x03'
	stlJustificationCodeUnchangedPresentation = '\x00'
)

// STL language codes
const (
	stlLanguageCodeChinese   = "75"
	stlLanguageCodeEnglish   = "09"
	stlLanguageCodeFrench    = "0F"
	stllanguageCodeJapanese  = "69"
	stlLanguageCodeNorwegian = "1E"
)

// STL language mapping
var stlLanguageMapping = astikit.NewBiMap().
	Set(stlLanguageCodeChinese, LanguageChinese).
	Set(stlLanguageCodeEnglish, LanguageEnglish).
	Set(stlLanguageCodeFrench, LanguageFrench).
	Set(stllanguageCodeJapanese, LanguageJapanese).
	Set(stlLanguageCodeNorwegian, LanguageNorwegian)

	// STL timecode status
const (
	stlTimecodeStatusNotIntendedForUse = "0"
	stlTimecodeStatusIntendedForUse    = "1"
)

// TTI Special Extension Block Number
const extensionBlockNumberReservedUserData = 0xfe

const stlLineSeparator = 0x8a

type STLPosition struct {
	VerticalPosition int
	MaxRows          int
	Rows             int
}

// STLOptions represents STL parsing options
type STLOptions struct {
	// IgnoreTimecodeStartOfProgramme - set STLTimecodeStartOfProgramme to zero before parsing
	IgnoreTimecodeStartOfProgramme bool
}

// ReadFromSTL parses an .stl content
func ReadFromSTL(i io.Reader, opts STLOptions) (o *Subtitles, err error) {
	// Init
	o = NewSubtitles()

	// Read GSI block
	var b []byte
	if b, err = readNBytes(i, stlBlockSizeGSI); err != nil {
		return
	}

	// Parse GSI block
	var g *gsiBlock
	if g, err = parseGSIBlock(b); err != nil {
		err = fmt.Errorf("astisub: building gsi block failed: %w", err)
		return
	}

	// Create character handler
	var ch *stlCharacterHandler
	if ch, err = newSTLCharacterHandler(g.characterCodeTableNumber); err != nil {
		err = fmt.Errorf("astisub: creating stl character handler failed: %w", err)
		return
	}

	// Update metadata
	// TODO Add more STL fields to metadata
	o.Metadata = &Metadata{
		Framerate:              g.framerate,
		STLCountryOfOrigin:     g.countryOfOrigin,
		STLCreationDate:        &g.creationDate,
		STLDisplayStandardCode: g.displayStandardCode,
		STLMaximumNumberOfDisplayableCharactersInAnyTextRow: astikit.IntPtr(g.maximumNumberOfDisplayableCharactersInAnyTextRow),
		STLMaximumNumberOfDisplayableRows:                   astikit.IntPtr(g.maximumNumberOfDisplayableRows),
		STLPublisher:                                        g.publisher,
		STLRevisionDate:                                     &g.revisionDate,
		STLSubtitleListReferenceCode:                        g.subtitleListReferenceCode,
		Title:                                               g.originalProgramTitle,
	}
	if !opts.IgnoreTimecodeStartOfProgramme {
		o.Metadata.STLTimecodeStartOfProgramme = g.timecodeStartOfProgramme
	}
	if v, ok := stlLanguageMapping.Get(g.languageCode); ok {
		o.Metadata.Language = v.(string)
	}

	// Parse Text and Timing Information (TTI) blocks.
	for {
		// Read TTI block
		if b, err = readNBytes(i, stlBlockSizeTTI); err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return
		}

		// Parse TTI block
		var t = parseTTIBlock(b, g.framerate)

		// Do not process reserved user data
		if t.extensionBlockNumber == extensionBlockNumberReservedUserData {
			continue
		}

		justification := parseSTLJustificationCode(t.justificationCode)
		rows := bytes.Split(t.text, []byte{stlLineSeparator})

		position := STLPosition{
			MaxRows:          g.maximumNumberOfDisplayableRows,
			Rows:             len(rows),
			VerticalPosition: t.verticalPosition,
		}

		styleAttributes := StyleAttributes{
			STLJustification: &justification,
			STLPosition:      &position,
		}
		styleAttributes.propagateSTLAttributes()

		// Create item
		var i = &Item{
			EndAt:       t.timecodeOut - o.Metadata.STLTimecodeStartOfProgramme,
			InlineStyle: &styleAttributes,
			StartAt:     t.timecodeIn - o.Metadata.STLTimecodeStartOfProgramme,
		}

		// Loop through rows
		for _, text := range bytes.Split(t.text, []byte{stlLineSeparator}) {
			if g.displayStandardCode == stlDisplayStandardCodeOpenSubtitling {
				err = parseOpenSubtitleRow(i, ch, func() styler { return newSTLStyler() }, text)
				if err != nil {
					return nil, err
				}
			} else {
				parseTeletextRow(i, ch, func() styler { return newSTLStyler() }, text)
			}
		}

		// Append item
		o.Items = append(o.Items, i)

	}
	return
}

// readNBytes reads n bytes
func readNBytes(i io.Reader, c int) (o []byte, err error) {
	o = make([]byte, c)
	var n int
	if n, err = i.Read(o); err != nil || n != len(o) {
		if err != nil {
			if err == io.EOF {
				return
			}
			err = fmt.Errorf("astisub: reading %d bytes failed: %w", c, err)
			return
		}
		err = fmt.Errorf("astisub: read %d bytes, should have read %d", n, c)
		return
	}
	return
}

// gsiBlock represents a GSI block
type gsiBlock struct {
	characterCodeTableNumber                         uint16
	codePageNumber                                   uint32
	countryOfOrigin                                  string
	creationDate                                     time.Time
	diskSequenceNumber                               int
	displayStandardCode                              string
	editorContactDetails                             string
	editorName                                       string
	framerate                                        int
	languageCode                                     string
	maximumNumberOfDisplayableCharactersInAnyTextRow int
	maximumNumberOfDisplayableRows                   int
	originalEpisodeTitle                             string
	originalProgramTitle                             string
	publisher                                        string
	revisionDate                                     time.Time
	revisionNumber                                   int
	subtitleListReferenceCode                        string
	timecodeFirstInCue                               time.Duration
	timecodeStartOfProgramme                         time.Duration
	timecodeStatus                                   string
	totalNumberOfDisks                               int
	totalNumberOfSubtitleGroups                      int
	totalNumberOfSubtitles                           int
	totalNumberOfTTIBlocks                           int
	translatedEpisodeTitle                           string
	translatedProgramTitle                           string
	translatorContactDetails                         string
	translatorName                                   string
	userDefinedArea                                  string
}

// newGSIBlock builds the subtitles GSI block
func newGSIBlock(s Subtitles) (g *gsiBlock) {
	// Init
	g = &gsiBlock{
		characterCodeTableNumber: stlCharacterCodeTableNumberLatin,
		codePageNumber:           stlCodePageNumberMultilingual,
		countryOfOrigin:          stlCountryCodeFrance,
		creationDate:             Now(),
		diskSequenceNumber:       1,
		displayStandardCode:      stlDisplayStandardCodeLevel1Teletext,
		framerate:                25,
		languageCode:             stlLanguageCodeFrench,
		maximumNumberOfDisplayableCharactersInAnyTextRow: 40,
		maximumNumberOfDisplayableRows:                   23,
		revisionDate:                                     Now(),
		subtitleListReferenceCode:                        "",
		timecodeStatus:                                   stlTimecodeStatusIntendedForUse,
		timecodeStartOfProgramme:                         0,
		totalNumberOfDisks:                               1,
		totalNumberOfSubtitleGroups:                      1,
		totalNumberOfSubtitles:                           len(s.Items),
		totalNumberOfTTIBlocks:                           len(s.Items),
	}

	// Add metadata
	if s.Metadata != nil {
		if s.Metadata.STLCreationDate != nil {
			g.creationDate = *s.Metadata.STLCreationDate
		}
		g.countryOfOrigin = s.Metadata.STLCountryOfOrigin
		g.displayStandardCode = s.Metadata.STLDisplayStandardCode
		g.framerate = s.Metadata.Framerate
		if v, ok := stlLanguageMapping.GetInverse(s.Metadata.Language); ok {
			g.languageCode = v.(string)
		}
		g.originalProgramTitle = s.Metadata.Title
		if s.Metadata.STLMaximumNumberOfDisplayableCharactersInAnyTextRow != nil {
			g.maximumNumberOfDisplayableCharactersInAnyTextRow = *s.Metadata.STLMaximumNumberOfDisplayableCharactersInAnyTextRow
		}
		if s.Metadata.STLMaximumNumberOfDisplayableRows != nil {
			g.maximumNumberOfDisplayableRows = *s.Metadata.STLMaximumNumberOfDisplayableRows
		}
		g.publisher = s.Metadata.STLPublisher
		if s.Metadata.STLRevisionDate != nil {
			g.revisionDate = *s.Metadata.STLRevisionDate
		}
		g.subtitleListReferenceCode = s.Metadata.STLSubtitleListReferenceCode
		g.timecodeStartOfProgramme = s.Metadata.STLTimecodeStartOfProgramme
	}

	// Timecode first in cue
	if len(s.Items) > 0 {
		g.timecodeFirstInCue = s.Items[0].StartAt
	}
	return
}

// parseGSIBlock parses a GSI block
func parseGSIBlock(b []byte) (g *gsiBlock, err error) {
	// Init
	g = &gsiBlock{
		characterCodeTableNumber:  binary.BigEndian.Uint16(b[12:14]),
		countryOfOrigin:           string(bytes.TrimSpace(b[274:277])),
		codePageNumber:            binary.BigEndian.Uint32(append([]byte{0x0}, b[0:3]...)),
		displayStandardCode:       string(bytes.TrimSpace([]byte{b[11]})),
		editorName:                string(bytes.TrimSpace(b[309:341])),
		editorContactDetails:      string(bytes.TrimSpace(b[341:373])),
		languageCode:              string(bytes.TrimSpace(b[14:16])),
		originalEpisodeTitle:      string(bytes.TrimSpace(b[48:80])),
		originalProgramTitle:      string(bytes.TrimSpace(b[16:48])),
		publisher:                 string(bytes.TrimSpace(b[277:309])),
		subtitleListReferenceCode: string(bytes.TrimSpace(b[208:224])),
		timecodeStatus:            string(bytes.TrimSpace([]byte{b[255]})),
		translatedEpisodeTitle:    string(bytes.TrimSpace(b[80:112])),
		translatedProgramTitle:    string(bytes.TrimSpace(b[112:144])),
		translatorContactDetails:  string(bytes.TrimSpace(b[176:208])),
		translatorName:            string(bytes.TrimSpace(b[144:176])),
		userDefinedArea:           string(bytes.TrimSpace(b[448:])),
	}

	// Framerate
	if v, ok := stlFramerateMapping.Get(string(b[3:11])); ok {
		g.framerate = v.(int)
	}

	// Creation date
	if v := strings.TrimSpace(string(b[224:230])); len(v) > 0 {
		if g.creationDate, err = time.Parse("060102", v); err != nil {
			err = fmt.Errorf("astisub: parsing date %s failed: %w", v, err)
			return
		}
	}

	// Revision date
	if v := strings.TrimSpace(string(b[230:236])); len(v) > 0 {
		if g.revisionDate, err = time.Parse("060102", v); err != nil {
			err = fmt.Errorf("astisub: parsing date %s failed: %w", v, err)
			return
		}
	}

	// Revision number
	if v := strings.TrimSpace(string(b[236:238])); len(v) > 0 {
		if g.revisionNumber, err = strconv.Atoi(v); err != nil {
			err = fmt.Errorf("astisub: atoi of %s failed: %w", v, err)
			return
		}
	}

	// Total number of TTI blocks
	if v := strings.TrimSpace(string(b[238:243])); len(v) > 0 {
		if g.totalNumberOfTTIBlocks, err = strconv.Atoi(v); err != nil {
			err = fmt.Errorf("astisub: atoi of %s failed: %w", v, err)
			return
		}
	}

	// Total number of subtitles
	if v := strings.TrimSpace(string(b[243:248])); len(v) > 0 {
		if g.totalNumberOfSubtitles, err = strconv.Atoi(v); err != nil {
			err = fmt.Errorf("astisub: atoi of %s failed: %w", v, err)
			return
		}
	}

	// Total number of subtitle groups
	if v := strings.TrimSpace(string(b[248:251])); len(v) > 0 {
		if g.totalNumberOfSubtitleGroups, err = strconv.Atoi(v); err != nil {
			err = fmt.Errorf("astisub: atoi of %s failed: %w", v, err)
			return
		}
	}

	// Maximum number of displayable characters in any text row
	if v := strings.TrimSpace(string(b[251:253])); len(v) > 0 {
		if g.maximumNumberOfDisplayableCharactersInAnyTextRow, err = strconv.Atoi(v); err != nil {
			err = fmt.Errorf("astisub: atoi of %s failed: %w", v, err)
			return
		}
	}

	// Maximum number of displayable rows
	if v := strings.TrimSpace(string(b[253:255])); len(v) > 0 {
		if g.maximumNumberOfDisplayableRows, err = strconv.Atoi(v); err != nil {
			err = fmt.Errorf("astisub: atoi of %s failed: %w", v, err)
			return
		}
	}

	// Timecode start of programme
	if v := strings.TrimSpace(string(b[256:264])); len(v) > 0 {
		if g.timecodeStartOfProgramme, err = parseDurationSTL(v, g.framerate); err != nil {
			err = fmt.Errorf("astisub: parsing of stl duration %s failed: %w", v, err)
			return
		}
	}

	// Timecode first in cue
	if v := strings.TrimSpace(string(b[264:272])); len(v) > 0 {
		if g.timecodeFirstInCue, err = parseDurationSTL(v, g.framerate); err != nil {
			err = fmt.Errorf("astisub: parsing of stl duration %s failed: %w", v, err)
			return
		}
	}

	// Total number of disks
	if v := strings.TrimSpace(string(b[272])); len(v) > 0 {
		if g.totalNumberOfDisks, err = strconv.Atoi(v); err != nil {
			err = fmt.Errorf("astisub: atoi of %s failed: %w", v, err)
			return
		}
	}

	// Disk sequence number
	if v := strings.TrimSpace(string(b[273])); len(v) > 0 {
		if g.diskSequenceNumber, err = strconv.Atoi(v); err != nil {
			err = fmt.Errorf("astisub: atoi of %s failed: %w", v, err)
			return
		}
	}
	return
}

// bytes transforms the GSI block into []byte
func (b gsiBlock) bytes() (o []byte) {
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, b.codePageNumber)
	o = append(o, astikit.BytesPad(bs[1:], ' ', 3, astikit.PadRight, astikit.PadCut)...) // Code page number
	// Disk format code
	var f string
	if v, ok := stlFramerateMapping.GetInverse(b.framerate); ok {
		f = v.(string)
	}
	o = append(o, astikit.BytesPad([]byte(f), ' ', 8, astikit.PadRight, astikit.PadCut)...)
	o = append(o, astikit.BytesPad([]byte(b.displayStandardCode), ' ', 1, astikit.PadRight, astikit.PadCut)...) // Display standard code
	binary.BigEndian.PutUint16(bs, b.characterCodeTableNumber)
	o = append(o, astikit.BytesPad(bs[:2], ' ', 2, astikit.PadRight, astikit.PadCut)...)                                                             // Character code table number
	o = append(o, astikit.BytesPad([]byte(b.languageCode), ' ', 2, astikit.PadRight, astikit.PadCut)...)                                             // Language code
	o = append(o, astikit.BytesPad([]byte(b.originalProgramTitle), ' ', 32, astikit.PadRight, astikit.PadCut)...)                                    // Original program title
	o = append(o, astikit.BytesPad([]byte(b.originalEpisodeTitle), ' ', 32, astikit.PadRight, astikit.PadCut)...)                                    // Original episode title
	o = append(o, astikit.BytesPad([]byte(b.translatedProgramTitle), ' ', 32, astikit.PadRight, astikit.PadCut)...)                                  // Translated program title
	o = append(o, astikit.BytesPad([]byte(b.translatedEpisodeTitle), ' ', 32, astikit.PadRight, astikit.PadCut)...)                                  // Translated episode title
	o = append(o, astikit.BytesPad([]byte(b.translatorName), ' ', 32, astikit.PadRight, astikit.PadCut)...)                                          // Translator's name
	o = append(o, astikit.BytesPad([]byte(b.translatorContactDetails), ' ', 32, astikit.PadRight, astikit.PadCut)...)                                // Translator's contact details
	o = append(o, astikit.BytesPad([]byte(b.subtitleListReferenceCode), ' ', 16, astikit.PadRight, astikit.PadCut)...)                               // Subtitle list reference code
	o = append(o, astikit.BytesPad([]byte(b.creationDate.Format("060102")), ' ', 6, astikit.PadRight, astikit.PadCut)...)                            // Creation date
	o = append(o, astikit.BytesPad([]byte(b.revisionDate.Format("060102")), ' ', 6, astikit.PadRight, astikit.PadCut)...)                            // Revision date
	o = append(o, astikit.BytesPad([]byte(strconv.Itoa(b.revisionNumber)), '0', 2, astikit.PadCut)...)                                               // Revision number
	o = append(o, astikit.BytesPad([]byte(strconv.Itoa(b.totalNumberOfTTIBlocks)), '0', 5, astikit.PadCut)...)                                       // Total number of TTI blocks
	o = append(o, astikit.BytesPad([]byte(strconv.Itoa(b.totalNumberOfSubtitles)), '0', 5, astikit.PadCut)...)                                       // Total number of subtitles
	o = append(o, astikit.BytesPad([]byte(strconv.Itoa(b.totalNumberOfSubtitleGroups)), '0', 3, astikit.PadCut)...)                                  // Total number of subtitle groups
	o = append(o, astikit.BytesPad([]byte(strconv.Itoa(b.maximumNumberOfDisplayableCharactersInAnyTextRow)), '0', 2, astikit.PadCut)...)             // Maximum number of displayable characters in any text row
	o = append(o, astikit.BytesPad([]byte(strconv.Itoa(b.maximumNumberOfDisplayableRows)), '0', 2, astikit.PadCut)...)                               // Maximum number of displayable rows
	o = append(o, astikit.BytesPad([]byte(b.timecodeStatus), ' ', 1, astikit.PadRight, astikit.PadCut)...)                                           // Timecode status
	o = append(o, astikit.BytesPad([]byte(formatDurationSTL(b.timecodeStartOfProgramme, b.framerate)), ' ', 8, astikit.PadRight, astikit.PadCut)...) // Timecode start of a programme
	o = append(o, astikit.BytesPad([]byte(formatDurationSTL(b.timecodeFirstInCue, b.framerate)), ' ', 8, astikit.PadRight, astikit.PadCut)...)       // Timecode first in cue
	o = append(o, astikit.BytesPad([]byte(strconv.Itoa(b.totalNumberOfDisks)), ' ', 1, astikit.PadRight, astikit.PadCut)...)                         // Total number of disks
	o = append(o, astikit.BytesPad([]byte(strconv.Itoa(b.diskSequenceNumber)), ' ', 1, astikit.PadRight, astikit.PadCut)...)                         // Disk sequence number
	o = append(o, astikit.BytesPad([]byte(b.countryOfOrigin), ' ', 3, astikit.PadRight, astikit.PadCut)...)                                          // Country of origin
	o = append(o, astikit.BytesPad([]byte(b.publisher), ' ', 32, astikit.PadRight, astikit.PadCut)...)                                               // Publisher
	o = append(o, astikit.BytesPad([]byte(b.editorName), ' ', 32, astikit.PadRight, astikit.PadCut)...)                                              // Editor's name
	o = append(o, astikit.BytesPad([]byte(b.editorContactDetails), ' ', 32, astikit.PadRight, astikit.PadCut)...)                                    // Editor's contact details
	o = append(o, astikit.BytesPad([]byte{}, ' ', 75+576, astikit.PadRight, astikit.PadCut)...)                                                      // Spare bytes + user defined area                                                                                           //                                                                                                                      // Editor's contact details
	return
}

// parseDurationSTL parses a STL duration
func parseDurationSTL(i string, framerate int) (d time.Duration, err error) {
	// Parse hours
	var hours, hoursString = 0, i[0:2]
	if hours, err = strconv.Atoi(hoursString); err != nil {
		err = fmt.Errorf("astisub: atoi of %s failed: %w", hoursString, err)
		return
	}

	// Parse minutes
	var minutes, minutesString = 0, i[2:4]
	if minutes, err = strconv.Atoi(minutesString); err != nil {
		err = fmt.Errorf("astisub: atoi of %s failed: %w", minutesString, err)
		return
	}

	// Parse seconds
	var seconds, secondsString = 0, i[4:6]
	if seconds, err = strconv.Atoi(secondsString); err != nil {
		err = fmt.Errorf("astisub: atoi of %s failed: %w", secondsString, err)
		return
	}

	// Parse frames
	var frames, framesString = 0, i[6:8]
	if frames, err = strconv.Atoi(framesString); err != nil {
		err = fmt.Errorf("astisub: atoi of %s failed: %w", framesString, err)
		return
	}

	// Set duration
	d = time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second + time.Duration(1e9*frames/framerate)*time.Nanosecond
	return
}

// formatDurationSTL formats a STL duration
func formatDurationSTL(d time.Duration, framerate int) (o string) {
	// Add hours
	if d.Hours() < 10 {
		o += "0"
	}
	var delta = int(math.Floor(d.Hours()))
	o += strconv.Itoa(delta)
	d -= time.Duration(delta) * time.Hour

	// Add minutes
	if d.Minutes() < 10 {
		o += "0"
	}
	delta = int(math.Floor(d.Minutes()))
	o += strconv.Itoa(delta)
	d -= time.Duration(delta) * time.Minute

	// Add seconds
	if d.Seconds() < 10 {
		o += "0"
	}
	delta = int(math.Floor(d.Seconds()))
	o += strconv.Itoa(delta)
	d -= time.Duration(delta) * time.Second

	// Add frames
	var frames = int(int(d.Nanoseconds()) * framerate / 1e9)
	if frames < 10 {
		o += "0"
	}
	o += strconv.Itoa(frames)
	return
}

// ttiBlock represents a TTI block
type ttiBlock struct {
	commentFlag          byte
	cumulativeStatus     byte
	extensionBlockNumber int
	justificationCode    byte
	subtitleGroupNumber  int
	subtitleNumber       int
	text                 []byte
	timecodeIn           time.Duration
	timecodeOut          time.Duration
	verticalPosition     int
}

// newTTIBlock builds an item TTI block
func newTTIBlock(i *Item, idx int) (t *ttiBlock) {
	// Init
	t = &ttiBlock{
		commentFlag:          stlCommentFlagTextContainsSubtitleData,
		cumulativeStatus:     stlCumulativeStatusSubtitleNotPartOfACumulativeSet,
		extensionBlockNumber: 255,
		justificationCode:    stlJustificationCodeLeftJustifiedText,
		subtitleGroupNumber:  0,
		subtitleNumber:       idx,
		timecodeIn:           i.StartAt,
		timecodeOut:          i.EndAt,
		verticalPosition:     stlVerticalPositionFromStyle(i.InlineStyle),
	}

	// Add text
	var lines []string
	for _, l := range i.Lines {
		var lineItems []string
		for _, li := range l.Items {
			lineItems = append(lineItems, li.STLString())
		}
		lines = append(lines, strings.Join(lineItems, " "))
	}
	t.text = []byte(strings.Join(lines, string(rune(stlLineSeparator))))
	return
}

func stlVerticalPositionFromStyle(sa *StyleAttributes) int {
	if sa != nil && sa.STLPosition != nil {
		return sa.STLPosition.VerticalPosition
	} else {
		return 20
	}
}

func (li LineItem) STLString() string {
	rs := li.Text
	if li.InlineStyle != nil {
		if li.InlineStyle.STLItalics != nil && *li.InlineStyle.STLItalics {
			rs = string(rune(0x80)) + rs + string(rune(0x81))
		}
		if li.InlineStyle.STLUnderline != nil && *li.InlineStyle.STLUnderline {
			rs = string(rune(0x82)) + rs + string(rune(0x83))
		}
		if li.InlineStyle.STLBoxing != nil && *li.InlineStyle.STLBoxing {
			rs = string(rune(0x84)) + rs + string(rune(0x85))
		}
	}
	return rs
}

// parseTTIBlock parses a TTI block
func parseTTIBlock(p []byte, framerate int) *ttiBlock {
	return &ttiBlock{
		commentFlag:          p[15],
		cumulativeStatus:     p[4],
		extensionBlockNumber: int(uint8(p[3])),
		justificationCode:    p[14],
		subtitleGroupNumber:  int(uint8(p[0])),
		subtitleNumber:       int(binary.LittleEndian.Uint16(p[1:3])),
		text:                 p[16:128],
		timecodeIn:           parseDurationSTLBytes(p[5:9], framerate),
		timecodeOut:          parseDurationSTLBytes(p[9:13], framerate),
		verticalPosition:     int(uint8(p[13])),
	}
}

// bytes transforms the TTI block into []byte
func (t *ttiBlock) bytes(g *gsiBlock) (o []byte) {
	o = append(o, byte(uint8(t.subtitleGroupNumber))) // Subtitle group number
	var b = make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(t.subtitleNumber))
	o = append(o, b...)                                                                                              // Subtitle number
	o = append(o, byte(uint8(t.extensionBlockNumber)))                                                               // Extension block number
	o = append(o, t.cumulativeStatus)                                                                                // Cumulative status
	o = append(o, formatDurationSTLBytes(t.timecodeIn, g.framerate)...)                                              // Timecode in
	o = append(o, formatDurationSTLBytes(t.timecodeOut, g.framerate)...)                                             // Timecode out
	o = append(o, validateVerticalPosition(t.verticalPosition, g.displayStandardCode))                               // Vertical position
	o = append(o, t.justificationCode)                                                                               // Justification code
	o = append(o, t.commentFlag)                                                                                     // Comment flag
	o = append(o, astikit.BytesPad(encodeTextSTL(string(t.text)), '\x8f', 112, astikit.PadRight, astikit.PadCut)...) // Text field
	return
}

// According to EBU 3264 (https://tech.ebu.ch/docs/tech/tech3264.pdf):
// page 12:
// for teletext subtitles, VP contains a value in the range 1-23 decimal (01h-17h)
// corresponding to theteletext row number of the first subtitle row.
// page 6:
// Teletext ("closed") subtitles are indicated via the Display Standard Code
// in the GSI block.
func validateVerticalPosition(vp int, dsc string) byte {
	closed := false
	switch dsc {
	case stlDisplayStandardCodeLevel1Teletext, stlDisplayStandardCodeLevel2Teletext:
		closed = true
	}
	if vp < 1 && closed {
		vp = 1
	}
	if vp > 23 && closed {
		vp = 23
	}
	return byte(uint8(vp))
}

// formatDurationSTLBytes formats a STL duration in bytes
func formatDurationSTLBytes(d time.Duration, framerate int) (o []byte) {
	// Add hours
	var hours = int(math.Floor(d.Hours()))
	o = append(o, byte(uint8(hours)))
	d -= time.Duration(hours) * time.Hour

	// Add minutes
	var minutes = int(math.Floor(d.Minutes()))
	o = append(o, byte(uint8(minutes)))
	d -= time.Duration(minutes) * time.Minute

	// Add seconds
	var seconds = int(math.Floor(d.Seconds()))
	o = append(o, byte(uint8(seconds)))
	d -= time.Duration(seconds) * time.Second

	// Add frames
	var frames = int(int(d.Nanoseconds()) * framerate / 1e9)
	o = append(o, byte(uint8(frames)))
	return
}

// parseDurationSTLBytes parses a STL duration in bytes
func parseDurationSTLBytes(b []byte, framerate int) time.Duration {
	return time.Duration(uint8(b[0]))*time.Hour + time.Duration(uint8(b[1]))*time.Minute + time.Duration(uint8(b[2]))*time.Second + time.Duration(1e9*int(uint8(b[3]))/framerate)*time.Nanosecond
}

type stlCharacterHandler struct {
	accent string
	c      uint16
	m      *astikit.BiMap
}

func newSTLCharacterHandler(characterCodeTable uint16) (*stlCharacterHandler, error) {
	if v, ok := stlCharacterCodeTables[characterCodeTable]; ok {
		return &stlCharacterHandler{
			c: characterCodeTable,
			m: v,
		}, nil
	}
	return nil, fmt.Errorf("astisub: table doesn't exist for character code table %d", characterCodeTable)
}

// TODO Use this instead of encodeTextSTL => use in teletext process like for decode
// TODO Test
func (h *stlCharacterHandler) encode(i []byte) byte {
	return ' '
}

func (h *stlCharacterHandler) decode(i byte) (o []byte) {
	k := int(i)
	vi, ok := h.m.Get(k)
	if !ok {
		return
	}
	v := vi.(string)
	if len(h.accent) > 0 {
		o = norm.NFC.Bytes([]byte(v + h.accent))
		h.accent = ""
		return
	} else if h.c == stlCharacterCodeTableNumberLatin && k >= 0xc0 && k <= 0xcf {
		h.accent = v
		return
	}
	return []byte(v)
}

type stlStyler struct {
	boxing    *bool
	italics   *bool
	underline *bool
}

func newSTLStyler() *stlStyler {
	return &stlStyler{}
}

func (s *stlStyler) parseSpacingAttribute(i byte) {
	switch i {
	case 0x80:
		s.italics = astikit.BoolPtr(true)
	case 0x81:
		s.italics = astikit.BoolPtr(false)
	case 0x82:
		s.underline = astikit.BoolPtr(true)
	case 0x83:
		s.underline = astikit.BoolPtr(false)
	case 0x84:
		s.boxing = astikit.BoolPtr(true)
	case 0x85:
		s.boxing = astikit.BoolPtr(false)
	}
}

func (s *stlStyler) hasBeenSet() bool {
	return s.italics != nil || s.boxing != nil || s.underline != nil
}

func (s *stlStyler) hasChanged(sa *StyleAttributes) bool {
	return s.boxing != sa.STLBoxing || s.italics != sa.STLItalics || s.underline != sa.STLUnderline
}

func (s *stlStyler) propagateStyleAttributes(sa *StyleAttributes) {
	sa.propagateSTLAttributes()
}

func (s *stlStyler) update(sa *StyleAttributes) {
	if s.boxing != nil && s.boxing != sa.STLBoxing {
		sa.STLBoxing = s.boxing
	}
	if s.italics != nil && s.italics != sa.STLItalics {
		sa.STLItalics = s.italics
	}
	if s.underline != nil && s.underline != sa.STLUnderline {
		sa.STLUnderline = s.underline
	}
}

// WriteToSTL writes subtitles in .stl format
func (s Subtitles) WriteToSTL(o io.Writer) (err error) {
	// Do not write anything if no subtitles
	if len(s.Items) == 0 {
		err = ErrNoSubtitlesToWrite
		return
	}

	// Write GSI block
	var g = newGSIBlock(s)
	if _, err = o.Write(g.bytes()); err != nil {
		err = fmt.Errorf("astisub: writing gsi block failed: %w", err)
		return
	}

	// Loop through items
	for idx, item := range s.Items {
		// Write tti block
		if _, err = o.Write(newTTIBlock(item, idx+1).bytes(g)); err != nil {
			err = fmt.Errorf("astisub: writing tti block #%d failed: %w", idx+1, err)
			return
		}
	}
	return
}

// TODO Remove below

// STL unicode diacritic
var stlUnicodeDiacritic = astikit.NewBiMap().
	Set(byte('\xc1'), "\u0300"). // Grave accent
	Set(byte('\xc2'), "\u0301"). // Acute accent
	Set(byte('\xc3'), "\u0302"). // Circumflex
	Set(byte('\xc4'), "\u0303"). // Tilde
	Set(byte('\xc5'), "\u0304"). // Macron
	Set(byte('\xc6'), "\u0306"). // Breve
	Set(byte('\xc7'), "\u0307"). // Dot
	Set(byte('\xc8'), "\u0308"). // Umlaut
	Set(byte('\xca'), "\u030a"). // Ring
	Set(byte('\xcb'), "\u0327"). // Cedilla
	Set(byte('\xcd'), "\u030B"). // Double acute accent
	Set(byte('\xce'), "\u0328"). // Ogonek
	Set(byte('\xcf'), "\u030c")  // Caron

// STL unicode mapping
var stlUnicodeMapping = astikit.NewBiMap().
	Set(byte('\x8a'), "\u000a"). // Line break
	Set(byte('\xa8'), "\u00a4"). // ¤
	Set(byte('\xa9'), "\u2018"). // ‘
	Set(byte('\xaa'), "\u201C"). // “
	Set(byte('\xab'), "\u00AB"). // «
	Set(byte('\xac'), "\u2190"). // ←
	Set(byte('\xad'), "\u2191"). // ↑
	Set(byte('\xae'), "\u2192"). // →
	Set(byte('\xaf'), "\u2193"). // ↓
	Set(byte('\xb4'), "\u00D7"). // ×
	Set(byte('\xb8'), "\u00F7"). // ÷
	Set(byte('\xb9'), "\u2019"). // ’
	Set(byte('\xba'), "\u201D"). // ”
	Set(byte('\xbc'), "\u00BC"). // ¼
	Set(byte('\xbd'), "\u00BD"). // ½
	Set(byte('\xbe'), "\u00BE"). // ¾
	Set(byte('\xbf'), "\u00BF"). // ¿
	Set(byte('\xd0'), "\u2015"). // ―
	Set(byte('\xd1'), "\u00B9"). // ¹
	Set(byte('\xd2'), "\u00AE"). // ®
	Set(byte('\xd3'), "\u00A9"). // ©
	Set(byte('\xd4'), "\u2122"). // ™
	Set(byte('\xd5'), "\u266A"). // ♪
	Set(byte('\xd6'), "\u00AC"). // ¬
	Set(byte('\xd7'), "\u00A6"). // ¦
	Set(byte('\xdc'), "\u215B"). // ⅛
	Set(byte('\xdd'), "\u215C"). // ⅜
	Set(byte('\xde'), "\u215D"). // ⅝
	Set(byte('\xdf'), "\u215E"). // ⅞
	Set(byte('\xe0'), "\u2126"). // Ohm Ω
	Set(byte('\xe1'), "\u00C6"). // Æ
	Set(byte('\xe2'), "\u0110"). // Đ
	Set(byte('\xe3'), "\u00AA"). // ª
	Set(byte('\xe4'), "\u0126"). // Ħ
	Set(byte('\xe6'), "\u0132"). // Ĳ
	Set(byte('\xe7'), "\u013F"). // Ŀ
	Set(byte('\xe8'), "\u0141"). // Ł
	Set(byte('\xe9'), "\u00D8"). // Ø
	Set(byte('\xea'), "\u0152"). // Œ
	Set(byte('\xeb'), "\u00BA"). // º
	Set(byte('\xec'), "\u00DE"). // Þ
	Set(byte('\xed'), "\u0166"). // Ŧ
	Set(byte('\xee'), "\u014A"). // Ŋ
	Set(byte('\xef'), "\u0149"). // ŉ
	Set(byte('\xf0'), "\u0138"). // ĸ
	Set(byte('\xf1'), "\u00E6"). // æ
	Set(byte('\xf2'), "\u0111"). // đ
	Set(byte('\xf3'), "\u00F0"). // ð
	Set(byte('\xf4'), "\u0127"). // ħ
	Set(byte('\xf5'), "\u0131"). // ı
	Set(byte('\xf6'), "\u0133"). // ĳ
	Set(byte('\xf7'), "\u0140"). // ŀ
	Set(byte('\xf8'), "\u0142"). // ł
	Set(byte('\xf9'), "\u00F8"). // ø
	Set(byte('\xfa'), "\u0153"). // œ
	Set(byte('\xfb'), "\u00DF"). // ß
	Set(byte('\xfc'), "\u00FE"). // þ
	Set(byte('\xfd'), "\u0167"). // ŧ
	Set(byte('\xfe'), "\u014B"). // ŋ
	Set(byte('\xff'), "\u00AD")  // Soft hyphen

// encodeTextSTL encodes the STL text
func encodeTextSTL(i string) (o []byte) {
	i = string(norm.NFD.Bytes([]byte(i)))
	for _, c := range i {
		if v, ok := stlUnicodeMapping.GetInverse(string(c)); ok {
			o = append(o, v.(byte))
		} else if v, ok := stlUnicodeDiacritic.GetInverse(string(c)); ok {
			o = append(o[:len(o)-1], v.(byte), o[len(o)-1])
		} else {
			o = append(o, byte(c))
		}
	}
	return
}

func parseSTLJustificationCode(i byte) Justification {
	switch i {
	case 0x00:
		return JustificationUnchanged
	case 0x01:
		return JustificationLeft
	case 0x02:
		return JustificationCentered
	case 0x03:
		return JustificationRight
	default:
		return JustificationUnchanged
	}
}

func isTeletextControlCode(i byte) (b bool) {
	return i <= 0x1f
}

func parseOpenSubtitleRow(i *Item, d decoder, fs func() styler, row []byte) error {
	// Loop through columns
	var l = Line{}
	var li = LineItem{InlineStyle: &StyleAttributes{}}
	var s styler
	for _, v := range row {
		// Create specific styler
		if fs != nil {
			s = fs()
		}

		if isTeletextControlCode(v) {
			return errors.New("teletext control code in open text")
		}
		if s != nil {
			s.parseSpacingAttribute(v)
		}

		// Style has been set
		if s != nil && s.hasBeenSet() {
			// Style has changed
			if s.hasChanged(li.InlineStyle) {
				if len(li.Text) > 0 {
					// Append line item
					appendOpenSubtitleLineItem(&l, li, s)

					// Create new line item
					sa := &StyleAttributes{}
					*sa = *li.InlineStyle
					li = LineItem{InlineStyle: sa}
				}
				s.update(li.InlineStyle)
			}
		} else {
			// Append text
			li.Text += string(d.decode(v))
		}
	}

	appendOpenSubtitleLineItem(&l, li, s)

	// Append line
	if len(l.Items) > 0 {
		i.Lines = append(i.Lines, l)
	}
	return nil
}

func appendOpenSubtitleLineItem(l *Line, li LineItem, s styler) {
	// There's some text
	if len(strings.TrimSpace(li.Text)) > 0 {
		// Make sure inline style exists
		if li.InlineStyle == nil {
			li.InlineStyle = &StyleAttributes{}
		}

		// Propagate style attributes
		if s != nil {
			s.propagateStyleAttributes(li.InlineStyle)
		}

		// Append line item
		li.Text = strings.TrimSpace(li.Text)
		l.Items = append(l.Items, li)
	}
}
