package requests

import (
	"gitlab.com/tokend/horizon/bridge"
	"net/http"
)

type GetVoteList struct {
	*base
	PollID     int64
	PageParams *bridge.CursorPageParams
}

func NewGetVoteList(r *http.Request) (*GetVoteList, error) {
	b, err := newBase(r, baseOpts{})
	if err != nil {
		return nil, err
	}

	pageParams, err := b.getCursorBasedPageParams()
	if err != nil {
		return nil, err
	}

	id, err := b.getUint64ID()
	if err != nil {
		return nil, err
	}

	return &GetVoteList{
		base:       b,
		PageParams: pageParams,
		PollID:     int64(id),
	}, nil
}
