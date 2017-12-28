package history

import (
	"time"

	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/db2"
)

// Represents Reviewable request
type ReviewableRequest struct {
	db2.TotalOrderID
	Requestor    string                    `db:"requestor"`
	Reviewer     string                    `db:"reviewer"`
	Reference    *string                   `db:"reference"`
	RejectReason string                    `db:"reject_reason"`
	RequestType  xdr.ReviewableRequestType `db:"request_type"`
	RequestState ReviewableRequestState    `db:"request_state"`
	Hash         string                    `db:"hash"`
	Details      []byte                    `db:"details"`
	CreatedAt    time.Time                 `db:"created_at"`
	UpdatedAt    time.Time                 `db:"updated_at"`
}
