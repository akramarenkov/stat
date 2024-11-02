package stat_test

import (
	"fmt"
	"os"

	"github.com/akramarenkov/stat"
)

func ExampleStat() {
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
