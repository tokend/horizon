package horizon

import (
	"encoding/json"
	"testing"

	"gitlab.com/distributed_lab/tokend/horizon/resource/operations"
	"gitlab.com/distributed_lab/tokend/horizon/test"
)

func TestOperationActions_Index(t *testing.T) {
	ht := StartHTTPTest(t, "base")
	defer ht.Finish()

	// no filter
	w := ht.Get("/operations")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(4, w.Body)
	}

	// filtered by ledger sequence
	w = ht.Get("/ledgers/1/operations")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(0, w.Body)
	}

	w = ht.Get("/ledgers/2/operations")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(3, w.Body)
	}

	w = ht.Get("/ledgers/3/operations")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(1, w.Body)
	}

	// filtered by account
	w = ht.Get("/accounts/GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H/operations")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(3, w.Body)
	}

	w = ht.Get("/accounts/GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2/operations")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(1, w.Body)
	}

	w = ht.Get("/accounts/GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU/operations")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(2, w.Body)
	}

	// filtered by transaction
	w = ht.Get("/transactions/2374e99349b9ef7dba9a5db3339b78fda8f34777b1af33ba468ad5c0df946d4d/operations")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(1, w.Body)
	}

	w = ht.Get("/transactions/164a5064eba64f2cdbadb856bf3448485fc626247ada3ed39cddf0f6902133b6/operations")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(1, w.Body)
	}

	// filtered by ledger
	w = ht.Get("/ledgers/3/operations")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(1, w.Body)
	}

	// missing ledger
	w = ht.Get("/ledgers/100/operations")
	ht.Assert.Equal(404, w.Code)
}

func TestOperationActions_Show(t *testing.T) {
	ht := StartHTTPTest(t, "base")
	defer ht.Finish()

	// exists
	w := ht.Get("/operations/8589938689")
	if ht.Assert.Equal(200, w.Code) {
		var result operations.Base
		err := json.Unmarshal(w.Body.Bytes(), &result)
		ht.Require.NoError(err, "failed to parse body")
		ht.Assert.Equal("8589938689", result.PT)
	}

	// doesn't exist
	w = ht.Get("/operations/9589938689")
	ht.Assert.Equal(404, w.Code)

	// before history
	ht.ReapHistory(1)
	w = ht.Get("/operations/8589938689")
	ht.Assert.Equal(410, w.Code)
}

func TestOperationActions_Regressions(t *testing.T) {
	ht := StartHTTPTest(t, "base")
	defer ht.Finish()

	// ensure that trying to stream ops from an account that doesn't exist
	// fails before streaming the hello message.  Regression test for #285
	w := ht.Get("/accounts/foo/operations?limit=1", test.RequestHelperStreaming)
	if ht.Assert.Equal(404, w.Code) {
		ht.Assert.ProblemType(w.Body, "not_found")
	}
}
