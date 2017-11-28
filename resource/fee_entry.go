package resource

import (
	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/horizon/db2/core"
)

type FeeEntry struct {
	Asset       string `json:"asset"`
	Fixed       string `json:"fixed"`
	Percent     string `json:"percent"`
	FeeType     int    `json:"fee_type"`
	Subtype     int64  `json:"subtype"`
	AccountID   string `json:"account_id"`
	AccountType int32  `json:"account_type"`
	LowerBound  string `json:"lower_bound"`
	UpperBound  string `json:"upper_bound"`
	Exists      bool   `json:"exists"`
}

// Populate fills out the resource's fields
func (fee *FeeEntry) Populate(cfee core.FeeEntry) error {
	fee.FeeType = cfee.FeeType
	fee.Asset = cfee.Asset
	fee.Fixed = amount.String(cfee.Fixed)
	fee.Percent = amount.String(cfee.Percent)
	fee.Subtype = cfee.Subtype
	fee.AccountID = cfee.AccountID
	fee.AccountType = cfee.AccountType
	fee.LowerBound = amount.String(cfee.LowerBound)
	fee.UpperBound = amount.String(cfee.UpperBound)
	fee.Exists = true
	return nil
}
