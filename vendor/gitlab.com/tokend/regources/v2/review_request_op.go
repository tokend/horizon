package regources

import "gitlab.com/tokend/go/xdr"

//ReviewRequestAttrs - details of corresponding op
type ReviewRequest struct {
	Key
	Attributes ReviewRequestAttrs `json:"attributes"`
	//TODO: add review request details as relation
}

//ReviewRequestAttrs - details of corresponding op
type ReviewRequestAttrs struct {
	Action          xdr.ReviewRequestOpAction `json:"action"`
	Reason          string                    `json:"reason"`
	RequestHash     string                    `json:"request_hash"`
	RequestID       int64                     `json:"request_id"`
	IsFulfilled     bool                      `json:"is_fulfilled"`
	AddedTasks      uint32                    `json:"added_tasks"`
	RemovedTasks    uint32                    `json:"removed_tasks"`
	ExternalDetails Details                   `json:"external_details"`
}
