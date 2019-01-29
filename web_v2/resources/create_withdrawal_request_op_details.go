package resources

import (
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

func newCreateWithdrawalRequestOp(id int64,
	details history2.CreateWithdrawRequestDetails) *regources.CreateWithdrawRequest {
	return &regources.CreateWithdrawRequest{
		Key: regources.NewKeyInt64(id, regources.TypeCreateWithdrawalRequest),
		Attributes: regources.CreateWithdrawRequestAttrs{
			BalanceAddress:  details.BalanceAddress,
			Amount:          details.Amount,
			Fee:             details.Fee,
			ExternalDetails: details.ExternalDetails,
		},
	}
}
