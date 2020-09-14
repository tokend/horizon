package urlval

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSort_Desc(t *testing.T) {
	t.Run("asc", func(t *testing.T) {
		require.False(t, Sort("price").Desc())
	})

	t.Run("desc", func(t *testing.T) {
		require.True(t, Sort("-price").Desc())
	})
}

func TestSort_Key(t *testing.T) {
	t.Run("asc", func(t *testing.T) {
		require.Equal(t, "price", Sort("price").Key())
	})

	t.Run("desc", func(t *testing.T) {
		require.Equal(t, "price", Sort("-price").Key())
	})
}
