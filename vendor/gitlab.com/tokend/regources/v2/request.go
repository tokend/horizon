package regources

//Request - details of the request created or reviewed via op
type Request struct {
	RequestID   int64 `json:"request_id,omitempty"`
	IsFulfilled bool  `json:"is_fulfilled"`
}

