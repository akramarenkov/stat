package stat

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLinearError(t *testing.T) {
	linear, err := NewLinear(1, 2, -1)
	require.Error(t, err)
	require.Nil(t, linear)

	linear, err = NewLinear(1, 2, 0)
	require.Error(t, err)
	require.Nil(t, linear)

	linear, err = NewLinear(2, 1, 1)
	require.Error(t, err)
	require.Nil(t, linear)
}
