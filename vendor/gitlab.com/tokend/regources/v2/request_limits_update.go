package regources

// LimitsUpdateRequest - represents details of the `limits update` reviewable request
type LimitsUpdateRequest struct {
	Key
	Attributes LimitsUpdateRequestAttrs `json:"attributes"`
}

// LimitsUpdateRequestAttrs - attributes of the `limits update` reviewable request
type LimitsUpdateRequestAttrs struct {
	DocumentHash string  `json:"document_hash"`
	Details      Details `json:"details"`
}
