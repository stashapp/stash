package errx

import "fmt"

// go2 errors
type Wrapper interface {
	Unwrap() error
}

// pkg/errors
type Causer interface {
	Cause() error
}

func Unwrap(err error) error {
	switch e := err.(type) {
	case Wrapper:
		return e.Unwrap()
	case Causer:
		return e.Cause()
	}
	return err
}

var Cause = Unwrap

func Wrap(err error, msg string) error {
	return wrapped{
		err: err,
		msg: msg,
	}
}

type wrapped struct {
	err error
	msg string
}

func (w wrapped) Error() string {
	return fmt.Sprintf("%s: %s", w.msg, w.err)
}

func (w wrapped) Unwrap() error {
	return w.err
}

func (w wrapped) Cause() error {
	return w.err
}
