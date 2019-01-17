package core2

import "gitlab.com/tokend/horizon/db2"

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
	TrailingDigits         int64       `db:"trailing_digits"`
}
