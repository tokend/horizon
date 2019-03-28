package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/internal"
	regources "gitlab.com/tokend/regources/v2/generated"
)

type createSaleRequestOpHandler struct {
	effectsProvider
}

// Details returns details about create sale request operation
func (h *createSaleRequestOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	createSaleRequest := op.Body.MustCreateSaleCreationRequestOp().Request

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateSaleRequest,
		CreateSaleRequest: &history2.CreateSaleRequestDetails{
			RequestID:         int64(opRes.MustCreateSaleCreationRequestResult().MustSuccess().RequestId),
			BaseAsset:         string(createSaleRequest.BaseAsset),
			DefaultQuoteAsset: string(createSaleRequest.DefaultQuoteAsset),
			StartTime:         internal.TimeFromXdr(createSaleRequest.StartTime),
			EndTime:           internal.TimeFromXdr(createSaleRequest.EndTime),
			HardCap:           regources.Amount(createSaleRequest.HardCap),
			SoftCap:           regources.Amount(createSaleRequest.SoftCap),
			CreatorDetails:    internal.MarshalCustomDetails(createSaleRequest.CreatorDetails),
		},
	}, nil
}
