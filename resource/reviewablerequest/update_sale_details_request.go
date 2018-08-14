package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/regources/reviewablerequest2"
)

func PopulateUpdateSaleDetailsRequest(histRequest history.UpdateSaleDetailsRequest) (
	r *reviewablerequest2.UpdateSaleDetailsRequest, err error,
) {
	r = &reviewablerequest2.UpdateSaleDetailsRequest{}
	r.SaleID = histRequest.SaleID
	r.NewDetails = histRequest.NewDetails
	return
}
