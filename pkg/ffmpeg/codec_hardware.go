package ffmpeg

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/logger"
)

var (
	// Hardware codec's
	VideoCodecN264  = makeVideoCodec("H264 NVENC", "h264_nvenc")
	VideoCodecN264H = makeVideoCodec("H264 NVENC HQ profile", "h264_nvenc")
	VideoCodecI264  = makeVideoCodec("H264 Intel Quick Sync Video (QSV)", "h264_qsv")
	VideoCodecI264C = makeVideoCodec("H264 Intel Quick Sync Video (QSV) Compatibility profile", "h264_qsv")
	VideoCodecA264  = makeVideoCodec("H264 Advanced Media Framework (AMF)", "h264_amf")
	VideoCodecM264  = makeVideoCodec("H264 VideoToolbox", "h264_videotoolbox")
	VideoCodecV264  = makeVideoCodec("H264 VAAPI", "h264_vaapi")
	VideoCodecR264  = makeVideoCodec("H264 V4L2M2M", "h264_v4l2m2m")
	VideoCodecO264  = makeVideoCodec("H264 OMX", "h264_omx")
	VideoCodecIVP9  = makeVideoCodec("VP9 Intel Quick Sync Video (QSV)", "vp9_qsv")
	VideoCodecVVP9  = makeVideoCodec("VP9 VAAPI", "vp9_vaapi")
	VideoCodecVVPX  = makeVideoCodec("VP8 VAAPI", "vp8_vaapi")
	VideoCodecRK264 = makeVideoCodec("H264 Rockchip MPP (rkmpp)", "h264_rkmpp")
)

const minHeight int = 480

// Tests all (given) hardware codec's
func (f *FFMpeg) InitHWSupport(ctx context.Context) {
	// do the hardware codec tests in a separate goroutine to avoid blocking
	done := make(chan struct{})
	go func() {
		f.initHWSupport(ctx)
		close(done)
	}()

	// log if the initialization takes too long
	const hwInitLogTimeoutSecondsDefault = 5
	hwInitLogTimeoutSeconds := hwInitLogTimeoutSecondsDefault * time.Second
	timer := time.NewTimer(hwInitLogTimeoutSeconds)

	go func() {
		select {
		case <-timer.C:
			logger.Warnf("[InitHWSupport] Hardware codec initialization is taking longer than %s...", hwInitLogTimeoutSeconds)
			logger.Info("[InitHWSupport] Hardware encoding will not be available until initialization is complete.")
		case <-done:
			if !timer.Stop() {
				<-timer.C
			}
		}
	}()
}

func (f *FFMpeg) initHWSupport(ctx context.Context) {
	var hwCodecSupport []VideoCodec

	// Note that the first compatible codec is returned, so order is important
	for _, codec := range []VideoCodec{
		VideoCodecN264H,
		VideoCodecN264,
		VideoCodecI264,
		VideoCodecI264C,
		VideoCodecV264,
		VideoCodecR264,
		VideoCodecRK264,
		VideoCodecIVP9,
		VideoCodecVVP9,
		VideoCodecM264,
	} {
		var args Args
		args = append(args, "-hide_banner")
		args = args.LogLevel(LogLevelWarning)
		args = f.HWDeviceInit(args, codec, false)
		args = args.Format("lavfi")
		width, height := 1280, 720
		args = args.Input(fmt.Sprintf("color=c=red:s=%dx%d", width, height))
		args = args.Duration(0.1)

		// Test scaling
		videoFilter := f.HWMaxResFilter(codec, width, height, minHeight, false)
		args = append(args, codec.ExtraArgs()...)
		args = args.VideoFilter(videoFilter)

		args = args.Format("null")
		args = args.Output("-")

		// #6064 - add timeout to context to prevent hangs
		const hwTestTimeoutSecondsDefault = 10
		hwTestTimeoutSeconds := hwTestTimeoutSecondsDefault * time.Second

		// allow timeout to be overridden with environment variable
		if timeout := os.Getenv("STASH_HW_TEST_TIMEOUT"); timeout != "" {
			if seconds, err := strconv.Atoi(timeout); err == nil {
				hwTestTimeoutSeconds = time.Duration(seconds) * time.Second
			}
		}

		testCtx, cancel := context.WithTimeout(ctx, hwTestTimeoutSeconds)
		defer cancel()

		cmd := f.Command(testCtx, args)
		cmd.WaitDelay = time.Second
		logger.Tracef("[InitHWSupport] Testing codec %s: %v", codec, cmd.Args)

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			if testCtx.Err() != nil {
				logger.Debugf("[InitHWSupport] Codec %s test timed out after %s", codec, hwTestTimeoutSeconds)
				continue
			}

			errOutput := stderr.String()

			if len(errOutput) == 0 {
				errOutput = err.Error()
			}

			logger.Debugf("[InitHWSupport] Codec %s not supported. Error output:\n%s", codec, errOutput)
		} else {
			hwCodecSupport = append(hwCodecSupport, codec)
		}
	}

	outstr := fmt.Sprintf("[InitHWSupport] Supported HW codecs [%d]:\n", len(hwCodecSupport))
	for _, codec := range hwCodecSupport {
		outstr += fmt.Sprintf("\t%s - %s\n", codec.Name, codec.CodeName)
	}
	logger.Info(outstr)

	f.hwCodecSupportMutex.Lock()
	defer f.hwCodecSupportMutex.Unlock()
	f.hwCodecSupport = hwCodecSupport
}

func (f *FFMpeg) HWCanFullHWTranscode(ctx context.Context, codec VideoCodec, path string, width int, height int, reqHeight int) bool {
	if codec == VideoCodecCopy {
		return false
	}

	var args Args
	args = append(args, "-hide_banner")
	args = args.LogLevel(LogLevelWarning)
	args = args.XError()
	args = f.HWDeviceInit(args, codec, true)
	args = args.Input(path)
	args = args.Duration(1)

	videoFilter := f.HWMaxResFilter(codec, width, height, reqHeight, true)
	args = append(args, codec.ExtraArgs()...)
	args = args.VideoFilter(videoFilter)

	args = args.Format("null")
	args = args.Output("-")

	cmd := f.Command(ctx, args)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errOutput := stderr.String()

		if len(errOutput) == 0 {
			errOutput = err.Error()
		}

		logger.Debugf("[InitHWSupport] Full hardware transcode for file %s not supported. Error output:\n%s", path, errOutput)
		return false
	}

	return true
}

// Prepend input for hardware encoding only
func (f *FFMpeg) HWDeviceInit(args Args, toCodec VideoCodec, fullhw bool) Args {
	switch toCodec {
	case VideoCodecN264,
		VideoCodecN264H:
		args = append(args, "-hwaccel_device")
		args = append(args, "0")
		if fullhw {
			args = append(args, "-threads")
			args = append(args, "1")
			args = append(args, "-hwaccel")
			args = append(args, "cuda")
			args = append(args, "-hwaccel_output_format")
			args = append(args, "cuda")
		}
	case VideoCodecV264,
		VideoCodecVVP9:
		args = append(args, "-vaapi_device")
		args = append(args, "/dev/dri/renderD128")
		if fullhw {
			args = append(args, "-hwaccel")
			args = append(args, "vaapi")
			args = append(args, "-hwaccel_output_format")
			args = append(args, "vaapi")
		}
	case VideoCodecI264,
		VideoCodecI264C,
		VideoCodecIVP9:
		if fullhw {
			args = append(args, "-hwaccel")
			args = append(args, "qsv")
			args = append(args, "-hwaccel_output_format")
			args = append(args, "qsv")
		} else {
			args = append(args, "-init_hw_device")
			args = append(args, "qsv=hw")
			args = append(args, "-filter_hw_device")
			args = append(args, "hw")
		}
	case VideoCodecM264:
		if fullhw {
			args = append(args, "-hwaccel")
			args = append(args, "videotoolbox")
			args = append(args, "-hwaccel_output_format")
			args = append(args, "videotoolbox_vld")
		} else {
			args = append(args, "-init_hw_device")
			args = append(args, "videotoolbox=vt")
		}
	case VideoCodecRK264:
		// Rockchip: always create rkmpp device and make it the filter device, so
		// scale_rkrga and subsequent hwupload/hwmap operate in the right context.
		args = append(args, "-init_hw_device")
		args = append(args, "rkmpp=rk")
		args = append(args, "-filter_hw_device")
		args = append(args, "rk")
		if fullhw {
			args = append(args, "-hwaccel")
			args = append(args, "rkmpp")
			args = append(args, "-hwaccel_output_format")
			args = append(args, "drm_prime")
		}
	}

	return args
}

// Initialise a video filter for HW encoding
func (f *FFMpeg) HWFilterInit(toCodec VideoCodec, fullhw bool) VideoFilter {
	var videoFilter VideoFilter
	switch toCodec {
	case VideoCodecV264,
		VideoCodecVVP9:
		if !fullhw {
			videoFilter = videoFilter.Append("format=nv12")
			videoFilter = videoFilter.Append("hwupload")
		}
	case VideoCodecN264, VideoCodecN264H:
		if !fullhw {
			videoFilter = videoFilter.Append("format=nv12")
			videoFilter = videoFilter.Append("hwupload_cuda")
		}
	case VideoCodecI264,
		VideoCodecI264C,
		VideoCodecIVP9:
		if !fullhw {
			videoFilter = videoFilter.Append("hwupload=extra_hw_frames=64")
			videoFilter = videoFilter.Append("format=qsv")
		}
	case VideoCodecM264:
		if !fullhw {
			videoFilter = videoFilter.Append("format=nv12")
			videoFilter = videoFilter.Append("hwupload")
		}
	case VideoCodecRK264:
		// For Rockchip full-hw, do NOT pre-map to rkrga here. scale_rkrga can
		// consume DRM_PRIME frames directly when filter_hw_device is set.
		// For non-fullhw, keep a sane software format.
		if !fullhw {
			videoFilter = videoFilter.Append("format=nv12")
			videoFilter = videoFilter.Append("hwupload")
		}
	}

	return videoFilter
}

var scaler_re = regexp.MustCompile(`scale=(?P<value>([-\d]+):([-\d]+))`)

func templateReplaceScale(input string, template string, match []int, width, height int, minusonehack bool) string {
	result := []byte{}

	if minusonehack {
		// Parse width and height
		w, err := strconv.Atoi(input[match[4]:match[5]])
		if err != nil {
			logger.Error("failed to parse width")
			return input
		}
		h, err := strconv.Atoi(input[match[6]:match[7]])
		if err != nil {
			logger.Error("failed to parse height")
			return input
		}

		// Calculate ratio
		ratio := float64(width) / float64(height)
		if w < 0 {
			w = int(math.Round(float64(h) * ratio))
		} else if h < 0 {
			h = int(math.Round(float64(w) / ratio))
		}

		// Fix not divisible by 2 errors
		if w%2 != 0 {
			w++
		}
		if h%2 != 0 {
			h++
		}

		template = strings.ReplaceAll(template, "$value", fmt.Sprintf("%d:%d", w, h))
	}

	res := string(scaler_re.ExpandString(result, template, input, match))

	matchStart := match[0]
	matchEnd := match[1]

	return input[0:matchStart] + res + input[matchEnd:]
}

// Replace video filter scaling with hardware scaling for full hardware transcoding (also fixes the format)
func (f *FFMpeg) HWCodecFilter(args VideoFilter, codec VideoCodec, width, height int, fullhw bool) VideoFilter {
	sargs := string(args)

	match := scaler_re.FindStringSubmatchIndex(sargs)
	if match == nil {
		return f.HWApplyFullHWFilter(args, codec, fullhw)
	}

	return f.HWApplyScaleTemplate(sargs, codec, match, width, height, fullhw)
}

// Apply format switching if applicable
func (f *FFMpeg) HWApplyFullHWFilter(args VideoFilter, codec VideoCodec, fullhw bool) VideoFilter {
	switch codec {
	case VideoCodecN264, VideoCodecN264H:
		if fullhw && f.version.Gteq(Version{major: 5}) { // Added in FFMpeg 5
			args = args.Append("scale_cuda=format=yuv420p")
		}
	case VideoCodecV264, VideoCodecVVP9:
		if fullhw && f.version.Gteq(Version{major: 3, minor: 1}) { // Added in FFMpeg 3.1
			args = args.Append("scale_vaapi=format=nv12")
		}
	case VideoCodecI264, VideoCodecI264C, VideoCodecIVP9:
		if fullhw && f.version.Gteq(Version{major: 3, minor: 3}) { // Added in FFMpeg 3.3
			args = args.Append("scale_qsv=format=nv12")
		}
	case VideoCodecRK264:
		// For Rockchip, no extra mapping here. If there is no scale filter,
		// leave frames in DRM_PRIME for the encoder.
	}

	return args
}

// Switch scaler
func (f *FFMpeg) HWApplyScaleTemplate(sargs string, codec VideoCodec, match []int, width, height int, fullhw bool) VideoFilter {
	var template string

	switch codec {
	case VideoCodecN264, VideoCodecN264H:
		template = "scale_cuda=$value"
		if fullhw && f.version.Gteq(Version{major: 5}) { // Added in FFMpeg 5
			template += ":format=yuv420p"
		}
	case VideoCodecV264, VideoCodecVVP9:
		template = "scale_vaapi=$value"
		if fullhw && f.version.Gteq(Version{major: 3, minor: 1}) { // Added in FFMpeg 3.1
			template += ":format=nv12"
		}
	case VideoCodecI264, VideoCodecI264C, VideoCodecIVP9:
		template = "scale_qsv=$value"
		if fullhw && f.version.Gteq(Version{major: 3, minor: 3}) { // Added in FFMpeg 3.3
			template += ":format=nv12"
		}
	case VideoCodecM264:
		template = "scale_vt=$value"
	case VideoCodecRK264:
		// The original filter chain is a fallback for maximum compatibility:
		// "scale_rkrga=$value:format=nv12,hwdownload,format=nv12,hwupload"
		// It avoids hwmap(rkrga→rkmpp) failures (-38/-12) seen on some builds
		// by downloading the scaled frame to system RAM and re-uploading it.
		// The filter chain below uses a zero-copy approach, passing the hardware-scaled
		// frame directly to the encoder. This is more efficient but may be less stable.
		template = "scale_rkrga=$value"
	default:
		return VideoFilter(sargs)
	}

	// BUG: [scale_qsv]: Size values less than -1 are not acceptable.
	isIntel := codec == VideoCodecI264 || codec == VideoCodecI264C || codec == VideoCodecIVP9
	// BUG: scale_vt doesn't call ff_scale_adjust_dimensions, thus cant accept negative size values
	isApple := codec == VideoCodecM264
	// Rockchip's scale_rkrga supports -1/-2; don't apply minus-one hack here.
	return VideoFilter(templateReplaceScale(sargs, template, match, width, height, isIntel || isApple))
}

// Returns the max resolution for a given codec, or a default
func (f *FFMpeg) HWCodecMaxRes(codec VideoCodec) (int, int) {
	switch codec {
	case VideoCodecRK264:
		return 8192, 8192
	case VideoCodecN264,
		VideoCodecN264H,
		VideoCodecI264,
		VideoCodecI264C:
		return 4096, 4096
	}

	return 0, 0
}

// Return a maxres filter
func (f *FFMpeg) HWMaxResFilter(toCodec VideoCodec, width int, height int, reqHeight int, fullhw bool) VideoFilter {
	if width == 0 || height == 0 {
		return ""
	}
	videoFilter := f.HWFilterInit(toCodec, fullhw)
	maxWidth, maxHeight := f.HWCodecMaxRes(toCodec)
	videoFilter = videoFilter.ScaleMaxLM(width, height, reqHeight, maxWidth, maxHeight)
	return f.HWCodecFilter(videoFilter, toCodec, width, height, fullhw)
}

// Return a hardware accelerated codec for HLS if available, otherwise the default
func (f *FFMpeg) HWCodecHLSCompatible(defaultCodec VideoCodec) VideoCodec {
	for _, element := range f.getHWCodecSupport() {
		switch element {
		case VideoCodecN264,
			VideoCodecN264H,
			VideoCodecI264,
			VideoCodecI264C,
			VideoCodecV264,
			VideoCodecR264,
			VideoCodecM264, // Note that the Apple encoder sucks at startup, thus HLS quality is crap
			VideoCodecRK264:
			return element
		}
	}
	return defaultCodec
}

// Return a hardware accelerated codec for MP4 if available, otherwise the default
func (f *FFMpeg) HWCodecMP4Compatible(defaultCodec VideoCodec) VideoCodec {
	for _, element := range f.getHWCodecSupport() {
		switch element {
		case VideoCodecN264,
			VideoCodecN264H,
			VideoCodecI264,
			VideoCodecI264C,
			VideoCodecM264,
			VideoCodecRK264:
			return element
		}
	}
	return defaultCodec
}

// Return a hardware accelerated codec for WebM if available, otherwise the default
func (f *FFMpeg) HWCodecWEBMCompatible(defaultCodec VideoCodec) VideoCodec {
	for _, element := range f.getHWCodecSupport() {
		switch element {
		case VideoCodecIVP9,
			VideoCodecVVP9:
			return element
		}
	}
	return defaultCodec
}
