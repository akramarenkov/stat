package stat

import (
	"golang.org/x/exp/constraints"
)

func searcher[Type constraints.Integer](item, target Item[Type]) int {
	switch {
	case item.Span.End < target.Span.Begin:
		return -1
	case item.Span.Begin > target.Span.Begin:
		return 1
	}

	return 0
}
