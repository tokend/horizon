package participants

import (
	"bullioncoin.githost.io/development/go/xdr"
)

type MatchesDetails struct {
	Source     xdr.AccountId
	Balance    xdr.BalanceId
	BaseAsset  string
	QuoteAsset string
	IsBuy      bool
	Matches    []Match
}

func NewMatchesDetails(source xdr.AccountId, balance xdr.BalanceId, baseAsset, quoteAsset string, isBuy bool) *MatchesDetails {
	return &MatchesDetails{
		Source:     source,
		Balance:    balance,
		BaseAsset:  baseAsset,
		QuoteAsset: quoteAsset,
		IsBuy:      isBuy,
		Matches:    []Match{},
	}
}

func (m *MatchesDetails) ToParticipant() Participant {
	return Participant{m.Source, &m.Balance, &map[string]interface{}{
		"base_asset":  m.BaseAsset,
		"quote_asset": m.QuoteAsset,
		"is_buy":      m.IsBuy,
		"matches":     m.Matches,
	}}
}

func (m *MatchesDetails) Add(match *Match) {
	for i := range m.Matches {
		if m.Matches[i].CanAdd(match) {
			m.Matches[i].Add(match)
			return
		}
	}

	m.Matches = append(m.Matches, *match)
}
