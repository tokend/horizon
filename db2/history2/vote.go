package history2

import "github.com/lib/pq"

// Vote - represents choice of voting campaign participant
type Vote struct {
	PollID  int64         `db:"poll_id"`
	VoterID string        `db:"voter_id"`
	Choices pq.Int64Array `db:"choices"`
}
