package ffmpeg

import (
	"fmt"
	"io"
	"strings"
)

const hlsSegmentLength = 10.0

// WriteHLSPlaylist writes a HLS playlist to w using baseUrl as the base URL for TS streams.
func WriteHLSPlaylist(duration float64, baseUrl string, w io.Writer) {
	fmt.Fprint(w, "#EXTM3U\n")
	fmt.Fprint(w, "#EXT-X-VERSION:3\n")
	fmt.Fprint(w, "#EXT-X-MEDIA-SEQUENCE:0\n")
	fmt.Fprint(w, "#EXT-X-ALLOW-CACHE:YES\n")
	fmt.Fprintf(w, "#EXT-X-TARGETDURATION:%d\n", int(hlsSegmentLength))
	fmt.Fprint(w, "#EXT-X-PLAYLIST-TYPE:VOD\n")

	leftover := duration
	upTo := 0.0

	i := strings.LastIndex(baseUrl, ".m3u8")
	tsURL := baseUrl[0:i] + ".ts"

	for leftover > 0 {
		thisLength := hlsSegmentLength
		if leftover < thisLength {
			thisLength = leftover
		}

		fmt.Fprintf(w, "#EXTINF: %f,\n", thisLength)
		fmt.Fprintf(w, "%s?start=%f\n", tsURL, upTo)

		leftover -= thisLength
		upTo += thisLength
	}

	fmt.Fprint(w, "#EXT-X-ENDLIST\n")
}
