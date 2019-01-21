package operations

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

type manageOfferOpHandler struct {
	pubKeyProvider IDProvider
}

// Details returns details about manage offer operation
func (h *manageOfferOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (regources.OperationDetails, error) {
	manageOfferOp := op.Body.MustManageOfferOp()
	manageOfferOpRes := opRes.MustManageOfferResult().MustSuccess()

	offerID := int64(manageOfferOp.OfferId)
	isDeleted := manageOfferOpRes.Offer.Effect == xdr.ManageOfferEffectDeleted
	if !isDeleted {
		offerID = int64(manageOfferOpRes.Offer.MustOffer().OfferId)
	}

	return regources.OperationDetails{
		Type: xdr.OperationTypeManageOffer,
		ManageOffer: &regources.ManageOfferDetails{
			OfferID:     offerID,
			OrderBookID: int64(manageOfferOp.OrderBookId),
			BaseAsset:   string(manageOfferOpRes.BaseAsset),
			QuoteAsset:  string(manageOfferOpRes.QuoteAsset),
			Amount:      regources.Amount(manageOfferOp.Amount),
			Price:       regources.Amount(manageOfferOp.Price),
			IsBuy:       manageOfferOp.IsBuy,
			Fee: regources.Fee{
				CalculatedPercent: regources.Amount(manageOfferOp.Fee),
			},
			IsDeleted: isDeleted,
		},
	}, nil
}

// ParticipantsEffects can return `matched` and `locked` effects if offer created
// returns `unlocked` effects if offer canceled (deleted by user)
func (h *manageOfferOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	manageOfferOp := opBody.MustManageOfferOp()
	manageOfferOpRes := opRes.MustManageOfferResult().MustSuccess()

	if manageOfferOp.Amount != 0 {
		return h.getNewOfferEffect(manageOfferOp, manageOfferOpRes, source), nil
	}

	deletedOfferEffects := h.getDeletedOffersEffect(ledgerChanges)
	if len(deletedOfferEffects) != 1 {
		return nil, errors.From(errors.New("Unexpected number of deleted offer for manage offer delete"), logan.F{
			"expected": 1,
			"actual":   len(deletedOfferEffects),
		})
	}

	return deletedOfferEffects, nil
}

func (h *manageOfferOpHandler) getNewOfferEffect(op xdr.ManageOfferOp,
	res xdr.ManageOfferSuccessResult, source history2.ParticipantEffect,
) []history2.ParticipantEffect {
	participants, _ := h.getMatchesEffects(res.OffersClaimed, offer{
		OrderBookID:         int64(op.OrderBookId),
		AccountID:           source.AccountID,
		BaseBalanceID:       h.pubKeyProvider.MustBalanceID(op.BaseBalance),
		BaseBalanceAddress:  op.BaseBalance.AsString(),
		QuoteBalanceID:      h.pubKeyProvider.MustBalanceID(op.QuoteBalance),
		QuoteBalanceAddress: op.QuoteBalance.AsString(),
		BaseAsset:           string(res.BaseAsset),
		QuoteAsset:          string(res.QuoteAsset),
		IsBuy:               op.IsBuy,
	})
	if res.Offer.Effect == xdr.ManageOfferEffectDeleted {
		return participants
	}

	// we need to handle amount which was not matched
	newOffer := res.Offer.MustOffer()
	source.AssetCode = new(string)
	source.BalanceID = new(uint64)
	source.Effect = regources.Effect{
		Type:   regources.EffectTypeLocked,
		Locked: &regources.BalanceChangeEffect{},
	}
	if newOffer.IsBuy {
		*source.BalanceID = h.pubKeyProvider.MustBalanceID(newOffer.QuoteBalance)
		*source.AssetCode = string(newOffer.Quote)
		source.Effect.Locked.Amount = regources.Amount(newOffer.QuoteAmount)
		source.Effect.Locked.Fee.CalculatedPercent = regources.Amount(newOffer.PercentFee)
	} else {
		*source.BalanceID = h.pubKeyProvider.MustBalanceID(newOffer.BaseBalance)
		*source.AssetCode = string(newOffer.Base)
		source.Effect.Locked.Amount = regources.Amount(newOffer.BaseAmount)
	}

	participants = append(participants, source)
	return participants
}

// getDeletedOffersEffect - creates participant effects for all offer entries with change type `State`
func (h *manageOfferOpHandler) getDeletedOffersEffect(ledgerChanges []xdr.LedgerEntryChange,
) []history2.ParticipantEffect {
	result := make([]history2.ParticipantEffect, 0)
	for _, change := range ledgerChanges {
		if change.Type != xdr.LedgerEntryChangeTypeState {
			continue
		}

		if change.MustState().Data.Type != xdr.LedgerEntryTypeOfferEntry {
			continue
		}

		deletedOffer := change.MustState().Data.MustOffer()

		participant := history2.ParticipantEffect{
			AccountID: h.pubKeyProvider.MustAccountID(deletedOffer.OwnerId),
			BalanceID: new(uint64),
			AssetCode: new(string),
			Effect: regources.Effect{
				Type:     regources.EffectTypeUnlocked,
				Unlocked: &regources.BalanceChangeEffect{},
			},
		}

		if deletedOffer.IsBuy {
			*participant.BalanceID = h.pubKeyProvider.MustBalanceID(deletedOffer.QuoteBalance)
			*participant.AssetCode = string(deletedOffer.Quote)
			participant.Effect.Unlocked.Amount = regources.Amount(deletedOffer.QuoteAmount)
			participant.Effect.Unlocked.Fee.CalculatedPercent = regources.Amount(deletedOffer.PercentFee)
		} else {
			*participant.BalanceID = h.pubKeyProvider.MustBalanceID(deletedOffer.BaseBalance)
			*participant.AssetCode = string(deletedOffer.Base)
			participant.Effect.Unlocked.Amount = regources.Amount(deletedOffer.BaseAmount)
		}

		result = append(result, participant)
	}

	return result
}

type offer struct {
	OrderBookID         int64
	AccountID           uint64
	BaseBalanceID       uint64
	BaseBalanceAddress  string
	QuoteBalanceID      uint64
	QuoteBalanceAddress string
	BaseAsset           string
	QuoteAsset          string
	IsBuy               bool
}

// getParticipantsEffects - returns participants effects based on the provided matches and total base amount
func (h *manageOfferOpHandler) getMatchesEffects(claimOfferAtoms []xdr.ClaimOfferAtom,
	sourceOffer offer) ([]history2.ParticipantEffect, uint64) {
	var totalBaseAmount uint64
	result := make([]history2.ParticipantEffect, 0, len(claimOfferAtoms)*4)
	for _, matchedOffer := range claimOfferAtoms {
		totalBaseAmount += uint64(matchedOffer.BaseAmount)

		result = h.addParticipantEffects(result, offer{
			OrderBookID:         sourceOffer.OrderBookID,
			AccountID:           h.pubKeyProvider.MustAccountID(matchedOffer.BAccountId),
			BaseBalanceID:       h.pubKeyProvider.MustBalanceID(matchedOffer.BaseBalance),
			QuoteBalanceID:      h.pubKeyProvider.MustBalanceID(matchedOffer.QuoteBalance),
			BaseBalanceAddress:  matchedOffer.BaseBalance.AsString(),
			QuoteBalanceAddress: matchedOffer.QuoteBalance.AsString(),
			BaseAsset:           sourceOffer.BaseAsset,
			QuoteAsset:          sourceOffer.QuoteAsset,
			IsBuy:               !sourceOffer.IsBuy,
		}, int64(matchedOffer.OfferId), matchedOffer.BaseAmount, matchedOffer.QuoteAmount, matchedOffer.CurrentPrice,
			matchedOffer.BFeePaid)

		result = h.addParticipantEffects(result, sourceOffer, 0, matchedOffer.BaseAmount,
			matchedOffer.QuoteAmount, matchedOffer.CurrentPrice, matchedOffer.AFeePaid)
	}

	return result, totalBaseAmount
}

func (h *manageOfferOpHandler) addParticipantEffects(participants []history2.ParticipantEffect,
	offer offer, id int64, baseAmount, quoteAmount, price, fee xdr.Int64) []history2.ParticipantEffect {
	baseBalanceEffect := regources.ParticularBalanceChangeEffect{
		BalanceAddress: offer.BaseBalanceAddress,
		AssetCode:      offer.BaseAsset,
		BalanceChangeEffect: regources.BalanceChangeEffect{
			Amount: regources.Amount(baseAmount),
		},
	}

	quoteBalanceEffect := regources.ParticularBalanceChangeEffect{
		BalanceAddress: offer.QuoteBalanceAddress,
		AssetCode:      offer.QuoteAsset,
		BalanceChangeEffect: regources.BalanceChangeEffect{
			Amount: regources.Amount(quoteAmount),
			Fee: regources.Fee{
				CalculatedPercent: regources.Amount(fee),
			},
		},
	}

	matchedOfferEffect := regources.Effect{
		Type: regources.EffectTypeMatched,
		Matched: &regources.MatchEffect{
			OfferID:     id,
			OrderBookID: offer.OrderBookID,
			Price:       regources.Amount(price),
		},
	}

	if offer.IsBuy {
		matchedOfferEffect.Matched.Funded = baseBalanceEffect
		matchedOfferEffect.Matched.Charged = quoteBalanceEffect
	} else {
		matchedOfferEffect.Matched.Funded = quoteBalanceEffect
		matchedOfferEffect.Matched.Charged = baseBalanceEffect
	}

	return append(participants, history2.ParticipantEffect{
		AccountID: offer.AccountID,
		BalanceID: &offer.BaseBalanceID,
		AssetCode: &offer.BaseAsset,
		Effect:    matchedOfferEffect,
	}, history2.ParticipantEffect{
		AccountID: offer.AccountID,
		BalanceID: &offer.QuoteBalanceID,
		AssetCode: &offer.QuoteAsset,
		Effect:    matchedOfferEffect,
	})
}
