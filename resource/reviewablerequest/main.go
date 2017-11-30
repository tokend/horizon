package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"strconv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// Represents Reviewable request
type ReviewableRequest struct {
	ID           string   `json:"id"`
	Requestor    string   `json:"requestor"`
	Reviewer     string   `json:"reviewer"`
	Reference    *string  `json:"reference"`
	RejectReason string   `json:"reject_reason"`
	Hash         string   `json:"hash"`
	Details      *Details `json:"details"`
	RequestState
}

func (r *ReviewableRequest) Populate(request *history.ReviewableRequest) error {
	r.ID = strconv.FormatUint(request.ID, 10)
	r.Requestor = request.Requestor
	r.Reviewer = request.Reviewer
	r.Reference = request.Reference
	r.RejectReason = request.RejectReason
	r.RequestState.Populate(request.RequestState)
	r.Hash = request.Hash
	r.Details = new(Details)
	err := r.Details.PopulateFromRawJSON(request.RequestType, request.Details)
	if err != nil {
		return errors.Wrap(err, "failed to populate reviewable request details")
	}

	return nil
}
