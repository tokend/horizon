package operaitons

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type manageOfferOpHandler struct {
	pubKeyProvider        publicKeyProvider
	ledgerChangesProvider ledgerChangesProvider
}

func (h *manageOfferOpHandler) OperationDetails(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageOfferOp := opBody.MustManageOfferOp()
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

func (h *manageOfferOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, baseSource history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	manageOfferOp := opBody.MustManageOfferOp()
	manageOfferOpRes := opRes.MustManageOfferResult().MustSuccess()

	baseBalanceID := h.pubKeyProvider.GetBalanceID(manageOfferOp.BaseBalance)
	quoteBalanceID := h.pubKeyProvider.GetBalanceID(manageOfferOp.QuoteBalance)

	affectedOffers := h.getUpdatedOffers()
	affectedOffers = append(affectedOffers, h.getMatchedOffers(manageOfferOp, manageOfferOpRes)...)

	participants := h.getParticipantsEffects(affectedOffers, baseSource, manageOfferOp, baseBalanceID, quoteBalanceID)

	baseSource.AssetCode = &manageOfferOpRes.BaseAsset
	baseSource.BalanceID = &baseBalanceID
	if manageOfferOp.IsBuy {
		baseSource.AssetCode = &manageOfferOpRes.QuoteAsset
		baseSource.BalanceID = &quoteBalanceID
	}

	if manageOfferOpRes.Offer.Effect != xdr.ManageOfferEffectDeleted {
		newOffer := manageOfferOpRes.Offer.MustOffer()

		baseSource.Effect = history2.Effect{
			Type: history2.EffectTypeLocked,
			Offer: &history2.OfferEffect{
				BaseBalanceID:  baseBalanceID,
				QuoteBalanceID: quoteBalanceID,
				BaseAsset:      manageOfferOpRes.BaseAsset,
				QuoteAsset:     manageOfferOpRes.QuoteAsset,
				BaseAmount:     amount.String(int64(newOffer.BaseAmount)),
				QuoteAmount:    amount.String(int64(newOffer.QuoteAmount)),
				IsBuy:          newOffer.IsBuy,
				Price:          amount.String(int64(newOffer.Price)),
			},
		}
		participants = append(participants, baseSource)
	}

	return participants, nil
}

func (h *manageOfferOpHandler) getUpdatedOffers() []xdr.OfferEntry {
	ledgerChanges := h.ledgerChangesProvider.GetLedgerChanges()

	var updatedOffers []xdr.OfferEntry

	for _, change := range ledgerChanges {
		if change.Type != xdr.LedgerEntryChangeTypeUpdated {
			continue
		}

		if change.MustUpdated().Data.Type == xdr.LedgerEntryTypeOfferEntry {
			updatedOffers = append(updatedOffers, change.MustUpdated().Data.MustOffer())
		}
	}

	return updatedOffers
}

func (h *manageOfferOpHandler) getMatchedOffers(manageOfferOp xdr.ManageOfferOp,
	manageOfferOpRes xdr.ManageOfferSuccessResult,
) []xdr.OfferEntry {
	var result []xdr.OfferEntry

	for _, offerAtom := range manageOfferOpRes.OffersClaimed {
		result = append(result, xdr.OfferEntry{
			OfferId:      offerAtom.OfferId,
			OrderBookId:  0,
			OwnerId:      offerAtom.BAccountId,
			IsBuy:        !manageOfferOp.IsBuy,
			Base:         manageOfferOpRes.BaseAsset,
			Quote:        manageOfferOpRes.QuoteAsset,
			BaseBalance:  offerAtom.BaseBalance,
			QuoteBalance: offerAtom.QuoteBalance,
			BaseAmount:   offerAtom.BaseAmount,
			QuoteAmount:  offerAtom.QuoteAmount,
			Price:        offerAtom.CurrentPrice,
		})
	}

	return result
}

func (h *manageOfferOpHandler) getParticipantsEffects(affectedOffers []xdr.OfferEntry,
	source history2.ParticipantEffect, manageOfferOp xdr.ManageOfferOp, baseSourceBalanceID, quoteSourceBalanceID int64,
) []history2.ParticipantEffect {
	var result []history2.ParticipantEffect

	for _, offer := range affectedOffers {
		baseBalanceID := h.pubKeyProvider.GetBalanceID(offer.BaseBalance)
		quoteBalanceID := h.pubKeyProvider.GetBalanceID(offer.QuoteBalance)

		counterpartyEffect := history2.Effect{
			Type: history2.EffectTypeMatched,
			Offer: &history2.OfferEffect{
				BaseBalanceID:  baseBalanceID,
				QuoteBalanceID: quoteBalanceID,
				BaseAsset:      offer.Base,
				QuoteAsset:     offer.Quote,
				BaseAmount:     amount.String(int64(offer.BaseAmount)),
				QuoteAmount:    amount.String(int64(offer.QuoteAmount)),
				IsBuy:          offer.IsBuy,
				Price:          amount.String(int64(offer.Price)),
			},
		}

		baseCounterparty := history2.ParticipantEffect{
			AccountID: h.pubKeyProvider.GetAccountID(offer.OwnerId),
			BalanceID: &baseBalanceID,
			AssetCode: &offer.Base,
			Effect:    counterpartyEffect,
		}

		quoteCounterparty := history2.ParticipantEffect{
			AccountID: h.pubKeyProvider.GetAccountID(offer.OwnerId),
			BalanceID: &quoteBalanceID,
			AssetCode: &offer.Quote,
			Effect:    counterpartyEffect,
		}

		sourceEffect := history2.Effect{
			Type: history2.EffectTypeMatched,
			Offer: &history2.OfferEffect{
				BaseBalanceID:  baseSourceBalanceID,
				QuoteBalanceID: quoteSourceBalanceID,
				BaseAsset:      offer.Base,
				QuoteAsset:     offer.Quote,
				BaseAmount:     amount.String(int64(offer.BaseAmount)),
				QuoteAmount:    amount.String(int64(offer.QuoteAmount)),
				IsBuy:          offer.IsBuy,
				Price:          amount.String(int64(offer.Price)),
			},
		}

		baseSource := history2.ParticipantEffect{
			AccountID: source.AccountID,
			BalanceID: &baseSourceBalanceID,
			AssetCode: &offer.Base,
			Effect:    sourceEffect,
		}

		quoteSource := history2.ParticipantEffect{
			AccountID: source.AccountID,
			BalanceID: &quoteSourceBalanceID,
			Effect:    sourceEffect,
		}

		result = append(result, baseCounterparty, quoteCounterparty, baseSource, quoteSource)
	}

	return result
}
