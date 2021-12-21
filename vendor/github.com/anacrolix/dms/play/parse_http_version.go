//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"net/http"
	"strings"
)

func main() {
	fmt.Println(http.ParseHTTPVersion(strings.TrimSpace("HTTP/1.1 ")))
}
