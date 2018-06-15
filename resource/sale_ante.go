package resource

import (
	"gitlab.com/swarmfund/horizon/db2/core"
	"strconv"
	"gitlab.com/tokend/go/amount"
)

type SaleAnte struct {
	PT                   string `json:"paging_token"`
	SaleID               string `json:"sale_id"`
	ParticipantBalanceID string `json:"participant_balance_id"`
	Amount               string `json:"amount"`
}

func (s *SaleAnte) Populate(raw *core.SaleAnte) {
	s.SaleID = strconv.FormatUint(raw.SaleID, 10)
	s.ParticipantBalanceID = raw.ParticipantBalanceID
	s.Amount = amount.StringU(raw.Amount)
}

func (s *SaleAnte) PagingToken() string {
	return s.PT
}
