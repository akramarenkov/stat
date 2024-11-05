package stat

import (
	"github.com/akramarenkov/span"
	"golang.org/x/exp/constraints"
)

func search[Type constraints.Integer](item, target Item[Type]) int {
	return span.SearchInc(item.Span, target.Span)
}
