package core2

import (
	"math"

	"gitlab.com/tokend/horizon/db2"
)

const maximumTrailingDigits uint32 = 6

// Asset - db representation of asset
type Asset struct {
	Code                   string      `db:"code"`
	Owner                  string      `db:"owner"`
	PreIssuanceAssetSigner string      `db:"preissued_asset_signer"`
	Details                db2.Details `db:"details"`
	MaxIssuanceAmount      int64       `db:"max_issuance_amount"`
	AvailableForIssuance   int64       `db:"available_for_issueance"`
	Issued                 int64       `db:"issued"`
	PendingIssuance        int64       `db:"pending_issuance"`
	Policies               int32       `db:"policies"`
	TrailingDigits         uint32      `db:"trailing_digits"`
	Type                   uint64      `db:"type"`
}

//GetMinimumAmount - returns min amount support for that asset
func (a Asset) GetMinimumAmount() int64 {
	nullDigits := maximumTrailingDigits - a.TrailingDigits
	if nullDigits < 0 {
		panic("Unexpected database state. Expected asset trailing digits be equal or less 6")
	}

	return int64(math.Pow10(int(nullDigits)))
}
