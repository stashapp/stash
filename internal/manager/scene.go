package manager

import (
	"fmt"
	"net/url"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/models"
)

func GetVideoFileContainer(file *file.VideoFile) (ffmpeg.Container, error) {
	var container ffmpeg.Container
	format := file.Format
	if format != "" {
		container = ffmpeg.Container(format)
	} else { // container isn't in the DB
		// shouldn't happen, fallback to ffprobe
		ffprobe := GetInstance().FFProbe
		tmpVideoFile, err := ffprobe.NewVideoFile(file.Path)
		if err != nil {
			return ffmpeg.Container(""), fmt.Errorf("error reading video file: %v", err)
		}

		return ffmpeg.MatchContainer(tmpVideoFile.Container, file.Path)
	}

	return container, nil
}

func includeSceneStreamPath(f *file.VideoFile, streamingResolution models.StreamingResolutionEnum, maxStreamingTranscodeSize models.StreamingResolutionEnum) bool {
	// convert StreamingResolutionEnum to ResolutionEnum so we can get the min
	// resolution
	convertedRes := models.ResolutionEnum(streamingResolution)

	minResolution := convertedRes.GetMinResolution()
	sceneResolution := f.GetMinResolution()

	// don't include if scene resolution is smaller than the streamingResolution
	if sceneResolution != 0 && sceneResolution < minResolution {
		return false
	}

	// if we always allow everything, then return true
	if maxStreamingTranscodeSize == models.StreamingResolutionEnumOriginal {
		return true
	}

	// convert StreamingResolutionEnum to ResolutionEnum
	maxStreamingResolution := models.ResolutionEnum(maxStreamingTranscodeSize)
	return maxStreamingResolution.GetMinResolution() >= minResolution
}

type SceneStreamEndpoint struct {
	URL      string  `json:"url"`
	MimeType *string `json:"mime_type"`
	Label    *string `json:"label"`
}

func makeDASHStreamEndpoint(directStreamURL *url.URL, label string, resolution models.StreamingResolutionEnum) *SceneStreamEndpoint {
	mimeType := ffmpeg.MimeDASH

	url := *directStreamURL
	url.Path += ".mpd"

	if resolution != "" {
		v := url.Query()
		v.Set("resolution", resolution.String())
		url.RawQuery = v.Encode()
	}

	return &SceneStreamEndpoint{
		URL:      url.String(),
		MimeType: &mimeType,
		Label:    &label,
	}
}

func makeHLSStreamEndpoint(directStreamURL *url.URL, label string, resolution models.StreamingResolutionEnum) *SceneStreamEndpoint {
	mimeType := ffmpeg.MimeHLS

	url := *directStreamURL
	url.Path += ".m3u8"

	if resolution != "" {
		v := url.Query()
		v.Set("resolution", resolution.String())
		url.RawQuery = v.Encode()
	}

	return &SceneStreamEndpoint{
		URL:      url.String(),
		MimeType: &mimeType,
		Label:    &label,
	}
}

func GetSceneStreamPaths(scene *models.Scene, directStreamURL *url.URL, maxStreamingTranscodeSize models.StreamingResolutionEnum) ([]*SceneStreamEndpoint, error) {
	if scene == nil {
		return nil, fmt.Errorf("nil scene")
	}

	pf := scene.Files.Primary()
	if pf == nil {
		return nil, fmt.Errorf("nil file")
	}

	var endpoints []*SceneStreamEndpoint

	// direct stream should only apply when the audio codec is supported
	audioCodec := ffmpeg.MissingUnsupported
	if pf.AudioCodec != "" {
		audioCodec = ffmpeg.ProbeAudioCodec(pf.AudioCodec)
	}

	// don't care if we can't get the container
	container, _ := GetVideoFileContainer(pf)

	if HasTranscode(scene, config.GetInstance().GetVideoFileNamingAlgorithm()) || ffmpeg.IsValidAudioForContainer(audioCodec, container) {
		label := "Direct stream"
		mimeType := ffmpeg.MimeMp4Video
		endpoints = append(endpoints, &SceneStreamEndpoint{
			URL:      directStreamURL.String(),
			MimeType: &mimeType,
			Label:    &label,
		})
	}

	// only add mkv hls stream endpoint if the scene container is an mkv already
	// hls also only supports H264 (with mpeg-ts) and AAC/MP3
	// this endpoint will copy the streams directly without re-encoding for hls,
	// so playing/scrubbing might not work perfectly everywhere
	if container == ffmpeg.Matroska && pf.VideoCodec == ffmpeg.H264 && (audioCodec == ffmpeg.Aac || audioCodec == ffmpeg.Mp3) {
		url := *directStreamURL
		url.Path += ".mkv"

		label := "MKV"
		mimeType := ffmpeg.MimeHLS
		endpoints = append(endpoints, &SceneStreamEndpoint{
			URL:      url.String(),
			MimeType: &mimeType,
			Label:    &label,
		})
	}

	dashLabel := "DASH (VP9)"
	dashLabelFourK := "DASH 4K (2160p)"         // "FOUR_K"
	dashLabelFullHD := "DASH Full HD (1080p)"   // "FULL_HD"
	dashLabelStandardHD := "DASH HD (720p)"     // "STANDARD_HD"
	dashLabelStandard := "DASH Standard (480p)" // "STANDARD"
	dashLabelLow := "DASH Low (240p)"           // "LOW"

	hlsLabel := "HLS (MP4)"
	hlsLabelFourK := "HLS 4K (2160p)"         // "FOUR_K"
	hlsLabelFullHD := "HLS Full HD (1080p)"   // "FULL_HD"
	hlsLabelStandardHD := "HLS HD (720p)"     // "STANDARD_HD"
	hlsLabelStandard := "HLS Standard (480p)" // "STANDARD"
	hlsLabelLow := "HLS Low (240p)"           // "LOW"

	dashStreams := []*SceneStreamEndpoint{
		makeDASHStreamEndpoint(directStreamURL, dashLabel, ""),
	}
	hlsStreams := []*SceneStreamEndpoint{
		makeHLSStreamEndpoint(directStreamURL, hlsLabel, ""),
	}

	if includeSceneStreamPath(pf, models.StreamingResolutionEnumFourK, maxStreamingTranscodeSize) {
		dashStreams = append(dashStreams, makeDASHStreamEndpoint(directStreamURL, dashLabelFourK, models.StreamingResolutionEnumFourK))
		hlsStreams = append(hlsStreams, makeHLSStreamEndpoint(directStreamURL, hlsLabelFourK, models.StreamingResolutionEnumFourK))
	}

	if includeSceneStreamPath(pf, models.StreamingResolutionEnumFullHd, maxStreamingTranscodeSize) {
		dashStreams = append(dashStreams, makeDASHStreamEndpoint(directStreamURL, dashLabelFullHD, models.StreamingResolutionEnumFullHd))
		hlsStreams = append(hlsStreams, makeHLSStreamEndpoint(directStreamURL, hlsLabelFullHD, models.StreamingResolutionEnumFullHd))
	}

	if includeSceneStreamPath(pf, models.StreamingResolutionEnumStandardHd, maxStreamingTranscodeSize) {
		dashStreams = append(dashStreams, makeDASHStreamEndpoint(directStreamURL, dashLabelStandardHD, models.StreamingResolutionEnumStandardHd))
		hlsStreams = append(hlsStreams, makeHLSStreamEndpoint(directStreamURL, hlsLabelStandardHD, models.StreamingResolutionEnumStandardHd))
	}

	if includeSceneStreamPath(pf, models.StreamingResolutionEnumStandard, maxStreamingTranscodeSize) {
		dashStreams = append(dashStreams, makeDASHStreamEndpoint(directStreamURL, dashLabelStandard, models.StreamingResolutionEnumStandard))
		hlsStreams = append(hlsStreams, makeHLSStreamEndpoint(directStreamURL, hlsLabelStandard, models.StreamingResolutionEnumStandard))
	}

	if includeSceneStreamPath(pf, models.StreamingResolutionEnumLow, maxStreamingTranscodeSize) {
		dashStreams = append(dashStreams, makeDASHStreamEndpoint(directStreamURL, dashLabelLow, models.StreamingResolutionEnumLow))
		hlsStreams = append(hlsStreams, makeHLSStreamEndpoint(directStreamURL, hlsLabelLow, models.StreamingResolutionEnumLow))
	}

	endpoints = append(endpoints, dashStreams...)
	endpoints = append(endpoints, hlsStreams...)

	return endpoints, nil
}

// HasTranscode returns true if a transcoded video exists for the provided
// scene. It will check using the OSHash of the scene first, then fall back
// to the checksum.
func HasTranscode(scene *models.Scene, fileNamingAlgo models.HashAlgorithm) bool {
	if scene == nil {
		return false
	}

	sceneHash := scene.GetHash(fileNamingAlgo)
	if sceneHash == "" {
		return false
	}

	transcodePath := instance.Paths.Scene.GetTranscodePath(sceneHash)
	ret, _ := fsutil.FileExists(transcodePath)
	return ret
}
