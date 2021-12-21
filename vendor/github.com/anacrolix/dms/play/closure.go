//go:build ignore
// +build ignore

package main

import (
	"fmt"
)

func main() {
	ch := make(chan struct{})
	for i := 0; i < 5; i++ {
		j := i
		go func() {
			fmt.Print(j)
			ch <- struct{}{}
		}()
	}
	for i := 5; i < 10; i++ {
		go func() {
			fmt.Print(i)
			ch <- struct{}{}
		}()
	}
	for i := 0; i < 10; i++ {
		<-ch
	}
}
