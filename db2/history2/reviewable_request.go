package history2

import (
	"gitlab.com/tokend/horizon/bridge"
	"time"

	"gitlab.com/tokend/go/xdr"
)

// ReviewableRequest - represents details of the reviewable request
type ReviewableRequest struct {
	bridge.TotalOrderID
	Requestor       string                    `db:"requestor"`
	Reviewer        string                    `db:"reviewer"`
	Reference       *string                   `db:"reference"`
	RejectReason    string                    `db:"reject_reason"`
	RequestType     xdr.ReviewableRequestType `db:"request_type"`
	RequestState    ReviewableRequestState    `db:"request_state"`
	Hash            string                    `db:"hash"`
	CreatedAt       time.Time                 `db:"created_at"`
	UpdatedAt       time.Time                 `db:"updated_at"`
	Details         ReviewableRequestDetails  `db:"details"`
	AllTasks        uint32                    `db:"all_tasks"`
	PendingTasks    uint32                    `db:"pending_tasks"`
	ExternalDetails bridge.Details            `db:"external_details"`
}
