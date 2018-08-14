package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/regources/reviewablerequest2"
)

// Populate - populates requestState from history.ReviewableRequestState
func PopulateRequestState(rawState history.ReviewableRequestState) (
	r reviewablerequest2.RequestState,
) {
	r = reviewablerequest2.RequestState{}
	r.RequestStateI = int32(rawState)
	r.RequestState = rawState.String()
	return
}
