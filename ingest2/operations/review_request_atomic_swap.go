package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	regources "gitlab.com/tokend/regources/v2/generated"
)

type atomicSwapHandler struct {
	effectsProvider
}

//PermanentReject - returns participants of fully rejected request
func (h *atomicSwapHandler) PermanentReject(details requestDetails) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.Participant(details.SourceAccountID)}, nil
}

//Fulfilled - returns slice of effects for participants of the operation
func (h *atomicSwapHandler) Fulfilled(details requestDetails) ([]history2.ParticipantEffect, error) {
	atomicSwapExtendedResult := details.Result.TypeExt.MustASwapExtended()

	participants := []history2.ParticipantEffect{
		h.BalanceEffect(atomicSwapExtendedResult.BidOwnerBaseBalanceId, &history2.Effect{
			Type: history2.EffectTypeChargedFromLocked,
			ChargedFromLocked: &history2.BalanceChangeEffect{
				Amount: regources.Amount(atomicSwapExtendedResult.BaseAmount),
			},
		}),
	}

	participants = append(participants,
		h.BalanceEffect(atomicSwapExtendedResult.PurchaserBaseBalanceId, &history2.Effect{
			Type: history2.EffectTypeFunded,
			Funded: &history2.BalanceChangeEffect{
				Amount: regources.Amount(atomicSwapExtendedResult.BaseAmount),
			},
		}))

	bid := h.getAtomicSwapBid(atomicSwapExtendedResult.BidId, details.Changes)
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
		participants = append(participants,
			h.BalanceEffect(atomicSwapExtendedResult.BidOwnerBaseBalanceId, &history2.Effect{
				Type: history2.EffectTypeUnlocked,
				Unlocked: &history2.BalanceChangeEffect{
					Amount: regources.Amount(bid.Amount),
				},
			}))
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
