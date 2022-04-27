package file

var (
	FingerprintTypeOshash = "oshash"
	FingerprintTypeMD5    = "md5"
	FingerprintTypePhash  = "phash"
)

// Fingerprint represents a fingerprint of a file.
type Fingerprint struct {
	Type        string
	Fingerprint interface{}
}

// FingerprintCalculator calculates a fingerprint for the provided file.
type FingerprintCalculator interface {
	CalculateFingerprints(f *BaseFile, o Opener) ([]Fingerprint, error)
}
