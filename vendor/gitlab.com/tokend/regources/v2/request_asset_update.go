package regources

// AssetUpdateRequest - represents details of the `asset update` reviewable request
type AssetUpdateRequest struct {
	Key
	Attributes    AssetUpdateRequestAttrs     `json:"attributes"`
	Relationships AssetUpdateRequestRelations `json:"relationships"`
}

// AssetUpdateRequestAttrs - attributes of the `asset update` reviewable request
type AssetUpdateRequestAttrs struct {
	Policies int32   `json:"policies"`
	Details  Details `json:"details"`
}

// AssetUpdateRequestRelations - attributes of the `asset update` reviewable request
type AssetUpdateRequestRelations struct {
	Asset *Relation `json:"asset"`
}
