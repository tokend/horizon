package core

import (
	"time"

	"gitlab.com/swarmfund/go/xdr"
)

// Account is a row of data from the `accounts` table
type Account struct {
	AccountID        string         `db:"accountid"`
	Thresholds       xdr.Thresholds `db:"thresholds"`
	AccountType      int32          `db:"account_type"`
	BlockReasons     int32          `db:"block_reasons"`
	Referrer         string         `db:"referrer"`
	ShareForReferrer xdr.Int64      `db:"share_for_referrer"`
	Policies         int32          `db:"policies"`
	CreatedAt        time.Time      `db:"created_at"`
	*Statistics
}
