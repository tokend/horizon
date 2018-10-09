package core

import (
	"gitlab.com/tokend/go/xdr"
)

// Account is a row of data from the `accounts` table
type Account struct {
	AccountID    string         `db:"accountid"`
	RecoveryID   string         `db:"recoveryid"`
	Thresholds   xdr.Thresholds `db:"thresholds"`
	AccountType  int32          `db:"account_type"`
	BlockReasons int32          `db:"block_reasons"`
	Referrer     string         `db:"referrer"`
	Policies     int32          `db:"policies"`
	KYCLevel     int32          `db:"kyc_level"`
	StatisticsV2 []StatisticsV2Entry
	*AccountKYC
}
