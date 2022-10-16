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

func makeStreamEndpoint(streamURL *url.URL, streamingResolution models.StreamingResolutionEnum, mimeType, label string) *SceneStreamEndpoint {
	urlCopy := *streamURL
	v := urlCopy.Query()
	v.Set("resolution", streamingResolution.String())
	urlCopy.RawQuery = v.Encode()

	return &SceneStreamEndpoint{
		URL:      urlCopy.String(),
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

	var ret []*SceneStreamEndpoint
	mimeHLS := ffmpeg.MimeHLS
	mimeMp4 := ffmpeg.MimeMp4

	// direct stream should only apply when the audio codec is supported
	audioCodec := ffmpeg.MissingUnsupported
	if pf.AudioCodec != "" {
		audioCodec = ffmpeg.ProbeAudioCodec(pf.AudioCodec)
	}

	// don't care if we can't get the container
	container, _ := GetVideoFileContainer(pf)

	replaceSuffix := func(suffix string) *url.URL {
		urlCopy := *directStreamURL
		urlCopy.Path += suffix
		return &urlCopy
	}

	if HasTranscode(scene, config.GetInstance().GetVideoFileNamingAlgorithm()) || ffmpeg.IsValidAudioForContainer(audioCodec, container) {
		label := "Direct stream"
		ret = append(ret, &SceneStreamEndpoint{
			URL:      directStreamURL.String(),
			MimeType: &mimeMp4,
			Label:    &label,
		})
	}

	// only add mkv stream endpoint if the scene container is an mkv already
	if container == ffmpeg.Matroska {
		label := "mkv"
		ret = append(ret, &SceneStreamEndpoint{
			URL: replaceSuffix(".mkv").String(),
			// set mkv to mp4 to trick the client, since many clients won't try mkv
			MimeType: &mimeMp4,
			Label:    &label,
		})
	}

	// Setup HLS Streams
	hlsLabelOriginal := "HLS (original)"
	hlsLabelFourK := "HLS 4K (2160p)"         // "FOUR_K"
	hlsLabelFullHD := "HLS Full HD (1080p)"   // "FULL_HD"
	hlsLabelStandardHD := "HLS HD (720p)"     // "STANDARD_HD"
	hlsLabelStandard := "HLS Standard (480p)" // "STANDARD"
	hlsLabelLow := "HLS Low (240p)"           // "LOW"

	hlsStreams := []*SceneStreamEndpoint{
		{
			URL:      replaceSuffix(".m3u8").String(),
			MimeType: &mimeHLS,
			Label:    &hlsLabelOriginal,
		},
	}

	hlsURL := replaceSuffix(".m3u8")

	if includeSceneStreamPath(pf, models.StreamingResolutionEnumFourK, maxStreamingTranscodeSize) {
		hlsStreams = append(hlsStreams, makeStreamEndpoint(hlsURL, models.StreamingResolutionEnumFourK, mimeHLS, hlsLabelFourK))
	}

	if includeSceneStreamPath(pf, models.StreamingResolutionEnumFullHd, maxStreamingTranscodeSize) {
		hlsStreams = append(hlsStreams, makeStreamEndpoint(hlsURL, models.StreamingResolutionEnumFullHd, mimeHLS, hlsLabelFullHD))
	}

	if includeSceneStreamPath(pf, models.StreamingResolutionEnumStandardHd, maxStreamingTranscodeSize) {
		hlsStreams = append(hlsStreams, makeStreamEndpoint(hlsURL, models.StreamingResolutionEnumStandardHd, mimeHLS, hlsLabelStandardHD))
	}

	if includeSceneStreamPath(pf, models.StreamingResolutionEnumStandard, maxStreamingTranscodeSize) {
		hlsStreams = append(hlsStreams, makeStreamEndpoint(hlsURL, models.StreamingResolutionEnumStandard, mimeHLS, hlsLabelStandard))
	}

	if includeSceneStreamPath(pf, models.StreamingResolutionEnumLow, maxStreamingTranscodeSize) {
		hlsStreams = append(hlsStreams, makeStreamEndpoint(hlsURL, models.StreamingResolutionEnumLow, mimeHLS, hlsLabelLow))
	}

	ret = append(ret, hlsStreams...)
	// TODO - change this to use hls when we move to videojs
	// still use HLS - but copy the video stream
	// URL:      directStreamURL + ".m3u8?videoCodec=copy",
	// MimeType: &mimeHLS,

	return ret, nil
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
