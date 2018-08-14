package reviewablerequest2

// RequestState - provides frontend friendly representation of history.ReviewableRequestState
type RequestState struct {
	// RequestStateI  - integer representation of request state
	RequestStateI int32 `json:"request_state_i"`
	// RequestState  - string representation of request state
	RequestState string `json:"request_state"`
}
