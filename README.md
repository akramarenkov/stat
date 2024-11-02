# Stat

[![Go Reference](https://pkg.go.dev/badge/github.com/akramarenkov/stat.svg)](https://pkg.go.dev/github.com/akramarenkov/stat)
[![Go Report Card](https://goreportcard.com/badge/github.com/akramarenkov/stat)](https://goreportcard.com/report/github.com/akramarenkov/stat)
[![Coverage Status](https://coveralls.io/repos/github/akramarenkov/stat/badge.svg)](https://coveralls.io/github/akramarenkov/stat)

## Purpose

Library that allows you to collect and display the quantity of occurrences of values ​​in given spans

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
    linear, err := stat.NewLinear(1, 100, 20)
    if err != nil {
        panic(err)
    }

    linear.Add(0)

    linear.Add(1)
    linear.Add(20)

    linear.Add(21)
    linear.Add(22)
    linear.Add(40)

    linear.Add(41)
    linear.Add(42)
    linear.Add(59)
    linear.Add(60)

    linear.Add(61)
    linear.Add(62)
    linear.Add(80)

    linear.Add(81)
    linear.Add(100)

    linear.Add(101)

    fmt.Println(linear.Stat().Graph(os.Stderr))

    // Output:
    // <nil>
}
```
