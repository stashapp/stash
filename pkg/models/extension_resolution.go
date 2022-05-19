package models

type ResolutionRange struct {
	min, max int
}

var resolutionRanges = map[ResolutionEnum]ResolutionRange{
	ResolutionEnumVeryLow:    {144, 239},
	ResolutionEnumLow:        {240, 359},
	ResolutionEnumR360p:      {360, 479},
	ResolutionEnumStandard:   {480, 539},
	ResolutionEnumWebHd:      {540, 719},
	ResolutionEnumStandardHd: {720, 1079},
	ResolutionEnumFullHd:     {1080, 1439},
	ResolutionEnumQuadHd:     {1440, 1919},
	ResolutionEnumVrHd:       {1920, 2159},
	ResolutionEnumFourK:      {2160, 2879},
	ResolutionEnumFiveK:      {2880, 3383},
	ResolutionEnumSixK:       {3384, 4319},
	ResolutionEnumEightK:     {4320, 8639},
}

// GetMaxResolution returns the maximum width or height that media must be
// to qualify as this resolution.
func (r *ResolutionEnum) GetMaxResolution() int {
	return resolutionRanges[*r].max
}

// GetMinResolution returns the minimum width or height that media must be
// to qualify as this resolution.
func (r ResolutionEnum) GetMinResolution() int {
	return resolutionRanges[r].min
}

var streamingResolutionMax = map[StreamingResolutionEnum]int{
	StreamingResolutionEnumLow:        resolutionRanges[ResolutionEnumLow].min,
	StreamingResolutionEnumStandard:   resolutionRanges[ResolutionEnumStandard].min,
	StreamingResolutionEnumStandardHd: resolutionRanges[ResolutionEnumStandardHd].min,
	StreamingResolutionEnumFullHd:     resolutionRanges[ResolutionEnumFullHd].min,
	StreamingResolutionEnumFourK:      resolutionRanges[ResolutionEnumFourK].min,
	StreamingResolutionEnumOriginal:   0,
}

func (r StreamingResolutionEnum) GetMaxResolution() int {
	return streamingResolutionMax[r]
}
