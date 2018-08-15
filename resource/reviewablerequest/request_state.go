package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/regources"
)

// Populate - populates requestState from history.ReviewableRequestState
func PopulateRequestState(rawState history.ReviewableRequestState) (
	r regources.RequestState,
) {
	r = regources.RequestState{}
	r.RequestStateI = int32(rawState)
	r.RequestState = rawState.String()
	return
}
