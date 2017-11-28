package participants

import "gitlab.com/tokend/go/xdr"

type MatchesDetailsByBalance struct {
	Details map[string]*MatchesDetails
}

func (m *MatchesDetailsByBalance) ToParticipants(participants []Participant) []Participant {
	for _, detail := range m.Details {
		participants = append(participants, detail.ToParticipant())
	}

	return participants
}

func NewMatchesDetailsByBalance() *MatchesDetailsByBalance {
	return &MatchesDetailsByBalance{
		Details: make(map[string]*MatchesDetails),
	}
}

func (m *MatchesDetailsByBalance) Add(source xdr.AccountId, balance xdr.BalanceId, baseAsset, quoteAsset xdr.AssetCode, isBuy bool, match *Match) {
	matches, ok := m.Details[balance.AsString()]
	if !ok {
		matches = NewMatchesDetails(source, balance, string(baseAsset), string(quoteAsset), isBuy)
		m.Details[balance.AsString()] = matches
	}

	matches.Add(match)
}
