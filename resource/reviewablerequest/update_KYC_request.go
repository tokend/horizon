package reviewablerequest

import (
	"encoding/json"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/db2/history"
)

type ChangeKYCRequest struct {
	UpdatedAccountID string          `json:"updated_account_id"`
	AccountTypeToSet xdr.AccountType `json:"account_type_to_set"`
	KYCData          string          `json:"KYC_data"`
	KYCLevel         xdr.Uint32      `json:"KYC_level"`
}

func (r *ChangeKYCRequest) Populate(histRequest history.ChangeKYCRequest) {
	r.UpdatedAccountID = histRequest.UpdatedAccountId
	r.AccountTypeToSet = histRequest.AccountTypeToSet
	r.KYCData = histRequest.KYCData
	r.KYCLevel = histRequest.KYCLevel

}

func (r *ChangeKYCRequest) PopulateFromRawJsonHistory(rawJson []byte) error {
	var histRequest history.ChangeKYCRequest
	err := json.Unmarshal(rawJson, &histRequest)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal history.ChangeKYCRequest")
	}

	r.Populate(histRequest)
	return nil
}
