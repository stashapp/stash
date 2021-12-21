//go:build ignore
// +build ignore

package main

import (
	"flag"
	"log"

	"github.com/anacrolix/ffprobe"
)

func main() {
	log.SetFlags(log.Llongfile)
	flag.Parse()
	for _, path := range flag.Args() {
		i, err := ffprobe.Probe(path)
		log.Printf("%#v %#v", i, err)
	}
}
