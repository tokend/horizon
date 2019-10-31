package operations

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	regources "gitlab.com/tokend/regources/generated"
)

type checkSaleStateOpHandler struct {
	manageOfferOpHandler *manageOfferOpHandler
}

// Details returns details about check sale state operation
func (h *checkSaleStateOpHandler) Details(op rawOperation,
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
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	res := opRes.MustCheckSaleStateResult().MustSuccess()

	switch res.Effect.Effect {
	case xdr.CheckSaleStateEffectCanceled, xdr.CheckSaleStateEffectUpdated:
		return h.manageOfferOpHandler.getDeletedOffersEffect(ledgerChanges), nil
	case xdr.CheckSaleStateEffectClosed:
		return h.getParticipationChanges(int64(opBody.MustCheckSaleStateOp().SaleId), res.Effect.MustSaleClosed(), ledgerChanges), nil
	default:
		return nil, errors.From(errors.New("unexpected check sale state result effect"), map[string]interface{}{
			"effect_i": int32(res.Effect.Effect),
			"sale_id":  uint64(res.SaleId),
		})
	}
}

func (h *checkSaleStateOpHandler) getParticipationChanges(orderBookID int64, closedRes xdr.CheckSaleClosedResult,
	ledgerChanges []xdr.LedgerEntryChange,
) []history2.ParticipantEffect {
	// TODO: we are not handling here cases that some parts of offers might be canceled due to rounding
	if len(closedRes.Results) == 0 {
		return nil
	}

	result := make([]history2.ParticipantEffect, 0)
	var totalBaseIssued uint64
	ownerID := h.manageOfferOpHandler.MustAccountID(closedRes.SaleOwner)
	// it does not matter which base balance is used as we are sure that the operation of distribution will be clean
	baseBalanceAddress := closedRes.Results[0].SaleBaseBalance.AsString()
	baseBalanceID := h.manageOfferOpHandler.MustBalanceID(closedRes.Results[0].SaleBaseBalance)
	baseAsset := string(closedRes.Results[0].SaleDetails.BaseAsset)
	removedOffers := h.getRemovedOfferEntries(ledgerChanges)
	for _, assetPairResult := range closedRes.Results {
		sourceOffer := offer{
			OrderBookID:         orderBookID,
			AccountID:           ownerID,
			BaseBalanceID:       baseBalanceID,
			BaseBalanceAddress:  baseBalanceAddress,
			QuoteBalanceID:      h.manageOfferOpHandler.MustBalanceID(assetPairResult.SaleQuoteBalance),
			QuoteBalanceAddress: assetPairResult.SaleQuoteBalance.AsString(),
			BaseAsset:           baseAsset,
			QuoteAsset:          string(assetPairResult.SaleDetails.QuoteAsset),
			IsBuy:               false,
		}
		assetPairMatches, baseIssued := h.manageOfferOpHandler.getMatchesEffects(
			assetPairResult.SaleDetails.OffersClaimed, sourceOffer)

		//We must filter offers deleted because of matches from the leftovers from sale, which were deleted
		for _, match := range assetPairMatches {
			if match.Effect == nil {
				continue
			}
			if match.Effect.Type != history2.EffectTypeMatched {
				continue
			}
		}

		totalBaseIssued += baseIssued
		result = append(result, assetPairMatches...)
	}

	removedOfferEffects := h.getUnlockedEffects(removedOffers)
	result = append(result, removedOfferEffects...)

	// we need to show explicitly that issuance has been perform to ensure that balance history is consistent
	issuanceEffect := history2.ParticipantEffect{
		AccountID: ownerID,
		BalanceID: &baseBalanceID,
		AssetCode: &baseAsset,
		Effect: &history2.Effect{
			Type: history2.EffectTypeIssued,
			Issued: &history2.BalanceChangeEffect{
				Amount: regources.Amount(totalBaseIssued),
			},
		},
	}

	// prepend
	result = append(result, result[0])
	result[0] = issuanceEffect

	return result
}

func (h *checkSaleStateOpHandler) getRemovedOfferEntries(changes []xdr.LedgerEntryChange) map[int64]xdr.OfferEntry {
	statedOffers := make(map[int64]xdr.OfferEntry)
	result := make(map[int64]xdr.OfferEntry)
	removedOffers := make(map[int64]bool)
	for _, change := range changes {
		switch change.Type {
		case xdr.LedgerEntryChangeTypeRemoved:
			{
				if change.Removed.Type != xdr.LedgerEntryTypeOfferEntry {
					continue
				}
				removedOffers[int64(change.Removed.MustOffer().OfferId)] = true
			}
		case xdr.LedgerEntryChangeTypeState:

			{
				if change.State.Data.Type != xdr.LedgerEntryTypeOfferEntry {
					continue
				}
				statedOffers[int64(change.State.Data.MustOffer().OfferId)] = change.State.Data.MustOffer()
			}
		}
	}

	for id := range removedOffers {
		if offer, ok := statedOffers[id]; ok {
			result[id] = offer
		}
	}

	return result
}

func (h *checkSaleStateOpHandler) getUnlockedEffects(removedOffers map[int64]xdr.OfferEntry) []history2.ParticipantEffect {
	result := map[offerID]history2.ParticipantEffect{}
	for _, off := range removedOffers {
		balanceID := off.BaseBalance
		if off.IsBuy {
			balanceID = off.QuoteBalance
		}
		unlockedAmount := off.BaseAmount
		if off.IsBuy {
			unlockedAmount = off.QuoteAmount
		}

		participant := h.manageOfferOpHandler.BalanceEffect(balanceID, &history2.Effect{
			Type: history2.EffectTypeUnlocked,
			Unlocked: &history2.BalanceChangeEffect{
				Amount: regources.Amount(unlockedAmount),
			},
		})

		if off.IsBuy {
			participant.Effect.Unlocked.Fee.CalculatedPercent = regources.Amount(off.Fee)
		}

		result[offerID{
			OrderBookID: uint64(off.OrderBookId),
			OfferID:     uint64(off.OfferId),
		}] = participant
	}

	uniqueResults := make([]history2.ParticipantEffect, 0, len(result))
	for _, participant := range result {
		uniqueResults = append(uniqueResults, participant)
	}

	return uniqueResults
}
