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
		mimeType:  ffmpeg.MimeMp3Audio,
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

	makeStreamEndpoint := func(t endpointType) *AudioStreamEndpoint {
		url := *directStreamURL
		url.Path += t.extension

		label := t.label

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
		endpoints = append(endpoints, makeStreamEndpoint(directAudioEndpointType))
	}

	mp3Streams := []*AudioStreamEndpoint{}
	hlsStreams := []*AudioStreamEndpoint{}
	dashStreams := []*AudioStreamEndpoint{}

	// TODO(audio): do we need the `if includeAudioStreamPath() {`?
	mp3Streams = append(mp3Streams, makeStreamEndpoint(mp3AudioEndpointType))
	hlsStreams = append(hlsStreams, makeStreamEndpoint(hlsEndpointType))
	dashStreams = append(dashStreams, makeStreamEndpoint(dashEndpointType))

	endpoints = append(endpoints, mp3Streams...)
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
