package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/internal"
)

type createManageLimitsRequestOpHandler struct {
	effectsProvider
}

// Details returns details about create limits request operation
func (h *createManageLimitsRequestOpHandler) Details(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	createManageLimitsRequestOp := op.Body.MustCreateManageLimitsRequestOp()

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateManageLimitsRequest,
		CreateManageLimitsRequest: &history2.CreateManageLimitsRequestDetails{
			CreatorDetails: internal.MarshalCustomDetails(
				createManageLimitsRequestOp.ManageLimitsRequest.CreatorDetails),
			RequestID: int64(opRes.MustCreateManageLimitsRequestResult().MustSuccess().ManageLimitsRequestId),
		},
	}, nil
}
