package operaitons

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type manageOfferOpHandler struct {
	pubKeyProvider publicKeyProvider
	offerHelper    offerHelper
}

func (h *manageOfferOpHandler) OperationDetails(op rawOperation, opRes xdr.OperationResultTr,
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

func (h *manageOfferOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	manageOfferOp := opBody.MustManageOfferOp()
	manageOfferOpRes := opRes.MustManageOfferResult().MustSuccess()

	baseBalanceID := h.pubKeyProvider.GetBalanceID(manageOfferOp.BaseBalance)
	quoteBalanceID := h.pubKeyProvider.GetBalanceID(manageOfferOp.QuoteBalance)

	participants, _ := h.offerHelper.getParticipantsEffects(manageOfferOpRes.OffersClaimed,
		offerDirection{
			BaseAsset:  manageOfferOpRes.BaseAsset,
			QuoteAsset: manageOfferOpRes.QuoteAsset,
			IsBuy:      manageOfferOp.IsBuy,
		},
		source.AccountID, baseBalanceID, quoteBalanceID)

	source.AssetCode = &manageOfferOpRes.BaseAsset
	source.BalanceID = &baseBalanceID
	if manageOfferOp.IsBuy {
		source.AssetCode = &manageOfferOpRes.QuoteAsset
		source.BalanceID = &quoteBalanceID
	}

	if manageOfferOpRes.Offer.Effect != xdr.ManageOfferEffectDeleted {
		newOffer := manageOfferOpRes.Offer.MustOffer()

		source.Effect = history2.Effect{
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
		participants = append(participants, source)
	}

	return participants, nil
}
