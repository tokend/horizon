package core2

import "gitlab.com/tokend/go/xdr"

// Account is a row of data from the `accounts` table
type Account struct {
	AccountID    string         `db:"accountid"`
	SequenceID   uint64         `db:"sequence_id"`
	RecoveryID   string         `db:"recoveryid"`
	Thresholds   xdr.Thresholds `db:"thresholds"`
	AccountType  int32          `db:"account_type"`
	BlockReasons int32          `db:"block_reasons"`
	Referrer     string         `db:"referrer"`
	Policies     int32          `db:"policies"`
	KYCLevel     int32          `db:"kyc_level"`
}
