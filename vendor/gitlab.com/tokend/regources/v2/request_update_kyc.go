package regources

import (
	"encoding/json"
	"gitlab.com/tokend/go/xdr"
)

// UpdateKYCRequest - represents details of the `update_kyc` reviewable request
type UpdateKYCRequest struct {
	Key
	Attributes    UpdateKYCRequestAttrs     `json:"attributes"`
	Relationships UpdateKYCRequestRelations `json:"relationships"`
}

// UpdateKYCRequestAttrs - attributes of the `update_kyc` reviewable request
type UpdateKYCRequestAttrs struct {
	AccountTypeToSet xdr.AccountType `json:"account_type_to_set"`
	KYCLevel         uint32          `json:"kyc_level"`
	KYCData          Details         `json:"kyc_data"`
	SequenceNumber   uint32          `json:"sequence_number"`
	ExternalDetails  []Details       `json:"external_details"`
}

func (r UpdateKYCRequestAttrs) MarshalJSON () ([]byte, error) {
	if r.ExternalDetails == nil {
		r.ExternalDetails = []Details{}
	}

	type temp UpdateKYCRequestAttrs
	return json.Marshal(temp(r))
}

// UpdateKYCRequestRelations - attributes of the `update_kyc` reviewable request
type UpdateKYCRequestRelations struct {
	AccountToUpdateKYC *Relation `json:"account_to_update_kyc"`
}
