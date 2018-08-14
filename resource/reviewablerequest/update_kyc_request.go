package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/regources/reviewablerequest2"
)

func PopulateUpdateKYCRequest(histRequest history.UpdateKYCRequest) (
	r *reviewablerequest2.UpdateKYCRequest, err error,
) {
	r = &reviewablerequest2.UpdateKYCRequest{}
	r.AccountToUpdateKYC = histRequest.AccountToUpdateKYC
	r.AccountTypeToSet = int32(histRequest.AccountTypeToSet)
	r.KYCLevel = histRequest.KYCLevel
	r.KYCData = histRequest.KYCData
	r.AllTasks = histRequest.AllTasks
	r.PendingTasks = histRequest.PendingTasks
	r.SequenceNumber = histRequest.SequenceNumber
	r.ExternalDetails = histRequest.ExternalDetails
	return
}
