package operations

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

type cancelAtomicSwapBidOpHandler struct {
	pubKeyProvider IDProvider
}

// Details returns details about cancel atomic swap bid operation
func (h *cancelAtomicSwapBidOpHandler) Details(op rawOperation,
	opRes xdr.OperationResultTr,
) (regources.OperationDetails, error) {
	return regources.OperationDetails{
		Type: xdr.OperationTypeCancelAswapBid,
		CancelAtomicSwapBid: &regources.CancelAtomicSwapBidDetails{
			BidID: int64(op.Body.MustCancelASwapBidOp().BidId),
		},
	}, nil
}

// ParticipantsEffects returns participants effects with source effect `unlocked`
// if atomic swap bid has zero locked amount
func (h *cancelAtomicSwapBidOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	atomicSwapBid := h.getAtomicSwapBid(opBody.MustCancelASwapBidOp().BidId, ledgerChanges)

	if atomicSwapBid == nil {
		return nil, errors.From(
			errors.New("expected atomic swap to be in STATE ledger changes"), map[string]interface{}{
				"bid_id": uint64(opBody.MustCancelASwapBidOp().BidId),
			})
	}

	// it means that there is pending atomic swap request,
	// so bid still exists
	// we must wait for review that request
	if atomicSwapBid.LockedAmount != 0 {
		return []history2.ParticipantEffect{source}, nil
	}

	balanceID := h.pubKeyProvider.MustBalanceID(atomicSwapBid.BaseBalance)

	source.BalanceID = &balanceID
	atomicSwapBidBaseAsset := string(atomicSwapBid.BaseAsset)
	source.AssetCode = &atomicSwapBidBaseAsset
	source.Effect = &regources.Effect{
		Type: regources.EffectTypeUnlocked,
		Unlocked: &regources.BalanceChangeEffect{
			Amount: regources.Amount(atomicSwapBid.Amount),
		},
	}

	return []history2.ParticipantEffect{source}, nil
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
