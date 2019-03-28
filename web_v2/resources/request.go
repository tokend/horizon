package resources

import (
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/rgenerated"
)

//NewRequestKey - creates new instance of request key
func NewRequestKey(requestId int64) rgenerated.Key {
	return rgenerated.NewKeyInt64(requestId, rgenerated.REQUESTS)
}

//NewRequest - creates new instance of reviewable request
func NewRequest(record history2.ReviewableRequest) rgenerated.ReviewableRequest {
	return rgenerated.ReviewableRequest{
		Key: NewRequestKey(record.ID),
		Attributes: rgenerated.ReviewableRequestAttributes{
			Reference:       record.Reference,
			RejectReason:    record.RejectReason,
			Hash:            record.Hash,
			AllTasks:        record.AllTasks,
			PendingTasks:    record.PendingTasks,
			ExternalDetails: record.ExternalDetails.AsRegourcesDetails(),
			CreatedAt:       record.CreatedAt,
			UpdatedAt:       record.UpdatedAt,
			XdrType:         record.RequestType,

			// TODO shouldn't those look like `state: {str: "", int: 2}`?
			//  or just drop int part?
			State:  record.RequestState.String(),
			StateI: int32(record.RequestState),
		},
		Relationships: rgenerated.ReviewableRequestRelationships{
			Requestor: NewAccountKey(record.Requestor).AsRelation(),
			Reviewer:  NewAccountKey(record.Reviewer).AsRelation(),
		},
	}
}
