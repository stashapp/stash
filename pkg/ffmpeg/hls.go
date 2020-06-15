package ffmpeg

import (
	"fmt"
	"io"
)

const hlsSegmentLength = 10.0

func WriteHLSPlaylist(probeResult VideoFile, baseUrl string, w io.Writer) {
	fmt.Fprint(w, "#EXTM3U\n")
	fmt.Fprint(w, "#EXT-X-VERSION:3\n")
	fmt.Fprint(w, "#EXT-X-MEDIA-SEQUENCE:0\n")
	fmt.Fprint(w, "#EXT-X-ALLOW-CACHE:YES\n")
	fmt.Fprintf(w, "#EXT-X-TARGETDURATION:%d\n", int(hlsSegmentLength))
	fmt.Fprint(w, "#EXT-X-PLAYLIST-TYPE:VOD\n")

	duration := probeResult.Duration

	leftover := duration
	upTo := 0.0

	for leftover > 0 {
		thisLength := hlsSegmentLength
		if leftover < thisLength {
			thisLength = leftover
		}

		fmt.Fprintf(w, "#EXTINF: %f,\n", thisLength)
		fmt.Fprintf(w, "%s&start=%f\n", baseUrl, upTo)

		leftover -= thisLength
		upTo += thisLength
	}

	fmt.Fprint(w, "#EXT-X-ENDLIST\n")
}
