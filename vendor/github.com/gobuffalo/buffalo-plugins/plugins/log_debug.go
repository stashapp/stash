//+build debug

package plugins

import (
	"fmt"
	"time"
)

func log(name string, fn func() error) error {
	start := time.Now()
	defer fmt.Println(name, time.Now().Sub(start))
	return fn()
}
