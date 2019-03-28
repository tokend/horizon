package resources

import (
	"strconv"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/rgenerated"
)

//NewParticipantEffectKey - returns new key for ID
func NewParticipantEffectKey(id int64) rgenerated.Key {
	return rgenerated.Key{
		Type: rgenerated.PARTICIPANT_EFFECTS,
		ID:   strconv.FormatInt(id, 10),
	}
}

//NewEffect - returns new instance of effect
func NewEffect(id int64, effect history2.Effect) rgenerated.Resource {
	switch effect.Type {
	case history2.EffectTypeFunded:
		return newBalanceChangeEffect(id, rgenerated.EFFECTS_FUNDED, *effect.Funded)
	case history2.EffectTypeIssued:
		return newBalanceChangeEffect(id, rgenerated.EFFECTS_ISSUED, *effect.Issued)
	case history2.EffectTypeCharged:
		return newBalanceChangeEffect(id, rgenerated.EFFECTS_CHARGED, *effect.Charged)
	case history2.EffectTypeWithdrawn:
		return newBalanceChangeEffect(id, rgenerated.EFFECTS_WITHDRAWN, *effect.Withdrawn)
	case history2.EffectTypeLocked:
		return newBalanceChangeEffect(id, rgenerated.EFFECTS_LOCKED, *effect.Locked)
	case history2.EffectTypeUnlocked:
		return newBalanceChangeEffect(id, rgenerated.EFFECTS_UNLOCKED, *effect.Unlocked)
	case history2.EffectTypeChargedFromLocked:
		return newBalanceChangeEffect(id, rgenerated.EFFECTS_CHARGED_FROM_LOCKED, *effect.ChargedFromLocked)
	case history2.EffectTypeMatched:
		return newMatchedEffect(id, *effect.Matched)
	default:
		panic(errors.From(errors.New("unexpected effect type"), logan.F{
			"type": effect.Type,
		}))
	}
}

func newMatchedEffect(id int64, effect history2.MatchEffect) *rgenerated.EffectMatched {
	return &rgenerated.EffectMatched{
		Key: rgenerated.Key{
			Type: rgenerated.EFFECTS_MATCHED,
			ID:   strconv.FormatInt(id, 10),
		},
		Attributes: rgenerated.EffectMatchedAttributes{
			OrderBookId: effect.OrderBookID,
			OfferId:     effect.OfferID,
			Price:       effect.Price,
			Charged:     newParticularBalanceChange(effect.Charged),
			Funded:      newParticularBalanceChange(effect.Funded),
		},
	}
}

func newParticularBalanceChange(effect history2.ParticularBalanceChangeEffect) rgenerated.ParticularBalanceChangeEffect {
	return rgenerated.ParticularBalanceChangeEffect{
		Amount:         effect.Amount,
		Fee:            effect.Fee,
		BalanceAddress: effect.BalanceAddress,
		AssetCode:      effect.AssetCode,
	}
}

func newBalanceChangeEffect(id int64, resourceType rgenerated.ResourceType,
	effect history2.BalanceChangeEffect) *rgenerated.EffectBalanceChange {
	return &rgenerated.EffectBalanceChange{
		Key: rgenerated.Key{
			Type: resourceType,
			ID:   strconv.FormatInt(id, 10),
		},
		Attributes: rgenerated.EffectBalanceChangeAttributes{
			Amount: effect.Amount,
			Fee:    effect.Fee,
		},
	}
}
