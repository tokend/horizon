package history2

import (
	"gitlab.com/tokend/regources/generated"
)

type DeferredPayment struct {
	ID                    int64            `db:"id"`
	Amount                regources.Amount `db:"amount"`
	SourceAccount         string           `db:"source_account"`
	SourceBalance         string           `db:"source_balance"`
	DestinationAccount    string           `db:"destination_account"`
	SourcePaysForDest     bool             `db:"source_pays_for_dest"`
	SourceFixedFee        regources.Amount `db:"source_fixed_fee"`
	SourcePercentFee      regources.Amount `db:"source_percent_fee"`
	DestinationFixedFee   regources.Amount `db:"destination_fixed_fee"`
	DestinationPercentFee regources.Amount `db:"destination_percent_fee"`
}
