package resource

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"

	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/resource/base"
)

type Sale struct {
	PT         string                 `json:"paging_token"`
	ID         string                 `json:"id"`
	OwnerID    string                 `json:"owner_id"`
	BaseAsset  string                 `json:"base_asset"`
	QuoteAsset string                 `json:"quote_asset"`
	StartTime  time.Time              `json:"start_time"`
	EndTime    time.Time              `json:"end_time"`
	Price      string                 `json:"price"`
	SoftCap    string                 `json:"soft_cap"`
	HardCap    string                 `json:"hard_cap"`
	CurrentCap string                 `json:"current_cap"`
	Details    map[string]interface{} `json:"details"`
	State      base.Flag              `json:"state"`
	Statistics SaleStatistics         `json:"statistics"`
}

type SaleStatistics struct {
	Investors     int    `json:"investors"`
	AverageAmount string `json:"average_amount"`
}

func (s *Sale) Populate(h *history.Sale) {
	s.PT = strconv.FormatUint(h.ID, 10)
	s.ID = strconv.FormatUint(h.ID, 10)
	s.OwnerID = h.OwnerID
	s.BaseAsset = h.BaseAsset
	s.QuoteAsset = h.QuoteAsset
	s.StartTime = h.StartTime
	s.EndTime = h.EndTime
	s.Price = amount.StringU(h.Price)
	s.SoftCap = amount.StringU(h.SoftCap)
	s.HardCap = amount.StringU(h.HardCap)
	s.CurrentCap = amount.StringU(h.CurrentCap)
	s.Details = h.Details
	s.State.Name = h.State.String()
	s.State.Value = int32(h.State)
}

func (s *Sale) PopulateStatistic(offers []core.Offer) {
	if len(offers) == 0 {
		return
	}

	sum := big.NewInt(0)
	uniqueInvestors := make(map[string]bool)
	for _, offer := range offers {
		sum = sum.Add(sum, big.NewInt(offer.QuoteAmount))
		uniqueInvestors[offer.OwnerID] = true
	}

	quantity := len(uniqueInvestors)
	average := sum.Div(sum, big.NewInt(int64(quantity)))

	averageStr := fmt.Sprintf("%d", math.MaxInt64)
	if average.IsInt64() {
		averageStr = amount.String(average.Int64())
	}

	s.Statistics = SaleStatistics{Investors: quantity, AverageAmount: averageStr}
}

// PagingToken implementation for hal.Pageable
func (s *Sale) PagingToken() string {
	return s.PT
}
