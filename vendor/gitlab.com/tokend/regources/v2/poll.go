package regources

import (
	"encoding/json"
	"time"

	"gitlab.com/tokend/go/xdr"
)

type PollState int

const (
	PollStateOpen PollState = iota + 1
	PollStateClosed
)

var pollStateMap = map[PollState]string{
	PollStateOpen:   "open",
	PollStateClosed: "closed",
}

func (s PollState) MarshalJSON() ([]byte, error) {
	return json.Marshal(Flag{
		Name:  pollStateMap[s],
		Value: int(s),
	})
}

//String - converts int enum to string
func (s PollState) String() string {
	return pollStateMap[s]
}

//PollResponse - response for poll handler
type PollResponse struct {
	Data     Poll     `json:"data"`
	Included Included `json:"included"`
}

type PollsResponse struct {
	Links    *Links   `json:"links"`
	Data     []Poll   `json:"data"`
	Included Included `json:"included"`
}

// Poll - Resource object representing PollEntry
type Poll struct {
	Key
	Attributes    PollAttrs     `json:"attributes"`
	Relationships PollRelations `json:"relationships"`
}

type PollAttrs struct {
	PollType                 xdr.PollType `json:"poll_type"`
	PermissionType           uint64       `json:"permission_type"`
	NumberOfChoices          uint64       `json:"number_of_choices"`
	StartTime                time.Time    `json:"start_time"`
	EndTime                  time.Time    `json:"end_time"`
	VoteConfirmationRequired bool         `json:"vote_confirmation_required"`
	Details                  Details      `json:"details"`
	PollState                PollState    `json:"poll_state"`
	VotesCount               []VoteCount  `json:"votes_count"`
}

type VoteCount struct {
	Choice uint64 `json:"choice"`
	Count  uint32 `json:"count"`
}

type PollRelations struct {
	Owner          Relation           `json:"owner"`
	ResultProvider Relation           `json:"result_provider"`
	Votes          RelationCollection `json:"votes"`
}
type PollData struct {
	Type xdr.PollType
}
