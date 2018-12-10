package operaitons

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type manageInvoiceRequestOpHandler struct {
}

// OperationDetails returns details about manage invoice request operation
func (h *manageInvoiceRequestOpHandler) OperationDetails(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageInvoiceRequestOp := op.Body.MustManageInvoiceRequestOp()
	manageInvoiceRequestOpRes := opRes.MustManageInvoiceRequestResult().MustSuccess()

	opDetails := history2.OperationDetails{
		Type: xdr.OperationTypeManageInvoiceRequest,
		ManageInvoiceRequest: &history2.ManageInvoiceRequestDetails{
			Action: manageInvoiceRequestOp.Details.Action,
		},
	}

	switch manageInvoiceRequestOp.Details.Action {
	case xdr.ManageInvoiceRequestActionCreate:
		creationDetails := manageInvoiceRequestOp.Details.MustInvoiceRequest()

		var contractID *int64
		if creationDetails.ContractId != nil {
			contractIDInt := int64(*creationDetails.ContractId)
			contractID = &contractIDInt
		}

		opDetails.ManageInvoiceRequest.Create = &history2.CreateInvoiceRequestDetails{
			Amount:     amount.StringU(uint64(creationDetails.Amount)),
			Sender:     creationDetails.Sender.Address(),
			RequestID:  int64(manageInvoiceRequestOpRes.Details.MustResponse().RequestId),
			Asset:      creationDetails.Asset,
			ContractID: contractID,
			Details:    customDetailsUnmarshal([]byte(creationDetails.Details)),
		}
	case xdr.ManageInvoiceRequestActionRemove:
		opDetails.ManageInvoiceRequest.Remove = &history2.RemoveInvoiceRequestDetails{
			RequestID: int64(manageInvoiceRequestOp.Details.MustRequestId()),
		}
	}

	return opDetails, nil
}

func (h *manageInvoiceRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	_ xdr.OperationResultTr, source history2.ParticipantEffect, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{source}, nil
}
