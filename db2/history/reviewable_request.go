package history

import "gitlab.com/swarmfund/go/xdr"

type ReviewableRequestState int

// Represents Reviewable request
type ReviewableRequest struct {
	ID          int64  `db:"id"`
	Requestor     string `db:"requestor"`
	Reviewer string `json:"reviewer"`
	Reference *string `json:"reference"`
	RejectReason string `json:"reject_reason"`
	RequestType xdr.ReviewableRequestType `json:"request_type"`

}
