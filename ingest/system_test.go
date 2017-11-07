package ingest

import (
	"testing"

	"bullioncoin.githost.io/development/go/network"
	"gitlab.com/distributed_lab/tokend/horizon/test"
)

func TestValidation(t *testing.T) {
	tt := test.Start(t).ScenarioWithoutHorizon("kahuna")
	defer tt.Finish()

	sys := New(network.TestNetworkPassphrase, "GDVEPIE37LURE2G4CJM5OJ6RAV2BKR4OKXCOYXJEXOBRPC2KACRLYTOT",
		"GAK744S3MBMXOMD7UMJ5SDXTD57QH6K4FRKPIC5BF5ESGSKHD4OYUBTA", "", tt.CoreRepo(), tt.HorizonRepo())

	// intact chain
	for i := int32(2); i <= 59; i++ {
		tt.Assert.NoError(sys.validateLedgerChain(i))
	}
	_, err := tt.CoreRepo().ExecRaw(
		`DELETE FROM ledgerheaders WHERE ledgerseq = ?`, 5,
	)
	tt.Require.NoError(err)

	// missing cur
	err = sys.validateLedgerChain(5)
	tt.Assert.Error(err)
	tt.Assert.Contains(err.Error(), "failed to load cur ledger")

	// missing prev
	err = sys.validateLedgerChain(6)
	tt.Assert.Error(err)
	tt.Assert.Contains(err.Error(), "failed to load prev ledger")

	// mismatched header
	_, err = tt.CoreRepo().ExecRaw(`
		UPDATE ledgerheaders
		SET ledgerhash = ?
		WHERE ledgerseq = ?`, "00000", 8)
	tt.Require.NoError(err)

	err = sys.validateLedgerChain(9)
	tt.Assert.Error(err)
	tt.Assert.Contains(err.Error(), "cur and prev ledger hashes don't match")
}
