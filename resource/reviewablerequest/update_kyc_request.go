package reviewablerequest

import (
	"gitlab.com/tokend/horizon/db2/history"
	"gitlab.com/tokend/regources"
)

func PopulateChangeRoleRequest(histRequest history.ChangeRoleRequest) (
	r *regources.ChangeRoleRequest, err error,
) {
	r = &regources.ChangeRoleRequest{}
	r.DestinationAccount = histRequest.DestinationAccount
	r.AccountRoleToSet = histRequest.AccountRoleToSet
	r.KYCData = histRequest.KYCData
	r.AllTasks = histRequest.AllTasks
	r.PendingTasks = histRequest.PendingTasks
	r.SequenceNumber = histRequest.SequenceNumber
	return
}
