package codes

import (
	"testing"

	"gitlab.com/swarmfund/go/xdr"
	. "github.com/smartystreets/goconvey/convey"
)


func TestCodeProviders(t *testing.T) {
	Convey("Code providers", t, func() {
		for _, opType := range xdr.OperationTypeAll {
			_, ok := codeProviders[opType]
			if !ok {
				t.Errorf("Failed to get code provider for %s", opType.ShortString())
			}
		}
	})
}
