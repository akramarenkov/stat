package stat

import (
	"testing"

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
