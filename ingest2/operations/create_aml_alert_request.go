package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/internal"
	"gitlab.com/tokend/regources/v2"
)

type createAMLAlertReqeustOpHandler struct {
	effectsProvider
}

// Details returns details about create AML alert request operation
func (h *createAMLAlertReqeustOpHandler) Details(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	amlAlertRequest := op.Body.MustCreateAmlAlertRequestOp().AmlAlertRequest

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateAmlAlert,
		CreateAMLAlertRequest: &history2.CreateAMLAlertRequestDetails{
			Amount:         regources.Amount(amlAlertRequest.Amount),
			BalanceAddress: amlAlertRequest.BalanceId.AsString(),
			CreatorDetails: internal.MarshalCustomDetails(amlAlertRequest.CreatorDetails),
		},
	}, nil
}

// ParticipantsEffects returns `locked` effect for account
// which is suspected in illegal obtaining of tokens
func (h *createAMLAlertReqeustOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	amlAlertRequest := opBody.MustCreateAmlAlertRequestOp().AmlAlertRequest
	result := h.BalanceEffectWithAccount(sourceAccountID, amlAlertRequest.BalanceId, &history2.Effect{
		Type: history2.EffectTypeLocked,
		Locked: &history2.BalanceChangeEffect{
			Amount: regources.Amount(amlAlertRequest.Amount),
		},
	})

	isFulfilled := opRes.CreateAmlAlertRequestResult.MustSuccess().Fulfilled
	// request was fulfilled to funds has been withdrawn
	if isFulfilled {
		result = append(result, h.BalanceEffect(amlAlertRequest.BalanceId, &history2.Effect{
			Type: history2.EffectTypeWithdrawn,
			Withdrawn: &history2.BalanceChangeEffect{
				Amount: regources.Amount(amlAlertRequest.Amount),
			},
		}))
	}
	return result, nil
}
