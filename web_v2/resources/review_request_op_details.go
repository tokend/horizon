package resources

import (
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

func newReviewRequestOpDetails(id int64, details history2.ReviewRequestDetails) *regources.ReviewRequest {
	return &regources.ReviewRequest{
		Key: regources.NewKeyInt64(id, regources.TypeReviewRequest),
		Attributes: regources.ReviewRequestAttrs{
			Action:          details.Action,
			Reason:          details.Reason,
			RequestHash:     details.RequestHash,
			RequestID:       details.RequestID,
			IsFulfilled:     details.IsFulfilled,
			AddedTasks:      details.AddedTasks,
			RemovedTasks:    details.RemovedTasks,
			ExternalDetails: details.ExternalDetails,
		},
	}
}
