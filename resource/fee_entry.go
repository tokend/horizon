package resource

import (
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/regources"
)

// Populate fills out the resource's fields
func Populate(fee regources.FeeEntry, cfee core.FeeEntry) regources.FeeEntry {
	fee.FeeType = cfee.FeeType
	fee.Asset = cfee.Asset
	fee.Fixed = amount.String(cfee.Fixed)
	fee.Percent = amount.String(cfee.Percent)
	fee.Subtype = cfee.Subtype
	fee.AccountID = cfee.AccountID
	fee.AccountType = cfee.AccountType
	fee.LowerBound = amount.String(cfee.LowerBound)
	fee.UpperBound = amount.String(cfee.UpperBound)
	fee.FeeAsset = cfee.FeeAsset
	if fee.FeeAsset == "" {
		fee.FeeAsset = cfee.Asset
	}
	fee.Exists = true
	return fee
}
