package operations

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type offerDirection struct {
	BaseAsset  xdr.AssetCode
	QuoteAsset xdr.AssetCode
	IsBuy      bool
}

type offerHelper struct {
	pubKeyProvider publicKeyProvider
}

func (h *offerHelper) getParticipantsEffects(claimOfferAtoms []xdr.ClaimOfferAtom,
	sourceOfferDirection offerDirection, sourceAccountID, baseSourceBalanceID, quoteSourceBalanceID int64,
) ([]history2.ParticipantEffect, uint64) {
	var result []history2.ParticipantEffect
	var totalBaseAmount uint64 = 0

	for _, offerAtom := range claimOfferAtoms {
		totalBaseAmount += uint64(offerAtom.BaseAmount)

		baseBalanceID := h.pubKeyProvider.GetBalanceID(offerAtom.BaseBalance)
		quoteBalanceID := h.pubKeyProvider.GetBalanceID(offerAtom.QuoteBalance)

		counterpartyEffect := history2.Effect{
			Type: history2.EffectTypeMatched,
			Offer: &history2.OfferEffect{
				BaseBalanceID:  baseBalanceID,
				QuoteBalanceID: quoteBalanceID,
				BaseAsset:      sourceOfferDirection.BaseAsset,
				QuoteAsset:     sourceOfferDirection.QuoteAsset,
				BaseAmount:     amount.String(int64(offerAtom.BaseAmount)),
				QuoteAmount:    amount.String(int64(offerAtom.QuoteAmount)),
				IsBuy:          !sourceOfferDirection.IsBuy,
				Price:          amount.String(int64(offerAtom.CurrentPrice)),
				FeePaid: history2.FeePaid{
					CalculatedPercent: amount.String(int64(offerAtom.BFeePaid)),
				},
			},
		}

		baseCounterparty := history2.ParticipantEffect{
			AccountID: h.pubKeyProvider.GetAccountID(offerAtom.BAccountId),
			BalanceID: &baseBalanceID,
			AssetCode: &sourceOfferDirection.BaseAsset,
			Effect:    counterpartyEffect,
		}

		quoteCounterparty := history2.ParticipantEffect{
			AccountID: h.pubKeyProvider.GetAccountID(offerAtom.BAccountId),
			BalanceID: &quoteBalanceID,
			AssetCode: &sourceOfferDirection.QuoteAsset,
			Effect:    counterpartyEffect,
		}

		sourceEffect := history2.Effect{
			Type: history2.EffectTypeMatched,
			Offer: &history2.OfferEffect{
				BaseBalanceID:  baseSourceBalanceID,
				QuoteBalanceID: quoteSourceBalanceID,
				BaseAsset:      sourceOfferDirection.BaseAsset,
				QuoteAsset:     sourceOfferDirection.QuoteAsset,
				BaseAmount:     amount.String(int64(offerAtom.BaseAmount)),
				QuoteAmount:    amount.String(int64(offerAtom.QuoteAmount)),
				IsBuy:          sourceOfferDirection.IsBuy,
				Price:          amount.String(int64(offerAtom.CurrentPrice)),
				FeePaid: history2.FeePaid{
					CalculatedPercent: amount.String(int64(offerAtom.AFeePaid)),
				},
			},
		}

		baseSource := history2.ParticipantEffect{
			AccountID: sourceAccountID,
			BalanceID: &baseSourceBalanceID,
			AssetCode: &sourceOfferDirection.BaseAsset,
			Effect:    sourceEffect,
		}

		quoteSource := history2.ParticipantEffect{
			AccountID: sourceAccountID,
			BalanceID: &quoteSourceBalanceID,
			AssetCode: &sourceOfferDirection.QuoteAsset,
			Effect:    sourceEffect,
		}

		result = append(result, baseCounterparty, quoteCounterparty, baseSource, quoteSource)
	}

	return result, totalBaseAmount
}

func (h *offerHelper) getStateOffers(ledgerChanges []xdr.LedgerEntryChange) []xdr.OfferEntry {
	var result []xdr.OfferEntry

	for _, change := range ledgerChanges {
		if change.Type != xdr.LedgerEntryChangeTypeState {
			continue
		}

		if change.MustState().Data.Type != xdr.LedgerEntryTypeOfferEntry {
			continue
		}

		result = append(result, change.MustState().Data.MustOffer())
	}

	return result
}
