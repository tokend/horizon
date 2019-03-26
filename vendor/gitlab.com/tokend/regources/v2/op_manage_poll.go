package regources

import (
	"gitlab.com/tokend/go/xdr"
)

//ManagePollOp - stores details of manage poll operation
type ManagePollOp struct {
	Key
	Attributes    ManagePollOpAttrs     `json:"attributes"`
	Relationships ManagePollOpRelations `json:"relationships"`
}

//ManagePollOpAttrs - details of ManagePollOp
type ManagePollOpAttrs struct {
	Action xdr.ManagePollAction `json:"action"`
	PollID int64                `json:"poll_id"`
	Close  *ClosePollOp         `json:"close,omitempty"`
}

type ClosePollOp struct {
	PollResult xdr.PollResult `json:"poll_result"`
	Details    Details        `json:"details"`
}

//ManagePollOpAttrs - relationships of ManageBalanceOp
type ManagePollOpRelations struct {
	Poll Relation
}
