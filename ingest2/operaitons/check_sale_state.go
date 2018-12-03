package operaitons

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type checkSaleStateOpHandler struct {
	pubKeyProvider        publicKeyProvider
	offerHelper           offerHelper
	ledgerChangesProvider ledgerChangesProvider
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

	switch res.Effect.Effect {
	case xdr.CheckSaleStateEffectUpdated:
		fallthrough
	case xdr.CheckSaleStateEffectCanceled:
		return h.getDeletedParticipants(), nil
	case xdr.CheckSaleStateEffectClosed:
		return h.getApprovedParticipants(res.Effect.MustSaleClosed()), nil
	default:
		return nil, errors.From(errors.New("unexpected check sale state result effect"), map[string]interface{}{
			"effect_i": int32(res.Effect.Effect),
		})
	}
}

func (h *checkSaleStateOpHandler) getApprovedParticipants(closedRes xdr.CheckSaleClosedResult) []history2.ParticipantEffect {
	if len(closedRes.Results) == 0 {
		return nil
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
			Funded: &history2.FundedEffect{
				Amount: amount.StringU(issuedAmount),
			},
		},
	})
}

func (h *checkSaleStateOpHandler) getDeletedParticipants() []history2.ParticipantEffect {
	var result []history2.ParticipantEffect

	deletedOffers := h.offerHelper.getStateOffers()

	for _, offer := range deletedOffers {
		baseBalanceID := h.pubKeyProvider.GetBalanceID(offer.BaseBalance)
		quoteBalanceID := h.pubKeyProvider.GetBalanceID(offer.QuoteBalance)

		balanceID := baseBalanceID
		asset := offer.Base
		unlockedAmount := offer.BaseAmount
		if offer.IsBuy {
			balanceID = quoteBalanceID
			asset = offer.Quote
			unlockedAmount = offer.QuoteAmount + offer.Fee
		}

		participantEffect := history2.ParticipantEffect{
			AccountID: h.pubKeyProvider.GetAccountID(offer.OwnerId),
			BalanceID: &balanceID,
			AssetCode: &asset,
			Effect: history2.Effect{
				Type: history2.EffectTypeUnlocked,
				Unlocked: &history2.UnlockedEffect{
					Amount: amount.String(int64(unlockedAmount)),
				},
			},
		}

		result = append(result, participantEffect)
	}

	return result
}
