/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

import (
	"time"

	"gitlab.com/tokend/go/xdr"
)

type VoteData struct {
	CreationTime *time.Time `json:"creation_time,omitempty"`
	CustomChoice *Details   `json:"custom_choice,omitempty"`
	// type of the poll
	PollType     xdr.PollType `json:"poll_type"`
	SingleChoice *uint64      `json:"single_choice,omitempty"`
}
