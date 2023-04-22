package mpd

type SegmentBase struct {
	Initialization           *URL     `xml:"Initialization,omitempty"`
	RepresentationIndex      *URL     `xml:"RepresentationIndex,omitempty"`
	Timescale                *uint32  `xml:"timescale,attr,omitempty"`
	PresentationTimeOffset   *uint64  `xml:"presentationTimeOffset,attr,omitempty"`
	IndexRange               *string  `xml:"indexRange,attr,omitempty"`
	IndexRangeExact          *bool    `xml:"indexRangeExact,attr,omitempty"`
	AvailabilityTimeOffset   *float32 `xml:"availabilityTimeOffset,attr,omitempty"`
	AvailabilityTimeComplete *bool    `xml:"availabilityTimeComplete,attr,omitempty"`
}

type MultipleSegmentBase struct {
	SegmentBase
	SegmentTimeline    *SegmentTimeline `xml:"SegmentTimeline,omitempty"`
	BitstreamSwitching *URL             `xml:"BitstreamSwitching,omitempty"`
	Duration           *uint32          `xml:"duration,attr,omitempty"`
	StartNumber        *uint32          `xml:"startNumber,attr,omitempty"`
}

type SegmentList struct {
	MultipleSegmentBase
	SegmentURLs []*SegmentURL `xml:"SegmentURL,omitempty"`
}

type SegmentURL struct {
	Media      *string `xml:"media,attr,omitempty"`
	MediaRange *string `xml:"mediaRange,attr,omitempty"`
	Index      *string `xml:"index,attr,omitempty"`
	IndexRange *string `xml:"indexRange,attr,omitempty"`
}

type SegmentTimeline struct {
	Segments []*SegmentTimelineSegment `xml:"S,omitempty"`
}

type SegmentTimelineSegment struct {
	StartTime   *uint64 `xml:"t,attr,omitempty"`
	Duration    uint64  `xml:"d,attr"`
	RepeatCount *int    `xml:"r,attr,omitempty"`
}

type URL struct {
	SourceURL *string `xml:"sourceURL,attr,omitempty"`
	Range     *string `xml:"range,attr,omitempty"`
}
