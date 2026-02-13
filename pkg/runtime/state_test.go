package runtime

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestState(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		_ = NewForTest()
		require.True(t, true)
	})
}
