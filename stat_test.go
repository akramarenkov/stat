package stat

import (
	"io"
	"math"
	"testing"

	"github.com/akramarenkov/span"
	"github.com/stretchr/testify/require"
)

func TestStatPanic(t *testing.T) {
	stat := New[int](nil)

	require.Panics(t, func() { stat.IncNegInf() })
	require.Panics(t, func() { stat.IncPosInf() })
	require.Panics(t, func() { stat.Inc(0) })

	require.Len(t, stat.Items(), 0)
}

func TestStatGraphError(t *testing.T) {
	stat := New([]span.Span[int]{{Begin: 0, End: 0}})

	stat.items[2].Quantity = math.MaxUint64
	require.Error(t, stat.Graph(io.Discard))

	stat.items[1].Quantity = math.MaxUint64
	require.Error(t, stat.Graph(io.Discard))

	stat.items[0].Quantity = math.MaxUint64
	require.Error(t, stat.Graph(io.Discard))
}

func BenchmarkStat(b *testing.B) {
	stat := New([]span.Span[int]{{Begin: 0, End: 0}})

	for range b.N {
		stat.IncNegInf()
		stat.IncPosInf()
		stat.Inc(0)
	}
}
