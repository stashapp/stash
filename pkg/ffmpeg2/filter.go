package ffmpeg2

import "fmt"

type VideoFilter string

func (f VideoFilter) Args() []string {
	if f == "" {
		return nil
	}

	return []string{"-vf", string(f)}
}

func (f VideoFilter) ScaleWidth(w int) VideoFilter {
	return f.Append(fmt.Sprintf("scale=%v:-2", w))
}

func (f VideoFilter) ScaleDimensions(maxDimensions int) VideoFilter {
	return f.Append(fmt.Sprintf("scale=%v:%v:force_original_aspect_ratio=decrease", maxDimensions, maxDimensions))
}

func (f VideoFilter) Fps(fps int) VideoFilter {
	return f.Append(fmt.Sprintf("fps=%v", fps))
}

func (f VideoFilter) Select(frame int) VideoFilter {
	return f.Append(fmt.Sprintf("select=eq(n\\,%d)", frame))
}

func (f VideoFilter) Append(s string) VideoFilter {
	// if filter is empty, then just set
	if f == "" {
		return VideoFilter(s)
	}

	return VideoFilter(fmt.Sprintf(",%s", s))
}
