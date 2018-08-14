package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/regources/reviewablerequest2"
)

func PopulateUpdateSaleEndTimeRequest(histRequest history.UpdateSaleEndTimeRequest) (
	r *reviewablerequest2.UpdateSaleEndTimeRequest, err error,
) {
	r = &reviewablerequest2.UpdateSaleEndTimeRequest{}
	r.SaleID = histRequest.SaleID
	r.NewEndTime = histRequest.NewEndTime
	return
}
