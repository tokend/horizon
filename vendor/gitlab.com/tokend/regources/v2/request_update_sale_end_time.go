package regources

import "time"

// UpdateSaleEndTimeRequest - represents details of the `update_sale_end_time` reviewable request
type UpdateSaleEndTimeRequest struct {
	Key
	Attributes    UpdateSaleEndTimeRequestAttrs     `json:"attributes"`
	Relationships UpdateSaleEndTimeRequestRelations `json:"relationships"`
}

// UpdateSaleEndTimeRequestAttrs - attributes of the `update_sale_end_time` reviewable request
type UpdateSaleEndTimeRequestAttrs struct {
	NewEndTime time.Time `json:"new_end_time"`
}

// UpdateSaleEndTimeRequestRelations - relationships of the `update_sale_end_time` reviewable request
type UpdateSaleEndTimeRequestRelations struct {
	Sale *Relation `json:"sale"`
}
