package ffmpeg

import (
	"context"
	"fmt"
)

var HWCodecSupport []StreamFormat

/*
Tests all hardware codec's to see if they work
*/
func FindHWCodecs(ctx context.Context, encoder FFMpeg) {
	for _, codec := range []StreamFormat{
		StreamFormatN264,
		StreamFormatI264,
		StreamFormatA264,
		StreamFormatV264,
		StreamFormatR264,
		StreamFormatO264,
		StreamFormatIVP9,
		StreamFormatVVP9,
	} {
		var args Args
		args = append(args, "-hide_banner")
		args = args.LogLevel(LogLevelQuiet)
		args = args.Format("lavfi")
		args = args.Input("color=c=red")
		args = args.Duration(0.1)

		args = args.VideoCodec(codec.codec)
		if len(codec.extraArgs) > 0 {
			args = append(args, codec.extraArgs...)
		}

		args = args.Format("null")
		args = args.Output("-")

		cmd := encoder.Command(ctx, args)

		if err := cmd.Run(); err == nil {
			HWCodecSupport = append(HWCodecSupport, codec)
		}
	}

	fmt.Println("Supported HW codecs:")
	for _, codec := range HWCodecSupport {
		fmt.Println("\t", codec.codec)
	}
}

func HWCodecH264Compatible() *StreamFormat {
	for _, element := range HWCodecSupport {
		switch element.codec {
		case VideoCodecLibN264,
			VideoCodecLibI264,
			VideoCodecLibA264,
			VideoCodecLibV264:
			return &element
		}
	}
	return nil
}

func HWCodecVP9Compatible() *StreamFormat {
	for _, element := range HWCodecSupport {
		switch element.codec {
		case VideoCodecIVP9,
			VideoCodecVVP9:
			return &element
		}
	}
	return nil
}
