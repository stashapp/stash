package ffmpeg

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/desktop"
	"github.com/stashapp/stash/pkg/logger"
)

type Container string
type AudioCodec string

const (
	Mp4                Container  = "mp4"
	M4v                Container  = "m4v"
	Mov                Container  = "mov"
	Wmv                Container  = "wmv"
	Webm               Container  = "webm"
	Matroska           Container  = "matroska"
	Avi                Container  = "avi"
	Flv                Container  = "flv"
	Mpegts             Container  = "mpegts"
	Aac                AudioCodec = "aac"
	Mp3                AudioCodec = "mp3"
	Opus               AudioCodec = "opus"
	Vorbis             AudioCodec = "vorbis"
	MissingUnsupported AudioCodec = ""
	Mp4Ffmpeg          string     = "mov,mp4,m4a,3gp,3g2,mj2" // browsers support all of them
	M4vFfmpeg          string     = "mov,mp4,m4a,3gp,3g2,mj2" // so we don't care that ffmpeg
	MovFfmpeg          string     = "mov,mp4,m4a,3gp,3g2,mj2" // can't differentiate between them
	WmvFfmpeg          string     = "asf"
	WebmFfmpeg         string     = "matroska,webm"
	MatroskaFfmpeg     string     = "matroska,webm"
	AviFfmpeg          string     = "avi"
	FlvFfmpeg          string     = "flv"
	MpegtsFfmpeg       string     = "mpegts"
	H264               string     = "h264"
	H265               string     = "h265" // found in rare cases from a faulty encoder
	Hevc               string     = "hevc"
	Vp8                string     = "vp8"
	Vp9                string     = "vp9"
	Mkv                string     = "mkv" // only used from the browser to indicate mkv support
	Hls                string     = "hls" // only used from the browser to indicate hls support
	MimeWebm           string     = "video/webm"
	MimeMkv            string     = "video/x-matroska"
	MimeMp4            string     = "video/mp4"
	MimeHLS            string     = "application/vnd.apple.mpegurl"
	MimeMpegts         string     = "video/MP2T"
)

// only support H264 by default, since Safari does not support VP8/VP9
var DefaultSupportedCodecs = []string{H264, H265}

var validForH264Mkv = []Container{Mp4, Matroska}
var validForH264 = []Container{Mp4}
var validForH265Mkv = []Container{Mp4, Matroska}
var validForH265 = []Container{Mp4}
var validForVp8 = []Container{Webm}
var validForVp9Mkv = []Container{Webm, Matroska}
var validForVp9 = []Container{Webm}
var validForHevcMkv = []Container{Mp4, Matroska}
var validForHevc = []Container{Mp4}

var validAudioForMkv = []AudioCodec{Aac, Mp3, Vorbis, Opus}
var validAudioForWebm = []AudioCodec{Vorbis, Opus}
var validAudioForMp4 = []AudioCodec{Aac, Mp3}

// ContainerToFfprobe maps user readable container strings to ffprobe's format_name.
// On some formats ffprobe can't differentiate
var ContainerToFfprobe = map[Container]string{
	Mp4:      Mp4Ffmpeg,
	M4v:      M4vFfmpeg,
	Mov:      MovFfmpeg,
	Wmv:      WmvFfmpeg,
	Webm:     WebmFfmpeg,
	Matroska: MatroskaFfmpeg,
	Avi:      AviFfmpeg,
	Flv:      FlvFfmpeg,
	Mpegts:   MpegtsFfmpeg,
}

var FfprobeToContainer = map[string]Container{
	Mp4Ffmpeg:      Mp4,
	WmvFfmpeg:      Wmv,
	AviFfmpeg:      Avi,
	FlvFfmpeg:      Flv,
	MpegtsFfmpeg:   Mpegts,
	MatroskaFfmpeg: Matroska,
}

func MatchContainer(format string, filePath string) Container { // match ffprobe string to our Container

	container := FfprobeToContainer[format]
	if container == Matroska {
		container = MagicContainer(filePath) // use magic number instead of ffprobe for matroska,webm
	}
	if container == "" { // if format is not in our Container list leave it as ffprobes reported format_name
		container = Container(format)
	}
	return container
}

func IsValidCodec(codecName string, supportedCodecs []string) bool {
	for _, c := range supportedCodecs {
		if c == codecName {
			return true
		}
	}
	return false
}

func IsValidAudio(audio AudioCodec, validCodecs []AudioCodec) bool {

	// if audio codec is missing or unsupported by ffmpeg we can't do anything about it
	// report it as valid so that the file can at least be streamed directly if the video codec is supported
	if audio == MissingUnsupported {
		return true
	}

	for _, c := range validCodecs {
		if c == audio {
			return true
		}
	}

	return false
}

func IsValidAudioForContainer(audio AudioCodec, format Container) bool {
	switch format {
	case Matroska:
		return IsValidAudio(audio, validAudioForMkv)
	case Webm:
		return IsValidAudio(audio, validAudioForWebm)
	case Mp4:
		return IsValidAudio(audio, validAudioForMp4)
	}
	return false

}

func IsValidForContainer(format Container, validContainers []Container) bool {
	for _, fmt := range validContainers {
		if fmt == format {
			return true
		}
	}
	return false
}

// IsValidCombo checks if a codec/container combination is valid.
// Returns true on validity, false otherwise
func IsValidCombo(codecName string, format Container, supportedVideoCodecs []string) bool {
	supportMKV := IsValidCodec(Mkv, supportedVideoCodecs)
	supportHEVC := IsValidCodec(Hevc, supportedVideoCodecs)

	switch codecName {
	case H264:
		if supportMKV {
			return IsValidForContainer(format, validForH264Mkv)
		}
		return IsValidForContainer(format, validForH264)
	case H265:
		if supportMKV {
			return IsValidForContainer(format, validForH265Mkv)
		}
		return IsValidForContainer(format, validForH265)
	case Vp8:
		return IsValidForContainer(format, validForVp8)
	case Vp9:
		if supportMKV {
			return IsValidForContainer(format, validForVp9Mkv)
		}
		return IsValidForContainer(format, validForVp9)
	case Hevc:
		if supportHEVC {
			if supportMKV {
				return IsValidForContainer(format, validForHevcMkv)
			}
			return IsValidForContainer(format, validForHevc)
		}
	}
	return false
}

func IsStreamable(videoCodec string, audioCodec AudioCodec, container Container) bool {
	supportedVideoCodecs := DefaultSupportedCodecs

	// check if the video codec matches the supported codecs
	return IsValidCodec(videoCodec, supportedVideoCodecs) && IsValidCombo(videoCodec, container, supportedVideoCodecs) && IsValidAudioForContainer(audioCodec, container)
}

type VideoFile struct {
	JSON        FFProbeJSON
	AudioStream *FFProbeStream
	VideoStream *FFProbeStream

	Path         string
	Title        string
	Comment      string
	Container    string
	Duration     float64
	StartTime    float64
	Bitrate      int64
	Size         int64
	CreationTime time.Time

	VideoCodec   string
	VideoBitrate int64
	Width        int
	Height       int
	FrameRate    float64
	Rotation     int64
	FrameCount   int64

	AudioCodec string
}

// FFProbe
type FFProbe string

// Execute exec command and bind result to struct.
func (f *FFProbe) NewVideoFile(videoPath string, stripExt bool) (*VideoFile, error) {
	args := []string{"-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", "-show_error", videoPath}
	cmd := exec.Command(string(*f), args...)
	desktop.HideExecShell(cmd)
	out, err := cmd.Output()

	if err != nil {
		return nil, fmt.Errorf("FFProbe encountered an error with <%s>.\nError JSON:\n%s\nError: %s", videoPath, string(out), err.Error())
	}

	probeJSON := &FFProbeJSON{}
	if err := json.Unmarshal(out, probeJSON); err != nil {
		return nil, fmt.Errorf("error unmarshalling video data for <%s>: %s", videoPath, err.Error())
	}

	return parse(videoPath, probeJSON, stripExt)
}

// GetReadFrameCount counts the actual frames of the video file
func (f *FFProbe) GetReadFrameCount(vf *VideoFile) (int64, error) {
	args := []string{"-v", "quiet", "-print_format", "json", "-count_frames", "-show_format", "-show_streams", "-show_error", vf.Path}
	out, err := exec.Command(string(*f), args...).Output()

	if err != nil {
		return 0, fmt.Errorf("FFProbe encountered an error with <%s>.\nError JSON:\n%s\nError: %s", vf.Path, string(out), err.Error())
	}

	probeJSON := &FFProbeJSON{}
	if err := json.Unmarshal(out, probeJSON); err != nil {
		return 0, fmt.Errorf("error unmarshalling video data for <%s>: %s", vf.Path, err.Error())
	}

	fc, err := parse(vf.Path, probeJSON, false)
	return fc.FrameCount, err
}

func parse(filePath string, probeJSON *FFProbeJSON, stripExt bool) (*VideoFile, error) {
	if probeJSON == nil {
		return nil, fmt.Errorf("failed to get ffprobe json for <%s>", filePath)
	}

	result := &VideoFile{}
	result.JSON = *probeJSON

	if result.JSON.Error.Code != 0 {
		return nil, fmt.Errorf("ffprobe error code %d: %s", result.JSON.Error.Code, result.JSON.Error.String)
	}

	result.Path = filePath
	result.Title = probeJSON.Format.Tags.Title

	if result.Title == "" {
		// default title to filename
		result.SetTitleFromPath(stripExt)
	}

	result.Comment = probeJSON.Format.Tags.Comment
	result.Bitrate, _ = strconv.ParseInt(probeJSON.Format.BitRate, 10, 64)

	result.Container = probeJSON.Format.FormatName
	duration, _ := strconv.ParseFloat(probeJSON.Format.Duration, 64)
	result.Duration = math.Round(duration*100) / 100
	fileStat, err := os.Stat(filePath)
	if err != nil {
		statErr := fmt.Errorf("error statting file <%s>: %w", filePath, err)
		logger.Errorf("%v", statErr)
		return nil, statErr
	}
	result.Size = fileStat.Size()
	result.StartTime, _ = strconv.ParseFloat(probeJSON.Format.StartTime, 64)
	result.CreationTime = probeJSON.Format.Tags.CreationTime.Time

	audioStream := result.GetAudioStream()
	if audioStream != nil {
		result.AudioCodec = audioStream.CodecName
		result.AudioStream = audioStream
	}

	videoStream := result.GetVideoStream()
	if videoStream != nil {
		result.VideoStream = videoStream
		result.VideoCodec = videoStream.CodecName
		result.FrameCount, _ = strconv.ParseInt(videoStream.NbFrames, 10, 64)
		if videoStream.NbReadFrames != "" { // if ffprobe counted the frames use that instead
			fc, _ := strconv.ParseInt(videoStream.NbReadFrames, 10, 64)
			if fc > 0 {
				result.FrameCount, _ = strconv.ParseInt(videoStream.NbReadFrames, 10, 64)
			} else {
				logger.Debugf("[ffprobe] <%s> invalid Read Frames count", videoStream.NbReadFrames)
			}
		}
		result.VideoBitrate, _ = strconv.ParseInt(videoStream.BitRate, 10, 64)
		var framerate float64
		if strings.Contains(videoStream.AvgFrameRate, "/") {
			frameRateSplit := strings.Split(videoStream.AvgFrameRate, "/")
			numerator, _ := strconv.ParseFloat(frameRateSplit[0], 64)
			denominator, _ := strconv.ParseFloat(frameRateSplit[1], 64)
			framerate = numerator / denominator
		} else {
			framerate, _ = strconv.ParseFloat(videoStream.AvgFrameRate, 64)
		}
		result.FrameRate = math.Round(framerate*100) / 100
		if rotate, err := strconv.ParseInt(videoStream.Tags.Rotate, 10, 64); err == nil && rotate != 180 {
			result.Width = videoStream.Height
			result.Height = videoStream.Width
		} else {
			result.Width = videoStream.Width
			result.Height = videoStream.Height
		}
	}

	return result, nil
}

func (v *VideoFile) GetAudioStream() *FFProbeStream {
	index := v.getStreamIndex("audio", v.JSON)
	if index != -1 {
		return &v.JSON.Streams[index]
	}
	return nil
}

func (v *VideoFile) GetVideoStream() *FFProbeStream {
	index := v.getStreamIndex("video", v.JSON)
	if index != -1 {
		return &v.JSON.Streams[index]
	}
	return nil
}

func (v *VideoFile) getStreamIndex(fileType string, probeJSON FFProbeJSON) int {
	for i, stream := range probeJSON.Streams {
		if stream.CodecType == fileType {
			return i
		}
	}

	return -1
}

func (v *VideoFile) SetTitleFromPath(stripExtension bool) {
	v.Title = filepath.Base(v.Path)
	if stripExtension {
		ext := filepath.Ext(v.Title)
		v.Title = strings.TrimSuffix(v.Title, ext)
	}

}
