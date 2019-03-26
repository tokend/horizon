package regources

import (
	"gitlab.com/tokend/go/xdr"
)

//ManageVoteOp - stores details of manage vote operation
type ManageVoteOp struct {
	Key
	Attributes    ManageVoteOpAttrs     `json:"attributes"`
	Relationships ManageVoteOpRelations `json:"relationships"`
}

//ManageCreateVoteRequestOpAttrs - details of ManageCreateVoteRequestOp
type ManageVoteOpAttrs struct {
	Action xdr.ManageVoteAction `json:"action"`
	Create *CreateVoteOp        `json:"create,omitempty"`
	Remove *RemoveVoteOp        `json:"remove,omitempty"`
}

type CreateVoteOp struct {
	PollID   int64        `json:"poll_id"`
	VoteData xdr.VoteData `json:"vote_data"`
}

type RemoveVoteOp struct {
	PollID int64 `json:"poll_id"`
}

//ManageVoteOpAttrs - relationships of ManageVoteOp
type ManageVoteOpRelations struct {
	Poll           Relation
	Voter          Relation
	ResultProvider *Relation
}
