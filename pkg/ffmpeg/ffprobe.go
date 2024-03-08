package ffmpeg

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	stashExec "github.com/stashapp/stash/pkg/exec"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
)

func ValidateFFProbe(ffprobePath string) error {
	cmd := stashExec.Command(ffprobePath, "-h")
	bytes, err := cmd.CombinedOutput()
	output := string(bytes)
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return fmt.Errorf("error running ffprobe: %v", output)
		}

		return fmt.Errorf("error running ffprobe: %v", err)
	}

	return nil
}

func LookPathFFProbe() string {
	ret, _ := exec.LookPath(getFFProbeFilename())

	if ret != "" {
		if err := ValidateFFProbe(ret); err != nil {
			logger.Warnf("ffprobe found in PATH (%s), but it is missing required flags: %v", ret, err)
			ret = ""
		}
	}

	return ret
}

func FindFFProbe(path string) string {
	ret := fsutil.FindInPaths([]string{path}, getFFProbeFilename())

	if ret != "" {
		if err := ValidateFFProbe(ret); err != nil {
			logger.Warnf("ffprobe found (%s), but it is missing required flags: %v", ret, err)
			ret = ""
		}
	}

	return ret
}

// ResolveFFMpeg attempts to resolve the path to the ffmpeg executable.
// It first looks in the provided path, then resolves from the environment, and finally looks in the fallback path.
// Returns an empty string if a valid ffmpeg cannot be found.
func ResolveFFProbe(path string, fallbackPath string) string {
	// look in the provided path first
	ret := FindFFProbe(path)
	if ret != "" {
		return ret
	}

	// then resolve from the environment
	ret = LookPathFFProbe()
	if ret != "" {
		return ret
	}

	// finally, look in the fallback path
	ret = FindFFProbe(fallbackPath)
	return ret
}

// VideoFile represents the ffprobe output for a video file.
type VideoFile struct {
	JSON        FFProbeJSON
	AudioStream *FFProbeStream
	VideoStream *FFProbeStream

	Path      string
	Title     string
	Comment   string
	Container string
	// FileDuration is the declared (meta-data) duration of the *file*.
	// In most cases (sprites, previews, etc.) we actually care about the duration of the video stream specifically,
	// because those two can differ slightly (e.g. audio stream longer than the video stream, making the whole file
	// longer).
	FileDuration        float64
	VideoStreamDuration float64
	StartTime           float64
	Bitrate             int64
	Size                int64
	CreationTime        time.Time

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

func (f *FFProbe) Path() string {
	return string(*f)
}

// NewVideoFile runs ffprobe on the given path and returns a VideoFile.
func (f *FFProbe) NewVideoFile(videoPath string) (*VideoFile, error) {
	args := []string{"-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", "-show_error", videoPath}
	cmd := stashExec.Command(string(*f), args...)
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
	out, err := stashExec.Command(string(*f), args...).Output()

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
	result.FileDuration = math.Round(duration*100) / 100
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
		if math.IsNaN(framerate) {
			framerate = 0
		}
		result.FrameRate = math.Round(framerate*100) / 100
		if rotate, err := strconv.ParseInt(videoStream.Tags.Rotate, 10, 64); err == nil && rotate != 180 {
			result.Width = videoStream.Height
			result.Height = videoStream.Width
		} else {
			result.Width = videoStream.Width
			result.Height = videoStream.Height
		}
		result.VideoStreamDuration, err = strconv.ParseFloat(videoStream.Duration, 64)
		if err != nil {
			// Revert to the historical behaviour, which is still correct in the vast majority of cases.
			result.VideoStreamDuration = result.FileDuration
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
	ret := -1
	for i, stream := range probeJSON.Streams {
		// skip cover art/thumbnails
		if stream.CodecType == fileType && stream.Disposition.AttachedPic == 0 {
			// prefer default stream
			if stream.Disposition.Default == 1 {
				return i
			}

			// backwards compatible behaviour - fallback to first matching stream
			if ret == -1 {
				ret = i
			}
		}
	}

	return ret
}
