package stat

import (
	"github.com/akramarenkov/span"
	"golang.org/x/exp/constraints"
)

// A function used to determine (at least approximately) the index of a span in a
// list of spans to which a value belongs.
type Predictor[Type constraints.Integer] func(value Type) uint64

// Item of statistics.
type Item[Type constraints.Integer] struct {
	// Kind (purpose) of item
	Kind ItemKind
	// Quantity of occurrences of a value belonging to a Span
	Quantity uint64
	// Span of values ​​for which the Quantity of occurrences is collected
	Span span.Span[Type]
}

// Kind (purpose) of item.
type ItemKind int

const (
	ItemKindRegular ItemKind = iota + 1
	ItemKindNegInf
	ItemKindPosInf
	ItemKindMissed
)

func (ik ItemKind) String() string {
	switch ik {
	case ItemKindRegular:
		return "regular"
	case ItemKindNegInf:
		return "-Inf"
	case ItemKindPosInf:
		return "+Inf"
	case ItemKindMissed:
		return "missed"
	}

	return "unexpected"
}
