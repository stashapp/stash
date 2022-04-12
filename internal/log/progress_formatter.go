package log

import (
	"github.com/sirupsen/logrus"
)

type ProgressFormatter struct{}

func (f *ProgressFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	msg := []byte("Processing --> " + entry.Message + "\r")
	return msg, nil
}
