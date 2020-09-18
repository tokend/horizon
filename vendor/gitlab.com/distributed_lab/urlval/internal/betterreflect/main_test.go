package betterreflect

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvertSlice(t *testing.T) {
	t.Run("converts string to alias string", func(t *testing.T) {
		type String string

		var source = []string{"1", "2", "foo", "bar"}

		result, err := parseSlice(reflect.TypeOf(String("")), source)

		require.NoError(t, err)
		require.IsType(t, []String{}, result)
		require.EqualValues(t, []String{"1", "2", "foo", "bar"}, result)
	})
}
