package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	regources "gitlab.com/tokend/regources/generated"
)

type createManageOfferRequestOpHandler struct {
	effectsProvider
}

// Details returns details about create limits request operation
func (h *createManageOfferRequestOpHandler) Details(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	createManageOfferRequestOp := op.Body.MustCreateManageOfferRequestOp()
	createManageOfferRequestRes := opRes.MustCreateManageOfferRequestResult().MustSuccess()
	manageOfferOp := createManageOfferRequestOp.Request.Op

	var allTasks *uint32
	if createManageOfferRequestOp.AllTasks != nil {
		tasks := uint32(*createManageOfferRequestOp.AllTasks)
		allTasks = &tasks
	}

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateManageOfferRequest,
		CreateManageOfferRequest: &history2.CreateManageOfferRequestDetails{
			ManageOfferDetails: history2.ManageOfferRequest{
				OfferID:     int64(manageOfferOp.OfferId),
				OrderBookID: int64(manageOfferOp.OrderBookId),
				Amount:      regources.Amount(manageOfferOp.Amount),
				Price:       regources.Amount(manageOfferOp.Price),
				IsBuy:       manageOfferOp.IsBuy,
				Fee: regources.Fee{
					CalculatedPercent: regources.Amount(manageOfferOp.Fee),
				},
			},
			AllTasks: allTasks,
			RequestDetails: history2.RequestDetails{
				RequestID:   int64(createManageOfferRequestRes.RequestId),
				IsFulfilled: createManageOfferRequestRes.Fulfilled,
			},
		},
	}, nil
}
