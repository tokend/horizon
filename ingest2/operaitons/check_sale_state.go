package operaitons

import (
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
		return []history2.ParticipantEffect{source}, nil
	}

	closedRes := res.Effect.MustSaleClosed()
	saleOwnerID := h.pubKeyProvider.GetAccountID(closedRes.SaleOwner)

	var participants []history2.ParticipantEffect

	for _, subRes := range closedRes.Results {
		baseBalanceID := h.pubKeyProvider.GetBalanceID(subRes.SaleBaseBalance)
		quoteBalanceID := h.pubKeyProvider.GetBalanceID(subRes.SaleQuoteBalance)

		participants = append(participants, h.offerHelper.getParticipantsEffects(subRes.SaleDetails.OffersClaimed,
			offerDirection{
				BaseAsset:  subRes.SaleDetails.BaseAsset,
				QuoteAsset: subRes.SaleDetails.QuoteAsset,
				IsBuy:      false,
			}, saleOwnerID, baseBalanceID, quoteBalanceID)...,
		)
	}

	return participants, nil
}
