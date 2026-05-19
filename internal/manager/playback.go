package manager

import (
	"context"
	"net/url"
	"strings"

	"github.com/stashapp/stash/pkg/models"
)

// MarkDefaultStream inspects endpoints against the stored playback defaults for
// the given user-agent string and sets Default on the matching endpoint. If no
// rule matches, all Default fields are left nil.
func MarkDefaultStream(ctx context.Context, endpoints []*SceneStreamEndpoint, userAgent string, store models.PlaybackDefaultReader) error {
	if userAgent == "" || store == nil {
		return nil
	}

	match, err := store.FindByUserAgent(ctx, userAgent)
	if err != nil || match == nil {
		return err
	}

	quality := ""
	if match.Quality != nil {
		quality = string(*match.Quality)
	}

	for _, ep := range endpoints {
		if matchesRule(ep, match.StreamType, quality) {
			t := true
			ep.Default = &t
			return nil
		}
	}

	return nil
}

func matchesRule(ep *SceneStreamEndpoint, streamType models.PlaybackStreamType, quality string) bool {
	u, err := url.Parse(ep.URL)
	if err != nil {
		return false
	}

	var pathMatch bool
	switch streamType {
	case models.PlaybackStreamTypeDirect:
		pathMatch = strings.HasSuffix(u.Path, "/stream")
	case models.PlaybackStreamTypeMP4:
		pathMatch = strings.HasSuffix(u.Path, "/stream.mp4")
	case models.PlaybackStreamTypeWEBM:
		pathMatch = strings.HasSuffix(u.Path, "/stream.webm")
	case models.PlaybackStreamTypeMKV:
		pathMatch = strings.HasSuffix(u.Path, "/stream.mkv")
	case models.PlaybackStreamTypeHLS:
		pathMatch = strings.HasSuffix(u.Path, "/stream.m3u8")
	case models.PlaybackStreamTypeDASH:
		pathMatch = strings.HasSuffix(u.Path, "/stream.mpd")
	}

	if !pathMatch {
		return false
	}

	epQuality := u.Query().Get("resolution")

	if quality == "" || quality == string(models.StreamingResolutionEnumOriginal) {
		return epQuality == "" || epQuality == string(models.StreamingResolutionEnumOriginal)
	}

	return epQuality == quality
}
