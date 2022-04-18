package ffmpeg

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/exec"
	"github.com/stashapp/stash/pkg/logger"
)

// VideoFile represents the ffprobe output for a video file.
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

// TranscodeScale calculates the dimension scaling for a transcode, where maxSize is the maximum size of the longest dimension of the input video.
// If no scaling is required, then returns 0, 0.
// Returns -2 for the dimension that will scale to maintain aspect ratio.
func (v *VideoFile) TranscodeScale(maxSize int) (int, int) {
	// get the smaller dimension of the video file
	videoSize := v.Height
	if v.Width < videoSize {
		videoSize = v.Width
	}

	// if our streaming resolution is larger than the video dimension
	// or we are streaming the original resolution, then just set the
	// input width
	if maxSize >= videoSize || maxSize == 0 {
		return 0, 0
	}

	// we're setting either the width or height
	// we'll set the smaller dimesion
	if v.Width > v.Height {
		// set the height
		return -2, maxSize
	}

	return maxSize, -2
}

// FFProbe provides an interface to the ffprobe executable.
type FFProbe string

// NewVideoFile runs ffprobe on the given path and returns a VideoFile.
func (f *FFProbe) NewVideoFile(videoPath string) (*VideoFile, error) {
	args := []string{"-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", "-show_error", videoPath}
	cmd := exec.Command(string(*f), args...)
	out, err := cmd.Output()

	if err != nil {
		return nil, fmt.Errorf("FFProbe encountered an error with <%s>.\nError JSON:\n%s\nError: %s", videoPath, string(out), err.Error())
	}

	probeJSON := &FFProbeJSON{}
	if err := json.Unmarshal(out, probeJSON); err != nil {
		return nil, fmt.Errorf("error unmarshalling video data for <%s>: %s", videoPath, err.Error())
	}

	return parse(videoPath, probeJSON)
}

// GetReadFrameCount counts the actual frames of the video file.
// Used when the frame count is missing or incorrect.
func (f *FFProbe) GetReadFrameCount(path string) (int64, error) {
	args := []string{"-v", "quiet", "-print_format", "json", "-count_frames", "-show_format", "-show_streams", "-show_error", path}
	out, err := exec.Command(string(*f), args...).Output()

	if err != nil {
		return 0, fmt.Errorf("FFProbe encountered an error with <%s>.\nError JSON:\n%s\nError: %s", path, string(out), err.Error())
	}

	probeJSON := &FFProbeJSON{}
	if err := json.Unmarshal(out, probeJSON); err != nil {
		return 0, fmt.Errorf("error unmarshalling video data for <%s>: %s", path, err.Error())
	}

	fc, err := parse(path, probeJSON)
	return fc.FrameCount, err
}

func parse(filePath string, probeJSON *FFProbeJSON) (*VideoFile, error) {
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

	audioStream := result.getAudioStream()
	if audioStream != nil {
		result.AudioCodec = audioStream.CodecName
		result.AudioStream = audioStream
	}

	videoStream := result.getVideoStream()
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

func (v *VideoFile) getAudioStream() *FFProbeStream {
	index := v.getStreamIndex("audio", v.JSON)
	if index != -1 {
		return &v.JSON.Streams[index]
	}
	return nil
}

func (v *VideoFile) getVideoStream() *FFProbeStream {
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
