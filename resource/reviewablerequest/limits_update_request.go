package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/regources/reviewablerequest2"
)

func PopulateLimitsUpdateRequest(histRequest history.LimitsUpdateRequest) (
	r *reviewablerequest2.LimitsUpdateRequest, err error,
) {
	r = &reviewablerequest2.LimitsUpdateRequest{}
	r.Details = histRequest.Details
	r.DocumentHash = histRequest.DocumentHash
	return
}
