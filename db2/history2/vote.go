package history2

import "gitlab.com/tokend/regources/v2"

// Vote - represents choice of voting campaign participant
type Vote struct {
	PollID   int64              `db:"poll_id"`
	VoterID  string             `db:"voter_id"`
	VoteData regources.VoteData `db:"data"`
}
