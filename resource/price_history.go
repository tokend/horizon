package resource

import "bullioncoin.githost.io/development/horizon/db2/history"

type PriceHistory struct {
	Prices []history.PricePoint `json:"prices"`
}
