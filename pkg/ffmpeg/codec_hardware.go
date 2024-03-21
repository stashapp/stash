package ffmpeg

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/stashapp/stash/pkg/logger"
)

var (
	// Hardware codec's
	VideoCodecN264 VideoCodec = "h264_nvenc"
	VideoCodecI264 VideoCodec = "h264_qsv"
	VideoCodecA264 VideoCodec = "h264_amf"
	VideoCodecM264 VideoCodec = "h264_videotoolbox"
	VideoCodecV264 VideoCodec = "h264_vaapi"
	VideoCodecR264 VideoCodec = "h264_v4l2m2m"
	VideoCodecO264 VideoCodec = "h264_omx"
	VideoCodecIVP9 VideoCodec = "vp9_qsv"
	VideoCodecVVP9 VideoCodec = "vp9_vaapi"
	VideoCodecVVPX VideoCodec = "vp8_vaapi"
)

// Tests all (given) hardware codec's
func (f *FFMpeg) InitHWSupport(ctx context.Context) {
	var hwCodecSupport []VideoCodec

	for _, codec := range []VideoCodec{
		VideoCodecN264,
		VideoCodecI264,
		VideoCodecV264,
		VideoCodecR264,
		VideoCodecIVP9,
		VideoCodecVVP9,
	} {
		var args Args
		args = append(args, "-hide_banner")
		args = args.LogLevel(LogLevelWarning)
		args = f.hwDeviceInit(args, codec)
		args = args.Format("lavfi")
		args = args.Input("color=c=red")
		args = args.Duration(0.1)

		videoFilter := f.hwFilterInit(codec)
		// Test scaling
		videoFilter = videoFilter.ScaleDimensions(-2, 160)
		videoFilter = f.hwCodecFilter(videoFilter, codec)
		args = append(args, CodecInit(codec)...)
		args = args.VideoFilter(videoFilter)

		args = args.Format("null")
		args = args.Output("-")

		cmd := f.Command(ctx, args)

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		if err := cmd.Start(); err != nil {
			logger.Debugf("[InitHWSupport] error starting command: %v", err)
			continue
		}

		if err := cmd.Wait(); err != nil {
			errOutput := stderr.String()

			if len(errOutput) == 0 {
				errOutput = err.Error()
			}

			logger.Debugf("[InitHWSupport] Codec %s not supported. Error output:\n%s", codec, errOutput)
		} else {
			hwCodecSupport = append(hwCodecSupport, codec)
		}
	}

	outstr := "[InitHWSupport] Supported HW codecs:\n"
	for _, codec := range hwCodecSupport {
		outstr += fmt.Sprintf("\t%s\n", codec)
	}
	logger.Info(outstr)

	f.hwCodecSupport = hwCodecSupport
}

// Prepend input for hardware encoding only
func (f *FFMpeg) hwDeviceInit(args Args, codec VideoCodec) Args {
	switch codec {
	case VideoCodecN264:
		args = append(args, "-hwaccel_device")
		args = append(args, "0")
	case VideoCodecV264,
		VideoCodecVVP9:
		args = append(args, "-vaapi_device")
		args = append(args, "/dev/dri/renderD128")
	case VideoCodecI264,
		VideoCodecIVP9:
		args = append(args, "-init_hw_device")
		args = append(args, "qsv=hw")
		args = append(args, "-filter_hw_device")
		args = append(args, "hw")
	}

	return args
}

// Initialise a video filter for HW encoding
func (f *FFMpeg) hwFilterInit(codec VideoCodec) VideoFilter {
	var videoFilter VideoFilter
	switch codec {
	case VideoCodecV264,
		VideoCodecVVP9:
		videoFilter = videoFilter.Append("format=nv12")
		videoFilter = videoFilter.Append("hwupload")
	case VideoCodecN264:
		videoFilter = videoFilter.Append("format=nv12")
		videoFilter = videoFilter.Append("hwupload_cuda")
	case VideoCodecI264,
		VideoCodecIVP9:
		videoFilter = videoFilter.Append("hwupload=extra_hw_frames=64")
		videoFilter = videoFilter.Append("format=qsv")
	}

	return videoFilter
}

// Replace video filter scaling with hardware scaling for full hardware transcoding
func (f *FFMpeg) hwCodecFilter(args VideoFilter, codec VideoCodec) VideoFilter {
	sargs := string(args)

	if strings.Contains(sargs, "scale=") {
		switch codec {
		case VideoCodecN264:
			args = VideoFilter(strings.Replace(sargs, "scale=", "scale_cuda=", 1))
		case VideoCodecV264,
			VideoCodecVVP9:
			args = VideoFilter(strings.Replace(sargs, "scale=", "scale_vaapi=", 1))
		case VideoCodecI264,
			VideoCodecIVP9:
			// BUG: [scale_qsv]: Size values less than -1 are not acceptable.
			// Fix: Replace all instances of -2 with -1 in a scale operation
			re := regexp.MustCompile(`(scale=)([\d:]*)(-2)(.*)`)
			sargs = re.ReplaceAllString(sargs, "scale=$2-1$4")
			args = VideoFilter(strings.Replace(sargs, "scale=", "scale_qsv=", 1))
		}
	}

	return args
}

// Returns the max resolution for a given codec, or a default
func (f *FFMpeg) hwCodecMaxRes(codec VideoCodec, dW int, dH int) (int, int) {
	if codec == VideoCodecN264 {
		return 4096, 4096
	}

	return dW, dH
}

// Return a maxres filter
func (f *FFMpeg) hwMaxResFilter(codec VideoCodec, width int, height int, max int) VideoFilter {
	videoFilter := f.hwFilterInit(codec)
	maxWidth, maxHeight := f.hwCodecMaxRes(codec, width, height)
	videoFilter = videoFilter.ScaleMaxLM(width, height, max, maxWidth, maxHeight)
	return f.hwCodecFilter(videoFilter, codec)
}

// Return if a hardware accelerated for HLS is available
func (f *FFMpeg) hwCodecHLSCompatible() *VideoCodec {
	for _, element := range f.hwCodecSupport {
		switch element {
		case VideoCodecN264,
			VideoCodecI264,
			VideoCodecV264,
			VideoCodecR264:
			return &element
		}
	}
	return nil
}

// Return if a hardware accelerated codec for MP4 is available
func (f *FFMpeg) hwCodecMP4Compatible() *VideoCodec {
	for _, element := range f.hwCodecSupport {
		switch element {
		case VideoCodecN264,
			VideoCodecI264:
			return &element
		}
	}
	return nil
}

// Return if a hardware accelerated codec for WebM is available
func (f *FFMpeg) hwCodecWEBMCompatible() *VideoCodec {
	for _, element := range f.hwCodecSupport {
		switch element {
		case VideoCodecIVP9,
			VideoCodecVVP9:
			return &element
		}
	}
	return nil
}
