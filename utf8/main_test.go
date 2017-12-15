package utf8

import (
	"testing"

	"gitlab.com/swarmfund/horizon/test"
)

func TestScrub(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()

	tt.Assert.Equal("c78c167633163b96b03d1681911b1f9e1babe6d30d3e537c1c7f50af8e33c49", Scrub("c78c167633163b96b03d1681911b1f9e1babe6d30d3e537c1c7f50af8e33c49\u0000"))
	tt.Assert.Equal("scott", Scrub("scott"))
	tt.Assert.Equal("scött", Scrub("scött"))
	tt.Assert.Equal("�(", Scrub(string([]byte{0xC3, 0x28})))
}
