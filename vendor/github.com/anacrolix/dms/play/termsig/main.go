package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {
	c := make(chan os.Signal, 0x100)
	signal.Notify(c)
	for i := range c {
		fmt.Println(i)
	}
}
