package stat

import (
	"github.com/akramarenkov/safe"
	"github.com/akramarenkov/span"
	"golang.org/x/exp/constraints"
)

// Statistics for a linear sequence of numbers linearly divided into spans.
type Linear[Type constraints.Integer] struct {
	lower Type
	upper Type
	width Type

	stat *Stat[Type]
}

// Creates an instance of linear statistics.
func NewLinear[Type constraints.Integer](lower, upper, width Type) (*Linear[Type], error) {
	if lower > upper {
		return nil, ErrLowerGreaterUpper
	}

	spans, err := span.Linear(lower, upper, width)
	if err != nil {
		return nil, err
	}

	lnr := &Linear[Type]{
		lower: lower,
		upper: upper,
		width: width,

		stat: New(spans),
	}

	return lnr, nil
}

// Takes into account the value in statistics.
func (lnr *Linear[Type]) Add(value Type) {
	if value < lnr.lower {
		lnr.stat.IncNegInf()
		return
	}

	if value > lnr.upper {
		lnr.stat.IncPosInf()
		return
	}

	id := safe.Dist(value, lnr.lower) / uint64(lnr.width)

	lnr.stat.Inc(id)
}

// Returns statistics.
func (lnr *Linear[Type]) Stat() *Stat[Type] {
	return lnr.stat
}
