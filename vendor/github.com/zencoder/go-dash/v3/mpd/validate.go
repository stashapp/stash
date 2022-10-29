package mpd

// Validate checks for incomplete MPD object
func (m *MPD) Validate() error {
	if m.Profiles == nil {
		return ErrNoDASHProfileSet
	}
	return nil
}
