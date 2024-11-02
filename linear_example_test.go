package stat_test

import (
	"fmt"
	"os"

	"github.com/akramarenkov/stat"
)

func ExampleLinear() {
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
