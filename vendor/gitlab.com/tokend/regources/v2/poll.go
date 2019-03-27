package regources

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gitlab.com/distributed_lab/logan/v3/errors"

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
	PollData                 PollData    `json:"poll_data"`
	PermissionType           uint64      `json:"permission_type"`
	NumberOfChoices          uint64      `json:"number_of_choices"`
	StartTime                time.Time   `json:"start_time"`
	EndTime                  time.Time   `json:"end_time"`
	VoteConfirmationRequired bool        `json:"vote_confirmation_required"`
	Details                  Details     `json:"details"`
	PollState                PollState   `json:"poll_state"`
	VotesCount               []VoteCount `json:"votes_count"`
}

type VoteCount struct {
	Choice uint64 `json:"choice"`
	Count  uint32 `json:"count"`
}

type PollRelations struct {
	Owner          *Relation           `json:"owner"`
	ResultProvider *Relation           `json:"result_provider"`
	Votes          *RelationCollection `json:"votes"`
}
type PollData struct {
	Type xdr.PollType `json:"type"`
}

//Value - implements db driver method for auto marshal
func (r PollData) Value() (driver.Value, error) {
	result, err := json.Marshal(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal poll data")
	}

	return result, nil
}

//Scan - implements db driver method for auto unmarshal
func (r *PollData) Scan(src interface{}) error {
	var data []byte
	switch rawData := src.(type) {
	case []byte:
		data = rawData
	case string:
		data = []byte(rawData)
	default:
		return errors.New("Unexpected type for jsonb")
	}

	err := json.Unmarshal(data, r)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal poll data")
	}

	return nil
}
