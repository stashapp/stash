// +build ignore

package main

import "fmt"
import "time"

func main() {
	now := time.Now().Format("2006-01-02 15:04:05")

	fmt.Printf("%s", now)
}
