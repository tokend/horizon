package requests

import (
	"gitlab.com/distributed_lab/kit/pgdb"
	"net/http"
	"time"
)

const (
	// FilterTypePollListOwner - defines if we need to filter response by owner
	FilterTypePollListOwner = "owner"

	FilterTypePollListResultProvider = "result_provider"

	// FilterTypePollListMaxEndTime - defines if we need to filter response by max_end_time
	FilterTypePollListMaxEndTime = "max_end_time"
	// FilterTypePollListMaxStartTime - defines if we need to filter response by max_start_time
	FilterTypePollListMaxStartTime = "max_start_time"
	// FilterTypePollListMinStartTime - defines if we need to filter response by min_start_time
	FilterTypePollListMinStartTime = "min_start_time"
	// FilterTypePollListMinEndTime - defines if we need to filter response by min_end_time
	FilterTypePollListMinEndTime = "min_end_time"
	// FilterTypePollListState - defines if we need to filter response by state
	FilterTypePollListState = "state"
	// FilterTypePollListPollType - defines if we need to filter response by poll_type
	FilterTypePollListPollType = "poll_type"

	FilterTypePollListPermissionType = "permission_type"

	FilterTypePollListVoteConfirmation = "vote_confirmation"
)

var filterTypePollListAll = map[string]struct{}{
	FilterTypePollListOwner:            {},
	FilterTypePollListMaxEndTime:       {},
	FilterTypePollListMaxStartTime:     {},
	FilterTypePollListMinStartTime:     {},
	FilterTypePollListMinEndTime:       {},
	FilterTypePollListState:            {},
	FilterTypePollListPollType:         {},
	FilterTypePollListPermissionType:   {},
	FilterTypePollListVoteConfirmation: {},
	FilterTypePollListResultProvider:   {},
}

type GetPollList struct {
	*base
	Filters struct {
		Owner            string     `fig:"owner"`
		ResultProvider   string     `fig:"result_provider"`
		MaxEndTime       *time.Time `fig:"max_end_time"`
		MaxStartTime     *time.Time `fig:"max_start_time"`
		MinStartTime     *time.Time `fig:"min_start_time"`
		MinEndTime       *time.Time `fig:"min_end_time"`
		State            int32      `fig:"state"`
		PollType         uint64     `fig:"poll_type"`
		PermissionType   uint32     `fig:"permission_type"`
		VoteConfirmation bool       `fig:"vote_confirmation"`
	}
	PageParams *pgdb.CursorPageParams
}

func NewGetPollList(r *http.Request) (*GetPollList, error) {
	b, err := newBase(r, baseOpts{
		supportedFilters: filterTypePollListAll,
	})
	if err != nil {
		return nil, err
	}

	pageParams, err := b.getCursorBasedPageParams()
	if err != nil {
		return nil, err
	}

	request := GetPollList{
		base:       b,
		PageParams: pageParams,
	}

	err = b.populateFilters(&request.Filters)
	if err != nil {
		return nil, err
	}

	return &request, nil
}
