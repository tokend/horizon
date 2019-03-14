package operations

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

type cancelAtomicSwapBidOpHandler struct {
	effectsProvider
}

// Details returns details about cancel atomic swap bid operation
func (h *cancelAtomicSwapBidOpHandler) Details(op rawOperation,
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
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	opResult := opRes.MustCancelASwapBidResult().MustSuccess()

	// it means that there is pending atomic swap request,
	// so bid still exists
	// we must wait for review that request
	if opResult.LockedAmount != 0 {
		return h.effectsProvider.ParticipantsEffects(opBody, opRes, sourceAccountID, ledgerChanges)
	}

	atomicSwapBid := h.getAtomicSwapBid(opBody.MustCancelASwapBidOp().BidId, ledgerChanges)

	if atomicSwapBid == nil {
		return nil, errors.From(
			errors.New("expected atomic swap to be in STATE ledger changes"), map[string]interface{}{
				"bid_id": uint64(opBody.MustCancelASwapBidOp().BidId),
			})
	}

	return []history2.ParticipantEffect{h.BalanceEffect(atomicSwapBid.BaseBalance, &history2.Effect{
		Type: history2.EffectTypeUnlocked,
		Unlocked: &history2.BalanceChangeEffect{
			Amount: regources.Amount(atomicSwapBid.Amount),
		},
	})}, nil
}

func (h *cancelAtomicSwapBidOpHandler) getAtomicSwapBid(bidID xdr.Uint64,
	ledgerChanges []xdr.LedgerEntryChange,
) *xdr.AtomicSwapBidEntry {
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
