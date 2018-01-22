package resource

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
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

func (s *Sale) PopulateStat(offers []core.Offer, balances []core.Balance, assetPair *core.AssetPair) error {
	if len(offers) == 0 && len(balances) == 0 {
		return nil
	}
	sum := big.NewInt(0)
	uniqueInvestors := make(map[string]bool)
	for _, offer := range offers {
		sum = sum.Add(sum, big.NewInt(offer.QuoteAmount))
		uniqueInvestors[offer.OwnerID] = true
	}

	balanceSum := big.NewInt(0)
	for _, balance := range balances {
		uniqueInvestors[balance.AccountID] = true
		if balance.Amount == 0 {
			continue
		}
		balanceSum = balanceSum.Add(balanceSum, big.NewInt(balance.Amount))
	}

	balanceSumConverted, isConverted, err := assetPair.ConvertToDestAsset(s.QuoteAsset, balanceSum.Int64())
	if err != nil {
		return errors.Wrap(err,
			"failed to convert balance summary",
			logan.F{
				"from":  assetPair.BaseAsset,
				"to":    assetPair.QuoteAsset,
				"price": assetPair.CurrentPrice,
			})
	}

	if !isConverted {
		return errors.Wrap(
			errors.New("failed to convert due to overflow"), "",
			logan.F{
				"from":  assetPair.BaseAsset,
				"to":    assetPair.QuoteAsset,
				"price": assetPair.CurrentPrice,
			})
	}
	sum = sum.Add(sum, big.NewInt(balanceSumConverted))

	quantity := len(uniqueInvestors)
	s.Statistics.Investors = quantity
	s.Statistics.AverageAmount = divToAmountStr(sum, quantity)
	return nil
}

func divToAmountStr(sum *big.Int, quantity int) string {
	result := sum.Div(sum, big.NewInt(int64(quantity)))
	resultStr := fmt.Sprintf("%d", math.MaxInt64)
	if result.IsInt64() {
		resultStr = amount.String(result.Int64())
	}
	return resultStr
}

// PagingToken implementation for hal.Pageable
func (s *Sale) PagingToken() string {
	return s.PT
}
