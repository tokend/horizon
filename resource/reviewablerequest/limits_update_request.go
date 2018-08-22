package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/regources"
)

func PopulateLimitsUpdateRequest(histRequest history.LimitsUpdateRequest) (
	r *regources.LimitsUpdateRequest, err error,
) {
	r = &regources.LimitsUpdateRequest{}
	r.Details = histRequest.Details
	r.DocumentHash = histRequest.DocumentHash
	return
}
