package reviewablerequest

import (
	"encoding/json"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/db2/history"
)

type UpdateKYCRequest struct {
	AccountToUpdateKYC string                   `json:"account_to_update_kyc"`
	AccountTypeToSet   xdr.AccountType          `json:"account_type_to_set"`
	KYCLevel           uint32                   `json:"kyc_level"`
	KYCData            map[string]interface{}   `json:"kyc_data"`
	AllTasks           uint32                   `json:"all_tasks"`
	PendingTasks       uint32                   `json:"pending_tasks"`
	SequenceNumber     uint32                   `json:"sequence_number"`
	ExternalDetails    []map[string]interface{} `json:"extrenal_details"`
}

func (r *UpdateKYCRequest) Populate(histRequest history.UpdateKYCRequest) {
	r.AccountToUpdateKYC = histRequest.AccountToUpdateKYC
	r.AccountTypeToSet = histRequest.AccountTypeToSet
	r.KYCLevel = histRequest.KYCLevel
	r.KYCData = histRequest.KYCData
	r.AllTasks = histRequest.AllTasks
	r.PendingTasks = histRequest.PendingTasks
	r.SequenceNumber = histRequest.SequenceNumber
	r.ExternalDetails = histRequest.ExternalDetails
}

func (r *UpdateKYCRequest) PopulateFromRawJsonHistory(rawJson []byte) error {
	var histRequest history.UpdateKYCRequest
	err := json.Unmarshal(rawJson, &histRequest)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal history.UpdateKYCRequest")
	}

	r.Populate(histRequest)
	return nil
}
