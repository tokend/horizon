package operaitons

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type manageOfferOpHandler struct {
	pubKeyProvider        publicKeyProvider
	offerHelper           offerHelper
	ledgerChangesProvider ledgerChangesProvider
	balanceProvider       balanceProvider
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

// ParticipantsEffects can return `matched` and `locked` effects if offer created
// returns `unlocked` effects if offer canceled (deleted by user)
func (h *manageOfferOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	manageOfferOp := opBody.MustManageOfferOp()
	manageOfferOpRes := opRes.MustManageOfferResult().MustSuccess()

	if manageOfferOp.Amount != 0 {
		var result []history2.ParticipantEffect
		if manageOfferOp.OrderBookId != 0 {
			var err error
			result, err = h.getSaleAnteEffects()
			if err != nil {
				return nil, errors.Wrap(err, "failed to get sale ante effects")
			}
		}

		return append(result, h.getNewOfferEffect(manageOfferOp, manageOfferOpRes, source)...), nil
	}

	source, err := h.getDeletedOfferEffect(source)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get source effect")
	}

	return []history2.ParticipantEffect{source}, nil
}

func (h *manageOfferOpHandler) getNewOfferEffect(op xdr.ManageOfferOp,
	res xdr.ManageOfferSuccessResult, source history2.ParticipantEffect,
) []history2.ParticipantEffect {
	baseBalanceID := h.pubKeyProvider.GetBalanceID(op.BaseBalance)
	quoteBalanceID := h.pubKeyProvider.GetBalanceID(op.QuoteBalance)

	participants, _ := h.offerHelper.getParticipantsEffects(res.OffersClaimed,
		offerDirection{
			BaseAsset:  res.BaseAsset,
			QuoteAsset: res.QuoteAsset,
			IsBuy:      op.IsBuy,
		},
		source.AccountID, baseBalanceID, quoteBalanceID)

	source.AssetCode = &res.BaseAsset
	source.BalanceID = &baseBalanceID
	if op.IsBuy {
		source.AssetCode = &res.QuoteAsset
		source.BalanceID = &quoteBalanceID
	}

	if res.Offer.Effect != xdr.ManageOfferEffectDeleted {
		newOffer := res.Offer.MustOffer()

		source.Effect = history2.Effect{
			Type: history2.EffectTypeLocked,
			Offer: &history2.OfferEffect{
				BaseBalanceID:  baseBalanceID,
				QuoteBalanceID: quoteBalanceID,
				BaseAsset:      res.BaseAsset,
				QuoteAsset:     res.QuoteAsset,
				BaseAmount:     amount.String(int64(newOffer.BaseAmount)),
				QuoteAmount:    amount.String(int64(newOffer.QuoteAmount)),
				IsBuy:          newOffer.IsBuy,
				Price:          amount.String(int64(newOffer.Price)),
			},
		}
		participants = append(participants, source)
	}

	return participants
}

func (h *manageOfferOpHandler) getDeletedOfferEffect(source history2.ParticipantEffect) (history2.ParticipantEffect, error) {
	offers := h.offerHelper.getStateOffers()
	if len(offers) != 1 {
		return history2.ParticipantEffect{}, errors.New("unexpected count of state offers")
	}

	offer := offers[0]

	baseBalanceID := h.pubKeyProvider.GetBalanceID(offer.BaseBalance)
	quoteBalanceID := h.pubKeyProvider.GetBalanceID(offer.QuoteBalance)

	source.BalanceID = &baseBalanceID
	source.AssetCode = &offer.Base
	unlockedAmount := offer.BaseAmount
	if offer.IsBuy {
		source.BalanceID = &quoteBalanceID
		source.AssetCode = &offer.Quote
		unlockedAmount = offer.QuoteAmount + offer.Fee
	}

	source.Effect = history2.Effect{
		Type: history2.EffectTypeUnlocked,
		Unlocked: &history2.UnlockedEffect{
			Amount: amount.String(int64(unlockedAmount)),
		},
	}

	return source, nil
}

func (h *manageOfferOpHandler) getSaleAnteEffects() ([]history2.ParticipantEffect, error) {
	var result []history2.ParticipantEffect

	createdSaleAntes, updatedSaleAntes, statedSaleAntes := h.getChangedSaleAntes()
	for _, saleAnte := range createdSaleAntes {
		balance := h.balanceProvider.GetBalanceByID(saleAnte.ParticipantBalanceId)

		result = append(result, history2.ParticipantEffect{
			AccountID: balance.AccountID,
			BalanceID: &balance.ID,
			AssetCode: &balance.AssetCode,
			Effect: history2.Effect{
				Type: history2.EffectTypeLocked,
				Locked: &history2.LockedEffect{
					Amount: amount.StringU(uint64(saleAnte.Amount)),
				},
			},
		})
	}

	// for now possible only one updated saleAnte
	for _, updatedSaleAnte := range updatedSaleAntes {
		if len(statedSaleAntes) == 0 {
			return nil, errors.New("unexpected state, updated sale ante exists, but stated not")
		}

		for _, statedSaleAnte := range statedSaleAntes {
			if (updatedSaleAnte.ParticipantBalanceId.AsString() != statedSaleAnte.ParticipantBalanceId.AsString()) ||
				(updatedSaleAnte.SaleId != statedSaleAnte.SaleId) {
				continue
			}

			balance := h.balanceProvider.GetBalanceByID(updatedSaleAnte.ParticipantBalanceId)

			result = append(result, history2.ParticipantEffect{
				AccountID: balance.AccountID,
				BalanceID: &balance.ID,
				AssetCode: &balance.AssetCode,
				Effect: history2.Effect{
					Type: history2.EffectTypeLocked,
					Locked: &history2.LockedEffect{
						Amount: amount.StringU(uint64(updatedSaleAnte.Amount - statedSaleAnte.Amount)),
					},
				},
			})
		}
	}

	return result, nil
}

func (h *manageOfferOpHandler) getChangedSaleAntes() ([]xdr.SaleAnteEntry, []xdr.SaleAnteEntry, []xdr.SaleAnteEntry) {
	var created, updated, stated []xdr.SaleAnteEntry

	ledgerChanges := h.ledgerChangesProvider.GetLedgerChanges()
	for _, change := range ledgerChanges {
		switch change.Type {
		case xdr.LedgerEntryChangeTypeCreated:
			if change.MustCreated().Data.Type != xdr.LedgerEntryTypeSaleAnte {
				continue
			}
			created = append(created, change.MustCreated().Data.MustSaleAnte())
		case xdr.LedgerEntryChangeTypeUpdated:
			if change.MustCreated().Data.Type != xdr.LedgerEntryTypeSaleAnte {
				continue
			}
			updated = append(updated, change.MustUpdated().Data.MustSaleAnte())
		case xdr.LedgerEntryChangeTypeState:
			if change.MustCreated().Data.Type != xdr.LedgerEntryTypeSaleAnte {
				continue
			}
			stated = append(stated, change.MustState().Data.MustSaleAnte())
		}
	}

	return created, updated, stated
}
