package ffmpeg2

type Error struct {
	Cause  error
	Output string
}

func (e Error) Error() string {
	return e.Cause.Error()
}
