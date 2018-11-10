package reviewablerequest

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/horizon/db2/history"
	"gitlab.com/tokend/regources"
)

func PopulateWithdrawalRequest(histRequest history.WithdrawalRequest) (
	r *regources.WithdrawalRequest, err error,
) {
	r = &regources.WithdrawalRequest{}
	r.BalanceID = histRequest.BalanceID
	r.Amount = regources.Amount(amount.MustParse(histRequest.Amount))
	r.FixedFee = regources.Amount(amount.MustParse(histRequest.FixedFee))
	r.PercentFee = regources.Amount(amount.MustParse(histRequest.PercentFee))
	r.ExternalDetails = histRequest.ExternalDetails
	r.DestAssetCode = histRequest.DestAssetCode
	r.DestAssetAmount = regources.Amount(amount.MustParse(histRequest.DestAssetAmount))
	r.ReviewerDetails = histRequest.ReviewerDetails
	r.PreConfirmationDetails = histRequest.PreConfirmationDetails
	return
}
