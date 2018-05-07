package reviewablerequest

import (
	"gitlab.com/tokend/go/xdr"
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
	ExternalDetails    []map[string]interface{} `json:"external_details"`
}

func (r *UpdateKYCRequest) Populate(histRequest history.UpdateKYCRequest) error {
	r.AccountToUpdateKYC = histRequest.AccountToUpdateKYC
	r.AccountTypeToSet = histRequest.AccountTypeToSet
	r.KYCLevel = histRequest.KYCLevel
	r.KYCData = histRequest.KYCData
	r.AllTasks = histRequest.AllTasks
	r.PendingTasks = histRequest.PendingTasks
	r.SequenceNumber = histRequest.SequenceNumber
	r.ExternalDetails = histRequest.ExternalDetails
	return nil
}