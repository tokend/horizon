package regources

import (
	"encoding/json"
	"gitlab.com/tokend/go/xdr"
	"time"
)

type SaleState int

const (
	SaleStateOpen SaleState = 1 << iota
	SaleStateClosed
	SaleStateCanceled
)

var saleStateMap = map[int]string{
	int(SaleStateOpen):     "open",
	int(SaleStateClosed):   "closed",
	int(SaleStateCanceled): "canceled",
}

func (s SaleState) MarshalJSON() ([]byte, error) {
	return json.Marshal(FlagsFromMask(int(s), saleStateMap))
}

//String - converts int enum to string
func (s SaleState) String() string {
	return saleStateMap[int(s)]
}

// SaleResponse - representation of response with sale resource
type SaleResponse struct {
	Data     Sale     `json:"data"`
	Included Included `json:"included"`
}

// SalesResponse - representation of response with sale collection
type SalesResponse struct {
	Links    *Links   `json:"links"`
	Data     []Sale   `json:"data"`
	Included Included `json:"included"`
}

func (r SalesResponse) MarshalJSON() ([]byte, error) {
	if r.Data == nil {
		r.Data = []Sale{}
	}

	type temp SalesResponse
	return json.Marshal(temp(r))
}

// Sale - represents details of the sale
type Sale struct {
	Key
	Attributes    SaleAttrs     `json:"attributes"`
	Relationships SaleRelations `json:"relationships"`
}

// SaleAttrs - attributes of the sale
type SaleAttrs struct {
	StartTime      time.Time    `json:"start_time"`
	EndTime        time.Time    `json:"end_time"`
	SaleType       xdr.SaleType `json:"sale_type"`
	SaleState      SaleState    `json:"sale_state"`
	InvestorsCount uint32       `json:"investors_count"`
	Details        Details      `json:"details"`
}

// SaleRelations - relationships of the sale
type SaleRelations struct {
	Owner             *Relation           `json:"owner"`
	BaseAsset         *Relation           `json:"base_asset"`
	QuoteAssets       *RelationCollection `json:"quote_assets"`
	DefaultQuoteAsset *Relation           `json:"default_quote_asset"`
}
