package operaitons

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type checkSaleStateOpHandler struct {
	pubKeyProvider  publicKeyProvider
	offerHelper     offerHelper
	balanceProvider balanceProvider
}

// OperationDetails returns details about check sale state operation
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

// ParticipantsEffects returns sale owner and participants `matched` effects if sale closed
// returns `unlocked` effects if sale canceled or updated
func (h *checkSaleStateOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	res := opRes.MustCheckSaleStateResult().MustSuccess()

	var result []history2.ParticipantEffect
	var err error
	switch res.Effect.Effect {
	case xdr.CheckSaleStateEffectCanceled:
		result, err = h.getSaleAntesEffects(history2.EffectTypeUnlocked, ledgerChanges)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get effects from sale antes")
		}
		fallthrough
	case xdr.CheckSaleStateEffectUpdated:
		return append(result, h.getDeletedParticipants(ledgerChanges)...), nil
	case xdr.CheckSaleStateEffectClosed:
		result, err = h.getSaleAntesEffects(history2.EffectTypeChargedFromLocked, ledgerChanges)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get effects from sale antes")
		}
		participants, err := h.getApprovedParticipants(res.Effect.MustSaleClosed())
		if err != nil {
			return nil, errors.Wrap(err, "failed to approved participants", map[string]interface{}{
				"sale_id": uint64(res.SaleId),
			})
		}

		return append(result, participants...), nil
	default:
		return nil, errors.From(errors.New("unexpected check sale state result effect"), map[string]interface{}{
			"effect_i": int32(res.Effect.Effect),
			"sale_id":  uint64(res.SaleId),
		})
	}
}

func (h *checkSaleStateOpHandler) getApprovedParticipants(closedRes xdr.CheckSaleClosedResult,
) ([]history2.ParticipantEffect, error) {
	if len(closedRes.Results) == 0 {
		return nil, errors.New("expected not empty results")
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
	}), nil
}

func (h *checkSaleStateOpHandler) getDeletedParticipants(ledgerChanges []xdr.LedgerEntryChange,
) []history2.ParticipantEffect {
	var result []history2.ParticipantEffect

	deletedOffers := h.offerHelper.getStateOffers(ledgerChanges)

	for _, offer := range deletedOffers {
		balanceID := h.pubKeyProvider.GetBalanceID(offer.BaseBalance)
		asset := offer.Base
		unlockedAmount := offer.BaseAmount
		if offer.IsBuy {
			balanceID = h.pubKeyProvider.GetBalanceID(offer.QuoteBalance)
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

func (h *checkSaleStateOpHandler) getSaleAntesEffects(effectType history2.EffectType,
	ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	var result []history2.ParticipantEffect

	saleAntes := h.getStatedSaleAntes(ledgerChanges)
	for _, saleAnte := range saleAntes {
		balance := h.balanceProvider.GetBalanceByID(saleAnte.ParticipantBalanceId)

		effect := history2.Effect{
			Type: effectType,
		}
		switch effect.Type {
		case history2.EffectTypeChargedFromLocked:
			effect.ChargedFromLocked = &history2.ChargedFromLockedEffect{
				Amount: amount.StringU(uint64(saleAnte.Amount)),
			}
		case history2.EffectTypeUnlocked:
			effect.Unlocked = &history2.UnlockedEffect{
				Amount: amount.StringU(uint64(saleAnte.Amount)),
			}
		default:
			return nil, errors.New("unexpected effect type for sale ante entry")
		}

		result = append(result, history2.ParticipantEffect{
			AccountID: balance.AccountID,
			BalanceID: &balance.ID,
			AssetCode: &balance.AssetCode,
			Effect:    effect,
		})
	}

	return result, nil
}

func (h *checkSaleStateOpHandler) getStatedSaleAntes(ledgerChanges []xdr.LedgerEntryChange,
) []xdr.SaleAnteEntry {
	var result []xdr.SaleAnteEntry

	for _, change := range ledgerChanges {
		if change.Type != xdr.LedgerEntryChangeTypeState {
			continue
		}

		if change.MustState().Data.Type != xdr.LedgerEntryTypeSaleAnte {
			continue
		}

		result = append(result, change.MustState().Data.MustSaleAnte())
	}

	return result
}
