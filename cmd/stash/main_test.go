package main

import (
	"os"
	"testing"

	"github.com/stashapp/stash/internal/log"
	"github.com/stashapp/stash/pkg/logger"
)

func TestExitError(t *testing.T) {
	origCode := exitCode
	origLogger := logger.Logger
	defer func() {
		exitCode = origCode
		logger.Logger = origLogger
	}()

	l := log.NewLogger()
	l.Init("", true, "Error", 0)
	logger.Logger = l

	exitCode = 0
	exitError(os.ErrInvalid)
	if exitCode != 1 {
		t.Errorf("exitError should set exitCode to 1, got %d", exitCode)
	}
}

func TestRecoverPanic(t *testing.T) {
	origCode := exitCode
	origLogger := logger.Logger
	defer func() {
		exitCode = origCode
		logger.Logger = origLogger
	}()

	l := log.NewLogger()
	l.Init("", true, "Error", 0)
	logger.Logger = l

	exitCode = 0
	func() {
		defer recoverPanic()
		panic("test panic")
	}()
	if exitCode != 1 {
		t.Errorf("recoverPanic should set exitCode to 1 after a panic, got %d", exitCode)
	}
}

func TestInitLogTemp(t *testing.T) {
	origLogger := logger.Logger
	defer func() { logger.Logger = origLogger }()

	logger.Logger = nil
	l := initLogTemp()
	if l == nil {
		t.Fatal("initLogTemp should return a non-nil logger")
	}
	if logger.Logger == nil {
		t.Error("initLogTemp should set the global logger.Logger")
	}
}
