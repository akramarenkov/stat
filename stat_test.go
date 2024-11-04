package stat

import (
	"io"
	"math"
	"os"
	"testing"

	"github.com/akramarenkov/safe"
	"github.com/akramarenkov/span"
	"github.com/stretchr/testify/require"
)

func TestStatLinear(t *testing.T) {
	stat, err := NewLinear(1, 100, 39)
	require.NoError(t, err)

	stat.Inc(-102)
	stat.Inc(-101)

	for value := range safe.Inc(1, 50) {
		stat.Inc(value)
	}

	for value := range safe.Inc(51, 100) {
		stat.Inc(value)
		stat.Inc(value)
	}

	stat.Inc(101)

	expected := []Item[int]{
		{
			Quantity: 2,
			Span: span.Span[int]{
				Begin: math.MinInt,
				End:   0,
			},
		},
		{
			Quantity: 39,
			Span: span.Span[int]{
				Begin: 1,
				End:   39,
			},
		},
		{
			Quantity: 67,
			Span: span.Span[int]{
				Begin: 40,
				End:   78,
			},
		},
		{
			Quantity: 44,
			Span: span.Span[int]{
				Begin: 79,
				End:   100,
			},
		},
		{
			Quantity: 1,
			Span: span.Span[int]{
				Begin: 101,
				End:   math.MaxInt,
			},
		},
	}

	require.Equal(t, expected, stat.Items())
	require.NoError(t, stat.Graph(io.Discard))
}

func TestStatLinearFullRange(t *testing.T) {
	stat, err := NewLinear[uint8](0, math.MaxUint8, 100)
	require.NoError(t, err)

	for value := range safe.Inc[uint8](0, math.MaxUint8) {
		stat.Inc(value)
	}

	expected := []Item[uint8]{
		{
			Quantity: 100,
			Span: span.Span[uint8]{
				Begin: 0,
				End:   99,
			},
		},
		{
			Quantity: 100,
			Span: span.Span[uint8]{
				Begin: 100,
				End:   199,
			},
		},
		{
			Quantity: 56,
			Span: span.Span[uint8]{
				Begin: 200,
				End:   255,
			},
		},
	}

	require.Equal(t, expected, stat.Items())
	require.NoError(t, stat.Graph(io.Discard))
}

func TestStatSearch(t *testing.T) {
	spans := []span.Span[int]{
		{Begin: 1, End: 2},
		{Begin: 3, End: 4},
		{Begin: 6, End: 8},
	}

	stat, err := New(spans, nil)
	require.NoError(t, err)

	stat.Inc(-1)
	stat.Inc(0)

	stat.Inc(1)
	stat.Inc(2)
	stat.Inc(2)

	stat.Inc(3)
	stat.Inc(4)
	stat.Inc(3)
	stat.Inc(4)
	stat.Inc(5)

	stat.Inc(6)
	stat.Inc(7)
	stat.Inc(8)
	stat.Inc(6)
	stat.Inc(7)

	stat.Inc(9)
	stat.Inc(10)
	stat.Inc(11)

	expected := []Item[int]{
		{
			Quantity: 1,
			Span:     span.Span[int]{},
		},
		{
			Quantity: 2,
			Span: span.Span[int]{
				Begin: math.MinInt,
				End:   0,
			},
		},
		{
			Quantity: 3,
			Span: span.Span[int]{
				Begin: 1,
				End:   2,
			},
		},
		{
			Quantity: 4,
			Span: span.Span[int]{
				Begin: 3,
				End:   4,
			},
		},
		{
			Quantity: 5,
			Span: span.Span[int]{
				Begin: 6,
				End:   8,
			},
		},
		{
			Quantity: 3,
			Span: span.Span[int]{
				Begin: 9,
				End:   math.MaxInt,
			},
		},
	}

	require.Equal(t, expected, stat.Items())
	require.NoError(t, stat.Graph(io.Discard))
}

func TestStatError(t *testing.T) {
	stat, err := New[int](nil, nil)
	require.Error(t, err)
	require.Nil(t, stat)
}

func TestStatGraphError(t *testing.T) {
	stat, err := New([]span.Span[int]{{Begin: 0, End: 0}}, nil)
	require.NoError(t, err)

	stat.missed.Quantity = math.MaxUint64
	stat.negInf.Quantity = 0
	stat.items[0].Quantity = 0
	stat.posInf.Quantity = 0
	require.Error(t, stat.Graph(io.Discard))

	stat.missed.Quantity = 0
	stat.negInf.Quantity = math.MaxUint64
	stat.items[0].Quantity = 0
	stat.posInf.Quantity = 0
	require.Error(t, stat.Graph(io.Discard))

	stat.negInf.Quantity = 0
	stat.items[0].Quantity = math.MaxUint64
	stat.posInf.Quantity = 0
	require.Error(t, stat.Graph(io.Discard))

	stat.negInf.Quantity = 0
	stat.items[0].Quantity = 0
	stat.posInf.Quantity = math.MaxUint64
	require.Error(t, stat.Graph(io.Discard))

	stdout := os.Stdout
	os.Stdout = nil

	require.Error(t, stat.Graph())

	os.Stdout = stdout
}

func BenchmarkStatLinear(b *testing.B) {
	stat, err := NewLinear(1, 80, 10)
	require.NoError(b, err)

	for range b.N {
		stat.Inc(0)
		stat.Inc(1)
		stat.Inc(11)
		stat.Inc(21)
		stat.Inc(31)
		stat.Inc(41)
		stat.Inc(51)
		stat.Inc(61)
		stat.Inc(71)
		stat.Inc(81)
	}
}

func BenchmarkStatSearch(b *testing.B) {
	spans, err := span.Linear(1, 80, 10)
	require.NoError(b, err)

	stat, err := New(spans, nil)
	require.NoError(b, err)

	for range b.N {
		stat.Inc(0)
		stat.Inc(1)
		stat.Inc(11)
		stat.Inc(21)
		stat.Inc(31)
		stat.Inc(41)
		stat.Inc(51)
		stat.Inc(61)
		stat.Inc(71)
		stat.Inc(81)
	}
}
