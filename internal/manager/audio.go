// TODO(audio): update this file
package manager

import (
	"fmt"
	"net/url"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/models"
)

type AudioStreamEndpoint struct {
	URL      string  `json:"url"`
	MimeType *string `json:"mime_type"`
	Label    *string `json:"label"`
}

var (
	// TODO(audio): figure out what stream types we need, and what we can support
	directAudioEndpointType = endpointType{
		label:     "Direct stream",
		mimeType:  ffmpeg.MimeMp4Audio,
		extension: "",
	}
	mp3AudioEndpointType = endpointType{
		label:     "MP3",
		mimeType:  ffmpeg.MimeMp3Audio,
		extension: ".mp3",
	}
)

func GetAudioFileContainer(file *models.AudioFile) (ffmpeg.Container, error) {
	var container ffmpeg.Container
	format := file.Format
	if format != "" {
		container = ffmpeg.Container(format)
	} else { // container isn't in the DB
		// shouldn't happen, fallback to ffprobe
		ffprobe := GetInstance().FFProbe
		tmpAudioFile, err := ffprobe.NewAudioFile(file.Path)
		if err != nil {
			return ffmpeg.Container(""), fmt.Errorf("error reading video file: %v", err)
		}

		return ffmpeg.MatchContainer(tmpAudioFile.Container, file.Path)
	}

	return container, nil
}

func GetAudioStreamPaths(audio *models.Audio, directStreamURL *url.URL, maxStreamingTranscodeSize models.StreamingResolutionEnum) ([]*AudioStreamEndpoint, error) {
	if audio == nil {
		return nil, fmt.Errorf("nil audio")
	}

	pf := audio.Files.Primary()
	if pf == nil {
		return nil, nil
	}

	// convert StreamingResolutionEnum to ResolutionEnum
	maxStreamingResolution := models.ResolutionEnum(maxStreamingTranscodeSize)
	audioResolution := models.GetMinResolution(pf)
	includeAudioStreamPath := func(streamingResolution models.StreamingResolutionEnum) bool {
		var minResolution int
		if streamingResolution == models.StreamingResolutionEnumOriginal {
			minResolution = audioResolution
		} else {
			// convert StreamingResolutionEnum to ResolutionEnum so we can get the min
			// resolution
			convertedRes := models.ResolutionEnum(streamingResolution)
			minResolution = convertedRes.GetMinResolution()

			// don't include if audio resolution is smaller than the streamingResolution
			if audioResolution != 0 && audioResolution < minResolution {
				return false
			}
		}

		// if we always allow everything, then return true
		if maxStreamingTranscodeSize == models.StreamingResolutionEnumOriginal {
			return true
		}

		return maxStreamingResolution.GetMinResolution() >= minResolution
	}

	makeStreamEndpoint := func(t endpointType, resolution models.StreamingResolutionEnum) *AudioStreamEndpoint {
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

		return &AudioStreamEndpoint{
			URL:      url.String(),
			MimeType: &t.mimeType,
			Label:    &label,
		}
	}

	var endpoints []*AudioStreamEndpoint

	// direct stream should only apply when the audio codec is supported
	audioCodec := ffmpeg.MissingUnsupported
	if pf.AudioCodec != "" {
		audioCodec = ffmpeg.ProbeAudioCodec(pf.AudioCodec)
	}

	// don't care if we can't get the container
	container, _ := GetAudioFileContainer(pf)

	if HasAudioTranscode(audio, config.GetInstance().GetAudioFileNamingAlgorithm()) || ffmpeg.IsValidAudioForContainer(audioCodec, container) {
		endpoints = append(endpoints, makeStreamEndpoint(directAudioEndpointType, ""))
	}

	// only add mkv stream endpoint if the audio container is an mkv already
	if container == ffmpeg.Matroska {
		endpoints = append(endpoints, makeStreamEndpoint(mkvAudioEndpointType, ""))
	}

	mp4Streams := []*AudioStreamEndpoint{}
	webmStreams := []*AudioStreamEndpoint{}
	hlsStreams := []*AudioStreamEndpoint{}
	dashStreams := []*AudioStreamEndpoint{}

	if includeAudioStreamPath(models.StreamingResolutionEnumOriginal) {
		mp4Streams = append(mp4Streams, makeStreamEndpoint(mp3AudioEndpointType, models.StreamingResolutionEnumOriginal))
		webmStreams = append(webmStreams, makeStreamEndpoint(webmAudioEndpointType, models.StreamingResolutionEnumOriginal))
		hlsStreams = append(hlsStreams, makeStreamEndpoint(hlsAudioEndpointType, models.StreamingResolutionEnumOriginal))
		dashStreams = append(dashStreams, makeStreamEndpoint(dashAudioEndpointType, models.StreamingResolutionEnumOriginal))
	}

	if includeAudioStreamPath(models.StreamingResolutionEnumFourK) {
		mp4Streams = append(mp4Streams, makeStreamEndpoint(mp3AudioEndpointType, models.StreamingResolutionEnumFourK))
		webmStreams = append(webmStreams, makeStreamEndpoint(webmAudioEndpointType, models.StreamingResolutionEnumFourK))
		hlsStreams = append(hlsStreams, makeStreamEndpoint(hlsAudioEndpointType, models.StreamingResolutionEnumFourK))
		dashStreams = append(dashStreams, makeStreamEndpoint(dashAudioEndpointType, models.StreamingResolutionEnumFourK))
	}

	if includeAudioStreamPath(models.StreamingResolutionEnumFullHd) {
		mp4Streams = append(mp4Streams, makeStreamEndpoint(mp3AudioEndpointType, models.StreamingResolutionEnumFullHd))
		webmStreams = append(webmStreams, makeStreamEndpoint(webmAudioEndpointType, models.StreamingResolutionEnumFullHd))
		hlsStreams = append(hlsStreams, makeStreamEndpoint(hlsAudioEndpointType, models.StreamingResolutionEnumFullHd))
		dashStreams = append(dashStreams, makeStreamEndpoint(dashAudioEndpointType, models.StreamingResolutionEnumFullHd))
	}

	if includeAudioStreamPath(models.StreamingResolutionEnumStandardHd) {
		mp4Streams = append(mp4Streams, makeStreamEndpoint(mp3AudioEndpointType, models.StreamingResolutionEnumStandardHd))
		webmStreams = append(webmStreams, makeStreamEndpoint(webmAudioEndpointType, models.StreamingResolutionEnumStandardHd))
		hlsStreams = append(hlsStreams, makeStreamEndpoint(hlsAudioEndpointType, models.StreamingResolutionEnumStandardHd))
		dashStreams = append(dashStreams, makeStreamEndpoint(dashAudioEndpointType, models.StreamingResolutionEnumStandardHd))
	}

	if includeAudioStreamPath(models.StreamingResolutionEnumStandard) {
		mp4Streams = append(mp4Streams, makeStreamEndpoint(mp3AudioEndpointType, models.StreamingResolutionEnumStandard))
		webmStreams = append(webmStreams, makeStreamEndpoint(webmAudioEndpointType, models.StreamingResolutionEnumStandard))
		hlsStreams = append(hlsStreams, makeStreamEndpoint(hlsAudioEndpointType, models.StreamingResolutionEnumStandard))
		dashStreams = append(dashStreams, makeStreamEndpoint(dashAudioEndpointType, models.StreamingResolutionEnumStandard))
	}

	if includeAudioStreamPath(models.StreamingResolutionEnumLow) {
		mp4Streams = append(mp4Streams, makeStreamEndpoint(mp3AudioEndpointType, models.StreamingResolutionEnumLow))
		webmStreams = append(webmStreams, makeStreamEndpoint(webmAudioEndpointType, models.StreamingResolutionEnumLow))
		hlsStreams = append(hlsStreams, makeStreamEndpoint(hlsAudioEndpointType, models.StreamingResolutionEnumLow))
		dashStreams = append(dashStreams, makeStreamEndpoint(dashAudioEndpointType, models.StreamingResolutionEnumLow))
	}

	endpoints = append(endpoints, mp4Streams...)
	endpoints = append(endpoints, webmStreams...)
	endpoints = append(endpoints, hlsStreams...)
	endpoints = append(endpoints, dashStreams...)

	return endpoints, nil
}

// HasAudioTranscode returns true if a transcoded video exists for the provided
// audio. It will check using the OSHash of the audio first, then fall back
// to the checksum.
func HasAudioTranscode(audio *models.Audio, fileNamingAlgo models.HashAlgorithm) bool {
	if audio == nil {
		return false
	}

	audioHash := audio.GetHash(fileNamingAlgo)
	if audioHash == "" {
		return false
	}

	transcodePath := instance.Paths.Audio.GetTranscodePath(audioHash)
	ret, _ := fsutil.FileExists(transcodePath)
	return ret
}
