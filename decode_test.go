package horizon

import (
	"testing"

	"github.com/davecgh/go-spew/spew"

	"gitlab.com/tokend/go/xdr"
)

func TestDecode(t *testing.T) {
	x := `AAAAAJF4bc3UwVvT5+TXWxbuLcccAY1LPU5QZPqOn5QJADbJAAAAAGpoZOYAAAAAAAAAAAAAAABceRByAAAAAAAAAAIAAAAAAAAAJgAAAAAAAAAATSWlHx+KnA3tHU9njVogaw38f7VONvqP3hZR9wJJ1YgAAAAAAAAAAgAAAAAAAAPoAAAAAAAAAAEAAAACe30AAAAAAAAAAAAAAAAAAAAAACYAAAACAAAAAJF4bc3UwVvT5+TXWxbuLcccAY1LPU5QZPqOn5QJADbJAAAAAAAAAAAAAAAAAAAAAQkANskAAABAICevv+KJcL9MbVuucdYPMiwxR/wgnC6Sri/RPAUCMyJLvmOcjx3klxFeyPcHatd11+R+HO4h5GjJjEuwcKS4Cg==`
	var d xdr.TransactionEnvelope
	err := xdr.SafeUnmarshalBase64(x, &d)
	if err != nil {
		t.Fatal(err)
	}
	spew.Dump(d)
}
