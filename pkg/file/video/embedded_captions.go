package video

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/models"
)

// embeddedSubsCache memoises the (relatively expensive) ffprobe of a video file
// for embedded subtitle detection. The scene caption resolver runs on each scene
// detail load, so caching avoids re-probing unchanged files. Entries are keyed by
// path and invalidated when the file's size or modification time changes.
type embeddedSubsCacheEntry struct {
	size    int64
	modTime time.Time
	subs    []EmbeddedSubtitle
}

var embeddedSubsCache *lru.Cache[string, embeddedSubsCacheEntry]

func init() {
	// size is a best-effort bound; errors only occur for a non-positive size.
	embeddedSubsCache, _ = lru.New[string, embeddedSubsCacheEntry](2048)
}

// imageSubtitleCodecs are subtitle codecs that are bitmap (image) based and so
// cannot be converted to text-based WebVTT. They are skipped when detecting
// embedded captions.
var imageSubtitleCodecs = map[string]bool{
	"hdmv_pgs_subtitle": true,
	"pgssub":            true,
	"dvd_subtitle":      true,
	"dvdsub":            true,
	"dvb_subtitle":      true,
	"dvbsub":            true,
	"dvb_teletext":      true,
	"xsub":              true,
}

// unknownSubtitleLangs are language tag values that should be treated as
// "no language" (mapped to LangUnknown).
var unknownSubtitleLangs = map[string]bool{
	"":        true,
	"und":     true,
	"unknown": true,
	"none":    true,
	"mis":     true,
	"mul":     true,
	"zxx":     true,
}

// EmbeddedSubtitle describes a text subtitle stream embedded in a video file.
type EmbeddedSubtitle struct {
	// Index is the absolute ffmpeg stream index within the container.
	Index int
	// LanguageCode is the normalised ISO-639 language code, or LangUnknown.
	LanguageCode string
}

// normalizeSubtitleLang converts an ffprobe language tag into the language code
// used by stash's caption system, returning LangUnknown when absent or invalid.
func normalizeSubtitleLang(raw string) string {
	lang := strings.ToLower(strings.TrimSpace(raw))
	if unknownSubtitleLangs[lang] {
		return LangUnknown
	}
	// ISO-639 codes are 2 or 3 letters; validate against the language database.
	if (len(lang) == 2 || len(lang) == 3) && IsValidLanguage(lang) {
		return lang
	}
	return LangUnknown
}

// GetEmbeddedSubtitles probes the given video file and returns its embedded
// text subtitle streams, de-duplicated to one per language (preferring the
// stream marked default, then forced, then the lowest stream index). Image-based
// subtitle streams are skipped as they cannot be converted to text.
func GetEmbeddedSubtitles(ffprobe *ffmpeg.FFProbe, videoPath string) ([]EmbeddedSubtitle, error) {
	if ffprobe == nil {
		return nil, nil
	}

	// return a cached result if the file is unchanged since it was last probed
	var stat os.FileInfo
	if fi, err := os.Stat(videoPath); err == nil {
		stat = fi
		if entry, ok := embeddedSubsCache.Get(videoPath); ok &&
			entry.size == fi.Size() && entry.modTime.Equal(fi.ModTime()) {
			return entry.subs, nil
		}
	}

	probeResult, err := ffprobe.NewVideoFile(videoPath)
	if err != nil {
		return nil, err
	}

	subs := selectEmbeddedSubtitles(probeResult.GetSubtitleStreams())

	if stat != nil {
		embeddedSubsCache.Add(videoPath, embeddedSubsCacheEntry{
			size:    stat.Size(),
			modTime: stat.ModTime(),
			subs:    subs,
		})
	}

	return subs, nil
}

// selectEmbeddedSubtitles filters and de-duplicates subtitle streams into the
// embedded captions stash will expose: image-based codecs are dropped, and at
// most one stream is kept per language (preferring default, then forced, then
// the lowest stream index).
func selectEmbeddedSubtitles(streams []*ffmpeg.FFProbeStream) []EmbeddedSubtitle {
	// order so that the preferred stream for a language is encountered first
	sort.SliceStable(streams, func(i, j int) bool {
		di, dj := streams[i].Disposition, streams[j].Disposition
		if di.Default != dj.Default {
			return di.Default > dj.Default
		}
		if di.Forced != dj.Forced {
			return di.Forced > dj.Forced
		}
		return streams[i].Index < streams[j].Index
	})

	seen := make(map[string]bool)
	var ret []EmbeddedSubtitle
	for _, s := range streams {
		if imageSubtitleCodecs[strings.ToLower(s.CodecName)] {
			continue
		}

		lang := normalizeSubtitleLang(s.Tags.Language)
		// stash stores a single caption per language; keep the first (preferred) one
		if seen[lang] {
			continue
		}
		seen[lang] = true

		ret = append(ret, EmbeddedSubtitle{
			Index:        s.Index,
			LanguageCode: lang,
		})
	}

	return ret
}

// EmbeddedCaptions probes the given video file and returns its embedded subtitle
// tracks as VideoCaption entries (with CaptionType == models.CaptionTypeEmbedded)
// suitable for inclusion in a scene's caption list.
func EmbeddedCaptions(ffprobe *ffmpeg.FFProbe, videoPath string) ([]*models.VideoCaption, error) {
	subs, err := GetEmbeddedSubtitles(ffprobe, videoPath)
	if err != nil {
		return nil, err
	}

	ret := make([]*models.VideoCaption, 0, len(subs))
	for _, s := range subs {
		ret = append(ret, &models.VideoCaption{
			LanguageCode: s.LanguageCode,
			CaptionType:  models.CaptionTypeEmbedded,
		})
	}
	return ret, nil
}

// ExtractEmbeddedSubtitle extracts the subtitle stream at the given absolute
// stream index from the input video and writes it to outputPath as WebVTT.
func ExtractEmbeddedSubtitle(ctx context.Context, encoder *ffmpeg.FFMpeg, input string, streamIndex int, outputPath string) error {
	if encoder == nil {
		return fmt.Errorf("ffmpeg not configured")
	}

	var args ffmpeg.Args
	args = args.LogLevel(ffmpeg.LogLevelError)
	args = args.Overwrite()
	args = args.Input(input)
	args = append(args, "-map", fmt.Sprintf("0:%d", streamIndex))
	args = append(args, "-f", "webvtt")
	args = args.Output(outputPath)

	return encoder.Generate(ctx, args)
}
