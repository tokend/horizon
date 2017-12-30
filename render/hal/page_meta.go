package hal

import "time"

type PageMeta struct {
	LatestLedger *LatestLedgerMeta `json:"latest_ledger"`
}

type LatestLedgerMeta struct {
	Sequence int32     `json:"sequence"`
	ClosedAt time.Time `json:"closed_at"`
}
