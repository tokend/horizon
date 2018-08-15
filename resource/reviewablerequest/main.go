package reviewablerequest

import (
	"strconv"
	"time"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/regources"
)

func PopulateReviewableRequest(request *history.ReviewableRequest) (
	r *regources.ReviewableRequest, err error,
) {
	r = &regources.ReviewableRequest{}
	r.ID = strconv.FormatInt(request.ID, 10)
	r.PT = request.PagingToken()
	r.Requestor = request.Requestor
	r.Reviewer = request.Reviewer
	r.Reference = request.Reference
	r.RejectReason = request.RejectReason
	r.RequestState = PopulateRequestState(request.RequestState)
	r.Hash = request.Hash
	r.CreatedAt = request.CreatedAt.Format(time.RFC3339)
	r.UpdatedAt = request.UpdatedAt.Format(time.RFC3339)

	r.Details, err = PopulateDetails(request.RequestType, request.Details)
	if err != nil {
		return nil, errors.Wrap(err, "failed to populate reviewable request details")
	}

	return
}
