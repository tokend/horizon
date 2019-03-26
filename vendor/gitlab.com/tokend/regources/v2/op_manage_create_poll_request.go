package regources

import (
	"time"

	"gitlab.com/tokend/go/xdr"
)

//ManageCreatePollRequestOp - stores details of manage create poll request poll operation
type ManageCreatePollRequestOp struct {
	Key
	Attributes    ManageCreatePollRequestOpAttrs     `json:"attributes"`
	Relationships ManageCreatePollRequestOpRelations `json:"relationships"`
}

//ManageCreatePollRequestOpAttrs - details of ManageCreatePollRequestOp
type ManageCreatePollRequestOpAttrs struct {
	Action xdr.ManageCreatePollRequestAction `json:"action"`
	Create *CreatePollRequestOp              `json:"create,omitempty"`
	Cancel *CancelPollRequestOp              `json:"cancel,omitempty"`
}

type CreatePollRequestOp struct {
	PermissionType           uint64    `json:"permission_type"`
	NumberOfChoices          uint64    `json:"number_of_choices"`
	CreatorDetails           Details   `json:"creator_details"`
	StartTime                time.Time `json:"start_time"`
	EndTime                  time.Time `json:"end_time"`
	ResultProviderID         string    `json:"result_provider_id"`
	VoteConfirmationRequired bool      `json:"vote_confirmation_required"`
	PollData                 PollData  `json:"poll_data"`
	AllTasks                 *uint32   `json:"all_tasks,omitempty"`
}

type CancelPollRequestOp struct {
	RequestID uint64 `json:"request_id"`
}

//ManageCreatePollRequestOpRelations - relationships of ManageCreatePollRequestOp
type ManageCreatePollRequestOpRelations struct {
	Request        *Relation
	ResultProvider *Relation
}
