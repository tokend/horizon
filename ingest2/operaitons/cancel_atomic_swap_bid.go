package operaitons

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type cancelAtomicSwapBidOpHandler struct {
	pubKeyProvider        publicKeyProvider
	ledgerChangesProvider ledgerChangesProvider
}

func (h *cancelAtomicSwapBidOpHandler) OperationDetails(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	return history2.OperationDetails{
		Type: xdr.OperationTypeCancelAswapBid,
		CancelAtomicSwapBid: &history2.CancelAtomicSwapBidDetails{
			BidID: int64(op.Body.MustCancelASwapBidOp().BidId),
		},
	}, nil
}

// ParticipantsEffects returns participants effects with source effect `unlocked`
// if atomic swap bid has zero locked amount
func (h *cancelAtomicSwapBidOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	atomicSwapBid := h.getAtomicSwapBid(opBody.MustCancelASwapBidOp().BidId)

	if atomicSwapBid == nil {
		return nil, nil
	}

	if atomicSwapBid.LockedAmount != 0 {
		return []history2.ParticipantEffect{source}, nil
	}

	balanceID := h.pubKeyProvider.GetBalanceID(atomicSwapBid.BaseBalance)

	source.BalanceID = &balanceID
	source.AssetCode = &atomicSwapBid.BaseAsset
	source.Effect = history2.Effect{
		Type: history2.EffectTypeUnlocked,
		Unlocked: &history2.UnlockedEffect{
			Amount: amount.StringU(uint64(atomicSwapBid.Amount)),
		},
	}

	return []history2.ParticipantEffect{source}, nil
}

func (h *cancelAtomicSwapBidOpHandler) getAtomicSwapBid(bidID xdr.Uint64) *xdr.AtomicSwapBidEntry {
	ledgerChanges := h.ledgerChangesProvider.GetLedgerChanges()

	for _, change := range ledgerChanges {
		if change.Type != xdr.LedgerEntryChangeTypeState {
			continue
		}

		if change.MustState().Data.Type != xdr.LedgerEntryTypeAtomicSwapBid {
			continue
		}

		atomicSwapBid := change.MustState().Data.MustAtomicSwapBid()

		if atomicSwapBid.BidId == bidID {
			return &atomicSwapBid
		}
	}

	return nil
}
