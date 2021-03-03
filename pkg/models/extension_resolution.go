package models

var resolutionMax = []int{
	240,
	360,
	480,
	540,
	720,
	1080,
	1440,
	1920,
	2160,
	2880,
	3384,
	4320,
	0,
}

// GetMaxResolution returns the maximum width or height that media must be
// to qualify as this resolution. A return value of 0 means that there is no
// maximum.
func (r *ResolutionEnum) GetMaxResolution() int {
	if !r.IsValid() {
		return 0
	}

	// sanity check - length of arrays must be the same
	if len(resolutionMax) != len(AllResolutionEnum) {
		panic("resolutionMax array length != AllResolutionEnum array length")
	}

	for i, rr := range AllResolutionEnum {
		if rr == *r {
			return resolutionMax[i]
		}
	}

	return 0
}

// GetMinResolution returns the minimum width or height that media must be
// to qualify as this resolution.
func (r *ResolutionEnum) GetMinResolution() int {
	if !r.IsValid() {
		return 0
	}

	// sanity check - length of arrays must be the same
	if len(resolutionMax) != len(AllResolutionEnum) {
		panic("resolutionMax array length != AllResolutionEnum array length")
	}

	// use the previous resolution max as this resolution min
	for i, rr := range AllResolutionEnum {
		if rr == *r {
			if i > 0 {
				return resolutionMax[i-1]
			}

			return 0
		}
	}

	return 0
}
