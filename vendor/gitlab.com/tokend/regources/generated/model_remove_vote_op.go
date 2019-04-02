package regources

type RemoveVoteOp struct {
	// id of the poll to remove vote from
	PollId int64 `json:"poll_id"`
}
