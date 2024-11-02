package stat

import (
	"github.com/akramarenkov/safe"
	"github.com/akramarenkov/span"
	"golang.org/x/exp/constraints"
)

// Creates linear statistics.
func NewLinear[Type constraints.Integer](lower, upper, width Type) (*Stat[Type], error) {
	if lower > upper {
		return nil, ErrLowerGreaterUpper
	}

	spans, err := span.Linear(lower, upper, width)
	if err != nil {
		return nil, err
	}

	predictor := func(value Type) uint64 {
		return safe.Dist(value, lower) / uint64(width)
	}

	return New(spans, predictor)
}
