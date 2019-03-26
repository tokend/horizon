package core2

import "gitlab.com/tokend/go/xdr"

type Vote struct {
	PollID  uint64       `db:"poll_id"`
	VoterID string       `db:"voter_id"`
	Choice  xdr.VoteData `db:"choice"`
}
