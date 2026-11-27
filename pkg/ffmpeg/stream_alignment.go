package ffmpeg

import "net/url"

const (
	AlignmentParamKey   = "align"
	StreamAlignment16x8 = "16x8"
)

func hlsSegmentQuery(requestURL *url.URL, resolution string) string {
	query := url.Values{}
	if resolution != "" {
		query.Set(resolutionParamKey, resolution)
	}
	if alignment := requestURL.Query().Get(AlignmentParamKey); alignment != "" {
		query.Set(AlignmentParamKey, alignment)
	}
	copyAuthParams(query, requestURL.Query())
	if len(query) == 0 {
		return ""
	}
	return "?" + query.Encode()
}

func ValidHLSAlignment(streamType *StreamType, alignment string) bool {
	if alignment == "" {
		return true
	}
	return streamType == StreamTypeHLS && alignment == StreamAlignment16x8
}

func streamDirWithAlignment(base, alignment string) string {
	if alignment == StreamAlignment16x8 {
		return base + "_align16x8"
	}
	return base
}

func videoFilterWithAlignment(filter VideoFilter, alignment string) VideoFilter {
	if alignment == StreamAlignment16x8 {
		return filter.Append("pad=ceil(iw/16)*16:ceil(ih/8)*8:x=0:y=0:color=black")
	}
	return filter
}

func hlsCodecWithAlignment(sm *StreamManager, name, alignment string) VideoCodec {
	if alignment == StreamAlignment16x8 {
		return VideoCodecLibX264
	}
	return HLSGetCodec(sm, name)
}

func fullHardwareWithAlignment(configured bool, alignment string) bool {
	return configured && alignment == ""
}
