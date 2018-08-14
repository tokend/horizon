package reviewablerequest2

type PromotionUpdateRequest struct {
	SaleID           uint64              `json:"sale_id"`
	NewPromotionData SaleCreationRequest `json:"new_promotion_data"`
}
