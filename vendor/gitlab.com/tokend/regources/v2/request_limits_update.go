package regources

// UpdateLimitsRequest - represents details of the `limits update` reviewable request
type UpdateLimitsRequest struct {
	Key
	Attributes UpdateLimitsRequestAttrs `json:"attributes"`
}

// UpdateLimitsRequestAttrs - attributes of the `limits update` reviewable request
type UpdateLimitsRequestAttrs struct {
	DocumentHash   string  `json:"document_hash"`
	CreatorDetails Details `json:"creator_details"`
}
