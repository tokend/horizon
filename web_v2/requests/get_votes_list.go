package requests

import (
	"net/http"

	"gitlab.com/tokend/horizon/db2"
)

type GetVoteList struct {
	*base
	PollID     int64
	PageParams *db2.OffsetPageParams
}

func NewGetVoteList(r *http.Request) (*GetVoteList, error) {
	b, err := newBase(r, baseOpts{})
	if err != nil {
		return nil, err
	}

	pageParams, err := b.getOffsetBasedPageParams()
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
