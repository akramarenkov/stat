package stat

import (
	"math"
	"testing"

	"github.com/akramarenkov/span"
	"github.com/stretchr/testify/require"
)

func TestLinearError(t *testing.T) {
	stat, err := NewLinear(1, 2, -1)
	require.Error(t, err)
	require.Nil(t, stat)

	stat, err = NewLinear(1, 2, 0)
	require.Error(t, err)
	require.Nil(t, stat)

	stat, err = NewLinear(2, 1, 1)
	require.Error(t, err)
	require.Nil(t, stat)
}

func TestLinearQ1(t *testing.T) {
	stat, err := NewLinearQ[int8](1, 100, 1)
	require.NoError(t, err)

	stat.Inc(math.MinInt8)
	stat.Inc(math.MinInt8 + 1)
	stat.Inc(-2)
	stat.Inc(-1)
	stat.Inc(0)

	stat.Inc(1)
	stat.Inc(2)
	stat.Inc(99)
	stat.Inc(100)

	stat.Inc(101)
	stat.Inc(102)
	stat.Inc(math.MaxInt8 - 1)
	stat.Inc(math.MaxInt8)

	expected := []Item[int8]{
		{
			Quantity: 5,
			Span: span.Span[int8]{
				Begin: math.MinInt8,
				End:   0,
			},
			Kind: ItemKindNegInf,
		},
		{
			Quantity: 4,
			Span: span.Span[int8]{
				Begin: 1,
				End:   100,
			},
			Kind: ItemKindRegular,
		},
		{
			Quantity: 4,
			Span: span.Span[int8]{
				Begin: 101,
				End:   math.MaxInt8,
			},
			Kind: ItemKindPosInf,
		},
	}

	require.Equal(t, expected, stat.Items())
}

func TestLinearQ1Max(t *testing.T) {
	stat, err := NewLinearQ[int8](math.MinInt8, math.MaxInt8, 1)
	require.NoError(t, err)

	stat.Inc(math.MinInt8)
	stat.Inc(math.MinInt8 + 1)

	stat.Inc(-2)
	stat.Inc(-1)
	stat.Inc(0)
	stat.Inc(1)
	stat.Inc(2)

	stat.Inc(math.MaxInt8 - 1)
	stat.Inc(math.MaxInt8)

	expected := []Item[int8]{
		{
			Quantity: 9,
			Span: span.Span[int8]{
				Begin: math.MinInt8,
				End:   math.MaxInt8,
			},
			Kind: ItemKindRegular,
		},
	}

	require.Equal(t, expected, stat.Items())
}

func TestLinearQ2(t *testing.T) {
	stat, err := NewLinearQ[int8](1, 100, 2)
	require.NoError(t, err)

	stat.Inc(math.MinInt8)
	stat.Inc(math.MinInt8 + 1)
	stat.Inc(-2)
	stat.Inc(-1)
	stat.Inc(0)

	stat.Inc(1)
	stat.Inc(2)
	stat.Inc(49)
	stat.Inc(50)

	stat.Inc(51)
	stat.Inc(52)
	stat.Inc(98)
	stat.Inc(99)
	stat.Inc(100)

	stat.Inc(101)
	stat.Inc(102)
	stat.Inc(math.MaxInt8 - 1)
	stat.Inc(math.MaxInt8)

	expected := []Item[int8]{
		{
			Quantity: 5,
			Span: span.Span[int8]{
				Begin: math.MinInt8,
				End:   0,
			},
			Kind: ItemKindNegInf,
		},
		{
			Quantity: 4,
			Span: span.Span[int8]{
				Begin: 1,
				End:   50,
			},
			Kind: ItemKindRegular,
		},
		{
			Quantity: 5,
			Span: span.Span[int8]{
				Begin: 51,
				End:   100,
			},
			Kind: ItemKindRegular,
		},
		{
			Quantity: 4,
			Span: span.Span[int8]{
				Begin: 101,
				End:   math.MaxInt8,
			},
			Kind: ItemKindPosInf,
		},
	}

	require.Equal(t, expected, stat.Items())
}

func TestLinearQ2Max(t *testing.T) {
	stat, err := NewLinearQ[int8](math.MinInt8, math.MaxInt8, 2)
	require.NoError(t, err)

	stat.Inc(math.MinInt8)
	stat.Inc(math.MinInt8 + 1)

	stat.Inc(-2)
	stat.Inc(-1)
	stat.Inc(0)
	stat.Inc(1)
	stat.Inc(2)

	stat.Inc(math.MaxInt8 - 1)
	stat.Inc(math.MaxInt8)

	expected := []Item[int8]{
		{
			Quantity: 4,
			Span: span.Span[int8]{
				Begin: math.MinInt8,
				End:   -1,
			},
			Kind: ItemKindRegular,
		},
		{
			Quantity: 5,
			Span: span.Span[int8]{
				Begin: 0,
				End:   math.MaxInt8,
			},
			Kind: ItemKindRegular,
		},
	}

	require.Equal(t, expected, stat.Items())
}

func TestLinearQ2MaxUns(t *testing.T) {
	stat, err := NewLinearQ[uint8](0, math.MaxUint8, 2)
	require.NoError(t, err)

	stat.Inc(0)
	stat.Inc(1)
	stat.Inc(2)
	stat.Inc(125)
	stat.Inc(126)
	stat.Inc(127)

	stat.Inc(128)
	stat.Inc(129)
	stat.Inc(130)

	stat.Inc(math.MaxUint8 - 1)
	stat.Inc(math.MaxUint8)

	expected := []Item[uint8]{
		{
			Quantity: 6,
			Span: span.Span[uint8]{
				Begin: 0,
				End:   127,
			},
			Kind: ItemKindRegular,
		},
		{
			Quantity: 5,
			Span: span.Span[uint8]{
				Begin: 128,
				End:   math.MaxUint8,
			},
			Kind: ItemKindRegular,
		},
	}

	require.Equal(t, expected, stat.Items())
}

func TestLinearQ10(t *testing.T) {
	stat, err := NewLinearQ(1, 100, 10)
	require.NoError(t, err)

	stat.Inc(-102)
	stat.Inc(-101)

	stat.Inc(1)
	stat.Inc(2)
	stat.Inc(11)
	stat.Inc(12)
	stat.Inc(21)
	stat.Inc(22)
	stat.Inc(31)
	stat.Inc(32)
	stat.Inc(41)
	stat.Inc(42)
	stat.Inc(51)
	stat.Inc(52)
	stat.Inc(61)
	stat.Inc(62)
	stat.Inc(71)
	stat.Inc(72)
	stat.Inc(81)
	stat.Inc(82)
	stat.Inc(91)
	stat.Inc(92)
	stat.Inc(99)
	stat.Inc(100)

	stat.Inc(101)

	expected := []Item[int]{
		{
			Quantity: 2,
			Span: span.Span[int]{
				Begin: math.MinInt,
				End:   0,
			},
			Kind: ItemKindNegInf,
		},
		{
			Quantity: 2,
			Span: span.Span[int]{
				Begin: 1,
				End:   10,
			},
			Kind: ItemKindRegular,
		},
		{
			Quantity: 2,
			Span: span.Span[int]{
				Begin: 11,
				End:   20,
			},
			Kind: ItemKindRegular,
		},
		{
			Quantity: 2,
			Span: span.Span[int]{
				Begin: 21,
				End:   30,
			},
			Kind: ItemKindRegular,
		},
		{
			Quantity: 2,
			Span: span.Span[int]{
				Begin: 31,
				End:   40,
			},
			Kind: ItemKindRegular,
		},
		{
			Quantity: 2,
			Span: span.Span[int]{
				Begin: 41,
				End:   50,
			},
			Kind: ItemKindRegular,
		},
		{
			Quantity: 2,
			Span: span.Span[int]{
				Begin: 51,
				End:   60,
			},
			Kind: ItemKindRegular,
		},
		{
			Quantity: 2,
			Span: span.Span[int]{
				Begin: 61,
				End:   70,
			},
			Kind: ItemKindRegular,
		},
		{
			Quantity: 2,
			Span: span.Span[int]{
				Begin: 71,
				End:   80,
			},
			Kind: ItemKindRegular,
		},
		{
			Quantity: 2,
			Span: span.Span[int]{
				Begin: 81,
				End:   90,
			},
			Kind: ItemKindRegular,
		},
		{
			Quantity: 4,
			Span: span.Span[int]{
				Begin: 91,
				End:   100,
			},
			Kind: ItemKindRegular,
		},
		{
			Quantity: 1,
			Span: span.Span[int]{
				Begin: 101,
				End:   math.MaxInt,
			},
			Kind: ItemKindPosInf,
		},
	}

	require.Equal(t, expected, stat.Items())
}

func TestLinearQError(t *testing.T) {
	stat, err := NewLinearQ(1, 2, -1)
	require.Error(t, err)
	require.Nil(t, stat)

	stat, err = NewLinearQ(1, 2, 0)
	require.Error(t, err)
	require.Nil(t, stat)

	stat, err = NewLinearQ(2, 1, 1)
	require.Error(t, err)
	require.Nil(t, stat)
}
