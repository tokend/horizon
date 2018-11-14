package history2

import "gitlab.com/tokend/go/xdr"


// ParticipantEffect - stores effect of operation done to particular account or balance
type ParticipantEffect struct {
	ID          int64          `db:"id"`
	AccountID   xdr.AccountId  `db:"account_id"`
	BalanceID   *xdr.BalanceId `db:"balance_id"`
	AssetCode   *xdr.AssetCode `db:"asset_code"`
	Effect      Effect         `db:"effect"`
	OperationID int64          `db:"operation_id"`
}
