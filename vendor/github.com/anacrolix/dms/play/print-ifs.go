//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"net"
)

func main() {
	ifs, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, if_ := range ifs {
		fmt.Printf("%#v\n", if_)
		addrs, err := if_.Addrs()
		if err != nil {
			panic(err)
		}
		for _, addr := range addrs {
			fmt.Printf("\t%s %s\n", addr.Network(), addr)
		}
		mcastAddrs, err := if_.MulticastAddrs()
		if err != nil {
			panic(err)
		}
		for _, addr := range mcastAddrs {
			fmt.Printf("\t%s\n", addr)
		}
	}
}
