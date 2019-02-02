package regources

// AmlAlertRequest - represents details of the `aml alert` reviewable request
type AmlAlertRequest struct {
	Key
	Attributes    AmlAlertRequestAttrs     `json:"attributes"`
	Relationships AmlAlertRequestRelations `json:"relationships"`
}

// AmlAlertRequestAttrs - attributes of the `aml alert` reviewable request
type AmlAlertRequestAttrs struct {
	Amount string `json:"amount"`
	Reason string `json:"reason"`
}

// AmlAlertRequestRelations - relationships of the `aml alert` reviewable request
type AmlAlertRequestRelations struct {
	Balance *Relation `json:"balance"`
}
