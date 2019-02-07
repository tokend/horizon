package resource

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/horizon/db2/core"
	"gitlab.com/tokend/horizon/fees"
	"gitlab.com/tokend/regources"
)

// Populate fills out the resource's fields
func NewFeeEntry(cfee core.FeeEntry) (fee regources.FeeEntry) {
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
	return fee
}

func NewFeeEntryFromWrapper(wrapper fees.FeeWrapper) (fee regources.FeeEntry) {
	fee = NewFeeEntry(wrapper.FeeEntry)
	fee.Exists = !wrapper.NotExists
	return fee
}
