package log

type Handler interface {
	Handle(Record)
}

type Record struct {
	Msg
	Level Level
	Names []string
}
