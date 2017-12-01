package reviewablerequest

import "gitlab.com/swarmfund/horizon/db2/history"

// RequestState - provides frontend friendly representation of history.ReviewableRequestState
type RequestState struct {
	// RequestStateI  - integer representation of request state
	RequestStateI int32 `json:"request_state_i"`
	// RequestState  - string representation of request state
	RequestState string `json:"request_state"`
}

// Populate - populates requestState from history.ReviewableRequestState
func (r *RequestState) Populate(rawState history.ReviewableRequestState) {
	r.RequestStateI = int32(rawState)
	r.RequestState = rawState.String()
}
