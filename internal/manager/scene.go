package manager

import (
	"fmt"
	"net/url"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/models"
)

type SceneStreamEndpoint struct {
	URL      string  `json:"url"`
	MimeType *string `json:"mime_type"`
	Label    *string `json:"label"`
}

type endpointType struct {
	label     string
	mimeType  string
	extension string
}

var (
	directEndpointType = endpointType{
		label:     "Direct stream",
		mimeType:  ffmpeg.MimeMp4Video,
		extension: "",
	}
	mp4EndpointType = endpointType{
		label:     "MP4",
		mimeType:  ffmpeg.MimeMp4Video,
		extension: ".mp4",
	}
	mkvEndpointType = endpointType{
		label: "MKV",
		// use mp4 mimetype to trick the client, since many clients won't try mkv
		mimeType:  ffmpeg.MimeMp4Video,
		extension: ".mkv",
	}
	webmEndpointType = endpointType{
		label:     "WEBM",
		mimeType:  ffmpeg.MimeWebmVideo,
		extension: ".webm",
	}
	hlsEndpointType = endpointType{
		label:     "HLS",
		mimeType:  ffmpeg.MimeHLS,
		extension: ".m3u8",
	}
	dashEndpointType = endpointType{
		label:     "DASH",
		mimeType:  ffmpeg.MimeDASH,
		extension: ".mpd",
	}
)

func GetVideoFileContainer(file *models.VideoFile) (ffmpeg.Container, error) {
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

func GetSceneStreamPaths(scene *models.Scene, directStreamURL *url.URL, maxStreamingTranscodeSize models.StreamingResolutionEnum) ([]*SceneStreamEndpoint, error) {
	if scene == nil {
		return nil, fmt.Errorf("nil scene")
	}

	pf := scene.Files.Primary()
	if pf == nil {
		return nil, nil
	}

	// convert StreamingResolutionEnum to ResolutionEnum
	maxStreamingResolution := models.ResolutionEnum(maxStreamingTranscodeSize)
	sceneResolution := models.GetMinResolution(pf)
	includeSceneStreamPath := func(streamingResolution models.StreamingResolutionEnum) bool {
		var minResolution int
		if streamingResolution == models.StreamingResolutionEnumOriginal {
			minResolution = sceneResolution
		} else {
			// convert StreamingResolutionEnum to ResolutionEnum so we can get the min
			// resolution
			convertedRes := models.ResolutionEnum(streamingResolution)
			minResolution = convertedRes.GetMinResolution()

			// don't include if scene resolution is smaller than the streamingResolution
			if sceneResolution != 0 && sceneResolution < minResolution {
				return false
			}
		}

		// if we always allow everything, then return true
		if maxStreamingTranscodeSize == models.StreamingResolutionEnumOriginal {
			return true
		}

		return maxStreamingResolution.GetMinResolution() >= minResolution
	}

	makeStreamEndpoint := func(t endpointType, resolution models.StreamingResolutionEnum) *SceneStreamEndpoint {
		url := *directStreamURL
		url.Path += t.extension

		label := t.label

		if resolution != "" {
			v := url.Query()
			v.Set("resolution", resolution.String())
			url.RawQuery = v.Encode()

			switch resolution {
			case models.StreamingResolutionEnumFourK:
				label += " 4K (2160p)"
			case models.StreamingResolutionEnumFullHd:
				label += " Full HD (1080p)"
			case models.StreamingResolutionEnumStandardHd:
				label += " HD (720p)"
			case models.StreamingResolutionEnumStandard:
				label += " Standard (480p)"
			case models.StreamingResolutionEnumLow:
				label += " Low (240p)"
			}
		}

		return &SceneStreamEndpoint{
			URL:      url.String(),
			MimeType: &t.mimeType,
			Label:    &label,
		}
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
		endpoints = append(endpoints, makeStreamEndpoint(directEndpointType, ""))
	}

	// only add mkv stream endpoint if the scene container is an mkv already
	if container == ffmpeg.Matroska {
		endpoints = append(endpoints, makeStreamEndpoint(mkvEndpointType, ""))
	}

	mp4Streams := []*SceneStreamEndpoint{}
	webmStreams := []*SceneStreamEndpoint{}
	hlsStreams := []*SceneStreamEndpoint{}
	dashStreams := []*SceneStreamEndpoint{}

	if includeSceneStreamPath(models.StreamingResolutionEnumOriginal) {
		mp4Streams = append(mp4Streams, makeStreamEndpoint(mp4EndpointType, models.StreamingResolutionEnumOriginal))
		webmStreams = append(webmStreams, makeStreamEndpoint(webmEndpointType, models.StreamingResolutionEnumOriginal))
		hlsStreams = append(hlsStreams, makeStreamEndpoint(hlsEndpointType, models.StreamingResolutionEnumOriginal))
		dashStreams = append(dashStreams, makeStreamEndpoint(dashEndpointType, models.StreamingResolutionEnumOriginal))
	}

	if includeSceneStreamPath(models.StreamingResolutionEnumFourK) {
		mp4Streams = append(mp4Streams, makeStreamEndpoint(mp4EndpointType, models.StreamingResolutionEnumFourK))
		webmStreams = append(webmStreams, makeStreamEndpoint(webmEndpointType, models.StreamingResolutionEnumFourK))
		hlsStreams = append(hlsStreams, makeStreamEndpoint(hlsEndpointType, models.StreamingResolutionEnumFourK))
		dashStreams = append(dashStreams, makeStreamEndpoint(dashEndpointType, models.StreamingResolutionEnumFourK))
	}

	if includeSceneStreamPath(models.StreamingResolutionEnumFullHd) {
		mp4Streams = append(mp4Streams, makeStreamEndpoint(mp4EndpointType, models.StreamingResolutionEnumFullHd))
		webmStreams = append(webmStreams, makeStreamEndpoint(webmEndpointType, models.StreamingResolutionEnumFullHd))
		hlsStreams = append(hlsStreams, makeStreamEndpoint(hlsEndpointType, models.StreamingResolutionEnumFullHd))
		dashStreams = append(dashStreams, makeStreamEndpoint(dashEndpointType, models.StreamingResolutionEnumFullHd))
	}

	if includeSceneStreamPath(models.StreamingResolutionEnumStandardHd) {
		mp4Streams = append(mp4Streams, makeStreamEndpoint(mp4EndpointType, models.StreamingResolutionEnumStandardHd))
		webmStreams = append(webmStreams, makeStreamEndpoint(webmEndpointType, models.StreamingResolutionEnumStandardHd))
		hlsStreams = append(hlsStreams, makeStreamEndpoint(hlsEndpointType, models.StreamingResolutionEnumStandardHd))
		dashStreams = append(dashStreams, makeStreamEndpoint(dashEndpointType, models.StreamingResolutionEnumStandardHd))
	}

	if includeSceneStreamPath(models.StreamingResolutionEnumStandard) {
		mp4Streams = append(mp4Streams, makeStreamEndpoint(mp4EndpointType, models.StreamingResolutionEnumStandard))
		webmStreams = append(webmStreams, makeStreamEndpoint(webmEndpointType, models.StreamingResolutionEnumStandard))
		hlsStreams = append(hlsStreams, makeStreamEndpoint(hlsEndpointType, models.StreamingResolutionEnumStandard))
		dashStreams = append(dashStreams, makeStreamEndpoint(dashEndpointType, models.StreamingResolutionEnumStandard))
	}

	if includeSceneStreamPath(models.StreamingResolutionEnumLow) {
		mp4Streams = append(mp4Streams, makeStreamEndpoint(mp4EndpointType, models.StreamingResolutionEnumLow))
		webmStreams = append(webmStreams, makeStreamEndpoint(webmEndpointType, models.StreamingResolutionEnumLow))
		hlsStreams = append(hlsStreams, makeStreamEndpoint(hlsEndpointType, models.StreamingResolutionEnumLow))
		dashStreams = append(dashStreams, makeStreamEndpoint(dashEndpointType, models.StreamingResolutionEnumLow))
	}

	endpoints = append(endpoints, mp4Streams...)
	endpoints = append(endpoints, webmStreams...)
	endpoints = append(endpoints, hlsStreams...)
	endpoints = append(endpoints, dashStreams...)

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
