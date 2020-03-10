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
)

var ValidCodecs = []string{"h264", "h265", "vp8", "vp9"}

type Container string

const (
	Mp4      Container = "mp4"
	M4v      Container = "m4v"
	Mov      Container = "mov"
	Wmv      Container = "wmv"
	Webm     Container = "webm"
	Matroska Container = "matroska"
	Avi      Container = "avi"
	Flv      Container = "flv"
	Mpegts   Container = "mpegts"
)

var validForH264 = []Container{Mp4}
var validForH265 = []Container{Mp4}
var validForVp8 = []Container{Webm}
var validForVp9 = []Container{Webm}

//maps user readable container strings to ffprobe's format_name
//on some formats ffprobe can't differentiate
var ContainerToFfprobe = map[Container]string{
	Mp4:      "mov,mp4,m4a,3gp,3g2,mj2",
	M4v:      "mov,mp4,m4a,3gp,3g2,mj2",
	Mov:      "mov,mp4,m4a,3gp,3g2,mj2",
	Wmv:      "asf",
	Webm:     "matroska,webm",
	Matroska: "matroska,webm",
	Avi:      "avi",
	Flv:      "flv",
	Mpegts:   "mpegts",
}

var FfprobeToContainer = map[string]Container{
	"mov,mp4,m4a,3gp,3g2,mj2": Mp4, // browsers support all of them so we don't  care
	"asf":                     Wmv,
	"avi":                     Avi,
	"flv":                     Flv,
	"mpegts":                  Mpegts,
	"matroska,webm":           Matroska,
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

func IsValidCodec(codecName string) bool {
	for _, c := range ValidCodecs {
		if c == codecName {
			return true
		}
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

//extend stream validation check to take into account container
func IsValidCombo(codecName string, format Container) bool {
	switch codecName {
	case "h264":
		return IsValidForContainer(format, validForH264)
	case "h265":
		return IsValidForContainer(format, validForH265)
	case "vp8":
		return IsValidForContainer(format, validForVp8)
	case "vp9":
		return IsValidForContainer(format, validForVp9)
	}
	return false
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

	AudioCodec string
}

// Execute exec command and bind result to struct.
func NewVideoFile(ffprobePath string, videoPath string) (*VideoFile, error) {
	args := []string{"-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", "-show_error", videoPath}
	//// Extremely slow on windows for some reason
	//if runtime.GOOS != "windows" {
	//	args = append(args, "-count_frames")
	//}
	out, err := exec.Command(ffprobePath, args...).Output()

	if err != nil {
		return nil, fmt.Errorf("FFProbe encountered an error with <%s>.\nError JSON:\n%s\nError: %s", videoPath, string(out), err.Error())
	}

	probeJSON := &FFProbeJSON{}
	if err := json.Unmarshal(out, probeJSON); err != nil {
		return nil, err
	}

	return parse(videoPath, probeJSON)
}

func parse(filePath string, probeJSON *FFProbeJSON) (*VideoFile, error) {
	if probeJSON == nil {
		return nil, fmt.Errorf("failed to get ffprobe json")
	}

	result := &VideoFile{}
	result.JSON = *probeJSON

	if result.JSON.Error.Code != 0 {
		return nil, fmt.Errorf("ffprobe error code %d: %s", result.JSON.Error.Code, result.JSON.Error.String)
	}
	//} else if (ffprobeResult.stderr.includes("could not find codec parameters")) {
	//	throw new Error(`FFProbe [${filePath}] -> Could not find codec parameters`);
	//} // TODO nil_or_unsupported.(video_stream) && nil_or_unsupported.(audio_stream)

	result.Path = filePath
	result.Title = probeJSON.Format.Tags.Title

	if result.Title == "" {
		// default title to filename
		result.SetTitleFromPath()
	}

	result.Comment = probeJSON.Format.Tags.Comment

	result.Bitrate, _ = strconv.ParseInt(probeJSON.Format.BitRate, 10, 64)
	result.Container = probeJSON.Format.FormatName
	duration, _ := strconv.ParseFloat(probeJSON.Format.Duration, 64)
	result.Duration = math.Round(duration*100) / 100
	fileStat, _ := os.Stat(filePath)
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

func (v *VideoFile) SetTitleFromPath() {
	v.Title = filepath.Base(v.Path)
}
