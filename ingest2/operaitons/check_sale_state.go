package operaitons

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type checkSaleStateOpHandler struct {
	pubKeyProvider publicKeyProvider
	offerHelper    offerHelper
}

func (h *checkSaleStateOpHandler) OperationDetails(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	return history2.OperationDetails{
		Type: xdr.OperationTypeCheckSaleState,
		CheckSaleState: &history2.CheckSaleStateDetails{
			SaleID: int64(op.Body.MustCheckSaleStateOp().SaleId),
			Effect: opRes.MustCheckSaleStateResult().MustSuccess().Effect.Effect,
		},
	}, nil
}

func (h *checkSaleStateOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	res := opRes.MustCheckSaleStateResult().MustSuccess()

	if res.Effect.Effect != xdr.CheckSaleStateEffectClosed {
		return nil, nil
	}

	closedRes := res.Effect.MustSaleClosed()
	if len(closedRes.Results) == 0 {
		return nil, nil
	}

	saleOwnerID := h.pubKeyProvider.GetAccountID(closedRes.SaleOwner)
	baseBalanceID := h.pubKeyProvider.GetBalanceID(closedRes.Results[0].SaleBaseBalance)
	baseAsset := closedRes.Results[0].SaleDetails.BaseAsset

	var participants []history2.ParticipantEffect
	var issuedAmount uint64 = 0

	for _, subRes := range closedRes.Results {
		quoteBalanceID := h.pubKeyProvider.GetBalanceID(subRes.SaleQuoteBalance)

		newParticipants, baseAmount := h.offerHelper.getParticipantsEffects(
			subRes.SaleDetails.OffersClaimed,
			offerDirection{
				BaseAsset:  subRes.SaleDetails.BaseAsset,
				QuoteAsset: subRes.SaleDetails.QuoteAsset,
				IsBuy:      false,
			}, saleOwnerID, baseBalanceID, quoteBalanceID)

		participants = append(participants, newParticipants...)

		issuedAmount += baseAmount
	}

	return append(participants, history2.ParticipantEffect{
		AccountID: saleOwnerID,
		BalanceID: &baseBalanceID,
		AssetCode: &baseAsset,
		Effect: history2.Effect{
			Type: history2.EffectTypeFunded,
			Issuance: &history2.IssuanceEffect{
				Amount: amount.StringU(issuedAmount),
			},
		},
	}), nil
}
