package operations

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)


type manageOfferOpHandler struct {
	pubKeyProvider publicKeyProvider
}

// Details returns details about manage offer operation
func (h *manageOfferOpHandler) Details(op RawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageOfferOp := op.Body.MustManageOfferOp()
	manageOfferOpRes := opRes.MustManageOfferResult().MustSuccess()

	offerID := int64(manageOfferOp.OfferId)
	isDeleted := manageOfferOpRes.Offer.Effect == xdr.ManageOfferEffectDeleted
	if !isDeleted {
		offerID = int64(manageOfferOpRes.Offer.MustOffer().OfferId)
	}

	return history2.OperationDetails{
		Type: xdr.OperationTypeManageOffer,
		ManageOffer: &history2.ManageOfferDetails{
			OfferID:     offerID,
			OrderBookID: int64(manageOfferOp.OrderBookId),
			BaseAsset:   manageOfferOpRes.BaseAsset,
			QuoteAsset:  manageOfferOpRes.QuoteAsset,
			Amount:      amount.String(int64(manageOfferOp.Amount)),
			Price:       amount.String(int64(manageOfferOp.Price)),
			IsBuy:       manageOfferOp.IsBuy,
			Fee:         amount.String(int64(manageOfferOp.Fee)),
			IsDeleted:   isDeleted,
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

	return []history2.ParticipantEffect{source}, nil
}

func (h *manageOfferOpHandler) getNewOfferEffect(op xdr.ManageOfferOp,
	res xdr.ManageOfferSuccessResult, source history2.ParticipantEffect,
) []history2.ParticipantEffect {
	offer := offer{
		AccountID:           source.AccountID,
		BaseBalanceID:       h.pubKeyProvider.GetBalanceID(op.BaseBalance),
		BaseBalanceAddress:  op.BaseBalance.AsString(),
		QuoteBalanceID:      h.pubKeyProvider.GetBalanceID(op.QuoteBalance),
		QuoteBalanceAddress: op.QuoteBalance.AsString(),
		BaseAsset:           string(res.BaseAsset),
		QuoteAsset:          string(res.QuoteAsset),
		IsBuy:               op.IsBuy,
	}

	participants, _ := h.getMatchesEffects(res.OffersClaimed, offer)
	if res.Offer.Effect == xdr.ManageOfferEffectDeleted {
		return participants
	}

	// we need to handle amount which was not matched
	newOffer := res.Offer.MustOffer()
	source.Effect = history2.Effect{
		Type: history2.EffectTypeLocked,
		Offer: &history2.OfferEffect{
			BaseBalanceAddress:  offer.BaseBalanceAddress,
			QuoteBalanceAddress: offer.QuoteBalanceAddress,
			BaseAsset:           offer.BaseAsset,
			QuoteAsset:          offer.QuoteAsset,
			BaseAmount:          amount.String(int64(newOffer.BaseAmount)),
			QuoteAmount:         amount.String(int64(newOffer.QuoteAmount)),
			IsBuy:               newOffer.IsBuy,
			Price:               amount.String(int64(newOffer.Price)),
		},
	}
	participants = append(participants, source)
	return participants
}

// getDeletedOffersEffect - creates participant effects for all offer entries with change type `State`
func (h *manageOfferOpHandler) getDeletedOffersEffect(ledgerChanges []xdr.LedgerEntryChange,
) []history2.ParticipantEffect {
	var result []history2.ParticipantEffect
	for _, change := range ledgerChanges {
		if change.Type != xdr.LedgerEntryChangeTypeState {
			continue
		}

		if change.MustState().Data.Type != xdr.LedgerEntryTypeOfferEntry {
			continue
		}

		offer := change.MustState().Data.MustOffer()

		participant := history2.ParticipantEffect{
			AccountID: h.pubKeyProvider.GetAccountID(offer.OwnerId),
			BalanceID: new(int64),
			AssetCode: new(string),
		}
		var unlockedAmount int64
		if offer.IsBuy {
			*participant.BalanceID = h.pubKeyProvider.GetBalanceID(offer.QuoteBalance)
			*participant.AssetCode = string(offer.Quote)
			unlockedAmount = int64(offer.QuoteAmount + offer.Fee)
		} else {
			*participant.BalanceID = h.pubKeyProvider.GetBalanceID(offer.BaseBalance)
			*participant.AssetCode = string(offer.Base)
			unlockedAmount = int64(offer.BaseAmount)
		}

		participant.Effect = history2.Effect{
			Type: history2.EffectTypeUnlocked,
			Unlocked: &history2.UnlockedEffect{
				Amount: amount.String(int64(unlockedAmount)),
			},
		}

		result = append(result, participant)
	}

	return result
}

func (h *manageOfferOpHandler) getStateOffers(ledgerChanges []xdr.LedgerEntryChange) []xdr.OfferEntry {
	var result []xdr.OfferEntry
	for _, change := range ledgerChanges {
		if change.Type != xdr.LedgerEntryChangeTypeState {
			continue
		}

		state := change.MustState()
		if state.Data.Type != xdr.LedgerEntryTypeOfferEntry {
			continue
		}

		result = append(result, state.Data.MustOffer())
	}

	return result
}

type offer struct {
	AccountID           int64
	BaseBalanceID       int64
	BaseBalanceAddress  string
	QuoteBalanceID      int64
	QuoteBalanceAddress string
	BaseAsset           string
	QuoteAsset          string
	IsBuy               bool
}

// getParticipantsEffects - returns participants effects based on the provided matches and total base amount
func (h *manageOfferOpHandler) getMatchesEffects(claimOfferAtoms []xdr.ClaimOfferAtom,
	offer offer) ([]history2.ParticipantEffect, uint64) {

	sourceOfferEffect := history2.OfferEffect{
		BaseBalanceAddress:  offer.BaseBalanceAddress,
		QuoteBalanceAddress: offer.QuoteBalanceAddress,
		BaseAsset:           offer.BaseAsset,
		QuoteAsset:          offer.QuoteAsset,
		IsBuy:               offer.IsBuy,
	}

	var totalBaseAmount uint64 = 0
	result := make([]history2.ParticipantEffect, 0, len(claimOfferAtoms)*4)
	for _, offerAtom := range claimOfferAtoms {
		totalBaseAmount += uint64(offerAtom.BaseAmount)

		counterpartyEffect := history2.Effect{
			Type: history2.EffectTypeMatched,
			Offer: &history2.OfferEffect{
				BaseBalanceAddress:  offerAtom.BaseBalance.AsString(),
				QuoteBalanceAddress: offerAtom.QuoteBalance.AsString(),
				BaseAsset:           offer.BaseAsset,
				QuoteAsset:          offer.QuoteAsset,
				BaseAmount:          amount.String(int64(offerAtom.BaseAmount)),
				QuoteAmount:         amount.String(int64(offerAtom.QuoteAmount)),
				IsBuy:               !offer.IsBuy,
				Price:               amount.String(int64(offerAtom.CurrentPrice)),
				FeePaid: history2.FeePaid{
					CalculatedPercent: amount.String(int64(offerAtom.BFeePaid)),
				},
			},
		}

		baseBalanceID := h.pubKeyProvider.GetBalanceID(offerAtom.BaseBalance)
		result = append(result, history2.ParticipantEffect{
			AccountID: h.pubKeyProvider.GetAccountID(offerAtom.BAccountId),
			BalanceID: &baseBalanceID,
			AssetCode: &offer.BaseAsset,
			Effect:    counterpartyEffect,
		})

		quoteBalanceID := h.pubKeyProvider.GetBalanceID(offerAtom.QuoteBalance)
		result = append(result, history2.ParticipantEffect{
			AccountID: h.pubKeyProvider.GetAccountID(offerAtom.BAccountId),
			BalanceID: &quoteBalanceID,
			AssetCode: &offer.QuoteAsset,
			Effect:    counterpartyEffect,
		})

		sourceOfferEffect.BaseAmount = amount.String(int64(offerAtom.BaseAmount))
		sourceOfferEffect.QuoteAmount = amount.String(int64(offerAtom.QuoteAmount))
		sourceOfferEffect.Price = amount.String(int64(offerAtom.CurrentPrice))
		sourceOfferEffect.FeePaid.CalculatedPercent = amount.String(int64(offerAtom.AFeePaid))
		sourceOfferEffectCopy := sourceOfferEffect
		sourceEffect := history2.Effect{
			Type:  history2.EffectTypeMatched,
			Offer: &sourceOfferEffectCopy,
		}

		result = append(result, history2.ParticipantEffect{
			AccountID: offer.AccountID,
			BalanceID: &offer.BaseBalanceID,
			AssetCode: &offer.BaseAsset,
			Effect:    sourceEffect,
		})

		result = append(result, history2.ParticipantEffect{
			AccountID: offer.AccountID,
			BalanceID: &offer.QuoteBalanceID,
			AssetCode: &offer.QuoteAsset,
			Effect:    sourceEffect,
		})
	}

	return result, totalBaseAmount
}
