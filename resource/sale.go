package resource

import (
	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/resource/base"
	"strconv"
	"time"
)

type Sale struct {
	PT         string      `json:"paging_token"`
	ID         string      `json:"id"`
	OwnerID    string      `json:"owner_id"`
	BaseAsset  string      `json:"base_asset"`
	QuoteAsset string      `json:"quote_asset"`
	StartTime  time.Time   `json:"start_time"`
	EndTime    time.Time   `json:"end_time"`
	Price      string      `json:"price"`
	SoftCap    string      `json:"soft_cap"`
	HardCap    string      `json:"hard_cap"`
	CurrentCap string      `json:"current_cap"`
	Details    SaleDetails `json:"details"`
	State      base.Flag   `json:"state"`
}

type SaleDetails struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	ShortDescription string `json:"short_description"`
	Logo             string `json:"logo"`
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
	s.Details.Name = h.Details.Name
	s.Details.Description = h.Details.Description
	s.Details.ShortDescription = h.Details.ShortDescription
	s.Details.Logo = h.Details.Logo
	s.State.Name = h.State.String()
	s.State.Value = int32(h.State)
}

// PagingToken implementation for hal.Pageable
func (s *Sale) PagingToken() string {
	return s.PT
}

