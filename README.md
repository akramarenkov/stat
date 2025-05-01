# Stat

[![Go Reference](https://pkg.go.dev/badge/github.com/akramarenkov/stat.svg)](https://pkg.go.dev/github.com/akramarenkov/stat)
[![Go Report Card](https://goreportcard.com/badge/github.com/akramarenkov/stat)](https://goreportcard.com/report/github.com/akramarenkov/stat)
[![Coverage Status](https://coveralls.io/repos/github/akramarenkov/stat/badge.svg)](https://coveralls.io/github/akramarenkov/stat)

## Purpose

Library that provides to collect and display the quantity of occurrences of
 values in given spans

## Usage

Example:

```go
package main

import (
    "fmt"
    "os"

    "github.com/akramarenkov/stat"
)

func main() {
    sts, err := stat.NewLinear(1, 100, 20)
    if err != nil {
        panic(err)
    }

    sts.Inc(0)

    sts.Inc(1)
    sts.Inc(20)

    sts.Inc(21)
    sts.Inc(22)
    sts.Inc(40)

    sts.Inc(41)
    sts.Inc(42)
    sts.Inc(59)
    sts.Inc(60)

    sts.Inc(61)
    sts.Inc(62)
    sts.Inc(80)

    sts.Inc(81)
    sts.Inc(100)

    sts.Inc(101)

    fmt.Println(sts.Graph(os.Stderr))
    // Output:
    // <nil>
}
```
