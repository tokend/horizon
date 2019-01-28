package regources

import (
	"encoding/json"
	"gitlab.com/tokend/go/xdr"
	"time"
)

type saleState int

const (
	SaleStateOpen saleState = iota
	SaleStateClosed
	SaleStateCanceled
)

var saleStateMap = map[saleState]string{
	SaleStateOpen:     "open",
	SaleStateClosed:   "closed",
	SaleStateCanceled: "canceled",
}

func (s saleState) MarshalJSON() ([]byte, error) {
	type flagValue struct {
		Name  string    `json:"name"`
		Value saleState `json:"value"`
	}

	type flag struct {
		Value saleState   `json:"int"`
		Flags []flagValue `json:"flags"`
	}

	return json.Marshal(flag{
		Value: s,
		Flags: []flagValue{
			{
				Name:  saleStateMap[s],
				Value: s,
			},
		},
	})
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
	SaleState      saleState    `json:"sale_state"`
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