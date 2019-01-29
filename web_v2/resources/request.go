package resources

import (
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

//NewRequestKey - creates new instance of request key
func NewRequestKey(requestId int64) regources.Key {
	return regources.NewKeyInt64(requestId, regources.TypeRequests)
}

//NewRequest - creates new instance of reviewable request
func NewRequest(record history2.ReviewableRequest) regources.ReviewableRequest {
	return regources.ReviewableRequest{
		Key: NewRequestKey(record.ID),
		Attributes: regources.ReviewableRequestAttrs{
			Reference:       record.Reference,
			RejectReason:    record.RejectReason,
			Hash:            record.Hash,
			AllTasks:        record.AllTasks,
			PendingTasks:    record.PendingTasks,
			ExternalDetails: record.ExternalDetails,
			CreatedAt:       record.CreatedAt,
			UpdatedAt:       record.UpdatedAt,
			State:           record.RequestState.String(),
			StateI:          int32(record.RequestState),
		},
		Relationships: regources.ReviewableRequestRelations{
			Requestor: NewAccountKey(record.Requestor).AsRelation(),
			Reviewer: NewAccountKey(record.Reviewer).AsRelation(),
		},
	}
}
