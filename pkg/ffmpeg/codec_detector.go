package ffmpeg

import (
	"context"
	"fmt"
)

var HWCodecSupport []VideoCodec

/*
Tests all hardware codec's to see if they work
*/
func FindHWCodecs(encoder FFMpeg, ctx context.Context) {
	//TODO: Should probably do a support matrix
	for _, codec := range []VideoCodec{
		VideoCodecLibN264,
		VideoCodecLibI264,
		VideoCodecLibA264,
		VideoCodecLibV264,
		VideoCodecVVP9,
		VideoCodecIVP9,
		VideoCodecVVPX,
	} {
		var args Args
		args = append(args, "-hide_banner")
		args = args.LogLevel(LogLevelQuiet)
		args = args.Format("lavfi")
		args = args.Input("color=c=red")
		args = args.Duration(0.1)
		args = args.VideoCodec(codec)

		args = args.Format("null")
		args = args.Output("-")

		cmd := encoder.Command(ctx, args)

		if err := cmd.Run(); err == nil {
			HWCodecSupport = append(HWCodecSupport, codec)
		}
	}

	fmt.Println("Supported HW codecs:")
	for _, codec := range HWCodecSupport {
		fmt.Println("\t", codec)
	}
}

func HWCodecCompatible(c VideoCodec) bool {
	for _, element := range HWCodecSupport {
		if element == c {
			return true
		}
	}
	return false
}
