package lorem

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashString(t *testing.T) {
	t.Run("deterministic", func(t *testing.T) {
		first := HashString("input")
		second := HashString("input")
		assert.Equal(t, first, second)
	})

	t.Run("random", func(t *testing.T) {
		first := HashString("foo")
		second := HashString("bar")
		assert.NotEqual(t, first, second)
	})

	t.Run("url safe", func(t *testing.T) {
		got := HashString("spamegg")
		escaped := url.QueryEscape(got)
		assert.Equal(t, got, escaped)
	})
}
