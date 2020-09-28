package history2

import regources "gitlab.com/tokend/regources/generated"

type DeferredPayment struct {
	ID                 int64            `db:"id"`
	Amount             regources.Amount `db:"amount"`
	SourceAccount      string           `db:"source_account"`
	SourceBalance      string           `db:"source_balance"`
	DestinationAccount string           `db:"destination_account"`
}
