package urlval

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrBadRequest_Error(t *testing.T) {
	t.Run("returns properly formatted error string", func(t *testing.T) {
		errs := errBadRequest{
			"field_1": errors.New("field is invalid for some reason"),
		}

		expected := "field_1: field is invalid for some reason"

		assert.Equal(t, expected, errs.Error())
	})
}

func TestErrBadRequest_Filter(t *testing.T) {
	t.Run("returns nil when empty", func(t *testing.T) {
		errs := errBadRequest{}
		assert.Nil(t, errs.Filter())
	})
	t.Run("returns nil when filled with empty errors", func(t *testing.T) {
		errs := errBadRequest{
			"field1": nil,
			"field2": nil,
			"field3": nil,
		}
		assert.Nil(t, errs.Filter())
	})
	t.Run("returns self when all filled with errors", func(t *testing.T) {
		errs := errBadRequest{
			"field": errors.New("field is invalid"),
		}
		assert.Equal(t, errs, errs.Filter())
	})
	t.Run("filters out nil errors when have it inside", func(t *testing.T) {
		errs := errBadRequest{
			"field":     errors.New("field is invalid"),
			"nil-field": nil,
		}

		filteredErrs := errs.Filter()
		expectedResult := errBadRequest{
			"field": errors.New("field is invalid"),
		}

		assert.Equal(t, filteredErrs, expectedResult)
		assert.Equal(t, errs, expectedResult)
	})
}
