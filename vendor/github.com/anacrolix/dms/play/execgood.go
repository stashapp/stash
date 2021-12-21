//go:build ignore
// +build ignore

package main

import (
	"io"
	"log"
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
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	r, w := io.Pipe()
	go func() {
		io.Copy(w, out)
		out.Close()
		w.Close()
		log.Println(cmd.Wait())
	}()
	time.Sleep(10 * time.Millisecond)
	if _, err := io.Copy(os.Stdout, r); err != nil {
		panic(err)
	}
}
