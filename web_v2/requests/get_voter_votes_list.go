package requests

import (
	"gitlab.com/tokend/horizon/bridge"
	"net/http"
)

const (
	IncludeTypeVotersVoteListAccount = "account"
	IncludeTypeVoterVoteListPolls    = "polls"
)

type GetVoterVoteList struct {
	*base
	VoterID    string
	PageParams *bridge.CursorPageParams
}

func NewGetVotersVotes(r *http.Request) (*GetVoterVoteList, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: map[string]struct{}{
			IncludeTypeVoterVoteListPolls:    {},
			IncludeTypeVotersVoteListAccount: {},
		},
	})
	if err != nil {
		return nil, err
	}

	pageParams, err := b.getCursorBasedPageParams()
	if err != nil {
		return nil, err
	}

	id := b.getString("voter")

	return &GetVoterVoteList{
		base:       b,
		PageParams: pageParams,
		VoterID:    id,
	}, nil
}
