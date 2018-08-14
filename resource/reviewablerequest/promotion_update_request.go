package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/regources/reviewablerequest2"
)

func PopulatePromotionUpdateRequest(histRequest history.PromotionUpdateRequest) (
	r *reviewablerequest2.PromotionUpdateRequest, err error,
) {
	r = &reviewablerequest2.PromotionUpdateRequest{}
	r.SaleID = histRequest.SaleID
	newPromotionData, err := PopulateSaleCreationRequest(histRequest.NewPromotionData)
	if newPromotionData != nil {
		r.NewPromotionData = *newPromotionData
	}
	return
}
