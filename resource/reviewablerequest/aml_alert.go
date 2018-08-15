package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/regources"
)

func PopulateAmlAlertRequest(histRequest history.AmlAlertRequest) (
	*regources.AMLAlertRequest, error,
) {
	return &regources.AMLAlertRequest{
		BalanceID: histRequest.BalanceID,
		Amount:    histRequest.Amount,
		Reason:    histRequest.Reason,
	}, nil
}
