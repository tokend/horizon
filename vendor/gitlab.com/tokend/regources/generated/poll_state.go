package regources

import "encoding/json"

type PollState int

const (
	PollStateOpen PollState = iota + 1
	PollStatePassed
	PollStateFailed
)

var pollStateStr = map[PollState]string{
	PollStateOpen:   "open",
	PollStatePassed: "passed",
	PollStateFailed: "failed",
}

func (s PollState) String() string {
	return pollStateStr[s]
}

func (s PollState) MarshalJSON() ([]byte, error) {
	return json.Marshal(Flag{
		Name:  pollStateStr[s],
		Value: int32(s),
	})
}
