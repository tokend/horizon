package reviewrequest

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type atomicSwapHandler struct {
	pubKeyProvider publicKeyProvider
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

func (h *atomicSwapHandler) specificParticipantsEffects(op xdr.ReviewRequestOp,
	res xdr.ReviewRequestSuccessResult, request xdr.ReviewableRequestEntry,
	source history2.ParticipantEffect, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	atomicSwapExtendedResult := res.TypeExt.MustASwapExtended()

	ownerBalanceID := h.pubKeyProvider.GetBalanceID(atomicSwapExtendedResult.BidOwnerBaseBalanceId)

	participants := []history2.ParticipantEffect{{
		AccountID: h.pubKeyProvider.GetAccountID(atomicSwapExtendedResult.BidOwnerId),
		BalanceID: &ownerBalanceID,
		AssetCode: &atomicSwapExtendedResult.BaseAsset,
		Effect: history2.Effect{
			Type: history2.EffectTypeChargedFromLocked,
			ChargedFromLocked: &history2.ChargedFromLockedEffect{
				Amount: amount.StringU(uint64(atomicSwapExtendedResult.BaseAmount)),
			},
		},
	}}

	purchaserBaseBalanceID := h.pubKeyProvider.GetBalanceID(atomicSwapExtendedResult.PurchaserBaseBalanceId)

	participants = append(participants, history2.ParticipantEffect{
		AccountID: h.pubKeyProvider.GetAccountID(atomicSwapExtendedResult.PurchaserId),
		BalanceID: &purchaserBaseBalanceID,
		AssetCode: &atomicSwapExtendedResult.BaseAsset,
		Effect: history2.Effect{
			Type: history2.EffectTypeFunded,
			Funded: &history2.FundedEffect{
				Amount: amount.StringU(uint64(atomicSwapExtendedResult.BaseAmount)),
			},
		},
	})

	bid := h.getAtomicSwapBid(atomicSwapExtendedResult.BidId, ledgerChanges)
	if bid == nil {
		return participants, nil
	}

	bidIsSoldOut := (bid.Amount == 0) && (bid.LockedAmount == 0)
	bidIsRemoved := bidIsSoldOut || (bid.IsCancelled && bid.LockedAmount == 0)
	if bidIsRemoved && (bid.Amount != 0) {
		participants = append(participants, history2.ParticipantEffect{
			AccountID: h.pubKeyProvider.GetAccountID(atomicSwapExtendedResult.BidOwnerId),
			BalanceID: &ownerBalanceID,
			AssetCode: &atomicSwapExtendedResult.BaseAsset,
			Effect: history2.Effect{
				Type: history2.EffectTypeUnlocked,
				Unlocked: &history2.UnlockedEffect{
					Amount: amount.StringU(uint64(bid.Amount)),
				},
			},
		})
	}

	return participants, nil
}
