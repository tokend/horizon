package history2

import "gitlab.com/tokend/regources/v2"

// ParticipantEffect - stores effect of operation done to particular account or balance
type ParticipantEffect struct {
	ID          int64            `db:"id"`
	AccountID   uint64           `db:"account_id"`
	BalanceID   *uint64          `db:"balance_id"`
	AssetCode   *string          `db:"asset_code"`
	Effect      regources.Effect `db:"effect"`
	OperationID int64            `db:"operation_id"`
}
