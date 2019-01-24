package operations

import (
	history "gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

func populateEffects(balance history.Balance, effect regources.Effect,
	source history.ParticipantEffect) []history.ParticipantEffect {

	if balance.AccountID == source.AccountID {
		source.BalanceID = &balance.ID
		source.AssetCode = &balance.AssetCode
		source.Effect = &effect
		return []history.ParticipantEffect{source}
	}

	return []history.ParticipantEffect{{
		AccountID: balance.AccountID,
		BalanceID: &balance.ID,
		AssetCode: &balance.AssetCode,
		Effect:    &effect,
	}, source}
}
