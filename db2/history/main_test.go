package history

import (
	"gitlab.com/distributed_lab/tokend/horizon/test"
	"testing"
)

func TestLatestLedger(t *testing.T) {
	tt := test.Start(t).Scenario("base")
	defer tt.Finish()
	q := &Q{tt.HorizonRepo()}

	var seq int
	err := q.LatestLedger(&seq)

	if tt.Assert.NoError(err) {
		tt.Assert.Equal(3, seq)
	}
}

func TestElderLedger(t *testing.T) {
	tt := test.Start(t).Scenario("base")
	defer tt.Finish()
	q := &Q{tt.HorizonRepo()}

	var seq int
	err := q.ElderLedger(&seq)

	if tt.Assert.NoError(err) {
		tt.Assert.Equal(1, seq)
	}
}
