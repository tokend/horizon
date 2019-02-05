package regources

import (
	"encoding/json"
)

// ChangeRoleRequest - represents details of the `update_kyc` reviewable request
type ChangeRoleRequest struct {
	Key
	Attributes    ChangeRoleRequestAttrs     `json:"attributes"`
	Relationships ChangeRoleRequestRelations `json:"relationships"`
}

// ChangeRoleRequestAttrs - attributes of the `update_kyc` reviewable request
type ChangeRoleRequestAttrs struct {
	AccountRoleToSet uint64    `json:"account_role_to_set"`
	KYCData          Details   `json:"kyc_data"`
	SequenceNumber   uint32    `json:"sequence_number"`
	CreatorDetails   []Details `json:"creator_details"`
}

func (r ChangeRoleRequestAttrs) MarshalJSON() ([]byte, error) {
	if r.CreatorDetails == nil {
		r.CreatorDetails = []Details{}
	}

	type temp ChangeRoleRequestAttrs
	return json.Marshal(temp(r))
}

// ChangeRoleRequestRelations - attributes of the `update_kyc` reviewable request
type ChangeRoleRequestRelations struct {
	AccountToUpdateRole *Relation `json:"account_to_update_role"`
}
