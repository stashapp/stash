package log

import (
	"io"

	"github.com/sirupsen/logrus"
)

type fileLogHook struct {
	Writer    io.Writer
	Formatter logrus.Formatter
}

func (hook *fileLogHook) Fire(entry *logrus.Entry) error {
	line, err := hook.Formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write(line)
	return err
}

func (hook *fileLogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
