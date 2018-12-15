package history2

// ParticipantEffect - stores effect of operation done to particular account or balance
type ParticipantEffect struct {
	ID          int64   `db:"id"`
	AccountID   int64   `db:"account_id"`
	BalanceID   *int64  `db:"balance_id"`
	AssetCode   *string `db:"asset_code"`
	Effect      Effect  `db:"effect"`
	OperationID int64   `db:"operation_id"`
}
