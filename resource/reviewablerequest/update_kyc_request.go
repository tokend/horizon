package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/regources"
)

func PopulateUpdateKYCRequest(histRequest history.UpdateKYCRequest) (
	r *regources.UpdateKYCRequest, err error,
) {
	r = &regources.UpdateKYCRequest{}
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
