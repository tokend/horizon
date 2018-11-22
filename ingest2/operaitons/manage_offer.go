package operaitons

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type manageOfferOpHandler struct {
	pubKeyProvider publicKeyProvider
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

	participants := h.participantEffects(
		baseSource, manageOfferOp.BaseBalance, manageOfferOp.QuoteBalance,
		manageOfferOpRes, int64(manageOfferOp.Amount))

	for _, claimedOffer := range manageOfferOpRes.OffersClaimed {
		participantBase := history2.ParticipantEffect{
			AccountID: h.pubKeyProvider.GetAccountID(claimedOffer.BAccountId),
		}

		participants = append(participants,
			h.participantEffects(
				participantBase, claimedOffer.BaseBalance, claimedOffer.QuoteBalance,
				manageOfferOpRes, int64(manageOfferOp.Amount))...)
	}

	return participants, nil
}

func (h *manageOfferOpHandler) participantEffects(participantBase history2.ParticipantEffect,
	baseBalance, quoteBalance xdr.BalanceId, manageOfferOpRes xdr.ManageOfferSuccessResult, baseAmount int64,
) []history2.ParticipantEffect {
	participantBaseBalanceID := h.pubKeyProvider.GetBalanceID(baseBalance)
	participantQuoteBalanceID := h.pubKeyProvider.GetBalanceID(quoteBalance)

	participantBase.Effect = history2.Effect{
		Type: history2.EffectTypeOffer,
		Offer: &history2.OfferEffect{
			BaseBalanceID:  participantBaseBalanceID,
			QuoteBalanceID: participantQuoteBalanceID,
			BaseAmount:     amount.String(baseAmount),
			BaseAsset:      manageOfferOpRes.BaseAsset,
			QuoteAsset:     manageOfferOpRes.QuoteAsset,
		},
	}

	participantQuote := participantBase
	participantBase.BalanceID = &participantBaseBalanceID
	participantBase.AssetCode = &manageOfferOpRes.BaseAsset
	participantQuote.BalanceID = &participantQuoteBalanceID
	participantQuote.AssetCode = &manageOfferOpRes.QuoteAsset

	return []history2.ParticipantEffect{participantBase, participantQuote}
}
