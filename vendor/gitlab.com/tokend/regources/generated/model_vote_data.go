/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

import "gitlab.com/tokend/go/xdr"

type VoteData struct {
	// type of the poll
	PollType     xdr.PollType `json:"poll_type"`
	CreationTime *uint64      `json:"creation_time,omitempty"`
	SingleChoice *uint64      `json:"single_choice,omitempty"`
}
