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
	// Quantity of occurrences of a value belonging to a Span
	Quantity uint64
	// Span of values ​​for which the Quantity of occurrences is collected
	Span span.Span[Type]
}
