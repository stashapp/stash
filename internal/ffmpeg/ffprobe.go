package ffmpeg

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type ffprobeExecutable struct {
	Path string
}

type FFProbeResult struct {
	JSON ffprobeJSON

	Path         string
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

func NewFFProbe(ffprobePath string) ffprobeExecutable {
	return ffprobeExecutable{
		Path: ffprobePath,
	}
}

// Execute exec command and bind result to struct.
func (ffp *ffprobeExecutable) ProbeVideo(filePath string) (*FFProbeResult, error) {
	args := []string{"-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", "-show_error", filePath}
	if runtime.GOOS == "windows" {
		args = append(args, "-count_frames")
	}
	out, err := exec.Command(ffp.Path, args...).Output()

	if err != nil {
		return nil, fmt.Errorf("FFProbe encountered an error with <%s>.\nError JSON:\n%s\nError: %s", filePath, string(out), err.Error())
	}

	probeJSON := &ffprobeJSON{}
	if err := json.Unmarshal(out, probeJSON); err != nil {
		return nil, err
	}

	result := ffp.newProbeResult(filePath, *probeJSON)
	return result, nil
}

func (ffp *ffprobeExecutable) newProbeResult(filePath string, probeJson ffprobeJSON) *FFProbeResult {
	videoStreamIndex := ffp.getStreamIndex("video", probeJson)
	audioStreamIndex := ffp.getStreamIndex("audio", probeJson)

	result := &FFProbeResult{}
	result.JSON = probeJson
	result.Path = filePath
	result.Container = probeJson.Format.FormatName
	duration, _ := strconv.ParseFloat(probeJson.Format.Duration, 64)
	result.Duration = math.Round(duration*100)/100
	result.StartTime, _ = strconv.ParseFloat(probeJson.Format.StartTime, 64)
	result.Bitrate, _ = strconv.ParseInt(probeJson.Format.BitRate, 10, 64)
	fileStat, _ := os.Stat(filePath)
	result.Size = fileStat.Size()
	result.CreationTime = probeJson.Format.Tags.CreationTime

	if videoStreamIndex != -1 {
		videoStream := probeJson.Streams[videoStreamIndex]
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
		result.FrameRate = math.Round(framerate*100)/100
		if rotate, err := strconv.ParseInt(videoStream.Tags.Rotate, 10, 64); err == nil && rotate != 180 {
			result.Width = videoStream.Height
			result.Height = videoStream.Width
		} else {
			result.Width = videoStream.Width
			result.Height = videoStream.Height
		}
	}

	if audioStreamIndex != -1 {
		result.AudioCodec = probeJson.Streams[audioStreamIndex].CodecName
	}

	return result
}

func (ffp *ffprobeExecutable) getStreamIndex(fileType string, probeJson ffprobeJSON) int {
	for i, stream := range probeJson.Streams {
		if stream.CodecType == fileType {
			return i
		}
	}

	return -1
}