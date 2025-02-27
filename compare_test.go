package stat

import (
	"slices"
	"testing"

	"github.com/akramarenkov/span"
	"github.com/stretchr/testify/require"
)

func TestSearch(t *testing.T) {
	items := []Item[int]{
		{Span: span.Span[int]{Begin: 1, End: 2}},
		{Span: span.Span[int]{Begin: 3, End: 4}},
		{Span: span.Span[int]{Begin: 6, End: 8}},
	}

	id, found := slices.BinarySearchFunc(
		items,
		Item[int]{Span: span.Span[int]{Begin: 0, End: 0}},
		search,
	)
	require.False(t, found)
	require.Equal(t, 0, id)

	id, found = slices.BinarySearchFunc(
		items,
		Item[int]{Span: span.Span[int]{Begin: 1, End: 1}},
		search,
	)
	require.True(t, found)
	require.Equal(t, 0, id)

	id, found = slices.BinarySearchFunc(
		items,
		Item[int]{Span: span.Span[int]{Begin: 2, End: 2}},
		search,
	)
	require.True(t, found)
	require.Equal(t, 0, id)

	id, found = slices.BinarySearchFunc(
		items,
		Item[int]{Span: span.Span[int]{Begin: 3, End: 3}},
		search,
	)
	require.True(t, found)
	require.Equal(t, 1, id)

	id, found = slices.BinarySearchFunc(
		items,
		Item[int]{Span: span.Span[int]{Begin: 4, End: 4}},
		search,
	)
	require.True(t, found)
	require.Equal(t, 1, id)

	id, found = slices.BinarySearchFunc(
		items,
		Item[int]{Span: span.Span[int]{Begin: 5, End: 5}},
		search,
	)
	require.False(t, found)
	require.Equal(t, 2, id)

	id, found = slices.BinarySearchFunc(
		items,
		Item[int]{Span: span.Span[int]{Begin: 6, End: 6}},
		search,
	)
	require.True(t, found)
	require.Equal(t, 2, id)

	id, found = slices.BinarySearchFunc(
		items,
		Item[int]{Span: span.Span[int]{Begin: 7, End: 7}},
		search,
	)
	require.True(t, found)
	require.Equal(t, 2, id)

	id, found = slices.BinarySearchFunc(
		items,
		Item[int]{Span: span.Span[int]{Begin: 8, End: 8}},
		search,
	)
	require.True(t, found)
	require.Equal(t, 2, id)

	id, found = slices.BinarySearchFunc(
		items,
		Item[int]{Span: span.Span[int]{Begin: 9, End: 9}},
		search,
	)
	require.False(t, found)
	require.Equal(t, 3, id)
}
