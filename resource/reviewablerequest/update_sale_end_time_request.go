package reviewablerequest

import (
	"gitlab.com/tokend/horizon/db2/history"
	"gitlab.com/tokend/regources"
)

func PopulateUpdateSaleEndTimeRequest(histRequest history.UpdateSaleEndTimeRequest) (
	r *regources.UpdateSaleEndTimeRequest, err error,
) {
	r = &regources.UpdateSaleEndTimeRequest{}
	r.SaleID = histRequest.SaleID
	r.NewEndTime = regources.Time(histRequest.NewEndTime)
	return
}
