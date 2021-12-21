//go:build ignore
// +build ignore

package main

import (
	"io"
	"os"
	"os/exec"
	"time"
)

func main() {
	cmd := exec.Command("ls")
	out, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	if err = cmd.Start(); err != nil {
		panic(err)
	}
	go cmd.Wait()
	time.Sleep(10 * time.Millisecond)
	_, err = io.Copy(os.Stdout, out)
	if err != nil {
		panic(err)
	}
}
