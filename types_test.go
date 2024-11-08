package stat

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestItemKind(t *testing.T) {
	require.Equal(t, "regular", ItemKindRegular.String())
	require.Equal(t, "-Inf", ItemKindNegInf.String())
	require.Equal(t, "+Inf", ItemKindPosInf.String())
	require.Equal(t, "missed", ItemKindMissed.String())
	require.Equal(t, "unexpected", ItemKind(0).String())
}
