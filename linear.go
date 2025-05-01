package stat

import (
	"github.com/akramarenkov/safe"
	"github.com/akramarenkov/span"
	"golang.org/x/exp/constraints"
)

// Creates a linear statistics whose items have the specified width.
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

// Creates a linear statistics with the specified quantity of items.
func NewLinearQ[Type constraints.Integer](lower, upper, quantity Type) (*Stat[Type], error) {
	if lower > upper {
		return nil, ErrLowerGreaterUpper
	}

	if quantity < 0 {
		return nil, ErrItemsQuantityNegative
	}

	if quantity == 0 {
		return nil, ErrItemsQuantityZero
	}

	if quantity == 1 {
		spans := []span.Span[Type]{{Begin: lower, End: upper}}

		predictor := func(Type) uint64 {
			return 0
		}

		return New(spans, predictor)
	}

	width, err := safe.AddOneSubDiv(upper, lower, quantity)
	if err != nil {
		// Given the checks above, an error can only occur for signed types,
		// only if the lower and upper values are equal to the minimum and maximum
		// values for the type used, and only if the quantity value is two
		spans := []span.Span[Type]{
			{Begin: lower, End: -Type(1)},
			{Begin: 0, End: upper},
		}

		predictor := func(value Type) uint64 {
			if value >= 0 {
				return 1
			}

			return 0
		}

		return New(spans, predictor)
	}

	return NewLinear(lower, upper, width)
}
