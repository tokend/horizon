package operations

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type atomicSwapHandler struct {
	pubKeyProvider IDProvider
}

//ParticipantsEffects - returns slice of effects for participants of the operation
func (h *atomicSwapHandler) ParticipantsEffects(op xdr.ReviewRequestOp,
	res xdr.ReviewRequestSuccessResult, request xdr.ReviewableRequestEntry,
	source history2.ParticipantEffect, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	atomicSwapExtendedResult := res.TypeExt.MustASwapExtended()

	ownerBalanceID := h.pubKeyProvider.MustBalanceID(atomicSwapExtendedResult.BidOwnerBaseBalanceId)

	baseAsset := string(atomicSwapExtendedResult.BaseAsset)
	participants := []history2.ParticipantEffect{{
		AccountID: h.pubKeyProvider.MustAccountID(atomicSwapExtendedResult.BidOwnerId),
		BalanceID: &ownerBalanceID,
		AssetCode: &baseAsset,
		Effect: history2.Effect{
			Type: history2.EffectTypeChargedFromLocked,
			ChargedFromLocked: &history2.BalanceChangeEffect{
				Amount: amount.StringU(uint64(atomicSwapExtendedResult.BaseAmount)),
			},
		},
	}}

	purchaserBaseBalanceID := h.pubKeyProvider.MustBalanceID(atomicSwapExtendedResult.PurchaserBaseBalanceId)

	participants = append(participants, history2.ParticipantEffect{
		AccountID: h.pubKeyProvider.MustAccountID(atomicSwapExtendedResult.PurchaserId),
		BalanceID: &purchaserBaseBalanceID,
		AssetCode: &baseAsset,
		Effect: history2.Effect{
			Type: history2.EffectTypeFunded,
			Funded: &history2.BalanceChangeEffect{
				Amount: amount.StringU(uint64(atomicSwapExtendedResult.BaseAmount)),
			},
		},
	})

	bid := h.getAtomicSwapBid(atomicSwapExtendedResult.BidId, ledgerChanges)
	// review of bid request has not affected bid, so there is no additional effects
	if bid == nil {
		return participants, nil
	}

	// no additional effects for the bid owner
	if bid.Amount == 0 {
		return participants, nil
	}

	bidIsSoldOut := (bid.Amount == 0) && (bid.LockedAmount == 0)
	bindIsRemovedOnReviewAfterCancel := bid.IsCancelled && bid.LockedAmount == 0
	bidIsRemoved := bidIsSoldOut || bindIsRemovedOnReviewAfterCancel
	// If bid was removed, but we had to unlock some amount
	if bidIsRemoved {
		participants = append(participants, history2.ParticipantEffect{
			AccountID: h.pubKeyProvider.MustAccountID(atomicSwapExtendedResult.BidOwnerId),
			BalanceID: &ownerBalanceID,
			AssetCode: &baseAsset,
			Effect: history2.Effect{
				Type: history2.EffectTypeUnlocked,
				Unlocked: &history2.BalanceChangeEffect{
					Amount: amount.StringU(uint64(bid.Amount)),
				},
			},
		})
	}

	return participants, nil
}

func (h *atomicSwapHandler) getAtomicSwapBid(bidID xdr.Uint64, ledgerChanges []xdr.LedgerEntryChange,
) *xdr.AtomicSwapBidEntry {
	for _, change := range ledgerChanges {
		if change.Type != xdr.LedgerEntryChangeTypeUpdated {
			continue
		}

		if change.MustUpdated().Data.Type != xdr.LedgerEntryTypeAtomicSwapBid {
			continue
		}

		atomicSwapBid := change.MustUpdated().Data.MustAtomicSwapBid()

		if atomicSwapBid.BidId == bidID {
			return &atomicSwapBid
		}
	}

	return nil
}
