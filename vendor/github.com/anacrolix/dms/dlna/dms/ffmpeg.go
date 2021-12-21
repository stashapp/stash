package dms

import (
	"os/exec"
	"runtime"
	"syscall"
)

func suppressFFmpegProbeDataErrors(_err error) (err error) {
	if _err == nil {
		return
	}
	err = _err
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		return
	}
	waitStat, ok := exitErr.Sys().(syscall.WaitStatus)
	if !ok {
		return
	}
	code := waitStat.ExitStatus()
	if runtime.GOOS == "windows" {
		if code == -1094995529 {
			err = nil
		}
	} else if code == 183 {
		err = nil
	}
	return
}
