package astits

import (
	"fmt"
	"time"

	"github.com/asticode/go-astikit"
)

// Audio types
// Page: 683 | https://books.google.fr/books?id=6dgWB3-rChYC&printsec=frontcover&hl=fr
const (
	AudioTypeCleanEffects             = 0x1
	AudioTypeHearingImpaired          = 0x2
	AudioTypeVisualImpairedCommentary = 0x3
)

// Data stream alignments
// Page: 85 | Chapter:2.6.11 | Link: http://ecee.colorado.edu/~ecen5653/ecen5653/papers/iso13818-1.pdf
const (
	DataStreamAligmentAudioSyncWord          = 0x1
	DataStreamAligmentVideoSliceOrAccessUnit = 0x1
	DataStreamAligmentVideoAccessUnit        = 0x2
	DataStreamAligmentVideoGOPOrSEQ          = 0x3
	DataStreamAligmentVideoSEQ               = 0x4
)

// Descriptor tags
// Chapter: 6.1 | Link: https://www.etsi.org/deliver/etsi_en/300400_300499/300468/01.15.01_60/en_300468v011501p.pdf
const (
	DescriptorTagAC3                        = 0x6a
	DescriptorTagAVCVideo                   = 0x28
	DescriptorTagComponent                  = 0x50
	DescriptorTagContent                    = 0x54
	DescriptorTagDataStreamAlignment        = 0x6
	DescriptorTagEnhancedAC3                = 0x7a
	DescriptorTagExtendedEvent              = 0x4e
	DescriptorTagExtension                  = 0x7f
	DescriptorTagISO639LanguageAndAudioType = 0xa
	DescriptorTagLocalTimeOffset            = 0x58
	DescriptorTagMaximumBitrate             = 0xe
	DescriptorTagNetworkName                = 0x40
	DescriptorTagParentalRating             = 0x55
	DescriptorTagPrivateDataIndicator       = 0xf
	DescriptorTagPrivateDataSpecifier       = 0x5f
	DescriptorTagRegistration               = 0x5
	DescriptorTagService                    = 0x48
	DescriptorTagShortEvent                 = 0x4d
	DescriptorTagStreamIdentifier           = 0x52
	DescriptorTagSubtitling                 = 0x59
	DescriptorTagTeletext                   = 0x56
	DescriptorTagVBIData                    = 0x45
	DescriptorTagVBITeletext                = 0x46
)

// Descriptor extension tags
// Chapter: 6.3 | Link: https://www.etsi.org/deliver/etsi_en/300400_300499/300468/01.15.01_60/en_300468v011501p.pdf
const (
	DescriptorTagExtensionSupplementaryAudio = 0x6
)

// Service types
// Chapter: 6.2.33 | Link: https://www.etsi.org/deliver/etsi_en/300400_300499/300468/01.15.01_60/en_300468v011501p.pdf
const (
	ServiceTypeDigitalTelevisionService = 0x1
)

// Teletext types
// Chapter: 6.2.43 | Link: https://www.etsi.org/deliver/etsi_en/300400_300499/300468/01.15.01_60/en_300468v011501p.pdf
const (
	TeletextTypeAdditionalInformationPage                    = 0x3
	TeletextTypeInitialTeletextPage                          = 0x1
	TeletextTypeProgramSchedulePage                          = 0x4
	TeletextTypeTeletextSubtitlePage                         = 0x2
	TeletextTypeTeletextSubtitlePageForHearingImpairedPeople = 0x5
)

// VBI data service id
// Chapter: 6.2.47 | Link: https://www.etsi.org/deliver/etsi_en/300400_300499/300468/01.15.01_60/en_300468v011501p.pdf
const (
	VBIDataServiceIDClosedCaptioning     = 0x6
	VBIDataServiceIDEBUTeletext          = 0x1
	VBIDataServiceIDInvertedTeletext     = 0x2
	VBIDataServiceIDMonochrome442Samples = 0x7
	VBIDataServiceIDVPS                  = 0x4
	VBIDataServiceIDWSS                  = 0x5
)

// Descriptor represents a descriptor
// TODO Handle UTF8
type Descriptor struct {
	AC3                        *DescriptorAC3
	AVCVideo                   *DescriptorAVCVideo
	Component                  *DescriptorComponent
	Content                    *DescriptorContent
	DataStreamAlignment        *DescriptorDataStreamAlignment
	EnhancedAC3                *DescriptorEnhancedAC3
	ExtendedEvent              *DescriptorExtendedEvent
	Extension                  *DescriptorExtension
	ISO639LanguageAndAudioType *DescriptorISO639LanguageAndAudioType
	Length                     uint8
	LocalTimeOffset            *DescriptorLocalTimeOffset
	MaximumBitrate             *DescriptorMaximumBitrate
	NetworkName                *DescriptorNetworkName
	ParentalRating             *DescriptorParentalRating
	PrivateDataIndicator       *DescriptorPrivateDataIndicator
	PrivateDataSpecifier       *DescriptorPrivateDataSpecifier
	Registration               *DescriptorRegistration
	Service                    *DescriptorService
	ShortEvent                 *DescriptorShortEvent
	StreamIdentifier           *DescriptorStreamIdentifier
	Subtitling                 *DescriptorSubtitling
	Tag                        uint8 // the tag defines the structure of the contained data following the descriptor length.
	Teletext                   *DescriptorTeletext
	Unknown                    *DescriptorUnknown
	UserDefined                []byte
	VBIData                    *DescriptorVBIData
	VBITeletext                *DescriptorTeletext
}

// DescriptorAC3 represents an AC3 descriptor
// Chapter: Annex D | Link: https://www.etsi.org/deliver/etsi_en/300400_300499/300468/01.15.01_60/en_300468v011501p.pdf
type DescriptorAC3 struct {
	AdditionalInfo   []byte
	ASVC             uint8
	BSID             uint8
	ComponentType    uint8
	HasASVC          bool
	HasBSID          bool
	HasComponentType bool
	HasMainID        bool
	MainID           uint8
}

func newDescriptorAC3(i *astikit.BytesIterator, offsetEnd int) (d *DescriptorAC3, err error) {
	// Get next byte
	var b byte
	if b, err = i.NextByte(); err != nil {
		err = fmt.Errorf("astits: fetching next byte failed: %w", err)
		return
	}

	// Create descriptor
	d = &DescriptorAC3{
		HasASVC:          uint8(b&0x10) > 0,
		HasBSID:          uint8(b&0x40) > 0,
		HasComponentType: uint8(b&0x80) > 0,
		HasMainID:        uint8(b&0x20) > 0,
	}

	// Component type
	if d.HasComponentType {
		if b, err = i.NextByte(); err != nil {
			err = fmt.Errorf("astits: fetching next byte failed: %w", err)
			return
		}
		d.ComponentType = uint8(b)
	}

	// BSID
	if d.HasBSID {
		if b, err = i.NextByte(); err != nil {
			err = fmt.Errorf("astits: fetching next byte failed: %w", err)
			return
		}
		d.BSID = uint8(b)
	}

	// Main ID
	if d.HasMainID {
		if b, err = i.NextByte(); err != nil {
			err = fmt.Errorf("astits: fetching next byte failed: %w", err)
			return
		}
		d.MainID = uint8(b)
	}

	// ASVC
	if d.HasASVC {
		if b, err = i.NextByte(); err != nil {
			err = fmt.Errorf("astits: fetching next byte failed: %w", err)
			return
		}
		d.ASVC = uint8(b)
	}

	// Additional info
	if i.Offset() < offsetEnd {
		if d.AdditionalInfo, err = i.NextBytes(offsetEnd - i.Offset()); err != nil {
			err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
			return
		}
	}
	return
}

// DescriptorAVCVideo represents an AVC video descriptor
// No doc found unfortunately, basing the implementation on https://github.com/gfto/bitstream/blob/master/mpeg/psi/desc_28.h
type DescriptorAVCVideo struct {
	AVC24HourPictureFlag bool
	AVCStillPresent      bool
	CompatibleFlags      uint8
	ConstraintSet0Flag   bool
	ConstraintSet1Flag   bool
	ConstraintSet2Flag   bool
	LevelIDC             uint8
	ProfileIDC           uint8
}

func newDescriptorAVCVideo(i *astikit.BytesIterator) (d *DescriptorAVCVideo, err error) {
	// Init
	d = &DescriptorAVCVideo{}

	// Get next byte
	var b byte
	if b, err = i.NextByte(); err != nil {
		err = fmt.Errorf("astits: fetching next byte failed: %w", err)
		return
	}

	// Profile idc
	d.ProfileIDC = uint8(b)

	// Get next byte
	if b, err = i.NextByte(); err != nil {
		err = fmt.Errorf("astits: fetching next byte failed: %w", err)
		return
	}

	// Flags
	d.ConstraintSet0Flag = b&0x80 > 0
	d.ConstraintSet1Flag = b&0x40 > 0
	d.ConstraintSet2Flag = b&0x20 > 0
	d.CompatibleFlags = b & 0x1f

	// Get next byte
	if b, err = i.NextByte(); err != nil {
		err = fmt.Errorf("astits: fetching next byte failed: %w", err)
		return
	}

	// Level idc
	d.LevelIDC = uint8(b)

	// Get next byte
	if b, err = i.NextByte(); err != nil {
		err = fmt.Errorf("astits: fetching next byte failed: %w", err)
		return
	}

	// AVC still present
	d.AVCStillPresent = b&0x80 > 0

	// AVC 24 hour picture flag
	d.AVC24HourPictureFlag = b&0x40 > 0
	return
}

// DescriptorComponent represents a component descriptor
// Chapter: 6.2.8 | Link: https://www.etsi.org/deliver/etsi_en/300400_300499/300468/01.15.01_60/en_300468v011501p.pdf
type DescriptorComponent struct {
	ComponentTag       uint8
	ComponentType      uint8
	ISO639LanguageCode []byte
	StreamContent      uint8
	StreamContentExt   uint8
	Text               []byte
}

func newDescriptorComponent(i *astikit.BytesIterator, offsetEnd int) (d *DescriptorComponent, err error) {
	// Init
	d = &DescriptorComponent{}

	// Get next byte
	var b byte
	if b, err = i.NextByte(); err != nil {
		err = fmt.Errorf("astits: fetching next byte failed: %w", err)
		return
	}

	// Stream content ext
	d.StreamContentExt = uint8(b >> 4)

	// Stream content
	d.StreamContent = uint8(b & 0xf)

	// Get next byte
	if b, err = i.NextByte(); err != nil {
		err = fmt.Errorf("astits: fetching next byte failed: %w", err)
		return
	}

	// Component type
	d.ComponentType = uint8(b)

	// Get next byte
	if b, err = i.NextByte(); err != nil {
		err = fmt.Errorf("astits: fetching next byte failed: %w", err)
		return
	}

	// Component tag
	d.ComponentTag = uint8(b)

	// ISO639 language code
	if d.ISO639LanguageCode, err = i.NextBytes(3); err != nil {
		err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
		return
	}

	// Text
	if i.Offset() < offsetEnd {
		if d.Text, err = i.NextBytes(offsetEnd - i.Offset()); err != nil {
			err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
			return
		}
	}
	return
}

// DescriptorContent represents a content descriptor
// Chapter: 6.2.9 | Link: https://www.etsi.org/deliver/etsi_en/300400_300499/300468/01.15.01_60/en_300468v011501p.pdf
type DescriptorContent struct {
	Items []*DescriptorContentItem
}

// DescriptorContentItem represents a content item descriptor
// Chapter: 6.2.9 | Link: https://www.etsi.org/deliver/etsi_en/300400_300499/300468/01.15.01_60/en_300468v011501p.pdf
type DescriptorContentItem struct {
	ContentNibbleLevel1 uint8
	ContentNibbleLevel2 uint8
	UserByte            uint8
}

func newDescriptorContent(i *astikit.BytesIterator, offsetEnd int) (d *DescriptorContent, err error) {
	// Init
	d = &DescriptorContent{}

	// Add items
	for i.Offset() < offsetEnd {
		// Get next bytes
		var bs []byte
		if bs, err = i.NextBytesNoCopy(2); err != nil {
			err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
			return
		}

		// Append item
		d.Items = append(d.Items, &DescriptorContentItem{
			ContentNibbleLevel1: uint8(bs[0] >> 4),
			ContentNibbleLevel2: uint8(bs[0] & 0xf),
			UserByte:            uint8(bs[1]),
		})
	}
	return
}

// DescriptorDataStreamAlignment represents a data stream alignment descriptor
type DescriptorDataStreamAlignment struct {
	Type uint8
}

func newDescriptorDataStreamAlignment(i *astikit.BytesIterator) (d *DescriptorDataStreamAlignment, err error) {
	var b byte
	if b, err = i.NextByte(); err != nil {
		err = fmt.Errorf("astits: fetching next byte failed: %w", err)
		return
	}
	d = &DescriptorDataStreamAlignment{Type: uint8(b)}
	return
}

// DescriptorEnhancedAC3 represents an enhanced AC3 descriptor
// Chapter: Annex D | Link: https://www.etsi.org/deliver/etsi_en/300400_300499/300468/01.15.01_60/en_300468v011501p.pdf
type DescriptorEnhancedAC3 struct {
	AdditionalInfo   []byte
	ASVC             uint8
	BSID             uint8
	ComponentType    uint8
	HasASVC          bool
	HasBSID          bool
	HasComponentType bool
	HasMainID        bool
	HasSubStream1    bool
	HasSubStream2    bool
	HasSubStream3    bool
	MainID           uint8
	MixInfoExists    bool
	SubStream1       uint8
	SubStream2       uint8
	SubStream3       uint8
}

func newDescriptorEnhancedAC3(i *astikit.BytesIterator, offsetEnd int) (d *DescriptorEnhancedAC3, err error) {
	// Get next byte
	var b byte
	if b, err = i.NextByte(); err != nil {
		err = fmt.Errorf("astits: fetching next byte failed: %w", err)
		return
	}

	// Create descriptor
	d = &DescriptorEnhancedAC3{
		HasASVC:          uint8(b&0x10) > 0,
		HasBSID:          uint8(b&0x40) > 0,
		HasComponentType: uint8(b&0x80) > 0,
		HasMainID:        uint8(b&0x20) > 0,
		HasSubStream1:    uint8(b&0x4) > 0,
		HasSubStream2:    uint8(b&0x2) > 0,
		HasSubStream3:    uint8(b&0x1) > 0,
		MixInfoExists:    uint8(b&0x8) > 0,
	}

	// Component type
	if d.HasComponentType {
		// Get next byte
		if b, err = i.NextByte(); err != nil {
			err = fmt.Errorf("astits: fetching next byte failed: %w", err)
			return
		}
		d.ComponentType = uint8(b)
	}

	// BSID
	if d.HasBSID {
		// Get next byte
		if b, err = i.NextByte(); err != nil {
			err = fmt.Errorf("astits: fetching next byte failed: %w", err)
			return
		}
		d.BSID = uint8(b)
	}

	// Main ID
	if d.HasMainID {
		// Get next byte
		if b, err = i.NextByte(); err != nil {
			err = fmt.Errorf("astits: fetching next byte failed: %w", err)
			return
		}
		d.MainID = uint8(b)
	}

	// ASVC
	if d.HasASVC {
		// Get next byte
		if b, err = i.NextByte(); err != nil {
			err = fmt.Errorf("astits: fetching next byte failed: %w", err)
			return
		}
		d.ASVC = uint8(b)
	}

	// Substream 1
	if d.HasSubStream1 {
		// Get next byte
		if b, err = i.NextByte(); err != nil {
			err = fmt.Errorf("astits: fetching next byte failed: %w", err)
			return
		}
		d.SubStream1 = uint8(b)
	}

	// Substream 2
	if d.HasSubStream2 {
		// Get next byte
		if b, err = i.NextByte(); err != nil {
			err = fmt.Errorf("astits: fetching next byte failed: %w", err)
			return
		}
		d.SubStream2 = uint8(b)
	}

	// Substream 3
	if d.HasSubStream3 {
		// Get next byte
		if b, err = i.NextByte(); err != nil {
			err = fmt.Errorf("astits: fetching next byte failed: %w", err)
			return
		}
		d.SubStream3 = uint8(b)
	}

	// Additional info
	if i.Offset() < offsetEnd {
		if d.AdditionalInfo, err = i.NextBytes(offsetEnd - i.Offset()); err != nil {
			err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
			return
		}
	}
	return
}

// DescriptorExtendedEvent represents an extended event descriptor
// Chapter: 6.2.15 | Link: https://www.etsi.org/deliver/etsi_en/300400_300499/300468/01.15.01_60/en_300468v011501p.pdf
type DescriptorExtendedEvent struct {
	ISO639LanguageCode   []byte
	Items                []*DescriptorExtendedEventItem
	LastDescriptorNumber uint8
	Number               uint8
	Text                 []byte
}

// DescriptorExtendedEventItem represents an extended event item descriptor
// Chapter: 6.2.15 | Link: https://www.etsi.org/deliver/etsi_en/300400_300499/300468/01.15.01_60/en_300468v011501p.pdf
type DescriptorExtendedEventItem struct {
	Content     []byte
	Description []byte
}

func newDescriptorExtendedEvent(i *astikit.BytesIterator) (d *DescriptorExtendedEvent, err error) {
	// Init
	d = &DescriptorExtendedEvent{}

	// Get next byte
	var b byte
	if b, err = i.NextByte(); err != nil {
		err = fmt.Errorf("astits: fetching next byte failed: %w", err)
		return
	}

	// Number
	d.Number = uint8(b >> 4)

	// Last descriptor number
	d.LastDescriptorNumber = uint8(b & 0xf)

	// ISO639 language code
	if d.ISO639LanguageCode, err = i.NextBytes(3); err != nil {
		err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
		return
	}

	// Get next byte
	if b, err = i.NextByte(); err != nil {
		err = fmt.Errorf("astits: fetching next byte failed: %w", err)
		return
	}

	// Items length
	itemsLength := int(b)

	// Items
	offsetEnd := i.Offset() + itemsLength
	for i.Offset() < offsetEnd {
		// Create item
		var item *DescriptorExtendedEventItem
		if item, err = newDescriptorExtendedEventItem(i); err != nil {
			err = fmt.Errorf("astits: creating extended event item failed: %w", err)
			return
		}

		// Append item
		d.Items = append(d.Items, item)
	}

	// Get next byte
	if b, err = i.NextByte(); err != nil {
		err = fmt.Errorf("astits: fetching next byte failed: %w", err)
		return
	}

	// Text length
	textLength := int(b)

	// Text
	if d.Text, err = i.NextBytes(textLength); err != nil {
		err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
		return
	}
	return
}

func newDescriptorExtendedEventItem(i *astikit.BytesIterator) (d *DescriptorExtendedEventItem, err error) {
	// Init
	d = &DescriptorExtendedEventItem{}

	// Get next byte
	var b byte
	if b, err = i.NextByte(); err != nil {
		err = fmt.Errorf("astits: fetching next byte failed: %w", err)
		return
	}

	// Description length
	descriptionLength := int(b)

	// Description
	if d.Description, err = i.NextBytes(descriptionLength); err != nil {
		err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
		return
	}

	// Get next byte
	if b, err = i.NextByte(); err != nil {
		err = fmt.Errorf("astits: fetching next byte failed: %w", err)
		return
	}

	// Content length
	contentLength := int(b)

	// Content
	if d.Content, err = i.NextBytes(contentLength); err != nil {
		err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
		return
	}
	return
}

// DescriptorExtension represents an extension descriptor
// Chapter: 6.2.16 | Link: https://www.etsi.org/deliver/etsi_en/300400_300499/300468/01.15.01_60/en_300468v011501p.pdf
type DescriptorExtension struct {
	SupplementaryAudio *DescriptorExtensionSupplementaryAudio
	Tag                uint8
	Unknown            *[]byte
}

func newDescriptorExtension(i *astikit.BytesIterator, offsetEnd int) (d *DescriptorExtension, err error) {
	// Get next byte
	var b byte
	if b, err = i.NextByte(); err != nil {
		err = fmt.Errorf("astits: fetching next byte failed: %w", err)
		return
	}

	// Create descriptor
	d = &DescriptorExtension{Tag: uint8(b)}

	// Switch on tag
	switch d.Tag {
	case DescriptorTagExtensionSupplementaryAudio:
		if d.SupplementaryAudio, err = newDescriptorExtensionSupplementaryAudio(i, offsetEnd); err != nil {
			err = fmt.Errorf("astits: parsing extension supplementary audio descriptor failed: %w", err)
			return
		}
	default:
		// Get next bytes
		var b []byte
		if b, err = i.NextBytes(offsetEnd - i.Offset()); err != nil {
			err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
			return
		}

		// Update unknown
		d.Unknown = &b
	}
	return
}

// DescriptorExtensionSupplementaryAudio represents a supplementary audio extension descriptor
// Chapter: 6.4.10 | Link: https://www.etsi.org/deliver/etsi_en/300400_300499/300468/01.15.01_60/en_300468v011501p.pdf
type DescriptorExtensionSupplementaryAudio struct {
	EditorialClassification uint8
	HasLanguageCode         bool
	LanguageCode            []byte
	MixType                 bool
	PrivateData             []byte
}

func newDescriptorExtensionSupplementaryAudio(i *astikit.BytesIterator, offsetEnd int) (d *DescriptorExtensionSupplementaryAudio, err error) {
	// Get next byte
	var b byte
	if b, err = i.NextByte(); err != nil {
		err = fmt.Errorf("astits: fetching next byte failed: %w", err)
		return
	}

	// Init
	d = &DescriptorExtensionSupplementaryAudio{
		EditorialClassification: uint8(b >> 2 & 0x1f),
		HasLanguageCode:         b&0x1 > 0,
		MixType:                 b&0x80 > 0,
	}

	// Language code
	if d.HasLanguageCode {
		if d.LanguageCode, err = i.NextBytes(3); err != nil {
			err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
			return
		}
	}

	// Private data
	if i.Offset() < offsetEnd {
		if d.PrivateData, err = i.NextBytes(offsetEnd - i.Offset()); err != nil {
			err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
			return
		}
	}
	return
}

// DescriptorISO639LanguageAndAudioType represents an ISO639 language descriptor
// https://github.com/gfto/bitstream/blob/master/mpeg/psi/desc_0a.h
// FIXME (barbashov) according to Chapter 2.6.18 ISO/IEC 13818-1:2015 there could be not one, but multiple such descriptors
type DescriptorISO639LanguageAndAudioType struct {
	Language []byte
	Type     uint8
}

// In some actual cases, the length is 3 and the language is described in only 2 bytes
func newDescriptorISO639LanguageAndAudioType(i *astikit.BytesIterator, offsetEnd int) (d *DescriptorISO639LanguageAndAudioType, err error) {
	// Get next bytes
	var bs []byte
	if bs, err = i.NextBytes(offsetEnd - i.Offset()); err != nil {
		err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
		return
	}

	// Create descriptor
	d = &DescriptorISO639LanguageAndAudioType{
		Language: bs[0 : len(bs)-1],
		Type:     uint8(bs[len(bs)-1]),
	}
	return
}

// DescriptorLocalTimeOffset represents a local time offset descriptor
// Chapter: 6.2.20 | Link: https://www.etsi.org/deliver/etsi_en/300400_300499/300468/01.15.01_60/en_300468v011501p.pdf
type DescriptorLocalTimeOffset struct {
	Items []*DescriptorLocalTimeOffsetItem
}

// DescriptorLocalTimeOffsetItem represents a local time offset item descriptor
// Chapter: 6.2.20 | Link: https://www.etsi.org/deliver/etsi_en/300400_300499/300468/01.15.01_60/en_300468v011501p.pdf
type DescriptorLocalTimeOffsetItem struct {
	CountryCode             []byte
	CountryRegionID         uint8
	LocalTimeOffset         time.Duration
	LocalTimeOffsetPolarity bool
	NextTimeOffset          time.Duration
	TimeOfChange            time.Time
}

func newDescriptorLocalTimeOffset(i *astikit.BytesIterator, offsetEnd int) (d *DescriptorLocalTimeOffset, err error) {
	// Init
	d = &DescriptorLocalTimeOffset{}

	// Add items
	for i.Offset() < offsetEnd {
		// Create item
		itm := &DescriptorLocalTimeOffsetItem{}

		// Country code
		if itm.CountryCode, err = i.NextBytes(3); err != nil {
			err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
			return
		}

		// Get next byte
		var b byte
		if b, err = i.NextByte(); err != nil {
			err = fmt.Errorf("astits: fetching next byte failed: %w", err)
			return
		}

		// Country region ID
		itm.CountryRegionID = uint8(b >> 2)

		// Local time offset polarity
		itm.LocalTimeOffsetPolarity = b&0x1 > 0

		// Local time offset
		if itm.LocalTimeOffset, err = parseDVBDurationMinutes(i); err != nil {
			err = fmt.Errorf("astits: parsing DVB durationminutes failed: %w", err)
			return
		}

		// Time of change
		if itm.TimeOfChange, err = parseDVBTime(i); err != nil {
			err = fmt.Errorf("astits: parsing DVB time failed: %w", err)
			return
		}

		// Next time offset
		if itm.NextTimeOffset, err = parseDVBDurationMinutes(i); err != nil {
			err = fmt.Errorf("astits: parsing DVB duration minutes failed: %w", err)
			return
		}

		// Append item
		d.Items = append(d.Items, itm)
	}
	return
}

// DescriptorMaximumBitrate represents a maximum bitrate descriptor
type DescriptorMaximumBitrate struct {
	Bitrate uint32 // In bytes/second
}

func newDescriptorMaximumBitrate(i *astikit.BytesIterator) (d *DescriptorMaximumBitrate, err error) {
	// Get next bytes
	var bs []byte
	if bs, err = i.NextBytesNoCopy(3); err != nil {
		err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
		return
	}

	// Create descriptor
	d = &DescriptorMaximumBitrate{Bitrate: (uint32(bs[0]&0x3f)<<16 | uint32(bs[1])<<8 | uint32(bs[2])) * 50}
	return
}

// DescriptorNetworkName represents a network name descriptor
// Chapter: 6.2.27 | Link: https://www.etsi.org/deliver/etsi_en/300400_300499/300468/01.15.01_60/en_300468v011501p.pdf
type DescriptorNetworkName struct {
	Name []byte
}

func newDescriptorNetworkName(i *astikit.BytesIterator, offsetEnd int) (d *DescriptorNetworkName, err error) {
	// Create descriptor
	d = &DescriptorNetworkName{}

	// Name
	if d.Name, err = i.NextBytes(offsetEnd - i.Offset()); err != nil {
		err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
		return
	}
	return
}

// DescriptorParentalRating represents a parental rating descriptor
// Chapter: 6.2.28 | Link: https://www.etsi.org/deliver/etsi_en/300400_300499/300468/01.15.01_60/en_300468v011501p.pdf
type DescriptorParentalRating struct {
	Items []*DescriptorParentalRatingItem
}

// DescriptorParentalRatingItem represents a parental rating item descriptor
// Chapter: 6.2.28 | Link: https://www.etsi.org/deliver/etsi_en/300400_300499/300468/01.15.01_60/en_300468v011501p.pdf
type DescriptorParentalRatingItem struct {
	CountryCode []byte
	Rating      uint8
}

// MinimumAge returns the minimum age for the parental rating
func (d DescriptorParentalRatingItem) MinimumAge() int {
	// Undefined or user defined ratings
	if d.Rating == 0 || d.Rating > 0x10 {
		return 0
	}
	return int(d.Rating) + 3
}

func newDescriptorParentalRating(i *astikit.BytesIterator, offsetEnd int) (d *DescriptorParentalRating, err error) {
	// Create descriptor
	d = &DescriptorParentalRating{}

	// Add items
	for i.Offset() < offsetEnd {
		// Get next bytes
		var bs []byte
		if bs, err = i.NextBytes(4); err != nil {
			err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
			return
		}

		// Append item
		d.Items = append(d.Items, &DescriptorParentalRatingItem{
			CountryCode: bs[:3],
			Rating:      uint8(bs[3]),
		})
	}
	return
}

// DescriptorPrivateDataIndicator represents a private data Indicator descriptor
type DescriptorPrivateDataIndicator struct {
	Indicator uint32
}

func newDescriptorPrivateDataIndicator(i *astikit.BytesIterator) (d *DescriptorPrivateDataIndicator, err error) {
	// Get next bytes
	var bs []byte
	if bs, err = i.NextBytesNoCopy(4); err != nil {
		err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
		return
	}

	// Create descriptor
	d = &DescriptorPrivateDataIndicator{Indicator: uint32(bs[0])<<24 | uint32(bs[1])<<16 | uint32(bs[2])<<8 | uint32(bs[3])}
	return
}

// DescriptorPrivateDataSpecifier represents a private data specifier descriptor
type DescriptorPrivateDataSpecifier struct {
	Specifier uint32
}

func newDescriptorPrivateDataSpecifier(i *astikit.BytesIterator) (d *DescriptorPrivateDataSpecifier, err error) {
	// Get next bytes
	var bs []byte
	if bs, err = i.NextBytesNoCopy(4); err != nil {
		err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
		return
	}

	// Create descriptor
	d = &DescriptorPrivateDataSpecifier{Specifier: uint32(bs[0])<<24 | uint32(bs[1])<<16 | uint32(bs[2])<<8 | uint32(bs[3])}
	return
}

// DescriptorRegistration represents a registration descriptor
// Page: 84 | http://ecee.colorado.edu/~ecen5653/ecen5653/papers/iso13818-1.pdf
type DescriptorRegistration struct {
	AdditionalIdentificationInfo []byte
	FormatIdentifier             uint32
}

func newDescriptorRegistration(i *astikit.BytesIterator, offsetEnd int) (d *DescriptorRegistration, err error) {
	// Get next bytes
	var bs []byte
	if bs, err = i.NextBytesNoCopy(4); err != nil {
		err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
		return
	}

	// Create descriptor
	d = &DescriptorRegistration{FormatIdentifier: uint32(bs[0])<<24 | uint32(bs[1])<<16 | uint32(bs[2])<<8 | uint32(bs[3])}

	// Additional identification info
	if i.Offset() < offsetEnd {
		if d.AdditionalIdentificationInfo, err = i.NextBytes(offsetEnd - i.Offset()); err != nil {
			err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
			return
		}
	}
	return
}

// DescriptorService represents a service descriptor
// Chapter: 6.2.33 | Link: https://www.etsi.org/deliver/etsi_en/300400_300499/300468/01.15.01_60/en_300468v011501p.pdf
type DescriptorService struct {
	Name     []byte
	Provider []byte
	Type     uint8
}

func newDescriptorService(i *astikit.BytesIterator) (d *DescriptorService, err error) {
	// Get next byte
	var b byte
	if b, err = i.NextByte(); err != nil {
		err = fmt.Errorf("astits: fetching next byte failed: %w", err)
		return
	}

	// Create descriptor
	d = &DescriptorService{Type: uint8(b)}

	// Get next byte
	if b, err = i.NextByte(); err != nil {
		err = fmt.Errorf("astits: fetching next byte failed: %w", err)
		return
	}

	// Provider length
	providerLength := int(b)

	// Provider
	if d.Provider, err = i.NextBytes(providerLength); err != nil {
		err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
		return
	}

	// Get next byte
	if b, err = i.NextByte(); err != nil {
		err = fmt.Errorf("astits: fetching next byte failed: %w", err)
		return
	}

	// Name length
	nameLength := int(b)

	// Name
	if d.Name, err = i.NextBytes(nameLength); err != nil {
		err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
		return
	}
	return
}

// DescriptorShortEvent represents a short event descriptor
// Chapter: 6.2.37 | Link: https://www.etsi.org/deliver/etsi_en/300400_300499/300468/01.15.01_60/en_300468v011501p.pdf
type DescriptorShortEvent struct {
	EventName []byte
	Language  []byte
	Text      []byte
}

func newDescriptorShortEvent(i *astikit.BytesIterator) (d *DescriptorShortEvent, err error) {
	// Create descriptor
	d = &DescriptorShortEvent{}

	// Language
	if d.Language, err = i.NextBytes(3); err != nil {
		err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
		return
	}

	// Get next byte
	var b byte
	if b, err = i.NextByte(); err != nil {
		err = fmt.Errorf("astits: fetching next byte failed: %w", err)
		return
	}

	// Event length
	eventLength := int(b)

	// Event name
	if d.EventName, err = i.NextBytes(eventLength); err != nil {
		err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
		return
	}

	// Get next byte
	if b, err = i.NextByte(); err != nil {
		err = fmt.Errorf("astits: fetching next byte failed: %w", err)
		return
	}

	// Text length
	textLength := int(b)

	// Text
	if d.Text, err = i.NextBytes(textLength); err != nil {
		err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
		return
	}
	return
}

// DescriptorStreamIdentifier represents a stream identifier descriptor
// Chapter: 6.2.39 | Link: https://www.etsi.org/deliver/etsi_en/300400_300499/300468/01.15.01_60/en_300468v011501p.pdf
type DescriptorStreamIdentifier struct {
	ComponentTag uint8
}

func newDescriptorStreamIdentifier(i *astikit.BytesIterator) (d *DescriptorStreamIdentifier, err error) {
	var b byte
	if b, err = i.NextByte(); err != nil {
		err = fmt.Errorf("astits: fetching next byte failed: %w", err)
		return
	}
	d = &DescriptorStreamIdentifier{ComponentTag: uint8(b)}
	return
}

// DescriptorSubtitling represents a subtitling descriptor
// Chapter: 6.2.41 | Link: https://www.etsi.org/deliver/etsi_en/300400_300499/300468/01.15.01_60/en_300468v011501p.pdf
type DescriptorSubtitling struct {
	Items []*DescriptorSubtitlingItem
}

// DescriptorSubtitlingItem represents subtitling descriptor item
// Chapter: 6.2.41 | Link: https://www.etsi.org/deliver/etsi_en/300400_300499/300468/01.15.01_60/en_300468v011501p.pdf
type DescriptorSubtitlingItem struct {
	AncillaryPageID   uint16
	CompositionPageID uint16
	Language          []byte
	Type              uint8
}

func newDescriptorSubtitling(i *astikit.BytesIterator, offsetEnd int) (d *DescriptorSubtitling, err error) {
	// Create descriptor
	d = &DescriptorSubtitling{}

	// Loop
	for i.Offset() < offsetEnd {
		// Create item
		itm := &DescriptorSubtitlingItem{}

		// Language
		if itm.Language, err = i.NextBytes(3); err != nil {
			err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
			return
		}

		// Get next byte
		var b byte
		if b, err = i.NextByte(); err != nil {
			err = fmt.Errorf("astits: fetching next byte failed: %w", err)
			return
		}

		// Type
		itm.Type = uint8(b)

		// Get next bytes
		var bs []byte
		if bs, err = i.NextBytesNoCopy(2); err != nil {
			err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
			return
		}

		// Composition page ID
		itm.CompositionPageID = uint16(bs[0])<<8 | uint16(bs[1])

		// Get next bytes
		if bs, err = i.NextBytesNoCopy(2); err != nil {
			err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
			return
		}

		// Ancillary page ID
		itm.AncillaryPageID = uint16(bs[0])<<8 | uint16(bs[1])

		// Append item
		d.Items = append(d.Items, itm)
	}
	return
}

// DescriptorTeletext represents a teletext descriptor
// Chapter: 6.2.43 | Link: https://www.etsi.org/deliver/etsi_en/300400_300499/300468/01.15.01_60/en_300468v011501p.pdf
type DescriptorTeletext struct {
	Items []*DescriptorTeletextItem
}

// DescriptorTeletextItem represents a teletext descriptor item
// Chapter: 6.2.43 | Link: https://www.etsi.org/deliver/etsi_en/300400_300499/300468/01.15.01_60/en_300468v011501p.pdf
type DescriptorTeletextItem struct {
	Language []byte
	Magazine uint8
	Page     uint8
	Type     uint8
}

func newDescriptorTeletext(i *astikit.BytesIterator, offsetEnd int) (d *DescriptorTeletext, err error) {
	// Create descriptor
	d = &DescriptorTeletext{}

	// Loop
	for i.Offset() < offsetEnd {
		// Create item
		itm := &DescriptorTeletextItem{}

		// Language
		if itm.Language, err = i.NextBytes(3); err != nil {
			err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
			return
		}

		// Get next byte
		var b byte
		if b, err = i.NextByte(); err != nil {
			err = fmt.Errorf("astits: fetching next byte failed: %w", err)
			return
		}

		// Type
		itm.Type = uint8(b) >> 3

		// Magazine
		itm.Magazine = uint8(b & 0x7)

		// Get next byte
		if b, err = i.NextByte(); err != nil {
			err = fmt.Errorf("astits: fetching next byte failed: %w", err)
			return
		}

		// Page
		itm.Page = uint8(b)>>4*10 + uint8(b&0xf)

		// Append item
		d.Items = append(d.Items, itm)
	}
	return
}

type DescriptorUnknown struct {
	Content []byte
	Tag     uint8
}

func newDescriptorUnknown(i *astikit.BytesIterator, tag, length uint8) (d *DescriptorUnknown, err error) {
	// Create descriptor
	d = &DescriptorUnknown{Tag: tag}

	// Get next bytes
	if d.Content, err = i.NextBytes(int(length)); err != nil {
		err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
		return
	}
	return
}

// DescriptorVBIData represents a VBI data descriptor
// Chapter: 6.2.47 | Link: https://www.etsi.org/deliver/etsi_en/300400_300499/300468/01.15.01_60/en_300468v011501p.pdf
type DescriptorVBIData struct {
	Services []*DescriptorVBIDataService
}

// DescriptorVBIDataService represents a vbi data service descriptor
// Chapter: 6.2.47 | Link: https://www.etsi.org/deliver/etsi_en/300400_300499/300468/01.15.01_60/en_300468v011501p.pdf
type DescriptorVBIDataService struct {
	DataServiceID uint8
	Descriptors   []*DescriptorVBIDataDescriptor
}

// DescriptorVBIDataItem represents a vbi data descriptor item
// Chapter: 6.2.47 | Link: https://www.etsi.org/deliver/etsi_en/300400_300499/300468/01.15.01_60/en_300468v011501p.pdf
type DescriptorVBIDataDescriptor struct {
	FieldParity bool
	LineOffset  uint8
}

func newDescriptorVBIData(i *astikit.BytesIterator, offsetEnd int) (d *DescriptorVBIData, err error) {
	// Create descriptor
	d = &DescriptorVBIData{}

	// Loop
	for i.Offset() < offsetEnd {
		// Create service
		srv := &DescriptorVBIDataService{}

		// Get next byte
		var b byte
		if b, err = i.NextByte(); err != nil {
			err = fmt.Errorf("astits: fetching next byte failed: %w", err)
			return
		}

		// Data service ID
		srv.DataServiceID = uint8(b)

		// Get next byte
		if b, err = i.NextByte(); err != nil {
			err = fmt.Errorf("astits: fetching next byte failed: %w", err)
			return
		}

		// Data service descriptor length
		dataServiceDescriptorLength := int(b)

		// Data service descriptor
		offsetDataEnd := i.Offset() + dataServiceDescriptorLength
		for i.Offset() < offsetDataEnd {
			if srv.DataServiceID == VBIDataServiceIDClosedCaptioning ||
				srv.DataServiceID == VBIDataServiceIDEBUTeletext ||
				srv.DataServiceID == VBIDataServiceIDInvertedTeletext ||
				srv.DataServiceID == VBIDataServiceIDMonochrome442Samples ||
				srv.DataServiceID == VBIDataServiceIDVPS ||
				srv.DataServiceID == VBIDataServiceIDWSS {
				// Get next byte
				if b, err = i.NextByte(); err != nil {
					err = fmt.Errorf("astits: fetching next byte failed: %w", err)
					return
				}

				// Append data
				srv.Descriptors = append(srv.Descriptors, &DescriptorVBIDataDescriptor{
					FieldParity: b&0x20 > 0,
					LineOffset:  uint8(b & 0x1f),
				})
			}
		}

		// Append service
		d.Services = append(d.Services, srv)
	}
	return
}

// parseDescriptors parses descriptors
func parseDescriptors(i *astikit.BytesIterator) (o []*Descriptor, err error) {
	// Get next 2 bytes
	var bs []byte
	if bs, err = i.NextBytesNoCopy(2); err != nil {
		err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
		return
	}

	// Get length
	length := int(uint16(bs[0]&0xf)<<8 | uint16(bs[1]))

	// Loop
	if length > 0 {
		offsetEnd := i.Offset() + length
		for i.Offset() < offsetEnd {
			// Get next 2 bytes
			if bs, err = i.NextBytesNoCopy(2); err != nil {
				err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
				return
			}

			// Create descriptor
			d := &Descriptor{
				Length: uint8(bs[1]),
				Tag:    uint8(bs[0]),
			}

			// Parse data
			if d.Length > 0 {
				// Unfortunately there's no way to be sure the real descriptor length is the same as the one indicated
				// previously therefore we must fetch bytes in descriptor functions and seek at the end
				offsetDescriptorEnd := i.Offset() + int(d.Length)

				// User defined
				if d.Tag >= 0x80 && d.Tag <= 0xfe {
					// Get next bytes
					if d.UserDefined, err = i.NextBytes(int(d.Length)); err != nil {
						err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
						return
					}
				} else {
					// Switch on tag
					switch d.Tag {
					case DescriptorTagAC3:
						if d.AC3, err = newDescriptorAC3(i, offsetDescriptorEnd); err != nil {
							err = fmt.Errorf("astits: parsing AC3 descriptor failed: %w", err)
							return
						}
					case DescriptorTagAVCVideo:
						if d.AVCVideo, err = newDescriptorAVCVideo(i); err != nil {
							err = fmt.Errorf("astits: parsing AVC Video descriptor failed: %w", err)
							return
						}
					case DescriptorTagComponent:
						if d.Component, err = newDescriptorComponent(i, offsetDescriptorEnd); err != nil {
							err = fmt.Errorf("astits: parsing Component descriptor failed: %w", err)
							return
						}
					case DescriptorTagContent:
						if d.Content, err = newDescriptorContent(i, offsetDescriptorEnd); err != nil {
							err = fmt.Errorf("astits: parsing Content descriptor failed: %w", err)
							return
						}
					case DescriptorTagDataStreamAlignment:
						if d.DataStreamAlignment, err = newDescriptorDataStreamAlignment(i); err != nil {
							err = fmt.Errorf("astits: parsing Data Stream Alignment descriptor failed: %w", err)
							return
						}
					case DescriptorTagEnhancedAC3:
						if d.EnhancedAC3, err = newDescriptorEnhancedAC3(i, offsetDescriptorEnd); err != nil {
							err = fmt.Errorf("astits: parsing Enhanced AC3 descriptor failed: %w", err)
							return
						}
					case DescriptorTagExtendedEvent:
						if d.ExtendedEvent, err = newDescriptorExtendedEvent(i); err != nil {
							err = fmt.Errorf("astits: parsing Extended event descriptor failed: %w", err)
							return
						}
					case DescriptorTagExtension:
						if d.Extension, err = newDescriptorExtension(i, offsetDescriptorEnd); err != nil {
							err = fmt.Errorf("astits: parsing Extension descriptor failed: %w", err)
							return
						}
					case DescriptorTagISO639LanguageAndAudioType:
						if d.ISO639LanguageAndAudioType, err = newDescriptorISO639LanguageAndAudioType(i, offsetDescriptorEnd); err != nil {
							err = fmt.Errorf("astits: parsing ISO639 Language and Audio Type descriptor failed: %w", err)
							return
						}
					case DescriptorTagLocalTimeOffset:
						if d.LocalTimeOffset, err = newDescriptorLocalTimeOffset(i, offsetDescriptorEnd); err != nil {
							err = fmt.Errorf("astits: parsing Local Time Offset descriptor failed: %w", err)
							return
						}
					case DescriptorTagMaximumBitrate:
						if d.MaximumBitrate, err = newDescriptorMaximumBitrate(i); err != nil {
							err = fmt.Errorf("astits: parsing Maximum Bitrate descriptor failed: %w", err)
							return
						}
					case DescriptorTagNetworkName:
						if d.NetworkName, err = newDescriptorNetworkName(i, offsetDescriptorEnd); err != nil {
							err = fmt.Errorf("astits: parsing Network Name descriptor failed: %w", err)
							return
						}
					case DescriptorTagParentalRating:
						if d.ParentalRating, err = newDescriptorParentalRating(i, offsetDescriptorEnd); err != nil {
							err = fmt.Errorf("astits: parsing Parental Rating descriptor failed: %w", err)
							return
						}
					case DescriptorTagPrivateDataIndicator:
						if d.PrivateDataIndicator, err = newDescriptorPrivateDataIndicator(i); err != nil {
							err = fmt.Errorf("astits: parsing Private Data Indicator descriptor failed: %w", err)
							return
						}
					case DescriptorTagPrivateDataSpecifier:
						if d.PrivateDataSpecifier, err = newDescriptorPrivateDataSpecifier(i); err != nil {
							err = fmt.Errorf("astits: parsing Private Data Specifier descriptor failed: %w", err)
							return
						}
					case DescriptorTagRegistration:
						if d.Registration, err = newDescriptorRegistration(i, offsetDescriptorEnd); err != nil {
							err = fmt.Errorf("astits: parsing Registration descriptor failed: %w", err)
							return
						}
					case DescriptorTagService:
						if d.Service, err = newDescriptorService(i); err != nil {
							err = fmt.Errorf("astits: parsing Service descriptor failed: %w", err)
							return
						}
					case DescriptorTagShortEvent:
						if d.ShortEvent, err = newDescriptorShortEvent(i); err != nil {
							err = fmt.Errorf("astits: parsing Short Event descriptor failed: %w", err)
							return
						}
					case DescriptorTagStreamIdentifier:
						if d.StreamIdentifier, err = newDescriptorStreamIdentifier(i); err != nil {
							err = fmt.Errorf("astits: parsing Stream Identifier descriptor failed: %w", err)
							return
						}
					case DescriptorTagSubtitling:
						if d.Subtitling, err = newDescriptorSubtitling(i, offsetDescriptorEnd); err != nil {
							err = fmt.Errorf("astits: parsing Subtitling descriptor failed: %w", err)
							return
						}
					case DescriptorTagTeletext:
						if d.Teletext, err = newDescriptorTeletext(i, offsetDescriptorEnd); err != nil {
							err = fmt.Errorf("astits: parsing Teletext descriptor failed: %w", err)
							return
						}
					case DescriptorTagVBIData:
						if d.VBIData, err = newDescriptorVBIData(i, offsetDescriptorEnd); err != nil {
							err = fmt.Errorf("astits: parsing VBI Date descriptor failed: %w", err)
							return
						}
					case DescriptorTagVBITeletext:
						if d.VBITeletext, err = newDescriptorTeletext(i, offsetDescriptorEnd); err != nil {
							err = fmt.Errorf("astits: parsing VBI Teletext descriptor failed: %w", err)
							return
						}
					default:
						if d.Unknown, err = newDescriptorUnknown(i, d.Tag, d.Length); err != nil {
							err = fmt.Errorf("astits: parsing unknown descriptor failed: %w", err)
							return
						}
					}
				}

				// Seek in iterator to make sure we move to the end of the descriptor since its content may be
				// corrupted
				i.Seek(offsetDescriptorEnd)
			}
			o = append(o, d)
		}
	}
	return
}

func calcDescriptorUserDefinedLength(d []byte) uint8 {
	return uint8(len(d))
}

func writeDescriptorUserDefined(w *astikit.BitsWriter, d []byte) error {
	b := astikit.NewBitsWriterBatch(w)

	b.Write(d)

	return b.Err()
}

func calcDescriptorAC3Length(d *DescriptorAC3) uint8 {
	ret := 1 // flags

	if d.HasComponentType {
		ret++
	}
	if d.HasBSID {
		ret++
	}
	if d.HasMainID {
		ret++
	}
	if d.HasASVC {
		ret++
	}

	ret += len(d.AdditionalInfo)

	return uint8(ret)
}

func writeDescriptorAC3(w *astikit.BitsWriter, d *DescriptorAC3) error {
	b := astikit.NewBitsWriterBatch(w)

	b.Write(d.HasComponentType)
	b.Write(d.HasBSID)
	b.Write(d.HasMainID)
	b.Write(d.HasASVC)
	b.WriteN(uint8(0xff), 4)

	if d.HasComponentType {
		b.Write(d.ComponentType)
	}
	if d.HasBSID {
		b.Write(d.BSID)
	}
	if d.HasMainID {
		b.Write(d.MainID)
	}
	if d.HasASVC {
		b.Write(d.ASVC)
	}
	b.Write(d.AdditionalInfo)

	return b.Err()
}

func calcDescriptorAVCVideoLength(d *DescriptorAVCVideo) uint8 {
	return 4
}

func writeDescriptorAVCVideo(w *astikit.BitsWriter, d *DescriptorAVCVideo) error {
	b := astikit.NewBitsWriterBatch(w)

	b.Write(d.ProfileIDC)

	b.Write(d.ConstraintSet0Flag)
	b.Write(d.ConstraintSet1Flag)
	b.Write(d.ConstraintSet2Flag)
	b.WriteN(d.CompatibleFlags, 5)

	b.Write(d.LevelIDC)

	b.Write(d.AVCStillPresent)
	b.Write(d.AVC24HourPictureFlag)
	b.WriteN(uint8(0xff), 6)

	return b.Err()
}

func calcDescriptorComponentLength(d *DescriptorComponent) uint8 {
	return uint8(6 + len(d.Text))
}

func writeDescriptorComponent(w *astikit.BitsWriter, d *DescriptorComponent) error {
	b := astikit.NewBitsWriterBatch(w)

	b.WriteN(d.StreamContentExt, 4)
	b.WriteN(d.StreamContent, 4)

	b.Write(d.ComponentType)
	b.Write(d.ComponentTag)

	b.WriteBytesN(d.ISO639LanguageCode, 3, 0)

	b.Write(d.Text)

	return b.Err()
}

func calcDescriptorContentLength(d *DescriptorContent) uint8 {
	return uint8(2 * len(d.Items))
}

func writeDescriptorContent(w *astikit.BitsWriter, d *DescriptorContent) error {
	b := astikit.NewBitsWriterBatch(w)

	for _, item := range d.Items {
		b.WriteN(item.ContentNibbleLevel1, 4)
		b.WriteN(item.ContentNibbleLevel2, 4)
		b.Write(item.UserByte)
	}

	return b.Err()
}

func calcDescriptorDataStreamAlignmentLength(d *DescriptorDataStreamAlignment) uint8 {
	return 1
}

func writeDescriptorDataStreamAlignment(w *astikit.BitsWriter, d *DescriptorDataStreamAlignment) error {
	b := astikit.NewBitsWriterBatch(w)

	b.Write(d.Type)

	return b.Err()
}

func calcDescriptorEnhancedAC3Length(d *DescriptorEnhancedAC3) uint8 {
	ret := 1 // flags

	if d.HasComponentType {
		ret++
	}
	if d.HasBSID {
		ret++
	}
	if d.HasMainID {
		ret++
	}
	if d.HasASVC {
		ret++
	}
	if d.HasSubStream1 {
		ret++
	}
	if d.HasSubStream2 {
		ret++
	}
	if d.HasSubStream3 {
		ret++
	}

	ret += len(d.AdditionalInfo)

	return uint8(ret)
}

func writeDescriptorEnhancedAC3(w *astikit.BitsWriter, d *DescriptorEnhancedAC3) error {
	b := astikit.NewBitsWriterBatch(w)

	b.Write(d.HasComponentType)
	b.Write(d.HasBSID)
	b.Write(d.HasMainID)
	b.Write(d.HasASVC)
	b.Write(d.MixInfoExists)
	b.Write(d.HasSubStream1)
	b.Write(d.HasSubStream2)
	b.Write(d.HasSubStream3)

	if d.HasComponentType {
		b.Write(d.ComponentType)
	}
	if d.HasBSID {
		b.Write(d.BSID)
	}
	if d.HasMainID {
		b.Write(d.MainID)
	}
	if d.HasASVC {
		b.Write(d.ASVC)
	}
	if d.HasSubStream1 {
		b.Write(d.SubStream1)
	}
	if d.HasSubStream2 {
		b.Write(d.SubStream2)
	}
	if d.HasSubStream3 {
		b.Write(d.SubStream3)
	}

	b.Write(d.AdditionalInfo)

	return b.Err()
}

func calcDescriptorExtendedEventLength(d *DescriptorExtendedEvent) (descriptorLength, lengthOfItems uint8) {
	ret := 1 + 3 + 1 // numbers, language and items length

	itemsRet := 0
	for _, item := range d.Items {
		itemsRet += 1 // description length
		itemsRet += len(item.Description)
		itemsRet += 1 // content length
		itemsRet += len(item.Content)
	}

	ret += itemsRet

	ret += 1 // text length
	ret += len(d.Text)

	return uint8(ret), uint8(itemsRet)
}

func writeDescriptorExtendedEvent(w *astikit.BitsWriter, d *DescriptorExtendedEvent) error {
	b := astikit.NewBitsWriterBatch(w)

	var lengthOfItems uint8

	_, lengthOfItems = calcDescriptorExtendedEventLength(d)

	b.WriteN(d.Number, 4)
	b.WriteN(d.LastDescriptorNumber, 4)

	b.WriteBytesN(d.ISO639LanguageCode, 3, 0)

	b.Write(lengthOfItems)
	for _, item := range d.Items {
		b.Write(uint8(len(item.Description)))
		b.Write(item.Description)
		b.Write(uint8(len(item.Content)))
		b.Write(item.Content)
	}

	b.Write(uint8(len(d.Text)))
	b.Write(d.Text)

	return b.Err()
}

func calcDescriptorExtensionSupplementaryAudioLength(d *DescriptorExtensionSupplementaryAudio) int {
	ret := 1
	if d.HasLanguageCode {
		ret += 3
	}
	ret += len(d.PrivateData)
	return ret
}

func calcDescriptorExtensionLength(d *DescriptorExtension) uint8 {
	ret := 1 // tag

	switch d.Tag {
	case DescriptorTagExtensionSupplementaryAudio:
		ret += calcDescriptorExtensionSupplementaryAudioLength(d.SupplementaryAudio)
	default:
		if d.Unknown != nil {
			ret += len(*d.Unknown)
		}
	}

	return uint8(ret)
}

func writeDescriptorExtensionSupplementaryAudio(w *astikit.BitsWriter, d *DescriptorExtensionSupplementaryAudio) error {
	b := astikit.NewBitsWriterBatch(w)

	b.Write(d.MixType)
	b.WriteN(d.EditorialClassification, 5)
	b.Write(true) // reserved
	b.Write(d.HasLanguageCode)

	if d.HasLanguageCode {
		b.WriteBytesN(d.LanguageCode, 3, 0)
	}

	b.Write(d.PrivateData)

	return b.Err()
}

func writeDescriptorExtension(w *astikit.BitsWriter, d *DescriptorExtension) error {
	b := astikit.NewBitsWriterBatch(w)

	b.Write(d.Tag)

	switch d.Tag {
	case DescriptorTagExtensionSupplementaryAudio:
		err := writeDescriptorExtensionSupplementaryAudio(w, d.SupplementaryAudio)
		if err != nil {
			return err
		}
	default:
		if d.Unknown != nil {
			b.Write(*d.Unknown)
		}
	}

	return b.Err()
}

func calcDescriptorISO639LanguageAndAudioTypeLength(d *DescriptorISO639LanguageAndAudioType) uint8 {
	return 3 + 1 // language code + type
}

func writeDescriptorISO639LanguageAndAudioType(w *astikit.BitsWriter, d *DescriptorISO639LanguageAndAudioType) error {
	b := astikit.NewBitsWriterBatch(w)

	b.WriteBytesN(d.Language, 3, 0)
	b.Write(d.Type)

	return b.Err()
}

func calcDescriptorLocalTimeOffsetLength(d *DescriptorLocalTimeOffset) uint8 {
	return uint8(13 * len(d.Items))
}

func writeDescriptorLocalTimeOffset(w *astikit.BitsWriter, d *DescriptorLocalTimeOffset) error {
	b := astikit.NewBitsWriterBatch(w)

	for _, item := range d.Items {
		b.WriteBytesN(item.CountryCode, 3, 0)

		b.WriteN(item.CountryRegionID, 6)
		b.WriteN(uint8(0xff), 1)
		b.Write(item.LocalTimeOffsetPolarity)

		if _, err := writeDVBDurationMinutes(w, item.LocalTimeOffset); err != nil {
			return err
		}
		if _, err := writeDVBTime(w, item.TimeOfChange); err != nil {
			return err
		}
		if _, err := writeDVBDurationMinutes(w, item.NextTimeOffset); err != nil {
			return err
		}
	}

	return b.Err()
}

func calcDescriptorMaximumBitrateLength(d *DescriptorMaximumBitrate) uint8 {
	return 3
}

func writeDescriptorMaximumBitrate(w *astikit.BitsWriter, d *DescriptorMaximumBitrate) error {
	b := astikit.NewBitsWriterBatch(w)

	b.WriteN(uint8(0xff), 2)
	b.WriteN(uint32(d.Bitrate/50), 22)

	return b.Err()
}

func calcDescriptorNetworkNameLength(d *DescriptorNetworkName) uint8 {
	return uint8(len(d.Name))
}

func writeDescriptorNetworkName(w *astikit.BitsWriter, d *DescriptorNetworkName) error {
	b := astikit.NewBitsWriterBatch(w)

	b.Write(d.Name)

	return b.Err()
}

func calcDescriptorParentalRatingLength(d *DescriptorParentalRating) uint8 {
	return uint8(4 * len(d.Items))
}

func writeDescriptorParentalRating(w *astikit.BitsWriter, d *DescriptorParentalRating) error {
	b := astikit.NewBitsWriterBatch(w)

	for _, item := range d.Items {
		b.WriteBytesN(item.CountryCode, 3, 0)
		b.Write(item.Rating)
	}

	return b.Err()
}

func calcDescriptorPrivateDataIndicatorLength(d *DescriptorPrivateDataIndicator) uint8 {
	return 4
}

func writeDescriptorPrivateDataIndicator(w *astikit.BitsWriter, d *DescriptorPrivateDataIndicator) error {
	b := astikit.NewBitsWriterBatch(w)

	b.Write(d.Indicator)

	return b.Err()
}

func calcDescriptorPrivateDataSpecifierLength(d *DescriptorPrivateDataSpecifier) uint8 {
	return 4
}

func writeDescriptorPrivateDataSpecifier(w *astikit.BitsWriter, d *DescriptorPrivateDataSpecifier) error {
	b := astikit.NewBitsWriterBatch(w)

	b.Write(d.Specifier)

	return b.Err()
}

func calcDescriptorRegistrationLength(d *DescriptorRegistration) uint8 {
	return uint8(4 + len(d.AdditionalIdentificationInfo))
}

func writeDescriptorRegistration(w *astikit.BitsWriter, d *DescriptorRegistration) error {
	b := astikit.NewBitsWriterBatch(w)

	b.Write(d.FormatIdentifier)
	b.Write(d.AdditionalIdentificationInfo)

	return b.Err()
}

func calcDescriptorServiceLength(d *DescriptorService) uint8 {
	ret := 3 // type and lengths
	ret += len(d.Name)
	ret += len(d.Provider)
	return uint8(ret)
}

func writeDescriptorService(w *astikit.BitsWriter, d *DescriptorService) error {
	b := astikit.NewBitsWriterBatch(w)

	b.Write(d.Type)
	b.Write(uint8(len(d.Provider)))
	b.Write(d.Provider)
	b.Write(uint8(len(d.Name)))
	b.Write(d.Name)

	return b.Err()
}

func calcDescriptorShortEventLength(d *DescriptorShortEvent) uint8 {
	ret := 3 + 1 + 1 // language code and lengths
	ret += len(d.EventName)
	ret += len(d.Text)
	return uint8(ret)
}

func writeDescriptorShortEvent(w *astikit.BitsWriter, d *DescriptorShortEvent) error {
	b := astikit.NewBitsWriterBatch(w)

	b.WriteBytesN(d.Language, 3, 0)

	b.Write(uint8(len(d.EventName)))
	b.Write(d.EventName)

	b.Write(uint8(len(d.Text)))
	b.Write(d.Text)

	return b.Err()
}

func calcDescriptorStreamIdentifierLength(d *DescriptorStreamIdentifier) uint8 {
	return 1
}

func writeDescriptorStreamIdentifier(w *astikit.BitsWriter, d *DescriptorStreamIdentifier) error {
	b := astikit.NewBitsWriterBatch(w)

	b.Write(d.ComponentTag)

	return b.Err()
}

func calcDescriptorSubtitlingLength(d *DescriptorSubtitling) uint8 {
	return uint8(8 * len(d.Items))
}

func writeDescriptorSubtitling(w *astikit.BitsWriter, d *DescriptorSubtitling) error {
	b := astikit.NewBitsWriterBatch(w)

	for _, item := range d.Items {
		b.WriteBytesN(item.Language, 3, 0)
		b.Write(item.Type)
		b.Write(item.CompositionPageID)
		b.Write(item.AncillaryPageID)
	}

	return b.Err()
}

func calcDescriptorTeletextLength(d *DescriptorTeletext) uint8 {
	return uint8(5 * len(d.Items))
}

func writeDescriptorTeletext(w *astikit.BitsWriter, d *DescriptorTeletext) error {
	b := astikit.NewBitsWriterBatch(w)

	for _, item := range d.Items {
		b.WriteBytesN(item.Language, 3, 0)
		b.WriteN(item.Type, 5)
		b.WriteN(item.Magazine, 3)
		b.WriteN(item.Page/10, 4)
		b.WriteN(item.Page%10, 4)
	}

	return b.Err()
}

func calcDescriptorVBIDataLength(d *DescriptorVBIData) uint8 {
	return uint8(3 * len(d.Services))
}

func writeDescriptorVBIData(w *astikit.BitsWriter, d *DescriptorVBIData) error {
	b := astikit.NewBitsWriterBatch(w)

	for _, item := range d.Services {
		b.Write(item.DataServiceID)

		if item.DataServiceID == VBIDataServiceIDClosedCaptioning ||
			item.DataServiceID == VBIDataServiceIDEBUTeletext ||
			item.DataServiceID == VBIDataServiceIDInvertedTeletext ||
			item.DataServiceID == VBIDataServiceIDMonochrome442Samples ||
			item.DataServiceID == VBIDataServiceIDVPS ||
			item.DataServiceID == VBIDataServiceIDWSS {

			b.Write(uint8(len(item.Descriptors))) // each descriptor is 1 byte
			for _, desc := range item.Descriptors {
				b.WriteN(uint8(0xff), 2)
				b.Write(desc.FieldParity)
				b.WriteN(desc.LineOffset, 5)
			}
		} else {
			// let's put one reserved byte
			b.Write(uint8(1))
			b.Write(uint8(0xff))
		}
	}

	return b.Err()
}

func calcDescriptorUnknownLength(d *DescriptorUnknown) uint8 {
	return uint8(len(d.Content))
}

func writeDescriptorUnknown(w *astikit.BitsWriter, d *DescriptorUnknown) error {
	b := astikit.NewBitsWriterBatch(w)

	b.Write(d.Content)

	return b.Err()
}

func calcDescriptorLength(d *Descriptor) uint8 {
	if d.Tag >= 0x80 && d.Tag <= 0xfe {
		return calcDescriptorUserDefinedLength(d.UserDefined)
	}

	switch d.Tag {
	case DescriptorTagAC3:
		return calcDescriptorAC3Length(d.AC3)
	case DescriptorTagAVCVideo:
		return calcDescriptorAVCVideoLength(d.AVCVideo)
	case DescriptorTagComponent:
		return calcDescriptorComponentLength(d.Component)
	case DescriptorTagContent:
		return calcDescriptorContentLength(d.Content)
	case DescriptorTagDataStreamAlignment:
		return calcDescriptorDataStreamAlignmentLength(d.DataStreamAlignment)
	case DescriptorTagEnhancedAC3:
		return calcDescriptorEnhancedAC3Length(d.EnhancedAC3)
	case DescriptorTagExtendedEvent:
		ret, _ := calcDescriptorExtendedEventLength(d.ExtendedEvent)
		return ret
	case DescriptorTagExtension:
		return calcDescriptorExtensionLength(d.Extension)
	case DescriptorTagISO639LanguageAndAudioType:
		return calcDescriptorISO639LanguageAndAudioTypeLength(d.ISO639LanguageAndAudioType)
	case DescriptorTagLocalTimeOffset:
		return calcDescriptorLocalTimeOffsetLength(d.LocalTimeOffset)
	case DescriptorTagMaximumBitrate:
		return calcDescriptorMaximumBitrateLength(d.MaximumBitrate)
	case DescriptorTagNetworkName:
		return calcDescriptorNetworkNameLength(d.NetworkName)
	case DescriptorTagParentalRating:
		return calcDescriptorParentalRatingLength(d.ParentalRating)
	case DescriptorTagPrivateDataIndicator:
		return calcDescriptorPrivateDataIndicatorLength(d.PrivateDataIndicator)
	case DescriptorTagPrivateDataSpecifier:
		return calcDescriptorPrivateDataSpecifierLength(d.PrivateDataSpecifier)
	case DescriptorTagRegistration:
		return calcDescriptorRegistrationLength(d.Registration)
	case DescriptorTagService:
		return calcDescriptorServiceLength(d.Service)
	case DescriptorTagShortEvent:
		return calcDescriptorShortEventLength(d.ShortEvent)
	case DescriptorTagStreamIdentifier:
		return calcDescriptorStreamIdentifierLength(d.StreamIdentifier)
	case DescriptorTagSubtitling:
		return calcDescriptorSubtitlingLength(d.Subtitling)
	case DescriptorTagTeletext:
		return calcDescriptorTeletextLength(d.Teletext)
	case DescriptorTagVBIData:
		return calcDescriptorVBIDataLength(d.VBIData)
	case DescriptorTagVBITeletext:
		return calcDescriptorTeletextLength(d.VBITeletext)
	}

	return calcDescriptorUnknownLength(d.Unknown)
}

func writeDescriptor(w *astikit.BitsWriter, d *Descriptor) (int, error) {
	b := astikit.NewBitsWriterBatch(w)
	length := calcDescriptorLength(d)

	b.Write(d.Tag)
	b.Write(length)

	if err := b.Err(); err != nil {
		return 0, err
	}

	written := int(length) + 2

	if d.Tag >= 0x80 && d.Tag <= 0xfe {
		return written, writeDescriptorUserDefined(w, d.UserDefined)
	}

	switch d.Tag {
	case DescriptorTagAC3:
		return written, writeDescriptorAC3(w, d.AC3)
	case DescriptorTagAVCVideo:
		return written, writeDescriptorAVCVideo(w, d.AVCVideo)
	case DescriptorTagComponent:
		return written, writeDescriptorComponent(w, d.Component)
	case DescriptorTagContent:
		return written, writeDescriptorContent(w, d.Content)
	case DescriptorTagDataStreamAlignment:
		return written, writeDescriptorDataStreamAlignment(w, d.DataStreamAlignment)
	case DescriptorTagEnhancedAC3:
		return written, writeDescriptorEnhancedAC3(w, d.EnhancedAC3)
	case DescriptorTagExtendedEvent:
		return written, writeDescriptorExtendedEvent(w, d.ExtendedEvent)
	case DescriptorTagExtension:
		return written, writeDescriptorExtension(w, d.Extension)
	case DescriptorTagISO639LanguageAndAudioType:
		return written, writeDescriptorISO639LanguageAndAudioType(w, d.ISO639LanguageAndAudioType)
	case DescriptorTagLocalTimeOffset:
		return written, writeDescriptorLocalTimeOffset(w, d.LocalTimeOffset)
	case DescriptorTagMaximumBitrate:
		return written, writeDescriptorMaximumBitrate(w, d.MaximumBitrate)
	case DescriptorTagNetworkName:
		return written, writeDescriptorNetworkName(w, d.NetworkName)
	case DescriptorTagParentalRating:
		return written, writeDescriptorParentalRating(w, d.ParentalRating)
	case DescriptorTagPrivateDataIndicator:
		return written, writeDescriptorPrivateDataIndicator(w, d.PrivateDataIndicator)
	case DescriptorTagPrivateDataSpecifier:
		return written, writeDescriptorPrivateDataSpecifier(w, d.PrivateDataSpecifier)
	case DescriptorTagRegistration:
		return written, writeDescriptorRegistration(w, d.Registration)
	case DescriptorTagService:
		return written, writeDescriptorService(w, d.Service)
	case DescriptorTagShortEvent:
		return written, writeDescriptorShortEvent(w, d.ShortEvent)
	case DescriptorTagStreamIdentifier:
		return written, writeDescriptorStreamIdentifier(w, d.StreamIdentifier)
	case DescriptorTagSubtitling:
		return written, writeDescriptorSubtitling(w, d.Subtitling)
	case DescriptorTagTeletext:
		return written, writeDescriptorTeletext(w, d.Teletext)
	case DescriptorTagVBIData:
		return written, writeDescriptorVBIData(w, d.VBIData)
	case DescriptorTagVBITeletext:
		return written, writeDescriptorTeletext(w, d.VBITeletext)
	}

	return written, writeDescriptorUnknown(w, d.Unknown)
}

func calcDescriptorsLength(ds []*Descriptor) uint16 {
	length := uint16(0)
	for _, d := range ds {
		length += 2 // tag and length
		length += uint16(calcDescriptorLength(d))
	}
	return length
}

func writeDescriptors(w *astikit.BitsWriter, ds []*Descriptor) (int, error) {
	written := 0

	for _, d := range ds {
		n, err := writeDescriptor(w, d)
		if err != nil {
			return 0, err
		}
		written += n
	}

	return written, nil
}

func writeDescriptorsWithLength(w *astikit.BitsWriter, ds []*Descriptor) (int, error) {
	length := calcDescriptorsLength(ds)
	b := astikit.NewBitsWriterBatch(w)

	b.WriteN(uint8(0xff), 4) // reserved
	b.WriteN(length, 12)     // program_info_length

	if err := b.Err(); err != nil {
		return 0, err
	}

	written, err := writeDescriptors(w, ds)
	return written + 2, err // 2 for length
}
