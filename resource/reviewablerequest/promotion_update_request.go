package reviewablerequest

import "gitlab.com/swarmfund/horizon/db2/history"

type PromotionUpdateRequest struct {
	SaleID           uint64              `json:"sale_id"`
	NewPromotionData SaleCreationRequest `json:"new_promotion_data"`
}

func (r *PromotionUpdateRequest) Populate(histRequest history.PromotionUpdateRequest) error {
	r.SaleID = histRequest.SaleID
	r.NewPromotionData.Populate(histRequest.NewPromotionData)
	return nil
}
