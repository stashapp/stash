package ffmpeg

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	stashExec "github.com/stashapp/stash/pkg/exec"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
)

const minimumFFProbeVersion = 5

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
type FFProbe struct {
	path    string
	version Version
}

func (f *FFProbe) Path() string {
	return f.path
}

var ffprobeVersionRE = regexp.MustCompile(`ffprobe version n?((\d+)\.(\d+)(?:\.(\d+))?)`)

func (f *FFProbe) getVersion() error {
	var args []string
	args = append(args, "-version")
	cmd := stashExec.Command(f.path, args...)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	var err error
	if err = cmd.Run(); err != nil {
		return err
	}

	stdoutStr := stdout.String()
	match := ffprobeVersionRE.FindStringSubmatchIndex(stdoutStr)
	if match == nil {
		return errors.New("version string malformed")
	}

	majorS := stdoutStr[match[4]:match[5]]
	minorS := stdoutStr[match[6]:match[7]]

	// patch is optional
	var patchS string
	if match[8] != -1 && match[9] != -1 {
		patchS = stdoutStr[match[8]:match[9]]
	}

	if i, err := strconv.Atoi(majorS); err == nil {
		f.version.major = i
	}
	if i, err := strconv.Atoi(minorS); err == nil {
		f.version.minor = i
	}
	if i, err := strconv.Atoi(patchS); err == nil {
		f.version.patch = i
	}
	logger.Debugf("FFProbe version %s detected", f.version.String())

	return nil
}

// Creates a new FFProbe instance.
func NewFFProbe(path string) *FFProbe {
	ret := &FFProbe{
		path: path,
	}
	if err := ret.getVersion(); err != nil {
		logger.Warnf("FFProbe version not detected %v", err)
	}

	if ret.version.major != 0 && ret.version.major < minimumFFProbeVersion {
		logger.Warnf("FFProbe version %d.%d.%d detected, but %d.x or later is required", ret.version.major, ret.version.minor, ret.version.patch, minimumFFProbeVersion)
	}

	return ret
}

// NewVideoFile runs ffprobe on the given path and returns a VideoFile.
func (f *FFProbe) NewVideoFile(videoPath string) (*VideoFile, error) {
	args := []string{
		"-v",
		"quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		"-show_error",
	}

	// show_entries stream_side_data=rotation requires 5.x or later ffprobe
	if f.version.major >= 5 {
		args = append(args, "-show_entries", "stream_side_data=rotation")
	}

	args = append(args, videoPath)

	cmd := stashExec.Command(f.path, args...)
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
	out, err := stashExec.Command(f.path, args...).Output()

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
		result.Width = videoStream.Width
		result.Height = videoStream.Height

		if isRotated(videoStream) {
			result.Width = videoStream.Height
			result.Height = videoStream.Width
		}

		result.VideoStreamDuration, err = strconv.ParseFloat(videoStream.Duration, 64)
		if err != nil {
			// Revert to the historical behaviour, which is still correct in the vast majority of cases.
			result.VideoStreamDuration = result.FileDuration
		}
	}

	return result, nil
}

func isRotated(s *FFProbeStream) bool {
	rotate, _ := strconv.ParseInt(s.Tags.Rotate, 10, 64)
	if rotate != 180 && rotate != 0 {
		return true
	}

	for _, sd := range s.SideDataList {
		r := sd.Rotation
		if r < 0 {
			r = -r
		}
		if r != 0 && r != 180 {
			return true
		}
	}

	return false
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
