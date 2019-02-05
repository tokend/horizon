package history2

import (
	"database/sql/driver"
	"time"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2"
	"gitlab.com/tokend/regources/v2"
)

// Sale - represents instance of compounding campaign
type Sale struct {
	ID                uint64       `db:"id"`
	SoftCap           uint64       `db:"soft_cap"`
	HardCap           uint64       `db:"hard_cap"`
	BaseCurrentCap    int64        `db:"base_current_cap"`
	BaseHardCap       int64        `db:"base_hard_cap"`
	SaleType          xdr.SaleType `db:"sale_type"`
	OwnerAddress      string       `db:"owner_address"`
	BaseAsset         string       `db:"base_asset"`
	DefaultQuoteAsset string       `db:"default_quote_asset"`
	StartTime         time.Time    `db:"start_time"`
	EndTime           time.Time    `db:"end_time"`
	CurrentCap        string
	Details           regources.Details `db:"details"`
	QuoteAssets       SaleQuoteAssets   `db:"quote_assets"`
	State             SaleState         `db:"state"`
}

//SaleQuoteAssets - assets allowed to invest in sale
type SaleQuoteAssets struct {
	QuoteAssets []SaleQuoteAsset `json:"quote_assets"`
}

//Value - implements db driver method for auto marshal
func (r SaleQuoteAssets) Value() (driver.Value, error) {
	result, err := db2.DriverValue(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal quote assets")
	}

	return result, nil
}

//Scan - implements db driver method for auto unmarshal
func (r *SaleQuoteAssets) Scan(src interface{}) error {
	err := db2.DriveScan(src, r)
	if err != nil {
		return errors.Wrap(err, "failed to scan quote assets")
	}

	return nil
}

//SaleQuoteAsset - asset allowed to invest into sale
type SaleQuoteAsset struct {
	Asset           string `json:"asset"`
	Price           string `json:"price"`
	QuoteBalanceID  string `json:"quote_balance_id"`
	CurrentCap      string `json:"current_cap"`
	TotalCurrentCap string `json:"total_current_cap,omitempty"`
	HardCap         string `json:"hard_cap,omitempty"`
}
