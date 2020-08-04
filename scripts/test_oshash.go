// +build ignore

package main

import (
	"fmt"
	"os"

	"github.com/stashapp/stash/pkg/utils"
)

func main() {
	hash, err := utils.OSHashFromFilePath(os.Args[1])

	if err != nil {
		panic(err)
	}

	fmt.Println(hash)
}
