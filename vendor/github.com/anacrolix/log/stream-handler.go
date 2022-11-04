package log

import (
	"fmt"
	"io"
	"runtime"
	"time"
)

type StreamHandler struct {
	W   io.Writer
	Fmt ByteFormatter
}

func (me StreamHandler) Handle(r Record) {
	r.Msg = r.Skip(1)
	me.W.Write(me.Fmt(r))
}

type ByteFormatter func(Record) []byte

func LineFormatter(msg Record) []byte {
	names := msg.Names
	if true || len(names) == 0 {
		var pc [1]uintptr
		msg.Callers(1, pc[:])
		names = pcNames(pc[0], names)
	}
	ret := []byte(fmt.Sprintf(
		"%s %s %s: %s",
		time.Now().Format("2006-01-02T15:04:05-0700"),
		msg.Level.LogString(),
		names,
		msg.Text(),
	))
	if ret[len(ret)-1] != '\n' {
		ret = append(ret, '\n')
	}
	return ret
}

func pcNames(pc uintptr, names []string) []string {
	if pc == 0 {
		panic(pc)
	}
	funcName, file, line := func() (string, string, int) {
		if false {
			// This seems to result in one less allocation, but doesn't handle inlining?
			func_ := runtime.FuncForPC(pc)
			file, line := func_.FileLine(pc)
			return func_.Name(), file, line
		} else {
			f, _ := runtime.CallersFrames([]uintptr{pc}).Next()
			return f.Function, f.File, f.Line
		}
	}()
	_ = file
	return append(names, fmt.Sprintf("%s:%v", funcName, line))
}
