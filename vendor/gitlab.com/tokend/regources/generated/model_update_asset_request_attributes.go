package regources

type UpdateAssetRequestAttributes struct {
	CreatorDetails Details `json:"creator_details"`
	// Policies specified for the asset creation
	Policies int32 `json:"policies"`
}
