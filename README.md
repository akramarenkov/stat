# Stat

[![Go Reference](https://pkg.go.dev/badge/github.com/akramarenkov/stat.svg)](https://pkg.go.dev/github.com/akramarenkov/stat)
[![Go Report Card](https://goreportcard.com/badge/github.com/akramarenkov/stat)](https://goreportcard.com/report/github.com/akramarenkov/stat)
[![Coverage Status](https://coveralls.io/repos/github/akramarenkov/stat/badge.svg)](https://coveralls.io/github/akramarenkov/stat)

## Purpose

Library that provides to collect and display the quantity of occurrences of values ​​in given spans

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
    stat, err := stat.NewLinear(1, 100, 20)
    if err != nil {
        panic(err)
    }

    stat.Inc(0)

    stat.Inc(1)
    stat.Inc(20)

    stat.Inc(21)
    stat.Inc(22)
    stat.Inc(40)

    stat.Inc(41)
    stat.Inc(42)
    stat.Inc(59)
    stat.Inc(60)

    stat.Inc(61)
    stat.Inc(62)
    stat.Inc(80)

    stat.Inc(81)
    stat.Inc(100)

    stat.Inc(101)

    fmt.Println(stat.Graph(os.Stderr))
    // Output:
    // <nil>
}
```
