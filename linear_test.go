package stat

import (
	"io"
	"math"
	"testing"

	"github.com/akramarenkov/safe"
	"github.com/akramarenkov/span"
	"github.com/stretchr/testify/require"
)

func TestLinear(t *testing.T) {
	linear, err := NewLinear(1, 100, 39)
	require.NoError(t, err)

	linear.Add(-102)
	linear.Add(-101)

	for value := range safe.Inc(1, 50) {
		linear.Add(value)
	}

	for value := range safe.Inc(51, 100) {
		linear.Add(value)
		linear.Add(value)
	}

	linear.Add(101)

	expected := []Item[int]{
		{
			Quantity: 2,
			Span: span.Span[int]{
				Begin: math.MinInt64,
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
				End:   math.MaxInt64,
			},
		},
	}

	require.Equal(t, expected, linear.Stat().Items())
	require.NoError(t, linear.Stat().Graph(io.Discard))
}

func TestLinearFullRange(t *testing.T) {
	linear, err := NewLinear[uint8](0, math.MaxUint8, 100)
	require.NoError(t, err)

	for value := range safe.Inc[uint8](0, math.MaxUint8) {
		linear.Add(value)
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

	require.Equal(t, expected, linear.Stat().Items())
	require.NoError(t, linear.Stat().Graph(io.Discard))
}

func TestLinearError(t *testing.T) {
	linear, err := NewLinear(1, 2, -1)
	require.Error(t, err)
	require.Nil(t, linear)

	linear, err = NewLinear(2, 1, 1)
	require.Error(t, err)
	require.Nil(t, linear)
}

func BenchmarkLinear(b *testing.B) {
	linear, err := NewLinear(1, 80, 10)
	require.NoError(b, err)

	for range b.N {
		linear.Add(0)
		linear.Add(1)
		linear.Add(11)
		linear.Add(21)
		linear.Add(31)
		linear.Add(41)
		linear.Add(51)
		linear.Add(61)
		linear.Add(71)
		linear.Add(81)
	}
}
