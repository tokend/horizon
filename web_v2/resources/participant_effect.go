package resources

import (
	"strconv"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

//NewParticipantEffectKey - returns new key for ID
func NewParticipantEffectKey(id int64) regources.Key {
	return regources.Key{
		Type: regources.TypeParticipantEffects,
		ID:   strconv.FormatInt(id, 10),
	}
}

//NewEffect - returns new instance of effect
func NewEffect(id int64, effect history2.Effect) regources.Resource {
	switch effect.Type {
	case history2.EffectTypeFunded:
		return newBalanceChangeEffect(id, regources.TypeEffectsFunded, *effect.Funded)
	case history2.EffectTypeIssued:
		return newBalanceChangeEffect(id, regources.TypeEffectsIssued, *effect.Issued)
	case history2.EffectTypeCharged:
		return newBalanceChangeEffect(id, regources.TypeEffectsCharged, *effect.Charged)
	case history2.EffectTypeWithdrawn:
		return newBalanceChangeEffect(id, regources.TypeEffectsWithdrawn, *effect.Withdrawn)
	case history2.EffectTypeLocked:
		return newBalanceChangeEffect(id, regources.TypeEffectsLocked, *effect.Locked)
	case history2.EffectTypeUnlocked:
		return newBalanceChangeEffect(id, regources.TypeEffectsUnlocked, *effect.Unlocked)
	case history2.EffectTypeChargedFromLocked:
		return newBalanceChangeEffect(id, regources.TypeEffectsChargedFromLocked, *effect.ChargedFromLocked)
	case history2.EffectTypeMatched:
		return newMatchedEffect(id, *effect.Matched)
	default:
		panic(errors.From(errors.New("unexpected effect type"), logan.F{
			"type": effect.Type,
		}))
	}
}

func newMatchedEffect(id int64, effect history2.MatchEffect) *regources.EffectMatched {
	return &regources.EffectMatched{
		Key: regources.Key{
			Type: regources.TypeEffectsMatched,
			ID:   strconv.FormatInt(id, 10),
		},
		Attributes: regources.EffectMatchedAttrs{
			OrderBookID: effect.OrderBookID,
			OfferID:     effect.OfferID,
			Price:       effect.Price,
			Charged:     newParticularBalanceChange(effect.Charged),
			Funded:      newParticularBalanceChange(effect.Funded),
		},
	}
}

func newParticularBalanceChange(effect history2.ParticularBalanceChangeEffect) regources.ParticularBalanceChangeEffect {
	return regources.ParticularBalanceChangeEffect{
		EffectBalanceChangeAttrs: regources.EffectBalanceChangeAttrs{
			Amount: effect.Amount,
			Fee:    effect.Fee,
		},
		BalanceAddress: effect.BalanceAddress,
		AssetCode:      effect.AssetCode,
	}
}

func newBalanceChangeEffect(id int64, resourceType regources.ResourceType,
	effect history2.BalanceChangeEffect) *regources.EffectBalanceChange {
	return &regources.EffectBalanceChange{
		Key: regources.Key{
			Type: resourceType,
			ID:   strconv.FormatInt(id, 10),
		},
		Attributes: regources.EffectBalanceChangeAttrs{
			Amount: effect.Amount,
			Fee:    effect.Fee,
		},
	}
}
