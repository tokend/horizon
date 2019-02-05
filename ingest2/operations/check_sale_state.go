package operations

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

type checkSaleStateOpHandler struct {
	manageOfferOpHandler *manageOfferOpHandler
}

// CreatorDetails returns details about check sale state operation
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
	opRes xdr.OperationResultTr, source history2.ParticipantEffect, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	res := opRes.MustCheckSaleStateResult().MustSuccess()

	switch res.Effect.Effect {
	case xdr.CheckSaleStateEffectCanceled, xdr.CheckSaleStateEffectUpdated:
		return h.manageOfferOpHandler.getDeletedOffersEffect(ledgerChanges), nil
	case xdr.CheckSaleStateEffectClosed:
		return h.getApprovedParticipants(int64(opBody.MustCheckSaleStateOp().SaleId), res.Effect.MustSaleClosed()), nil
	default:
		return nil, errors.From(errors.New("unexpected check sale state result effect"), map[string]interface{}{
			"effect_i": int32(res.Effect.Effect),
			"sale_id":  uint64(res.SaleId),
		})
	}
}

func (h *checkSaleStateOpHandler) getApprovedParticipants(orderBookID int64, closedRes xdr.CheckSaleClosedResult,
) []history2.ParticipantEffect {
	// TODO: we are not handling here cases that some parts of offers might be canceled due to rounding
	if len(closedRes.Results) == 0 {
		return nil
	}

	result := make([]history2.ParticipantEffect, 0)
	var totalBaseIssued uint64
	ownerID := h.manageOfferOpHandler.pubKeyProvider.MustAccountID(closedRes.SaleOwner)
	// it does not matter which base balance is used as we are sure that the operation of distribution will be clean
	baseBalanceAddress := closedRes.Results[0].SaleBaseBalance.AsString()
	baseBalanceID := h.manageOfferOpHandler.pubKeyProvider.MustBalanceID(closedRes.Results[0].SaleBaseBalance)
	baseAsset := string(closedRes.Results[0].SaleDetails.BaseAsset)
	for _, assetPairResult := range closedRes.Results {
		sourceOffer := offer{
			OrderBookID:         orderBookID,
			AccountID:           ownerID,
			BaseBalanceID:       baseBalanceID,
			BaseBalanceAddress:  baseBalanceAddress,
			QuoteBalanceID:      h.manageOfferOpHandler.pubKeyProvider.MustBalanceID(assetPairResult.SaleQuoteBalance),
			QuoteBalanceAddress: assetPairResult.SaleQuoteBalance.AsString(),
			BaseAsset:           baseAsset,
			QuoteAsset:          string(assetPairResult.SaleDetails.QuoteAsset),
			IsBuy:               false,
		}
		assetPairMatches, baseIssued := h.manageOfferOpHandler.getMatchesEffects(
			assetPairResult.SaleDetails.OffersClaimed, sourceOffer)

		totalBaseIssued += baseIssued
		result = append(result, assetPairMatches...)
	}

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
