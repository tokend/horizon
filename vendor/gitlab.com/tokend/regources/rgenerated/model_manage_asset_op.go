package rgenerated

type ManageAssetOp struct {
	Key
	Attributes    *ManageAssetOpAttributes    `json:"attributes,omitempty"`
	Relationships *ManageAssetOpRelationships `json:"relationships,omitempty"`
}
type ManageAssetOpResponse struct {
	Data     ManageAssetOp `json:"data"`
	Included Included      `json:"included"`
}

type ManageAssetOpsResponse struct {
	Data     []ManageAssetOp `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
}

// MustManageAssetOp - returns ManageAssetOp from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustManageAssetOp(key Key) *ManageAssetOp {
	var manageAssetOp ManageAssetOp
	if c.tryFindEntry(key, &manageAssetOp) {
		return &manageAssetOp
	}
	return nil
}
