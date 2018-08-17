package reviewablerequest

import (
			"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/regources"
)

func PopulateReviewableRequest(request *history.ReviewableRequest) (r *regources.ReviewableRequest, err error) {
	r = &regources.ReviewableRequest{}

	r.ID = uint64(request.ID)
	r.PT = request.PagingToken()
	r.Requestor = request.Requestor
	r.Reviewer = request.Reviewer
	r.Reference = request.Reference
	r.RejectReason = request.RejectReason
	r.StateName = request.RequestState.String()
	r.State = int32(request.RequestState)
	r.Hash = request.Hash
	r.CreatedAt = regources.Time(request.CreatedAt)
	r.UpdatedAt = regources.Time(request.UpdatedAt)

	r.Details, err = PopulateDetails(request.RequestType, request.Details)
	if err != nil {
		return nil, errors.Wrap(err, "failed to populate reviewable request details")
	}

	return
}
