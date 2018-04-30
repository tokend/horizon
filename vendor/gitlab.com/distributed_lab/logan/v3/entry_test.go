package logan

import (
	"github.com/sirupsen/logrus"
	"testing"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func TestEntry(t *testing.T) {
	err := errors.From(errors.New("Error."), F{
		"key": "value",
	})

	cases := []struct {
		actual   *Entry
		expected *Entry
	}{
		// 1
		{New().WithField("key", "value"), &Entry{
			entry: &logrus.Entry{
				Data: map[string]interface{}{
					"key": "value",
				},
			},
		}},
		// 2
		{New().WithField("key", "value").WithField("key2", "value2"),
			New().WithFields(F{
				"key":  "value",
				"key2": "value2",
			})},
		{New().WithField("my_err", err),
			New().WithFields(F{
				"my_err":  err,
				"key": "value",
			})},
	}

	for _, tc := range cases {
		for expectedK, expectedV := range tc.expected.entry.Data {
			value, ok := tc.actual.entry.Data[expectedK]
			if !ok {
				t.Errorf("Missing key: %s", expectedK)
				continue
			}

			if value != expectedV {
				t.Errorf("Wrong value for %s: expected (%s) got (%s)", expectedK, expectedV, value)
			}
		}
	}

}
