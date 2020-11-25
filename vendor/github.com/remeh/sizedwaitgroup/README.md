# SizedWaitGroup

[![GoDoc](https://godoc.org/github.com/remeh/sizedwaitgroup?status.svg)](https://godoc.org/github.com/remeh/sizedwaitgroup)

`SizedWaitGroup` has the same role and API as `sync.WaitGroup` but it adds a limit of the amount of goroutines started concurrently.

`SizedWaitGroup` adds the feature of limiting the maximum number of concurrently started routines. It could for example be used to start multiples routines querying a database but without sending too much queries in order to not overload the given database.

# Example

```
package main

import (
        "fmt"
        "math/rand"
        "time"

        "github.com/remeh/sizedwaitgroup"
)

func main() {
        rand.Seed(time.Now().UnixNano())

        // Typical use-case:
        // 50 queries must be executed as quick as possible
        // but without overloading the database, so only
        // 8 routines should be started concurrently.
        swg := sizedwaitgroup.New(8)
        for i := 0; i < 50; i++ {
                swg.Add()
                go func(i int) {
                        defer swg.Done()
                        query(i)
                }(i)
        }

        swg.Wait()
}

func query(i int) {
        fmt.Println(i)
        ms := i + 500 + rand.Intn(500)
        time.Sleep(time.Duration(ms) * time.Millisecond)
}
```

# License

MIT

# Copyright

Rémy Mathieu © 2016
