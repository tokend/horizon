package core

import (
	"gitlab.com/swarmfund/go/xdr"
)

// Account is a row of data from the `accounts` table
type Account struct {
	AccountID        string         `db:"accountid"`
	RecoveryID		 string 		`db:"recoveryid"`
	Thresholds       xdr.Thresholds `db:"thresholds"`
	AccountType      int32          `db:"account_type"`
	BlockReasons     int32          `db:"block_reasons"`
	Policies         int32          `db:"policies"`
	*Statistics
}
