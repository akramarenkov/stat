package stat_test

import (
	"fmt"
	"os"

	"github.com/akramarenkov/stat"
)

func ExampleStat() {
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
