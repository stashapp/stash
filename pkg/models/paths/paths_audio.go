package paths

type audioPaths struct {
	generatedPaths
}

func newAudioPaths(p Paths) *audioPaths {
	sp := audioPaths{
		generatedPaths: *p.Generated,
	}
	return &sp
}

func (sp *audioPaths) GetStreamPath(audioPath string, checksum string) string {
	// No Transcodes at this time
	return audioPath
}
