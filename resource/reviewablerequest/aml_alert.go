package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/regources/reviewablerequest2"
)

func PopulateAmlAlertRequest(histRequest history.AmlAlertRequest) (
	*reviewablerequest2.AmlAlertRequest, error,
) {
	return &reviewablerequest2.AmlAlertRequest{
		BalanceID: histRequest.BalanceID,
		Amount:    histRequest.Amount,
		Reason:    histRequest.Reason,
	}, nil
}
