package internal

import (
	"fmt"

	"github.com/corona10/goimagehash"
	"github.com/disintegration/imaging"

	// Decoder registrations. imaging.Open dispatches to image.Decode, which needs
	// the format decoders registered via their package init(). jpeg/png/gif come
	// in through imaging's own imports; webp/bmp are added explicitly so the
	// worker can phash the formats stash itself indexes.
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/webp"
)

// Image perceptual hashing for stash images. Unlike a SCENE (where the phash is
// computed over a 25-frame montage sampled from the video), an IMAGE's phash is
// computed directly on the decoded image itself — exactly what stash's
// pkg/hash/imagephash/phash.go does: decode → goimagehash.PerceptionHash → GetHash.
//
// Fidelity note: stash decodes the image with its own pipeline, then runs the
// same goimagehash.PerceptionHash (which internally resizes to 64x64 and DCTs).
// We rely on goimagehash's own internal resize, so the only variable is the
// decoder. imaging.Open uses the stdlib image decoders (plus the x/image webp/bmp
// decoders registered above), matching stash's decode path for the common
// formats. If a bulk run ever shows mismatches against native image phashes,
// that's where to look first — but there is no per-image montage step to diverge
// on the way video phash has.

// GenerateImagePhash decodes the image at path and returns stash's perceptual
// hash as a raw uint64. The caller formats it as lowercase hex for the API
// (SetFilePhash does that). path is the WORKER-side path — translate any
// stash-side path through the media rewriter BEFORE calling this.
func GenerateImagePhash(path string) (uint64, error) {
	// imaging.Open auto-orients via EXIF and decodes jpeg/png/gif/tiff/bmp/webp
	// (the latter two via the blank imports above). AutoOrientation keeps the
	// pixels consistent regardless of EXIF rotation flags.
	img, err := imaging.Open(path, imaging.AutoOrientation(true))
	if err != nil {
		return 0, fmt.Errorf("decode image %s: %w", path, err)
	}
	hash, err := goimagehash.PerceptionHash(img)
	if err != nil {
		return 0, fmt.Errorf("computing phash for %s: %w", path, err)
	}
	return hash.GetHash(), nil
}
